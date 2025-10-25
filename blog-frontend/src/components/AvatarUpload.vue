<template>
  <div class="avatar-upload">
    <div class="avatar-container" @click="triggerFileInput">
      <!-- 自定义圆形头像 -->
      <div 
        class="avatar-preview" 
        :style="{ 
          width: size + 'px', 
          height: size + 'px',
          backgroundImage: currentAvatar ? `url(${currentAvatar})` : 'none'
        }"
      >
        <!-- 显示默认文字（当没有头像时） -->
        <span v-if="!currentAvatar" class="avatar-text">{{ defaultText }}</span>
      </div>
      
      <!-- 悬停遮罩层 -->
      <div class="upload-overlay" :style="{ width: size + 'px', height: size + 'px' }">
        <!-- 相机图标 -->
        <svg class="camera-icon" viewBox="0 0 512 512" width="24" height="24">
          <path fill="currentColor" d="M512 144v288c0 26.5-21.5 48-48 48H48c-26.5 0-48-21.5-48-48V144c0-26.5 21.5-48 48-48h88l12.3-32.9c7-18.7 24.9-31.1 44.9-31.1h125.5c20 0 37.9 12.4 44.9 31.1L376 96h88c26.5 0 48 21.5 48 48zM376 288c0-66.2-53.8-120-120-120s-120 53.8-120 120 53.8 120 120 120 120-53.8 120-120zm-32 0c0 48.5-39.5 88-88 88s-88-39.5-88-88 39.5-88 88-88 88 39.5 88 88z"/>
        </svg>
        <span class="upload-text">{{ uploading ? '上传中...' : '点击上传' }}</span>
      </div>
    </div>
    
    <!-- 隐藏的文件输入框 -->
    <input
      ref="fileInputRef"
      type="file"
      accept="image/jpeg,image/jpg,image/png,image/gif,image/webp"
      style="display: none"
      @change="handleFileChange"
    />
    
    <div v-if="tip" class="upload-tip">{{ tip }}</div>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import { useAuthStore } from '@/stores'

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
const fileInputRef = ref<HTMLInputElement | null>(null)

// 监听 modelValue 变化
watch(
  () => props.modelValue,
  (newValue) => {
    currentAvatar.value = newValue
  }
)

// 触发文件选择
function triggerFileInput() {
  if (!uploading.value) {
    fileInputRef.value?.click()
  }
}

// 处理文件选择
async function handleFileChange(event: Event) {
  const target = event.target as HTMLInputElement
  const file = target.files?.[0]
  
  if (!file) return
  
  // 检查文件类型
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif', 'image/webp']
  if (!allowedTypes.includes(file.type)) {
    message.error('只支持上传 jpg、png、gif、webp 格式的图片')
    target.value = '' // 重置 input
    return
  }
  
  // 检查文件大小（5MB）
  const maxSize = 5 * 1024 * 1024
  if (file.size > maxSize) {
    message.error('图片大小不能超过 5MB')
    target.value = '' // 重置 input
    return
  }
  
  // 开始上传
  await uploadFile(file)
  target.value = '' // 重置 input，允许重复上传同一文件
}

// 上传文件
async function uploadFile(file: File) {
  uploading.value = true
  
  try {
    const formData = new FormData()
    formData.append('file', file)
    
    const baseURL = import.meta.env.VITE_API_BASE_URL || 'http://localhost:8080'
    const uploadURL = `${baseURL}/api/upload/avatar`
    
    const response = await fetch(uploadURL, {
      method: 'POST',
      headers: {
        'Authorization': `Bearer ${authStore.token}`
      },
      body: formData
    })
    
    const result = await response.json()
    
    if (result.code === 200 && result.data?.url) {
      const fullURL = baseURL + result.data.url
      
      currentAvatar.value = fullURL
      emit('update:modelValue', fullURL)
      emit('success', fullURL)
      message.success('头像上传成功')
    } else {
      message.error(result.message || '上传失败')
    }
  } catch (error) {
    console.error('Upload error:', error)
    message.error('上传失败，请重试')
  } finally {
    uploading.value = false
  }
}
</script>

<style scoped>
.avatar-upload {
  display: inline-block;
  text-align: center;
}

.avatar-container {
  position: relative;
  display: inline-block;
  cursor: pointer;
}

.avatar-preview {
  border-radius: 50%;
  background-color: #f0f0f0;
  background-size: cover;
  background-position: center;
  background-repeat: no-repeat;
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  transition: all 0.3s;
  border: 2px solid #e0e0e0;
}

.avatar-text {
  font-size: 36px;
  font-weight: bold;
  color: #999;
  user-select: none;
}

.upload-overlay {
  position: absolute;
  top: 0;
  left: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  opacity: 0;
  transition: opacity 0.3s;
  border-radius: 50%;
  pointer-events: none;
}

.avatar-container:hover .upload-overlay {
  opacity: 1;
}

.camera-icon {
  margin-bottom: 4px;
}

.upload-text {
  font-size: 12px;
  white-space: nowrap;
}

.upload-tip {
  margin-top: 8px;
  font-size: 12px;
  color: #999;
}

/* 禁用状态 */
.avatar-container.uploading {
  cursor: not-allowed;
  opacity: 0.7;
}
</style>

