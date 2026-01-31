/*
 * 项目名称：blog-backend
 * 文件名称：chat.go
 * 创建时间：2026-01-31 16:34:35
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：聊天室业务逻辑层，提供WebSocket实时聊天、消息管理、在线用户管理等功能
 */
package service

import (
	"blog-backend/model"
	"blog-backend/repository"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Client WebSocket客户端
type Client struct {
	ID       string          // 客户端唯一标识
	Conn     *websocket.Conn // WebSocket连接
	Hub      *Hub            // 所属Hub
	Send     chan []byte     // 发送消息通道
	UserID   *uint           // 登录用户ID，可为空
	Username string          // 用户名
	Avatar   string          // 头像
	IP       string          // IP地址
	Role     string          // 角色：admin/user/guest
}

// Hub WebSocket Hub，管理所有客户端
type Hub struct {
	Clients     map[*Client]bool // 注册的客户端
	Broadcast   chan []byte      // 广播消息通道
	Register    chan *Client     // 注册客户端通道
	Unregister  chan *Client     // 注销客户端通道
	mutex       sync.RWMutex     // 读写锁
	Repo        *repository.ChatRepository
	SettingRepo *repository.SettingRepository
}

// WebSocketMessage WebSocket消息结构
type WebSocketMessage struct {
	Type      string      `json:"type"`      // 消息类型：message, user_join, user_leave, user_list, history
	Data      interface{} `json:"data"`      // 消息内容
	Timestamp int64       `json:"timestamp"` // 时间戳
}

// UserInfo 用户信息
type UserInfo struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Avatar   string `json:"avatar"`
}

// NewHub 创建新的Hub
func NewHub() *Hub {
	return &Hub{
		Clients:     make(map[*Client]bool),
		Broadcast:   make(chan []byte, 256),
		Register:    make(chan *Client),
		Unregister:  make(chan *Client),
		Repo:        repository.NewChatRepository(),
		SettingRepo: repository.NewSettingRepository(),
	}
}

// Run 启动Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			h.Clients[client] = true
			h.mutex.Unlock()

			// 发送历史消息给新连接的客户端
			go h.sendHistory(client)

			// 广播用户加入消息
			h.broadcastUserJoin(client)

			// 广播最新的在线用户列表给所有人
			go h.broadcastUserList()

		case client := <-h.Unregister:
			h.mutex.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				close(client.Send)

				// 广播用户离开消息
				h.broadcastUserLeave(client)
			}
			h.mutex.Unlock()

			// 广播最新的在线用户列表给所有人
			go h.broadcastUserList()

		case message := <-h.Broadcast:
			h.mutex.RLock()
			for client := range h.Clients {
				select {
				case client.Send <- message:
				default:
					close(client.Send)
					delete(h.Clients, client)
				}
			}
			h.mutex.RUnlock()
		}
	}
}

// sendHistory 发送历史消息
func (h *Hub) sendHistory(client *Client) {
	messages, err := h.Repo.GetRecentMessages(50)
	if err != nil {
		log.Printf("获取历史消息失败: %v", err)
		return
	}

	// 为历史消息添加client_id字段（设为nil，因为用户可能已离线）
	messagesWithClientID := make([]map[string]interface{}, len(messages))
	for i, msg := range messages {
		messagesWithClientID[i] = map[string]interface{}{
			"id":           msg.ID,
			"content":      msg.Content,
			"user_id":      msg.UserID,
			"username":     msg.Username,
			"avatar":       msg.Avatar,
			"client_id":    nil, // 历史消息用户可能已离线
			"is_broadcast": msg.IsBroadcast,
			"target":       msg.Target,
			"status":       msg.Status,
			"created_at":   msg.CreatedAt,
			"updated_at":   msg.UpdatedAt,
		}
	}

	wsMsg := WebSocketMessage{
		Type:      "history",
		Data:      messagesWithClientID,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("序列化历史消息失败: %v", err)
		return
	}

	select {
	case client.Send <- data:
	default:
	}
}

