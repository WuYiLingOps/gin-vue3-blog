/*
 * @ProjectName: go-vue3-blog
 * @FileName: common.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 通用类型定义，包含分页、API响应、表单规则等通用类型
 */

// 通用类型定义

// 分页参数
export interface PaginationParams {
  page: number
  page_size: number
}

// 分页参数（别名）
export interface PageParams {
  page?: number
  page_size?: number
}

// 分页数据
export interface PageData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
}

// 分页结果（别名）
export interface PageResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
  total_pages?: number
}

// 分页结果（别名2）
export interface PaginationResult<T> {
  list: T[]
  total: number
  page: number
  page_size: number
  total_pages?: number
}

// API 响应
export interface ApiResponse<T = any> {
  code: number
  message: string
  data?: T
}

// 表单规则
export interface FormRule {
  required?: boolean
  message?: string
  trigger?: string | string[]
  min?: number
  max?: number
  pattern?: RegExp
  validator?: (rule: any, value: any) => Promise<void> | boolean
}

