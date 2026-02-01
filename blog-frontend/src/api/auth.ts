/*
 * 项目名称：blog-frontend
 * 文件名称：auth.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：用户认证相关 API 接口定义，包括登录、注册、登出、用户信息管理、密码重置、邮箱修改等功能。
 */

import { request } from '@/utils/request'
import type { LoginForm, RegisterForm, LoginResponse, User, ProfileForm, PasswordForm, CaptchaResponse } from '@/types/auth'

/**
 * 获取图形验证码
 * @returns 返回验证码图片和验证码ID
 */
export function getCaptcha() {
  return request.get<CaptchaResponse>('/captcha')
}

/**
 * 发送注册验证码邮件
 * @param data 邮箱数据
 * @param data.email 用户邮箱地址
 * @returns 返回发送结果
 */
export function sendRegisterCode(data: { email: string }) {
  return request.post('/auth/send-register-code', data)
}

/**
 * 用户注册
 * @param data 注册表单数据
 * @returns 返回注册后的用户信息
 */
export function register(data: RegisterForm) {
  return request.post<User>('/auth/register', data)
}

/**
 * 用户登录
 * @param data 登录表单数据（用户名/邮箱、密码、验证码等）
 * @returns 返回登录响应，包含用户信息和访问令牌
 */
export function login(data: LoginForm) {
  return request.post<LoginResponse>('/auth/login', data)
}

/**
 * 用户登出
 * @returns 返回登出结果
 */
export function logout() {
  return request.post('/auth/logout')
}

/**
 * 获取当前登录用户信息
 * @returns 返回用户详细信息
 */
export function getProfile() {
  return request.get<User>('/auth/profile')
}

/**
 * 更新用户信息
 * @param data 用户信息表单数据
 * @returns 返回更新后的用户信息
 */
export function updateProfile(data: ProfileForm) {
  return request.put<User>('/auth/profile', data)
}

/**
 * 修改用户密码
 * @param data 密码表单数据（旧密码、新密码等）
 * @returns 返回修改结果
 */
export function updatePassword(data: PasswordForm) {
  return request.put('/auth/password', data)
}

/**
 * 修改用户邮箱
 * @param data 邮箱修改数据
 * @param data.new_email 新邮箱地址
 * @returns 返回修改结果
 */
export function updateEmail(data: { new_email: string }) {
  return request.put('/auth/email', data)
}

/**
 * 获取邮箱修改信息（包括修改次数、剩余次数等）
 * @returns 返回邮箱修改相关信息
 */
export function getEmailChangeInfo() {
  return request.get<{
    change_count: number        // 已修改次数
    remaining_times: number     // 剩余可修改次数
    can_change: boolean         // 是否可以修改
  }>('/auth/email-change-info')
}

/**
 * 刷新访问令牌
 * @returns 返回新的访问令牌
 */
export function refreshToken() {
  return request.post<{ token: string }>('/auth/refresh')
}

/**
 * 忘记密码 - 发送重置密码验证码邮件
 * @param data 邮箱数据
 * @param data.email 用户邮箱地址
 * @returns 返回发送结果
 */
export function forgotPassword(data: { email: string }) {
  return request.post('/auth/forgot-password', data)
}

/**
 * 重置密码
 * @param data 重置密码数据
 * @param data.email 用户邮箱地址
 * @param data.code 验证码
 * @param data.new_password 新密码
 * @returns 返回重置结果
 */
export function resetPassword(data: {
  email: string
  code: string
  new_password: string
}) {
  return request.post('/auth/reset-password', data)
}

