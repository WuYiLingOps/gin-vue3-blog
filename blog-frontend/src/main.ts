import { createApp } from 'vue'
import { createPinia } from 'pinia'
import piniaPluginPersistedstate from 'pinia-plugin-persistedstate'
import App from './App.vue'
import router from './router'

// 样式
import './assets/styles/global.css'

const app = createApp(App)
const pinia = createPinia()

// 使用 Pinia 持久化插件
pinia.use(piniaPluginPersistedstate)

app.use(pinia)
app.use(router)

app.mount('#app')

