/*
 * @ProjectName: go-vue3-blog
 * @FileName: storage.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 本地存储工具，提供Token和用户信息的本地存储操作
 */

// 本地存储工具

const TOKEN_KEY = 'blog_token'
const USER_KEY = 'blog_user'

export const storage = {
  // Token 操作
  getToken(): string | null {
    return localStorage.getItem(TOKEN_KEY)
  },

  setToken(token: string): void {
    localStorage.setItem(TOKEN_KEY, token)
  },

  removeToken(): void {
    localStorage.removeItem(TOKEN_KEY)
  },

  // 用户信息操作
  getUser(): any {
    const user = localStorage.getItem(USER_KEY)
    return user ? JSON.parse(user) : null
  },

  setUser(user: any): void {
    localStorage.setItem(USER_KEY, JSON.stringify(user))
  },

  removeUser(): void {
    localStorage.removeItem(USER_KEY)
  },

  // 清除所有
  clear(): void {
    localStorage.removeItem(TOKEN_KEY)
    localStorage.removeItem(USER_KEY)
  },

  // 通用方法
  get(key: string): any {
    const value = localStorage.getItem(key)
    try {
      return value ? JSON.parse(value) : null
    } catch {
      return value
    }
  },

  set(key: string, value: any): void {
    localStorage.setItem(key, typeof value === 'string' ? value : JSON.stringify(value))
  },

  remove(key: string): void {
    localStorage.removeItem(key)
  }
}

