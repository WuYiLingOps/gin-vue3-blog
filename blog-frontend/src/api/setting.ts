// 系统配置相关 API

import { request } from '@/utils/request'

export interface SiteSettings {
  site_name?: string
  site_icp?: string
  site_police?: string
  // 社交链接
  social_github?: string
  social_gitee?: string
  social_email?: string
  social_rss?: string
  social_qq?: string
  social_wechat?: string
}

export interface UploadSettings {
  storage_type?: string  // 'local' | 'oss'
  oss_endpoint?: string
  oss_access_key_id?: string
  oss_access_key_secret?: string
  oss_bucket_name?: string
  oss_domain?: string
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

// 获取上传配置（管理员）
export function getUploadSettings() {
  return request.get<UploadSettings>('/settings/upload')
}

// 更新上传配置（管理员）
export function updateUploadSettings(data: Record<string, string>) {
  return request.put('/settings/upload', data)
}

export interface FriendLinkInfo {
  name?: string
  desc?: string
  url?: string
  avatar?: string
  screenshot?: string
  rss?: string
}

// 获取我的友链信息（公开接口）
export function getFriendLinkInfo() {
  return request.get<FriendLinkInfo>('/settings/friendlink-info')
}

// 更新我的友链信息（管理员）
export function updateFriendLinkInfo(data: FriendLinkInfo) {
  return request.put('/settings/friendlink-info', data)
}

