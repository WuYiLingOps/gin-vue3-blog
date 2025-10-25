<template>
  <div class="about-page">
    <n-spin :show="loading">
      <div class="about-container">
        <!-- 头部卡片 -->
        <div class="about-header-card">
          <div class="header-bg"></div>
          <div class="header-content">
            <div class="about-avatar">
              <img 
                v-if="settings.about_avatar" 
                :src="settings.about_avatar" 
                alt="头像"
                @error="handleImageError"
              />
              <div v-else class="avatar-placeholder">
                <span>👤</span>
              </div>
            </div>
            <h1 class="about-title">{{ settings.about_title || '关于我' }}</h1>
            <p class="about-subtitle">{{ settings.about_intro || '欢迎来到我的个人博客' }}</p>
          </div>
        </div>

        <!-- 富文本内容卡片 -->
        <div v-if="settings.about_content" class="content-card">
          <markdown-preview :content="settings.about_content" />
        </div>

        <!-- 技能标签卡片 -->
        <div v-if="skills.length > 0" class="skills-card">
          <h2 class="card-title">
            <span class="title-icon">💻</span>
            技术栈
          </h2>
          <div class="skills-grid">
            <n-tag 
              v-for="skill in skills" 
              :key="skill" 
              :bordered="false"
              size="large"
              class="skill-tag"
            >
              {{ skill }}
            </n-tag>
          </div>
        </div>

        <!-- 联系方式卡片 -->
        <div class="contact-card">
          <h2 class="card-title">
            <span class="title-icon">📧</span>
            联系方式
          </h2>
          <div class="contact-list">
            <a v-if="settings.about_email" :href="`mailto:${settings.about_email}`" class="contact-item">
              <n-icon :component="MailOutline" size="20" />
              <span>{{ settings.about_email }}</span>
            </a>
            <a v-if="settings.about_github" :href="settings.about_github" target="_blank" class="contact-item">
              <n-icon :component="LogoGithub" size="20" />
              <span>{{ settings.about_github }}</span>
            </a>
          </div>
        </div>
      </div>
    </n-spin>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { MailOutline, LogoGithub } from '@vicons/ionicons5'
import { getAboutSettings } from '@/api/setting'
import type { AboutSettings } from '@/api/setting'
import MarkdownPreview from '@/components/MarkdownPreview.vue'

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
    const settingsData: AboutSettings = ('data' in data ? data.data : data) as AboutSettings
    
    console.log('About settings loaded:', settingsData)
    console.log('Avatar URL:', settingsData?.about_avatar)
    Object.assign(settings, settingsData)
  } catch (error) {
    console.error('Failed to load about settings:', error)
  } finally {
    loading.value = false
  }
})
</script>

<style scoped>
.about-page {
  max-width: 1000px;
  margin: 0 auto;
  padding: 40px 20px;
}

.about-container {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

/* 头部卡片 */
.about-header-card {
  position: relative;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 24px;
  overflow: hidden;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
}

html.dark .about-header-card {
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

.header-bg {
  height: 200px;
  background: linear-gradient(135deg, #10b981 0%, #06b6d4 100%);
  position: relative;
}

.header-bg::after {
  content: '';
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 100px;
  background: linear-gradient(to bottom, transparent, rgba(255, 255, 255, 0.9));
}

html.dark .header-bg::after {
  background: linear-gradient(to bottom, transparent, rgba(30, 41, 59, 0.8));
}

.header-content {
  position: relative;
  text-align: center;
  padding: 0 40px 40px;
  margin-top: -80px;
}

.about-avatar {
  width: 160px;
  height: 160px;
  margin: 0 auto 24px;
  border-radius: 50%;
  overflow: hidden;
  border: 6px solid white;
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.15);
  transition: all 0.3s ease;
}

html.dark .about-avatar {
  border-color: rgba(30, 41, 59, 0.9);
}

.about-avatar:hover {
  transform: scale(1.05) rotate(5deg);
  box-shadow: 0 12px 32px rgba(16, 185, 129, 0.3);
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
  background: linear-gradient(135deg, #10b981 0%, #06b6d4 100%);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 64px;
}

.about-title {
  font-size: 36px;
  font-weight: 800;
  margin: 0 0 12px 0;
  background: linear-gradient(135deg, #10b981 0%, #06b6d4 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.about-subtitle {
  font-size: 16px;
  color: #64748b;
  margin: 0;
  line-height: 1.6;
}

html.dark .about-subtitle {
  color: #94a3b8;
}

/* 内容卡片 */
.content-card,
.skills-card,
.contact-card {
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 16px;
  padding: 32px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

html.dark .content-card,
html.dark .skills-card,
html.dark .contact-card {
  background: rgba(30, 41, 59, 0.8);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.3);
}

.content-card:hover,
.skills-card:hover,
.contact-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 24px rgba(16, 185, 129, 0.15);
}

.card-title {
  font-size: 24px;
  font-weight: 700;
  margin: 0 0 24px 0;
  color: #1a202c;
  display: flex;
  align-items: center;
  gap: 12px;
}

html.dark .card-title {
  color: #e5e5e5;
}

.title-icon {
  font-size: 28px;
}

/* 技能网格 */
.skills-grid {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
}

.skill-tag {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1) 0%, rgba(6, 182, 212, 0.1) 100%) !important;
  color: #10b981 !important;
  font-weight: 500;
  padding: 8px 16px;
  transition: all 0.2s ease;
}

.skill-tag:hover {
  transform: translateY(-2px);
  box-shadow: 0 4px 12px rgba(16, 185, 129, 0.2);
}

/* 联系方式列表 */
.contact-list {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.contact-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 16px;
  background: rgba(16, 185, 129, 0.05);
  border-radius: 12px;
  color: #1a202c;
  text-decoration: none;
  transition: all 0.2s ease;
}

html.dark .contact-item {
  background: rgba(16, 185, 129, 0.1);
  color: #e5e5e5;
}

.contact-item:hover {
  background: rgba(16, 185, 129, 0.1);
  transform: translateX(8px);
}

html.dark .contact-item:hover {
  background: rgba(16, 185, 129, 0.15);
}

.contact-item span {
  font-size: 15px;
}

/* 响应式 */
@media (max-width: 768px) {
  .about-page {
    padding: 24px 16px;
  }

  .header-content {
    padding: 0 24px 32px;
  }

  .about-avatar {
    width: 120px;
    height: 120px;
  }

  .about-title {
    font-size: 28px;
  }

  .content-card,
  .skills-card,
  .contact-card {
    padding: 24px;
  }

  .card-title {
    font-size: 20px;
  }
}
</style>

