import { request } from '@/utils/request'

export interface DashboardStats {
  posts: number
  users: number
  comments: number
  views: number
}

export interface CategoryStat {
  name: string
  value: number
  color: string
}

/**
 * 获取仪表盘统计数据
 */
export function getDashboardStats() {
  return request.get<DashboardStats>('/admin/dashboard/stats')
}

/**
 * 获取分类统计数据
 */
export function getCategoryStats() {
  return request.get<CategoryStat[]>('/admin/dashboard/category-stats')
}

