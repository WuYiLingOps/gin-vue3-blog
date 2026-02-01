/*
 * 项目名称：blog-frontend
 * 文件名称：dashboard.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：管理后台仪表盘相关 API 接口定义，包括统计数据、分类统计、访问量统计等功能。
 */

import { request } from '@/utils/request'

/**
 * 仪表盘统计数据接口
 */
export interface DashboardStats {
  posts: number        // 文章总数
  users: number        // 用户总数
  comments: number     // 评论总数
  views: number        // 访问量总数
}

/**
 * 分类统计数据接口
 */
export interface CategoryStat {
  name: string        // 分类名称
  value: number       // 该分类下的文章数量
  color: string       // 图表显示颜色
}

/**
 * 访问量统计数据接口
 */
export interface VisitStat {
  date: string        // 日期
  count: number       // 该日期的访问量
}

/**
 * 获取仪表盘统计数据
 * @returns 返回文章、用户、评论、访问量等统计数据
 */
export function getDashboardStats() {
  return request.get<DashboardStats>('/admin/dashboard/stats')
}

/**
 * 获取分类统计数据
 * @returns 返回各分类下的文章数量统计
 */
export function getCategoryStats() {
  return request.get<CategoryStat[]>('/admin/dashboard/category-stats')
}

/**
 * 获取最近 N 天访问量统计（按天）
 * @param days 最近天数，默认 7 天，最大 30 天
 * @returns 返回每天的访问量统计数据
 */
export function getVisitStats(days = 7) {
  return request.get<VisitStat[]>('/admin/dashboard/visit-stats', {
    params: { days }
  })
}

