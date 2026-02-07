<!--
  项目名称：blog-frontend
  文件名称：MarkdownEditor.vue
  创建时间：2026-02-01 20:03:19

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：Markdown编辑器组件，提供Markdown编辑和实时预览功能，支持代码高亮、图片上传插入、代码块复制等功能。
-->
<template>
  <div class="markdown-editor" ref="editorRef">
    <v-md-editor
      v-model="content"
      :height="height"
      :disabled-menus="[]"
      :subfield="subfield"
      :mode="mode"
      @upload-image="handleUploadImage"
      @change="handleChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, onMounted, nextTick } from 'vue'
import { uploadImage } from '@/api/upload'
import { useMessage } from 'naive-ui'

interface Props {
  modelValue?: string
  height?: string
  subfield?: boolean
  mode?: 'edit' | 'preview' | 'editable'
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  height: '500px',
  subfield: true,
  mode: 'editable'
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

