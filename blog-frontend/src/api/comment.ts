// 评论相关 API

import { request } from '@/utils/request'
import type { Comment, CommentForm } from '@/types/blog'
import type { PageData } from '@/types/common'

// 获取文章的评论列表
export function getCommentsByPostId(postId: number) {
  return request.get<Comment[]>(`/comments/post/${postId}`)
}

// 根据评论类型和目标ID获取评论列表（用于友链等特殊页面）
export function getCommentsByTypeAndTarget(commentType: string, targetId: number = 0) {
  return request.get<Comment[]>(`/comments/type?type=${commentType}&target_id=${targetId}`)
}

// 创建评论
export function createComment(data: CommentForm) {
  return request.post<Comment>('/comments', data)
}

// 更新评论
export function updateComment(id: number, data: { content: string }) {
  return request.put<Comment>(`/comments/${id}`, data)
}

// 删除评论
export function deleteComment(id: number) {
  return request.delete(`/comments/${id}`)
}

// 获取所有评论（管理后台）
export function getAllComments(params: { page?: number; page_size?: number }) {
  return request.get<PageData<Comment>>('/admin/comments', { params })
}

// 更新评论状态
export function updateCommentStatus(id: number, status: number) {
  return request.put(`/admin/comments/${id}/status`, { status })
}

