/*
 * 项目名称：blog-frontend
 * 文件名称：post.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文章相关 API 接口定义，包括文章的增删改查、导出、点赞、归档、热门文章、最新文章等功能。
 */

import type { AxiosResponse } from 'axios'
import service, { request } from '@/utils/request'
import type { Post, PostForm, PostQuery } from '@/types/blog'
import type { PageData } from '@/types/common'

/**
 * 获取文章列表
 * @param params 查询参数（分页、分类、标签、关键词等）
 * @returns 返回分页的文章列表
 */
export function getPosts(params: PostQuery) {
  return request.get<PageData<Post>>('/posts', { params })
}

/**
 * 获取文章详情（支持ID或slug）
 * @param idOrSlug 文章ID或slug标识符
 * @returns 返回文章详情
 */
export function getPostById(idOrSlug: number | string) {
  return request.get<Post>(`/posts/${idOrSlug}`)
}

/**
 * 获取文章详情（通过slug，向后兼容）
 * @param slug 文章slug标识符
 * @returns 返回文章详情
 */
export function getPostBySlug(slug: string) {
  return request.get<Post>(`/posts/${slug}`)
}

/**
 * 创建文章
 * @param data 文章表单数据
 * @returns 返回创建的文章数据
 */
export function createPost(data: PostForm) {
  return request.post<Post>('/posts', data)
}

/**
 * 更新文章
 * @param id 文章ID
 * @param data 要更新的文章数据（部分字段）
 * @returns 返回更新后的文章数据
 */
export function updatePost(id: number, data: Partial<PostForm>) {
  return request.put<Post>(`/posts/${id}`, data)
}

/**
 * 删除文章
 * @param id 文章ID
 * @returns 返回删除结果
 */
export function deletePost(id: number) {
  return request.delete(`/posts/${id}`)
}

/**
 * 导出文章（Markdown格式）
 * @param id 文章ID
 * @param format 导出格式，默认为'md'
 * @returns 返回文件下载响应
 */
export function exportPost(id: number, format: 'md' = 'md'): Promise<AxiosResponse<Blob>> {
  return service.get<Blob>(`/admin/posts/${id}/export`, {
    params: { format },
    responseType: 'blob'
  })
}

/**
 * 点赞文章
 * @param id 文章ID
 * @returns 返回点赞结果
 */
export function likePost(id: number) {
  return request.post(`/posts/${id}/like`)
}

/**
 * 获取归档（按月份分组）
 * @returns 返回文章归档数据
 */
export function getArchives() {
  return request.get('/posts/archives')
}

/**
 * 获取热门文章
 * @param limit 返回数量限制，默认为10
 * @returns 返回热门文章列表
 */
export function getHotPosts(limit = 10) {
  return request.get<Post[]>('/posts/hot', { params: { limit } })
}

/**
 * 获取最新文章
 * @param limit 返回数量限制，默认为10
 * @returns 返回最新文章列表
 */
export function getRecentPosts(limit = 10) {
  return request.get<Post[]>('/posts/recent', { params: { limit } })
}

