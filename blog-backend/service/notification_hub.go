/*
 * 项目名称：blog-backend
 * 文件名称：notification_hub.go
 * 创建时间：2026-06-21 10:30:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：通知WebSocket Hub，提供定向推送给特定用户的实时通知功能
 */
package service

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// NotificationClient 通知WebSocket客户端
type NotificationClient struct {
	ID     string          // 客户端唯一标识
	Conn   *websocket.Conn // WebSocket连接
	Hub    *NotificationHub // 所属Hub
	Send   chan []byte     // 发送消息通道
	UserID uint            // 用户ID（必须登录）
}

// NotificationHub 通知WebSocket Hub，管理所有通知客户端
type NotificationHub struct {
	Clients    map[uint]map[*NotificationClient]bool // userID -> 客户端集合
	Register   chan *NotificationClient               // 注册客户端通道
	Unregister chan *NotificationClient               // 注销客户端通道
	mutex      sync.RWMutex                          // 读写锁
}

// NotificationWebSocketMessage 通知WebSocket消息结构
type NotificationWebSocketMessage struct {
	Type      string      `json:"type"`      // 消息类型：notification
	Data      interface{} `json:"data"`      // 消息内容
	Timestamp int64       `json:"timestamp"` // 时间戳
}

// NewNotificationHub 创建新的通知Hub
func NewNotificationHub() *NotificationHub {
	return &NotificationHub{
		Clients:    make(map[uint]map[*NotificationClient]bool),
		Register:   make(chan *NotificationClient),
		Unregister: make(chan *NotificationClient),
	}
}

// Run 启动通知Hub
func (h *NotificationHub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.mutex.Lock()
			if _, exists := h.Clients[client.UserID]; !exists {
				h.Clients[client.UserID] = make(map[*NotificationClient]bool)
			}
			h.Clients[client.UserID][client] = true
			h.mutex.Unlock()

			log.Printf("用户 %d 的通知客户端已连接", client.UserID)

		case client := <-h.Unregister:
			h.mutex.Lock()
			if clients, exists := h.Clients[client.UserID]; exists {
				if _, ok := clients[client]; ok {
					delete(clients, client)
					close(client.Send)
					if len(clients) == 0 {
						delete(h.Clients, client.UserID)
					}
				}
			}
			h.mutex.Unlock()

			log.Printf("用户 %d 的通知客户端已断开", client.UserID)
		}
	}
}

// BroadcastToUser 向指定用户的所有客户端广播消息
func (h *NotificationHub) BroadcastToUser(userID uint, message []byte) {
	h.mutex.RLock()
	defer h.mutex.RUnlock()

	if clients, exists := h.Clients[userID]; exists {
		for client := range clients {
			select {
			case client.Send <- message:
			default:
				// 发送失败，关闭连接
				close(client.Send)
				delete(clients, client)
			}
		}
	}
}

// GetOnlineUserCount 获取在线用户数量
func (h *NotificationHub) GetOnlineUserCount() int {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	return len(h.Clients)
}

// IsUserOnline 检查用户是否在线
func (h *NotificationHub) IsUserOnline(userID uint) bool {
	h.mutex.RLock()
	defer h.mutex.RUnlock()
	clients, exists := h.Clients[userID]
	return exists && len(clients) > 0
}

// SendNotificationToUser 发送通知给指定用户
func (h *NotificationHub) SendNotificationToUser(userID uint, notification interface{}) {
	wsMsg := NotificationWebSocketMessage{
		Type:      "notification",
		Data:      notification,
		Timestamp: time.Now().Unix(),
	}

	data, err := json.Marshal(wsMsg)
	if err != nil {
		log.Printf("序列化通知消息失败: %v", err)
		return
	}

	h.BroadcastToUser(userID, data)
}

// NotificationReadPump 从客户端读取消息
func (c *NotificationClient) NotificationReadPump() {
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
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("通知WebSocket错误: %v", err)
			}
			break
		}
		// 通知客户端不需要处理接收到的消息
	}
}

// NotificationWritePump 向客户端写入消息
func (c *NotificationClient) NotificationWritePump() {
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
