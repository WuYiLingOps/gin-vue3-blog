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

// 获取博主资料和统计数据
export function getAuthorProfile() {
  return request.get<AuthorProfile>('/blog/author')
}

