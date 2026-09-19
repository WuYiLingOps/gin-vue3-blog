<!--
 * @ProjectName: go-vue3-blog
 * @FileName: MaximCard.vue
 * @CreateTime: 2026-09-19
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于页"追求"卡片（设计蓝本：anheyu 线上版 MaximCard），
 *               小字引导语 + 左对齐多行大字（正文深色），末行内嵌轮换词：
 *               每 2s 当前词上滑退出、下一词自下滑入（0.5s ease-in-out），
 *               词色按位置循环四种 45° 渐变（青/绿/紫/红）。
 -->
<template>
  <AboutCard :tips="config.tips" center>
    <div class="maxim-body">
      <p v-for="(line, i) in topLines" :key="`t-${i}`" class="maxim-line">{{ line }}</p>
      <p class="maxim-line maxim-word-line">
        <template v-if="words.length">
          <span class="word-mask">
            <Transition name="word-slide">
              <span :key="currentIndex" class="word" :class="`word-c${currentIndex % 4}`">
                {{ words[currentIndex] }}
              </span>
            </Transition>
          </span>
        </template>
        <span v-else-if="config.bottom" class="word word-static">{{ config.bottom }}</span>
      </p>
    </div>
  </AboutCard>
</template>

<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import AboutCard from './AboutCard.vue'
import type { AboutQuote } from '@/types/about'

const props = defineProps<{
  config: AboutQuote
}>()

const currentIndex = ref(0)
const words = computed(() => (props.config.word || []).filter(w => w.trim()))
let timer: ReturnType<typeof setInterval> | null = null

// 正文按换行拆分为多行，轮换词独占末行
const topLines = computed(() => (props.config.top || '').split('\n'))

function startRotation() {
  stopRotation()
  if (words.value.length <= 1) return
  timer = setInterval(() => {
    currentIndex.value = (currentIndex.value + 1) % words.value.length
  }, 2000)
}

function stopRotation() {
  if (timer) {
    clearInterval(timer)
    timer = null
  }
}

watch(
  () => props.config.word,
  () => {
    currentIndex.value = 0
    startRotation()
  },
  { deep: true }
)

onMounted(startRotation)
onUnmounted(stopRotation)
</script>

<style scoped>
.maxim-body {
  padding: 4px 0;
}

.maxim-line {
  margin: 0;
  font-size: 40px;
  font-weight: 700;
  line-height: 1.35;
  letter-spacing: 0.5px;
  color: #1a202c;
  text-align: left;
}

html.dark .maxim-line {
  color: #e5e5e5;
}

/* 轮换词遮罩：独占一行，单行高度，溢出隐藏 */
.word-mask {
  position: relative;
  display: inline-block;
  overflow: hidden;
  min-height: 1.2em;
}

.word {
  display: inline-block;
  background: linear-gradient(45deg, #0ecffe 50%, #07a6f1);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  transition: transform 0.5s ease-in-out, opacity 0.5s ease-in-out;
}

/* 四种渐变词色（对齐 anheyu wordColor1-4），用 background-image 避免重置 background-clip */
.word-c0 {
  background-image: linear-gradient(45deg, #0ecffe 50%, #07a6f1);
}

.word-c1 {
  background-image: linear-gradient(45deg, #18e198 50%, #0ec15d);
}

.word-c2 {
  background-image: linear-gradient(45deg, #8a7cfb 50%, #633e9c);
}

.word-c3 {
  background-image: linear-gradient(45deg, #fa7671 50%, #f45f7f);
}

.word-static {
  background-image: linear-gradient(45deg, #8a7cfb 50%, #633e9c);
}

/* 轮换过渡：当前词上滑退出，下一词自下滑入 */
.word-slide-enter-active,
.word-slide-leave-active {
  transition:
    transform 0.5s ease-in-out,
    opacity 0.5s ease-in-out;
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
  .maxim-line {
    font-size: 28px;
  }
}
</style>
