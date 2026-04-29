/*
 * 项目名称：blog-frontend
 * 文件名称：diff.ts
 * 创建时间：2026-04-29 15:00:00
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：文本差异对比工具，用于生成可视化的文本差异
 */

import DiffMatchPatch from 'diff-match-patch'

export interface DiffResult {
  type: 'equal' | 'delete' | 'insert'
  text: string
}

/**
 * 计算两个文本之间的差异
 * @param oldText 原文本
 * @param newText 新文本
 * @returns 差异结果数组
 */
export function computeDiff(oldText: string, newText: string): DiffResult[] {
  const dmp = new DiffMatchPatch()
  const diffs = dmp.diff_main(oldText, newText)
  dmp.diff_cleanupSemantic(diffs)

  return diffs.map(([type, text]) => ({
    type: type === 0 ? 'equal' : type === -1 ? 'delete' : 'insert',
    text
  }))
}

/**
 * 判断是否为长文本（需要使用 diff 对比）
 * @param text 文本内容
 * @returns 是否为长文本
 */
export function isLongText(text: string): boolean {
  return text.length > 200
}

/**
 * 获取差异的上下文（前后各 N 个字符）
 * @param diffs 差异结果
 * @param contextSize 上下文大小
 * @returns 带上下文的差异结果
 */
export function getDiffWithContext(diffs: DiffResult[], contextSize: number = 50): DiffResult[] {
  const result: DiffResult[] = []
  let hasChanges = false

  for (let i = 0; i < diffs.length; i++) {
    const diff = diffs[i]

    if (diff.type !== 'equal') {
      hasChanges = true
      // 添加前面的上下文
      if (i > 0 && diffs[i - 1].type === 'equal') {
        const prevText = diffs[i - 1].text
        if (prevText.length > contextSize) {
          result.push({
            type: 'equal',
            text: '...' + prevText.slice(-contextSize)
          })
        } else {
          result.push(diffs[i - 1])
        }
      }

      // 添加当前差异
      result.push(diff)

      // 添加后面的上下文
      if (i < diffs.length - 1 && diffs[i + 1].type === 'equal') {
        const nextText = diffs[i + 1].text
        if (nextText.length > contextSize) {
          result.push({
            type: 'equal',
            text: nextText.slice(0, contextSize) + '...'
          })
        } else {
          result.push(diffs[i + 1])
        }
        i++ // 跳过下一个 equal
      }
    }
  }

  return hasChanges ? result : diffs
}
