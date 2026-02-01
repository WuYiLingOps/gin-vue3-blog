/*
 * 项目名称：blog-frontend
 * 文件名称：comment.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：评论相关 API 接口定义，包括评论的增删改查、评论状态管理等操作。
 */

import { request } from '@/utils/request'
import type { Comment, CommentForm } from '@/types/blog'
import type { PageData } from '@/types/common'

/**
 * 获取文章的评论列表
 * @param postId 文章ID
 * @returns 返回该文章的所有评论列表
 */
export function getCommentsByPostId(postId: number) {
  return request.get<Comment[]>(`/comments/post/${postId}`)
}

/**
 * 根据评论类型和目标ID获取评论列表（用于友链等特殊页面）
 * @param commentType 评论类型
 * @param targetId 目标ID，默认为0
 * @returns 返回符合条件的评论列表
 */
export function getCommentsByTypeAndTarget(commentType: string, targetId: number = 0) {
  return request.get<Comment[]>(`/comments/type?type=${commentType}&target_id=${targetId}`)
}

/**
 * 创建评论
 * @param data 评论表单数据
 * @returns 返回创建的评论数据
 */
export function createComment(data: CommentForm) {
  return request.post<Comment>('/comments', data)
}

/**
 * 更新评论
 * @param id 评论ID
 * @param data 要更新的评论数据
 * @param data.content 评论内容
 * @returns 返回更新后的评论数据
 */
export function updateComment(id: number, data: { content: string }) {
  return request.put<Comment>(`/comments/${id}`, data)
}

/**
 * 删除评论
 * @param id 评论ID
 * @returns 返回删除结果
 */
export function deleteComment(id: number) {
  return request.delete(`/comments/${id}`)
}

/**
 * 获取所有评论（管理后台）
 * @param params 分页参数
 * @param params.page 页码（可选）
 * @param params.page_size 每页数量（可选）
 * @returns 返回分页的评论列表
 */
export function getAllComments(params: { page?: number; page_size?: number }) {
  return request.get<PageData<Comment>>('/admin/comments', { params })
}

/**
 * 更新评论状态（管理员操作）
 * @param id 评论ID
 * @param status 评论状态
 * @returns 返回更新结果
 */
export function updateCommentStatus(id: number, status: number) {
  return request.put(`/admin/comments/${id}/status`, { status })
}

