// 应用全局状态管理

import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore(
  'app',
  () => {
    // 状态
    const theme = ref<'light' | 'dark'>('light')
    const sidebarCollapsed = ref(false)
    const loading = ref(false)

    // 切换主题
    function toggleTheme() {
      theme.value = theme.value === 'light' ? 'dark' : 'light'
    }

    // 设置主题
    function setTheme(value: 'light' | 'dark') {
      theme.value = value
    }

    // 切换侧边栏
    function toggleSidebar() {
      sidebarCollapsed.value = !sidebarCollapsed.value
    }

    // 设置加载状态
    function setLoading(value: boolean) {
      loading.value = value
    }

    return {
      theme,
      sidebarCollapsed,
      loading,
      toggleTheme,
      setTheme,
      toggleSidebar,
      setLoading
    }
  },
  {
    // 配置持久化
    persist: {
      key: 'blog-app',
      storage: localStorage,
      pick: ['theme', 'sidebarCollapsed']
    }
  }
)

