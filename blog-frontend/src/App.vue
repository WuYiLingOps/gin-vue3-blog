<!--
 * @ProjectName: go-vue3-blog
 * @FileName: App.vue
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 应用根组件，提供全局配置和主题管理
 -->
<template>
  <n-config-provider :theme="theme" :locale="zhCN" :date-locale="dateZhCN">
    <n-loading-bar-provider>
      <n-message-provider>
        <n-notification-provider>
          <n-dialog-provider>
            <router-view />
          </n-dialog-provider>
        </n-notification-provider>
      </n-message-provider>
    </n-loading-bar-provider>
  </n-config-provider>
</template>

<script setup lang="ts">
import { computed, watch, onMounted } from 'vue'
import { darkTheme, zhCN, dateZhCN } from 'naive-ui'
import { useAppStore } from '@/stores'

const appStore = useAppStore()
const theme = computed(() => (appStore.theme === 'dark' ? darkTheme : null))

// 监听主题变化，切换 html class
watch(() => appStore.theme, (newTheme) => {
  if (newTheme === 'dark') {
    document.documentElement.classList.add('dark')
  } else {
    document.documentElement.classList.remove('dark')
  }
}, { immediate: true })

// 初始化时应用主题
onMounted(() => {
  if (appStore.theme === 'dark') {
    document.documentElement.classList.add('dark')
  }
})
</script>

