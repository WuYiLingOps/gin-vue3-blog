<!--
  项目名称：blog-frontend
  文件名称：CommentContent.vue
  创建时间：2026-02-01 20:03:19

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：评论内容渲染组件，用于将Markdown格式的评论内容渲染为HTML，支持代码高亮、代码块复制功能，自动处理代码块滚动位置。
-->
<template>
  <div class="comment-content" ref="previewRef">
    <!-- 使用Markdown预览组件，支持代码高亮和图片 -->
    <v-md-preview :text="content" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { useMessage } from 'naive-ui'
import { VMdPreview } from '@/plugins/v-md-editor'

interface Props {
  content: string
}

const props = defineProps<Props>()
const message = useMessage()
const previewRef = ref<HTMLElement>()

const ALERT_MAP: Record<string, { title: string; icon: string }> = {
  note: { title: '注意', icon: '📝' },
  tip: { title: '技巧提示', icon: '💡' },
  important: { title: '重要', icon: '❗' },
  warning: { title: '注意事项', icon: '⚠️' },
  caution: { title: '危险', icon: '🚨' },
  danger: { title: '严重警告', icon: '⛔' },
  info: { title: '补充说明', icon: 'ℹ️' },
  example: { title: '示例说明', icon: '📋' },
  question: { title: '提出问题', icon: '❓' },
  quote: { title: '引用内容', icon: '💬' }
}

const ALERT_PATTERN = /^\[!(note|tip|important|warning|caution|danger|info|example|question|quote)\]\s*/i

function processAlerts() {
  if (!previewRef.value) return

  const blockquotes = previewRef.value.querySelectorAll('blockquote')
  blockquotes.forEach(blockquote => {
    if ((blockquote as HTMLElement).classList.contains('markdown-alert')) return

    const firstP = blockquote.querySelector('p')
    if (!firstP) return

    const text = firstP.textContent || ''
    const match = text.match(ALERT_PATTERN)
    if (!match) return

    const type = match[1].toLowerCase()
    const alertInfo = ALERT_MAP[type]
    if (!alertInfo) return

    blockquote.classList.add('markdown-alert', `markdown-alert-${type}`)

    const titleEl = document.createElement('p')
    titleEl.className = 'markdown-alert-title'
    titleEl.setAttribute('dir', 'auto')
    titleEl.textContent = `${alertInfo.icon} ${alertInfo.title}`
    blockquote.insertBefore(titleEl, firstP)

    firstP.textContent = text.replace(ALERT_PATTERN, '').trim()
    if (!firstP.textContent) {
      blockquote.removeChild(firstP)
    }
  })
}

// 添加复制按钮到代码块
function addCopyButtons() {
  if (!previewRef.value) return

  const codeBlocks = previewRef.value.querySelectorAll('pre code')

  codeBlocks.forEach(codeBlock => {
    const pre = codeBlock.parentElement
    if (!pre || pre.querySelector('.copy-code-btn')) return

    const button = document.createElement('button')
    button.className = 'copy-code-btn'
    button.textContent = '复制'
    button.onclick = () => {
      const code = codeBlock.textContent || ''
      navigator.clipboard
        .writeText(code)
        .then(() => {
          button.textContent = '已复制!'
          message.success('代码已复制到剪贴板')
          setTimeout(() => {
            button.textContent = '复制'
          }, 2000)
        })
        .catch(() => {
          message.error('复制失败，请手动复制')
        })
    }

    pre.style.position = 'relative'
    pre.appendChild(button)
  })
}

