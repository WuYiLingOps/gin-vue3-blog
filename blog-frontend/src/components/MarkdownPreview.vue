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
let observer: MutationObserver | null = null

// 添加复制按钮到代码块并确保滚动位置正确
function addCopyButtons() {
  if (!previewRef.value) return

  const codeBlocks = previewRef.value.querySelectorAll('pre code')
  
  codeBlocks.forEach((codeBlock) => {
    const pre = codeBlock.parentElement as HTMLElement
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
  
  // 确保所有代码块滚动到最左侧，显示完整内容
  // 使用 setTimeout 确保 DOM 完全渲染后再设置滚动位置
  setTimeout(() => {
    const allPreElements = previewRef.value?.querySelectorAll('pre') || []
    allPreElements.forEach((pre) => {
      const preElement = pre as HTMLElement
      // 强制设置滚动位置为 0，确保从左侧开始显示
      preElement.scrollLeft = 0
      // 如果内容宽度超过容器，确保可以滚动
      if (preElement.scrollWidth > preElement.clientWidth) {
        preElement.style.overflowX = 'auto'
      }
    })
  }, 100)
}

// 确保代码块滚动位置正确的函数
function ensureCodeBlockScrollPosition() {
  if (!previewRef.value) return
  
  const setScrollPosition = () => {
    const allPreElements = previewRef.value?.querySelectorAll('pre') || []
    allPreElements.forEach((pre) => {
      const preElement = pre as HTMLElement
      
      // 确保代码块宽度限制在容器内
      preElement.style.maxWidth = '100%'
      preElement.style.width = '100%'
      preElement.style.boxSizing = 'border-box'
      preElement.style.overflowX = 'auto'
      preElement.style.overflowY = 'hidden'
      preElement.style.left = '0'
      preElement.style.marginLeft = '0'
      preElement.style.transform = 'translateX(0)'
      
      // 强制设置滚动位置为 0，确保从左侧开始显示
      // 使用多种方法确保滚动位置正确
      preElement.scrollLeft = 0
      preElement.scrollTo(0, 0)
      
      // 确保代码内容从左侧开始，不被裁剪
      const codeElement = preElement.querySelector('code') as HTMLElement
      if (codeElement) {
        // 重置所有可能影响位置的样式
        codeElement.style.left = '0'
        codeElement.style.marginLeft = '0'
        codeElement.style.paddingLeft = '0'
        codeElement.style.position = 'relative'
        codeElement.style.boxSizing = 'content-box'
        codeElement.style.transform = 'translateX(0)'
        codeElement.style.textAlign = 'left'
        codeElement.style.direction = 'ltr'
        codeElement.style.float = 'none'
        codeElement.style.clear = 'both'
        codeElement.style.textIndent = '0'
        
        // 确保代码内容宽度正确
        // 获取 pre 的实际可用宽度（减去 padding）
        const prePaddingLeft = parseInt(window.getComputedStyle(preElement).paddingLeft) || 0
        const prePaddingRight = parseInt(window.getComputedStyle(preElement).paddingRight) || 0
        const preAvailableWidth = preElement.clientWidth - prePaddingLeft - prePaddingRight
        
        codeElement.style.width = 'max-content'
        // 使用 calc 确保 minWidth 考虑 padding
        codeElement.style.minWidth = `${preAvailableWidth}px`
        // 如果计算失败，使用 100% 作为后备
        if (preAvailableWidth <= 0) {
          codeElement.style.minWidth = '100%'
        }
        
        // 重置代码块内所有子元素的位置
        const allChildren = codeElement.querySelectorAll('*')
        allChildren.forEach((child) => {
          const childElement = child as HTMLElement
          childElement.style.marginLeft = '0'
          childElement.style.paddingLeft = '0'
          childElement.style.left = '0'
          childElement.style.transform = 'translateX(0)'
        })
        
        // 再次强制设置滚动位置，确保代码内容可见
        requestAnimationFrame(() => {
          preElement.scrollLeft = 0
          preElement.scrollTo(0, 0)
          // 双重检查，确保滚动位置正确
          if (preElement.scrollLeft !== 0) {
            preElement.scrollLeft = 0
          }
        })
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
    ensureCodeBlockScrollPosition()
    setupCodeBlockObserver()
  })
})

onBeforeUnmount(() => {
  if (observer) {
    observer.disconnect()
    observer = null
  }
})

watch(() => props.content, () => {
  nextTick(() => {
    addCopyButtons()
    ensureCodeBlockScrollPosition()
    // 重新设置 observer
    if (observer) {
      observer.disconnect()
      observer = null
    }
    setupCodeBlockObserver()
  })
})
</script>

<style scoped>
.markdown-preview {
  width: 100%;
  max-width: 100%;
  overflow-x: hidden !important;
  box-sizing: border-box;
  /* 确保不会超出父容器 */
  position: relative;
}

.markdown-preview :deep(.vuepress-markdown-body) {
  padding: 16px 0;
  background: transparent !important;
  /* 确保 markdown 容器不会溢出 */
  overflow-x: hidden !important;
  max-width: 100% !important;
  width: 100% !important;
  box-sizing: border-box !important;
  /* 确保所有子元素都不会超出 */
  position: relative;
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
  padding: 12px 16px;
  background: #f5f5f5;
  border: 1px solid rgba(0, 0, 0, 0.1);
  /* 关键：代码块本身必须限制在容器内 */
  max-width: 100% !important;
  width: 100% !important;
  box-sizing: border-box !important;
  /* 代码块内部可以滚动 */
  overflow-x: auto !important;
  overflow-y: hidden;
  white-space: pre !important;
  word-wrap: normal !important;
  word-break: normal !important;
  /* 移动端优化：确保可以水平滚动 */
  -webkit-overflow-scrolling: touch;
  scrollbar-width: thin;
  /* 确保内容从左侧开始显示，不被裁剪 */
  direction: ltr;
  text-align: left;
  /* 强制初始滚动位置为 0 */
  scroll-behavior: auto;
  /* 确保 padding 计算在宽度内 */
  padding-left: 16px !important;
  padding-right: 16px !important;
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
  display: block;
  white-space: pre !important;
  word-wrap: normal !important;
  word-break: normal !important;
  padding: 0 !important;
  margin: 0 !important;
  /* 关键：代码内容可以超出 pre 的宽度，但 pre 本身限制在容器内 */
  width: max-content;
  min-width: 100%;
  box-sizing: content-box;
  overflow-x: visible;
  /* 确保内容从左侧开始显示 */
  direction: ltr;
  text-align: left;
  /* 确保代码内容不被裁剪，从最左侧开始 */
  position: relative;
  left: 0 !important;
  margin-left: 0 !important;
  padding-left: 0 !important;
  transform: translateX(0) !important;
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
    overflow-x: hidden !important;
    width: 100% !important;
    max-width: 100% !important;
    box-sizing: border-box !important;
  }

  .markdown-preview :deep(.vuepress-markdown-body) {
    overflow-x: hidden !important;
    max-width: 100% !important;
    width: 100% !important;
    word-wrap: break-word;
    box-sizing: border-box !important;
    padding-left: 0 !important;
    padding-right: 0 !important;
  }
  
  .markdown-preview :deep(pre) {
    margin: 12px 0 !important;
    padding: 10px 12px !important;
    border-radius: 4px;
    font-size: 13px;
    /* 关键：代码块本身必须限制在容器内 */
    max-width: 100% !important;
    width: 100% !important;
    box-sizing: border-box !important;
    /* 代码块内部可以滚动 */
    overflow-x: auto !important;
    overflow-y: hidden;
    -webkit-overflow-scrolling: touch;
    /* 确保左侧内容不被裁剪 */
    position: relative;
    /* 强制初始滚动位置为 0 */
    scroll-behavior: auto;
    /* 确保 padding 计算在宽度内 */
    padding-left: 12px !important;
    padding-right: 12px !important;
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
    overflow-x: hidden !important;
    width: 100% !important;
    max-width: 100% !important;
    box-sizing: border-box !important;
    /* 移除所有可能的 padding 和 margin */
    padding: 0 !important;
    margin: 0 !important;
  }

  .markdown-preview :deep(.vuepress-markdown-body) {
    overflow-x: hidden !important;
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
    margin: 8px 0 !important;
    /* 进一步减小 padding，为代码内容留出更多空间 */
    padding: 6px 4px !important;
    font-size: 11px;
    /* 关键：代码块本身必须限制在容器内，使用 calc 确保宽度计算正确 */
    max-width: 100% !important;
    width: 100% !important;
    box-sizing: border-box !important;
    /* 代码块内部可以滚动 */
    overflow-x: auto !important;
    overflow-y: hidden !important;
    -webkit-overflow-scrolling: touch;
    /* 确保左侧内容不被裁剪 */
    position: relative !important;
    /* 强制初始滚动位置为 0 */
    scroll-behavior: auto !important;
    /* 最小化 padding，确保最大可用宽度 */
    padding-left: 4px !important;
    padding-right: 4px !important;
    /* 确保代码块从左侧开始 */
    left: 0 !important;
    margin-left: 0 !important;
    margin-right: 0 !important;
    transform: translateX(0) !important;
    /* 确保滚动条可见且可用 */
    scrollbar-width: thin;
    scrollbar-color: rgba(255, 255, 255, 0.3) transparent;
  }

  .markdown-preview :deep(pre code) {
    font-size: 11px;
    line-height: 1.4;
    /* 代码内容可以超出 pre 的宽度，使用 max-content 确保完整显示 */
    display: block !important;
    width: max-content !important;
    min-width: calc(100% - 8px) !important; /* 减去左右 padding */
    padding: 0 !important;
    margin: 0 !important;
    box-sizing: content-box !important;
    /* 确保代码内容从左侧开始，不被裁剪 */
    position: relative !important;
    left: 0 !important;
    margin-left: 0 !important;
    margin-right: 0 !important;
    padding-left: 0 !important;
    padding-right: 0 !important;
    transform: translateX(0) !important;
    /* 强制重置所有可能影响位置的属性 */
    float: none !important;
    clear: both !important;
    text-indent: 0 !important;
    /* 确保文本方向正确 */
    direction: ltr !important;
    text-align: left !important;
    /* 确保不会被裁剪 */
    overflow: visible !important;
    white-space: pre !important;
    word-wrap: normal !important;
    word-break: normal !important;
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

/* 极小屏幕优化（420px-480px） */
@media (min-width: 420px) and (max-width: 480px) {
  .markdown-preview :deep(pre) {
    margin: 10px 0;
    padding: 8px 10px !important;
    font-size: 12px;
    /* 确保左侧内容完整显示 */
    overflow-x: auto !important;
    overflow-y: hidden;
    max-width: 100%;
    box-sizing: border-box;
    /* 强制初始滚动位置为 0 */
    scroll-behavior: auto;
    /* 确保左侧 padding 不会导致内容被裁剪 */
    padding-left: 10px !important;
    padding-right: 10px !important;
  }

  .markdown-preview {
    overflow-x: hidden !important;
    width: 100% !important;
    max-width: 100% !important;
  }

  .markdown-preview :deep(.vuepress-markdown-body) {
    overflow-x: hidden !important;
    max-width: 100% !important;
    width: 100% !important;
  }

  .markdown-preview :deep(pre) {
    max-width: 100% !important;
    width: 100% !important;
    box-sizing: border-box !important;
  }

  .markdown-preview :deep(pre code) {
    font-size: 12px;
    line-height: 1.4;
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

  /* 复制按钮在移动端始终显示 */
  .markdown-preview :deep(.copy-code-btn) {
    opacity: 1;
    padding: 6px 10px;
    font-size: 11px;
  }
}
</style>

