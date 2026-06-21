<template>
  <n-popover trigger="click" placement="bottom" :width="380" :style="{ padding: 0 }">
    <template #trigger>
      <n-badge :value="notificationStore.unreadCount" :max="99" :offset="[-4, 0]">
        <n-button text @click="handleOpen">
          <template #icon>
            <n-icon :component="NotificationsOutline" size="20" />
          </template>
        </n-button>
      </n-badge>
    </template>

    <div class="notification-panel">
      <div class="notification-header">
        <span class="notification-title">通知中心</span>
        <n-button text size="small" type="primary" @click="handleMarkAllRead" :disabled="notificationStore.unreadCount === 0">
          全部已读
        </n-button>
      </div>

      <n-scrollbar style="max-height: 400px">
        <div v-if="loading && notifications.length === 0" class="notification-loading">
          <n-spin size="small" />
        </div>
        <div v-else-if="notifications.length === 0" class="notification-empty">
          <n-empty description="暂无通知" size="small" />
        </div>
        <div v-else>
          <div
            v-for="item in notifications"
            :key="item.id"
            class="notification-item"
            :class="{ unread: !item.is_read }"
            @click="handleRead(item)"
          >
            <div class="notification-item-icon">
              <n-icon :component="getIcon(item.notification.type)" :color="getColor(item.notification.type)" size="20" />
            </div>
            <div class="notification-item-content">
              <div class="notification-item-title">{{ item.notification.title }}</div>
              <div class="notification-item-text">{{ item.notification.content }}</div>
              <div class="notification-item-time">{{ formatTime(item.created_at) }}</div>
            </div>
          </div>
          <div v-if="hasMore" class="notification-load-more">
            <n-button text size="small" @click="loadMore">加载更多</n-button>
          </div>
          <div v-else-if="notifications.length > 0" class="notification-end">
            暂无更多通知
          </div>
        </div>
      </n-scrollbar>
    </div>
  </n-popover>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { useRouter } from 'vue-router'
import { NPopover, NBadge, NButton, NIcon, NScrollbar, NEmpty, NSpin } from 'naive-ui'
import {
  NotificationsOutline,
  ChatbubbleEllipsesOutline,
  DocumentTextOutline,
  MegaphoneOutline,
  LinkOutline
} from '@vicons/ionicons5'
import { useNotificationStore } from '@/stores/notification'
import { useAuthStore } from '@/stores/auth'
import { getNotifications, markAsRead, markAllAsRead } from '@/api/notification'
import type { Notification } from '@/api/notification'

const router = useRouter()
const notificationStore = useNotificationStore()
const authStore = useAuthStore()

const notifications = ref<Notification[]>([])
const loading = ref(false)
const page = ref(1)
const hasMore = ref(true)
let ws: WebSocket | null = null
let pollTimer: ReturnType<typeof setInterval> | null = null

function getIcon(type: string) {
  switch (type) {
    case 'comment_reply': return ChatbubbleEllipsesOutline
    case 'comment_new': return DocumentTextOutline
    case 'article_push': return MegaphoneOutline
    case 'friend_apply': return LinkOutline
    default: return NotificationsOutline
  }
}

function getColor(type: string) {
  switch (type) {
    case 'comment_reply': return '#3b82f6'
    case 'comment_new': return '#22c55e'
    case 'article_push': return '#f97316'
    case 'friend_apply': return '#a855f7'
    default: return '#6b7280'
  }
}

function getRoute(type: string): string {
  switch (type) {
    case 'comment_reply':
    case 'comment_new':
      return '/admin/comments'
    case 'article_push':
      return '/admin/posts'
    case 'friend_apply':
      return '/admin/friend-links'
    default:
      return '/admin/dashboard'
  }
}

function formatTime(time: string) {
  const now = new Date()
  const date = new Date(time)
  const diff = now.getTime() - date.getTime()
  const minutes = Math.floor(diff / 60000)
  if (minutes < 1) return '刚刚'
  if (minutes < 60) return `${minutes}分钟前`
  const hours = Math.floor(minutes / 60)
  if (hours < 24) return `${hours}小时前`
  const days = Math.floor(hours / 24)
  if (days < 30) return `${days}天前`
  return time.split(' ')[0]
}

