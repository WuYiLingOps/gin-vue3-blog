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
import App from './App.vue'
import router from './router'
import { VMdEditor, VMdPreview } from './plugins/v-md-editor'

// 样式
import './assets/styles/global.css'

const app = createApp(App)
const pinia = createPinia()

// 使用 Pinia 持久化插件
pinia.use(piniaPluginPersistedstate)

app.use(pinia)
app.use(router)
app.use(VMdEditor)
app.use(VMdPreview)

app.mount('#app')

