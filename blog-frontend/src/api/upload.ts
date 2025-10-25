// 上传相关 API

import { request, getFileUrl } from '@/utils/request'

export interface UploadResponse {
  url: string
}

// 上传头像
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

// 上传图片
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

