<template>
  <div class="comment-markdown-editor" ref="editorRef">
    <v-md-editor
      v-model="content"
      :height="height"
      :mode="mode"
      :disabled-menus="disabledMenus"
      :toolbar="toolbar"
      @upload-image="handleUploadImage"
      @change="handleChange"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, watch, computed } from 'vue'
import VMdEditor from '@kangc/v-md-editor'
import '@kangc/v-md-editor/lib/style/base-editor.css'
import vuepressTheme from '@kangc/v-md-editor/lib/theme/vuepress.js'
import '@kangc/v-md-editor/lib/theme/style/vuepress.css'
import Prism from 'prismjs'
import { uploadImage } from '@/api/upload'
import { useMessage } from 'naive-ui'

// 引入代码高亮语言包（精简版，只保留常用语言）
import 'prismjs/components/prism-json'
import 'prismjs/components/prism-bash'
import 'prismjs/components/prism-typescript'
import 'prismjs/components/prism-javascript'
import 'prismjs/components/prism-go'
import 'prismjs/components/prism-python'
import 'prismjs/components/prism-sql'
import 'prismjs/components/prism-markdown'
import 'prismjs/components/prism-css'

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
  mode: 'editable', // 默认双栏模式（编辑+预览）
  placeholder: '写下你的评论...支持 Markdown 语法',
  maxLength: 5000
})

const emit = defineEmits<Emits>()
const message = useMessage()
const editorRef = ref<HTMLElement>()

const content = ref(props.modelValue)

// 精简的工具栏配置（只保留常用功能）
const toolbar = computed(() => [
  'bold',
  'italic',
  'strike',
  'quote',
  '|',
  'link',
  'code',
  'table',
  '|',
  'image',
  '|',
  'unordered-list',
  'ordered-list',
  '|',
  'preview',
  'fullscreen'
])

// 禁用的菜单项（移除不常用的功能）
const disabledMenus = computed(() => [
  'h1',
  'h2',
  'h3',
  'h4',
  'h5',
  'h6',
  'hr',
  'save'
])

watch(
  () => props.modelValue,
  (newValue) => {
    content.value = newValue
  }
)

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

.comment-markdown-editor :deep(.v-md-editor__toolbar) {
  border-bottom: 1px solid var(--n-border-color);
  padding: 8px;
}

.comment-markdown-editor :deep(.v-md-editor__toolbar-item) {
  margin: 0 4px;
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
  .comment-markdown-editor :deep(.v-md-editor) {
    font-size: 14px;
  }
  
  .comment-markdown-editor :deep(.v-md-editor__toolbar) {
    padding: 6px;
  }
  
  .comment-markdown-editor :deep(.v-md-editor__toolbar-item) {
    margin: 0 2px;
    padding: 4px;
  }
}
</style>

