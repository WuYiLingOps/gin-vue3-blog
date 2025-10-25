// 文章相关 API

import { request } from '@/utils/request'
import type { Post, PostForm, PostQuery } from '@/types/blog'
import type { PageData } from '@/types/common'

// 获取文章列表
export function getPosts(params: PostQuery) {
  return request.get<PageData<Post>>('/api/posts', { params })
}

// 获取文章详情
export function getPostById(id: number) {
  return request.get<Post>(`/api/posts/${id}`)
}

// 创建文章
export function createPost(data: PostForm) {
  return request.post<Post>('/api/posts', data)
}

// 更新文章
export function updatePost(id: number, data: Partial<PostForm>) {
  return request.put<Post>(`/api/posts/${id}`, data)
}

// 删除文章
export function deletePost(id: number) {
  return request.delete(`/api/posts/${id}`)
}

// 点赞文章
export function likePost(id: number) {
  return request.post(`/api/posts/${id}/like`)
}

// 获取归档
export function getArchives() {
  return request.get('/api/posts/archives')
}

// 获取热门文章
export function getHotPosts(limit = 10) {
  return request.get<Post[]>('/api/posts/hot', { params: { limit } })
}

// 获取最新文章
export function getRecentPosts(limit = 10) {
  return request.get<Post[]>('/api/posts/recent', { params: { limit } })
}

