// 上传相关 API

import { request } from '@/utils/request'

export interface UploadResponse {
  url: string
}

// 上传头像
export function uploadAvatar(file: File) {
  const formData = new FormData()
  formData.append('file', file)

  return request.post<UploadResponse>('/api/upload/avatar', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

// 上传图片
export function uploadImage(file: File) {
  const formData = new FormData()
  formData.append('file', file)

  return request.post<UploadResponse>('/api/upload/image', formData, {
    headers: {
      'Content-Type': 'multipart/form-data'
    }
  })
}

