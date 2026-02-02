<!--
 * @ProjectName: go-vue3-blog
 * @FileName: PostEdit.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 文章编辑页面组件，提供文章的创建和编辑功能
 -->
<template>
  <div class="post-edit-page">
    <n-page-header @back="handleBack">
      <template #title>{{ isEdit ? '编辑文章' : '新建文章' }}</template>
    </n-page-header>

    <n-card class="edit-card">
      <n-spin :show="loading">
        <n-form
          ref="formRef"
          :model="formData"
          :rules="rules"
          label-placement="top"
        >
          <n-form-item label="文章标题" path="title">
            <n-input
              v-model:value="formData.title"
              placeholder="请输入文章标题"
              size="large"
            />
          </n-form-item>

          <n-form-item label="分类" path="category_id">
            <n-select
              v-model:value="formData.category_id"
              :options="categoryOptions"
              placeholder="请选择分类"
            />
          </n-form-item>

          <n-form-item label="标签" path="tag_ids">
            <n-select
              v-model:value="formData.tag_ids"
              :options="tagOptionsWithSearchHint"
              multiple
              filterable
              clearable
              :filter="filterTagOption"
              :on-search="handleTagSearch"
              placeholder="请输入关键词搜索标签"
            />
          </n-form-item>

          <n-form-item label="文章摘要" path="summary">
            <n-input
              v-model:value="formData.summary"
              type="textarea"
              :rows="3"
              placeholder="请输入文章摘要"
            />
          </n-form-item>

          <n-form-item label="封面图" path="cover">
            <n-space vertical style="width: 100%">
              <image-upload
                v-model="formData.cover"
                :width="600"
                :height="340"
                :max-size-m-b="5"
                alt="文章封面"
                @success="handleCoverSuccess"
              />
              <n-text depth="3" style="font-size: 12px">
                建议尺寸：1200x680 或 16:9 比例，支持 jpg、png、gif 格式
              </n-text>
            </n-space>
          </n-form-item>

          <n-form-item label="文章内容" path="content">
            <markdown-editor v-model="formData.content" :height="isMobile ? '400px' : '600px'" />
          </n-form-item>

          <n-form-item label="状态" path="status">
            <n-radio-group v-model:value="formData.status">
              <n-radio :value="1">发布</n-radio>
              <n-radio :value="0">草稿</n-radio>
            </n-radio-group>
          </n-form-item>

          <n-form-item label="可见性" path="visibility">
            <n-radio-group v-model:value="formData.visibility">
              <n-radio :value="1">公开</n-radio>
              <n-radio :value="0">私密（仅管理员可见）</n-radio>
            </n-radio-group>
          </n-form-item>

          <n-form-item label="置顶" path="is_top">
            <n-switch v-model:value="formData.is_top" />
          </n-form-item>

          <n-space :vertical="isMobile">
            <n-button type="primary" :size="isMobile ? 'medium' : 'large'" :block="isMobile" :loading="submitting" @click="handleSubmit">
              {{ isEdit ? '保存修改' : '发布文章' }}
            </n-button>
            <n-button :size="isMobile ? 'medium' : 'large'" :block="isMobile" @click="handleBack">取消</n-button>
          </n-space>
        </n-form>
      </n-spin>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, watch, onMounted, onUnmounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useMessage } from 'naive-ui'
import type { FormInst, FormRules, SelectOption } from 'naive-ui'
import { getPostById, createPost, updatePost } from '@/api/post'
import { useBlogStore } from '@/stores'
import type { PostForm } from '@/types/blog'
import MarkdownEditor from '@/components/MarkdownEditor.vue'
import ImageUpload from '@/components/ImageUpload.vue'

const route = useRoute()
const router = useRouter()
const message = useMessage()
const blogStore = useBlogStore()

const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const submitting = ref(false)
const isMobile = ref(false)

const isEdit = computed(() => !!route.params.id)

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

const formData = reactive<PostForm>({
  title: '',
  content: '',
  summary: '',
  cover: '',
  category_id: null,
  tag_ids: [],
  status: 1,
  visibility: 1,
  is_top: false
})

// 保存用户之前选择的可见性（用于从草稿切换回发布时恢复）
let previousVisibility = 1

const rules: FormRules = {
  title: [{ required: true, message: '请输入文章标题', trigger: 'blur' }],
  content: [{ required: true, message: '请输入文章内容', trigger: 'blur' }],
  category_id: [
    { 
      required: true, 
      message: '请选择分类', 
      trigger: ['blur', 'change'],
      validator: (_rule, value) => {
        if (value === null || value === undefined || value === 0) {
          return new Error('请选择分类')
        }
        return true
      }
    }
  ]
}

const categoryOptions = computed<SelectOption[]>(() => {
  return blogStore.categories.map((cat) => ({
    label: cat.name,
    value: cat.id
  }))
})

