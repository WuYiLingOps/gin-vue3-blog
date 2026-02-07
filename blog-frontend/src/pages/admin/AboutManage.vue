<!--
 * @ProjectName: go-vue3-blog
 * @FileName: AboutManage.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于我管理页面组件，提供关于我内容的编辑功能
 -->
<template>
  <div class="about-manage-page">
    <n-card title="关于我信息管理">
      <n-spin :show="loading">
        <n-form
          :model="formData"
          label-placement="top"
        >
          <n-form-item label="关于我内容（支持 Markdown 格式）" path="content">
            <markdown-editor 
              v-model="formData.content" 
              :height="isMobile ? '500px' : '700px'" 
              :subfield="!isMobile"
              :mode="isSmallScreen ? 'edit' : 'editable'"
            />
            <template #feedback>
              <n-text depth="3" style="font-size: 12px; margin-top: 8px; display: block">
                提示：支持 Markdown 语法，内容将显示在"关于我"页面中
              </n-text>
            </template>
          </n-form-item>

          <n-space :vertical="isMobile" style="margin-top: 24px">
            <n-button 
              type="primary" 
              :size="isMobile ? 'medium' : 'large'" 
              :block="isMobile" 
              :loading="submitting" 
              @click="handleSubmit"
            >
              保存
            </n-button>
            <n-button 
              :size="isMobile ? 'medium' : 'large'" 
              :block="isMobile" 
              @click="handleReset"
            >
              重置
            </n-button>
          </n-space>
        </n-form>
      </n-spin>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'
import { getAboutInfo, updateAboutInfo } from '@/api/setting'
import MarkdownEditor from '@/components/MarkdownEditor.vue'

const message = useMessage()
const loading = ref(false)
const submitting = ref(false)
const isMobile = ref(false)
const isSmallScreen = ref(false)

const formData = reactive({
  content: ''
})

const originalContent = ref('')

// 检测移动设备
function checkMobile() {
  const width = window.innerWidth
  isMobile.value = width <= 1100
  isSmallScreen.value = width <= 600
}

// 获取关于我信息
async function fetchAboutInfo() {
  try {
    loading.value = true
    const res = await getAboutInfo()
    if (res.data) {
      formData.content = res.data.content || ''
      originalContent.value = res.data.content || ''
    }
  } catch (error: any) {
    message.error(error.message || '获取关于我信息失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 保存
async function handleSubmit() {
  try {
    submitting.value = true
    await updateAboutInfo(formData.content)
    message.success('保存成功')
    originalContent.value = formData.content
  } catch (error: any) {
    message.error(error.message || '保存失败')
    console.error(error)
  } finally {
    submitting.value = false
  }
}

// 重置
function handleReset() {
  formData.content = originalContent.value
  message.info('已重置为上次保存的内容')
}

onMounted(() => {
  checkMobile()
  fetchAboutInfo()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.about-manage-page {
  padding: 20px;
}

@media (max-width: 1100px) {
  .about-manage-page {
    padding: 12px;
  }
}
</style>
