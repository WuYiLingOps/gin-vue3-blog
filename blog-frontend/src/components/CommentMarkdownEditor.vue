<!--
  项目名称：blog-frontend
  文件名称：CommentMarkdownEditor.vue
  创建时间：2026-02-01 20:03:19

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：评论Markdown编辑器组件，提供Markdown编辑和实时预览功能，包含自定义工具栏（粗体、斜体、链接、图片），支持图片上传插入、代码高亮等功能。
-->
<template>
  <div class="comment-markdown-editor" ref="editorRef">
    <!-- 自定义工具栏 -->
    <div class="custom-toolbar" v-if="showCustomToolbar">
      <n-space size="small" align="center">
        <n-button 
          size="small" 
          quaternary 
          @click="insertMarkdown('bold')"
          title="粗体"
        >
          <strong>B</strong>
        </n-button>
        <n-button 
          size="small" 
          quaternary 
          @click="insertMarkdown('italic')"
          title="斜体"
        >
          <em>I</em>
        </n-button>
        <n-divider vertical />
        <n-button 
          size="small" 
          quaternary 
          @click="insertMarkdown('link')"
          title="链接"
        >
          🔗
        </n-button>
        <n-button 
          size="small" 
          quaternary 
          @click="triggerImageUpload"
          title="图片"
        >
          🖼️
        </n-button>
      </n-space>
    </div>
    
    <v-md-editor
      v-model="content"
      :height="height"
      :mode="mode"
      :disabled-menus="disabledMenus"
      :toolbar="[]"
      @upload-image="handleUploadImage"
      @change="handleChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed, onMounted, nextTick } from 'vue'
import { NButton, NSpace, NDivider } from 'naive-ui'
import { uploadImage } from '@/api/upload'
import { useMessage } from 'naive-ui'

interface Props {
  modelValue?: string
  height?: string
  mode?: 'edit' | 'editable' | 'preview'
  placeholder?: string
  maxLength?: number
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'change', value: string): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  height: '250px',
  mode: 'edit', // 单栏编辑模式，取消预览
  placeholder: '写下你的评论...支持 Markdown 语法',
  maxLength: 5000
})

const emit = defineEmits<Emits>()
const message = useMessage()
const editorRef = ref<HTMLElement>()

const content = ref(props.modelValue)
const showCustomToolbar = ref(true)

// 禁用的菜单项（移除不常用的功能，但保留code功能，只是不在工具栏显示）
const disabledMenus = computed(() => [
  'h1',
  'h2',
  'h3',
  'h4',
  'h5',
  'h6',
  'hr',
  'save',
  'strike',
  'quote',
  'code', // 禁用代码按钮（但支持直接输入代码块语法）
  'table',
  'unordered-list',
  'ordered-list',
  'preview',
  'fullscreen'
])

watch(
  () => props.modelValue,
  (newValue) => {
    content.value = newValue
  }
)

// 插入Markdown语法
function insertMarkdown(type: 'bold' | 'italic' | 'link') {
  if (!editorRef.value) return
  
  const textarea = editorRef.value.querySelector('textarea') as HTMLTextAreaElement
  if (!textarea) return
  
  const start = textarea.selectionStart
  const end = textarea.selectionEnd
  const selectedText = content.value.substring(start, end)
  
  let insertText = ''
  switch (type) {
    case 'bold':
      insertText = selectedText ? `**${selectedText}**` : '****'
      break
    case 'italic':
      insertText = selectedText ? `*${selectedText}*` : '**'
      break
    case 'link':
      insertText = selectedText ? `[${selectedText}](url)` : '[链接文本](url)'
      break
  }
  
  const newContent = 
    content.value.substring(0, start) + 
    insertText + 
    content.value.substring(end)
  
  content.value = newContent
  emit('update:modelValue', newContent)
  emit('change', newContent)
  
  // 恢复焦点和光标位置
  nextTick(() => {
    textarea.focus()
    const newPosition = type === 'link' && !selectedText 
      ? start + insertText.indexOf('url')
      : start + insertText.length - (type === 'bold' && !selectedText ? 2 : 0)
    textarea.setSelectionRange(newPosition, newPosition)
  })
}