// 修复代码块间距问题
function fixCodeBlockSpacing() {
  if (!previewRef.value) return

  const codeBlocks = previewRef.value.querySelectorAll('pre code')

  codeBlocks.forEach(codeBlock => {
    // 移除所有 token 和 span 的 margin 和 padding
    const allElements = codeBlock.querySelectorAll('*')
    allElements.forEach(element => {
      const el = element as HTMLElement
      // 强制设置所有样式属性
      el.style.margin = '0'
      el.style.padding = '0'
      el.style.wordSpacing = 'normal'
      el.style.letterSpacing = 'normal'
      el.style.wordWrap = 'normal'
      el.style.wordBreak = 'normal'
      el.style.whiteSpace = 'inherit'
      el.style.display = 'inline'
      el.style.lineHeight = 'inherit'
      el.style.verticalAlign = 'baseline'
      el.style.fontSize = 'inherit'
      el.style.border = 'none'
      el.style.background = 'transparent'
      el.style.boxSizing = 'border-box'
    })

    // 确保 code 元素本身没有额外的样式
    const codeElement = codeBlock as HTMLElement
    codeElement.style.wordSpacing = 'normal'
    codeElement.style.letterSpacing = 'normal'
    codeElement.style.wordWrap = 'normal'
    codeElement.style.wordBreak = 'normal'
    codeElement.style.margin = '0'
    codeElement.style.padding = '0'
    codeElement.style.whiteSpace = 'pre'
    codeElement.style.fontVariantLigatures = 'none'
    codeElement.style.textRendering = 'auto'
    // 关键：使用 max-content 确保内容不被压缩
    codeElement.style.width = 'max-content'
    codeElement.style.minWidth = '100%'

    // 确保 pre 元素也没有额外的样式
    const preElement = codeBlock.parentElement as HTMLElement
    if (preElement) {
      preElement.style.wordSpacing = 'normal'
      preElement.style.letterSpacing = 'normal'
      preElement.style.wordWrap = 'normal'
      preElement.style.wordBreak = 'normal'
      preElement.style.whiteSpace = 'pre'
      preElement.style.overflowX = 'auto'
      preElement.style.overflowY = 'hidden'
      preElement.style.width = '100%'
      // 移除可能存在的 max-width 限制
      preElement.style.maxWidth = 'none'
    }
  })
}

// 处理图片点击放大
function handleImageClick() {
  if (!previewRef.value) return

  const images = previewRef.value.querySelectorAll('img')
  images.forEach(imgElement => {
    const img = imgElement as HTMLImageElement
    if (img.hasAttribute('data-lightbox')) return

    img.setAttribute('data-lightbox', 'true')
    img.style.cursor = 'pointer'
    img.onclick = () => {
      // 创建图片预览模态框
      const modal = document.createElement('div')
      modal.style.cssText = `
        position: fixed;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background: rgba(0, 0, 0, 0.8);
        display: flex;
        align-items: center;
        justify-content: center;
        z-index: 9999;
        cursor: pointer;
      `

      const imgClone = img.cloneNode(true) as HTMLImageElement
      imgClone.style.cssText = `
        max-width: 90%;
        max-height: 90%;
        object-fit: contain;
      `

      modal.appendChild(imgClone)
      document.body.appendChild(modal)

      modal.onclick = () => {
        document.body.removeChild(modal)
      }
    }
  })
}

onMounted(() => {
  nextTick(() => {
    addCopyButtons()
    processAlerts()
    handleImageClick()
    fixCodeBlockSpacing()

    // 使用 MutationObserver 监听代码块变化，确保间距修复持续生效
    if (previewRef.value) {
      const observer = new MutationObserver(() => {
        fixCodeBlockSpacing()
      })

      observer.observe(previewRef.value, {
        childList: true,
        subtree: true,
        attributes: true,
        attributeFilter: ['class', 'style']
      })

      // 在组件卸载时断开观察
      return () => {
        observer.disconnect()
      }
    }
  })
})

watch(
  () => props.content,
  () => {
    nextTick(() => {
      addCopyButtons()
      processAlerts()
      handleImageClick()
      // 延迟执行，确保 Prism.js 完成渲染
      setTimeout(() => {
        fixCodeBlockSpacing()
      }, 100)
    })
  }
)
</script>

<style scoped>
.comment-content {
  width: 100%;
  max-width: 100%;
  word-wrap: break-word;
  word-break: break-word;
  /* 确保正常文字的空格宽度正常 */
  word-spacing: normal;
  letter-spacing: normal;
  white-space: normal;
  /* 防止溢出 */
  overflow-x: hidden;
  box-sizing: border-box;
}

/* 确保代码块不受父容器的 word-wrap 影响 */
.comment-content :deep(pre),
.comment-content :deep(pre code) {
  word-wrap: normal !important;
  word-break: normal !important;
}

