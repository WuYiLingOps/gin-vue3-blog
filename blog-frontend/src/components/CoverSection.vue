<!--
  项目名称：blog-frontend
  文件名称：CoverSection.vue
  创建时间：2026-04-13 17:24:45

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：首页全屏封面组件，展示博客标题、打字机副标题、社交图标和下拉箭头，
           背景为后台配置的壁纸图片叠加 Canvas 粒子动画效果
-->
<template>
  <section class="cover-section">
    <!-- 半透明遮罩（确保文字可读，背景图由 GlobalBackground 提供） -->
    <div class="cover-overlay"></div>

    <!-- 层级3：内容区 -->
    <div class="cover-content">
      <h1 class="cover-title">{{ siteName }}</h1>
      <div class="cover-subtitle" v-if="subtitles.length">
        <span class="typewriter-text">{{ displayedText }}</span>
        <span class="typewriter-cursor">|</span>
      </div>
    </div>

    <!-- 层级3：下拉箭头 -->
    <div class="cover-arrow" @click="scrollToContent">
      <svg
        xmlns="http://www.w3.org/2000/svg"
        width="36"
        height="36"
        viewBox="0 0 24 24"
        fill="none"
        stroke="currentColor"
        stroke-width="2"
        stroke-linecap="round"
        stroke-linejoin="round"
      >
        <polyline points="6 9 12 15 18 9"></polyline>
      </svg>
    </div>
  </section>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { getPublicSettings } from '@/api/setting'
import { useTypewriter } from '@/composables/useTypewriter'

const props = defineProps<{
  scrollElement?: HTMLElement | null
}>()

// 站点数据
const siteName = ref('菱风叙')

// 副标题数据（从后台设置读取，按换行分隔）
const subtitleRaw = ref('')
const subtitles = computed(() => {
  if (!subtitleRaw.value) return []
  return subtitleRaw.value
    .split('\n')
    .map(s => s.trim())
    .filter(Boolean)
})

// 打字机效果（传入 computed ref，数据加载后自动启动）
const effectiveSubtitles = computed(() =>
  subtitles.value.length ? subtitles.value : ['欢迎来到我的博客']
)
const { displayedText } = useTypewriter(effectiveSubtitles, {
  typeSpeed: 120,
  deleteSpeed: 60,
  pauseTime: 2000
})

// 下拉到内容区
function scrollToContent() {
  const el = props.scrollElement
  if (el) {
    el.scrollTo({ top: window.innerHeight, behavior: 'smooth' })
  }
}

// 获取站点设置
async function fetchSettings() {
  try {
    const res = await getPublicSettings()
    if (!res.data) return

    siteName.value = res.data.site_name || '菱风叙'
    subtitleRaw.value = res.data.cover_subtitle || ''
    // 背景图已由 DefaultLayout → appStore.bgImage 全局管理，此处不再加载
  } catch (err) {
    console.error('CoverSection: 获取设置失败', err)
  }
}

onMounted(() => {
  fetchSettings()
})
</script>

<style scoped>
/* 封面根容器：全屏高度，父容器无 padding，直接 100% 宽铺满 */
.cover-section {
  position: relative;
  height: 100vh;
  width: 100%;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  z-index: 0;
}

/* 半透明遮罩 */
.cover-overlay {
  position: absolute;
  inset: 0;
  z-index: 2;
  background: linear-gradient(
    to bottom,
    rgba(0, 0, 0, 0.15) 0%,
    rgba(0, 0, 0, 0.35) 50%,
    rgba(0, 0, 0, 0.45) 100%
  );
}
html.dark .cover-overlay {
  background: linear-gradient(
    to bottom,
    rgba(0, 0, 0, 0.35) 0%,
    rgba(0, 0, 0, 0.5) 50%,
    rgba(0, 0, 0, 0.6) 100%
  );
}

/* 内容区 */
.cover-content {
  position: relative;
  z-index: 3;
  text-align: center;
  color: #fff;
  padding: 0 24px;
}

/* 标题 */
.cover-title {
  font-size: 56px;
  font-weight: 800;
  letter-spacing: 0.02em;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  margin: 0 0 20px 0;
  animation: coverFadeInUp 1s ease-out;
}

/* 副标题 + 打字机光标 */
.cover-subtitle {
  font-size: 22px;
  font-weight: 400;
  text-shadow: 0 2px 10px rgba(0, 0, 0, 0.3);
  min-height: 34px;
  margin-bottom: 32px;
  animation: coverFadeInUp 1s ease-out 0.3s both;
}
.typewriter-cursor {
  display: inline-block;
  animation: cursorBlink 1s step-end infinite;
  margin-left: 2px;
  font-weight: 300;
}

/* 下拉箭头 */
.cover-arrow {
  position: absolute;
  bottom: 40px;
  left: 50%;
  transform: translateX(-50%);
  z-index: 3;
  color: #fff;
  cursor: pointer;
  animation: coverBounce 2s ease-in-out infinite;
  opacity: 0.8;
  transition: opacity 0.3s;
}
.cover-arrow:hover {
  opacity: 1;
}

/* 动画定义 */
@keyframes coverFadeInUp {
  from {
    opacity: 0;
    transform: translateY(30px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}
@keyframes cursorBlink {
  0%,
  100% {
    opacity: 1;
  }
  50% {
    opacity: 0;
  }
}
@keyframes coverBounce {
  0%,
  20%,
  50%,
  80%,
  100% {
    transform: translateX(-50%) translateY(0);
  }
  40% {
    transform: translateX(-50%) translateY(-12px);
  }
  60% {
    transform: translateX(-50%) translateY(-6px);
  }
}

/* 响应式：平板 */
@media (max-width: 768px) {
  .cover-title {
    font-size: 32px;
  }
  .cover-subtitle {
    font-size: 16px;
    min-height: 26px;
  }
  .cover-social-icon {
    width: 36px;
    height: 36px;
  }
  .cover-arrow {
    bottom: 24px;
  }
}

/* 响应式：手机 */
@media (max-width: 480px) {
  .cover-title {
    font-size: 26px;
  }
  .cover-subtitle {
    font-size: 14px;
    min-height: 22px;
  }
  .cover-social {
    gap: 12px;
  }
}
</style>
