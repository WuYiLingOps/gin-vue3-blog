/*
 * 项目名称：blog-frontend
 * 文件名称：about.ts
 * 创建时间：2026-09-18
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：关于我页面卡片配置的类型定义与默认值，对应 settings 表 site 组的
 *           about_* 系列键（JSON 字符串存储，参照 cover_bg_images 先例）。
 */

/** 打字词轮播配置（about_site_tips） */
export interface AboutSiteTips {
  tips: string // 卡片顶部小字引导语
  title1: string // 打字词前缀文字
  title2: string // 打字词后缀文字
  word: string[] // 轮播词列表
}

/** 技能配置（about_skills） */
export interface AboutSkillItem {
  name: string
  icon?: string // tech-icons 注册表键名，如 vue/java/docker，缺省显示首字母方块
  color?: string // 图标方块底色，缺省用品牌色
}

export interface AboutSkills {
  tips: string
  title: string
  left: string[] // 头像左侧浮动标签
  right: string[] // 头像右侧浮动标签
  list: AboutSkillItem[] // 技能卡图标列表
}

/** 作者介绍卡配置（about_intro） */
export interface AboutIntro {
  tips: string // 问候语，如"你好，很高兴认识你 👋"
  namePrefix: string // 名字前缀，如"我叫"
  descPrefix: string // 角色前缀，如"是一名"
  roles: string[] // 角色列表，以、连接展示
}

/** 个人信息卡配置（about_self_info），与 anheyu self_info 同款三栏结构 */
export interface AboutSelfInfoItem {
  tips: string // 栏目小字，如"生于"
  value: string // 大字内容，如"2003"
}

export interface AboutSelfInfo {
  item1: AboutSelfInfoItem
  item2: AboutSelfInfoItem
  item3: AboutSelfInfoItem
}

/** MBTI 性格卡配置（about_personality），与 anheyu personalities 同款结构 */
export interface AboutPersonality {
  tips: string // 引导语，如"性格"
  type: string // 人格名称 + 代码，如"执政官 ESFJ-A"
  color?: string // 人格文字颜色，缺省 #ac899c
  link?: string // "了解更多"链接地址，缺省指向 16personalities 官网
  img?: string // 人格形象图URL（16personalities 角色图，显示在卡片右上，悬浮旋转）
}

/** 职业生涯时间线配置（about_careers） */
export interface AboutCareerItem {
  time: string
  title: string
  desc: string
  color?: string // 时间线节点颜色，缺省用主题色
}

export interface AboutCareers {
  tips: string
  title: string
  list: AboutCareerItem[]
}

/** 座右铭 / 追求卡配置（about_maxim），bottom 为无轮换词时的静态降级展示 */
export interface AboutQuote {
  tips: string
  top: string
  bottom?: string
  word?: string[] // 末行轮换词（anheyu 同款动态效果，按位置循环四种渐变色）
}

/** 地理位置卡配置（about_map），地图由高德瓦片实时渲染（经纬度可配） */
export interface AboutMapConfig {
  title: string // 位置条前缀文案，如"我现在住在"
  location: string // 位置文案（加粗展示），如"中国，桂林市"
  lat: number // 纬度
  lng: number // 经度
  zoom: number // 地图缩放级别 3-17（渲染时自动 +1 加密瓦片保证高清）
}

/** 解析 settings 中的 JSON 配置，失败时返回默认值 */
export function parseAboutJson<T>(raw: string | undefined | null, fallback: T): T {
  if (!raw) return fallback
  try {
    const parsed = JSON.parse(raw)
    if (parsed && typeof parsed === 'object') {
      return { ...fallback, ...parsed } as T
    }
    return fallback
  } catch {
    return fallback
  }
}

export const DEFAULT_ABOUT_SITE_TIPS: AboutSiteTips = {
  tips: '一个热爱技术与分享的博主',
  title1: '当前状态',
  title2: '专注方向',
  word: ['后端开发', '云原生运维', '开源爱好者', '持续学习中']
}

export const DEFAULT_ABOUT_SKILLS: AboutSkills = {
  tips: '技能',
  title: '开启创造力',
  left: ['Java', 'Go', 'Vue3', 'TypeScript'],
  right: ['Docker', 'Kubernetes', 'Linux', 'PostgreSQL'],
  list: [
    { name: 'Java', icon: 'java', color: '#f89820' },
    { name: 'JavaScript', icon: 'javascript' },
    { name: 'TypeScript', icon: 'typescript' },
    { name: 'Vue3', icon: 'vue' },
    { name: 'Python', icon: 'python' },
    { name: 'Go', icon: 'go' },
    { name: 'Spring Boot', icon: 'springboot' },
    { name: 'Docker', icon: 'docker' },
    { name: 'Kubernetes', icon: 'kubernetes' },
    { name: 'Linux', icon: 'linux' },
    { name: 'MySQL', icon: 'mysql' },
    { name: 'Redis', icon: 'redis' },
    { name: 'Nginx', icon: 'nginx' },
    { name: 'PostgreSQL', icon: 'postgresql' },
    { name: 'Git', icon: 'git' },
    { name: 'Jenkins', icon: 'jenkins' }
  ]
}

export const DEFAULT_ABOUT_INTRO: AboutIntro = {
  tips: '你好，很高兴认识你 👋',
  namePrefix: '我叫',
  descPrefix: '是一名',
  roles: ['学生', '后端开发', '运维折腾爱好者', '技术博主']
}

export const DEFAULT_ABOUT_SELF_INFO: AboutSelfInfo = {
  item1: { tips: '毕业时间', value: '2026' },
  item2: { tips: '桂林电子科技大学', value: '大数据管理与应用' },
  item3: { tips: '现在职业', value: '运维工程师' }
}

export const DEFAULT_ABOUT_PERSONALITY: AboutPersonality = {
  tips: '性格',
  type: '执政官 ESFJ-A',
  color: '#ac899c',
  img: 'https://npm.elemecdn.com/anzhiyu-blog@2.0.8/img/svg/ESFJ-A.svg'
}

export const DEFAULT_ABOUT_CAREERS: AboutCareers = {
  tips: '一路走来的足迹',
  title: '成长轨迹',
  list: [
    { time: '2022', title: '初识代码', desc: '接触编程与 Linux，从双系统到云服务器，折腾之路开始', color: '#0891b2' },
    { time: '2024', title: '沉迷云原生', desc: 'Kubernetes、CI/CD、监控告警越折腾越上头，确定运维方向', color: '#059669' },
    { time: '2026', title: '踏入职场', desc: '成为一名运维工程师，把折腾变成热爱的事业，继续前行', color: '#9a60b4' }
  ]
}

export const DEFAULT_ABOUT_MAXIM: AboutQuote = {
  tips: '追求',
  top: '源于\n热爱而去 感受',
  bottom: '程序',
  word: ['学习', '生活', '程序', '体验']
}

export const DEFAULT_ABOUT_MAP: AboutMapConfig = {
  title: '我现在住在',
  location: '中国，桂林市',
  lat: 25.274,
  lng: 110.29,
  zoom: 10
}
