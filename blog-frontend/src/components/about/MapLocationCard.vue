<!--
 * @ProjectName: go-vue3-blog
 * @FileName: MapLocationCard.vue
 * @CreateTime: 2026-09-19
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于页地理位置卡片（设计蓝本：anheyu MapAndInfoCard 的地图部分），
 *               根据配置的经纬度实时渲染高德瓦片。瓦片网格 3×2、按接近原生尺寸
 *               显示以保证地名文字可读；悬浮时 4s 缓动放大至 120%、位置条下滑
 *               隐藏，移开 1s 恢复。
 -->
<template>
  <AboutCard flush>
    <div class="map-wrap">
      <div v-if="!tilesFailed" class="map-tiles">
        <img
          v-for="tile in tiles"
          :key="`${tile.x}-${tile.y}`"
          :src="tile.url"
          alt=""
          loading="lazy"
          @error="tilesFailed = true"
        />
      </div>
      <div v-else class="map-fallback" />

      <span class="map-attribution">地图数据 © 高德地图</span>

      <div class="map-bar">
        <svg class="map-pin" viewBox="0 0 24 24" fill="none" aria-hidden="true">
          <path
            d="M12 21s-7-5.5-7-11a7 7 0 1 1 14 0c0 5.5-7 11-7 11Z"
            stroke="currentColor"
            stroke-width="1.8"
          />
          <circle cx="12" cy="10" r="2.6" stroke="currentColor" stroke-width="1.8" />
        </svg>
        <span class="map-bar-text">
          {{ config.title }}
          <strong>{{ config.location }}</strong>
        </span>
      </div>
    </div>
  </AboutCard>
</template>

<script setup lang="ts">
import { computed, ref } from 'vue'
import AboutCard from './AboutCard.vue'
import type { AboutMapConfig } from '@/types/about'

const props = defineProps<{
  config: AboutMapConfig
}>()

const tilesFailed = ref(false)

const DEFAULT_TILE_URL =
  'https://webrd0{s}.is.autonavi.com/appmaptile?lang=zh_cn&size=1&scale=1&style=8&x={x}&y={y}&z={z}'

// 经纬度 → 瓦片坐标
function lngToTileX(lng: number, z: number): number {
  return Math.floor(((lng + 180) / 360) * 2 ** z)
}

function latToTileY(lat: number, z: number): number {
  const rad = (lat * Math.PI) / 180
  return Math.floor(((1 - Math.log(Math.tan(rad) + 1 / Math.cos(rad)) / Math.PI) / 2) * 2 ** z)
}

// 渲染 3x2 瓦片网格：每块瓦片以接近原生尺寸显示，保证地名文字清晰可读
const tiles = computed(() => {
  const cfg = props.config
  const z = Math.min(17, Math.max(3, Math.round(cfg.zoom || 10)))
  const cols = 3
  const rows = 2
  const centerX = lngToTileX(cfg.lng, z)
  const centerY = latToTileY(cfg.lat, z)
  const max = 2 ** z

  const list: Array<{ x: number; y: number; url: string }> = []
  for (let r = 0; r < rows; r++) {
    for (let c = 0; c < cols; c++) {
      const x = centerX - 1 + c
      const y = centerY + r
      if (x < 0 || x >= max || y < 0 || y >= max) continue
      const url = DEFAULT_TILE_URL.replace('{s}', String(((x + y) % 4) + 1))
        .replace('{x}', String(x))
        .replace('{y}', String(y))
        .replace('{z}', String(z))
      list.push({ x, y, url })
    }
  }
  return list
})
</script>

<style scoped>
.map-wrap {
  position: relative;
  width: 100%;
  height: 330px;
  overflow: hidden;
}

.map-tiles {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  background: rgba(8, 145, 178, 0.08);
  /* 悬浮缓动放大（复刻 anheyu .map:hover：4s 放大至 120%，离开 1s 恢复） */
  transform: scale(1);
  transform-origin: 0 36%;
  transition: transform 1s ease-in-out;
}

.map-tiles img {
  display: block;
  width: 100%;
  aspect-ratio: 1;
  object-fit: cover;
}

.map-fallback {
  height: 240px;
  background: linear-gradient(135deg, rgba(8, 145, 178, 0.15), rgba(5, 150, 105, 0.15));
}

html.dark .map-tiles img {
  /* 暗色模式反色处理，得到深色地图 */
  filter: invert(0.9) hue-rotate(180deg) contrast(0.92) brightness(0.92);
}

.map-wrap:hover .map-tiles {
  transform: scale(1.2);
  transition: transform 4s ease-in-out;
}

.map-attribution {
  position: absolute;
  right: 6px;
  top: 6px;
  z-index: 2;
  padding: 2px 8px;
  font-size: 11px;
  color: #475569;
  background: rgba(255, 255, 255, 0.7);
  border-radius: 6px;
  pointer-events: none;
}

html.dark .map-attribution {
  color: #cbd5e1;
  background: rgba(15, 23, 42, 0.7);
}

.map-bar {
  position: absolute;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 3;
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 14px 20px;
  background: rgba(255, 255, 255, 0.65);
  backdrop-filter: blur(16px) saturate(160%);
  -webkit-backdrop-filter: blur(16px) saturate(160%);
  border-top: 1px solid rgba(255, 255, 255, 0.4);
  transition: all 1s ease-in-out;
}

/* 悬浮时位置条下滑隐藏（复刻 anheyu .map:hover .mapTitle） */
.map-wrap:hover .map-bar {
  bottom: -100%;
}

html.dark .map-bar {
  background: rgba(30, 41, 59, 0.65);
  border-top-color: rgba(255, 255, 255, 0.08);
}

.map-pin {
  width: 22px;
  height: 22px;
  flex-shrink: 0;
  color: #0891b2;
}

html.dark .map-pin {
  color: #38bdf8;
}

.map-bar-text {
  font-size: 18px;
  color: #1a202c;
}

html.dark .map-bar-text {
  color: #e5e5e5;
}

.map-bar-text strong {
  font-weight: 700;
}

@media (max-width: 767px) {
  .map-bar {
    padding: 11px 14px;
  }

  .map-bar-text {
    font-size: 15px;
  }
}
</style>
