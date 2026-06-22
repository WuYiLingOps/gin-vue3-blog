package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"blog-backend/service"
	"blog-backend/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type NotificationHandler struct {
	service *service.NotificationService
	hub     *service.NotificationHub
}

func NewNotificationHandler(hub *service.NotificationHub, svc *service.NotificationService) *NotificationHandler {
	return &NotificationHandler{
		service: svc,
		hub:     hub,
	}
}

var notificationUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (h *NotificationHandler) GetNotifications(c *gin.Context) {
	userID := c.GetUint("user_id")
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 50 {
		pageSize = 20
	}

	notifications, total, err := h.service.GetNotifications(userID, page, pageSize)
	if err != nil {
		util.ServerError(c, "获取通知列表失败")
		return
	}

	util.Success(c, gin.H{
		"list":      notifications,
		"total":     total,
		"page":      page,
		"page_size": pageSize,
	})
}

func (h *NotificationHandler) GetUnreadCount(c *gin.Context) {
	userID := c.GetUint("user_id")
	count, err := h.service.GetUnreadCount(userID)
	if err != nil {
		util.ServerError(c, "获取未读数失败")
		return
	}
	util.Success(c, gin.H{"count": count})
}

func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的通知ID")
		return
	}

	if err := h.service.MarkAsRead(uint(id), userID); err != nil {
		util.ServerError(c, "标记已读失败")
		return
	}
	util.Success(c, nil)
}

func (h *NotificationHandler) MarkAllAsRead(c *gin.Context) {
	userID := c.GetUint("user_id")
	if err := h.service.MarkAllAsRead(userID); err != nil {
		util.ServerError(c, "全部已读失败")
		return
	}
	util.Success(c, nil)
}

func (h *NotificationHandler) DeleteNotification(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		util.BadRequest(c, "无效的通知ID")
		return
	}

	if err := h.service.DeleteNotification(uint(id)); err != nil {
		util.ServerError(c, "删除通知失败")
		return
	}
	util.Success(c, nil)
}

func (h *NotificationHandler) HandleWebSocket(c *gin.Context) {
	userID := c.GetUint("user_id")
	if userID == 0 {
		util.Unauthorized(c, "请先登录")
		return
	}

	conn, err := notificationUpgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		util.ServerError(c, "WebSocket升级失败")
		return
	}

	client := &service.NotificationClient{
		Conn:   conn,
		Hub:    h.hub,
		Send:   make(chan []byte, 256),
		UserID: userID,
	}

	h.hub.Register <- client

	go client.NotificationWritePump()
	go client.NotificationReadPump()

	go func() {
		time.Sleep(500 * time.Millisecond)
		count, err := h.service.GetUnreadCount(userID)
		if err == nil {
			msg := service.NotificationWebSocketMessage{
				Type:      "unread_count",
				Data:      map[string]int64{"count": count},
				Timestamp: time.Now().Unix(),
			}
			data, _ := json.Marshal(msg)
			select {
			case client.Send <- data:
			default:
			}
		}
	}()
}
