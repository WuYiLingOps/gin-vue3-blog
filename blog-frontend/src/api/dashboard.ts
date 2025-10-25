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

export interface VisitStat {
  date: string  // 日期格式：MM-DD
  count: number // 访问量
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

/**
 * 获取最近7天访问量统计
 */
export function getVisitStats() {
  return request.get<VisitStat[]>('/admin/dashboard/visit-stats')
}

