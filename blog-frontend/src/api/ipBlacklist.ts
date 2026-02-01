/*
 * 项目名称：blog-frontend
 * 文件名称：ipBlacklist.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：IP黑名单管理相关 API 接口定义，包括IP黑名单的增删改查、IP状态检查、过期记录清理等功能（管理员操作）。
 */

import { request } from '@/utils/request'
import type { PageData } from '@/types/common'

/**
 * IP黑名单接口
 */
export interface IPBlacklist {
  id: number                // 记录ID
  ip: string                // IP地址
  reason: string            // 封禁原因
  ban_type: number          // 封禁类型：1表示自动封禁，2表示手动封禁
  expire_at: string | null  // 过期时间，null表示永久封禁
  created_at: string        // 创建时间
  updated_at: string        // 更新时间
}

/**
 * IP检查结果接口
 */
export interface IPCheckResult {
  banned: boolean           // 是否被封禁
  info?: IPBlacklist        // 封禁信息（如果被封禁）
}

/**
 * 获取IP黑名单列表（管理员）
 * @param params 分页参数
 * @param params.page 页码（可选）
 * @param params.page_size 每页数量（可选）
 * @returns 返回分页的IP黑名单列表
 */
export function getIPBlacklist(params: { page?: number; page_size?: number }) {
  return request.get<PageData<IPBlacklist>>('/admin/ip-blacklist', { params })
}

/**
 * 添加IP到黑名单（管理员）
 * @param data IP黑名单数据
 * @param data.ip IP地址（必填）
 * @param data.reason 封禁原因（可选）
 * @param data.duration 封禁时长（小时），0表示永久（可选）
 * @returns 返回创建的IP黑名单记录
 */
export function addIPBlacklist(data: {
  ip: string
  reason?: string
  duration?: number // 封禁时长（小时），0表示永久
}) {
  return request.post<IPBlacklist>('/admin/ip-blacklist', data)
}

/**
 * 删除IP黑名单（管理员）
 * @param id 记录ID
 * @returns 返回删除结果
 */
export function deleteIPBlacklist(id: number) {
  return request.delete(`/admin/ip-blacklist/${id}`)
}

/**
 * 检查IP是否在黑名单中（管理员）
 * @param ip IP地址
 * @returns 返回IP检查结果
 */
export function checkIPStatus(ip: string) {
  return request.get<IPCheckResult>('/admin/ip-blacklist/check', { params: { ip } })
}

/**
 * 清理过期的黑名单记录（管理员）
 * @returns 返回清理结果，包含删除的记录数量
 */
export function cleanExpiredIPBlacklist() {
  return request.post<{ deleted_count: number }>('/admin/ip-blacklist/clean-expired')
}

