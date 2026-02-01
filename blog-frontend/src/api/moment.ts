/*
 * 项目名称：blog-frontend
 * 文件名称：moment.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：说说（动态）相关 API 接口定义，包括说说的增删改查、点赞、最新说说获取等功能。
 */

import request from '@/utils/request'
import type { PageParams, PageResult } from '@/types/common'

/**
 * 说说（动态）接口
 */
export interface Moment {
  id: number                // 说说ID
  content: string            // 说说内容
  images: string            // 图片列表（JSON字符串）
  user_id: number           // 用户ID
  status: number            // 状态（0:待审核 1:已发布 2:已删除）
  like_count: number        // 点赞数
  created_at: string        // 创建时间
  updated_at: string        // 更新时间
  liked?: boolean           // 前端标记是否已点赞
  user: {                   // 发布者信息
    id: number
    username: string
    nickname: string
    avatar: string
  }
}

/**
 * 说说表单数据接口
 */
export interface MomentForm {
  content: string      // 说说内容
  images?: string      // 图片列表（JSON字符串，可选）
  status?: number      // 状态（可选）
}

/**
 * 说说查询参数接口
 */
export interface MomentParams extends PageParams {
  status?: number      // 状态筛选（可选）
}

/**
 * 获取说说列表
 * @param params 查询参数（分页、状态筛选等）
 * @returns 返回分页的说说列表
 */
export function getMoments(params?: MomentParams) {
  return request.get<PageResult<Moment>>('/moments', { params })
}

/**
 * 获取说说详情
 * @param id 说说ID
 * @returns 返回说说详情
 */
export function getMomentById(id: number) {
  return request.get<Moment>(`/moments/${id}`)
}

/**
 * 获取最新说说
 * @param limit 返回数量限制（可选）
 * @returns 返回最新的说说列表
 */
export function getRecentMoments(limit?: number) {
  return request.get<Moment[]>('/moments/recent', {
    params: { limit }
  })
}

/**
 * 创建说说
 * @param data 说说表单数据
 * @returns 返回创建的说说数据
 */
export function createMoment(data: MomentForm) {
  return request.post<Moment>('/moments', data)
}

/**
 * 更新说说
 * @param id 说说ID
 * @param data 要更新的说说数据（部分字段）
 * @returns 返回更新结果
 */
export function updateMoment(id: number, data: Partial<MomentForm>) {
  return request.put<null>(`/moments/${id}`, data)
}

/**
 * 删除说说
 * @param id 说说ID
 * @returns 返回删除结果
 */
export function deleteMoment(id: number) {
  return request.delete<null>(`/moments/${id}`)
}

/**
 * 点赞说说
 * @param id 说说ID
 * @returns 返回点赞结果
 */
export function likeMoment(id: number) {
  return request.post<null>(`/moments/${id}/like`)
}

/**
 * 管理员获取说说列表
 * @param params 查询参数（分页、状态筛选等）
 * @returns 返回分页的说说列表（包含所有状态）
 */
export function getAdminMoments(params?: MomentParams) {
  return request.get<PageResult<Moment>>('/admin/moments', { params })
}

