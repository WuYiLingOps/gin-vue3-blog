// 格式化工具

import dayjs from 'dayjs'
import relativeTime from 'dayjs/plugin/relativeTime'
import 'dayjs/locale/zh-cn'

dayjs.extend(relativeTime)
dayjs.locale('zh-cn')

// 格式化日期
export function formatDate(date: string | Date, format = 'YYYY-MM-DD HH:mm:ss'): string {
  return dayjs(date).format(format)
}

// 相对时间
export function formatRelativeTime(date: string | Date): string {
  return dayjs(date).fromNow()
}

// 格式化数字
export function formatNumber(num: number): string {
  if (num >= 10000) {
    return (num / 10000).toFixed(1) + 'w'
  } else if (num >= 1000) {
    return (num / 1000).toFixed(1) + 'k'
  }
  return num.toString()
}

// 截取文本
export function truncate(text: string, length: number, suffix = '...'): string {
  if (text.length <= length) {
    return text
  }
  return text.substring(0, length) + suffix
}

// 移除 HTML 标签
export function stripHtml(html: string): string {
  const div = document.createElement('div')
  div.innerHTML = html
  return div.textContent || div.innerText || ''
}

