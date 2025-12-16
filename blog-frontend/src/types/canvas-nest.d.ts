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

