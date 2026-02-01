/*
 * 项目名称：blog-frontend
 * 文件名称：blog.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：博客相关 API 接口定义，包括博主信息、公告、标签统计、归档、网站资讯等公开接口。
 */

import { request } from '@/utils/request'

/**
 * 博主基本信息接口
 */
export interface AuthorInfo {
  id: number              // 用户ID
  username: string        // 用户名
  nickname: string        // 昵称
  avatar: string          // 头像URL
  bio: string             // 个人简介
}

/**
 * 博客统计数据接口
 */
export interface BlogStats {
  posts: number           // 文章数
  tags: number            // 标签数
  categories: number      // 分类数
}

/**
 * 博主资料和统计数据接口
 */
export interface AuthorProfile {
  author: AuthorInfo      // 博主基本信息
  stats: BlogStats        // 博客统计数据
}

/**
 * 公告/系统广播接口
 */
export interface Announcement {
  id: number                                    // 公告ID
  content: string                               // 公告内容
  title?: string                                // 公告标题（可选）
  username: string                              // 发布者用户名
  avatar?: string                               // 发布者头像（可选）
  priority: number                              // 优先级
  is_broadcast: boolean                         // 是否为广播
  target?: 'announcement' | 'chat' | 'both'     // 目标位置：公告栏、聊天室或两者
  created_at: string                            // 创建时间
  updated_at: string                            // 更新时间
}

/**
 * 获取博主资料和统计数据
 * @returns 返回博主信息和博客统计数据
 */
export function getAuthorProfile() {
  return request.get<AuthorProfile>('/blog/author')
}

/**
 * 获取公告列表
 * @param limit 返回数量限制，默认为3
 * @returns 返回公告列表
 */
export function getAnnouncements(limit = 3) {
  return request.get<Announcement[]>('/blog/announcements', { params: { limit } })
}

/**
 * 获取公告详情
 * @param id 公告ID
 * @returns 返回公告详情
 */
export function getAnnouncementDetail(id: number) {
  return request.get<Announcement>(`/blog/announcements/${id}`)
}

/**
 * 标签统计数据接口
 */
export interface TagStat {
  name: string           // 标签名称
  value: number          // 标签使用次数
}

/**
 * 归档数据接口（按月统计）
 */
export interface ArchiveStat {
  month: string          // 月份，格式：YYYY-MM-DD HH:mm:ss
  count: number          // 该月文章数量
}

/**
 * 获取标签统计（TOP10，公开接口）
 * @returns 返回标签统计数据数组
 */
export function getPublicTagStats() {
  return request.get<TagStat[]>('/blog/tag-stats')
}

/**
 * 网站资讯接口
 */
export interface WebsiteInfo {
  total_words: number       // 本站总字数
  total_visitors: number    // 本站访客数
  total_views: number       // 本站总访问量
  last_update_time: string  // 最后更新时间
}

/**
 * 获取网站资讯（公开接口）
 * @returns 返回网站资讯数据
 */
export function getWebsiteInfo() {
  return request.get<WebsiteInfo>('/blog/website-info')
}