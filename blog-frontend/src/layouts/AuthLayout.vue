<!--
 * @ProjectName: go-vue3-blog
 * @FileName: AuthLayout.vue
 * @CreateTime: 2026-02-02 11:40:46
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 认证布局组件，用于登录、注册等认证相关页面的布局容器
 -->
<template>
  <div class="auth-layout">
    <!-- 全局背景图组件 -->
    <GlobalBackground />

    <!-- 左侧装饰区 -->
    <div class="auth-brand-section">
      <div class="brand-content">
        <!-- Logo 区域 -->
        <div class="brand-logo">
          <div class="logo-wrapper">
            <div class="modern-logo">
              <div class="logo-icon">
                <svg viewBox="0 0 100 100" xmlns="http://www.w3.org/2000/svg">
                  <!-- 代码符号 </> -->
                  <text
                    x="50"
                    y="65"
                    font-size="48"
                    font-weight="bold"
                    text-anchor="middle"
                    fill="currentColor"
                  >
                    &lt;/&gt;
                  </text>
                </svg>
              </div>
            </div>
          </div>
          <h1 class="brand-title gradient-text">{{ siteName }}</h1>
        </div>

        <!-- Slogan -->
        <p class="brand-slogan">{{ displayedSlogan }}<span class="cursor">|</span></p>

        <!-- 装饰圆圈 -->
        <div class="decorative-elements">
          <div class="deco-circle deco-1"></div>
          <div class="deco-circle deco-2"></div>
          <div class="deco-circle deco-3"></div>
        </div>

        <!-- 功能图标 -->
        <div class="feature-icons">
          <div class="feature-item">
            <n-icon :component="CreateOutline" size="32" />
            <span>创作</span>
          </div>
          <div class="feature-item">
            <n-icon :component="ShareSocialOutline" size="32" />
            <span>分享</span>
          </div>
          <div class="feature-item">
            <n-icon :component="TrendingUpOutline" size="32" />
            <span>成长</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 右侧表单区 -->
    <div class="auth-form-section">
      <div class="auth-container">
        <div class="auth-content">
          <router-view />
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { CreateOutline, ShareSocialOutline, TrendingUpOutline } from '@vicons/ionicons5'
import GlobalBackground from '@/components/GlobalBackground.vue'
import { getPublicSettings } from '@/api/setting'

const siteName = ref('博客系统')
const displayedSlogan = ref('')
const fullSlogan = '记录思考，分享知识，连接世界'

// 打字机效果
const typeWriter = () => {
  let index = 0
  const timer = setInterval(() => {
    if (index < fullSlogan.length) {
      displayedSlogan.value += fullSlogan[index]
      index++
    } else {
      clearInterval(timer)
    }
  }, 150) // 每个字符间隔 150ms
}

onMounted(async () => {
  try {
    const res = await getPublicSettings()
    if (res.data?.site_name) {
      siteName.value = res.data.site_name
    }
  } catch (error) {
    console.error('获取网站配置失败:', error)
  }

  // 延迟 500ms 后开始打字机效果
  setTimeout(() => {
    typeWriter()
  }, 500)
})
</script>

<style scoped>
.auth-layout {
  min-height: 100vh;
  display: grid;
  grid-template-columns: 45% 55%;
  position: relative;
  overflow: hidden;
}

/* 左侧装饰区 */
.auth-brand-section {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
  position: relative;
  z-index: 1;
}

.brand-content {
  max-width: 480px;
  width: 100%;
  text-align: center;
}

/* Logo 区域 */
.brand-logo {
  margin-bottom: 32px;
}

.logo-wrapper {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 120px;
  height: 120px;
  margin-bottom: 24px;
  position: relative;
}

/* 现代简约 Logo */
.modern-logo {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.15);
  backdrop-filter: blur(20px);
  border-radius: 24px;
  border: 2px solid rgba(255, 255, 255, 0.3);
  box-shadow:
    0 8px 32px rgba(0, 0, 0, 0.15),
    inset 0 1px 0 rgba(255, 255, 255, 0.3);
  position: relative;
  overflow: hidden;
  transition: all 0.3s ease;
}

