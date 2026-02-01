/*
 * 项目名称：blog-frontend
 * 文件名称：chat.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：聊天室相关 API 接口定义，包括消息管理、在线用户管理、系统广播、用户封禁、聊天室配置等功能。
 */

import request from '@/utils/request'
import type { PaginationParams, PaginationResult } from '@/types/common'

/**
 * 聊天消息接口
 */
export interface ChatMessage {
  id: number                                    // 消息ID
  content: string                               // 消息内容
  user_id?: number                              // 用户ID（可选）
  username: string                              // 用户名
  avatar?: string                               // 用户头像（可选）
  ip?: string                                   // IP地址（可选）
  client_id?: string                            // WebSocket客户端ID，用于踢出用户
  priority?: number                             // 优先级
  is_broadcast?: boolean                        // 是否为广播消息
  target?: 'announcement' | 'chat' | 'both'     // 目标位置：公告栏、聊天室或两者
  status: number                                // 消息状态
  created_at: string                            // 创建时间
  updated_at: string                            // 更新时间
}

/**
 * 在线用户信息接口
 */
export interface OnlineUser {
  id: string            // 用户ID（WebSocket客户端ID）
  username: string      // 用户名
  avatar?: string       // 用户头像（可选）
}

/**
 * 在线信息接口
 */
export interface OnlineInfo {
  online_count: number        // 在线用户数量
  online_users: OnlineUser[]  // 在线用户列表
}

/**
 * WebSocket消息接口
 */
export interface WebSocketMessage {
  type: 'message' | 'history' | 'user_join' | 'user_leave' | 'user_list' | 'system'  // 消息类型
  data: any                                                                          // 消息数据
  timestamp: number                                                                  // 时间戳
}

/**
 * 获取消息列表（公开接口）
 * @param params 分页参数
 * @returns 返回分页的消息列表
 */
export function getChatMessages(params: PaginationParams) {
  return request.get<PaginationResult<ChatMessage>>('/chat/messages', { params })
}

/**
 * 获取在线信息
 * @returns 返回在线用户数量和列表
 */
export function getOnlineInfo() {
  return request.get<OnlineInfo>('/chat/online')
}

/**
 * 管理员：获取消息列表
 * @param params 分页参数
 * @returns 返回分页的消息列表
 */
export function adminGetMessages(params: PaginationParams) {
  return request.get<PaginationResult<ChatMessage>>('/admin/chat/messages', { params })
}

/**
 * 管理员：删除消息
 * @param id 消息ID
 * @returns 返回删除结果
 */
export function adminDeleteMessage(id: number) {
  return request.delete(`/admin/chat/messages/${id}`)
}

/**
 * 管理员：发送系统广播
 * @param content 广播内容
 * @param priority 优先级，默认为0
 * @param target 目标位置，默认为'both'（公告栏和聊天室）
 * @returns 返回发送结果
 */
export function adminBroadcastMessage(content: string, priority = 0, target: 'announcement' | 'chat' | 'both' = 'both') {
  return request.post('/admin/chat/broadcast', { content, priority, target })
}

/**
 * 管理员：踢出用户
 * @param client_id WebSocket客户端ID
 * @param reason 踢出原因（可选）
 * @returns 返回踢出结果
 */
export function adminKickUser(client_id: string, reason?: string) {
  return request.post('/admin/chat/kick', { client_id, reason })
}

/**
 * 管理员：封禁IP
 * @param client_id WebSocket客户端ID
 * @param reason 封禁原因（可选）
 * @param duration 封禁时长（小时，可选）
 * @returns 返回封禁结果
 */
export function adminBanIP(client_id: string, reason?: string, duration?: number) {
  return request.post('/admin/chat/ban', { client_id, reason, duration })
}

/**
 * 聊天室配置接口
 */
export interface ChatSettings {
  chat_mute_all: string  // 是否全员禁言，'0'表示否，'1'表示是
}

/**
 * 公开获取聊天室配置（前端展示禁言状态）
 * @returns 返回聊天室配置
 */
export function getChatSettings() {
  return request.get<ChatSettings>('/chat/settings')
}

/**
 * 管理员更新聊天室配置
 * @param data 聊天室配置数据
 * @returns 返回更新结果
 */
export function updateChatSettings(data: ChatSettings) {
  return request.put('/admin/chat/settings', data)
}