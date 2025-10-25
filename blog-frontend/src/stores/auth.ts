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
    const isAdmin = computed(() => user.value?.role === 'admin')

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
      isAdmin,
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

