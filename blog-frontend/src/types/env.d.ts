/*
 * @ProjectName: go-vue3-blog
 * @FileName: env.d.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 环境变量类型定义，引用Vite客户端类型
 */

/// <reference types="vite/client" />

declare module '*.vue' {
  import type { DefineComponent } from 'vue'
  const component: DefineComponent<{}, {}, any>
  export default component
}

declare interface Window {
  // 这里可以添加全局 window 对象的类型定义
}