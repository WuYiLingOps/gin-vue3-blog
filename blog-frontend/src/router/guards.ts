// 路由守卫

import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores'
import { getPublicSettings } from '@/api/setting'

// 站点名称缓存，避免每次路由切换都请求
let cachedSiteName = '無以菱'
let siteNameLoaded = false
let siteNameLoading: Promise<void> | null = null

async function ensureSiteName() {
  if (siteNameLoaded) return
  if (siteNameLoading) {
    await siteNameLoading
    return
  }
  siteNameLoading = getPublicSettings()
    .then((res) => {
      if (res.data?.site_name) {
        cachedSiteName = res.data.site_name
      }
      siteNameLoaded = true
    })
    .catch(() => {
      // 保持默认缓存名称
    })
    .finally(() => {
      siteNameLoading = null
    })
  await siteNameLoading
}

export function setupRouterGuards(router: Router) {
  // 全局前置守卫
  router.beforeEach((to, _from, next) => {
    const authStore = useAuthStore()

    // 先用缓存的站点名设置标题，再异步刷新一次
    const baseTitle = cachedSiteName || '情迁阁'
    document.title = to.meta.title ? `${to.meta.title} - ${baseTitle}` : baseTitle
    ensureSiteName().then(() => {
      const latest = cachedSiteName || '情迁阁'
      document.title = to.meta.title ? `${to.meta.title} - ${latest}` : latest
    })

    // 检查是否需要认证
    if (to.meta.requiresAuth && !authStore.isLoggedIn) {
      next({
        name: 'Login',
        query: { redirect: to.fullPath }
      })
      return
    }

    // 检查是否需要管理员权限
    if (to.meta.requiresAdmin && !authStore.isAdmin) {
      next({ name: 'Home' })
      return
    }

    // 如果已登录访问登录页，重定向到首页
    if (to.name === 'Login' && authStore.isLoggedIn) {
      next({ name: 'Home' })
      return
    }

    next()
  })

  // 全局后置钩子
  router.afterEach(() => {
    // 可以在这里添加页面加载完成后的逻辑
  })
}

