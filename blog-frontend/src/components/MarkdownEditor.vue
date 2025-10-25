<template>
  <div class="markdown-editor" ref="editorRef">
    <v-md-editor
      v-model="content"
      :height="height"
      :disabled-menus="[]"
      @upload-image="handleUploadImage"
      @change="handleChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, nextTick } from 'vue'
import VMdEditor from '@kangc/v-md-editor'
import '@kangc/v-md-editor/lib/style/base-editor.css'
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js'
import '@kangc/v-md-editor/lib/theme/style/vuepress.css'
import Prism from 'prismjs'
import { uploadImage } from '@/api/upload'
import { useMessage } from 'naive-ui'

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

// 配置编辑器主题
VMdEditor.use(vuepressTheme, {
  Prism,
  codeHighlightExtensionMap: {
    vue: 'html',
  }
})

interface Props {
  modelValue?: string
  height?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  height: '500px'
})

const emit = defineEmits<Emits>()
const message = useMessage()
const editorRef = ref<HTMLElement>()

const content = ref(props.modelValue)

// 添加复制按钮到代码块
function addCopyButtons() {
  if (!editorRef.value) return

  const codeBlocks = editorRef.value.querySelectorAll('.v-md-editor__right-area pre code')
  
  codeBlocks.forEach((codeBlock) => {
    const pre = codeBlock.parentElement
    if (!pre || pre.querySelector('.copy-code-btn')) return

    const button = document.createElement('button')
    button.className = 'copy-code-btn'
    button.textContent = '复制'
    button.onclick = (e) => {
      e.stopPropagation()
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

watch(
  () => props.modelValue,
  (newValue) => {
    content.value = newValue
  }
)

watch(content, () => {
  nextTick(() => {
    addCopyButtons()
  })
})

onMounted(() => {
  nextTick(() => {
    addCopyButtons()
  })
})

function handleChange(text: string) {
  emit('update:modelValue', text)
  emit('change', text)
  nextTick(() => {
    addCopyButtons()
  })
}

async function handleUploadImage(
  _event: Event,
  insertImage: (arg: { url: string; desc?: string; width?: string; height?: string }) => void,
  files: File[]
) {
  try {
    const file = files[0]
    if (!file) return

    // 上传图片（已在 API 中自动拼接完整 URL）
    const res = await uploadImage(file)
    const imageUrl = res.data?.url || ''

    // 插入图片到编辑器
    insertImage({
      url: imageUrl,
      desc: file.name
    })

    message.success('图片上传成功')
  } catch (error: any) {
    message.error(error.message || '图片上传失败')
  }
}
</script>

<style scoped>
.markdown-editor {
  width: 100%;
}

/* 代码块样式优化 */
.markdown-editor :deep(pre) {
  position: relative;
  border-radius: 6px;
}

.markdown-editor :deep(pre code) {
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 14px;
  line-height: 1.6;
}

/* 行内代码样式 */
.markdown-editor :deep(code:not(pre code)) {
  background: rgba(150, 150, 150, 0.1);
  padding: 2px 6px;
  border-radius: 3px;
  font-family: 'Consolas', 'Monaco', 'Courier New', monospace;
  font-size: 0.9em;
  color: #e83e8c;
}

/* 编辑器和预览区代码高亮 */
.markdown-editor :deep(.v-md-editor__left-area pre),
.markdown-editor :deep(.v-md-editor__right-area pre) {
  margin: 16px 0;
}

/* 复制按钮样式 */
.markdown-editor :deep(.copy-code-btn) {
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

.markdown-editor :deep(.copy-code-btn:hover) {
  background: #fff;
  border-color: #18a058;
  color: #18a058;
}

.markdown-editor :deep(.v-md-editor__right-area pre:hover .copy-code-btn) {
  opacity: 1;
}

/* 暗色主题下的按钮样式 */
@media (prefers-color-scheme: dark) {
  .markdown-editor :deep(.copy-code-btn) {
    background: rgba(40, 40, 40, 0.9);
    border-color: #555;
    color: #ccc;
  }

  .markdown-editor :deep(.copy-code-btn:hover) {
    background: #333;
    border-color: #18a058;
    color: #18a058;
  }
}
</style>

