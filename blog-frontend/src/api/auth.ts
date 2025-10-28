// 认证相关 API

import { request } from '@/utils/request'
import type { LoginForm, RegisterForm, LoginResponse, User, ProfileForm, PasswordForm, CaptchaResponse } from '@/types/auth'

// 获取验证码
export function getCaptcha() {
  return request.get<CaptchaResponse>('/captcha')
}

// 用户注册
export function register(data: RegisterForm) {
  return request.post<User>('/auth/register', data)
}

// 用户登录
export function login(data: LoginForm) {
  return request.post<LoginResponse>('/auth/login', data)
}

// 用户登出
export function logout() {
  return request.post('/auth/logout')
}

// 获取用户信息
export function getProfile() {
  return request.get<User>('/auth/profile')
}

// 更新用户信息
export function updateProfile(data: ProfileForm) {
  return request.put<User>('/auth/profile', data)
}

// 修改密码
export function updatePassword(data: PasswordForm) {
  return request.put('/auth/password', data)
}

// 修改邮箱
export function updateEmail(data: { new_email: string }) {
  return request.put('/auth/email', data)
}

// 获取邮箱修改信息
export function getEmailChangeInfo() {
  return request.get<{
    change_count: number
    remaining_times: number
    can_change: boolean
  }>('/auth/email-change-info')
}

// 刷新 Token
export function refreshToken() {
  return request.post<{ token: string }>('/auth/refresh')
}

// 忘记密码 - 发送验证码
export function forgotPassword(data: { email: string }) {
  return request.post('/auth/forgot-password', data)
}

// 重置密码
export function resetPassword(data: {
  email: string
  code: string
  new_password: string
}) {
  return request.post('/auth/reset-password', data)
}

