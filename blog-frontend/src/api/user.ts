/*
 * 项目名称：blog-frontend
 * 文件名称：user.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：用户管理相关 API 接口定义，包括用户列表查询、用户状态管理、用户删除、注册配置等功能（管理员操作）。
 */

import { request } from '@/utils/request'
import type { User } from '@/types/auth'
import type { PageData } from '@/types/common'

/**
 * 获取用户列表（管理员）
 * @param params 分页参数
 * @param params.page 页码（可选）
 * @param params.page_size 每页数量（可选）
 * @returns 返回分页的用户列表
 */
export function getUsers(params: { page?: number; page_size?: number }) {
  return request.get<PageData<User>>('/admin/users', { params })
}

/**
 * 获取用户详情（管理员）
 * @param id 用户ID
 * @returns 返回用户详细信息
 */
export function getUserById(id: number) {
  return request.get<User>(`/admin/users/${id}`)
}

/**
 * 更新用户状态（管理员）
 * @param id 用户ID
 * @param status 用户状态（0:禁用 1:启用）
 * @returns 返回更新结果
 */
export function updateUserStatus(id: number, status: number) {
  return request.put(`/admin/users/${id}/status`, { status })
}

/**
 * 更新用户角色（仅超级管理员）
 * @param id 用户ID
 * @param role 用户角色（super_admin/admin/user）
 * @returns 返回更新结果
 */
export function updateUserRole(id: number, role: 'super_admin' | 'admin' | 'user') {
  return request.put(`/admin/users/${id}/role`, { role })
}

/**
 * 删除用户（管理员）
 * @param id 用户ID
 * @returns 返回删除结果
 */
export function deleteUser(id: number) {
  return request.delete(`/admin/users/${id}`)
}

/**
 * 获取注册配置（管理员）
 * @returns 返回注册配置信息
 */
export function getRegisterSettings() {
  return request.get<{ disable_register: string }>('/admin/settings/register')
}

/**
 * 更新注册配置（管理员）
 * @param data 注册配置数据
 * @param data.disable_register 是否禁用注册：'0'表示否，'1'表示是
 * @returns 返回更新结果
 */
export function updateRegisterSettings(data: { disable_register: string }) {
  return request.put('/admin/settings/register', data)
}
