<template>
  <div class="canvas-bg" ref="container"></div>
</template>

<script setup lang="ts">
import { onBeforeUnmount, onMounted, ref } from 'vue'
import CanvasNest from 'canvas-nest.js'

const container = ref<HTMLElement | null>(null)
let instance: CanvasNest | null = null

onMounted(() => {
  const el = container.value || document.body
  instance = new CanvasNest(el, {
    color: '24,170,204', // 主色调 RGB
    opacity: 0.8,        // 线条透明度
    count: 120,          // 线条数量
    zIndex: -1           // 置于内容下方
  })
})

onBeforeUnmount(() => {
  instance?.destroy()
  instance = null
})
</script>

<style scoped>
.canvas-bg {
  position: fixed;
  inset: 0;
  pointer-events: none;
  z-index: -1;
}
</style>

