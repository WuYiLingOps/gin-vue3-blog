<!--
 * @ProjectName: go-vue3-blog
 * @FileName: PersonalityCard.vue
 * @CreateTime: 2026-09-19
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于页 MBTI 性格卡片（设计蓝本：anheyu PersonalityCard），
 *               小字引导语 + 彩色人格类型大字 + 底部"了解更多"小字链接。
 -->
<template>
  <AboutCard :tips="config.tips" center>
    <div class="personality-body" :class="{ 'has-img': showImg && !imgFailed }">
      <p class="personality-type" :style="{ color: config.color || '#ac899c' }">
        {{ config.type }}
      </p>
      <p class="personality-posttips">
        在
        <a :href="config.link || 'https://www.16personalities.com/'" target="_blank" rel="noopener nofollow">
          16personalities
        </a>
        了解更多
      </p>
      <div v-if="showImg && !imgFailed" class="personality-img">
        <img :src="config.img" alt="人格形象" loading="lazy" @error="imgFailed = true" />
      </div>
    </div>
  </AboutCard>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import AboutCard from './AboutCard.vue'
import type { AboutPersonality } from '@/types/about'

const props = defineProps<{
  config: AboutPersonality
}>()

const imgFailed = ref(false)
const showImg = computed(() => !!props.config.img?.trim())
</script>

<style scoped>
.personality-body {
  padding: 2px 0;
}

/* 有人格形象图时正文右侧留出空间 */
.personality-body.has-img {
  padding-right: 96px;
}

.personality-img {
  position: absolute;
  right: 18px;
  top: 0;
  bottom: 0;
  z-index: 1;
  display: flex;
  align-items: center;
  pointer-events: none;
  transition: transform 2s cubic-bezier(0.13, 0.45, 0.21, 1.02);
}

.about-card:hover .personality-img {
  transform: rotate(-10deg);
}

.personality-img img {
  max-height: 100%;
  max-width: 100%;
  display: block;
  margin: 0 auto;
}

.personality-type {
  margin: 0;
  font-size: 26px;
  font-weight: 700;
  line-height: 1.25;
  white-space: nowrap;
  overflow: hidden;
}

.personality-posttips {
  margin: 6px 0 0;
  font-size: 12px;
  color: #94a3b8;
}

html.dark .personality-posttips {
  color: #94a3b8;
}

.personality-posttips a {
  color: inherit;
  text-decoration: none;
  border-bottom: 1px dashed currentColor;
}

.personality-posttips a:hover {
  color: #0891b2;
  border-bottom-color: transparent;
}

html.dark .personality-posttips a:hover {
  color: #38bdf8;
}
</style>
