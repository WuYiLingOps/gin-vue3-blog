import { request } from '@/utils/request'

export interface Notification {
  id: number
  notification_id: number
  user_id: number
  is_read: boolean
  email_sent: boolean
  created_at: string
  notification: {
    id: number
    type: string
    title: string
    content: string
    sender_id: number | null
    target_type: string
    target_id: number | null
    extra: string | null
    status: number
    created_at: string
    updated_at: string
    sender: {
      id: number
      username: string
      nickname: string
      avatar: string
    } | null
  }
}

export interface NotificationListResponse {
  list: Notification[]
  total: number
  page: number
  page_size: number
}

export function getNotifications(params: { page?: number; page_size?: number }) {
  return request.get<NotificationListResponse>('/notifications', { params })
}

export function getUnreadCount() {
  return request.get<{ count: number }>('/notifications/unread-count')
}

export function markAsRead(id: number) {
  return request.put(`/notifications/${id}/read`)
}

export function markAllAsRead() {
  return request.put('/notifications/read-all')
}

export function deleteNotification(id: number) {
  return request.delete(`/notifications/${id}`)
}