.comment-content :deep(.vuepress-markdown-body) {
  padding: 0;
  background: transparent !important;
  font-size: 14px;
  line-height: 1.6;
  color: inherit;
  /* 防止溢出 */
  overflow-x: hidden;
  width: 100%;
  max-width: 100%;
  /* 确保正常文字的空格宽度正常 */
  word-spacing: normal;
  letter-spacing: normal;
  white-space: normal;
}

/* 确保代码块容器不受父容器宽度限制 */
.comment-content :deep(.vuepress-markdown-body pre) {
  max-width: none !important;
  width: 100% !important;
}

/* 评论内容样式优化 */
.comment-content :deep(.vuepress-markdown-body p) {
  margin: 8px 0;
  /* 确保段落中的空格宽度正常 */
  word-spacing: normal;
  letter-spacing: normal;
  white-space: normal;
}

.comment-content :deep(.vuepress-markdown-body p:first-child) {
  margin-top: 0;
}

.comment-content :deep(.vuepress-markdown-body p:last-child) {
  margin-bottom: 0;
}

/* 链接样式 */
.comment-content :deep(.vuepress-markdown-body a) {
  color: #18a058;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: all 0.2s;
  word-break: break-all;
}

.comment-content :deep(.vuepress-markdown-body a:hover) {
  border-bottom-color: #18a058;
}

/* 代码块样式 - 确保长代码行不换行 */
.comment-content :deep(pre) {
  position: relative;
  border-radius: 4px;
  margin: 8px 0;
  padding: 8px 12px;
  /* 使用更深的背景色，提高对比度和可读性 */
  background: #2d2d2d;
  border: 1px solid rgba(0, 0, 0, 0.1);
  overflow-x: auto !important;
  overflow-y: hidden;
  white-space: pre !important;
  word-spacing: normal !important;
  letter-spacing: normal !important;
  word-wrap: normal !important;
  word-break: normal !important;
  font-size: 13px;
  line-height: 1.4;
  width: 100%;
  max-width: 100%;
  box-sizing: border-box;
  /* 确保代码块可以水平滚动 */
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
}

.comment-content :deep(pre code) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.4;
  /* 使用浅色文字，与深色背景形成良好对比 */
  color: #f8f8f2;
  display: block;
  white-space: pre !important;
  word-spacing: normal !important;
  letter-spacing: normal !important;
  word-wrap: normal !important;
  word-break: normal !important;
  padding: 0;
  margin: 0;
  /* 关键：使用 max-content 确保内容不被压缩，但最小宽度为100% */
  width: max-content;
  min-width: 100%;
  box-sizing: border-box;
  /* 确保代码内容可以完整显示，不会被裁剪 */
  display: inline-block;
}

/* 强制重置所有 Prism.js token 样式，彻底消除间距，但保持空格正常显示 */
.comment-content :deep(pre code .token),
.comment-content :deep(pre code span),
.comment-content :deep(pre code .token.string),
.comment-content :deep(pre code .token.keyword),
.comment-content :deep(pre code .token.function),
.comment-content :deep(pre code .token.operator),
.comment-content :deep(pre code .token.punctuation),
.comment-content :deep(pre code .token.property),
.comment-content :deep(pre code .token.attr-name),
.comment-content :deep(pre code .token.attr-value),
.comment-content :deep(pre code .token.comment),
.comment-content :deep(pre code .token.number),
.comment-content :deep(pre code .token.boolean),
.comment-content :deep(pre code .token.variable),
.comment-content :deep(pre code *) {
  margin: 0 !important;
  padding: 0 !important;
  word-spacing: normal !important;
  letter-spacing: normal !important;
  word-wrap: normal !important;
  word-break: normal !important;
  white-space: inherit !important;
  display: inline !important;
  border: none !important;
  background: transparent !important;
  box-sizing: border-box !important;
  line-height: inherit !important;
  vertical-align: baseline !important;
  font-size: inherit !important;
}