// broadcastUserJoin 广播用户加入
func (h *Hub) broadcastUserJoin(client *Client) {
	wsMsg := WebSocketMessage{
		Type: "user_join",
		Data: UserInfo{
			ID:       client.ID,
			Username: client.Username,
			Avatar:   client.Avatar,
		},
		Timestamp: time.Now().Unix(),
	}

	data, _ := json.Marshal(wsMsg)
	h.Broadcast <- data
}

// broadcastUserLeave 广播用户离开
func (h *Hub) broadcastUserLeave(client *Client) {
	wsMsg := WebSocketMessage{
		Type: "user_leave",
		Data: UserInfo{
			ID:       client.ID,
			Username: client.Username,
		},
		Timestamp: time.Now().Unix(),
	}

	data, _ := json.Marshal(wsMsg)
	h.Broadcast <- data
}

// sendUserList 发送在线用户列表给单个客户端（去重后）
func (h *Hub) sendUserList(client *Client) {
	// 使用 GetOnlineUsers 获取去重后的用户列表
	userList := h.GetOnlineUsers()

	wsMsg := WebSocketMessage{
		Type:      "user_list",
		Data:      userList,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		return
	}

	select {
	case client.Send <- data:
	default:
	}
}

// broadcastUserList 广播在线用户列表给所有客户端（去重后）
func (h *Hub) broadcastUserList() {
	// 使用 GetOnlineUsers 获取去重后的用户列表
	userList := h.GetOnlineUsers()

	wsMsg := WebSocketMessage{
		Type:      "user_list",
		Data:      userList,
		Timestamp: time.Now().Unix(),
	}

	data, _ := json.Marshal(wsMsg)
	h.Broadcast <- data
}

// GetOnlineCount 获取在线人数（按用户去重）
func (h *Hub) GetOnlineCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// 使用 map 去重，key 为 user_id（登录用户）或 username（匿名用户）
	uniqueUsers := make(map[string]bool)

	for client := range h.Clients {
		var key string
		if client.UserID != nil {
			// 登录用户，使用 user_id 作为唯一标识
			key = fmt.Sprintf("user_%d", *client.UserID)
		} else {
			// 匿名用户，使用 username 作为唯一标识
			key = fmt.Sprintf("anonymous_%s", client.Username)
		}
		uniqueUsers[key] = true
	}

	return len(uniqueUsers)
}

// GetOnlineUsers 获取在线用户列表（按用户去重）
func (h *Hub) GetOnlineUsers() []UserInfo {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	// 使用 map 去重
	uniqueUsersMap := make(map[string]UserInfo)

	for client := range h.Clients {
		var key string
		if client.UserID != nil {
			// 登录用户，使用 user_id 作为唯一标识
			key = fmt.Sprintf("user_%d", *client.UserID)
		} else {
			// 匿名用户，使用 username 作为唯一标识
			key = fmt.Sprintf("anonymous_%s", client.Username)
		}

		// 如果已存在，保留第一个连接的信息（或者可以更新为最新的）
		if _, exists := uniqueUsersMap[key]; !exists {
			uniqueUsersMap[key] = UserInfo{
				ID:       client.ID,
				Username: client.Username,
				Avatar:   client.Avatar,
			}
		}
	}

	// 转换为数组
	users := make([]UserInfo, 0, len(uniqueUsersMap))
	for _, user := range uniqueUsersMap {
		users = append(users, user)
	}

	return users
}

// IsChatMuted 是否开启全员禁言（默认关闭）
func (h *Hub) IsChatMuted() bool {
	setting, err := h.SettingRepo.GetByKey("chat_mute_all")
	if err != nil || setting == nil {
		return false
	}
	return setting.Value == "1"
}

// KickClient 踢出客户端
func (h *Hub) KickClient(clientID string, reason string) bool {
	h.mutex.Lock()
	defer h.mutex.Unlock()

	for client := range h.Clients {
		if client.ID == clientID {
			// 发送被踢出消息
			wsMsg := WebSocketMessage{
				Type: "kick",
				Data: map[string]interface{}{
					"reason": reason,
				},
				Timestamp: time.Now().Unix(),
			}
			data, _ := json.Marshal(wsMsg)

			select {
			case client.Send <- data:
			default:
			}

			// 关闭连接
			time.AfterFunc(500*time.Millisecond, func() {
				client.Conn.Close()
			})

			return true
		}
	}
	return false
}

