<!--
  项目名称：blog-frontend
  文件名称：GlobalBackground.vue
  创建时间：2026-04-14 15:40:07

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：全局固定背景图组件，从 app store 读取后台配置的背景图 URL 数组，
           每次刷新时随机选择一张壁纸渲染为视口固定层，叠加亮色/暗色半透明遮罩，支持图片预加载淡入。
           仅当后台设置了背景图时渲染，未设置时不渲染，由 body CSS 渐变兜底。
-->
<template>
  <div class="global-bg">
    <!-- 背景图片层：预加载完成后切换，优先后台配置，未配置时使用默认背景图 -->
    <img :src="bgSrc" alt="" class="global-bg-img" :class="{ loaded: imgLoaded }" />
    <!-- 半透明遮罩层：保证内容可读 -->
    <div class="global-bg-overlay"></div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch } from 'vue'
import { useAppStore } from '@/stores'

const DEFAULT_BG = '/default_background.webp'

const appStore = useAppStore()
const currentSrc = ref('') // 已确认加载完成的背景图 URL
const imgLoaded = ref(false)

// 展示源：优先后台选定的壁纸，未就绪时回落到本地默认图
const bgSrc = computed(() => appStore.pickedBgUrl || DEFAULT_BG)

// 预加载完成后再切换显示源，避免切换期间出现灰色窗口期
function preloadAndApply(url: string) {
  const img = new Image()
  img.onload = () => {
    currentSrc.value = url
    imgLoaded.value = true
  }
  img.onerror = () => {
    // 后台配置的图加载失败：清空列表回退默认图（保持原有行为）
    if (appStore.bgImages.length > 0) {
      appStore.setBgImages([])
    } else if (url !== DEFAULT_BG) {
      preloadAndApply(DEFAULT_BG)
    }
    // 默认图也失败时保持灰色兜底（body 渐变），不再重试避免死循环
  }
  img.src = url
}

// 展示源变化时触发预加载
watch(
  bgSrc,
  url => {
    if (url !== currentSrc.value) preloadAndApply(url)
  },
  { immediate: true }
)
</script>

<style scoped>
.global-bg {
  position: fixed;
  inset: 0;
  z-index: 0;
  pointer-events: none;
}

.global-bg-img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  display: block;
  opacity: 0;
  transition: opacity 0.8s ease;
}

.global-bg-img.loaded {
  opacity: 1;
}

/* 亮色模式遮罩：轻微提亮，保证玻璃态卡片文字可读 */
.global-bg-overlay {
  position: absolute;
  inset: 0;
  background: rgba(255, 255, 255, 0.15);
  transition: background 0.3s ease;
}

/* 暗色模式遮罩：压暗背景，适配深色主题 */
html.dark .global-bg-overlay {
  background: rgba(0, 0, 0, 0.45);
}
</style>
