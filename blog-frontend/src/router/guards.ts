/*
 * @ProjectName: go-vue3-blog
 * @FileName: guards.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 路由守卫配置，提供路由权限控制和页面标题设置功能
 */

// 路由守卫

import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores'
import { getPublicSettings } from '@/api/setting'

// 站点名称缓存，避免每次路由切换都请求
let cachedSiteName = '菱风叙'
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

    // 角色白名单（更细粒度控制）
    if (to.meta.roles && to.meta.roles.length > 0) {
      if (!authStore.hasRole(to.meta.roles)) {
        next({ name: 'Home' })
        return
      }
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