// 触发图片上传
function triggerImageUpload() {
  if (!editorRef.value) return
  
  const input = document.createElement('input')
  input.type = 'file'
  input.accept = 'image/*'
  input.onchange = async (e) => {
    const file = (e.target as HTMLInputElement).files?.[0]
    if (!file) return
    
    try {
      const res = await uploadImage(file)
      const imageUrl = res.data?.url || ''
      
      if (!imageUrl) {
        message.error('图片上传失败')
        return
      }
      
      // 插入图片Markdown语法
      const textarea = editorRef.value?.querySelector('textarea') as HTMLTextAreaElement
      if (textarea) {
        const start = textarea.selectionStart
        const insertText = `![${file.name}](${imageUrl})`
        const newContent = 
          content.value.substring(0, start) + 
          insertText + 
          content.value.substring(start)
        
        content.value = newContent
        emit('update:modelValue', newContent)
        emit('change', newContent)
        
        nextTick(() => {
          textarea.focus()
          textarea.setSelectionRange(start + insertText.length, start + insertText.length)
        })
      }
      
      message.success('图片上传成功')
    } catch (error: any) {
      message.error(error.message || '图片上传失败')
    }
  }
  input.click()
}

onMounted(() => {
  // 隐藏编辑器自带的工具栏
  nextTick(() => {
    if (editorRef.value) {
      const toolbar = editorRef.value.querySelector('.v-md-editor__toolbar')
      if (toolbar) {
        ;(toolbar as HTMLElement).style.display = 'none'
      }
    }
  })
})

function handleChange(text: string) {
  // 检查长度限制
  if (props.maxLength && text.length > props.maxLength) {
    message.warning(`评论内容不能超过 ${props.maxLength} 个字符`)
    content.value = text.substring(0, props.maxLength)
    return
  }
  
  emit('update:modelValue', text)
  emit('change', text)
}

