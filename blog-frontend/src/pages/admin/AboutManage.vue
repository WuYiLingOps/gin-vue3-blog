<!--
 * @ProjectName: go-vue3-blog
 * @FileName: AboutManage.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 关于我管理页面（标签页布局）：正文 Markdown 编辑独立成卡；
 *               页面卡片配置按板块分为 7 个标签页（作者介绍 / 个人信息 /
 *               MBTI 性格 / 打字词轮播 / 技能 / 成长轨迹 / 追求卡），
 *               保存与重置按钮作用于全部板块，结构化配置存 settings 表
 *               site 组的 about_* 系列键（JSON 字符串）。
 -->
<template>
  <div class="about-manage-page">
    <!-- 正文内容 -->
    <n-card title="关于我正文" class="manage-card">
      <n-spin :show="loading">
        <n-form :model="contentForm" label-placement="top">
          <n-form-item label="正文内容（支持 Markdown 格式）" path="content">
            <markdown-editor
              v-model="contentForm.content"
              :height="isMobile ? '500px' : '700px'"
              :subfield="!isMobile"
              :mode="isMobile ? 'edit' : 'editable'"
            />
            <template #feedback>
              <n-text depth="3" style="font-size: 12px; margin-top: 8px; display: block">
                提示：支持 Markdown 语法，内容将显示在"关于我"页面的"自我介绍"卡片中
              </n-text>
            </template>
          </n-form-item>

          <n-space :vertical="isMobile" style="margin-top: 24px">
            <n-button
              type="primary"
              :size="isMobile ? 'medium' : 'large'"
              :block="isMobile"
              :loading="submitting"
              @click="handleSubmitContent"
            >
              保存正文
            </n-button>
            <n-button
              :size="isMobile ? 'medium' : 'large'"
              :block="isMobile"
              @click="handleResetContent"
            >
              重置正文
            </n-button>
          </n-space>
        </n-form>
      </n-spin>
    </n-card>

    <!-- 卡片配置 -->
    <n-card title="页面卡片配置" class="manage-card">
      <n-spin :show="configLoading">
        <n-tabs v-model:value="activeTab" type="line" animated display-directive="show">
          <!-- 作者介绍卡 -->
          <n-tab-pane name="intro" tab="作者介绍">
            <div class="form-grid">
              <n-form-item label="问候语" path="tips">
                <n-input v-model:value="introForm.tips" placeholder="如：你好，很高兴认识你 👋" />
              </n-form-item>
              <n-form-item label="名字前缀" path="namePrefix">
                <n-input v-model:value="introForm.namePrefix" placeholder="如：我叫" />
              </n-form-item>
              <n-form-item label="角色前缀" path="descPrefix">
                <n-input v-model:value="introForm.descPrefix" placeholder="如：是一名" />
              </n-form-item>
            </div>
            <n-form-item label="角色列表（以、连接展示）" path="roles">
              <n-dynamic-tags v-model:value="introForm.roles" />
            </n-form-item>
          </n-tab-pane>

          <!-- 个人信息卡 -->
          <n-tab-pane name="selfinfo" tab="个人信息">
            <div class="form-grid">
              <n-form-item label="栏目一小字" path="tips1">
                <n-input v-model:value="selfInfoForm.item1.tips" placeholder="如：毕业时间" />
              </n-form-item>
              <n-form-item label="栏目一大字" path="value1">
                <n-input v-model:value="selfInfoForm.item1.value" placeholder="如：2026" />
              </n-form-item>
              <n-form-item label="栏目二小字" path="tips2">
                <n-input v-model:value="selfInfoForm.item2.tips" placeholder="如：桂林电子科技大学" />
              </n-form-item>
              <n-form-item label="栏目二大字" path="value2">
                <n-input v-model:value="selfInfoForm.item2.value" placeholder="如：大数据管理与应用" />
              </n-form-item>
              <n-form-item label="栏目三小字" path="tips3">
                <n-input v-model:value="selfInfoForm.item3.tips" placeholder="如：现在职业" />
              </n-form-item>
              <n-form-item label="栏目三大字" path="value3">
                <n-input v-model:value="selfInfoForm.item3.value" placeholder="如：运维工程师" />
              </n-form-item>
            </div>
          </n-tab-pane>

          <!-- MBTI 性格卡 -->
          <n-tab-pane name="personality" tab="MBTI 性格">
            <div class="form-grid">
              <n-form-item label="引导语" path="tips">
                <n-input v-model:value="personalityForm.tips" placeholder="如：性格" />
              </n-form-item>
              <n-form-item label="人格名称 + 代码" path="type">
                <n-input v-model:value="personalityForm.type" placeholder="如：执政官 ESFJ-A" />
              </n-form-item>
              <n-form-item label="人格文字颜色" path="color">
                <n-input v-model:value="personalityForm.color" placeholder="留空使用 #ac899c" />
              </n-form-item>
              <n-form-item label="了解更多链接（可选）" path="link">
                <n-input v-model:value="personalityForm.link" placeholder="留空指向 16personalities 官网" />
              </n-form-item>
              <n-form-item label="人格形象图URL（可选，SVG/PNG）" path="img">
                <n-input v-model:value="personalityForm.img" placeholder="留空则不显示形象图" />
              </n-form-item>
            </div>
          </n-tab-pane>

          <!-- 地理位置 -->
          <n-tab-pane name="map" tab="地理位置">
            <div class="form-grid">
              <n-form-item label="位置条前缀文案" path="title">
                <n-input v-model:value="mapForm.title" placeholder="如：我现在住在" />
              </n-form-item>
              <n-form-item label="位置文案（加粗展示）" path="location">
                <n-input v-model:value="mapForm.location" placeholder="如：中国，桂林市" />
              </n-form-item>
              <n-form-item label="纬度" path="lat">
                <n-input-number
                  v-model:value="mapForm.lat"
                  :precision="4"
                  :step="0.01"
                  :min="-85"
                  :max="85"
                  style="width: 100%"
                />
              </n-form-item>
              <n-form-item label="经度" path="lng">
                <n-input-number
                  v-model:value="mapForm.lng"
                  :precision="4"
                  :step="0.01"
                  :min="-180"
                  :max="180"
                  style="width: 100%"
                />
              </n-form-item>
              <n-form-item label="缩放级别（3-16，越大越详细）" path="zoom">
                <n-input-number
                  v-model:value="mapForm.zoom"
                  :precision="0"
                  :min="3"
                  :max="16"
                  style="width: 100%"
                />
              </n-form-item>
            </div>
          </n-tab-pane>

          <!-- 打字词轮播 -->
          <n-tab-pane name="tips" tab="打字词轮播">
            <div class="form-grid">
              <n-form-item label="引导语" path="tips">
                <n-input v-model:value="tipsForm.tips" placeholder="如：一个热爱技术与分享的博主" />
              </n-form-item>
              <n-form-item label="轮播词前缀" path="title1">
                <n-input v-model:value="tipsForm.title1" placeholder="如：当前状态" />
              </n-form-item>
              <n-form-item label="轮播词后缀" path="title2">
                <n-input v-model:value="tipsForm.title2" placeholder="如：专注方向" />
              </n-form-item>
            </div>
            <n-form-item label="轮播词列表" path="word">
              <n-dynamic-tags v-model:value="tipsForm.word" />
            </n-form-item>
          </n-tab-pane>

          <!-- 技能 -->
          <n-tab-pane name="skills" tab="技能">
            <div class="form-grid">
              <n-form-item label="引导语" path="tips">
                <n-input v-model:value="skillsForm.tips" placeholder="如：技能熟练度" />
              </n-form-item>
              <n-form-item label="卡片标题" path="title">
                <n-input v-model:value="skillsForm.title" placeholder="如：开启创造力" />
              </n-form-item>
            </div>
            <n-form-item label="头像左侧浮动标签" path="left">
              <n-dynamic-tags v-model:value="skillsForm.left" />
            </n-form-item>
            <n-form-item label="头像右侧浮动标签" path="right">
              <n-dynamic-tags v-model:value="skillsForm.right" />
            </n-form-item>
            <n-form-item label="技能图标列表（颜色可留空，留空用品牌色）" path="list">
              <div class="row-list">
                <div v-for="(skill, index) in skillsForm.list" :key="index" class="row-item">
                  <n-input
                    v-model:value="skill.name"
                    placeholder="技能名称，如 Java"
                    class="row-main-input"
                  />
                  <n-select
                    v-model:value="skill.icon"
                    :options="techIconOptions"
                    placeholder="图标"
                    filterable
                    clearable
                    class="row-icon-select"
                  />
                  <n-input
                    v-model:value="skill.color"
                    placeholder="底色可留空"
                    class="row-color-input"
                  />
                  <n-button quaternary type="error" @click="removeListItem(skillsForm.list, index)">
                    删除
                  </n-button>
                </div>
                <n-button
                  dashed
                  block
                  @click="skillsForm.list.push({ name: '', icon: undefined, color: '' })"
                >
                  添加技能
                </n-button>
              </div>
            </n-form-item>
          </n-tab-pane>

          <!-- 成长轨迹 -->
          <n-tab-pane name="careers" tab="成长轨迹">
            <div class="form-grid">
              <n-form-item label="引导语" path="tips">
                <n-input v-model:value="careersForm.tips" placeholder="如：一路走来的足迹" />
              </n-form-item>
              <n-form-item label="卡片标题" path="title">
                <n-input v-model:value="careersForm.title" placeholder="如：成长轨迹" />
              </n-form-item>
            </div>
            <n-form-item label="时间线列表（按时间顺序展示）" path="list">
              <div class="row-list">
                <div v-for="(item, index) in careersForm.list" :key="index" class="row-item">
                  <n-input v-model:value="item.time" placeholder="时间" class="row-time-input" />
                  <n-input v-model:value="item.title" placeholder="标题" class="row-title-input" />
                  <n-input v-model:value="item.desc" placeholder="描述" class="row-main-input" />
                  <n-input v-model:value="item.color" placeholder="节点颜色，可留空" class="row-color-input" />
                  <n-button quaternary type="error" @click="removeListItem(careersForm.list, index)">
                    删除
                  </n-button>
                </div>
                <n-button
                  dashed
                  block
                  @click="careersForm.list.push({ time: '', title: '', desc: '', color: '' })"
                >
                  添加时间线节点
                </n-button>
              </div>
            </n-form-item>
          </n-tab-pane>

          <!-- 追求卡 -->
          <n-tab-pane name="maxim" tab="追求卡">
            <div class="form-grid">
              <n-form-item label="引导语" path="tips">
                <n-input v-model:value="maximForm.tips" placeholder="如：追求" />
              </n-form-item>
            </div>
            <div class="form-grid">
              <n-form-item label="正文（支持换行，深色大字）" path="top">
                <n-input
                  v-model:value="maximForm.top"
                  type="textarea"
                  :autosize="{ minRows: 1, maxRows: 4 }"
                  placeholder="如：源于&#10;热爱而去 感受"
                />
              </n-form-item>
            </div>
            <n-form-item label="末行轮换词列表（每 2 秒轮换并按位置循环四种渐变色；留空则显示静态点缀词）" path="word">
              <n-dynamic-tags v-model:value="maximForm.word" />
            </n-form-item>
            <n-form-item label="静态点缀词（仅当轮换词列表为空时显示）" path="bottom">
              <n-input v-model:value="maximForm.bottom" placeholder="如：程序" />
            </n-form-item>
          </n-tab-pane>
        </n-tabs>

        <!-- 悬浮保存条：长表单滚动时保持可见 -->
        <div class="config-actions">
          <n-space :vertical="isMobile">
            <n-button
              type="primary"
              :size="isMobile ? 'medium' : 'large'"
              :block="isMobile"
              :loading="configSubmitting"
              @click="handleSaveConfig"
            >
              保存卡片配置
            </n-button>
            <n-button
              :size="isMobile ? 'medium' : 'large'"
              :block="isMobile"
              @click="handleResetConfig"
            >
              重置卡片配置
            </n-button>
          </n-space>
        </div>
      </n-spin>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { useMessage } from 'naive-ui'
