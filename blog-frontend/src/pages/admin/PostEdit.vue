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
              :options="tagOptions"
              multiple
              placeholder="请选择标签"
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
            <markdown-editor v-model="formData.content" height="600px" />
          </n-form-item>

          <n-form-item label="状态" path="status">
            <n-radio-group v-model:value="formData.status">
              <n-radio :value="1">发布</n-radio>
              <n-radio :value="0">草稿</n-radio>
            </n-radio-group>
          </n-form-item>

          <n-form-item label="置顶" path="is_top">
            <n-switch v-model:value="formData.is_top" />
          </n-form-item>

          <n-space>
            <n-button type="primary" size="large" :loading="submitting" @click="handleSubmit">
              {{ isEdit ? '保存修改' : '发布文章' }}
            </n-button>
            <n-button size="large" @click="handleBack">取消</n-button>
          </n-space>
        </n-form>
      </n-spin>
    </n-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
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

const isEdit = computed(() => !!route.params.id)

const formData = reactive<PostForm>({
  title: '',
  content: '',
  summary: '',
  cover: '',
  category_id: null,
  tag_ids: [],
  status: 1,
  is_top: false
})

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

const tagOptions = computed<SelectOption[]>(() => {
  return blogStore.tags.map((tag) => ({
    label: tag.name,
    value: tag.id
  }))
})

onMounted(async () => {
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
      await updatePost(id, formData)
      message.success('文章更新成功')
    } else {
      // 创建文章
      await createPost(formData)
      message.success('文章发布成功')
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
</script>

<style scoped>
.post-edit-page {
  padding: 0;
}

.edit-card {
  max-width: 1200px;
  margin: 0 auto;
}
</style>

