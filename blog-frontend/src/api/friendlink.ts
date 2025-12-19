import { request } from '@/utils/request'

export interface FriendLinkCategory {
  id: number
  name: string
  description?: string
  sort_order: number
  created_at: string
  updated_at: string
}

export interface FriendLink {
  id: number
  name: string
  url: string
  icon?: string
  description?: string
  screenshot?: string
  atom_url?: string
  category_id: number
  category?: FriendLinkCategory
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
  category_id: number
  sort_order: number
  status: number
}

export interface FriendLinkCategoryForm {
  name: string
  description?: string
  sort_order: number
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

// =============================================================================
// 友链分类相关API
// =============================================================================

// 获取友链分类列表（公开）
export function getFriendLinkCategories() {
  return request.get<FriendLinkCategory[]>('/blog/friend-link-categories')
}

// 管理员：获取友链分类列表
export function getFriendLinkCategoriesAdmin() {
  return request.get<FriendLinkCategory[]>('/admin/friend-link-categories')
}

// 管理员：获取友链分类详情
export function getFriendLinkCategoryById(id: number) {
  return request.get<FriendLinkCategory>(`/admin/friend-link-categories/${id}`)
}

// 管理员：创建友链分类
export function createFriendLinkCategory(data: FriendLinkCategoryForm) {
  return request.post<FriendLinkCategory>('/admin/friend-link-categories', data)
}

// 管理员：更新友链分类
export function updateFriendLinkCategory(id: number, data: Partial<FriendLinkCategoryForm>) {
  return request.put<FriendLinkCategory>(`/admin/friend-link-categories/${id}`, data)
}

// 管理员：删除友链分类
export function deleteFriendLinkCategory(id: number) {
  return request.delete(`/admin/friend-link-categories/${id}`)
}

