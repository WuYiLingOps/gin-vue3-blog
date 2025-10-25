<template>
  <div class="setting-manage-page">
    <h1>关于页面配置</h1>

    <n-spin :show="loading">
      <n-card>
        <n-form ref="formRef" :model="formData" label-placement="left" label-width="120px">
          <n-form-item label="标题">
            <n-input v-model:value="formData.about_title" placeholder="例如：👋 你好" />
          </n-form-item>

          <n-form-item label="个人简介">
            <n-input
              v-model:value="formData.about_intro"
              type="textarea"
              :rows="3"
              placeholder="介绍一下自己"
            />
          </n-form-item>

          <n-form-item label="个人头像">
            <n-space vertical>
              <avatar-upload
                v-model="formData.about_avatar"
                :size="100"
                default-text="头像"
                @success="handleAvatarSuccess"
              />
              <n-text depth="3" style="font-size: 12px">点击头像上传新图片</n-text>
            </n-space>
          </n-form-item>

          <n-form-item label="技术栈">
            <n-dynamic-tags v-model:value="skills" />
          </n-form-item>

          <n-form-item label="联系邮箱">
            <n-input
              v-model:value="formData.about_email"
              placeholder="your-email@example.com"
            />
          </n-form-item>

          <n-form-item label="GitHub">
            <n-input
              v-model:value="formData.about_github"
              placeholder="github.com/yourname"
            />
          </n-form-item>

          <n-form-item label="关于本站">
            <n-input
              v-model:value="formData.about_site_intro"
              type="textarea"
              :rows="6"
              placeholder="介绍一下本站"
            />
          </n-form-item>

          <n-form-item>
            <n-space>
              <n-button type="primary" size="large" :loading="submitting" @click="handleSubmit">
                保存配置
              </n-button>
              <n-button size="large" @click="loadSettings">重置</n-button>
            </n-space>
          </n-form-item>
        </n-form>
      </n-card>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { useMessage } from 'naive-ui'
import type { FormInst } from 'naive-ui'
import { getAboutSettings, updateAboutSettings } from '@/api/setting'
import type { AboutSettings } from '@/api/setting'
import AvatarUpload from '@/components/AvatarUpload.vue'

const message = useMessage()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const submitting = ref(false)

const formData = reactive<AboutSettings>({
  about_title: '',
  about_intro: '',
  about_avatar: '',
  about_skills: '[]',
  about_email: '',
  about_github: '',
  about_site_intro: ''
})

const skills = ref<string[]>([])

// 监听 skills 变化，同步到 formData
watch(
  skills,
  (newSkills) => {
    formData.about_skills = JSON.stringify(newSkills)
  },
  { deep: true }
)

onMounted(() => {
  loadSettings()
})

async function loadSettings() {
  try {
    loading.value = true
    const response = await getAboutSettings()
    
    console.log('API response:', response)
    
    // 响应拦截器返回的是完整的响应对象，需要从 data 中获取实际数据
    const data = response.data || response
    
    console.log('Settings data:', data)
    
    formData.about_title = data.about_title || ''
    formData.about_intro = data.about_intro || ''
    formData.about_avatar = data.about_avatar || ''
    formData.about_skills = data.about_skills || '[]'
    formData.about_email = data.about_email || ''
    formData.about_github = data.about_github || ''
    formData.about_site_intro = data.about_site_intro || ''

    // 解析技能标签
    try {
      skills.value = JSON.parse(data.about_skills || '[]')
    } catch (error) {
      console.error('Failed to parse skills:', error)
      skills.value = []
    }
  } catch (error) {
    console.error('Failed to load settings:', error)
    message.error('加载配置失败')
  } finally {
    loading.value = false
  }
}

function handleAvatarSuccess(url: string) {
  formData.about_avatar = url
}

async function handleSubmit() {
  try {
    submitting.value = true

    // 确保技能标签转换为 JSON
    formData.about_skills = JSON.stringify(skills.value)

    const submitData: Record<string, string> = {}
    Object.keys(formData).forEach((key) => {
      const value = formData[key as keyof AboutSettings]
      if (value !== undefined && value !== null) {
        submitData[key] = String(value)
      }
    })

    console.log('Submitting data:', submitData)

    await updateAboutSettings(submitData)
    message.success('配置保存成功')

    // 重新加载配置以确认更新
    await loadSettings()
  } catch (error) {
    console.error('Failed to save settings:', error)
    message.error('保存配置失败')
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.setting-manage-page {
  padding: 0;
  max-width: 1000px;
}

.setting-manage-page h1 {
  margin: 0 0 24px 0;
  font-size: 24px;
  font-weight: 700;
  color: #1a202c;
}

html.dark .setting-manage-page h1 {
  color: #e5e5e5;
}

.setting-manage-page :deep(.n-card) {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(8, 145, 178, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  transition: all 0.3s;
}

html.dark .setting-manage-page :deep(.n-card) {
  background: rgba(30, 41, 59, 0.85);
  border: 1px solid rgba(255, 255, 255, 0.08);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.setting-manage-page :deep(.n-form-item) {
  margin-bottom: 24px;
}
</style>

