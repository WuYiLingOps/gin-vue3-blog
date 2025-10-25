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
      :custom-request="customRequest"
      :show-file-list="false"
      accept="image/*"
      @before-upload="handleBeforeUpload"
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
import { ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import type { UploadFileInfo, UploadCustomRequestOptions } from 'naive-ui'
import { CloudUploadOutline, EyeOutline, TrashOutline } from '@vicons/ionicons5'
import { uploadImage } from '@/api/upload'

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

const imageUrl = ref(props.modelValue)
const showPreview = ref(false)

// 监听外部传入的值变化
watch(() => props.modelValue, (newVal) => {
  imageUrl.value = newVal
})

function handleBeforeUpload(data: { file: UploadFileInfo }) {
  const file = data.file.file
  
  if (!file) {
    message.error('文件读取失败')
    return false
  }
  
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

// 使用自定义上传请求
async function customRequest(options: UploadCustomRequestOptions) {
  const { file, onFinish, onError } = options
  
  try {
    const result = await uploadImage(file.file as File)
    
    if (result.data?.url) {
      imageUrl.value = result.data.url
      emit('update:modelValue', result.data.url)
      emit('success', result.data.url)
      message.success('图片上传成功')
      onFinish()
    }
  } catch (error: any) {
    console.error('上传失败:', error)
    message.error(error.message || '上传失败，请重试')
    onError()
  }
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

