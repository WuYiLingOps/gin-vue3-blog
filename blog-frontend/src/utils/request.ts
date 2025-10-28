// Axios 请求封装

import axios, { type AxiosInstance, type AxiosRequestConfig, type AxiosResponse } from 'axios'
import type { ApiResponse } from '@/types/common'
import { useAuthStore } from '@/stores/auth'

// 创建 axios 实例
const service: AxiosInstance = axios.create({
  baseURL: '/api',
  timeout: 15000,
  headers: {
    'Content-Type': 'application/json'
  }
})

// 请求拦截器
service.interceptors.request.use(
  (config) => {
    // 添加 Token
    const authStore = useAuthStore()
    if (authStore.token) {
      config.headers.Authorization = `Bearer ${authStore.token}`
    }
    return config
  },
  (error) => {
    console.error('Request error:', error)
    return Promise.reject(error)
  }
)

// 响应拦截器
service.interceptors.response.use(
  (response: AxiosResponse<ApiResponse>): any => {
    const res = response.data

    // 如果返回的状态码不是 200，则判断为错误
    if (res.code !== 200) {
      // 401: 未授权
      if (res.code === 401) {
        const authStore = useAuthStore()
        authStore.logout()
        window.location.href = '/login'
      }

      return Promise.reject(new Error(res.message || 'Error'))
    }

    return res
  },
  (error): any => {
    console.error('Response error:', error)

    if (error.response) {
      const { status, data } = error.response
      
      switch (status) {
        case 400:
          // 参数错误，使用后端返回的具体错误信息
          if (data?.message) {
            error.message = data.message
          } else {
            error.message = '请求参数错误，请检查输入内容'
          }
          break
        case 401:
          const authStore = useAuthStore()
          authStore.logout()
          window.location.href = '/login'
          error.message = '登录已过期，请重新登录'
          break
        case 403:
          error.message = data?.message || '没有权限访问'
          break
        case 404:
          error.message = data?.message || '请求的资源不存在'
          break
        case 500:
          error.message = data?.message || '服务器错误，请稍后重试'
          break
        default:
          error.message = data?.message || '网络错误，请稍后重试'
      }
    } else if (error.request) {
      // 请求已发出，但没有收到响应
      error.message = '网络连接失败，请检查网络设置'
    } else {
      // 其他错误
      error.message = error.message || '请求失败'
    }

    return Promise.reject(error)
  }
)

// 封装请求方法
export const request = {
  get<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return service.get(url, config)
  },

  post<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return service.post(url, data, config)
  },

  put<T = any>(url: string, data?: any, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return service.put(url, data, config)
  },

  delete<T = any>(url: string, config?: AxiosRequestConfig): Promise<ApiResponse<T>> {
    return service.delete(url, config)
  }
}

/**
 * 获取完整的静态文件 URL
 * 用于拼接上传后返回的相对路径
 * @param path 相对路径，如 /uploads/avatars/xxx.jpg
 * @returns 完整 URL，如 http://localhost:8080/uploads/avatars/xxx.jpg
 */
export function getFileUrl(path: string): string {
  if (!path) return ''
  
  // 如果已经是完整 URL，直接返回
  if (path.startsWith('http://') || path.startsWith('https://')) {
    return path
  }
  
  // 拼接基础 URL
  const baseURL = import.meta.env.VITE_API_BASE_URL || ''
  return baseURL + path
}

export default service

