<!--
  项目名称：blog-frontend
  文件名称：MultiImageUpload.vue
  创建时间：2026-02-01 20:03:19

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：多图片上传组件，支持批量上传多张图片，显示已上传图片列表，支持删除单张图片，最多上传指定数量的图片。
-->
<template>
  <div class="multi-image-upload">
    <n-space :size="12">
      <!-- 已上传的图片列表 -->
      <div v-for="(url, index) in imageList" :key="index" class="image-item">
        <n-image
          :src="url"
          width="120"
          height="120"
          object-fit="cover"
          class="preview-image"
        />
        <div class="image-mask">
          <n-space>
            <n-button size="tiny" circle @click="handleRemove(index)">
              <template #icon>
                <n-icon :component="TrashOutline" />
              </template>
            </n-button>
          </n-space>
        </div>
      </div>

      <!-- 上传按钮 -->
      <n-upload
        v-if="imageList.length < maxCount"
        :custom-request="customRequest"
        :show-file-list="false"
        accept="image/*"
        :multiple="false"
        @before-upload="handleBeforeUpload"
      >
        <div class="upload-button">
          <n-icon size="32" :component="CloudUploadOutline" />
          <n-text style="font-size: 12px; margin-top: 8px; display: block">
            上传图片
          </n-text>
          <n-text depth="3" style="font-size: 11px; display: block">
            {{ imageList.length }}/{{ maxCount }}
          </n-text>
        </div>
      </n-upload>
    </n-space>

    <n-text depth="3" style="font-size: 12px; margin-top: 8px; display: block">
      支持 jpg、png、gif 格式，单张图片不超过 5MB，最多上传 {{ maxCount }} 张
    </n-text>
  </div>
</template>

<script setup lang="ts">
import { ref, watch } from 'vue'
import { useMessage } from 'naive-ui'
import type { UploadFileInfo, UploadCustomRequestOptions } from 'naive-ui'
import { CloudUploadOutline, TrashOutline } from '@vicons/ionicons5'
import { uploadImage } from '@/api/upload'

interface Props {
  modelValue?: string | string[]
  maxCount?: number
}

interface Emits {
  (e: 'update:modelValue', value: string): void
  (e: 'success', urls: string[]): void
}

const props = withDefaults(defineProps<Props>(), {
  modelValue: '',
  maxCount: 9
})

const emit = defineEmits<Emits>()
const message = useMessage()

const imageList = ref<string[]>([])

// 监听外部传入的值变化
watch(
  () => props.modelValue,
  (newVal) => {
    if (typeof newVal === 'string' && newVal) {
      try {
        imageList.value = JSON.parse(newVal)
      } catch {
        imageList.value = []
      }
    } else if (Array.isArray(newVal)) {
      imageList.value = newVal
    } else {
      imageList.value = []
    }
  },
  { immediate: true }
)

// 上传前验证
function handleBeforeUpload(data: { file: UploadFileInfo }): boolean {
  const file = data.file.file
  if (!file) return false

  // 验证文件类型
  const validTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
  if (!validTypes.includes(file.type)) {
    message.error('只支持 jpg、png、gif 格式的图片')
    return false
  }

  // 验证文件大小（5MB）
  const maxSize = 5 * 1024 * 1024
  if (file.size > maxSize) {
    message.error('图片大小不能超过 5MB')
    return false
  }

  // 验证数量
  if (imageList.value.length >= props.maxCount) {
    message.error(`最多只能上传 ${props.maxCount} 张图片`)
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
      imageList.value.push(result.data.url)
      const jsonStr = JSON.stringify(imageList.value)
      emit('update:modelValue', jsonStr)
      emit('success', imageList.value)
      message.success('图片上传成功')
      onFinish()
    }
  } catch (error: any) {
    console.error('上传失败:', error)
    message.error(error.message || '上传失败，请重试')
    onError()
  }
}

// 删除图片
function handleRemove(index: number) {
  imageList.value.splice(index, 1)
  const jsonStr = JSON.stringify(imageList.value)
  emit('update:modelValue', jsonStr)
  emit('success', imageList.value)
}
</script>

<style scoped>
.multi-image-upload {
  width: 100%;
}

.image-item {
  position: relative;
  width: 120px;
  height: 120px;
  border-radius: 8px;
  overflow: hidden;
  border: 1px solid #e5e7eb;
}

html.dark .image-item {
  border-color: #374151;
}

.preview-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-mask {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  opacity: 0;
  transition: opacity 0.3s;
}

.image-item:hover .image-mask {
  opacity: 1;
}

.upload-button {
  width: 120px;
  height: 120px;
  border: 2px dashed #d1d5db;
  border-radius: 8px;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  cursor: pointer;
  transition: all 0.3s;
  background: #f9fafb;
}

.upload-button:hover {
  border-color: #0891b2;
  background: #f0f9ff;
}

html.dark .upload-button {
  background: #1a1a1a;
  border-color: #374151;
}

html.dark .upload-button:hover {
  border-color: #38bdf8;
  background: #1e293b;
}
</style>

