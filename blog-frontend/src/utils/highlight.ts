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

  const lowerText = text.toLowerCase()
  const lowerKeyword = keyword.toLowerCase()
  const index = lowerText.indexOf(lowerKeyword)

  if (index === -1) {
    // 如果没找到关键词，返回开头部分
    return text.slice(0, maxLength)
  }

  // 计算截取的起始位置（关键词前后各留一些字符）
  const halfLength = Math.floor(maxLength / 2)
  const start = Math.max(0, index - halfLength)
  const end = Math.min(text.length, index + keyword.length + halfLength)

  let snippet = text.slice(start, end)

  // 添加省略号
  if (start > 0) {
    snippet = '...' + snippet
  }
  if (end < text.length) {
    snippet = snippet + '...'
  }

  return snippet
}

