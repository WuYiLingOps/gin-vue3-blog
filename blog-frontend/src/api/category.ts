/*
 * 项目名称：blog-frontend
 * 文件名称：category.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文章分类相关 API 接口定义，包括分类的增删改查操作。
 */

import { request } from '@/utils/request'
import type { Category, CategoryForm } from '@/types/blog'

/**
 * 获取分类列表
 * @returns 返回所有分类列表
 */
export function getCategories() {
  return request.get<Category[]>('/categories')
}

/**
 * 获取分类详情
 * @param id 分类ID
 * @returns 返回分类详情
 */
export function getCategoryById(id: number) {
  return request.get<Category>(`/categories/${id}`)
}

/**
 * 创建分类
 * @param data 分类表单数据
 * @returns 返回创建的分类数据
 */
export function createCategory(data: CategoryForm) {
  return request.post<Category>('/categories', data)
}

/**
 * 更新分类
 * @param id 分类ID
 * @param data 要更新的分类数据（部分字段）
 * @returns 返回更新后的分类数据
 */
export function updateCategory(id: number, data: Partial<CategoryForm>) {
  return request.put<Category>(`/categories/${id}`, data)
}

/**
 * 删除分类
 * @param id 分类ID
 * @returns 返回删除结果
 */
export function deleteCategory(id: number) {
  return request.delete(`/categories/${id}`)
}

