/*
 * @Author: 無以菱 2794998160@qq.com
 * @Date: 2025-12-14 15:02:03
 * @LastEditors: 無以菱 2794998160@qq.com
 * @LastEditTime: 2025-12-14 15:20:54
 * @FilePath: \go-vue3-blog\blog-frontend\src\api\ipWhitelist.ts
 * @Description: 这是默认设置,请设置`customMade`, 打开koroFileHeader查看配置 进行设置: https://github.com/OBKoro1/koro1FileHeader/wiki/%E9%85%8D%E7%BD%AE
 */
// IP白名单管理相关 API

import { request } from '@/utils/request'
import type { PageData } from '@/types/common'

export interface IPWhitelist {
  id: number
  ip: string
  reason: string
  expire_at: string | null // 过期时间，null表示永久有效
  created_at: string
  updated_at: string
}

export interface IPWhitelistCheckResult {
  whitelisted: boolean
  info?: IPWhitelist
}

// 获取IP白名单列表
export function getIPWhitelist(params: { page?: number; page_size?: number }) {
  return request.get<PageData<IPWhitelist>>('/admin/ip-whitelist', { params })
}

// 添加IP到白名单
export function addIPWhitelist(data: {
  ip: string
  reason?: string
  duration?: number // 有效期（小时），0表示永久
}) {
  return request.post<IPWhitelist>('/admin/ip-whitelist', data)
}

// 删除IP白名单
export function deleteIPWhitelist(id: number) {
  return request.delete(`/admin/ip-whitelist/${id}`)
}

// 检查IP是否在白名单中
export function checkIPWhitelistStatus(ip: string) {
  return request.get<IPWhitelistCheckResult>('/admin/ip-whitelist/check', { params: { ip } })
}

// 清理过期的白名单记录
export function cleanExpiredIPWhitelist() {
  return request.post<{ deleted_count: number }>('/admin/ip-whitelist/clean-expired')
}