.modern-logo::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, rgba(8, 145, 178, 0.1), rgba(5, 150, 105, 0.1));
  opacity: 0;
  transition: opacity 0.3s ease;
}

.modern-logo:hover::before {
  opacity: 1;
}

.logo-icon {
  width: 70%;
  height: 70%;
  color: #fff;
  filter: drop-shadow(0 4px 12px rgba(0, 0, 0, 0.3));
  transition: transform 0.3s ease;
}

.modern-logo:hover .logo-icon {
  transform: scale(1.05);
}

.brand-title {
  font-size: 42px;
  font-weight: 700;
  margin: 0;
  letter-spacing: 2px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.brand-slogan {
  font-size: 18px;
  color: #fff;
  margin: 16px 0 48px 0;
  line-height: 1.8;
  font-weight: 400;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
  min-height: 32px;
}

.cursor {
  display: inline-block;
  margin-left: 2px;
  animation: blink 1s step-end infinite;
}

@keyframes blink {
  0%,
  50% {
    opacity: 1;
  }
  51%,
  100% {
    opacity: 0;
  }
}

/* 装饰圆圈 */
.decorative-elements {
  position: relative;
  height: 200px;
  margin-bottom: 48px;
}

.deco-circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.15;
  animation: float 6s ease-in-out infinite;
}

.deco-1 {
  width: 120px;
  height: 120px;
  background: linear-gradient(135deg, #0891b2, #059669);
  left: 10%;
  top: 20%;
  animation-delay: 0s;
}

.deco-2 {
  width: 80px;
  height: 80px;
  background: linear-gradient(135deg, #059669, #0891b2);
  right: 15%;
  top: 50%;
  animation-delay: 2s;
}

.deco-3 {
  width: 60px;
  height: 60px;
  background: linear-gradient(135deg, #0891b2, #10b981);
  left: 50%;
  bottom: 10%;
  animation-delay: 4s;
}

@keyframes float {
  0%,
  100% {
    transform: translateY(0) scale(1);
  }
  50% {
    transform: translateY(-20px) scale(1.05);
  }
}

/* 功能图标 */
.feature-icons {
  display: flex;
  justify-content: center;
  gap: 48px;
  margin-top: 32px;
}

.feature-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  color: #fff;
  transition: all 0.3s ease;
  cursor: default;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.3);
}

.feature-item:hover {
  transform: translateY(-4px);
  color: #0891b2;
}

.feature-item span {
  font-size: 14px;
  font-weight: 500;
}

/* 右侧表单区 */
.auth-form-section {
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px;
  position: relative;
  z-index: 1;
}

.auth-container {
  width: 100%;
  max-width: 460px;
}

.auth-content {
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 16px;
  padding: 40px 48px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  position: relative;
}

/* 移除顶部装饰条 */

/* 响应式设计 */
@media (max-width: 768px) {
  .auth-layout {
    grid-template-columns: 1fr;
  }

  .auth-brand-section {
    display: none;
  }

  .auth-form-section {
    padding: 24px 16px;
  }

  .auth-container {
    max-width: 100%;
  }

  .auth-content {
    padding: 28px 24px;
  }
}

/* 深色模式支持 */
html.dark .brand-slogan {
  color: #e5e5e5;
}

html.dark .modern-logo {
  background: rgba(30, 30, 46, 0.4);
  border-color: rgba(255, 255, 255, 0.2);
}

html.dark .modern-logo::before {
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.15), rgba(74, 222, 128, 0.15));
}

html.dark .feature-item {
  color: #e5e5e5;
}

html.dark .feature-item:hover {
  color: #38bdf8;
}

html.dark .auth-content {
  background: rgba(30, 30, 46, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}
</style>
