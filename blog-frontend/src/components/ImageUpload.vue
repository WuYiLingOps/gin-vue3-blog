<template>
  <div
    class="image-upload"
    ref="rootRef"
    @mouseenter="handleMouseEnter"
    @mouseleave="handleMouseLeave"
  >
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

      <!-- 粘贴提示：悬停时允许 Ctrl+V 替换图片 -->
      <div v-if="isHover && isVisible" class="paste-hint">
        <n-text depth="3" style="font-size: 12px">悬停时按 Ctrl+V 可替换图片</n-text>
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
        <div
          ref="pasteTargetRef"
          class="upload-area"
          :class="{ paste_active: isHover && isVisible }"
          :style="{ width: width + 'px', height: height + 'px' }"
          tabindex="0"
          @paste="handlePaste"
        >
          <n-icon size="48" :component="CloudUploadOutline" />
          <n-text style="margin-top: 12px; display: block">
            点击或拖拽图片上传
          </n-text>
          <n-text depth="3" style="font-size: 12px; margin-top: 8px; display: block">
            支持 jpg、png、gif 格式，文件大小不超过 {{ maxSizeMB }}MB
          </n-text>
          <n-text depth="3" style="font-size: 12px; margin-top: 6px; display: block">
            悬停此区域可 Ctrl+V 粘贴上传（支持截图/复制图片）
          </n-text>
          <n-text v-if="!isVisible" depth="3" style="font-size: 12px; margin-top: 6px; display: block">
            组件不在可视区域时不响应粘贴
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
import { ref, watch, onBeforeUnmount, onMounted } from 'vue'
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
const uploading = ref(false)
const rootRef = ref<HTMLElement | null>(null)
const pasteTargetRef = ref<HTMLElement | null>(null)
const isHover = ref(false)
const isVisible = ref(true)
let io: IntersectionObserver | null = null

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

function handleMouseEnter() {
  isHover.value = true
  // 让上传区域更容易接收到粘贴事件（需要聚焦）。预览状态无此 ref，但有全局 paste 兜底。
  pasteTargetRef.value?.focus?.()
}

function handleMouseLeave() {
  isHover.value = false
}

function validateImageFile(file: File): boolean {
  if (!file.type.startsWith('image/')) {
    message.error('只能上传图片文件')
    return false
  }

  const maxSize = props.maxSizeMB * 1024 * 1024
  if (file.size > maxSize) {
    message.error(`图片大小不能超过 ${props.maxSizeMB}MB`)
    return false
  }

  return true
}

async function uploadAndSet(file: File) {
  if (uploading.value) return
  if (!validateImageFile(file)) return

  uploading.value = true
  try {
    const result = await uploadImage(file)
    if (result.data?.url) {
      imageUrl.value = result.data.url
      emit('update:modelValue', result.data.url)
      emit('success', result.data.url)
      message.success('图片上传成功')
    } else {
      message.error('上传失败：未获取到图片地址')
    }
  } catch (error: any) {
    console.error('粘贴上传失败:', error)
    message.error(error.message || '上传失败，请重试')
  } finally {
    uploading.value = false
  }
}

async function handlePaste(event: ClipboardEvent) {
  if (!isHover.value || !isVisible.value) return
  if (uploading.value) return

  const items = event.clipboardData?.items
  if (!items || items.length === 0) return

  // 仅处理剪贴板中的图片（如截图、复制图片）
  const imageItem = Array.from(items).find((it) => it.type?.startsWith('image/'))
  if (!imageItem) return

  const file = imageItem.getAsFile()
  if (!file) {
    message.error('读取剪贴板图片失败')
    return
  }

  event.preventDefault()
  await uploadAndSet(file)
}

function handleGlobalPaste(event: ClipboardEvent) {
  // 预览状态下无法保证有聚焦区域，因此用全局粘贴兜底；但必须悬停才处理，避免误触
  if (!isHover.value || !isVisible.value) return
  void handlePaste(event)
}

onMounted(() => {
  window.addEventListener('paste', handleGlobalPaste)

  // 组件可见性：只有在视口内才允许响应粘贴（双保险）
  if (!('IntersectionObserver' in window)) {
    isVisible.value = true
    return
  }

  io = new IntersectionObserver(
    (entries) => {
      const entry = entries[0]
      // intersectionRatio 有些浏览器/场景会抖动，这里用 isIntersecting + ratio 双判断更稳
      isVisible.value = !!entry && (entry.isIntersecting || entry.intersectionRatio > 0)
    },
    {
      root: null,
      // 只要进入视口就算可见
      threshold: [0, 0.01]
    }
  )

  if (rootRef.value) {
    io.observe(rootRef.value)
  }
})

onBeforeUnmount(() => {
  window.removeEventListener('paste', handleGlobalPaste)
  if (io) {
    io.disconnect()
    io = null
  }
})

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
  outline: none;
}

.upload-area.paste_active {
  box-shadow: 0 0 0 2px rgba(8, 145, 178, 0.35) inset;
}

.paste-hint {
  margin-top: 8px;
  text-align: center;
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