import {
  getAboutInfo,
  updateAboutInfo,
  getSiteSettings,
  updateSiteSettings
} from '@/api/setting'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import {
  parseAboutJson,
  DEFAULT_ABOUT_SITE_TIPS,
  DEFAULT_ABOUT_SKILLS,
  DEFAULT_ABOUT_CAREERS,
  DEFAULT_ABOUT_MAXIM,
  DEFAULT_ABOUT_MAP,
  DEFAULT_ABOUT_INTRO,
  DEFAULT_ABOUT_SELF_INFO,
  DEFAULT_ABOUT_PERSONALITY,
  type AboutSiteTips,
  type AboutSkills,
  type AboutCareers,
  type AboutQuote,
  type AboutMapConfig,
  type AboutIntro,
  type AboutSelfInfo,
  type AboutPersonality
} from '@/types/about'
import { TECH_ICONS } from '@/utils/tech-icons'

const message = useMessage()
const loading = ref(false)
const submitting = ref(false)
const configLoading = ref(false)
const configSubmitting = ref(false)
const isMobile = ref(false)
const isSmallScreen = ref(false)
const activeTab = ref('intro')

// 正文（settings 表 about 组，走 /admin/about）
const contentForm = reactive({
  content: ''
})
const originalContent = ref('')

// 卡片配置（settings 表 site 组 about_* 键，JSON 字符串）
const tipsForm = ref<AboutSiteTips>(cloneDefault(DEFAULT_ABOUT_SITE_TIPS))
const skillsForm = ref<AboutSkills>(cloneDefault(DEFAULT_ABOUT_SKILLS))
const careersForm = ref<AboutCareers>(cloneDefault(DEFAULT_ABOUT_CAREERS))
const maximForm = ref<AboutQuote>(cloneDefault(DEFAULT_ABOUT_MAXIM))
const mapForm = ref<AboutMapConfig>(cloneDefault(DEFAULT_ABOUT_MAP))
const introForm = ref<AboutIntro>(cloneDefault(DEFAULT_ABOUT_INTRO))
const selfInfoForm = ref<AboutSelfInfo>(cloneDefault(DEFAULT_ABOUT_SELF_INFO))
const personalityForm = ref<AboutPersonality>(cloneDefault(DEFAULT_ABOUT_PERSONALITY))
const originalConfig = ref('')

