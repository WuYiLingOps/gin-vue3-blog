// WebSocket管理工具

import type { WebSocketMessage } from '@/api/chat'

export type WebSocketEventCallback = (data: any) => void

export class ChatWebSocket {
  private ws: WebSocket | null = null
  private url: string
  private reconnectTimer: number | null = null
  private heartbeatTimer: number | null = null
  private reconnectAttempts = 0
  private maxReconnectAttempts = 5
  private reconnectDelay = 3000
  private eventHandlers: Map<string, WebSocketEventCallback[]> = new Map()
  private isManualClose = false

  constructor(url: string) {
    this.url = url
  }

  // 连接WebSocket
  connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        this.ws = new WebSocket(this.url)
        this.isManualClose = false

        this.ws.onopen = () => {
          console.log('WebSocket连接成功')
          this.reconnectAttempts = 0
          this.startHeartbeat()
          this.emit('open', null)
          resolve()
        }

        this.ws.onmessage = (event) => {
          try {
            const messages = event.data.split('\n')
            messages.forEach((msgStr: string) => {
              if (msgStr.trim()) {
                const message: WebSocketMessage = JSON.parse(msgStr)
                this.handleMessage(message)
              }
            })
          } catch (error) {
            console.error('解析消息失败:', error)
          }
        }

        this.ws.onerror = (error) => {
          console.error('WebSocket错误:', error)
          this.emit('error', error)
          reject(error)
        }

        this.ws.onclose = () => {
          console.log('WebSocket连接关闭')
          this.stopHeartbeat()
          this.emit('close', null)

          // 如果不是手动关闭，尝试重连
          if (!this.isManualClose && this.reconnectAttempts < this.maxReconnectAttempts) {
            this.scheduleReconnect()
          }
        }
      } catch (error) {
        reject(error)
      }
    })
  }

  // 发送消息
  send(type: string, data: any) {
    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
      const message = JSON.stringify({ type, ...data })
      this.ws.send(message)
    } else {
      console.error('WebSocket未连接')
    }
  }

  // 发送聊天消息
  sendMessage(content: string) {
    this.send('message', { content })
  }

  // 处理接收到的消息
  private handleMessage(message: WebSocketMessage) {
    const { type, data } = message
    // 触发特定类型的事件
    this.emit(type, data)
    // 触发通用事件（使用不同的事件名，避免和type='message'冲突）
    this.emit('ws:message', message)
  }

  // 注册事件监听
  on(event: string, callback: WebSocketEventCallback) {
    if (!this.eventHandlers.has(event)) {
      this.eventHandlers.set(event, [])
    }
    this.eventHandlers.get(event)?.push(callback)
  }

  // 移除事件监听
  off(event: string, callback?: WebSocketEventCallback) {
    if (!callback) {
      this.eventHandlers.delete(event)
    } else {
      const handlers = this.eventHandlers.get(event)
      if (handlers) {
        const index = handlers.indexOf(callback)
        if (index > -1) {
          handlers.splice(index, 1)
        }
      }
    }
  }

  // 触发事件
  private emit(event: string, data: any) {
    const handlers = this.eventHandlers.get(event)
    if (handlers) {
      handlers.forEach(callback => callback(data))
    }
  }

  // 心跳检测
  private startHeartbeat() {
    this.heartbeatTimer = window.setInterval(() => {
      if (this.ws && this.ws.readyState === WebSocket.OPEN) {
        // 发送ping消息（可选）
      }
    }, 30000)
  }

  private stopHeartbeat() {
    if (this.heartbeatTimer) {
      clearInterval(this.heartbeatTimer)
      this.heartbeatTimer = null
    }
  }

  // 重连机制
  private scheduleReconnect() {
    if (this.reconnectTimer) {
      return
    }

    this.reconnectAttempts++
    console.log(`尝试重连 (${this.reconnectAttempts}/${this.maxReconnectAttempts})...`)

    this.reconnectTimer = window.setTimeout(() => {
      this.reconnectTimer = null
      this.connect().catch(() => {
        // 重连失败，会在 onclose 中继续尝试
      })
    }, this.reconnectDelay)
  }

  // 关闭连接
  close() {
    this.isManualClose = true
    if (this.reconnectTimer) {
      clearTimeout(this.reconnectTimer)
      this.reconnectTimer = null
    }
    this.stopHeartbeat()
    if (this.ws) {
      this.ws.close()
      this.ws = null
    }
  }

  // 获取连接状态
  isConnected(): boolean {
    return this.ws !== null && this.ws.readyState === WebSocket.OPEN
  }
}

// 创建聊天WebSocket连接
export function createChatWebSocket(username?: string, avatar?: string, token?: string): ChatWebSocket {
  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const host = import.meta.env.VITE_API_BASE_URL?.replace(/^https?:\/\//, '') || window.location.host
  
  // 构建WebSocket URL
  let wsUrl = `${protocol}//${host}/api/chat/ws`
  
  // 添加查询参数
  const params = new URLSearchParams()
  if (username) {
    params.append('username', username)
  }
  if (avatar) {
    params.append('avatar', avatar)
  }
  // 添加token（如果有）
  if (token) {
    params.append('token', token)
  }
  
  if (params.toString()) {
    wsUrl += `?${params.toString()}`
  }

  return new ChatWebSocket(wsUrl)
}