/* 移除所有相邻元素之间的间距 */
.comment-content :deep(pre code .token + .token),
.comment-content :deep(pre code span + span),
.comment-content :deep(pre code .token + span),
.comment-content :deep(pre code span + .token),
.comment-content :deep(pre code * + *) {
  margin-left: 0 !important;
  margin-right: 0 !important;
  padding-left: 0 !important;
  padding-right: 0 !important;
  border-left: none !important;
  border-right: none !important;
}

/* 确保代码块内文本紧凑 */
.comment-content :deep(pre code) {
  font-variant-ligatures: none;
  font-feature-settings: normal;
  text-rendering: auto;
}

/* 行内代码样式 - 优化间距，让看起来不那么密集 */
.comment-content :deep(code:not(pre code)) {
  background: rgba(150, 150, 150, 0.1);
  padding: 3px 8px;
  margin: 0 2px;
  border-radius: 3px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9em;
  color: #e83e8c;
  white-space: nowrap;
  word-break: keep-all;
  display: inline-block;
  vertical-align: baseline;
}

/* 图片样式 */
.comment-content :deep(.vuepress-markdown-body img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 8px 0;
  cursor: pointer;
  transition: opacity 0.2s;
  display: block;
}

.comment-content :deep(.vuepress-markdown-body img:hover) {
  opacity: 0.9;
}

/* 复制按钮样式 */
.comment-content :deep(.copy-code-btn) {
  position: absolute;
  top: 8px;
  right: 8px;
  padding: 4px 12px;
  background: rgba(255, 255, 255, 0.9);
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 12px;
  color: #333;
  cursor: pointer;
  opacity: 0;
  transition: all 0.3s;
  z-index: 10;
}

.comment-content :deep(.copy-code-btn:hover) {
  background: #fff;
  border-color: #18a058;
  color: #18a058;
}

.comment-content :deep(pre:hover .copy-code-btn) {
  opacity: 1;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .comment-content :deep(.vuepress-markdown-body) {
    font-size: 15px;
    line-height: 1.7;
  }

  .comment-content :deep(.vuepress-markdown-body img) {
    margin: 12px 0;
    border-radius: 6px;
  }

  .comment-content :deep(pre) {
    padding: 6px 8px;
    font-size: 12px;
    margin: 6px 0;
  }

  .comment-content :deep(pre code) {
    font-size: 12px;
  }
}

/* 小屏幕移动端优化（小于420px） */
@media (max-width: 420px) {
  .comment-content {
    padding: 0;
    /* 允许代码块溢出以便滚动查看 */
    overflow-x: auto;
    width: 100%;
    max-width: 100%;
  }

  .comment-content :deep(.vuepress-markdown-body) {
    font-size: 14px;
    /* 允许代码块溢出以便滚动查看 */
    overflow-x: auto;
    padding: 0;
    width: 100%;
    max-width: 100%;
  }

  .comment-content :deep(pre) {
    padding: 6px 12px;
    font-size: 11px;
    /* 使用负margin向右扩展，利用父容器的padding空间 */
    margin: 6px -8px;
    border-radius: 4px;
    /* 代码块宽度扩展到父容器宽度加上左右负margin */
    width: calc(100% + 16px);
    min-width: calc(100% + 16px);
    max-width: calc(100% + 16px);
    box-sizing: border-box;
    /* 确保代码块可以水平滚动 */
    overflow-x: auto;
    overflow-y: hidden;
    /* 代码块向左对齐，确保左侧内容可见 */
    display: block;
    /* 确保滚动条显示正常 */
    -webkit-overflow-scrolling: touch;
    scrollbar-width: thin;
    /* 确保代码块从左侧开始显示 */
    direction: ltr;
    text-align: left;
  }

  .comment-content :deep(pre code) {
    font-size: 11px;
    padding: 0;
    /* 确保代码内容可以完整显示，不被压缩 */
    display: block;
    width: max-content;
    min-width: 100%;
    white-space: pre !important;
  }

  /* 确保代码块容器可以扩展 */
  .comment-content :deep(.vuepress-markdown-body pre) {
    max-width: calc(100% + 16px);
    width: calc(100% + 16px);
    min-width: calc(100% + 16px);
  }
}

