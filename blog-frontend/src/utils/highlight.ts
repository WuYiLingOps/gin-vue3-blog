/*
 * @ProjectName: go-vue3-blog
 * @FileName: highlight.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 文本高亮工具，提供关键词高亮和Markdown文本处理功能
 */

/**
 * 高亮匹配的关键词
 * @param text 原始文本
 * @param keyword 搜索关键词
 * @returns 带有高亮标记的HTML字符串
 */
export function highlightKeyword(text: string, keyword: string): string {
  if (!text || !keyword) {
    return text || ''
  }

  // 转义特殊字符，防止正则表达式错误
  const escapedKeyword = keyword.replace(/[.*+?^${}()|[\]\\]/g, '\\$&')
  
  // 创建全局不区分大小写的正则表达式
  const regex = new RegExp(`(${escapedKeyword})`, 'gi')
  
  // 替换匹配的关键词为带有mark标签的文本
  return text.replace(regex, '<mark class="search-highlight">$1</mark>')
}

/**
 * 清理 Markdown 语法标记
 * @param text Markdown 文本
 * @returns 纯文本
 */
export function stripMarkdown(text: string): string {
  if (!text) return ''
  
  return text
    // 提取代码块内容（保留代码内容，只移除```标记）
    .replace(/```[a-z]*\n?([\s\S]*?)```/g, ' $1 ')
    // 提取行内代码内容（保留代码内容，只移除`标记）
    .replace(/`([^`]+)`/g, ' $1 ')
    // 移除标题标记
    .replace(/#{1,6}\s+/g, '')
    // 移除粗体和斜体
    .replace(/\*\*\*(.+?)\*\*\*/g, '$1')
    .replace(/\*\*(.+?)\*\*/g, '$1')
    .replace(/\*(.+?)\*/g, '$1')
    .replace(/___(.+?)___/g, '$1')
    .replace(/__(.+?)__/g, '$1')
    .replace(/_(.+?)_/g, '$1')
    // 移除删除线
    .replace(/~~(.+?)~~/g, '$1')
    // 移除链接，保留文本
    .replace(/\[([^\]]+)\]\([^\)]+\)/g, '$1')
    // 移除图片
    .replace(/!\[([^\]]*)\]\([^\)]+\)/g, '')
    // 移除引用标记
    .replace(/^\s*>\s+/gm, '')
    // 移除列表标记
    .replace(/^\s*[-*+]\s+/gm, '')
    .replace(/^\s*\d+\.\s+/gm, '')
    // 移除水平线
    .replace(/^\s*[-*_]{3,}\s*$/gm, '')
    // 移除多余的空白字符
    .replace(/\n{2,}/g, ' ')
    .replace(/\s{2,}/g, ' ')
    .trim()
}

/**
 * 截取包含关键词的文本片段
 * @param text 原始文本
 * @param keyword 搜索关键词
 * @param maxLength 最大长度
 * @returns 截取后的文本
 */
export function extractHighlightSnippet(text: string, keyword: string, maxLength: number = 150): string {
  if (!text || !keyword) {
    return text?.slice(0, maxLength) || ''
  }

  // 先清理 Markdown 语法
  const cleanText = stripMarkdown(text)
  
  const lowerText = cleanText.toLowerCase()
  const lowerKeyword = keyword.toLowerCase()
  const index = lowerText.indexOf(lowerKeyword)

  if (index === -1) {
    // 如果没找到关键词，返回开头部分
    return cleanText.slice(0, maxLength)
  }

  // 计算截取的起始位置（关键词前后各留一些字符）
  const halfLength = Math.floor(maxLength / 2)
  const start = Math.max(0, index - halfLength)
  const end = Math.min(cleanText.length, index + keyword.length + halfLength)

  let snippet = cleanText.slice(start, end)

  // 添加省略号
  if (start > 0) {
    snippet = '...' + snippet
  }
  if (end < cleanText.length) {
    snippet = snippet + '...'
  }

  return snippet
}