async function handleUploadImage(
  _event: Event,
  insertImage: (arg: { url: string; desc?: string; width?: string; height?: string }) => void,
  files: File[]
) {
  try {
    const file = files[0]
    if (!file) return

    // 检查文件类型
    if (!file.type.startsWith('image/')) {
      message.error('只能上传图片文件')
      return
    }

    // 检查文件大小（限制为 5MB）
    const maxSize = 5 * 1024 * 1024
    if (file.size > maxSize) {
      message.error('图片大小不能超过 5MB')
      return
    }

    // 上传图片
    const res = await uploadImage(file)
    const imageUrl = res.data?.url || ''

    if (!imageUrl) {
      message.error('图片上传失败')
      return
    }

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
.comment-markdown-editor {
  width: 100%;
}

/* 评论编辑器样式优化 */
.comment-markdown-editor :deep(.v-md-editor) {
  border-radius: 6px;
  border: 1px solid var(--n-border-color);
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(10px);
  box-shadow: inset 0 1px 0 rgba(255, 255, 255, 0.15);
  transition: background 0.2s ease, border-color 0.2s ease, box-shadow 0.2s ease;
}

html.dark .comment-markdown-editor :deep(.v-md-editor) {
  /* 使用 !important 覆盖编辑器内置的白色背景 */
  background: rgba(15, 23, 42, 0.96) !important;
  border-color: rgba(56, 189, 248, 0.18);
  box-shadow: 0 12px 30px rgba(0, 0, 0, 0.35);
}

.comment-markdown-editor :deep(.v-md-editor:hover) {
  border-color: #22c55e;
}

.comment-markdown-editor :deep(.v-md-editor:focus-within) {
  border-color: #22c55e;
  box-shadow: 0 0 0 2px rgba(34, 197, 94, 0.25), inset 0 1px 0 rgba(255, 255, 255, 0.2);
}

/* 隐藏编辑器自带的工具栏 */
.comment-markdown-editor :deep(.v-md-editor__toolbar) {
  display: none !important;
}

/* 自定义工具栏样式 */
.custom-toolbar {
  padding: 8px 12px;
  border-bottom: 1px solid var(--n-border-color);
  background: var(--n-color);
  border-radius: 6px 6px 0 0;
}

/* 编辑区域样式 */
.comment-markdown-editor :deep(.v-md-editor__left-area) {
  font-size: 14px;
  line-height: 1.6;
  background: transparent;
}

.comment-markdown-editor :deep(.v-md-editor__left-area textarea) {
  font-size: 14px;
  line-height: 1.7;
  padding: 12px;
  background: transparent;
  color: #0f172a;
  border: none;
  box-shadow: none;
}

.comment-markdown-editor :deep(.v-md-editor__left-area textarea::placeholder) {
  color: rgba(15, 23, 42, 0.45);
}

html.dark .comment-markdown-editor :deep(.v-md-editor__left-area textarea) {
  color: #f9fafb;
  font-weight: 500;
  caret-color: #38bdf8;
  border: none;
  background: transparent;
}

html.dark .comment-markdown-editor :deep(.v-md-editor__left-area textarea::placeholder) {
  color: rgba(226, 232, 240, 0.55);
}

/* 预览区域样式 */
.comment-markdown-editor :deep(.v-md-editor__right-area) {
  font-size: 14px;
  line-height: 1.6;
  background: transparent;
}

.comment-markdown-editor :deep(.v-md-editor__right-area .vuepress-markdown-body) {
  padding: 12px;
  font-size: 14px;
  color: #0f172a;
  background: transparent;
}

html.dark .comment-markdown-editor :deep(.v-md-editor__right-area .vuepress-markdown-body) {
  color: #f9fafb;
  background: transparent;
}

/* 夜间模式下预览文字与标题对比度增强 */
html.dark .comment-markdown-editor :deep(.vuepress-markdown-body p),
html.dark .comment-markdown-editor :deep(.vuepress-markdown-body li) {
  color: #f9fafb;
}

html.dark .comment-markdown-editor :deep(.vuepress-markdown-body h1),
html.dark .comment-markdown-editor :deep(.vuepress-markdown-body h2),
html.dark .comment-markdown-editor :deep(.vuepress-markdown-body h3),
html.dark .comment-markdown-editor :deep(.vuepress-markdown-body h4),
html.dark .comment-markdown-editor :deep(.vuepress-markdown-body h5),
html.dark .comment-markdown-editor :deep(.vuepress-markdown-body h6) {
  color: #f9fafb;
}

/* 确保编辑区域整体为暗色背景 */
html.dark .comment-markdown-editor :deep(.v-md-editor__main),
html.dark .comment-markdown-editor :deep(.v-md-editor__editor-wrapper),
html.dark .comment-markdown-editor :deep(.v-md-editor__content),
html.dark .comment-markdown-editor :deep(.v-md-textarea-editor),
html.dark .comment-markdown-editor :deep(.v-md-textarea-editor pre),
html.dark .comment-markdown-editor :deep(.v-md-textarea-editor textarea) {
  background: rgba(15, 23, 42, 0.96) !important;
}

/* 夜间模式下进一步强制编辑区文字颜色和粗细，避免被库默认样式覆盖 */
html.dark .comment-markdown-editor :deep(.v-md-textarea-editor textarea),
html.dark .comment-markdown-editor :deep(.v-md-textarea-editor pre) {
  color: #f9fafb !important;
  font-weight: 500;
  text-shadow: 0 0 1px rgba(0, 0, 0, 0.6);
}

/* 代码块样式 */
.comment-markdown-editor :deep(pre) {
  border-radius: 4px;
  margin: 8px 0;
}

.comment-markdown-editor :deep(pre code) {
  font-size: 13px;
  line-height: 1.5;
}

/* 行内代码样式 */
.comment-markdown-editor :deep(code:not(pre code)) {
  background: rgba(150, 150, 150, 0.1);
  padding: 2px 4px;
  border-radius: 3px;
  font-size: 0.9em;
}

/* 移动端优化 */
@media (max-width: 768px) {
  .comment-markdown-editor {
    font-size: 16px; /* 移动端防止自动缩放 */
  }
  
  .comment-markdown-editor :deep(.v-md-editor) {
    font-size: 14px;
    min-height: 150px;
  }
  
  .custom-toolbar {
    padding: 6px 8px;
  }
  
  .comment-markdown-editor :deep(.v-md-editor__left-area textarea) {
    font-size: 16px; /* 移动端防止自动缩放 */
    padding: 10px;
    line-height: 1.6;
  }
}
</style>

