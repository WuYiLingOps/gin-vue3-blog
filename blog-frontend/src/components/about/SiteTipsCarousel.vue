<!--
 * @ProjectName: go-vue3-blog
 * @FileName: SiteTipsCarousel.vue
 * @CreateTime: 2026-09-18
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于页打字词轮播组件，大号词语每 2s 上下位移轮播
 *               （设计蓝本：anheyu AboutSiteTips 的 data-show/data-up 轮播）
 -->
<template>
  <div class="site-tips">
    <p v-if="tips" class="tips-line">{{ tips }}</p>
    <div v-if="words.length" class="word-line">
      <span v-if="title1" class="word-fix" :class="{ light }">{{ title1 }}</span>
      <span class="word-carousel">
        <Transition name="word-slide">
          <span :key="currentIndex" class="word-current" :class="{ light }">{{
            words[currentIndex]
          }}</span>
        </Transition>
      </span>
      <span v-if="title2" class="word-fix" :class="{ light }">{{ title2 }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, onUnmounted, ref, watch } from 'vue'

const props = withDefaults(
  defineProps<{
    tips?: string
    title1?: string
    title2?: string
    words?: string[]
    /** light 变体用于渐变/深色背景上，文字改白色 */
    light?: boolean
  }>(),
  {
    tips: '',
    title1: '',
    title2: '',
    words: () => [],
    light: false
  }
)

const currentIndex = ref(0)
let timer: ReturnType<typeof setInterval> | null = null

function startCarousel() {
  stopCarousel()
  if (props.words.length <= 1) return
  timer = setInterval(() => {
    currentIndex.value = (currentIndex.value + 1) % props.words.length
  }, 2000)
}

function stopCarousel() {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

watch(
  () => props.words,
  () => {
    currentIndex.value = 0
    startCarousel()
  },
  { deep: true }
)

onMounted(startCarousel)
onUnmounted(stopCarousel)
</script>

<style scoped>
.site-tips {
  margin-top: 4px;
}

.tips-line {
  margin: 0 0 6px;
  font-size: 14px;
  color: #64748b;
}

.tips-line.light {
  color: rgba(255, 255, 255, 0.85);
}

html.dark .tips-line:not(.light) {
  color: #94a3b8;
}

.word-line {
  display: flex;
  align-items: baseline;
  gap: 10px;
  font-size: 28px;
  font-weight: 700;
  line-height: 1.3;
}

.word-fix {
  color: #1a202c;
}

.word-fix.light {
  color: #fff;
}

html.dark .word-fix:not(.light) {
  color: #e5e5e5;
}

.word-carousel {
  position: relative;
  display: inline-block;
  min-height: 1.3em;
  overflow: hidden;
}

.word-current {
  display: inline-block;
  background: linear-gradient(135deg, #0891b2 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.word-current.light {
  background: none;
  -webkit-text-fill-color: #fff;
  color: #fff;
}

html.dark .word-current:not(.light) {
  background: linear-gradient(135deg, #38bdf8 0%, #4ade80 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* 轮播过渡：旧词上滑淡出，新词自下滑入 */
.word-slide-enter-active,
.word-slide-leave-active {
  transition:
    transform 0.45s cubic-bezier(0.4, 0, 0.2, 1),
    opacity 0.45s ease;
}

.word-slide-enter-from {
  opacity: 0;
  transform: translateY(100%);
}

.word-slide-leave-to {
  opacity: 0;
  transform: translateY(-100%);
}

.word-slide-leave-active {
  position: absolute;
  left: 0;
  top: 0;
}

@media (max-width: 767px) {
  .word-line {
    font-size: 20px;
  }
}
</style>
