// 系统配置相关 API

import { request } from '@/utils/request'

export interface AboutSettings {
  about_title?: string
  about_intro?: string
  about_avatar?: string
  about_content?: string  // 富文本内容
  about_skills?: string
  about_email?: string
  about_github?: string
  about_site_intro?: string
}

export interface SiteSettings {
  site_name?: string
  site_icp?: string
  site_police?: string
}

// 获取关于页面配置
export function getAboutSettings() {
  return request.get<AboutSettings>('/api/settings/about')
}

// 更新关于页面配置（管理员）
export function updateAboutSettings(data: Record<string, string>) {
  return request.put('/api/settings/about', data)
}

// 获取公开的网站配置
export function getPublicSettings() {
  return request.get<SiteSettings>('/api/settings/public')
}

// 获取网站配置（管理员）
export function getSiteSettings() {
  return request.get<SiteSettings>('/api/settings/site')
}

// 更新网站配置（管理员）
export function updateSiteSettings(data: Record<string, string>) {
  return request.put('/api/settings/site', data)
}

