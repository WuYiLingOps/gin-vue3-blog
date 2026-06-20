/*
 * 项目名称：blog-frontend
 * 文件名称：subscribe.ts
 * 创建时间：2026-04-26 10:30:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：邮件订阅相关 API 接口定义
 */

import { request } from '@/utils/request'

/**
 * 订阅统计信息接口
 */
export interface SubscribeStats {
  total_count: number // 累积订阅总数（包括已退订）
  active_count: number // 当前活跃订阅者数量
}

/**
 * 推送历史记录接口
 */
export interface PushHistory {
  id: number
  post_id: number
  post_title: string
  total_count: number
  success_count: number
  failed_count: number
  status: number // 0-进行中 1-已完成 2-部分失败
  started_at: string
  completed_at: string
  created_at: string
}

/**
 * 推送详情接口
 */
export interface PushDetail {
  id: number
  push_history_id: number
  subscriber_id: number
  subscriber_email: string
  status: number // 0-待发送 1-成功 2-失败
  error_message: string
  sent_at: string
  created_at: string
}

/**
 * 推送统计信息接口
 */
export interface PushStats {
  total_count: number
  total_success: number
  total_failed: number
  last_push_at: string
}

/**
 * 获取订阅统计信息（公开接口）
 * @returns 返回订阅统计信息
 */
export function getSubscribeStats() {
  return request.get<SubscribeStats>('/subscribe/stats')
}

/**
 * 获取推送历史记录列表（管理员）
 */
export function getPushHistories(page: number, pageSize: number) {
  return request.get<{ list: PushHistory[]; total: number }>('/admin/subscribers/push-histories', {
    params: { page, page_size: pageSize }
  })
}

/**
 * 获取推送历史详情（管理员）
 */
export function getPushHistoryDetail(id: number, page: number, pageSize: number) {
  return request.get<{ history: PushHistory; details: PushDetail[]; total: number }>(
    `/admin/subscribers/push-histories/${id}`,
    {
      params: { page, page_size: pageSize }
    }
  )
}

/**
 * 删除推送历史记录（管理员）
 */
export function deletePushHistory(id: number) {
  return request.delete(`/admin/subscribers/push-histories/${id}`)
}

/**
 * 获取推送统计信息（管理员）
 */
export function getPushStats() {
  return request.get<PushStats>('/admin/subscribers/push-stats')
}
