// 聊天室相关 API

import request from '@/utils/request'
import type { PaginationParams, PaginationResult } from '@/types/common'

// 聊天消息接口
export interface ChatMessage {
  id: number
  content: string
  user_id?: number
  username: string
  avatar?: string
  ip?: string
  client_id?: string  // WebSocket客户端ID，用于踢出用户
  priority?: number
  is_broadcast?: boolean
  target?: 'announcement' | 'chat' | 'both'
  status: number
  created_at: string
  updated_at: string
}

// 在线用户信息
export interface OnlineUser {
  id: string
  username: string
  avatar?: string
}

// 在线信息
export interface OnlineInfo {
  online_count: number
  online_users: OnlineUser[]
}

// WebSocket消息
export interface WebSocketMessage {
  type: 'message' | 'history' | 'user_join' | 'user_leave' | 'user_list' | 'system'
  data: any
  timestamp: number
}

// 获取消息列表
export function getChatMessages(params: PaginationParams) {
  return request.get<PaginationResult<ChatMessage>>('/chat/messages', { params })
}

// 获取在线信息
export function getOnlineInfo() {
  return request.get<OnlineInfo>('/chat/online')
}

// 管理员：获取消息列表
export function adminGetMessages(params: PaginationParams) {
  return request.get<PaginationResult<ChatMessage>>('/admin/chat/messages', { params })
}

// 管理员：删除消息
export function adminDeleteMessage(id: number) {
  return request.delete(`/admin/chat/messages/${id}`)
}

// 管理员：发送系统广播
export function adminBroadcastMessage(content: string, priority = 0, target: 'announcement' | 'chat' | 'both' = 'both') {
  return request.post('/admin/chat/broadcast', { content, priority, target })
}

// 管理员：踢出用户
export function adminKickUser(client_id: string, reason?: string) {
  return request.post('/admin/chat/kick', { client_id, reason })
}

// 管理员：封禁IP
export function adminBanIP(client_id: string, reason?: string, duration?: number) {
  return request.post('/admin/chat/ban', { client_id, reason, duration })
}

// 聊天室配置
export interface ChatSettings {
  chat_mute_all: string // 0/1
}

// 公开获取聊天室配置（前端展示禁言状态）
export function getChatSettings() {
  return request.get<ChatSettings>('/chat/settings')
}

// 管理员更新聊天室配置
export function updateChatSettings(data: ChatSettings) {
  return request.put('/admin/chat/settings', data)
}