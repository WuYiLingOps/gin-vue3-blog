/*
 * 项目名称：blog-frontend
 * 文件名称：ipWhitelist.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：IP白名单管理相关 API 接口定义，包括IP白名单的增删改查、IP状态检查、过期记录清理等功能（管理员操作）。
 */

import { request } from '@/utils/request'
import type { PageData } from '@/types/common'

/**
 * IP白名单接口
 */
export interface IPWhitelist {
  id: number                // 记录ID
  ip: string                // IP地址
  reason: string            // 添加原因
  expire_at: string | null  // 过期时间，null表示永久有效
  created_at: string        // 创建时间
  updated_at: string        // 更新时间
}

/**
 * IP白名单检查结果接口
 */
export interface IPWhitelistCheckResult {
  whitelisted: boolean       // 是否在白名单中
  info?: IPWhitelist         // 白名单信息（如果在白名单中）
}

/**
 * 获取IP白名单列表（管理员）
 * @param params 分页参数
 * @param params.page 页码（可选）
 * @param params.page_size 每页数量（可选）
 * @returns 返回分页的IP白名单列表
 */
export function getIPWhitelist(params: { page?: number; page_size?: number }) {
  return request.get<PageData<IPWhitelist>>('/admin/ip-whitelist', { params })
}

/**
 * 添加IP到白名单（管理员）
 * @param data IP白名单数据
 * @param data.ip IP地址（必填）
 * @param data.reason 添加原因（可选）
 * @param data.duration 有效期（小时），0表示永久（可选）
 * @returns 返回创建的IP白名单记录
 */
export function addIPWhitelist(data: {
  ip: string
  reason?: string
  duration?: number // 有效期（小时），0表示永久
}) {
  return request.post<IPWhitelist>('/admin/ip-whitelist', data)
}

/**
 * 删除IP白名单（管理员）
 * @param id 记录ID
 * @returns 返回删除结果
 */
export function deleteIPWhitelist(id: number) {
  return request.delete(`/admin/ip-whitelist/${id}`)
}

/**
 * 检查IP是否在白名单中（管理员）
 * @param ip IP地址
 * @returns 返回IP检查结果
 */
export function checkIPWhitelistStatus(ip: string) {
  return request.get<IPWhitelistCheckResult>('/admin/ip-whitelist/check', { params: { ip } })
}

/**
 * 清理过期的白名单记录（管理员）
 * @returns 返回清理结果，包含删除的记录数量
 */
export function cleanExpiredIPWhitelist() {
  return request.post<{ deleted_count: number }>('/admin/ip-whitelist/clean-expired')
}



