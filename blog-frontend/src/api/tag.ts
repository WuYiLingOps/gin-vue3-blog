/*
 * 项目名称：blog-frontend
 * 文件名称：tag.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：标签相关 API 接口定义，包括标签的增删改查、获取标签下的文章列表等操作。
 */

import { request } from '@/utils/request'
import type { Tag, TagForm, Post } from '@/types/blog'
import type { PageData } from '@/types/common'

/**
 * 获取标签列表
 * @returns 返回所有标签列表
 */
export function getTags() {
  return request.get<Tag[]>('/tags')
}

/**
 * 获取标签详情
 * @param id 标签ID
 * @returns 返回标签详情
 */
export function getTagById(id: number) {
  return request.get<Tag>(`/tags/${id}`)
}

/**
 * 获取标签下的文章列表
 * @param id 标签ID
 * @param params 分页参数
 * @param params.page 页码（可选）
 * @param params.page_size 每页数量（可选）
 * @returns 返回该标签下的分页文章列表
 */
export function getPostsByTag(id: number, params: { page?: number; page_size?: number }) {
  return request.get<PageData<Post>>(`/tags/${id}/posts`, { params })
}

/**
 * 创建标签
 * @param data 标签表单数据
 * @returns 返回创建的标签数据
 */
export function createTag(data: TagForm) {
  return request.post<Tag>('/tags', data)
}

/**
 * 更新标签
 * @param id 标签ID
 * @param data 要更新的标签数据（部分字段）
 * @returns 返回更新后的标签数据
 */
export function updateTag(id: number, data: Partial<TagForm>) {
  return request.put<Tag>(`/tags/${id}`, data)
}

/**
 * 删除标签
 * @param id 标签ID
 * @returns 返回删除结果
 */
export function deleteTag(id: number) {
  return request.delete(`/tags/${id}`)
}

