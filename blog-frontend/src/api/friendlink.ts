/*
 * 项目名称：blog-frontend
 * 文件名称：friendlink.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：友链和友链分类相关 API 接口定义，包括友链和分类的增删改查操作，支持公开访问和管理员管理功能。
 */

import { request } from '@/utils/request'

/**
 * 友链分类接口
 */
export interface FriendLinkCategory {
  id: number              // 分类ID
  name: string            // 分类名称
  description?: string    // 分类描述（可选）
  sort_order: number      // 排序顺序
  created_at: string      // 创建时间
  updated_at: string      // 更新时间
}

/**
 * 友链接口
 */
export interface FriendLink {
  id: number                        // 友链ID
  name: string                      // 友链名称
  url: string                       // 友链URL
  icon?: string                     // 友链图标（可选）
  description?: string              // 友链描述（可选）
  screenshot?: string               // 友链截图（可选）
  atom_url?: string                 // Atom/RSS订阅地址（可选）
  category_id: number               // 所属分类ID
  category?: FriendLinkCategory     // 所属分类信息（可选）
  sort_order: number                // 排序顺序
  status: number                    // 状态（0:待审核 1:已通过 2:已拒绝）
  created_at: string                // 创建时间
  updated_at: string                // 更新时间
}

/**
 * 友链表单数据接口
 */
export interface FriendLinkForm {
  name: string          // 友链名称
  url: string           // 友链URL
  icon?: string         // 友链图标（可选）
  description?: string  // 友链描述（可选）
  screenshot?: string   // 友链截图（可选）
  atom_url?: string     // Atom/RSS订阅地址（可选）
  category_id: number   // 所属分类ID
  sort_order: number    // 排序顺序
  status: number        // 状态
}

/**
 * 友链分类表单数据接口
 */
export interface FriendLinkCategoryForm {
  name: string          // 分类名称
  description?: string  // 分类描述（可选）
  sort_order: number    // 排序顺序
}

/**
 * 获取公开的友链列表（前端用）
 * @returns 返回已通过审核的友链列表
 */
export function getFriendLinks() {
  return request.get<FriendLink[]>('/blog/friend-links')
}

/**
 * 管理员：获取友链列表
 * @param page 页码，默认为1
 * @param pageSize 每页数量，默认为10
 * @returns 返回分页的友链列表数据
 */
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

/**
 * 管理员：获取友链详情
 * @param id 友链ID
 * @returns 返回友链详情
 */
export function getFriendLinkById(id: number) {
  return request.get<FriendLink>(`/admin/friend-links/${id}`)
}

/**
 * 管理员：创建友链
 * @param data 友链表单数据
 * @returns 返回创建的友链数据
 */
export function createFriendLink(data: FriendLinkForm) {
  return request.post<FriendLink>('/admin/friend-links', data)
}

/**
 * 管理员：更新友链
 * @param id 友链ID
 * @param data 要更新的友链数据（部分字段）
 * @returns 返回更新后的友链数据
 */
export function updateFriendLink(id: number, data: Partial<FriendLinkForm>) {
  return request.put<FriendLink>(`/admin/friend-links/${id}`, data)
}

/**
 * 管理员：删除友链
 * @param id 友链ID
 * @returns 返回删除结果
 */
export function deleteFriendLink(id: number) {
  return request.delete(`/admin/friend-links/${id}`)
}

// =============================================================================
// 友链分类相关API
// =============================================================================

/**
 * 获取友链分类列表（公开）
 * @returns 返回所有友链分类列表
 */
export function getFriendLinkCategories() {
  return request.get<FriendLinkCategory[]>('/blog/friend-link-categories')
}

/**
 * 管理员：获取友链分类列表
 * @returns 返回所有友链分类列表
 */
export function getFriendLinkCategoriesAdmin() {
  return request.get<FriendLinkCategory[]>('/admin/friend-link-categories')
}

/**
 * 管理员：获取友链分类详情
 * @param id 分类ID
 * @returns 返回友链分类详情
 */
export function getFriendLinkCategoryById(id: number) {
  return request.get<FriendLinkCategory>(`/admin/friend-link-categories/${id}`)
}

/**
 * 管理员：创建友链分类
 * @param data 友链分类表单数据
 * @returns 返回创建的友链分类数据
 */
export function createFriendLinkCategory(data: FriendLinkCategoryForm) {
  return request.post<FriendLinkCategory>('/admin/friend-link-categories', data)
}

/**
 * 管理员：更新友链分类
 * @param id 分类ID
 * @param data 要更新的友链分类数据（部分字段）
 * @returns 返回更新后的友链分类数据
 */
export function updateFriendLinkCategory(id: number, data: Partial<FriendLinkCategoryForm>) {
  return request.put<FriendLinkCategory>(`/admin/friend-link-categories/${id}`, data)
}

/**
 * 管理员：删除友链分类
 * @param id 分类ID
 * @returns 返回删除结果
 */
export function deleteFriendLinkCategory(id: number) {
  return request.delete(`/admin/friend-link-categories/${id}`)
}

