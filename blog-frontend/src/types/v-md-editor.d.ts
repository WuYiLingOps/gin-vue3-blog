declare module '@kangc/v-md-editor' {
  import { DefineComponent } from 'vue'
  const VMdEditor: DefineComponent<any, any, any>
  export default VMdEditor
}

declare module '@kangc/v-md-editor/lib/preview' {
  import { DefineComponent } from 'vue'
  const VMdPreview: DefineComponent<any, any, any>
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

