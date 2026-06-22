/*
 * 项目名称：blog-frontend
 * 文件名称：rss.ts
 * 创建时间：2026-02-02 15:30:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：RSS 订阅相关 API 接口定义，包括 RSS 配置管理、Feed 预览、缓存清理等功能。
 */

import { request } from '@/utils/request'

/**
 * RSS 配置接口
 */
export interface RSSConfig {
  enabled?: boolean // 是否启用 RSS
  title?: string // RSS 标题
  description?: string // RSS 描述
  link?: string // 网站链接
  author_name?: string // 作者名称
  author_email?: string // 作者邮箱
  language?: string // 语言（如 zh-CN）
  copyright?: string // 版权信息
  item_limit?: number // 每个 Feed 的文章数量限制
  cache_duration?: number // 缓存时长（分钟）
}

/**
 * RSS 统计信息接口
 */
export interface RSSStats {
  posts_count?: number // 文章数量
  moments_count?: number // 说说数量
  categories_count?: number // 分类数量
  tags_count?: number // 标签数量
  last_update?: string // 最后更新时间
}

/**
 * 获取 RSS 配置（管理员）
 * @returns 返回 RSS 配置信息
 */
export function getRSSConfig() {
  return request.get<RSSConfig>('/admin/rss/config')
}

/**
 * 更新 RSS 配置（管理员）
 * @param data RSS 配置数据
 * @returns 返回更新结果
 */
export function updateRSSConfig(data: RSSConfig) {
  return request.put('/admin/rss/config', data)
}

/**
 * 预览 RSS Feed（管理员）
 * @param type Feed 类型：posts | moments | all
 * @returns 返回 RSS XML 内容
 */
export function previewRSSFeed(type: 'posts' | 'moments' | 'all') {
  return request.get<string>('/admin/rss/preview', { params: { type } })
}

/**
 * 清除 RSS 缓存（管理员）
 * @returns 返回清除结果
 */
export function clearRSSCache() {
  return request.post('/admin/rss/clear-cache')
}

/**
 * 获取 RSS 统计信息（管理员）
 * @returns 返回 RSS 统计信息
 */
export function getRSSStats() {
  return request.get<RSSStats>('/admin/rss/stats')
}

/**
 * 获取公开的 RSS 配置状态
 * @returns 返回 RSS 是否启用
 */
export function getRSSStatus() {
  return request.get<{ enabled: boolean }>('/rss/status')
}
