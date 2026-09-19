/*
 * 项目名称：blog-frontend
 * 文件名称：tech-icons.ts
 * 创建时间：2026-09-19
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：技能卡技术栈图标注册表，基于 simple-icons（品牌色 + 单色路径），
 *           供关于页技能跑马灯卡片按 slug 取用。
 */
import type { SimpleIcon } from 'simple-icons'
import {
  siJavascript,
  siTypescript,
  siVuedotjs,
  siNodedotjs,
  siPython,
  siFlutter,
  siReact,
  siGit,
  siGithub,
  siGitee,
  siDocker,
  siNginx,
  siMysql,
  siRedis,
  siLinux,
  siUbuntu,
  siSpring,
  siSpringboot,
  siGo,
  siKubernetes,
  siPostgresql,
  siHtml5,
  siCss,
  siOpenjdk,
  siIntellijidea,
  siJenkins,
  siApachemaven,
  siGradle,
  siVite,
  siMarkdown
} from 'simple-icons'

export interface TechIcon extends SimpleIcon {
  /** 注册键，后台配置技能时填写 */
  key: string
}

/** 有序注册表：key → 图标（title/path/hex 来自 simple-icons） */
export const TECH_ICONS: Record<string, TechIcon> = (
  [
    ['java', siOpenjdk],
    ['javascript', siJavascript],
    ['typescript', siTypescript],
    ['vue', siVuedotjs],
    ['nodejs', siNodedotjs],
    ['python', siPython],
    ['flutter', siFlutter],
    ['react', siReact],
    ['git', siGit],
    ['github', siGithub],
    ['gitee', siGitee],
    ['docker', siDocker],
    ['nginx', siNginx],
    ['mysql', siMysql],
    ['redis', siRedis],
    ['linux', siLinux],
    ['ubuntu', siUbuntu],
    ['spring', siSpring],
    ['springboot', siSpringboot],
    ['go', siGo],
    ['kubernetes', siKubernetes],
    ['postgresql', siPostgresql],
    ['html', siHtml5],
    ['css', siCss],
    ['intellijidea', siIntellijidea],
    ['jenkins', siJenkins],
    ['maven', siApachemaven],
    ['gradle', siGradle],
    ['vite', siVite],
    ['markdown', siMarkdown]
  ] as Array<[string, SimpleIcon]>
).map(([key, icon]) => ({ ...icon, key }))
  .reduce((acc, icon) => {
    acc[icon.key] = icon
    return acc
  }, {} as Record<string, TechIcon>)

/** 技能项解析：icon 键匹配注册表则返回品牌图标，否则返回 null（由调用方降级为首字母方块） */
export function getTechIcon(key: string | undefined): TechIcon | null {
  if (!key) return null
  return TECH_ICONS[key.trim().toLowerCase()] || null
}
