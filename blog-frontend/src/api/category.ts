// 分类相关 API

import { request } from '@/utils/request'
import type { Category, CategoryForm } from '@/types/blog'

// 获取分类列表
export function getCategories() {
  return request.get<Category[]>('/categories')
}

// 获取分类详情
export function getCategoryById(id: number) {
  return request.get<Category>(`/categories/${id}`)
}

// 创建分类
export function createCategory(data: CategoryForm) {
  return request.post<Category>('/categories', data)
}

// 更新分类
export function updateCategory(id: number, data: Partial<CategoryForm>) {
  return request.put<Category>(`/categories/${id}`, data)
}

// 删除分类
export function deleteCategory(id: number) {
  return request.delete(`/categories/${id}`)
}

