/*
 * 项目名称：blog-frontend
 * 文件名称：operationLog.ts
 * 创建时间：2026-02-06 22:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：操作日志相关 API 接口定义，包括操作日志列表查询、详情查询等功能（仅超级管理员）。
 */

import { request } from '@/utils/request'
import type { PageData } from '@/types/common'

/**
 * 操作日志接口类型定义
 */
export interface OperationLog {
  id: number
  user_id: number
  username: string
  action: string
  module: string
  target_type: string
  target_id: number | null
  target_name: string
  description: string
  ip: string
  user_agent: string
  created_at: string
  user?: {
    id: number
    username: string
    nickname: string
    avatar: string
  }
}

/**
 * 操作日志查询参数
 */
export interface OperationLogParams {
  page?: number
  page_size?: number
  module?: string
  action?: string
  username?: string
}

/**
 * 获取操作日志列表（仅超级管理员）
 * @param params 查询参数（分页、模块、操作类型、用户名等）
 * @returns 返回分页的操作日志列表
 */
export function getOperationLogs(params?: OperationLogParams) {
  return request.get<PageData<OperationLog>>('/admin/operation-logs', { params })
}

/**
 * 获取操作日志详情（仅超级管理员）
 * @param id 操作日志ID
 * @returns 返回操作日志详细信息
 */
export function getOperationLog(id: number) {
  return request.get<OperationLog>(`/admin/operation-logs/${id}`)
}

/**
 * 删除操作日志（仅超级管理员）
 * @param id 操作日志ID
 * @returns 返回删除结果
 */
export function deleteOperationLog(id: number) {
  return request.delete(`/admin/operation-logs/${id}`)
}

/**
 * 批量删除操作日志（仅超级管理员）
 * @param ids 操作日志ID数组
 * @returns 返回删除结果
 */
export function batchDeleteOperationLogs(ids: number[]) {
  return request.post('/admin/operation-logs/batch-delete', { ids })
}
