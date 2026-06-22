<!--
 * @ProjectName: go-vue3-blog
 * @FileName: App.vue
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 应用根组件，提供全局配置和主题管理
 -->
<template>
  <n-config-provider :theme="theme" :locale="zhCN" :date-locale="dateZhCN" :hljs="hljs">
    <n-loading-bar-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-dialog-provider>
            <!-- 全局加载遮罩 -->
            <transition name="fade">
              <div v-if="appStore.loading" class="global-loading">
                <!-- 动态粒子背景 -->
                <div class="particles">
                  <div
                    class="particle particle-up"
                    v-for="i in 150"
                    :key="'up-' + i"
                    :style="getParticleStyle(i)"
                  ></div>
                  <div
                    class="particle particle-down"
                    v-for="i in 150"
                    :key="'down-' + i"
                    :style="getParticleStyle(i)"
                  ></div>
                </div>

                <div class="loading-content">
                  <h1 class="site-title">{{ appStore.siteName || '加载中' }}</h1>
                  <div class="cube">
                    <div class="cube-face front"></div>
                    <div class="cube-face back"></div>
                    <div class="cube-face left"></div>
                    <div class="cube-face right"></div>
                    <div class="cube-face top"></div>
                    <div class="cube-face bottom"></div>
                  </div>
                  <p class="loading-text">加载中...</p>
                  <div class="progress-bar">
                    <div class="progress-fill"></div>
                  </div>
                </div>
              </div>
            </transition>
            <router-view />
          </n-dialog-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed, watch, onMounted } from 'vue'
import { darkTheme, zhCN, dateZhCN } from 'naive-ui'
import hljs from 'highlight.js'
import { useAppStore } from '@/stores'

const appStore = useAppStore()
const theme = computed(() => (appStore.theme === 'dark' ? darkTheme : null))

// 监听主题变化，切换 html class
watch(
  () => appStore.theme,
  newTheme => {
    if (newTheme === 'dark') {
      document.documentElement.classList.add('dark')
    } else {
      document.documentElement.classList.remove('dark')
    }
  },
  { immediate: true }
)

// 初始化时应用主题
onMounted(() => {
  if (appStore.theme === 'dark') {
    document.documentElement.classList.add('dark')
  }
})

// 生成粒子样式
function getParticleStyle(_index: number) {
  const size = Math.random() * 3 + 1
  const left = Math.random() * 100
  const animationDuration = Math.random() * 10 + 15
  const animationDelay = Math.random() * 5

  return {
    width: `${size}px`,
    height: `${size}px`,
    left: `${left}%`,
    animationDuration: `${animationDuration}s`,
    animationDelay: `${animationDelay}s`
  }
}
</script>

<style scoped>
/* 全局加载遮罩 */
.global-loading {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(
    135deg,
    #0a0e27 0%,
    #1a1f3a 25%,
    #0f1729 50%,
    #1e2139 75%,
    #0a0e27 100%
  );
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 9999;
  overflow: hidden;
}

/* 动态粒子背景 */
.particles {
  position: absolute;
  width: 100%;
  height: 100%;
  overflow: hidden;
}

.particle {
  position: absolute;
  background: rgba(255, 255, 255, 0.9);
  border-radius: 50%;
  box-shadow: 0 0 6px rgba(255, 255, 255, 1);
  animation: float linear infinite;
}

.particle-up {
  bottom: -10px;
}

.particle-down {
  top: -10px;
  animation: floatDown linear infinite;
}

