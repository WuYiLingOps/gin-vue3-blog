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
    const pickedBgUrl = ref('') // 当前选中的背景图 URL（缓存后供 index.html 预加载）
    const siteName = ref('') // 网站名称（从后台设置获取，不持久化）

    const BG_CACHE_KEY = 'blog-bg-cache' // 与 index.html 内联脚本保持一致

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
    // keepPick: 配置刷新时保持当前选中的壁纸不变（避免回访用户重复下载新图）
    function setBgImages(urls: string[], options?: { keepPick?: boolean }) {
      bgImages.value = urls
      const keep = !!options?.keepPick && !!pickedBgUrl.value && urls.includes(pickedBgUrl.value)
      if (keep) {
        // 保持当前选中
      } else if (urls.length) {
        pickedBgUrl.value = urls[Math.floor(Math.random() * urls.length)]
      } else {
        pickedBgUrl.value = ''
      }
      // 持久化到 localStorage，供下次启动与 index.html 预加载使用
      try {
        if (urls.length) {
          localStorage.setItem(
            BG_CACHE_KEY,
            JSON.stringify({ urls: bgImages.value, picked: pickedBgUrl.value })
          )
        } else {
          localStorage.removeItem(BG_CACHE_KEY)
        }
      } catch {
        /* localStorage 不可用时静默降级 */
      }
    }

    // 启动时读取缓存的背景图列表与选中项，让壁纸下载与配置请求并行
    function loadCachedBgImages() {
      try {
        const raw = localStorage.getItem(BG_CACHE_KEY)
        if (!raw) return
        const cached = JSON.parse(raw)
        if (Array.isArray(cached.urls) && cached.urls.length) {
          bgImages.value = cached.urls
          pickedBgUrl.value =
            typeof cached.picked === 'string' && cached.urls.includes(cached.picked)
              ? cached.picked
              : cached.urls[0]
        }
      } catch {
        /* 缓存损坏时静默忽略 */
      }
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
      pickedBgUrl,
      siteName,
      toggleTheme,
      setTheme,
      toggleSidebar,
      setLoading,
      setBgImages,
      loadCachedBgImages,
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
