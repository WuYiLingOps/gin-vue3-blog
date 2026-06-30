<!--
  项目名称：blog-frontend
  文件名称：MarkdownPreview.vue
  创建时间：2026-02-01 20:03:19

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：Markdown预览组件，用于渲染Markdown内容为HTML，支持代码高亮、代码块复制功能，自动处理代码块滚动位置。
-->
<template>
  <div class="markdown-preview" ref="previewRef">
    <v-md-preview :text="content" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick, onBeforeUnmount } from 'vue'
import { useMessage } from 'naive-ui'
import { VMdPreview } from '@/plugins/v-md-editor'

interface Props {
  content: string
}

const props = defineProps<Props>()
const message = useMessage()
const previewRef = ref<HTMLElement>()
let observer: MutationObserver | null = null

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

// 添加复制按钮到代码块并确保滚动位置正确
function addCopyButtons() {
  if (!previewRef.value) return

  const codeBlocks = previewRef.value.querySelectorAll('pre code')

  codeBlocks.forEach(codeBlock => {
    const pre = codeBlock.parentElement as HTMLElement
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

// 确保代码块滚动位置正确的函数
function ensureCodeBlockScrollPosition() {
  if (!previewRef.value) return

  // 确保代码块在容器内正确显示且不丢失内容
  const setScrollPosition = () => {
    if (!previewRef.value) return
    const preElements = previewRef.value.querySelectorAll('pre')
    preElements.forEach(pre => {
      const preElement = pre as HTMLElement
      const codeElement = preElement.querySelector('code') as HTMLElement

      if (codeElement) {
        // 确保代码块宽度正确
        codeElement.style.width = 'max-content'
        codeElement.style.minWidth = '100%'

        // 重置滚动位置到最左侧
        preElement.scrollLeft = 0
      }
    })
  }

  // 立即执行一次
  setScrollPosition()

  // 延迟执行，确保 DOM 完全渲染和代码高亮完成
  setTimeout(setScrollPosition, 50)
  setTimeout(setScrollPosition, 100)
  setTimeout(setScrollPosition, 200)
  setTimeout(setScrollPosition, 500)
  setTimeout(setScrollPosition, 1000)
  // 针对小屏幕，增加额外的检查
  setTimeout(setScrollPosition, 1500)
  setTimeout(setScrollPosition, 2000)
  // 针对超小屏幕，增加更多检查
  setTimeout(setScrollPosition, 2500)
  setTimeout(setScrollPosition, 3000)
}

// 设置 MutationObserver 监听代码块变化
function setupCodeBlockObserver() {
  if (!previewRef.value || observer) return

  observer = new MutationObserver(() => {
    ensureCodeBlockScrollPosition()
  })

  observer.observe(previewRef.value, {
    childList: true,
    subtree: true,
    attributes: false
  })
}

onMounted(() => {
  nextTick(() => {
    addCopyButtons()
    processAlerts()
    ensureCodeBlockScrollPosition()
    setupCodeBlockObserver()
    window.addEventListener('resize', handleResize)
  })
})

onBeforeUnmount(() => {
  if (observer) {
    observer.disconnect()
    observer = null
  }
  window.removeEventListener('resize', handleResize)
})

// 处理窗口大小变化
function handleResize() {
  ensureCodeBlockScrollPosition()
}

watch(
  () => props.content,
  () => {
    nextTick(() => {
      addCopyButtons()
      processAlerts()
      ensureCodeBlockScrollPosition()
      // 重新设置 observer
      if (observer) {
        observer.disconnect()
        observer = null
      }
      setupCodeBlockObserver()
    })
  }
)
</script>

<style scoped>
.markdown-preview {
  width: 100%;
  max-width: 100%;
  overflow-x: visible;
  box-sizing: border-box;
  /* 确保不会超出父容器 */
  position: relative;
}

.markdown-preview :deep(.vuepress-markdown-body) {
  padding: 16px 0;
  background: transparent !important;
  /* 确保 markdown 容器可以容纳滚动内容 */
  overflow-x: visible;
  max-width: 100% !important;
  width: 100% !important;
  box-sizing: border-box !important;
  /* 确保所有子元素都不会超出 */
  position: relative;
}

/* 段落排版：解决全局 reset 导致的段落“挤在一起”，并增加首行缩进 */
.markdown-preview :deep(.vuepress-markdown-body > p) {
  margin: 0 0 1em 0;
  line-height: 1.9;
  text-align: justify;
  word-break: break-word;
}

.markdown-preview :deep(.vuepress-markdown-body > p:last-child) {
  margin-bottom: 0;
}

/* 暗黑模式下的 markdown 内容样式 */
html.dark .markdown-preview :deep(.vuepress-markdown-body) {
  background: transparent !important;
  color: #d1d5db !important;
}

/* 暗黑模式下的标题颜色 */
html.dark .markdown-preview :deep(.vuepress-markdown-body h1),
html.dark .markdown-preview :deep(.vuepress-markdown-body h2),
html.dark .markdown-preview :deep(.vuepress-markdown-body h3),
html.dark .markdown-preview :deep(.vuepress-markdown-body h4),
html.dark .markdown-preview :deep(.vuepress-markdown-body h5),
html.dark .markdown-preview :deep(.vuepress-markdown-body h6) {
  color: #e5e5e5 !important;
  border-bottom-color: rgba(255, 255, 255, 0.1) !important;
}

/* 暗黑模式下的链接颜色 */
html.dark .markdown-preview :deep(.vuepress-markdown-body a) {
  color: #38bdf8 !important;
}

/* 日间模式下的引用块样式 */
.markdown-preview :deep(.vuepress-markdown-body blockquote) {
  color: #4b5563 !important;
  border-left: 4px solid #0891b2 !important;
  background: rgba(8, 145, 178, 0.08) !important;
  padding: 12px 16px !important;
  margin: 16px 0 !important;
  border-radius: 4px !important;
  font-style: normal !important;
}

.markdown-preview :deep(.vuepress-markdown-body blockquote p) {
  margin: 0 !important;
  color: #4b5563 !important;
}

/* 暗黑模式下的引用块 */
html.dark .markdown-preview :deep(.vuepress-markdown-body blockquote) {
  color: #9ca3af !important;
  border-left: 4px solid rgba(56, 189, 248, 0.5) !important;
  background: rgba(56, 189, 248, 0.05) !important;
  padding: 12px 16px !important;
  margin: 16px 0 !important;
  border-radius: 4px !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body blockquote p) {
  color: #9ca3af !important;
}

/* 日间模式下的表格样式 */
.markdown-preview :deep(.vuepress-markdown-body table) {
  border-collapse: collapse !important;
  border: 1px solid #e5e7eb !important;
  width: 100% !important;
  margin: 16px auto !important;
  background: #ffffff !important;
  border-radius: 8px !important;
  overflow: hidden !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.05) !important;
  display: table !important;
}

.markdown-preview :deep(.vuepress-markdown-body th) {
  background: #dbeafe !important;
  color: #1e40af !important;
  font-weight: 600 !important;
  padding: 12px 16px !important;
  text-align: left !important;
  border: 1px solid #bfdbfe !important;
  font-size: 14px !important;
}

.markdown-preview :deep(.vuepress-markdown-body td) {
  padding: 12px 16px !important;
  border: 1px solid #e5e7eb !important;
  color: #374151 !important;
  font-size: 14px !important;
  line-height: 1.6 !important;
}

.markdown-preview :deep(.vuepress-markdown-body tr) {
  background: #ffffff !important;
  border-top: 1px solid #e5e7eb !important;
  transition: background-color 0.2s ease !important;
}

.markdown-preview :deep(.vuepress-markdown-body tbody tr:hover) {
  background: rgba(8, 145, 178, 0.05) !important;
}

.markdown-preview :deep(.vuepress-markdown-body tr:nth-child(2n)) {
  background: #f9fafb !important;
}

.markdown-preview :deep(.vuepress-markdown-body tbody tr:nth-child(2n):hover) {
  background: rgba(8, 145, 178, 0.08) !important;
}

/* 暗黑模式下的表格 */
html.dark .markdown-preview :deep(.vuepress-markdown-body table) {
  border-color: rgba(255, 255, 255, 0.15) !important;
  background: rgba(30, 41, 59, 0.6) !important;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.3) !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body th) {
  background: rgba(56, 189, 248, 0.2) !important;
  color: #7dd3fc !important;
  font-weight: 600 !important;
  border: 1px solid rgba(56, 189, 248, 0.3) !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body td) {
  border: 1px solid rgba(255, 255, 255, 0.15) !important;
  color: #d1d5db !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body tr) {
  background: rgba(30, 41, 59, 0.4) !important;
  border-top-color: rgba(255, 255, 255, 0.1) !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body tbody tr:hover) {
  background: rgba(56, 189, 248, 0.1) !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body tr:nth-child(2n)) {
  background: rgba(255, 255, 255, 0.05) !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body tbody tr:nth-child(2n):hover) {
  background: rgba(56, 189, 248, 0.15) !important;
}

/* 暗黑模式下的分隔线 */
html.dark .markdown-preview :deep(.vuepress-markdown-body hr) {
  border-color: rgba(255, 255, 255, 0.1) !important;
  background: rgba(255, 255, 255, 0.1) !important;
}

/* 暗黑模式下的列表 */
html.dark .markdown-preview :deep(.vuepress-markdown-body li) {
  color: #d1d5db !important;
}

/* 暗黑模式下的强调文本 */
html.dark .markdown-preview :deep(.vuepress-markdown-body strong) {
  color: #e5e5e5 !important;
}

/* 暗黑模式下的图片 */
html.dark .markdown-preview :deep(.vuepress-markdown-body img) {
  opacity: 0.9;
  border-radius: 8px;
}

/* 代码块样式优化 */
.markdown-preview :deep(pre) {
  position: relative;
  border-radius: 12px;
  margin: 20px auto !important;
  padding: 12px 16px;
  background: #f8f9fa;
  border: 1px solid rgba(0, 0, 0, 0.05);
  /* 居中显示优化 */
  width: 98% !important;
  max-width: calc(100% - 16px) !important; /* 预留左右缓冲空间 */
  box-sizing: border-box !important;
  /* 内部滚动 */
  overflow-x: auto !important;
  overflow-y: hidden;
  white-space: pre !important;
  word-wrap: normal !important;
  word-break: normal !important;
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
  direction: ltr;
  text-align: left;
  /* 阴影效果增加质感 */
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}

/* 暗黑模式下的代码块背景 */
html.dark .markdown-preview :deep(pre) {
  background: rgba(15, 23, 42, 0.9) !important;
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
}

.markdown-preview :deep(pre code) {
  font-family: 'Fira Code', 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
  display: block;
  width: max-content;
  min-width: 100%;
  padding: 0 !important;
  margin: 0 !important;
  background: transparent !important;
}

/* 行内代码样式 */
.markdown-preview :deep(code:not(pre code)) {
  background: rgba(150, 150, 150, 0.1);
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9em;
  color: #e83e8c;
}

/* 暗黑模式下的行内代码 */
html.dark .markdown-preview :deep(code:not(pre code)) {
  background: rgba(56, 189, 248, 0.15) !important;
  color: #38bdf8 !important;
}

/* 复制按钮样式 */
.markdown-preview :deep(.copy-code-btn) {
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

.markdown-preview :deep(.copy-code-btn:hover) {
  background: #fff;
  border-color: #18a058;
  color: #18a058;
}

.markdown-preview :deep(pre:hover .copy-code-btn) {
  opacity: 1;
}

/* 暗黑模式下的复制按钮样式 */
html.dark .markdown-preview :deep(.copy-code-btn) {
  background: rgba(30, 41, 59, 0.9) !important;
  border-color: rgba(255, 255, 255, 0.2) !important;
  color: #d1d5db !important;
}

html.dark .markdown-preview :deep(.copy-code-btn:hover) {
  background: rgba(51, 65, 85, 0.95) !important;
  border-color: #38bdf8 !important;
  color: #38bdf8 !important;
}

/* 移动端代码块优化 */
@media (max-width: 768px) {
  /* 确保所有容器都不溢出 */
  .markdown-preview {
    overflow-x: visible;
    width: 100% !important;
    max-width: 100% !important;
    box-sizing: border-box !important;
  }

  .markdown-preview :deep(.vuepress-markdown-body) {
    overflow-x: visible;
    max-width: 100% !important;
    width: 100% !important;
    word-wrap: break-word;
    box-sizing: border-box !important;
    padding-left: 0 !important;
    padding-right: 0 !important;
  }

  .markdown-preview :deep(pre) {
    margin: 12px auto !important;
    padding: 10px 12px !important;
    border-radius: 12px;
    font-size: 13px;
    /* 关键：代码块居中并预留缓冲空间 */
    width: 96% !important;
    max-width: calc(100% - 12px) !important;
    box-sizing: border-box !important;
    /* 代码块内部可以滚动 */
    overflow-x: auto !important;
    overflow-y: hidden;
    -webkit-overflow-scrolling: touch;
    /* 确保左侧内容不被裁剪 */
    position: relative;
    /* 强制初始滚动位置为 0 */
    scroll-behavior: auto;
  }

  .markdown-preview :deep(pre code) {
    font-size: 13px;
    line-height: 1.5;
    /* 代码内容可以超出 pre 的宽度 */
    display: block;
    width: max-content;
    min-width: 100%;
    padding: 0 !important;
    margin: 0 !important;
    box-sizing: content-box;
    /* 确保代码内容从左侧开始，不被裁剪 */
    position: relative;
    left: 0 !important;
    margin-left: 0 !important;
    padding-left: 0 !important;
    transform: translateX(0) !important;
  }

  /* 表格在移动端也需要处理溢出 */
  .markdown-preview :deep(.vuepress-markdown-body table) {
    display: block;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    max-width: 100%;
  }
}

/* 超小屏幕优化（小于420px） */
@media (max-width: 419px) {
  /* 确保所有容器都不溢出 */
  .markdown-preview {
    overflow-x: hidden;
    width: 100% !important;
    max-width: 100% !important;
    box-sizing: border-box !important;
    /* 移除所有可能的 padding 和 margin */
    padding: 0 !important;
    margin: 0 !important;
  }

  .markdown-preview :deep(.vuepress-markdown-body) {
    overflow-x: hidden;
    max-width: 100% !important;
    width: 100% !important;
    word-wrap: break-word;
    box-sizing: border-box !important;
    /* 移除左右 padding，最大化可用空间 */
    padding-left: 0 !important;
    padding-right: 0 !important;
    padding-top: 12px !important;
    padding-bottom: 12px !important;
  }

  .markdown-preview :deep(pre) {
    margin: 12px auto !important;
    /* 进一步减小 padding，为代码内容留出更多空间 */
    padding: 8px 10px !important;
    font-size: 11px;
    /* 关键：居中对齐 + 缓冲空间，彻底解决左侧内容丢失 */
    width: 95% !important;
    max-width: calc(100% - 10px) !important;
    box-sizing: border-box !important;
    /* 代码块内部可以滚动 */
    overflow-x: auto !important;
    overflow-y: hidden !important;
    -webkit-overflow-scrolling: touch;
    /* 确保滚动条可见且可用 */
    scrollbar-width: thin;
    scrollbar-color: rgba(255, 255, 255, 0.3) transparent;
    border-radius: 12px;
    position: relative;
    /* 确保左侧对齐 */
    text-align: left;
    direction: ltr;
    /* 强制居中 */
    left: 50%;
    transform: translateX(-50%);
  }

  .markdown-preview :deep(pre code) {
    font-size: 11px;
    line-height: 1.4;
    /* 代码内容可以超出 pre 的宽度，使用 max-content 确保完整显示 */
    display: block !important;
    width: max-content !important;
    min-width: 100% !important;
    padding: 0 !important;
    margin: 0 !important;
    background: transparent !important;
    /* 强制左对齐 */
    text-align: left !important;
    position: relative;
    left: 0 !important;
  }

  /* 确保代码块内的所有子元素都不影响位置 */
  .markdown-preview :deep(pre code *) {
    margin-left: 0 !important;
    margin-right: 0 !important;
    padding-left: 0 !important;
    padding-right: 0 !important;
    left: 0 !important;
    transform: translateX(0) !important;
    position: relative !important;
    float: none !important;
    clear: both !important;
  }

  /* 确保代码块内的 token 元素不影响位置 */
  .markdown-preview :deep(pre code .token) {
    margin-left: 0 !important;
    margin-right: 0 !important;
    padding-left: 0 !important;
    padding-right: 0 !important;
    left: 0 !important;
    transform: translateX(0) !important;
    position: relative !important;
  }

  /* 表格在超小屏幕也需要处理溢出 */
  .markdown-preview :deep(.vuepress-markdown-body table) {
    display: block;
    overflow-x: auto;
    -webkit-overflow-scrolling: touch;
    max-width: 100%;
    font-size: 11px;
  }
}

/* GitHub 风格提示块样式 - 使用 !important 确保覆盖 vuepress 主题 */
.markdown-preview :deep(blockquote.markdown-alert) {
  padding: 0.75rem 1rem !important;
  margin: 1rem 0 !important;
  border-left: 0.25rem solid !important;
  border-radius: 0.375rem !important;
  background-color: rgba(0, 0, 0, 0.02) !important;
  color: inherit !important;
}

.markdown-preview :deep(.markdown-alert-title) {
  font-weight: 600 !important;
  margin-bottom: 0.5rem !important;
  display: flex !important;
  align-items: center !important;
  gap: 0.5rem !important;
}

.markdown-preview :deep(.markdown-alert p) {
  margin: 0.5rem 0 0 0 !important;
}

.markdown-preview :deep(.markdown-alert p:first-child) {
  margin-top: 0 !important;
}

/* NOTE 提示块 */
.markdown-preview :deep(blockquote.markdown-alert-note) {
  border-left-color: #0969da !important;
  background-color: #ddf4ff !important;
}

.markdown-preview :deep(.markdown-alert-note .markdown-alert-title) {
  color: #0969da !important;
}

/* TIP 提示块 */
.markdown-preview :deep(blockquote.markdown-alert-tip) {
  border-left-color: #1a7f37 !important;
  background-color: #dafbe1 !important;
}

.markdown-preview :deep(.markdown-alert-tip .markdown-alert-title) {
  color: #1a7f37 !important;
}

/* IMPORTANT 提示块 */
.markdown-preview :deep(blockquote.markdown-alert-important) {
  border-left-color: #8250df !important;
  background-color: #fbefff !important;
}

.markdown-preview :deep(.markdown-alert-important .markdown-alert-title) {
  color: #8250df !important;
}

/* WARNING 提示块 */
.markdown-preview :deep(blockquote.markdown-alert-warning) {
  border-left-color: #9a6700 !important;
  background-color: #fff8c5 !important;
}

.markdown-preview :deep(.markdown-alert-warning .markdown-alert-title) {
  color: #9a6700 !important;
}

/* CAUTION 提示块 */
.markdown-preview :deep(blockquote.markdown-alert-caution) {
  border-left-color: #cf222e !important;
  background-color: #ffeef0 !important;
}

.markdown-preview :deep(.markdown-alert-caution .markdown-alert-title) {
  color: #cf222e !important;
}

/* DANGER 提示块 */
.markdown-preview :deep(blockquote.markdown-alert-danger) {
  border-left-color: #cf222e !important;
  background-color: #ffeef0 !important;
}

.markdown-preview :deep(.markdown-alert-danger .markdown-alert-title) {
  color: #cf222e !important;
}

/* INFO 提示块 */
.markdown-preview :deep(blockquote.markdown-alert-info) {
  border-left-color: #0550ae !important;
  background-color: #e6f0ff !important;
}

.markdown-preview :deep(.markdown-alert-info .markdown-alert-title) {
  color: #0969da !important;
}

/* EXAMPLE 提示块 */
.markdown-preview :deep(blockquote.markdown-alert-example) {
  border-left-color: #8250df !important;
  background-color: #fbefff !important;
}

.markdown-preview :deep(.markdown-alert-example .markdown-alert-title) {
  color: #8250df !important;
}

/* QUESTION 提示块 */
.markdown-preview :deep(blockquote.markdown-alert-question) {
  border-left-color: #9a6700 !important;
  background-color: #fff8c5 !important;
}

.markdown-preview :deep(.markdown-alert-question .markdown-alert-title) {
  color: #9a6700 !important;
}

/* QUOTE 提示块 */
.markdown-preview :deep(blockquote.markdown-alert-quote) {
  border-left-color: #656d76 !important;
  background-color: #f6f8fa !important;
}

.markdown-preview :deep(.markdown-alert-quote .markdown-alert-title) {
  color: #656d76 !important;
}

/* 暗黑模式下的提示块样式 */
html.dark .markdown-preview :deep(.markdown-alert) {
  background-color: rgba(255, 255, 255, 0.04);
}

html.dark .markdown-preview :deep(blockquote.markdown-alert-note) {
  border-left-color: #4493f8 !important;
  background-color: rgba(68, 147, 248, 0.15) !important;
}

html.dark .markdown-preview :deep(.markdown-alert-note .markdown-alert-title) {
  color: #4493f8 !important;
}

html.dark .markdown-preview :deep(blockquote.markdown-alert-tip) {
  border-left-color: #3fb950 !important;
  background-color: rgba(63, 185, 80, 0.15) !important;
}

html.dark .markdown-preview :deep(.markdown-alert-tip .markdown-alert-title) {
  color: #3fb950 !important;
}

html.dark .markdown-preview :deep(blockquote.markdown-alert-important) {
  border-left-color: #a371f7 !important;
  background-color: rgba(163, 113, 247, 0.15) !important;
}

html.dark .markdown-preview :deep(.markdown-alert-important .markdown-alert-title) {
  color: #a371f7 !important;
}

html.dark .markdown-preview :deep(blockquote.markdown-alert-warning) {
  border-left-color: #d29922 !important;
  background-color: rgba(210, 153, 34, 0.15) !important;
}

html.dark .markdown-preview :deep(.markdown-alert-warning .markdown-alert-title) {
  color: #d29922 !important;
}

html.dark .markdown-preview :deep(blockquote.markdown-alert-caution) {
  border-left-color: #f85149 !important;
  background-color: rgba(248, 81, 73, 0.15) !important;
}

html.dark .markdown-preview :deep(.markdown-alert-caution .markdown-alert-title) {
  color: #f85149 !important;
}

html.dark .markdown-preview :deep(blockquote.markdown-alert-danger) {
  border-left-color: #f85149 !important;
  background-color: rgba(248, 81, 73, 0.15) !important;
}

html.dark .markdown-preview :deep(.markdown-alert-danger .markdown-alert-title) {
  color: #f85149 !important;
}

html.dark .markdown-preview :deep(blockquote.markdown-alert-info) {
  border-left-color: #58a6ff !important;
  background-color: rgba(88, 166, 255, 0.15) !important;
}

html.dark .markdown-preview :deep(.markdown-alert-info .markdown-alert-title) {
  color: #4493f8 !important;
}

html.dark .markdown-preview :deep(blockquote.markdown-alert-example) {
  border-left-color: #a371f7 !important;
  background-color: rgba(163, 113, 247, 0.15) !important;
}

html.dark .markdown-preview :deep(.markdown-alert-example .markdown-alert-title) {
  color: #a371f7 !important;
}

html.dark .markdown-preview :deep(blockquote.markdown-alert-question) {
  border-left-color: #d29922 !important;
  background-color: rgba(210, 153, 34, 0.15) !important;
}

html.dark .markdown-preview :deep(.markdown-alert-question .markdown-alert-title) {
  color: #d29922 !important;
}

html.dark .markdown-preview :deep(blockquote.markdown-alert-quote) {
  border-left-color: #8b949e !important;
  background-color: rgba(139, 148, 158, 0.15) !important;
}

html.dark .markdown-preview :deep(.markdown-alert-quote .markdown-alert-title) {
  color: #8b949e !important;
}

/* 极小屏幕优化（420px-480px） */
@media (min-width: 420px) and (max-width: 480px) {
  .markdown-preview :deep(pre) {
    margin: 16px auto !important;
    padding: 10px 12px !important;
    font-size: 12px;
    overflow-x: auto !important;
    /* 居中对齐 */
    width: 96% !important;
    max-width: calc(100% - 12px) !important;
    box-sizing: border-box;
    border-radius: 12px;
    left: 50%;
    transform: translateX(-50%);
  }

  .markdown-preview :deep(pre code) {
    font-size: 12px;
    line-height: 1.4;
    display: block;
    width: max-content;
    min-width: 100%;
    padding: 0 !important;
    margin: 0 !important;
    background: transparent !important;
    position: relative;
    left: 0 !important;
  }
}

/* 复制按钮默认隐藏，仅在悬停时显示（包括移动端，移动端通过点击触发悬停态） */
.markdown-preview :deep(.copy-code-btn) {
  opacity: 0;
  padding: 6px 10px;
  font-size: 11px;
}

.markdown-preview :deep(pre:hover .copy-code-btn) {
  opacity: 1;
}
</style>