@keyframes float {
  0% {
    transform: translateY(0) translateX(0);
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  90% {
    opacity: 1;
  }
  100% {
    transform: translateY(-100vh) translateX(20px);
    opacity: 0;
  }
}

@keyframes floatDown {
  0% {
    transform: translateY(0) translateX(0);
    opacity: 0;
  }
  10% {
    opacity: 1;
  }
  90% {
    opacity: 1;
  }
  100% {
    transform: translateY(100vh) translateX(-20px);
    opacity: 0;
  }
}

.loading-content {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  z-index: 1;
}

/* 网站标题 */
.site-title {
  font-size: 48px;
  font-weight: 700;
  color: rgba(255, 255, 255, 1);
  margin: 0 0 60px 0;
  letter-spacing: 4px;
  text-shadow:
    0 0 20px rgba(255, 255, 255, 0.5),
    0 0 40px rgba(255, 255, 255, 0.3);
  animation: titleGlow 2s ease-in-out infinite;
}

@keyframes titleGlow {
  0%,
  100% {
    text-shadow:
      0 0 20px rgba(255, 255, 255, 0.5),
      0 0 40px rgba(255, 255, 255, 0.3);
  }
  50% {
    text-shadow:
      0 0 30px rgba(255, 255, 255, 0.8),
      0 0 60px rgba(255, 255, 255, 0.5);
  }
}

/* 3D 立方体容器 */
.cube {
  position: relative;
  width: 80px;
  height: 80px;
  margin: 0 auto;
  transform-style: preserve-3d;
  animation: rotate3d 4s linear infinite;
}

/* 立方体六个面 */
.cube-face {
  position: absolute;
  width: 80px;
  height: 80px;
  background: rgba(255, 255, 255, 0.05);
  backdrop-filter: blur(10px);
  border: 1.5px solid rgba(255, 255, 255, 0.6);
  box-shadow:
    0 0 30px rgba(255, 255, 255, 0.4),
    inset 0 0 30px rgba(255, 255, 255, 0.15);
  border-radius: 2px;
}

.cube-face.front {
  transform: translateZ(40px);
}

.cube-face.back {
  transform: rotateY(180deg) translateZ(40px);
}

.cube-face.left {
  transform: rotateY(-90deg) translateZ(40px);
}

.cube-face.right {
  transform: rotateY(90deg) translateZ(40px);
}

.cube-face.top {
  transform: rotateX(90deg) translateZ(40px);
}

.cube-face.bottom {
  transform: rotateX(-90deg) translateZ(40px);
}

/* 3D 旋转动画 */
@keyframes rotate3d {
  0% {
    transform: rotateX(0deg) rotateY(0deg);
  }
  100% {
    transform: rotateX(360deg) rotateY(360deg);
  }
}

.loading-text {
  color: rgba(255, 255, 255, 0.9);
  font-size: 16px;
  font-weight: 400;
  margin: 40px 0 20px 0;
  letter-spacing: 3px;
  animation: fadeInOut 2s ease-in-out infinite;
}

@keyframes fadeInOut {
  0%,
  100% {
    opacity: 0.6;
  }
  50% {
    opacity: 1;
  }
}

/* 进度条 */
.progress-bar {
  width: 200px;
  height: 3px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 2px;
  overflow: hidden;
  position: relative;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(
    90deg,
    rgba(255, 255, 255, 0.8) 0%,
    rgba(255, 255, 255, 1) 50%,
    rgba(255, 255, 255, 0.8) 100%
  );
  border-radius: 2px;
  animation: progress 2s ease-in-out infinite;
  box-shadow: 0 0 10px rgba(255, 255, 255, 0.8);
}

@keyframes progress {
  0% {
    width: 0%;
    transform: translateX(0);
  }
  50% {
    width: 70%;
  }
  100% {
    width: 100%;
    transform: translateX(0);
  }
}

/* 淡入淡出过渡 */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.5s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* 移动端响应式适配 */
@media (max-width: 768px) {
  /* 隐藏网站标题 */
  .site-title {
    display: none;
  }

  /* 调整立方体大小 */
  .cube {
    width: 60px;
    height: 60px;
  }

  .cube-face {
    width: 60px;
    height: 60px;
  }

  .cube-face.front {
    transform: translateZ(30px);
  }

  .cube-face.back {
    transform: rotateY(180deg) translateZ(30px);
  }

  .cube-face.left {
    transform: rotateY(-90deg) translateZ(30px);
  }

  .cube-face.right {
    transform: rotateY(90deg) translateZ(30px);
  }

  .cube-face.top {
    transform: rotateX(90deg) translateZ(30px);
  }

  .cube-face.bottom {
    transform: rotateX(-90deg) translateZ(30px);
  }

  /* 调整加载文字 */
  .loading-text {
    font-size: 14px;
    margin: 30px 0 15px 0;
  }

  /* 调整进度条 */
  .progress-bar {
    width: 150px;
  }
}
</style>
