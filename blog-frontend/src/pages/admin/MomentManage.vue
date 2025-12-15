<template>
  <div class="moment-manage-page">
    <n-card title="说说管理" :bordered="false">
      <template #header-extra>
        <n-button type="primary" :size="isMobile ? 'small' : 'medium'" @click="showCreateModal = true">
          <template #icon>
            <n-icon :component="AddOutline" />
          </template>
          <span v-if="!isMobile">发布说说</span>
          <span v-else>发布</span>
        </n-button>
      </template>

      <!-- 筛选栏 -->
      <n-space :size="16" style="margin-bottom: 16px">
        <n-select
          v-model:value="filterStatus"
          :options="statusOptions"
          placeholder="选择状态"
          :style="{ width: isMobile ? '100%' : '150px' }"
          @update:value="handleSearch"
        />
      </n-space>

      <!-- 说说列表 -->
      <n-spin :show="loading">
        <n-space vertical :size="16">
          <div v-for="moment in moments" :key="moment.id" class="moment-item">
            <div class="moment-header">
              <n-space :size="12">
                <n-avatar :src="moment.user.avatar" :size="40" round />
                <div>
                  <div class="moment-user">{{ moment.user.nickname }}</div>
                  <div class="moment-time">{{ formatDate(moment.created_at, 'YYYY-MM-DD HH:mm') }}</div>
                </div>
              </n-space>
              <n-space>
                <n-tag :type="moment.status === 1 ? 'success' : 'default'" size="small">
                  {{ moment.status === 1 ? '公开' : '私密' }}
                </n-tag>
              </n-space>
            </div>

            <div class="moment-content">
              {{ moment.content }}
            </div>

            <div v-if="moment.images" class="moment-images">
              <n-space :size="8">
                <n-image
                  v-for="(img, index) in parseImages(moment.images)"
                  :key="index"
                  :src="img"
                  width="80"
                  height="80"
                  object-fit="cover"
                  style="border-radius: 4px"
                />
              </n-space>
            </div>

            <div class="moment-footer">
              <n-space :size="16">
                <span class="stat-item">
                  <n-icon :component="HeartOutline" />
                  {{ moment.like_count }}
                </span>
              </n-space>

              <n-space :size="12">
                <n-button size="small" @click="handleEdit(moment)">编辑</n-button>
                <n-popconfirm @positive-click="handleDelete(moment.id)">
                  <template #trigger>
                    <n-button size="small" type="error">删除</n-button>
                  </template>
                  确定要删除这条说说吗？
                </n-popconfirm>
              </n-space>
            </div>
          </div>
        </n-space>

        <n-empty v-if="!loading && moments.length === 0" description="暂无说说" style="margin: 40px 0" />
      </n-spin>

      <!-- 分页 -->
      <div v-if="total > 0" style="margin-top: 24px; display: flex; justify-content: flex-end">
        <n-pagination
          v-model:page="currentPage"
          :page-count="Math.ceil(total / pageSize)"
          :page-size="pageSize"
          show-size-picker
          :page-sizes="[10, 20, 30]"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </div>
    </n-card>

    <!-- 创建/编辑说说对话框 -->
    <n-modal
      v-model:show="showCreateModal"
      preset="card"
      :title="editingMoment ? '编辑说说' : '发布说说'"
      :style="{ width: isMobile ? '95%' : '600px', maxWidth: isMobile ? '95vw' : '600px' }"
      :bordered="false"
    >
      <n-form ref="formRef" :model="formData" :rules="rules" label-placement="top">
        <n-form-item label="内容" path="content">
          <n-input
            v-model:value="formData.content"
            type="textarea"
            placeholder="分享你的想法..."
            :rows="6"
            maxlength="500"
            show-count
          />
        </n-form-item>

        <n-form-item label="图片">
          <MultiImageUpload
            v-model="formData.images"
            :max-count="9"
            @success="handleImageSuccess"
          />
        </n-form-item>

        <n-form-item label="状态" path="status">
          <n-radio-group v-model:value="formData.status">
            <n-radio :value="1">公开</n-radio>
            <n-radio :value="0">私密</n-radio>
          </n-radio-group>
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showCreateModal = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">
            {{ editingMoment ? '更新' : '发布' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted } from 'vue'
import { AddOutline, HeartOutline } from '@vicons/ionicons5'
import { useMessage } from 'naive-ui'
import { getAdminMoments, createMoment, updateMoment, deleteMoment } from '@/api/moment'
import type { Moment, MomentForm } from '@/api/moment'
import { formatDate } from '@/utils/format'
import MultiImageUpload from '@/components/MultiImageUpload.vue'

const message = useMessage()
const loading = ref(false)
const submitting = ref(false)
const showCreateModal = ref(false)
const editingMoment = ref<Moment | null>(null)
const moments = ref<Moment[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const filterStatus = ref<number | null>(null)
const isMobile = ref(false)

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

const statusOptions = [
  { label: '全部', value: null as any },
  { label: '公开', value: 1 },
  { label: '私密', value: 0 }
]

const formData = reactive<MomentForm>({
  content: '',
  images: '',
  status: 1
})

const rules = {
  content: {
    required: true,
    message: '请输入说说内容',
    trigger: 'blur'
  }
}

const formRef = ref()

// 解析图片
function parseImages(images: string): string[] {
  try {
    return images ? JSON.parse(images) : []
  } catch {
    return []
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchMoments()
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

// 获取说说列表
async function fetchMoments() {
  try {
    loading.value = true
    const params: any = {
      page: currentPage.value,
      page_size: pageSize.value
    }
    
    if (filterStatus.value !== null) {
      params.status = filterStatus.value
    }

    const res = await getAdminMoments(params)
    if (res.data) {
      moments.value = res.data.list || []
      total.value = res.data.total || 0
    }
  } catch (error) {
    message.error('获取说说列表失败')
  } finally {
    loading.value = false
  }
}

// 搜索
function handleSearch() {
  currentPage.value = 1
  fetchMoments()
}

// 图片上传成功
function handleImageSuccess(_urls: string[]) {
  // MultiImageUpload组件已经返回JSON字符串，不需要再次转换
  // formData.images已经通过v-model自动更新
}

// 编辑说说
function handleEdit(moment: Moment) {
  editingMoment.value = moment
  formData.content = moment.content
  formData.images = moment.images
  formData.status = moment.status
  showCreateModal.value = true
}

// 提交表单
async function handleSubmit() {
  try {
    await formRef.value?.validate()
    submitting.value = true

    const normalizedStatus = formData.status === undefined || formData.status === null
      ? 1
      : Number(formData.status)

    if (editingMoment.value) {
      await updateMoment(editingMoment.value.id, {
        content: formData.content,
        images: formData.images,
        status: normalizedStatus
      })
      message.success('更新成功')
    } else {
      const payload: MomentForm = {
        content: formData.content,
        images: formData.images,
        status: normalizedStatus
      }
      await createMoment(payload)
      message.success('发布成功')
    }

    showCreateModal.value = false
    resetForm()
    fetchMoments()
  } catch (error: any) {
    if (error?.message) {
      message.error(error.message)
    }
  } finally {
    submitting.value = false
  }
}

// 删除说说
async function handleDelete(id: number) {
  try {
    await deleteMoment(id)
    message.success('删除成功')
    fetchMoments()
  } catch (error) {
    message.error('删除失败')
  }
}

// 页码变化
function handlePageChange(page: number) {
  currentPage.value = page
  fetchMoments()
}

// 每页数量变化
function handlePageSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  fetchMoments()
}

// 重置表单
function resetForm() {
  formData.content = ''
  formData.images = ''
  formData.status = 1
  editingMoment.value = null
}

onMounted(() => {
  fetchMoments()
})
</script>

<style scoped>
.moment-manage-page {
  padding: 0;
  max-width: 1000px;
}

.moment-item {
  padding: 16px;
  background: #f9fafb;
  border-radius: 8px;
  border: 1px solid #e5e7eb;
}

html.dark .moment-item {
  background: #1a1a1a;
  border-color: #2a2a2a;
}

.moment-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.moment-user {
  font-size: 14px;
  font-weight: 600;
  color: #1a202c;
}

html.dark .moment-user {
  color: #e5e5e5;
}

.moment-time {
  font-size: 12px;
  color: #94a3b8;
  margin-top: 2px;
}

.moment-content {
  font-size: 14px;
  line-height: 1.6;
  color: #334155;
  margin-bottom: 12px;
  white-space: pre-wrap;
  word-break: break-word;
}

html.dark .moment-content {
  color: #cbd5e1;
}

.moment-images {
  margin-bottom: 12px;
}

.moment-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
  border-top: 1px solid #e5e7eb;
}

html.dark .moment-footer {
  border-top-color: #374151;
}

.stat-item {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  color: #64748b;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .moment-item {
    padding: 12px;
  }
  
  .moment-header {
    flex-wrap: wrap;
    gap: 8px;
  }
  
  .moment-content {
    font-size: 13px;
  }
  
  .moment-footer {
    flex-wrap: wrap;
    gap: 8px;
  }
}
</style>

