<template>
  <div class="comment-content" ref="previewRef">
    <!-- 使用Markdown预览组件，支持代码高亮和图片 -->
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
import 'prismjs/components/prism-yaml'

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

// 修复代码块间距问题
function fixCodeBlockSpacing() {
  if (!previewRef.value) return

  const codeBlocks = previewRef.value.querySelectorAll('pre code')
  
  codeBlocks.forEach((codeBlock) => {
    // 移除所有 token 和 span 的 margin 和 padding
    const allElements = codeBlock.querySelectorAll('*')
    allElements.forEach((element) => {
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
  images.forEach((imgElement) => {
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

watch(() => props.content, () => {
  nextTick(() => {
    addCopyButtons()
    handleImageClick()
    // 延迟执行，确保 Prism.js 完成渲染
    setTimeout(() => {
      fixCodeBlockSpacing()
    }, 100)
  })
})
</script>

<style scoped>
.comment-content {
  width: 100%;
  word-wrap: break-word;
  word-break: break-word;
  /* 确保正常文字的空格宽度正常 */
  word-spacing: normal;
  letter-spacing: normal;
  white-space: normal;
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
  /* 确保代码块容器不受限制 */
  overflow-x: visible;
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
  /* 关键：使用 max-content 确保内容不被压缩 */
  width: max-content;
  min-width: 100%;
  box-sizing: border-box;
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
</style>