// 标签搜索关键词（仅用于展示当前输入的提示项）
const tagSearchKeyword = ref('')

const tagOptions = computed<SelectOption[]>(() => {
  return blogStore.tags.map((tag) => ({
    label: tag.name,
    value: tag.id
  }))
})

// 在下拉列表顶部显示“当前输入：xxx”的只读提示项，方便在无匹配结果时仍能看到自己的输入
const tagOptionsWithSearchHint = computed<SelectOption[]>(() => {
  const base = tagOptions.value
  const keyword = tagSearchKeyword.value.trim()
  if (!keyword) return base

  return [
    {
      label: `当前输入：${keyword}`,
      value: -1,
      disabled: true
    },
    ...base
  ]
})

function handleTagSearch(value: string) {
  tagSearchKeyword.value = value
}

function filterTagOption(pattern: string, option: SelectOption) {
  const keyword = (pattern || '').trim().toLowerCase()
  if (!keyword) return true

  const label = String(option.label ?? '').toLowerCase()
  // 模糊查询：包含匹配（足够解决“标签过多难找”的问题）
  return label.includes(keyword)
}

// 监听状态变化，当选择草稿时自动设置可见性为私密
watch(
  () => formData.status,
  (newStatus, oldStatus) => {
    if (newStatus === 0) {
      // 草稿状态，保存当前可见性（如果之前是发布状态），然后自动设置为私密
      if (oldStatus === 1) {
        previousVisibility = formData.visibility
      }
      formData.visibility = 0
    } else if (newStatus === 1 && oldStatus === 0) {
      // 从草稿切换回发布状态，恢复之前保存的可见性
      formData.visibility = previousVisibility
    }
  },
  { immediate: false }
)

onMounted(async () => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  
  // 加载分类和标签
  await blogStore.fetchCategories()
  await blogStore.fetchTags()

  // 如果是编辑模式，加载文章数据
  if (isEdit.value) {
    await loadPost()
  }
})

async function loadPost() {
  try {
    loading.value = true
    const id = Number(route.params.id)
    const res = await getPostById(id)
    const post = res.data

    if (!post) {
      throw new Error('文章不存在')
    }

    formData.title = post.title
    formData.content = post.content
    formData.summary = post.summary
    formData.cover = post.cover || ''
    formData.category_id = post.category_id
    formData.tag_ids = post.tags?.map((tag) => tag.id) || []
    formData.status = post.status
    formData.visibility = post.visibility ?? 1
    formData.is_top = post.is_top
  } catch (error: any) {
    message.error(error.message || '加载文章失败')
    handleBack()
  } finally {
    loading.value = false
  }
}

function handleCoverSuccess(url: string) {
  formData.cover = url
  message.success('封面图上传成功')
}

async function handleSubmit() {
  try {
    // 先进行前端表单验证
    await formRef.value?.validate()
    
    submitting.value = true

    if (isEdit.value) {
      // 更新文章
      const id = Number(route.params.id)
      // 构建更新数据，确保 category_id 字段被正确传递
      // 表单验证已确保 category_id 不为 null，但为了安全起见，再次检查
      if (formData.category_id === null || formData.category_id === undefined) {
        message.error('请选择分类')
        return
      }
      const updateData: Partial<PostForm> = {
        title: formData.title,
        content: formData.content,
        summary: formData.summary,
        cover: formData.cover,
        category_id: formData.category_id, // 确保是有效数字
        tag_ids: formData.tag_ids,
        status: formData.status,
        visibility: formData.visibility,
        is_top: formData.is_top
      }
      await updatePost(id, updateData)
      message.success(formData.status === 0 ? '草稿保存成功' : '文章更新成功')
    } else {
      // 创建文章
      await createPost(formData)
      message.success(formData.status === 0 ? '草稿保存成功' : '文章发布成功')
    }

    handleBack()
  } catch (error: any) {
    // 如果是表单验证错误，不显示错误提示（Naive UI 会自动显示）
    if (error?.errors) {
      return
    }
    
    // 处理后端返回的错误
    let errorMessage = '保存失败'
    if (error.response?.data?.message) {
      errorMessage = error.response.data.message
    } else if (error.message) {
      errorMessage = error.message
    }
    
    message.error(errorMessage)
  } finally {
    submitting.value = false
  }
}

function handleBack() {
  router.push({ name: 'PostManage' })
}

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.post-edit-page {
  padding: 0;
}

.edit-card {
  max-width: 1200px;
  margin: 0 auto;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .post-edit-page :deep(.n-page-header) {
    padding: 12px;
  }
  
  .edit-card {
    margin: 0;
  }
  
  .post-edit-page :deep(.n-form-item) {
    margin-bottom: 20px;
  }
}
</style>

