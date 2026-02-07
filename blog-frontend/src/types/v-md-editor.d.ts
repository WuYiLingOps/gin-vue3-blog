/*
 * @ProjectName: go-vue3-blog
 * @FileName: v-md-editor.d.ts
 * @CreateTime: 2026-02-02 11:54:38
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: Markdown编辑器库的类型定义
 */

declare module '@kangc/v-md-editor' {
  import { DefineComponent, App } from 'vue'
  const VMdEditor: DefineComponent<any, any, any> & {
    install: (app: App) => void
    use: (theme: any, config?: any) => void
  }
  export default VMdEditor
}

declare module '@kangc/v-md-editor/lib/preview' {
  import { DefineComponent, App } from 'vue'
  const VMdPreview: DefineComponent<any, any, any> & {
    install: (app: App) => void
    use: (theme: any, config?: any) => void
  }
  export default VMdPreview
}

declare module '@kangc/v-md-editor/lib/theme/vuepress.js' {
  const vuepressTheme: any
  export default vuepressTheme
}

declare module 'prismjs' {
  const Prism: any
  export default Prism
}

