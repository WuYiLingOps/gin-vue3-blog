import { defineStore } from 'pinia'
import { ref } from 'vue'
import { getUnreadCount } from '@/api/notification'

export const useNotificationStore = defineStore('notification', () => {
  const unreadCount = ref(0)
  const isConnected = ref(false)

  async function fetchUnreadCount() {
    try {
      const res = await getUnreadCount()
      if (res.data) {
        unreadCount.value = res.data.count
      }
    } catch {
      // ignore
    }
  }

  function setUnreadCount(count: number) {
    unreadCount.value = count
  }

  function incrementUnread() {
    unreadCount.value++
  }

  function decrementUnread() {
    if (unreadCount.value > 0) {
      unreadCount.value--
    }
  }

  function setConnected(value: boolean) {
    isConnected.value = value
  }

  return {
    unreadCount,
    isConnected,
    fetchUnreadCount,
    setUnreadCount,
    incrementUnread,
    decrementUnread,
    setConnected
  }
})
