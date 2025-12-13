// IP黑名单管理相关 API

import { request } from '@/utils/request'
import type { PageData } from '@/types/common'

export interface IPBlacklist {
  id: number
  ip: string
  reason: string
  ban_type: number // 1:自动封禁 2:手动封禁
  expire_at: string | null // 过期时间，null表示永久封禁
  created_at: string
  updated_at: string
}

export interface IPCheckResult {
  banned: boolean
  info?: IPBlacklist
}

// 获取IP黑名单列表
export function getIPBlacklist(params: { page?: number; page_size?: number }) {
  return request.get<PageData<IPBlacklist>>('/admin/ip-blacklist', { params })
}

// 添加IP到黑名单
export function addIPBlacklist(data: {
  ip: string
  reason?: string
  duration?: number // 封禁时长（小时），0表示永久
}) {
  return request.post<IPBlacklist>('/admin/ip-blacklist', data)
}

// 删除IP黑名单
export function deleteIPBlacklist(id: number) {
  return request.delete(`/admin/ip-blacklist/${id}`)
}

// 检查IP是否在黑名单中
export function checkIPStatus(ip: string) {
  return request.get<IPCheckResult>('/admin/ip-blacklist/check', { params: { ip } })
}

// 清理过期的黑名单记录
export function cleanExpiredIPBlacklist() {
  return request.post<{ deleted_count: number }>('/admin/ip-blacklist/clean-expired')
}

