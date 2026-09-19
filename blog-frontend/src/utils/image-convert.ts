/*
 * 项目名称：blog-frontend
 * 文件名称：image-convert.ts
 * 创建时间：2026-09-19
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：客户端图片转 WebP 工具。上传前用 Canvas 将 PNG/JPEG 重编码为
 *           WebP 并可选降采样，减小大图体积；编码不受支持或无体积收益时
 *           自动回退原始文件。GIF（可能含动画）不做处理。
 */

export interface ConvertToWebPResult {
  file: File
  /** 是否发生了转换（false 表示回退为原始文件） */
  converted: boolean
  /** 原始大小（字节） */
  originalSize: number
  /** 转换后大小（字节），未转换时与 originalSize 相同 */
  convertedSize: number
}

interface ConvertOptions {
  /** 最大宽度，超过时等比降采样，默认 2560 */
  maxWidth?: number
  /** WebP 编码质量 0-1，默认 0.85 */
  quality?: number
}

/** 判断文件是否适合转换：仅处理 PNG/JPEG 位图，GIF/WebP 等跳过 */
function isConvertible(file: File): boolean {
  return /^image\/(png|jpeg)$/.test(file.type)
}

/**
 * 将 PNG/JPEG 图片转换为 WebP
 * @returns 转换结果（converted 为 false 时 file 为原始文件）
 */
export async function convertToWebP(
  file: File,
  options: ConvertOptions = {}
): Promise<ConvertToWebPResult> {
  const { maxWidth = 2560, quality = 0.85 } = options
  const originalSize = file.size

  // 浏览器不支持 WebP 编码（旧 Safari 的 toBlob 会回退为 PNG）时直接返回原文件
  if (!isConvertible(file) || typeof createImageBitmap !== 'function') {
    return { file, converted: false, originalSize, convertedSize: originalSize }
  }

  try {
    const bitmap = await createImageBitmap(file)
    const scale = Math.min(1, maxWidth / bitmap.width)
    const width = Math.max(1, Math.round(bitmap.width * scale))
    const height = Math.max(1, Math.round(bitmap.height * scale))

    const canvas = document.createElement('canvas')
    canvas.width = width
    canvas.height = height
    const ctx = canvas.getContext('2d')
    if (!ctx) {
      bitmap.close()
      return { file, converted: false, originalSize, convertedSize: originalSize }
    }
    ctx.drawImage(bitmap, 0, 0, width, height)
    bitmap.close()

    const blob = await new Promise<Blob | null>(resolve =>
      canvas.toBlob(resolve, 'image/webp', quality)
    )
    canvas.width = canvas.height = 0

    // 编码失败、类型回退（非 webp）或体积不降反升时，保留原始文件
    if (!blob || blob.type !== 'image/webp' || blob.size >= originalSize) {
      return { file, converted: false, originalSize, convertedSize: originalSize }
    }

    const webpFile = new File([blob], file.name.replace(/\.(png|jpe?g)$/i, '') + '.webp', {
      type: 'image/webp'
    })
    return { file: webpFile, converted: true, originalSize, convertedSize: webpFile.size }
  } catch {
    // 解码失败等异常时保留原始文件
    return { file, converted: false, originalSize, convertedSize: originalSize }
  }
}
