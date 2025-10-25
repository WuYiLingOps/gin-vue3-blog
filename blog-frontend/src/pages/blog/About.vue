<template>
  <div class="about-page">
    <n-spin :show="loading">
      <n-card>
        <div class="about-content">
          <div class="about-header">
            <n-avatar :size="120" :src="settings.about_avatar || undefined">
              {{ settings.about_avatar ? '' : '头像' }}
            </n-avatar>
            <h1>关于我</h1>
          </div>

          <n-divider />

          <div class="about-section">
            <h2>{{ settings.about_title || '👋 你好' }}</h2>
            <p>
              {{ settings.about_intro || '欢迎来到我的个人博客！' }}
            </p>
          </div>

          <div v-if="skills.length > 0" class="about-section">
            <h2>💻 技术栈</h2>
            <n-space>
              <n-tag v-for="skill in skills" :key="skill" type="info">{{ skill }}</n-tag>
            </n-space>
          </div>

          <div class="about-section">
            <h2>📧 联系方式</h2>
            <n-space vertical>
              <div v-if="settings.about_email">
                <n-icon :component="MailOutline" />
                <span style="margin-left: 8px">Email: {{ settings.about_email }}</span>
              </div>
              <div v-if="settings.about_github">
                <n-icon :component="LogoGithub" />
                <span style="margin-left: 8px">GitHub: {{ settings.about_github }}</span>
              </div>
            </n-space>
          </div>

          <div v-if="settings.about_site_intro" class="about-section">
            <h2>🎯 关于本站</h2>
            <p style="white-space: pre-line">{{ settings.about_site_intro }}</p>
          </div>
        </div>
      </n-card>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { MailOutline, LogoGithub } from '@vicons/ionicons5'
import { getAboutSettings } from '@/api/setting'
import type { AboutSettings } from '@/api/setting'

const loading = ref(false)
const settings = reactive<AboutSettings>({})

const skills = computed(() => {
  if (!settings.about_skills) return []
  try {
    return JSON.parse(settings.about_skills)
  } catch (error) {
    return []
  }
})

onMounted(async () => {
  try {
    loading.value = true
    const response = await getAboutSettings()
    
    // 响应拦截器返回的是完整的响应对象，需要从 data 中获取实际数据
    const data = response.data || response
    
    console.log('About settings loaded:', data)
    Object.assign(settings, data)
  } catch (error) {
    console.error('Failed to load about settings:', error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.about-page {
  max-width: 900px;
  margin: 0 auto;
}

.about-content {
  padding: 24px;
}

.about-header {
  text-align: center;
  margin-bottom: 32px;
}

.about-header h1 {
  margin-top: 16px;
  font-size: 32px;
}

.about-section {
  margin: 32px 0;
}

.about-section h2 {
  font-size: 24px;
  margin-bottom: 16px;
}

.about-section p {
  line-height: 1.8;
  color: #666;
  margin: 12px 0;
}
</style>

