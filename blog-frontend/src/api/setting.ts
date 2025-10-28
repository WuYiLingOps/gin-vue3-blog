// 系统配置相关 API

import { request } from '@/utils/request'

export interface SiteSettings {
  site_name?: string
  site_icp?: string
  site_police?: string
}

// 获取公开的网站配置
export function getPublicSettings() {
  return request.get<SiteSettings>('/settings/public')
}

// 获取网站配置（管理员）
export function getSiteSettings() {
  return request.get<SiteSettings>('/settings/site')
}

// 更新网站配置（管理员）
export function updateSiteSettings(data: Record<string, string>) {
  return request.put('/settings/site', data)
}

