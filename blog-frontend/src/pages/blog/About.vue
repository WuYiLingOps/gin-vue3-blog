<template>
  <div class="about-page">
    <n-spin :show="loading">
      <n-card>
        <div class="about-content">
          <div class="about-header">
            <!-- 自定义头像显示 -->
            <div class="about-avatar">
              <img 
                v-if="settings.about_avatar" 
                :src="settings.about_avatar" 
                alt="头像"
                @error="handleImageError"
              />
              <div v-else class="avatar-placeholder">
                <span>头像</span>
              </div>
            </div>
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

// 处理图片加载错误
function handleImageError(e: Event) {
  console.error('Avatar image load failed:', settings.about_avatar)
  // 图片加载失败时隐藏图片，显示占位符
  ;(e.target as HTMLImageElement).style.display = 'none'
}

onMounted(async () => {
  try {
    loading.value = true
    const response = await getAboutSettings()
    
    // 响应拦截器返回的是完整的响应对象，需要从 data 中获取实际数据
    const data = response.data || response
    
    console.log('About settings loaded:', data)
    console.log('Avatar URL:', data.about_avatar)
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

/* 自定义头像样式 */
.about-avatar {
  width: 120px;
  height: 120px;
  margin: 0 auto 16px;
  border-radius: 50%;
  overflow: hidden;
  border: 3px solid rgba(8, 145, 178, 0.2);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.about-avatar:hover {
  transform: scale(1.05);
  border-color: rgba(8, 145, 178, 0.4);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
}

.about-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
}

.avatar-placeholder {
  width: 100%;
  height: 100%;
  background: linear-gradient(135deg, #e0f7fa 0%, #b2ebf2 100%);
  display: flex;
  align-items: center;
  justify-content: center;
}

.avatar-placeholder span {
  font-size: 18px;
  font-weight: 600;
  color: #0891b2;
}

html.dark .about-avatar {
  border-color: rgba(56, 189, 248, 0.3);
}

html.dark .about-avatar:hover {
  border-color: rgba(56, 189, 248, 0.5);
}

html.dark .avatar-placeholder {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 100%);
}

html.dark .avatar-placeholder span {
  color: #38bdf8;
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

