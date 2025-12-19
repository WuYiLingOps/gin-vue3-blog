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
  
  // 如果已经是完整 URL，检查是否需要转换
  if (url.startsWith('http://') || url.startsWith('https://')) {
    // 如果是当前页面的域名，直接返回
    try {
      const urlObj = new URL(url)
      if (urlObj.host === window.location.host) {
        return url
      }
      // 如果是其他域名，提取路径部分，使用当前页面域名
      const relativePath = urlObj.pathname
      return `${window.location.protocol}//${window.location.host}${relativePath}`
    } catch (e) {
      console.warn('Invalid URL:', url)
      return url
    }
  }
  
  // 如果是 /uploads/ 开头的相对路径，使用当前页面的域名
  if (url.startsWith('/uploads/')) {
    return `${window.location.protocol}//${window.location.host}${url}`
  }
  
  // 其他相对路径，使用配置的基础 URL（通常是API路径）
  const currentBaseUrl = import.meta.env.VITE_API_BASE_URL || ''
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

