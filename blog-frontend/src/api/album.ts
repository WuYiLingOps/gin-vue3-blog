/*
 * 项目名称：blog-frontend
 * 文件名称：album.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：相册相关 API 接口定义，包括相册照片的增删改查操作，支持公开访问和管理员管理功能。
 */

import { request } from '@/utils/request'

/**
 * 相册照片数据接口
 */
export interface Album {
  id: number                    // 相册照片ID
  image_url: string             // 图片URL地址
  title?: string                // 照片标题（可选）
  description?: string          // 照片描述（可选）
  sort_order: number            // 排序顺序
  created_at: string            // 创建时间
  updated_at: string            // 更新时间
}

/**
 * 获取公开的相册列表（前端用）
 * @returns 返回相册照片数组
 */
export function getPublicAlbums() {
  return request.get<Album[]>('/blog/albums')
}

/**
 * 获取相册列表（管理员用）
 * @param page 页码，默认为1
 * @param pageSize 每页数量，默认为10
 * @returns 返回分页的相册列表数据
 */
export function getAlbums(page = 1, pageSize = 10) {
  return request.get<{ list: Album[]; total: number; page: number; page_size: number; total_pages: number }>('/admin/albums', {
    params: { page, page_size: pageSize }
  })
}

/**
 * 获取相册照片详情
 * @param id 相册照片ID
 * @returns 返回相册照片详情
 */
export function getAlbum(id: number) {
  return request.get<Album>(`/admin/albums/${id}`)
}

/**
 * 创建相册照片
 * @param data 相册照片数据
 * @param data.image_url 图片URL地址（必填）
 * @param data.title 照片标题（可选）
 * @param data.description 照片描述（可选）
 * @param data.sort_order 排序顺序（可选）
 * @returns 返回创建的相册照片数据
 */
export function createAlbum(data: {
  image_url: string
  title?: string
  description?: string
  sort_order?: number
}) {
  return request.post<Album>('/admin/albums', data)
}

/**
 * 更新相册照片
 * @param id 相册照片ID
 * @param data 要更新的相册照片数据
 * @param data.image_url 图片URL地址（可选）
 * @param data.title 照片标题（可选）
 * @param data.description 照片描述（可选）
 * @param data.sort_order 排序顺序（可选）
 * @returns 返回更新后的相册照片数据
 */
export function updateAlbum(id: number, data: {
  image_url?: string
  title?: string
  description?: string
  sort_order?: number
}) {
  return request.put<Album>(`/admin/albums/${id}`, data)
}

/**
 * 删除相册照片
 * @param id 相册照片ID
 * @returns 返回删除结果
 */
export function deleteAlbum(id: number) {
  return request.delete(`/admin/albums/${id}`)
}
