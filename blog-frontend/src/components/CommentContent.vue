<template>
  <div class="comment-content" ref="previewRef">
    <v-md-preview :text="content" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, watch, nextTick } from 'vue'
import { useMessage } from 'naive-ui'
import VMdPreview from '@kangc/v-md-editor/lib/preview'
import '@kangc/v-md-editor/lib/style/preview.css'
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js'
import '@kangc/v-md-editor/lib/theme/style/vuepress.css'
import Prism from 'prismjs'

// 引入代码高亮语言包（精简版）
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-sql'
import 'prismjs/components/prism-markdown'
import 'prismjs/components/prism-css'

// 配置主题
VMdPreview.use(vuepressTheme, {
  Prism,
  codeHighlightExtensionMap: {
    vue: 'html',
  }
})

interface Props {
  content: string
}

const props = defineProps<Props>()
const message = useMessage()
const previewRef = ref<HTMLElement>()

// 添加复制按钮到代码块
function addCopyButtons() {
  if (!previewRef.value) return

  const codeBlocks = previewRef.value.querySelectorAll('pre code')
  
  codeBlocks.forEach((codeBlock) => {
    const pre = codeBlock.parentElement
    if (!pre || pre.querySelector('.copy-code-btn')) return

    const button = document.createElement('button')
    button.className = 'copy-code-btn'
    button.textContent = '复制'
    button.onclick = () => {
      const code = codeBlock.textContent || ''
      navigator.clipboard.writeText(code).then(() => {
        button.textContent = '已复制!'
        message.success('代码已复制到剪贴板')
        setTimeout(() => {
          button.textContent = '复制'
        }, 2000)
      }).catch(() => {
        message.error('复制失败，请手动复制')
      })
    }

    pre.style.position = 'relative'
    pre.appendChild(button)
  })
}

// 处理图片点击放大
function handleImageClick() {
  if (!previewRef.value) return

  const images = previewRef.value.querySelectorAll('img')
  images.forEach((img) => {
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
    handleImageClick()
  })
})

watch(() => props.content, () => {
  nextTick(() => {
    addCopyButtons()
    handleImageClick()
  })
})
</script>

<style scoped>
.comment-content {
  width: 100%;
  word-wrap: break-word;
  word-break: break-word;
}

.comment-content :deep(.vuepress-markdown-body) {
  padding: 0;
  background: transparent !important;
  font-size: 14px;
  line-height: 1.6;
  color: inherit;
}

/* 评论内容样式优化 */
.comment-content :deep(.vuepress-markdown-body p) {
  margin: 8px 0;
}

.comment-content :deep(.vuepress-markdown-body p:first-child) {
  margin-top: 0;
}

.comment-content :deep(.vuepress-markdown-body p:last-child) {
  margin-bottom: 0;
}

/* 标题样式（评论中较少使用，但保留） */
.comment-content :deep(.vuepress-markdown-body h1),
.comment-content :deep(.vuepress-markdown-body h2),
.comment-content :deep(.vuepress-markdown-body h3),
.comment-content :deep(.vuepress-markdown-body h4),
.comment-content :deep(.vuepress-markdown-body h5),
.comment-content :deep(.vuepress-markdown-body h6) {
  margin: 12px 0 8px 0;
  font-weight: 600;
}

.comment-content :deep(.vuepress-markdown-body h1) {
  font-size: 1.5em;
}

.comment-content :deep(.vuepress-markdown-body h2) {
  font-size: 1.3em;
}

.comment-content :deep(.vuepress-markdown-body h3) {
  font-size: 1.1em;
}

/* 链接样式 */
.comment-content :deep(.vuepress-markdown-body a) {
  color: #18a058;
  text-decoration: none;
  border-bottom: 1px solid transparent;
  transition: all 0.2s;
}

.comment-content :deep(.vuepress-markdown-body a:hover) {
  border-bottom-color: #18a058;
}

/* 引用块样式 */
.comment-content :deep(.vuepress-markdown-body blockquote) {
  margin: 8px 0;
  padding: 8px 12px;
  border-left: 3px solid #18a058;
  background: rgba(24, 160, 88, 0.05);
  color: #666;
  border-radius: 4px;
}

/* 代码块样式 */
.comment-content :deep(pre) {
  position: relative;
  border-radius: 4px;
  margin: 8px 0;
  padding: 12px;
  background: #f5f5f5;
  overflow-x: auto;
}

.comment-content :deep(pre code) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 13px;
  line-height: 1.5;
  color: #333;
}

/* 行内代码样式 */
.comment-content :deep(code:not(pre code)) {
  background: rgba(150, 150, 150, 0.1);
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9em;
  color: #e83e8c;
}

/* 列表样式 */
.comment-content :deep(.vuepress-markdown-body ul),
.comment-content :deep(.vuepress-markdown-body ol) {
  margin: 8px 0;
  padding-left: 24px;
}

.comment-content :deep(.vuepress-markdown-body li) {
  margin: 4px 0;
}

/* 表格样式 */
.comment-content :deep(.vuepress-markdown-body table) {
  margin: 8px 0;
  border-collapse: collapse;
  width: 100%;
  overflow-x: auto;
  display: block;
}

.comment-content :deep(.vuepress-markdown-body table thead),
.comment-content :deep(.vuepress-markdown-body table tbody) {
  display: table;
  width: 100%;
}

.comment-content :deep(.vuepress-markdown-body th),
.comment-content :deep(.vuepress-markdown-body td) {
  padding: 6px 12px;
  border: 1px solid #ddd;
}

.comment-content :deep(.vuepress-markdown-body th) {
  background: #f5f5f5;
  font-weight: 600;
}

/* 图片样式 */
.comment-content :deep(.vuepress-markdown-body img) {
  max-width: 100%;
  height: auto;
  border-radius: 4px;
  margin: 8px 0;
  cursor: pointer;
  transition: opacity 0.2s;
}

.comment-content :deep(.vuepress-markdown-body img:hover) {
  opacity: 0.9;
}

/* 分隔线 */
.comment-content :deep(.vuepress-markdown-body hr) {
  margin: 12px 0;
  border: none;
  border-top: 1px solid #ddd;
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

/* 暗色模式支持 */
html.dark .comment-content :deep(.vuepress-markdown-body) {
  color: #d1d5db !important;
}

html.dark .comment-content :deep(.vuepress-markdown-body a) {
  color: #38bdf8 !important;
}

html.dark .comment-content :deep(.vuepress-markdown-body blockquote) {
  background: rgba(56, 189, 248, 0.05) !important;
  border-left-color: #38bdf8 !important;
  color: #9ca3af !important;
}

html.dark .comment-content :deep(pre) {
  background: rgba(15, 23, 42, 0.8) !important;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

html.dark .comment-content :deep(pre code) {
  color: #d1d5db !important;
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
</style>

