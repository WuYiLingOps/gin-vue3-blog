<template>
  <div class="markdown-preview" ref="previewRef">
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

// 引入代码高亮语言包
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-java'
import 'prismjs/components/prism-c'
import 'prismjs/components/prism-cpp'
import 'prismjs/components/prism-sql'
import 'prismjs/components/prism-yaml'
import 'prismjs/components/prism-markdown'
import 'prismjs/components/prism-css'
import 'prismjs/components/prism-scss'

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

onMounted(() => {
  nextTick(() => {
    addCopyButtons()
  })
})

watch(() => props.content, () => {
  nextTick(() => {
    addCopyButtons()
  })
})
</script>

<style scoped>
.markdown-preview {
  width: 100%;
}

.markdown-preview :deep(.vuepress-markdown-body) {
  padding: 16px 0;
  background: transparent !important;
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

/* 暗黑模式下的引用块 */
html.dark .markdown-preview :deep(.vuepress-markdown-body blockquote) {
  color: #9ca3af !important;
  border-left-color: rgba(56, 189, 248, 0.5) !important;
  background: rgba(56, 189, 248, 0.05) !important;
}

/* 暗黑模式下的表格 */
html.dark .markdown-preview :deep(.vuepress-markdown-body table) {
  border-color: rgba(255, 255, 255, 0.1) !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body th),
html.dark .markdown-preview :deep(.vuepress-markdown-body td) {
  border-color: rgba(255, 255, 255, 0.1) !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body tr) {
  background: transparent !important;
  border-top-color: rgba(255, 255, 255, 0.1) !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body tr:nth-child(2n)) {
  background: rgba(255, 255, 255, 0.02) !important;
}

html.dark .markdown-preview :deep(.vuepress-markdown-body th) {
  background: rgba(255, 255, 255, 0.05) !important;
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
  border-radius: 6px;
  margin: 16px 0;
}

/* 暗黑模式下的代码块背景 */
html.dark .markdown-preview :deep(pre) {
  background: rgba(15, 23, 42, 0.8) !important;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.markdown-preview :deep(pre code) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
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
</style>

