// 相册相关 API

import { request } from '@/utils/request'

export interface Album {
  id: number
  image_url: string
  title?: string
  description?: string
  sort_order: number
  created_at: string
  updated_at: string
}

// 获取公开的相册列表（前端用）
export function getPublicAlbums() {
  return request.get<Album[]>('/blog/albums')
}

// 获取相册列表（管理员用）
export function getAlbums(page = 1, pageSize = 10) {
  return request.get<{ list: Album[]; total: number; page: number; page_size: number; total_pages: number }>('/admin/albums', {
    params: { page, page_size: pageSize }
  })
}

// 获取相册照片详情
export function getAlbum(id: number) {
  return request.get<Album>(`/admin/albums/${id}`)
}

// 创建相册照片
export function createAlbum(data: {
  image_url: string
  title?: string
  description?: string
  sort_order?: number
}) {
  return request.post<Album>('/admin/albums', data)
}

// 更新相册照片
export function updateAlbum(id: number, data: {
  image_url?: string
  title?: string
  description?: string
  sort_order?: number
}) {
  return request.put<Album>(`/admin/albums/${id}`, data)
}

// 删除相册照片
export function deleteAlbum(id: number) {
  return request.delete(`/admin/albums/${id}`)
}
