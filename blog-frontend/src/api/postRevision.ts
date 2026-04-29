/*
 * 项目名称：blog-frontend
 * 文件名称：postRevision.ts
 * 创建时间：2026-04-28 21:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文章修订审批相关 API 接口定义，包括待审批列表、审批通过/拒绝、修改对比等功能。
 */

import { request } from '@/utils/request'
import type { PageData } from '@/types/common'
import type { PostRevision, RevisionDiff } from '@/types/blog'

/**
 * 修订版本查询参数
 */
export interface PostRevisionQuery {
  page?: number
  page_size?: number
  status?: string
}

/**
 * 拒绝请求参数
 */
export interface RejectRequest {
  reject_reason: string
}

/**
 * 获取待审批列表
 * @param params 查询参数（分页）
 * @returns 返回分页的待审批列表
 */
export function getPendingRevisions(params: PostRevisionQuery) {
  return request.get<PageData<PostRevision>>('/admin/revisions/pending', { params })
}

/**
 * 获取待审批数量
 * @returns 返回待审批数量
 */
export function getPendingCount() {
  return request.get<{ count: number }>('/admin/revisions/pending/count')
}

/**
 * 获取修订版本详情
 * @param id 修订版本ID
 * @returns 返回修订版本详情
 */
export function getRevisionDetail(id: number) {
  return request.get<PostRevision>(`/admin/revisions/${id}`)
}

/**
 * 获取修改对比
 * @param id 修订版本ID
 * @returns 返回修改对比数据
 */
export function getRevisionDiff(id: number) {
  return request.get<RevisionDiff>(`/admin/revisions/${id}/diff`)
}

/**
 * 审批通过
 * @param id 修订版本ID
 * @returns 返回审批结果
 */
export function approveRevision(id: number) {
  return request.post(`/admin/revisions/${id}/approve`)
}

/**
 * 审批拒绝
 * @param id 修订版本ID
 * @param data 拒绝原因
 * @returns 返回审批结果
 */
export function rejectRevision(id: number, data: RejectRequest) {
  return request.post(`/admin/revisions/${id}/reject`, data)
}

/**
 * 获取文章的修订历史
 * @param postId 文章ID
 * @returns 返回修订历史列表
 */
export function getPostRevisions(postId: number) {
  return request.get<PostRevision[]>(`/admin/posts/${postId}/revisions`)
}

/**
 * 撤回修订
 * @param id 修订版本ID
 * @returns 返回撤回结果
 */
export function withdrawRevision(id: number) {
  return request.post(`/admin/revisions/${id}/withdraw`)
}

/**
 * 获取当前用户的修订记录
 * @param params 查询参数（分页、状态筛选）
 * @returns 返回当前用户的修订记录列表
 */
export function getMyRevisions(params: PostRevisionQuery) {
  return request.get<PageData<PostRevision>>('/admin/revisions/my', { params })
}
