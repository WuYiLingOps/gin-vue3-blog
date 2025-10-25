<template>
  <div class="image-upload">
    <div v-if="imageUrl" class="image-preview">
      <n-image
        :src="imageUrl"
        :width="width"
        :height="height"
        object-fit="cover"
        :alt="alt"
      />
      <div class="image-actions">
        <n-space>
          <n-button size="small" @click="handlePreview">
            <template #icon>
              <n-icon :component="EyeOutline" />
            </template>
            预览
          </n-button>
          <n-button size="small" @click="handleRemove">
            <template #icon>
              <n-icon :component="TrashOutline" />
            </template>
            删除
          </n-button>
        </n-space>
      </div>
    </div>

    <n-upload
      v-else
      :action="uploadUrl"
      :headers="uploadHeaders"
      :show-file-list="false"
      accept="image/*"
      @before-upload="handleBeforeUpload"
      @finish="handleFinish"
      @error="handleError"
    >
      <n-upload-dragger>
        <div class="upload-area" :style="{ width: width + 'px', height: height + 'px' }">
          <n-icon size="48" :component="CloudUploadOutline" />
          <n-text style="margin-top: 12px; display: block">
            点击或拖拽图片上传
          </n-text>
          <n-text depth="3" style="font-size: 12px; margin-top: 8px; display: block">
            支持 jpg、png、gif 格式，文件大小不超过 {{ maxSizeMB }}MB
          </n-text>
        </div>
      </n-upload-dragger>
    </n-upload>

    <!-- 图片预览弹窗 -->
    <n-modal v-model:show="showPreview" preset="card" style="width: 90%; max-width: 1000px">
      <template #header>
        图片预览
      </template>
      <div style="text-align: center">
        <img :src="imageUrl" style="max-width: 100%; height: auto" :alt="alt" />
      </div>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useMessage } from 'naive-ui'
import type { UploadFileInfo } from 'naive-ui'
import { CloudUploadOutline, EyeOutline, TrashOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores/auth'

interface Props {
  modelValue?: string
  width?: number
  height?: number
  maxSizeMB?: number
  alt?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'success', url: string): void
  (e: 'remove'): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  width: 400,
  height: 250,
  maxSizeMB: 5,
  alt: '图片'
})

const emit = defineEmits<Emits>()
const message = useMessage()
const authStore = useAuthStore()

const imageUrl = ref(props.modelValue)
const showPreview = ref(false)

// 监听外部传入的值变化
watch(() => props.modelValue, (newVal) => {
  imageUrl.value = newVal
})

const uploadUrl = computed(() => {
  return `${import.meta.env.VITE_API_BASE_URL}/api/upload/image`
})

const uploadHeaders = computed(() => {
  const token = authStore.token
  if (!token) {
    console.warn('未找到认证 token')
    return {}
  }
  return { Authorization: `Bearer ${token}` }
})

function handleBeforeUpload(data: { file: UploadFileInfo }) {
  const file = data.file.file
  
  // 检查是否有 token
  if (!authStore.token) {
    message.error('请先登录')
    return false
  }
  
  if (!file) {
    message.error('文件读取失败')
    return false
  }
  
  console.log('准备上传文件:', file.name, file.type, file.size)
  console.log('使用 token:', authStore.token ? '已设置' : '未设置')
  
  // 检查文件类型
  if (!file.type.startsWith('image/')) {
    message.error('只能上传图片文件')
    return false
  }

  // 检查文件大小
  const maxSize = props.maxSizeMB * 1024 * 1024
  if (file.size > maxSize) {
    message.error(`图片大小不能超过 ${props.maxSizeMB}MB`)
    return false
  }

  return true
}

function handleFinish({ event }: any) {
  try {
    const responseText = (event?.target as XMLHttpRequest).response
    console.log('上传响应:', responseText)
    
    const response = JSON.parse(responseText)
    console.log('解析后的响应:', response)
    
    // 后端返回 code 200 表示成功
    if (response.code === 200 && response.data?.url) {
      let url = response.data.url
      
      // 如果是相对路径，补充完整 URL
      if (url.startsWith('/')) {
        url = `${import.meta.env.VITE_API_BASE_URL}${url}`
      }
      
      console.log('最终图片 URL:', url)
      
      imageUrl.value = url
      emit('update:modelValue', url)
      emit('success', url)
      message.success('图片上传成功')
    } else {
      console.error('上传失败，响应:', response)
      message.error(response.message || '上传失败')
    }
  } catch (error) {
    console.error('解析上传响应失败:', error)
    message.error('上传失败')
  }
}

function handleError(error: any) {
  console.error('上传失败:', error)
  message.error('上传失败，请重试')
}

function handlePreview() {
  showPreview.value = true
}

function handleRemove() {
  imageUrl.value = ''
  emit('update:modelValue', '')
  emit('remove')
  message.success('图片已删除')
}
</script>

<style scoped>
.image-upload {
  width: 100%;
}

.image-preview {
  position: relative;
  display: inline-block;
}

.image-preview:hover .image-actions {
  opacity: 1;
}

.image-actions {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  padding: 12px;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.7), transparent);
  display: flex;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.upload-area {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 24px;
  cursor: pointer;
  transition: all 0.3s;
  color: #666;
}

.upload-area:hover {
  color: #0891b2;
}

html.dark .upload-area {
  color: #94a3b8;
}

html.dark .upload-area:hover {
  color: #38bdf8;
}

.image-upload :deep(.n-upload-dragger) {
  padding: 0;
}

.image-upload :deep(.n-upload-trigger) {
  width: 100%;
}
</style>

