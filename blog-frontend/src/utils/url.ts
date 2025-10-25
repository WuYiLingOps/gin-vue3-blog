/**
 * URL 处理工具函数
 */

/**
 * 标准化图片 URL
 * 用于处理后端返回的可能包含旧域名的 URL
 * @param url 原始 URL
 * @returns 标准化后的完整 URL
 */
export function normalizeImageUrl(url: string): string {
  if (!url) return ''
  
  // 如果已经是当前环境的完整 URL，直接返回
  const currentBaseUrl = import.meta.env.VITE_API_BASE_URL || ''
  if (url.startsWith(currentBaseUrl)) {
    return url
  }
  
  // 如果是其他域名的完整 URL（如 http://localhost:8080），提取相对路径
  if (url.startsWith('http://') || url.startsWith('https://')) {
    try {
      const urlObj = new URL(url)
      // 提取路径部分（如 /uploads/avatars/xxx.jpg）
      const relativePath = urlObj.pathname
      return currentBaseUrl + relativePath
    } catch (e) {
      console.warn('Invalid URL:', url)
      return url
    }
  }
  
  // 如果是相对路径，拼接当前环境的基础 URL
  return currentBaseUrl + url
}

/**
 * 批量标准化图片 URL 数组
 * @param urls URL 数组或 JSON 字符串
 * @returns 标准化后的 URL 数组
 */
export function normalizeImageUrls(urls: string | string[]): string[] {
  if (!urls) return []
  
  // 如果是 JSON 字符串，先解析
  let urlArray: string[] = []
  if (typeof urls === 'string') {
    try {
      urlArray = JSON.parse(urls)
    } catch {
      urlArray = [urls]
    }
  } else {
    urlArray = urls
  }
  
  return urlArray.map(normalizeImageUrl)
}

