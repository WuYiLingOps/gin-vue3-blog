/*
 * @ProjectName: go-vue3-blog
 * @FileName: constants.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 常量定义，包含文章状态、用户角色、分页等常量
 */

// 常量定义

// 文章状态
export const POST_STATUS = {
  DRAFT: 0,
  PUBLISHED: 1,
  DELETED: -1
}

export const POST_STATUS_TEXT = {
  [POST_STATUS.DRAFT]: '草稿',
  [POST_STATUS.PUBLISHED]: '已发布',
  [POST_STATUS.DELETED]: '已删除'
}

// 用户角色
export const USER_ROLE = {
  USER: 'user',
  ADMIN: 'admin'
}

export const USER_ROLE_TEXT = {
  [USER_ROLE.USER]: '普通用户',
  [USER_ROLE.ADMIN]: '管理员'
}

// 用户状态
export const USER_STATUS = {
  DISABLED: 0,
  NORMAL: 1
}

export const USER_STATUS_TEXT = {
  [USER_STATUS.DISABLED]: '已禁用',
  [USER_STATUS.NORMAL]: '正常'
}

// 评论状态
export const COMMENT_STATUS = {
  HIDDEN: 0,
  NORMAL: 1
}

export const COMMENT_STATUS_TEXT = {
  [COMMENT_STATUS.HIDDEN]: '已隐藏',
  [COMMENT_STATUS.NORMAL]: '正常'
}

// 分页大小选项
export const PAGE_SIZE_OPTIONS = [10, 20, 30, 50, 100]

// 默认分页大小
export const DEFAULT_PAGE_SIZE = 10

// 默认颜色
export const DEFAULT_COLORS = [
  '#2196F3',
  '#4CAF50',
  '#FF9800',
  '#9C27B0',
  '#F44336',
  '#00BCD4',
  '#FFEB3B',
  '#795548'
]

