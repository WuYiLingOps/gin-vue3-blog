/*
 * 项目名称：blog-frontend
 * 文件名称：imageAsset.ts
 * 创建时间：2026-11-20 10:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：图片管理相关 API 接口定义，包括全站图片盘点和图床间迁移功能。
 */

import { request } from '@/utils/request'

/**
 * 图片使用位置明细
 */
export interface ImageUsageLocation {
  table: string // 所在表名
  field: string // 所在字段（settings 表为配置 key）
  table_label: string // 位置中文说明（如：文章封面、用户头像）
  row_id: number // 所在记录 ID
  kind: string // 字段类型：url/markdown/json_array/settings
}

/**
 * 图片资产接口（盘点结果）
 */
export interface ImageAsset {
  url: string // 图片 URL（库内原样）
  storage_type: 'local' | 'oss' | 'cos' | 'easyimage' | 'external' // 来源存储类型
  locations: ImageUsageLocation[] // 使用位置明细
  ref_count: number // 引用次数
  deletable: boolean // 源文件是否支持 API 删除
  alive?: boolean | null // 文件有效性：true-可访问 false-已失效 null-未检测
}

/**
 * 图片迁移日志接口
 */
export interface ImageMigrateLog {
  id: number
  operator_id?: number // 操作用户ID
  operator_name: string // 操作用户名
  target_storage: string // 目标存储类型
  delete_source: boolean // 是否删除了源文件
  total_count: number // 迁移图片总数
  success_count: number // 成功数
  failed_count: number // 失败数
  updated_rows: number // 数据库替换更新的记录行数
  deleted_files: number // 实际删除的源文件数
  mappings: { old_url: string; new_url: string }[] // 成功映射（兼容历史 {old: new} 对象格式）
  failed_items: { url: string; error: string }[] // 失败明细
  message: string // 任务结束汇总信息
  duration_sec: number // 任务总耗时（秒）
  created_at: string // 任务完成时间
}

/**
 * 迁移任务状态接口
 */
export interface MigrateStatus {
  running: boolean // 是否执行中
  total: number // 总数
  migrated: number // 已完成数
  failed: { url: string; error: string }[] // 失败列表
  success: { old_url: string; new_url: string }[] // 成功映射列表
  finished: boolean // 是否已有任务结束（区分从未运行和已结束）
  message: string // 结束时的汇总信息
  current_url?: string // 正在处理的图片 URL（running 时可见）
  current_step?: string // 当前步骤描述（下载中/上传中/替换数据库/删除源文件）
  elapsed_sec?: number // 任务已进行秒数
}

/**
 * 发起迁移请求参数
 */
export interface MigrateImageRequest {
  urls: string[] // 要迁移的图片 URL 列表
  target_storage: 'local' | 'oss' | 'cos' | 'easyimage' // 目标存储类型
  delete_source?: boolean // 迁移成功后是否删除源文件
}

/**
 * 全站图片盘点
 * 返回全部图片资产（后端 Redis 缓存 10 分钟，切换筛选请用前端本地过滤，不要反复调接口）
 * @param refresh 是否强制刷新缓存重新扫描数据库
 */
export function scanImageAssets(refresh: boolean = false) {
  return request.get<ImageAsset[]>('/admin/images/assets', {
    params: refresh ? { refresh: 'true' } : {}
  })
}

/**
 * 发起图片迁移任务
 */
export function migrateImages(data: MigrateImageRequest) {
  return request.post<null>('/admin/images/assets/migrate', data)
}

/**
 * 获取迁移任务进度
 */
export function getMigrateStatus() {
  return request.get<MigrateStatus>('/admin/images/assets/migrate/status')
}

/**
 * 分页查询图片迁移日志
 */
export function getMigrateLogs(page: number = 1, pageSize: number = 10) {
  return request.get<{ list: ImageMigrateLog[]; total: number; page: number; page_size: number; total_pages: number }>(
    '/admin/images/migrate-logs',
    { params: { page, page_size: pageSize } }
  )
}

/**
 * 删除迁移日志（ids 为空数组时清空全部）
 */
export function deleteMigrateLogs(ids: number[] = []) {
  return request.post<null>('/admin/images/migrate-logs/delete', { ids })
}
