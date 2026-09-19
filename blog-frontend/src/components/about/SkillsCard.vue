<!--
 * @ProjectName: go-vue3-blog
 * @FileName: SkillsCard.vue
 * @CreateTime: 2026-09-18
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于页技能卡片（设计蓝本：anheyu SkillsCard 的图标跑马灯），
 *               品牌色圆角方块 + 白色单色图标，双行反向无缝滚动，悬浮暂停。
 -->
<template>
  <AboutCard :tips="config.tips" :title="config.title">
    <div class="skills-marquee">
      <div
        v-for="(row, rowIndex) in rows"
        :key="rowIndex"
        class="marquee"
        :class="{ reverse: rowIndex % 2 === 1 }"
      >
        <div
          class="track"
          :style="{ animationDuration: `${18 + rowIndex * 5}s` }"
        >
          <div
            v-for="(item, i) in loopItems(row)"
            :key="`${rowIndex}-${i}`"
            class="skill-tile"
            :style="{ background: tileColor(item) }"
            :title="item.name"
          >
            <svg v-if="tileIcon(item)" viewBox="0 0 24 24" aria-hidden="true">
              <path :d="tileIcon(item)!.path" fill="#fff" />
            </svg>
            <span v-else class="tile-letter">{{ item.name.charAt(0).toUpperCase() }}</span>
          </div>
        </div>
      </div>
    </div>
  </AboutCard>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import AboutCard from './AboutCard.vue'
import { getTechIcon } from '@/utils/tech-icons'
import type { AboutSkills, AboutSkillItem } from '@/types/about'

const props = defineProps<{
  config: AboutSkills
}>()

// 技能列表拆成上下两行
const rows = computed(() => {
  const list = props.config.list
  if (list.length <= 4) return [list]
  const half = Math.ceil(list.length / 2)
  return [list.slice(0, half), list.slice(half)]
})

// 无缝滚动：每行内容复制一份首尾相接
function loopItems(row: AboutSkillItem[]): AboutSkillItem[] {
  if (!row.length) return []
  return [...row, ...row]
}

function tileIcon(item: AboutSkillItem) {
  return getTechIcon(item.icon)
}

function tileColor(item: AboutSkillItem): string {
  return item.color || '#' + (getTechIcon(item.icon)?.hex || '0891b2')
}
</script>

<style scoped>
.skills-marquee {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.marquee {
  overflow: hidden;
  /* 两端淡出，弱化裁切感 */
  mask-image: linear-gradient(90deg, transparent, #000 6%, #000 94%, transparent);
  -webkit-mask-image: linear-gradient(90deg, transparent, #000 6%, #000 94%, transparent);
}

.track {
  display: flex;
  width: max-content;
  /* 覆盖全局 * { max-width: 100% } 规则，保证轨道完整铺开实现无缝滚动 */
  max-width: none;
  animation: marquee-left 18s linear infinite;
}

.marquee.reverse .track {
  animation-name: marquee-right;
}

/* 悬浮暂停，方便查看图标 */
.marquee:hover .track {
  animation-play-state: paused;
}

.skill-tile {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 64px;
  height: 64px;
  margin-right: 16px;
  border-radius: 16px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
  transition: transform 0.3s ease;
}

.skill-tile:hover {
  transform: translateY(-4px) scale(1.05);
}

.skill-tile svg {
  width: 30px;
  height: 30px;
}

.tile-letter {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
}

@keyframes marquee-left {
  from {
    transform: translateX(0);
  }
  to {
    transform: translateX(-50%);
  }
}

@keyframes marquee-right {
  from {
    transform: translateX(-50%);
  }
  to {
    transform: translateX(0);
  }
}

@media (prefers-reduced-motion: reduce) {
  .track {
    animation: none;
  }
}

@media (max-width: 767px) {
  .skill-tile {
    width: 52px;
    height: 52px;
    margin-right: 12px;
    border-radius: 13px;
  }

  .skill-tile svg {
    width: 25px;
    height: 25px;
  }
}
</style>
