// 用户管理相关 API

import { request } from '@/utils/request'
import type { User } from '@/types/auth'
import type { PageData } from '@/types/common'

// 获取用户列表
export function getUsers(params: { page?: number; page_size?: number }) {
  return request.get<PageData<User>>('/admin/users', { params })
}

// 获取用户详情
export function getUserById(id: number) {
  return request.get<User>(`/admin/users/${id}`)
}

// 更新用户状态
export function updateUserStatus(id: number, status: number) {
  return request.put(`/admin/users/${id}/status`, { status })
}

// 删除用户
export function deleteUser(id: number) {
  return request.delete(`/admin/users/${id}`)
}

