<!--
 * @ProjectName: go-vue3-blog
 * @FileName: AuthorIntroCard.vue
 * @CreateTime: 2026-09-19
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于页作者介绍卡片（设计蓝本：anheyu AuthorPageContent），
 *               站点主题色渐变底 + 问候语 + "我叫 XX" + "是一名 XX、XX" +
 *               打字词轮播（浅色变体）。
 -->
<template>
  <section class="intro-card">
    <p class="intro-tips">{{ config.tips }}</p>
    <h2 class="intro-name">
      {{ config.namePrefix }}
      <strong>{{ authorName }}</strong>
    </h2>
    <p v-if="rolesText" class="intro-roles">
      {{ config.descPrefix }}
      <span class="roles-text">{{ rolesText }}</span>
    </p>
    <div v-if="siteTips.word.length" class="intro-carousel">
      <SiteTipsCarousel
        :title1="siteTips.title1"
        :title2="siteTips.title2"
        :words="siteTips.word"
        light
      />
    </div>
  </section>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import SiteTipsCarousel from './SiteTipsCarousel.vue'
import type { AboutIntro, AboutSiteTips } from '@/types/about'

const props = defineProps<{
  config: AboutIntro
  authorName: string
  /** 打字词轮播配置（about_site_tips），在渐变卡内以浅色变体展示 */
  siteTips: AboutSiteTips
}>()

const rolesText = computed(() =>
  props.config.roles.filter(r => r.trim()).join('、')
)
</script>

<style scoped>
.intro-card {
  position: relative;
  width: 100%;
  padding: 2rem 2.5rem;
  overflow: hidden;
  color: #fff;
  background: linear-gradient(120deg, #0891b2 0%, #2563eb 55%, #7c3aed 100%);
  border-radius: 16px;
  box-shadow: 0 8px 32px rgba(8, 145, 178, 0.25);
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  animation: intro-slide-in 0.6s backwards;
}

.intro-card::after {
  content: '';
  position: absolute;
  right: -40px;
  top: -60px;
  width: 220px;
  height: 220px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.08);
  pointer-events: none;
}

.intro-card::before {
  content: '';
  position: absolute;
  right: 90px;
  bottom: -80px;
  width: 160px;
  height: 160px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
  pointer-events: none;
}

.intro-card:hover {
  box-shadow: 0 16px 48px rgba(8, 145, 178, 0.35);
}

.intro-tips {
  position: relative;
  z-index: 1;
  margin: 0 0 8px;
  font-size: 15px;
  color: rgba(255, 255, 255, 0.85);
}

.intro-name {
  position: relative;
  z-index: 1;
  margin: 0 0 10px;
  font-size: 34px;
  font-weight: 700;
  line-height: 1.25;
  color: #fff;
}

.intro-name strong {
  font-weight: 800;
}

.intro-roles {
  position: relative;
  z-index: 1;
  margin: 0;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.9);
}

.roles-text {
  font-weight: 600;
}

.intro-carousel {
  position: relative;
  z-index: 1;
  margin-top: 18px;
}

@keyframes intro-slide-in {
  from {
    opacity: 0;
    transform: translateY(24px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

@media (max-width: 767px) {
  .intro-card {
    padding: 1.4rem 1.25rem;
  }

  .intro-name {
    font-size: 26px;
  }

  .intro-roles {
    font-size: 14px;
  }
}
</style>
