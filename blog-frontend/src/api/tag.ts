// 标签相关 API

import { request } from '@/utils/request'
import type { Tag, TagForm, Post } from '@/types/blog'
import type { PageData } from '@/types/common'

// 获取标签列表
export function getTags() {
  return request.get<Tag[]>('/tags')
}

// 获取标签详情
export function getTagById(id: number) {
  return request.get<Tag>(`/tags/${id}`)
}

// 获取标签下的文章
export function getPostsByTag(id: number, params: { page?: number; page_size?: number }) {
  return request.get<PageData<Post>>(`/tags/${id}/posts`, { params })
}

// 创建标签
export function createTag(data: TagForm) {
  return request.post<Tag>('/tags', data)
}

// 更新标签
export function updateTag(id: number, data: Partial<TagForm>) {
  return request.put<Tag>(`/tags/${id}`, data)
}

// 删除标签
export function deleteTag(id: number) {
  return request.delete(`/tags/${id}`)
}