/* GitHub 风格提示块样式 */
.comment-content :deep(blockquote.markdown-alert) {
  padding: 0.75rem 1rem !important;
  margin: 1rem 0 !important;
  border-left: 0.25rem solid !important;
  border-radius: 0.375rem !important;
  background-color: rgba(0, 0, 0, 0.02) !important;
  color: inherit !important;
}

.comment-content :deep(.markdown-alert-title) {
  font-weight: 600 !important;
  margin-bottom: 0.5rem !important;
  display: flex !important;
  align-items: center !important;
  gap: 0.5rem !important;
}

.comment-content :deep(.markdown-alert p) {
  margin: 0.5rem 0 0 0 !important;
}

.comment-content :deep(.markdown-alert p:first-child) {
  margin-top: 0 !important;
}

/* NOTE 提示块 */
.comment-content :deep(blockquote.markdown-alert-note) {
  border-left-color: #0969da !important;
  background-color: #ddf4ff !important;
}

.comment-content :deep(.markdown-alert-note .markdown-alert-title) {
  color: #0969da !important;
}

/* TIP 提示块 */
.comment-content :deep(blockquote.markdown-alert-tip) {
  border-left-color: #1a7f37 !important;
  background-color: #dafbe1 !important;
}

.comment-content :deep(.markdown-alert-tip .markdown-alert-title) {
  color: #1a7f37 !important;
}

/* IMPORTANT 提示块 */
.comment-content :deep(blockquote.markdown-alert-important) {
  border-left-color: #8250df !important;
  background-color: #fbefff !important;
}

.comment-content :deep(.markdown-alert-important .markdown-alert-title) {
  color: #8250df !important;
}

/* WARNING 提示块 */
.comment-content :deep(blockquote.markdown-alert-warning) {
  border-left-color: #9a6700 !important;
  background-color: #fff8c5 !important;
}

.comment-content :deep(.markdown-alert-warning .markdown-alert-title) {
  color: #9a6700 !important;
}

/* CAUTION 提示块 */
.comment-content :deep(blockquote.markdown-alert-caution) {
  border-left-color: #cf222e !important;
  background-color: #ffeef0 !important;
}

.comment-content :deep(.markdown-alert-caution .markdown-alert-title) {
  color: #cf222e !important;
}

/* DANGER 提示块 */
.comment-content :deep(blockquote.markdown-alert-danger) {
  border-left-color: #cf222e !important;
  background-color: #ffeef0 !important;
}

.comment-content :deep(.markdown-alert-danger .markdown-alert-title) {
  color: #cf222e !important;
}

/* INFO 提示块 */
.comment-content :deep(blockquote.markdown-alert-info) {
  border-left-color: #0550ae !important;
  background-color: #e6f0ff !important;
}

.comment-content :deep(.markdown-alert-info .markdown-alert-title) {
  color: #0969da !important;
}

/* EXAMPLE 提示块 */
.comment-content :deep(blockquote.markdown-alert-example) {
  border-left-color: #8250df !important;
  background-color: #fbefff !important;
}

.comment-content :deep(.markdown-alert-example .markdown-alert-title) {
  color: #8250df !important;
}

/* QUESTION 提示块 */
.comment-content :deep(blockquote.markdown-alert-question) {
  border-left-color: #9a6700 !important;
  background-color: #fff8c5 !important;
}

.comment-content :deep(.markdown-alert-question .markdown-alert-title) {
  color: #9a6700 !important;
}

/* QUOTE 提示块 */
.comment-content :deep(blockquote.markdown-alert-quote) {
  border-left-color: #656d76 !important;
  background-color: #f6f8fa !important;
}

.comment-content :deep(.markdown-alert-quote .markdown-alert-title) {
  color: #656d76 !important;
}

/* 暗色模式支持 */
html.dark .comment-content :deep(.vuepress-markdown-body) {
  color: #d1d5db !important;
}

html.dark .comment-content :deep(.vuepress-markdown-body a) {
  color: #38bdf8 !important;
}

html.dark .comment-content :deep(.vuepress-markdown-body a:hover) {
  border-bottom-color: #38bdf8;
}

