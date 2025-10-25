<template>
  <div class="avatar-upload">
    <n-upload
      :action="uploadAction"
      :headers="uploadHeaders"
      :max="1"
      :show-file-list="false"
      :on-before-upload="beforeUpload"
      :on-finish="handleFinish"
      :on-error="handleError"
      @change="handleChange"
    >
      <n-avatar
        round
        :size="size"
        :src="currentAvatar"
        class="avatar-preview"
      >
        {{ defaultText }}
      </n-avatar>
      <div class="upload-overlay">
        <n-icon size="24" :component="CameraOutline" />
        <span>{{ uploading ? '上传中...' : '点击上传' }}</span>
      </div>
    </n-upload>
    <div v-if="tip" class="upload-tip">{{ tip }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { CameraOutline } from '@vicons/ionicons5'
import { useAuthStore } from '@/stores'
import type { UploadFileInfo } from 'naive-ui'

interface Props {
  modelValue?: string
  size?: number
  tip?: string
  defaultText?: string
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'success', url: string): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  size: 100,
  tip: '支持 jpg、png、gif 格式，文件大小不超过 5MB',
  defaultText: ''
})

const emit = defineEmits<Emits>()

const message = useMessage()
const authStore = useAuthStore()

const uploading = ref(false)
const currentAvatar = ref(props.modelValue)

// 监听 modelValue 变化
watch(
  () => props.modelValue,
  (newValue) => {
    currentAvatar.value = newValue
  }
)

const uploadAction = computed(() => {
  const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
  return `${baseURL}/api/upload/avatar`
})

const uploadHeaders = computed(() => {
  return {
    Authorization: `Bearer ${authStore.token}`
  }
})

function beforeUpload(data: { file: UploadFileInfo; fileList: UploadFileInfo[] }) {
  const file = data.file.file
  if (!file) return false

  // 检查文件类型
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
  if (!allowedTypes.includes(file.type)) {
    message.error('只支持上传 jpg、png、gif、webp 格式的图片')
    return false
  }

  // 检查文件大小（5MB）
  const maxSize = 5 * 1024 * 1024
  if (file.size > maxSize) {
    message.error('图片大小不能超过 5MB')
    return false
  }

  uploading.value = true
  return true
}

function handleChange() {
  // 可以在这里添加额外的逻辑
}

function handleFinish({ event }: { file: UploadFileInfo; event?: ProgressEvent }) {
  uploading.value = false

  try {
    const response = JSON.parse((event?.target as XMLHttpRequest).response)
    
    if (response.code === 200 && response.data?.url) {
      const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
      const fullURL = baseURL + response.data.url
      
      currentAvatar.value = fullURL
      emit('update:modelValue', fullURL)
      emit('success', fullURL)
      message.success('头像上传成功')
    } else {
      message.error(response.message || '上传失败')
    }
  } catch (error) {
    console.error('Upload response parse error:', error)
    message.error('上传失败')
  }
}

function handleError({ event }: { file: UploadFileInfo; event?: ProgressEvent }) {
  uploading.value = false
  console.error('Upload error:', event)
  message.error('上传失败，请重试')
}
</script>

<style scoped>
.avatar-upload {
  display: inline-block;
  text-align: center;
}

.avatar-upload :deep(.n-upload) {
  position: relative;
  display: inline-block;
}

.avatar-preview {
  cursor: pointer;
  transition: all 0.3s;
}

.upload-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.5);
  color: white;
  opacity: 0;
  transition: opacity 0.3s;
  border-radius: 50%;
  cursor: pointer;
}

.avatar-upload :deep(.n-upload:hover) .upload-overlay {
  opacity: 1;
}

.upload-overlay span {
  margin-top: 4px;
  font-size: 12px;
}

.upload-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #999;
}
</style>

