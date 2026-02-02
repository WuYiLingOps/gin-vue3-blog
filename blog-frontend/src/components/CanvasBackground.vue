<!--
  项目名称：blog-frontend
  文件名称：CanvasBackground.vue
  创建时间：2026-02-01 20:03:19

  系统用户：Administrator
  作　　者：無以菱
  联系邮箱：huangjing510@126.com
  功能描述：Canvas 背景动画组件，使用 canvas-nest.js 库实现动态线条背景效果，组件卸载时自动清理资源。
-->
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

