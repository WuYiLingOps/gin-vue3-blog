/*
 * @ProjectName: go-vue3-blog
 * @FileName: auth.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 认证状态管理，管理用户登录、注册、登出等认证相关状态
 */

// 认证状态管理

import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import type { User, LoginForm, RegisterForm } from '@/types/auth'
import { login as loginApi, register as registerApi, getProfile, logout as logoutApi } from '@/api/auth'

export const useAuthStore = defineStore(
  'auth',
  () => {
    // 状态
    const token = ref<string | null>(null)
    const user = ref<User | null>(null)

    // 计算属性
    const isLoggedIn = computed(() => !!token.value)
    const isSuperAdmin = computed(() => user.value?.role === 'super_admin')
    const isAdmin = computed(() => user.value?.role === 'admin' || user.value?.role === 'super_admin')

    function hasRole(roles: User['role'][]): boolean {
      const r = user.value?.role
      if (!r) return false
      return roles.includes(r)
    }

    // 登录
    async function login(form: LoginForm) {
      const res = await loginApi(form)
      console.log(res)
      if (res.data) {
        token.value = res.data.token
        user.value = res.data.user
      }
      return res
    }

    // 注册
    async function register(form: RegisterForm) {
      const res = await registerApi(form)
      return res
    }

    // 登出
    async function logout() {
      try {
        await logoutApi()
      } finally {
        token.value = null
        user.value = null
      }
    }

    // 获取用户信息
    async function fetchUserInfo() {
      const res = await getProfile()
      if (res.data) {
        user.value = res.data
      }
      return res
    }

    return {
      token,
      user,
      isLoggedIn,
      isSuperAdmin,
      isAdmin,
      hasRole,
      login,
      register,
      logout,
      fetchUserInfo
    }
  },
  {
    // 配置持久化
    persist: {
      key: 'blog-auth',
      storage: localStorage,
      pick: ['token', 'user']
    }
  }
)

