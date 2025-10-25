import request from '@/utils/request'
import type { ApiResponse, PageParams, PageResult } from '@/types/common'

export interface Moment {
  id: number
  content: string
  images: string
  user_id: number
  status: number
  like_count: number
  created_at: string
  updated_at: string
  liked?: boolean  // 前端标记是否已点赞
  user: {
    id: number
    username: string
    nickname: string
    avatar: string
  }
}

export interface MomentForm {
  content: string
  images?: string
  status?: number
}

export interface MomentParams extends PageParams {
  status?: number
}

// 获取说说列表
export function getMoments(params?: MomentParams) {
  return request.get<ApiResponse<PageResult<Moment>>>('/api/moments', { params })
}

// 获取说说详情
export function getMomentById(id: number) {
  return request.get<ApiResponse<Moment>>(`/api/moments/${id}`)
}

// 获取最新说说
export function getRecentMoments(limit?: number) {
  return request.get<ApiResponse<Moment[]>>('/api/moments/recent', {
    params: { limit }
  })
}

// 创建说说
export function createMoment(data: MomentForm) {
  return request.post<ApiResponse<Moment>>('/api/moments', data)
}

// 更新说说
export function updateMoment(id: number, data: Partial<MomentForm>) {
  return request.put<ApiResponse<null>>(`/api/moments/${id}`, data)
}

// 删除说说
export function deleteMoment(id: number) {
  return request.delete<ApiResponse<null>>(`/api/moments/${id}`)
}

// 点赞说说
export function likeMoment(id: number) {
  return request.post<ApiResponse<null>>(`/api/moments/${id}/like`)
}

// 管理员获取说说列表
export function getAdminMoments(params?: MomentParams) {
  return request.get<ApiResponse<PageResult<Moment>>>('/api/admin/moments', { params })
}