// GetClientByID 根据ID获取客户端
func (h *Hub) GetClientByID(clientID string) *Client {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	for client := range h.Clients {
		if client.ID == clientID {
			return client
		}
	}
	return nil
}

// ReadPump 从客户端读取消息
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.Unregister <- c
		c.Conn.Close()
	}()

	c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.Conn.SetPongHandler(func(string) error {
		c.Conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket错误: %v", err)
			}
			break
		}

		// 解析消息
		var msg map[string]interface{}
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("解析消息失败: %v", err)
			continue
		}

		// 处理不同类型的消息
		msgType, _ := msg["type"].(string)
		switch msgType {
		case "message":
			content, _ := msg["content"].(string)
			if content == "" {
				continue
			}

			// 全员禁言校验：仅管理员可发言
			if c.Role != "admin" && c.Hub.IsChatMuted() {
				wsMsg := WebSocketMessage{
					Type: "system",
					Data: map[string]interface{}{
						"message": "当前已开启全员禁言，只有管理员可发言",
					},
					Timestamp: time.Now().Unix(),
				}
				if data, err := json.Marshal(wsMsg); err == nil {
					select {
					case c.Send <- data:
					default:
					}
				}
				continue
			}

			// 保存消息到数据库
			// 确保IP地址不为空
			ip := c.IP
			if ip == "" {
				ip = "unknown"
				log.Printf("警告: 保存消息时客户端IP为空，使用默认值: %s", ip)
			}

			chatMsg := &model.ChatMessage{
				Content:  content,
				UserID:   c.UserID,
				Username: c.Username,
				Avatar:   c.Avatar,
				IP:       ip,
				Status:   1,
			}

			if err := c.Hub.Repo.Create(chatMsg); err != nil {
				log.Printf("保存消息失败: %v", err)
				continue
			}

			// 构建包含client_id的消息响应
			messageData := map[string]interface{}{
				"id":         chatMsg.ID,
				"content":    chatMsg.Content,
				"user_id":    chatMsg.UserID,
				"username":   chatMsg.Username,
				"avatar":     chatMsg.Avatar,
				"client_id":  c.ID, // 添加client_id用于管理员踢出功能
				"status":     chatMsg.Status,
				"created_at": chatMsg.CreatedAt,
				"updated_at": chatMsg.UpdatedAt,
			}

			// 广播消息
			wsMsg := WebSocketMessage{
				Type:      "message",
				Data:      messageData,
				Timestamp: time.Now().Unix(),
			}

			data, _ := json.Marshal(wsMsg)
			c.Hub.Broadcast <- data
		}
	}
}

// WritePump 向客户端写入消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 批量发送排队的消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ChatService 聊天服务
type ChatService struct {
	repo *repository.ChatRepository
	hub  *Hub
}

// NewChatService 创建聊天服务
func NewChatService(hub *Hub) *ChatService {
	return &ChatService{
		repo: repository.NewChatRepository(),
		hub:  hub,
	}
}

// GetMessages 获取消息列表（分页）
func (s *ChatService) GetMessages(page, pageSize int, includeAnnouncementOnly bool) ([]model.ChatMessage, int64, error) {
	return s.repo.GetMessages(page, pageSize, includeAnnouncementOnly)
}

// DeleteMessage 删除消息
func (s *ChatService) DeleteMessage(id uint) error {
	return s.repo.Delete(id)
}

// GetOnlineCount 获取在线人数
func (s *ChatService) GetOnlineCount() int {
	return s.hub.GetOnlineCount()
}

// GetOnlineUsers 获取在线用户列表
func (s *ChatService) GetOnlineUsers() []UserInfo {
	return s.hub.GetOnlineUsers()
}
