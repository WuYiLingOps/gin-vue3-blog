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
import VMdEditor from '@kangc/v-md-editor'
import '@kangc/v-md-editor/lib/style/base-editor.css'
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js'
import '@kangc/v-md-editor/lib/theme/style/vuepress.css'
import Prism from 'prismjs'
import { uploadImage } from '@/api/upload'
import { useMessage } from 'naive-ui'

// 配置编辑器主题（评论编辑器不需要代码高亮，但主题需要 Prism）
VMdEditor.use(vuepressTheme, {
  Prism,
  codeHighlightExtensionMap: {
    vue: 'html',
  }
})

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
}

.comment-markdown-editor :deep(.v-md-editor__left-area textarea) {
  font-size: 14px;
  line-height: 1.6;
  padding: 12px;
}

/* 预览区域样式 */
.comment-markdown-editor :deep(.v-md-editor__right-area) {
  font-size: 14px;
  line-height: 1.6;
}

.comment-markdown-editor :deep(.v-md-editor__right-area .vuepress-markdown-body) {
  padding: 12px;
  font-size: 14px;
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

