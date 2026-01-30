// 博客相关 API

import { request } from '@/utils/request'

export interface AuthorInfo {
  id: number
  username: string
  nickname: string
  avatar: string
  bio: string
}

export interface BlogStats {
  posts: number      // 文章数
  tags: number       // 标签数
  categories: number  // 分类数
}

export interface AuthorProfile {
  author: AuthorInfo
  stats: BlogStats
}

// 公告/系统广播
export interface Announcement {
  id: number
  content: string
  title?: string
  username: string
  avatar?: string
  priority: number
  is_broadcast: boolean
  target?: 'announcement' | 'chat' | 'both'
  created_at: string
  updated_at: string
}

// 获取博主资料和统计数据
export function getAuthorProfile() {
  return request.get<AuthorProfile>('/blog/author')
}

// 获取公告列表
export function getAnnouncements(limit = 3) {
  return request.get<Announcement[]>('/blog/announcements', { params: { limit } })
}

// 获取公告详情
export function getAnnouncementDetail(id: number) {
  return request.get<Announcement>(`/blog/announcements/${id}`)
}

// 标签统计数据
export interface TagStat {
  name: string
  value: number
}

// 归档数据（按月）
export interface ArchiveStat {
  month: string  // 格式：YYYY-MM-DD HH:mm:ss
  count: number
}

// 获取标签统计（TOP10，公开接口）
export function getPublicTagStats() {
  return request.get<TagStat[]>('/blog/tag-stats')
}

// 网站资讯
export interface WebsiteInfo {
  total_words: number       // 本站总字数
  total_visitors: number    // 本站访客数
  total_views: number       // 本站总访问量
  last_update_time: string  // 最后更新时间
}

// 获取网站资讯（公开接口）
export function getWebsiteInfo() {
  return request.get<WebsiteInfo>('/blog/website-info')
}