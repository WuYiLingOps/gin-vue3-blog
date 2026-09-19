<!--
 * @ProjectName: go-vue3-blog
 * @FileName: SelfInfoCard.vue
 * @CreateTime: 2026-09-19
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于页个人信息卡片（设计蓝本：anheyu MapAndInfoCard 的 selfInfo 部分），
 *               三栏"小字标签 + 彩色大字"结构（生于 / 就读于 / 现在职业）。
 -->
<template>
  <AboutCard center>
    <div class="self-info">
      <div v-for="(item, index) in items" :key="index" class="info-item">
        <span class="info-tips">{{ item.tips }}</span>
        <strong class="info-value" :style="{ color: colors[index % colors.length], ...valueStyle(item.value) }">
          {{ item.value }}
        </strong>
      </div>
    </div>
  </AboutCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AboutCard from './AboutCard.vue'
import type { AboutSelfInfo } from '@/types/about'

const props = defineProps<{
  config: AboutSelfInfo
}>()

// 三栏配色对齐 anheyu：青 / 橙 / 紫
const colors = ['#43a6c6', '#c69043', '#b04fe6']

// 按内容长度自适应字号；小字号值补偿顶部偏移，使三栏大字基线光学对齐
function valueStyle(value: string): { fontSize: string; marginTop: string } {
  const len = [...value].length
  if (len <= 4) return { fontSize: '24px', marginTop: '0px' }
  if (len <= 8) return { fontSize: '20px', marginTop: '3px' }
  return { fontSize: '16px', marginTop: '6px' }
}

const items = computed(() =>
  [props.config.item1, props.config.item2, props.config.item3].map(item => ({
    ...item
  }))
)
</script>

<style scoped>
.self-info {
  display: grid;
  /* 中栏（校名/专业）占剩余宽列，两侧按内容自适应 */
  grid-template-columns: auto 1fr auto;
  gap: 4px 16px;
  align-items: start;
}

/* 栏宽随内容，三栏标签同行对齐 */
.info-item {
  display: flex;
  flex-direction: column;
  gap: 6px;
  min-width: 0;
}

.info-tips {
  font-size: 14px;
  color: #64748b;
}

html.dark .info-tips {
  color: #94a3b8;
}

.info-value {
  font-weight: 700;
  line-height: 1.25;
  letter-spacing: 0;
  white-space: nowrap;
  overflow: hidden;
}

@media (max-width: 767px) {
  /* 移动端改为纵向三行列表，避免窄屏挤压 */
  .self-info {
    grid-template-columns: 1fr;
    gap: 14px;
  }

  .info-value {
    font-size: 19px !important;
  }
}
</style>
