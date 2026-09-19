<!--
 * @ProjectName: go-vue3-blog
 * @FileName: CareersCard.vue
 * @CreateTime: 2026-09-18
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于页职业生涯时间线卡片，竖向轴线 + 彩色节点
 -->
<template>
  <AboutCard :tips="config.tips" :title="config.title">
    <div class="timeline">
      <div
        v-for="(item, index) in config.list"
        :key="`${item.time}-${index}`"
        class="timeline-item"
        :style="{ '--node-color': item.color || '#0891b2' }"
      >
        <div class="timeline-rail">
          <i class="timeline-node" />
          <i v-if="index < config.list.length - 1" class="timeline-line" />
        </div>
        <div class="timeline-body">
          <div class="timeline-head">
            <span class="timeline-time">{{ item.time }}</span>
            <strong class="timeline-title">{{ item.title }}</strong>
          </div>
          <p class="timeline-desc">{{ item.desc }}</p>
        </div>
      </div>
    </div>
  </AboutCard>
</template>

<script setup lang="ts">
import AboutCard from './AboutCard.vue'
import type { AboutCareers } from '@/types/about'

defineProps<{
  config: AboutCareers
}>()
</script>

<style scoped>
.timeline {
  display: flex;
  flex-direction: column;
}

.timeline-item {
  display: flex;
  gap: 14px;
}

.timeline-rail {
  display: flex;
  flex-direction: column;
  align-items: center;
  width: 14px;
  flex-shrink: 0;
}

.timeline-node {
  width: 12px;
  height: 12px;
  margin-top: 5px;
  border-radius: 50%;
  background: var(--node-color);
  box-shadow: 0 0 0 4px color-mix(in srgb, var(--node-color) 18%, transparent);
}

.timeline-line {
  flex: 1;
  width: 2px;
  min-height: 18px;
  background: linear-gradient(180deg, var(--node-color), rgba(8, 145, 178, 0.15));
  border-radius: 1px;
  opacity: 0.5;
}

.timeline-body {
  flex: 1;
  min-width: 0;
  padding-bottom: 20px;
}

.timeline-head {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.timeline-time {
  padding: 2px 10px;
  font-size: 12px;
  font-weight: 600;
  color: #fff;
  background: var(--node-color);
  border-radius: 999px;
}

.timeline-title {
  font-size: 16px;
  font-weight: 600;
  color: #1a202c;
}

html.dark .timeline-title {
  color: #e5e5e5;
}

.timeline-desc {
  margin: 6px 0 0;
  font-size: 14px;
  line-height: 1.7;
  color: #64748b;
}

html.dark .timeline-desc {
  color: #94a3b8;
}

.timeline-item:last-child .timeline-body {
  padding-bottom: 0;
}
</style>
