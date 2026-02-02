/*
 * @ProjectName: go-vue3-blog
 * @FileName: canvas-nest.d.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: Canvas Nest库的类型定义
 */

declare module 'canvas-nest.js' {
  interface CanvasNestConfig {
    color?: string
    pointColor?: string
    opacity?: number
    count?: number
    zIndex?: number
  }

  export default class CanvasNest {
    constructor(el?: HTMLElement | string, config?: CanvasNestConfig)
    destroy(): void
  }
}

