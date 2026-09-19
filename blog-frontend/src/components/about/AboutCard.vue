<!--
 * @ProjectName: go-vue3-blog
 * @FileName: AboutCard.vue
 * @CreateTime: 2026-09-18
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于页通用卡片基座，提供玻璃态卡片 + itemTips小字引导 + itemTitle大标题
 *               三件套结构（设计蓝本：anheyu 关于本站的卡片范式）
 -->
<template>
  <section class="about-card" :class="{ flush, center }">
    <p v-if="tips" class="item-tips">{{ tips }}</p>
    <h3 v-if="title" class="item-title">{{ title }}</h3>
    <div class="card-body">
      <slot />
    </div>
  </section>
</template>

<script setup lang="ts">
withDefaults(
  defineProps<{
    tips?: string
    title?: string
    /** 无内边距模式，内容贴卡片边缘（用于图片/地图类卡片） */
    flush?: boolean
    /** 垂直居中模式（用于成对行内较矮的内容卡） */
    center?: boolean
  }>(),
  {
    tips: '',
    title: '',
    flush: false,
    center: false
  }
)
</script>

<style scoped>
.about-card {
  position: relative;
  width: 100%;
  padding: 1.25rem 2rem;
  background: rgba(255, 255, 255, 0.9);
  backdrop-filter: blur(20px) saturate(180%);
  -webkit-backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
  border-radius: 16px;
  transition: all 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  animation: card-slide-in 0.6s backwards;
}

.about-card:hover {
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.12);
  border-color: rgba(8, 145, 178, 0.3);
}

html.dark .about-card {
  background: rgba(30, 41, 59, 0.7);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

html.dark .about-card:hover {
  box-shadow: 0 16px 48px rgba(0, 0, 0, 0.5);
  border-color: rgba(56, 189, 248, 0.3);
}

.item-tips {
  margin: 0 0 0.25rem;
  font-size: 0.8rem;
  color: #64748b;
  opacity: 0.85;
}

html.dark .item-tips {
  color: #94a3b8;
}

.item-title {
  margin: 0 0 1rem;
  font-size: 34px;
  font-weight: 700;
  line-height: 1.2;
  color: #1a202c;
}

html.dark .item-title {
  color: #e5e5e5;
}

.card-body {
  width: 100%;
}

/* 垂直居中模式：成对行内较矮内容卡在拉伸高度下居中 */
.about-card.center {
  display: flex;
  flex-direction: column;
  justify-content: center;
}

.about-card.center .card-body {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

/* 贴边模式：去内边距，圆角内裁剪 */
.about-card.flush {
  padding: 0;
  overflow: hidden;
}

.about-card.flush .card-body {
  height: 100%;
}

@keyframes card-slide-in {
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
  .about-card {
    padding: 1rem;
  }

  .item-title {
    font-size: 26px;
  }
}
</style>
