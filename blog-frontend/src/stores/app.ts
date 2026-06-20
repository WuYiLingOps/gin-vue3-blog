/*
 * @ProjectName: go-vue3-blog
 * @FileName: app.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 应用全局状态管理，管理主题、侧边栏、加载状态等全局状态
 */

// 应用全局状态管理

import { defineStore } from 'pinia'
import { ref } from 'vue'

export const useAppStore = defineStore(
  'app',
  () => {
    // 状态
    const theme = ref<'light' | 'dark'>('light')
    const sidebarCollapsed = ref(false)
    const loading = ref(true) // 默认为 true，应用启动时显示加载动画
    const bgImages = ref<string[]>([]) // 全局背景图 URL 数组（从后台设置获取，不持久化）
    const siteName = ref('') // 网站名称（从后台设置获取，不持久化）

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

    // 设置全局背景图数组
    function setBgImages(urls: string[]) {
      bgImages.value = urls
    }

    // 设置网站名称
    function setSiteName(name: string) {
      siteName.value = name
    }

    return {
      theme,
      sidebarCollapsed,
      loading,
      bgImages,
      siteName,
      toggleTheme,
      setTheme,
      toggleSidebar,
      setLoading,
      setBgImages,
      setSiteName
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
