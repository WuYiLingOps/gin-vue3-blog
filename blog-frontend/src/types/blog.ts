// 博客相关类型

import type { User } from './auth'

// 分类
export interface Category {
  id: number
  name: string
  description: string
  color: string
  sort: number
  post_count: number
  created_at: string
  updated_at: string
}

// 标签
export interface Tag {
  id: number
  name: string
  color: string
  text_color?: string
  font_size?: number
  post_count: number
  created_at: string
  updated_at: string
}

// 文章
export interface Post {
  id: number
  title: string
  content: string
  summary: string
  cover: string
  status: number
  is_top: boolean
  view_count: number
  like_count: number
  liked?: boolean
  user_id: number
  category_id: number
  published_at: string
  created_at: string
  updated_at: string
  user: User
  category: Category
  tags: Tag[]
  comments?: Comment[]
}

// 评论
export interface Comment {
  id: number
  content: string
  post_id: number
  user_id: number
  parent_id?: number
  status: number
  created_at: string
  updated_at: string
  user: User
  parent?: Comment
  children?: Comment[]
  post?: {
    id: number
    title: string
  }
}

// 文章表单
export interface PostForm {
  title: string
  content: string
  summary: string
  cover: string
  category_id: number
  tag_ids: number[]
  status: number
  is_top: boolean
}

// 分类表单
export interface CategoryForm {
  name: string
  description: string
  color: string
  sort: number
}

// 标签表单
export interface TagForm {
  name: string
  color: string
  text_color?: string
  font_size?: number
}

// 评论表单
export interface CommentForm {
  content: string
  post_id: number
  parent_id?: number
}

// 文章查询参数
export interface PostQuery {
  page?: number
  page_size?: number
  category_id?: number
  keyword?: string
  status?: number
}

