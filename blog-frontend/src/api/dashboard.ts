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
  date: string
  count: number
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
 * 获取最近 N 天访问量统计（按天）
 * @param days 最近天数，默认 7 天，最大 30 天
 */
export function getVisitStats(days = 7) {
  return request.get<VisitStat[]>('/admin/dashboard/visit-stats', {
    params: { days }
  })
}

