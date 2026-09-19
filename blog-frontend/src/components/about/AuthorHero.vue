<!--
 * @ProjectName: go-vue3-blog
 * @FileName: AuthorHero.vue
 * @CreateTime: 2026-09-18
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于页头像框组件，居中大头像 + 左右交错漂浮的技能标签
 -->
<template>
  <div class="author-hero">
    <!-- 左侧浮动标签 -->
    <div v-if="leftTags.length" class="tag-column">
      <span
        v-for="(tag, i) in leftTags"
        :key="`l-${tag}`"
        class="float-tag"
        :style="{ animationDelay: `${i * 0.6}s` }"
      >
        {{ tag }}
      </span>
    </div>

    <!-- 头像 -->
    <div class="avatar-section">
      <n-avatar
        :src="avatar || ''"
        :size="avatarSize"
        round
        :fallback-src="defaultAvatar"
        class="hero-avatar"
      >
        <template v-if="!avatar">{{ fallbackText }}</template>
      </n-avatar>
    </div>

    <!-- 右侧浮动标签 -->
    <div v-if="rightTags.length" class="tag-column">
      <span
        v-for="(tag, i) in rightTags"
        :key="`r-${tag}`"
        class="float-tag"
        :style="{ animationDelay: `${i * 0.6 + 0.3}s` }"
      >
        {{ tag }}
      </span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'

const props = withDefaults(
  defineProps<{
    avatar?: string
    fallbackText?: string
    leftTags?: string[]
    rightTags?: string[]
  }>(),
  {
    avatar: '',
    fallbackText: '博主',
    leftTags: () => [],
    rightTags: () => []
  }
)

const defaultAvatar = '/default-avatar.png'
const avatarSize = computed(() => 120)
</script>

<style scoped>
.author-hero {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 24px;
  margin: 0 0 8px;
  animation: hero-slide-in 0.6s 0.1s backwards;
}

.tag-column {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 14px;
}

.tag-column:last-child {
  align-items: flex-start;
}

.float-tag {
  padding: 5px 14px;
  font-size: 13px;
  color: #475569;
  white-space: nowrap;
  background: rgba(255, 255, 255, 0.75);
  backdrop-filter: blur(10px);
  -webkit-backdrop-filter: blur(10px);
  border: 1px solid rgba(8, 145, 178, 0.18);
  border-radius: 999px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.06);
  animation: tag-float 6s ease-in-out infinite;
  transition: all 0.3s ease;
}

html.dark .float-tag {
  color: #cbd5e1;
  background: rgba(30, 41, 59, 0.75);
  border-color: rgba(56, 189, 248, 0.2);
}

.float-tag:hover {
  color: #0891b2;
  border-color: rgba(8, 145, 178, 0.45);
  transform: translateY(-2px);
}

html.dark .float-tag:hover {
  color: #38bdf8;
}

.avatar-section :deep(.n-avatar) {
  box-shadow: 0 8px 24px rgba(8, 145, 178, 0.15);
  transition: all 0.3s;
  border: 2px solid rgba(8, 145, 178, 0.1);
}

.avatar-section :deep(.n-avatar:hover) {
  transform: translateY(-3px) scale(1.05);
  box-shadow: 0 14px 32px rgba(8, 145, 178, 0.25);
  border-color: rgba(8, 145, 178, 0.3);
}

@keyframes tag-float {
  0%,
  100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-9px);
  }
}

@keyframes hero-slide-in {
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
  .author-hero {
    flex-direction: column;
    gap: 16px;
  }

  /* 移动端隐藏左右浮动标签，避免窄屏拥挤 */
  .tag-column {
    display: none;
  }
}
</style>
