/**
 * 项目名称：blog-frontend
 * 文件名称：useTypewriter.ts
 * 创建时间：2026-04-13 17:24:45
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：打字机效果 composable，支持多条文本循环打字、删除、切换，
 *          支持响应式数据源（Ref/Computed），无外部依赖
 */

import { ref, watch, onBeforeUnmount, isRef, type Ref } from 'vue'

interface TypewriterOptions {
  /** 每个字符打字间隔（毫秒） */
  typeSpeed?: number
  /** 每个字符删除间隔（毫秒） */
  deleteSpeed?: number
  /** 打字完成后暂停时间（毫秒） */
  pauseTime?: number
  /** 首次开始前的延迟（毫秒） */
  startDelay?: number
}

/**
 * useTypewriter 提供打字机效果的响应式文本
 * @param source 要循环展示的文本数组（支持 Ref<string[]> 或普通数组）
 * @param options 速度和延迟配置
 * @returns displayedText — 当前显示的文本（响应式 ref）
 */
export function useTypewriter(source: string[] | Ref<string[]>, options: TypewriterOptions = {}) {
  const { typeSpeed = 120, deleteSpeed = 60, pauseTime = 2000, startDelay = 500 } = options

  const displayedText = ref('')
  let currentIndex = 0
  let charIndex = 0
  let isDeleting = false
  let timer: ReturnType<typeof setTimeout> | null = null

  // 获取当前文本列表
  function getStrings(): string[] {
    return isRef(source) ? source.value : source
  }

  function stop() {
    if (timer) {
      clearTimeout(timer)
      timer = null
    }
  }

  function tick() {
    const strings = getStrings()
    if (!strings.length) return

    // 防止 index 越界（列表可能动态变化）
    if (currentIndex >= strings.length) {
      currentIndex = 0
    }

    const current = strings[currentIndex]

    if (!isDeleting) {
      charIndex++
      displayedText.value = current.slice(0, charIndex)

      if (charIndex === current.length) {
        timer = setTimeout(() => {
          isDeleting = true
          tick()
        }, pauseTime)
        return
      }
      timer = setTimeout(tick, typeSpeed)
    } else {
      charIndex--
      displayedText.value = current.slice(0, charIndex)

      if (charIndex === 0) {
        isDeleting = false
        currentIndex = (currentIndex + 1) % strings.length
        timer = setTimeout(tick, typeSpeed)
        return
      }
      timer = setTimeout(tick, deleteSpeed)
    }
  }

  function start() {
    stop()
    currentIndex = 0
    charIndex = 0
    isDeleting = false
    displayedText.value = ''
    timer = setTimeout(tick, startDelay)
  }

  // 监听数据源变化，重新开始打字
  if (isRef(source)) {
    watch(
      source,
      newVal => {
        if (newVal.length > 0) {
          start()
        }
      },
      { immediate: true }
    )
  } else if (source.length > 0) {
    // 静态数组，直接启动
    start()
  }

  onBeforeUnmount(() => {
    stop()
  })

  return { displayedText }
}
