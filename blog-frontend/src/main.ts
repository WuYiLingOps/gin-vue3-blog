/*
 * @ProjectName: go-vue3-blog
 * @FileName: main.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 应用入口文件，初始化Vue应用、Pinia状态管理和路由
 */

import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import hljs from 'highlight.js'
import App from './App.vue'
import router from './router'
// v-md-editor 移除全局注册，改为按需导入
import { getPublicSettings } from './api/setting'

// 样式
import './assets/styles/global.css'

const app = createApp(App)
const pinia = createPinia()

// 使用 Pinia 持久化插件
pinia.use(piniaPluginPersistedstate)

app.use(pinia)
app.use(router)

// 配置 Naive UI 的 hljs（用于 NCode 组件代码高亮）
app.provide('hljs', hljs)
// v-md-editor 已移除全局注册，在需要的组件中按需导入

// 先挂载应用（此时 loading 默认为 true，会显示加载动画）
app.mount('#app')

// 预加载关键配置
async function initApp() {
  const { useAppStore, useBlogStore } = await import('./stores')
  const appStore = useAppStore()
  const blogStore = useBlogStore()

  try {
    // 预加载网站配置
    const res = await getPublicSettings()
    if (res.data) {
      // 解析并存储背景图数组
      let bgImages: string[] = []
      if (res.data.cover_bg_images) {
        try {
          bgImages = JSON.parse(res.data.cover_bg_images)
        } catch (e) {
          console.error('解析背景图数组失败:', e)
        }
      }
      appStore.setBgImages(bgImages)

      // 存储网站名称
      if (res.data.site_name) {
        appStore.setSiteName(res.data.site_name)
      }
    }

    // 预加载博客基础数据（分类、标签）
    await blogStore.init()
  } catch (error) {
    console.error('预加载配置失败:', error)
  } finally {
    // 最少显示 1 秒加载动画
    setTimeout(() => {
      appStore.setLoading(false)
    }, 1000)
  }
}

initApp()
