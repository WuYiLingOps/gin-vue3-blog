/*
 * 项目名称：blog-frontend
 * 文件名称：setting.ts
 * 创建时间：2026-02-01 19:52:47
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：系统配置相关 API 接口定义，包括网站配置、上传配置、通知配置、友链信息、关于我等功能的配置管理。
 */

import { request } from '@/utils/request'
import type { ApiResponse } from '@/types/common'

/**
 * 网站设置接口
 */
export interface SiteSettings {
  site_name?: string // 网站名称
  site_icp?: string // ICP备案号
  site_police?: string // 公安备案号
  // 社交链接
  social_github?: string // GitHub链接
  social_gitee?: string // Gitee链接
  social_email?: string // 邮箱
  social_qq?: string // QQ号
  social_wechat?: string // 微信号
  social_csdn?: string // CSDN链接
  social_link_order?: string // 社交链接排序顺序，逗号分隔的类型列表，如 "github,gitee,email,csdn,qq,wechat"
  reward_wechat?: string // 微信收款码图片URL
  reward_alipay?: string // 支付宝收款码图片URL
  // 封面设置
  cover_subtitle?: string // 封面副标题（多条用换行分隔，打字机效果轮播）
  cover_bg_image?: string // 封面背景图URL（已废弃，保留用于兼容）
  cover_bg_images?: string // 封面背景图URL数组（JSON字符串格式，最多3张）
  // 网站运行时间
  site_start_date?: string // 网站成立时间（用于计算运行时长）
  // 关于页卡片配置（JSON字符串格式，类型见 @/types/about.ts）
  about_site_tips?: string // 打字词轮播配置
  about_skills?: string // 技能配置（头像浮动标签 + 技能卡）
  about_careers?: string // 职业生涯时间线配置
  about_maxim?: string // 座右铭配置
  about_map?: string // 地理位置配置（JSON字符串，类型见 @/types/about.ts）
  about_intro?: string // 作者介绍卡配置（JSON字符串，类型见 @/types/about.ts）
  about_self_info?: string // 个人信息卡配置（JSON字符串，类型见 @/types/about.ts）
  about_personality?: string // MBTI 性格卡配置（JSON字符串，类型见 @/types/about.ts）
}

/**
 * 上传设置接口
 */
export interface UploadSettings {
  storage_type?: string // 存储类型：'local' | 'oss'
  oss_endpoint?: string // OSS端点地址
  oss_access_key_id?: string // OSS访问密钥ID
  oss_access_key_secret?: string // OSS访问密钥Secret
  oss_bucket_name?: string // OSS存储桶名称
  oss_domain?: string // OSS自定义域名
}

/**
 * 通知设置接口
 */
export interface NotificationSettings {
  notify_admin_on_comment?: string // 评论时通知管理员：'0'表示否，'1'表示是
}

/**
 * 获取公开的网站配置
 * 首屏多个组件会并发请求同一份配置，这里做 3 秒内的请求合并与短缓存，
 * 调用 updateSiteSettings 后自动失效
 */
let publicSettingsCache: Promise<ApiResponse<SiteSettings>> | null = null
let publicSettingsCachedAt = 0

export function getPublicSettings(): Promise<ApiResponse<SiteSettings>> {
  let p = publicSettingsCache
  const now = Date.now()
  if (!p || now - publicSettingsCachedAt > 3000) {
    publicSettingsCachedAt = now
    p = request.get<SiteSettings>('/settings/public').catch(err => {
      // 请求失败时清除缓存，允许下次重试
      if (publicSettingsCache === p) publicSettingsCache = null
      throw err
    })
    publicSettingsCache = p
  }
  return p
}

/**
 * 获取网站配置（管理员）
 * @returns 返回完整的网站配置信息
 */
export function getSiteSettings() {
  return request.get<SiteSettings>('/settings/site')
}

/**
 * 更新网站配置（管理员）
 * @param data 网站配置数据（键值对）
 * @returns 返回更新结果
 */
export function updateSiteSettings(data: Record<string, string>) {
  // 配置变更后使公开配置的短缓存失效，前台下次请求拉取最新值
  publicSettingsCache = null
  publicSettingsCachedAt = 0
  return request.put('/settings/site', data)
}

/**
 * 获取上传配置（管理员）
 * @returns 返回上传配置信息
 */
export function getUploadSettings() {
  return request.get<UploadSettings>('/settings/upload')
}

/**
 * 更新上传配置（管理员）
 * @param data 上传配置数据（键值对）
 * @returns 返回更新结果
 */
export function updateUploadSettings(data: Record<string, string>) {
  return request.put('/settings/upload', data)
}

/**
 * 友链信息接口
 */
export interface FriendLinkInfo {
  name?: string // 网站名称
  desc?: string // 网站描述
  url?: string // 网站URL
  avatar?: string // 网站头像
  screenshot?: string // 网站截图
  rss?: string // RSS订阅地址
}

/**
 * 获取我的友链信息（公开接口）
 * @returns 返回博主的友链信息
 */
export function getFriendLinkInfo() {
  return request.get<FriendLinkInfo>('/settings/friendlink-info')
}

/**
 * 更新我的友链信息（管理员）
 * @param data 友链信息数据
 * @returns 返回更新结果
 */
export function updateFriendLinkInfo(data: FriendLinkInfo) {
  return request.put('/settings/friendlink-info', data)
}

/**
 * 获取通知配置（管理员）
 * @returns 返回通知配置信息
 */
export function getNotificationSettings() {
  return request.get<NotificationSettings>('/settings/notification')
}

/**
 * 更新通知配置（管理员）
 * @param data 通知配置数据（键值对）
 * @returns 返回更新结果
 */
export function updateNotificationSettings(data: Record<string, string>) {
  return request.put('/settings/notification', data)
}

/**
 * 获取关于我信息（管理员）
 * @returns 返回关于我内容
 */
export function getAboutInfo() {
  return request.get<{ content: string }>('/admin/about')
}

/**
 * 更新关于我信息（管理员）
 * @param content 关于我内容
 * @returns 返回更新结果
 */
export function updateAboutInfo(content: string) {
  return request.put('/admin/about', { content })
}

/**
 * 公开接口：获取关于我信息
 * @returns 返回关于我内容
 */
export function getPublicAboutInfo() {
  return request.get<{ content: string }>('/blog/about')
}

/**
 * 显示设置接口
 */
export interface DisplaySettings {
  home_page_size?: string // 首页文章列表每页数量
  recent_posts_limit?: string // 最新发布文章卡片显示数量
}

/**
 * 公开接口：获取显示配置（首页使用）
 * @returns 返回显示配置
 */
export function getPublicDisplaySettings() {
  return request.get<DisplaySettings>('/settings/public/display')
}

/**
 * 获取显示配置（管理员）
 * @returns 返回显示配置
 */
export function getDisplaySettings() {
  return request.get<DisplaySettings>('/settings/display')
}

/**
 * 更新显示配置（管理员）
 * @param data 显示配置数据（键值对）
 * @returns 返回更新结果
 */
export function updateDisplaySettings(data: Record<string, string>) {
  return request.put('/settings/display', data)
}
