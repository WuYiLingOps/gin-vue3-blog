import { defineConfig, loadEnv } from 'vite'
import vue from '@vitejs/plugin-vue'
import { resolve } from 'path'
import AutoImport from 'unplugin-auto-import/vite'
import Components from 'unplugin-vue-components/vite'
import { NaiveUiResolver } from 'unplugin-vue-components/resolvers'

// https://vitejs.dev/config/
export default defineConfig(({ mode }) => {
  // 读取本地环境变量（.env.development / .env.production 等）
  const env = loadEnv(mode, process.cwd())
  const apiTarget = env.VITE_API_BASE_URL || 'http://localhost:8080'

  return {
    plugins: [
      vue(),
      AutoImport({
        imports: [
          'vue',
          'vue-router',
          'pinia',
          {
            'naive-ui': [
              'useDialog',
              'useMessage',
              'useNotification',
              'useLoadingBar'
            ]
          }
        ],
        dts: 'src/types/auto-imports.d.ts'
      }),
      Components({
        resolvers: [NaiveUiResolver()],
        dts: 'src/types/components.d.ts'
      })
    ],
    resolve: {
      alias: {
        '@': resolve(__dirname, 'src')
      }
    },
    server: {
      // 监听 0.0.0.0 以便在局域网 / 容器环境中通过 IP 访问开发服务器
      host: '0.0.0.0',
      port: 3000,
      proxy: {
        '/api': {
          target: apiTarget,
          changeOrigin: true,
          ws: true // 支持 WebSocket 代理
        },
        '/uploads': {
          target: apiTarget,
          changeOrigin: true
        }
      }
    },
    build: {
      rollupOptions: {
        output: {
          manualChunks: {
            vendor: ['vue', 'vue-router', 'pinia'],
            ui: ['naive-ui'],
            editor: ['@kangc/v-md-editor']
          }
        }
      }
    }
  }
})

