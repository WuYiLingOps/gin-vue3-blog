/*
 * 项目名称：blog-frontend
 * 文件名称：upload.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文件上传相关 API 接口定义，包括头像上传、图片上传等功能，支持本地存储和OSS存储。
 */

import { request, getFileUrl } from '@/utils/request'

/**
 * 上传响应接口
 */
export interface UploadResponse {
  url: string  // 上传后的文件URL地址
}

/**
 * 上传头像
 * @param file 头像文件对象
 * @returns 返回上传结果，包含文件URL（已转换为完整URL）
 */
export async function uploadAvatar(file: File) {
  const formData = new FormData()
  formData.append('file', file)

  const result = await request.post<UploadResponse>('/upload/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })

  // 将相对路径转换为完整 URL
  if (result.data?.url) {
    result.data.url = getFileUrl(result.data.url)
  }

  return result
}

/**
 * 上传图片
 * @param file 图片文件对象
 * @returns 返回上传结果，包含文件URL（已转换为完整URL）
 */
export async function uploadImage(file: File) {
  const formData = new FormData()
  formData.append('file', file)

  const result = await request.post<UploadResponse>('/upload/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })

  // 将相对路径转换为完整 URL
  if (result.data?.url) {
    result.data.url = getFileUrl(result.data.url)
  }

  return result
}

