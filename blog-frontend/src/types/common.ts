// 通用类型定义

// 分页参数
export interface PaginationParams {
  page: number
  page_size: number
}

// 分页数据
export interface PageData<T> {
  list: T[]
  total: number
  page: number
  page_size: number
  total_pages: number
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