// 技能图标下拉选项
const techIconOptions = Object.values(TECH_ICONS).map(icon => ({
  label: `${icon.title}（${icon.key}）`,
  value: icon.key
}))

function cloneDefault<T>(value: T): T {
  return JSON.parse(JSON.stringify(value)) as T
}

function removeListItem(list: Array<Record<string, string>>, index: number) {
  list.splice(index, 1)
}

// 检测移动设备
function checkMobile() {
  const width = window.innerWidth
  isMobile.value = width <= 1100
  isSmallScreen.value = width <= 600
}

// 获取正文内容
async function fetchAboutInfo() {
  try {
    loading.value = true
    const res = await getAboutInfo()
    if (res.data) {
      contentForm.content = res.data.content || ''
      originalContent.value = res.data.content || ''
    }
  } catch (error: any) {
    message.error(error.message || '获取关于我信息失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 保存正文
async function handleSubmitContent() {
  try {
    submitting.value = true
    await updateAboutInfo(contentForm.content)
    message.success('保存成功')
    originalContent.value = contentForm.content
  } catch (error: any) {
    message.error(error.message || '保存失败')
    console.error(error)
  } finally {
    submitting.value = false
  }
}

// 重置正文
function handleResetContent() {
  contentForm.content = originalContent.value
  message.info('已重置为上次保存的内容')
}

// 获取卡片配置
async function fetchConfig() {
  try {
    configLoading.value = true
    const res = await getSiteSettings()
    if (res.data) {
      tipsForm.value = parseAboutJson(res.data.about_site_tips, DEFAULT_ABOUT_SITE_TIPS)
      skillsForm.value = parseAboutJson(res.data.about_skills, DEFAULT_ABOUT_SKILLS)
      careersForm.value = parseAboutJson(res.data.about_careers, DEFAULT_ABOUT_CAREERS)
      maximForm.value = parseAboutJson(res.data.about_maxim, DEFAULT_ABOUT_MAXIM)
      mapForm.value = parseAboutJson(res.data.about_map, DEFAULT_ABOUT_MAP)
      introForm.value = parseAboutJson(res.data.about_intro, DEFAULT_ABOUT_INTRO)
      selfInfoForm.value = parseAboutJson(res.data.about_self_info, DEFAULT_ABOUT_SELF_INFO)
      personalityForm.value = parseAboutJson(res.data.about_personality, DEFAULT_ABOUT_PERSONALITY)
    }
    originalConfig.value = JSON.stringify(snapshotConfig())
  } catch (error: any) {
    message.error(error.message || '获取卡片配置失败')
    console.error(error)
  } finally {
    configLoading.value = false
  }
}

// 组装当前配置（过滤空白项后序列化）
function snapshotConfig(): Record<string, string> {
  return {
    about_site_tips: JSON.stringify({
      ...tipsForm.value,
      word: tipsForm.value.word.filter(w => w.trim())
    }),
    about_skills: JSON.stringify({
      ...skillsForm.value,
      left: skillsForm.value.left.filter(t => t.trim()),
      right: skillsForm.value.right.filter(t => t.trim()),
      list: skillsForm.value.list.filter(s => s.name.trim())
    }),
    about_careers: JSON.stringify({
      ...careersForm.value,
      list: careersForm.value.list.filter(c => c.title.trim() || c.desc.trim())
    }),
    about_maxim: JSON.stringify({
      ...maximForm.value,
      word: maximForm.value.word?.filter(w => w.trim())
    }),
    about_map: JSON.stringify(mapForm.value),
    about_intro: JSON.stringify({
      ...introForm.value,
      roles: introForm.value.roles.filter(r => r.trim())
    }),
    about_self_info: JSON.stringify(selfInfoForm.value),
    about_personality: JSON.stringify({
      ...personalityForm.value,
      link: personalityForm.value.link?.trim() || undefined,
      img: personalityForm.value.img?.trim() || undefined
    })
  }
}

// 保存卡片配置
async function handleSaveConfig() {
  try {
    configSubmitting.value = true
    await updateSiteSettings(snapshotConfig())
    message.success('卡片配置保存成功')
    originalConfig.value = JSON.stringify(snapshotConfig())
  } catch (error: any) {
    message.error(error.message || '保存失败')
    console.error(error)
  } finally {
    configSubmitting.value = false
  }
}

// 重置卡片配置
function handleResetConfig() {
  try {
    const snap = JSON.parse(originalConfig.value)
    tipsForm.value = snap.about_site_tips
      ? parseAboutJson(snap.about_site_tips, DEFAULT_ABOUT_SITE_TIPS)
      : cloneDefault(DEFAULT_ABOUT_SITE_TIPS)
    skillsForm.value = snap.about_skills
      ? parseAboutJson(snap.about_skills, DEFAULT_ABOUT_SKILLS)
      : cloneDefault(DEFAULT_ABOUT_SKILLS)
    careersForm.value = snap.about_careers
      ? parseAboutJson(snap.about_careers, DEFAULT_ABOUT_CAREERS)
      : cloneDefault(DEFAULT_ABOUT_CAREERS)
    maximForm.value = snap.about_maxim
      ? parseAboutJson(snap.about_maxim, DEFAULT_ABOUT_MAXIM)
      : cloneDefault(DEFAULT_ABOUT_MAXIM)
    mapForm.value = snap.about_map
      ? parseAboutJson(snap.about_map, DEFAULT_ABOUT_MAP)
      : cloneDefault(DEFAULT_ABOUT_MAP)
    introForm.value = snap.about_intro
      ? parseAboutJson(snap.about_intro, DEFAULT_ABOUT_INTRO)
      : cloneDefault(DEFAULT_ABOUT_INTRO)
    selfInfoForm.value = snap.about_self_info
      ? parseAboutJson(snap.about_self_info, DEFAULT_ABOUT_SELF_INFO)
      : cloneDefault(DEFAULT_ABOUT_SELF_INFO)
    personalityForm.value = snap.about_personality
      ? parseAboutJson(snap.about_personality, DEFAULT_ABOUT_PERSONALITY)
      : cloneDefault(DEFAULT_ABOUT_PERSONALITY)
    message.info('已重置为上次保存的配置')
  } catch {
    message.error('重置失败')
  }
}

onMounted(() => {
  checkMobile()
  fetchAboutInfo()
  fetchConfig()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.about-manage-page {
  padding: 20px;
}

.manage-card {
  margin-bottom: 20px;
}

/* 标签页内容区最小高度，减少切换时的高度跳动 */
.manage-card :deep(.n-tabs .n-tab-pane) {
  min-height: 220px;
  padding-top: 20px;
}

/* 悬浮保存条：长表单滚动到任意位置都可见 */
.config-actions {
  position: sticky;
  bottom: 12px;
  z-index: 10;
  display: flex;
  justify-content: flex-end;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.92);
  backdrop-filter: blur(8px);
  border: 1px solid rgba(227, 232, 247, 1);
  border-radius: 8px;
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.08);
}

.form-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 0 16px;
}

.row-list {
  display: flex;
  flex-direction: column;
  gap: 10px;
  width: 100%;
}

.row-item {
  display: flex;
  align-items: center;
  gap: 8px;
  flex-wrap: wrap;
}

.row-main-input {
  flex: 1;
  min-width: 180px;
}

.row-title-input {
  width: 160px;
}

.row-time-input {
  width: 110px;
}

.row-color-input {
  width: 150px;
}

.row-icon-select {
  width: 170px;
}

@media (max-width: 1100px) {
  .about-manage-page {
    padding: 12px;
  }

  .row-time-input {
    width: 96px;
  }

  .row-title-input {
    width: 130px;
  }
}
</style>
