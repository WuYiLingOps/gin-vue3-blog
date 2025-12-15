package handler

import (
	"blog-backend/model"
	"blog-backend/service"
	"blog-backend/util"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有跨域请求，生产环境应该设置具体的域名
	},
}

// ChatHandler 聊天处理器
type ChatHandler struct {
	service *service.ChatService
	hub     *service.Hub
}

// NewChatHandler 创建聊天处理器
func NewChatHandler(hub *service.Hub) *ChatHandler {
	return &ChatHandler{
		service: service.NewChatService(hub),
		hub:     hub,
	}
}

// HandleWebSocket 处理WebSocket连接
func (h *ChatHandler) HandleWebSocket(c *gin.Context) {
	// 升级HTTP连接为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket升级失败: %v, IP: %s", err, util.GetClientIP(c))
		util.ServerError(c, "WebSocket升级失败")
		return
	}

	log.Printf("WebSocket连接成功, IP: %s", util.GetClientIP(c))

	// 获取用户信息
	var userID *uint
	var username string
	var avatar string

	// 尝试从JWT获取用户信息（注意：键名是 user_id 不是 userID）
	if userIDInterface, exists := c.Get("user_id"); exists {
		if uid, ok := userIDInterface.(uint); ok {
			userID = &uid

			// 从JWT获取用户信息
			if usernameInterface, exists := c.Get("username"); exists {
				if uname, ok := usernameInterface.(string); ok {
					username = uname
				}
			}

			// 从JWT获取头像（如果有的话，也可以从数据库查询）
			// 这里简化处理
		}
	}

	// 如果是匿名用户，从查询参数获取昵称
	if username == "" {
		username = c.Query("username")
		if username == "" {
			username = "访客" + strconv.FormatInt(time.Now().Unix()%10000, 10)
		}
	}

	// 获取头像
	avatar = c.Query("avatar")
	if avatar == "" {
		// 使用默认头像
		avatar = ""
	}

	// 获取IP地址
	ip := util.GetClientIP(c)

	// 创建客户端
	client := &service.Client{
		ID:       uuid.New().String(),
		Conn:     conn,
		Hub:      h.hub,
		Send:     make(chan []byte, 256),
		UserID:   userID,
		Username: username,
		Avatar:   avatar,
		IP:       ip,
	}

	// 注册客户端
	h.hub.Register <- client

	// 启动读写goroutine
	go client.WritePump()
	go client.ReadPump()
}

// GetMessages 获取消息列表
func (h *ChatHandler) GetMessages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "50"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 50
	}

	messages, total, err := h.service.GetMessages(page, pageSize, false)
	if err != nil {
		util.ServerError(c, "获取消息列表失败")
		return
	}

	util.Success(c, gin.H{
		"list":      messages,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// DeleteMessage 删除消息
func (h *ChatHandler) DeleteMessage(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		util.BadRequest(c, "无效的消息ID")
		return
	}

	if err := h.service.DeleteMessage(uint(id)); err != nil {
		util.ServerError(c, "删除消息失败")
		return
	}

	util.Success(c, nil)
}

// GetOnlineInfo 获取在线信息
func (h *ChatHandler) GetOnlineInfo(c *gin.Context) {
	util.Success(c, gin.H{
		"online_count": h.service.GetOnlineCount(),
		"online_users": h.service.GetOnlineUsers(),
	})
}

// GetStats 获取聊天室统计信息
func (h *ChatHandler) GetStats(c *gin.Context) {
	onlineCount := h.service.GetOnlineCount()

	// 这里可以添加更多统计信息
	util.Success(c, gin.H{
		"online_count": onlineCount,
	})
}

// AdminListMessages 管理员获取消息列表
func (h *ChatHandler) AdminListMessages(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	messages, total, err := h.service.GetMessages(page, pageSize, true)
	if err != nil {
		util.ServerError(c, "获取消息列表失败")
		return
	}

	util.Success(c, gin.H{
		"list":      messages,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

// CreateAnonymousToken 创建匿名用户令牌（可选功能）
func (h *ChatHandler) CreateAnonymousToken(c *gin.Context) {
	var req struct {
		Username string `json:"username" binding:"required,min=2,max=20"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "用户名格式错误")
		return
	}

	// 生成临时令牌
	token := uuid.New().String()

	util.Success(c, gin.H{
		"token":    token,
		"username": req.Username,
	})
}

// BroadcastSystemMessage 发送系统广播消息（管理员功能）
func (h *ChatHandler) BroadcastSystemMessage(c *gin.Context) {
	var req struct {
		Content  string `json:"content" binding:"required"`
		Priority *int   `json:"priority"`         // 0:普通 1:置顶
		Target   string `json:"target,omitempty"` // announcement / chat / both
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "消息内容不能为空")
		return
	}

	priority := 0
	if req.Priority != nil && *req.Priority == 1 {
		priority = 1
	}

	target := "both"
	switch strings.ToLower(strings.TrimSpace(req.Target)) {
	case "announcement":
		target = "announcement"
	case "chat":
		target = "chat"
	case "both", "":
		target = "both"
	default:
		util.BadRequest(c, "target 参数不合法，可选 announcement/chat/both")
		return
	}

	// 创建系统消息
	message := &model.ChatMessage{
		Content:     req.Content,
		Username:    "系统消息",
		Avatar:      "",
		Priority:    priority,
		Target:      target,
		IsBroadcast: true,
		Status:      1,
	}

	// 保存到数据库
	if err := h.hub.Repo.Create(message); err != nil {
		util.ServerError(c, "发送消息失败")
		return
	}

	// 如选择投递到聊天室，才走 WebSocket 广播
	if target == "chat" || target == "both" {
		wsMsg := service.WebSocketMessage{
			Type:      "system",
			Data:      message,
			Timestamp: time.Now().Unix(),
		}

		data, _ := json.Marshal(wsMsg)
		h.hub.Broadcast <- data
	}

	util.Success(c, nil)
}

// KickUser 踢出用户（管理员功能）
func (h *ChatHandler) KickUser(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Reason   string `json:"reason"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	if req.Reason == "" {
		req.Reason = "违反聊天室规则"
	}

	// 踢出用户
	if h.hub.KickClient(req.ClientID, req.Reason) {
		util.Success(c, nil)
	} else {
		util.Error(c, 404, "用户不在线")
	}
}

// BanIP 封禁IP（管理员功能）
func (h *ChatHandler) BanIP(c *gin.Context) {
	var req struct {
		ClientID string `json:"client_id" binding:"required"`
		Reason   string `json:"reason"`
		Duration int    `json:"duration"` // 封禁时长（小时），0表示永久
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		util.BadRequest(c, "参数错误")
		return
	}

	// 获取客户端信息
	client := h.hub.GetClientByID(req.ClientID)
	if client == nil {
		util.Error(c, 404, "用户不在线")
		return
	}

	// 添加到IP黑名单
	// 这里需要调用IP黑名单服务
	// TODO: 实现IP封禁逻辑

	// 先踢出用户
	h.hub.KickClient(req.ClientID, "您已被封禁："+req.Reason)

	util.Success(c, gin.H{
		"ip": client.IP,
	})
}
