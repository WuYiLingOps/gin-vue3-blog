// 路由守卫

import type { Router } from 'vue-router'
import { useAuthStore } from '@/stores'

export function setupRouterGuards(router: Router) {
  // 全局前置守卫
  router.beforeEach((to, _from, next) => {
    const authStore = useAuthStore()

    // 设置页面标题
    document.title = to.meta.title ? `${to.meta.title} - 我的博客` : '我的博客'

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