async function fetchNotifications(reset = false) {
  if (reset) {
    page.value = 1
    hasMore.value = true
  }
  loading.value = true
  try {
    const res = await getNotifications({ page: page.value, page_size: 20 })
    if (res.data) {
      if (reset) {
        notifications.value = res.data.list
      } else {
        notifications.value.push(...res.data.list)
      }
      hasMore.value = notifications.value.length < res.data.total
    }
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

function loadMore() {
  page.value++
  fetchNotifications()
}

function handleOpen() {
  fetchNotifications(true)
  notificationStore.fetchUnreadCount()
}

async function handleRead(item: Notification) {
  if (!item.is_read) {
    if (item.notification_id > 0) {
      await markAsRead(item.notification_id)
    }
    item.is_read = true
    notificationStore.decrementUnread()
  }
  router.push(getRoute(item.notification.type))
}

async function handleMarkAllRead() {
  await markAllAsRead()
  notifications.value.forEach(n => { n.is_read = true })
  notificationStore.setUnreadCount(0)
}

function connectWebSocket() {
  if (!authStore.isLoggedIn || !authStore.token) return

  const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
  const wsUrl = `${protocol}//${window.location.host}/api/notifications/ws?token=${authStore.token}`

  try {
    ws = new WebSocket(wsUrl)

    ws.onopen = () => {
      notificationStore.setConnected(true)
    }

    ws.onmessage = (event) => {
      try {
        const msg = JSON.parse(event.data)
        if (msg.type === 'notification') {
          notificationStore.incrementUnread()
          notifications.value.unshift({
            id: Date.now(),
            notification_id: msg.data.id || 0,
            user_id: 0,
            is_read: false,
            email_sent: false,
            created_at: msg.data.created_at || new Date().toISOString(),
            notification: {
              id: msg.data.id || 0,
              type: msg.data.type,
              title: msg.data.title,
              content: msg.data.content,
              sender_id: msg.data.sender?.id || null,
              target_type: '',
              target_id: null,
              extra: null,
              status: 1,
              created_at: msg.data.created_at || new Date().toISOString(),
              updated_at: msg.data.created_at || new Date().toISOString(),
              sender: msg.data.sender || null
            }
          })
        } else if (msg.type === 'unread_count') {
          notificationStore.setUnreadCount(msg.data.count)
        }
      } catch {
        // ignore
      }
    }

    ws.onclose = () => {
      notificationStore.setConnected(false)
      if (authStore.isLoggedIn) {
        setTimeout(connectWebSocket, 30000)
      }
    }

    ws.onerror = () => {
      notificationStore.setConnected(false)
    }
  } catch {
    // ignore
  }
}

onMounted(() => {
  if (authStore.isLoggedIn) {
    notificationStore.fetchUnreadCount()
    connectWebSocket()
    pollTimer = setInterval(() => {
      notificationStore.fetchUnreadCount()
    }, 30000)
  }
})

onUnmounted(() => {
  if (ws) {
    ws.close()
    ws = null
  }
  if (pollTimer) {
    clearInterval(pollTimer)
    pollTimer = null
  }
})
</script>

<style scoped>
.notification-panel {
  padding: 0;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 12px 16px;
  border-bottom: 1px solid var(--n-border-color);
}

.notification-title {
  font-size: 16px;
  font-weight: 600;
}

.notification-loading {
  display: flex;
  justify-content: center;
  padding: 40px 0;
}

.notification-empty {
  padding: 40px 0;
}

.notification-item {
  display: flex;
  padding: 12px 16px;
  cursor: pointer;
  transition: background-color 0.2s;
  border-bottom: 1px solid rgba(0, 0, 0, 0.05);
}

.notification-item:hover {
  background-color: rgba(0, 0, 0, 0.03);
}

.notification-item.unread {
  background-color: rgba(59, 130, 246, 0.05);
}

.notification-item.unread:hover {
  background-color: rgba(59, 130, 246, 0.08);
}

.notification-item-icon {
  flex-shrink: 0;
  margin-right: 12px;
  margin-top: 2px;
}

.notification-item-content {
  flex: 1;
  min-width: 0;
}

.notification-item-title {
  font-size: 14px;
  font-weight: 500;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-item-text {
  font-size: 13px;
  color: #6b7280;
  margin-bottom: 4px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.notification-item-time {
  font-size: 12px;
  color: #9ca3af;
}

.notification-load-more,
.notification-end {
  text-align: center;
  padding: 12px;
  font-size: 13px;
  color: #9ca3af;
}
</style>