html.dark .comment-content :deep(pre) {
  background: #1e1e1e !important;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

html.dark .comment-content :deep(pre code) {
  color: #d4d4d4 !important;
}

html.dark .comment-content :deep(pre code .token) {
  margin: 0;
  padding: 0;
  word-spacing: normal;
  letter-spacing: normal;
}

html.dark .comment-content :deep(code:not(pre code)) {
  background: rgba(56, 189, 248, 0.15) !important;
  color: #38bdf8 !important;
}

html.dark .comment-content :deep(.copy-code-btn) {
  background: rgba(30, 41, 59, 0.9) !important;
  border-color: rgba(255, 255, 255, 0.2) !important;
  color: #d1d5db !important;
}

html.dark .comment-content :deep(.copy-code-btn:hover) {
  background: rgba(51, 65, 85, 0.95) !important;
  border-color: #38bdf8 !important;
  color: #38bdf8 !important;
}

/* 暗黑模式下的提示块样式 */
html.dark .comment-content :deep(blockquote.markdown-alert) {
  background-color: rgba(255, 255, 255, 0.04) !important;
}

html.dark .comment-content :deep(blockquote.markdown-alert-note) {
  border-left-color: #4493f8 !important;
  background-color: rgba(68, 147, 248, 0.15) !important;
}

html.dark .comment-content :deep(.markdown-alert-note .markdown-alert-title) {
  color: #4493f8 !important;
}

html.dark .comment-content :deep(blockquote.markdown-alert-tip) {
  border-left-color: #3fb950 !important;
  background-color: rgba(63, 185, 80, 0.15) !important;
}

html.dark .comment-content :deep(.markdown-alert-tip .markdown-alert-title) {
  color: #3fb950 !important;
}

html.dark .comment-content :deep(blockquote.markdown-alert-important) {
  border-left-color: #a371f7 !important;
  background-color: rgba(163, 113, 247, 0.15) !important;
}

html.dark .comment-content :deep(.markdown-alert-important .markdown-alert-title) {
  color: #a371f7 !important;
}

html.dark .comment-content :deep(blockquote.markdown-alert-warning) {
  border-left-color: #d29922 !important;
  background-color: rgba(210, 153, 34, 0.15) !important;
}

html.dark .comment-content :deep(.markdown-alert-warning .markdown-alert-title) {
  color: #d29922 !important;
}

html.dark .comment-content :deep(blockquote.markdown-alert-caution) {
  border-left-color: #f85149 !important;
  background-color: rgba(248, 81, 73, 0.15) !important;
}

html.dark .comment-content :deep(.markdown-alert-caution .markdown-alert-title) {
  color: #f85149 !important;
}

html.dark .comment-content :deep(blockquote.markdown-alert-danger) {
  border-left-color: #f85149 !important;
  background-color: rgba(248, 81, 73, 0.15) !important;
}

html.dark .comment-content :deep(.markdown-alert-danger .markdown-alert-title) {
  color: #f85149 !important;
}

html.dark .comment-content :deep(blockquote.markdown-alert-info) {
  border-left-color: #58a6ff !important;
  background-color: rgba(88, 166, 255, 0.15) !important;
}

html.dark .comment-content :deep(.markdown-alert-info .markdown-alert-title) {
  color: #4493f8 !important;
}

html.dark .comment-content :deep(blockquote.markdown-alert-example) {
  border-left-color: #a371f7 !important;
  background-color: rgba(163, 113, 247, 0.15) !important;
}

html.dark .comment-content :deep(.markdown-alert-example .markdown-alert-title) {
  color: #a371f7 !important;
}

html.dark .comment-content :deep(blockquote.markdown-alert-question) {
  border-left-color: #d29922 !important;
  background-color: rgba(210, 153, 34, 0.15) !important;
}

html.dark .comment-content :deep(.markdown-alert-question .markdown-alert-title) {
  color: #d29922 !important;
}

html.dark .comment-content :deep(blockquote.markdown-alert-quote) {
  border-left-color: #8b949e !important;
  background-color: rgba(139, 148, 158, 0.15) !important;
}

html.dark .comment-content :deep(.markdown-alert-quote .markdown-alert-title) {
  color: #8b949e !important;
}
</style>
