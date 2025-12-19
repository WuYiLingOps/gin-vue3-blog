import { request } from '@/utils/request'

export interface FriendLink {
  id: number
  name: string
  url: string
  icon?: string
  description?: string
  screenshot?: string
  atom_url?: string
  sort_order: number
  status: number
  created_at: string
  updated_at: string
}

export interface FriendLinkForm {
  name: string
  url: string
  icon?: string
  description?: string
  screenshot?: string
  atom_url?: string
  sort_order: number
  status: number
}

// 获取公开的友链列表（前端用）
export function getFriendLinks() {
  return request.get<FriendLink[]>('/blog/friend-links')
}

// 管理员：获取友链列表
export function getFriendLinksAdmin(page = 1, pageSize = 10) {
  return request.get<{
    list: FriendLink[]
    total: number
    page: number
    page_size: number
  }>('/admin/friend-links', {
    params: { page, page_size: pageSize }
  })
}

// 管理员：获取友链详情
export function getFriendLinkById(id: number) {
  return request.get<FriendLink>(`/admin/friend-links/${id}`)
}

// 管理员：创建友链
export function createFriendLink(data: FriendLinkForm) {
  return request.post<FriendLink>('/admin/friend-links', data)
}

// 管理员：更新友链
export function updateFriendLink(id: number, data: Partial<FriendLinkForm>) {
  return request.put<FriendLink>(`/admin/friend-links/${id}`, data)
}

// 管理员：删除友链
export function deleteFriendLink(id: number) {
  return request.delete(`/admin/friend-links/${id}`)
}

