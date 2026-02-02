/*
 * @ProjectName: go-vue3-blog
 * @FileName: validator.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 表单验证工具，提供邮箱、用户名、密码、URL等验证功能
 */

// 表单验证工具

// 验证邮箱
export function validateEmail(email: string): boolean {
  const reg = /^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$/
  return reg.test(email)
}

// 验证用户名（3-20个字符，字母、数字、下划线）
export function validateUsername(username: string): boolean {
  const reg = /^[a-zA-Z0-9_]{3,20}$/
  return reg.test(username)
}

// 验证密码（至少6个字符）
export function validatePassword(password: string): boolean {
  return password.length >= 6
}

// 验证 URL
export function validateUrl(url: string): boolean {
  const reg =
    /^(https?:\/\/)?([\da-z.-]+)\.([a-z.]{2,6})([/\w .-]*)*\/?$/
  return reg.test(url)
}

// 验证手机号
export function validatePhone(phone: string): boolean {
  const reg = /^1[3-9]\d{9}$/
  return reg.test(phone)
}

