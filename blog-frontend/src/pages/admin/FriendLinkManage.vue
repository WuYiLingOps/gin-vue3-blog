<template>
  <div class="friendlink-manage-page">
    <div class="header">
      <h1>友链管理</h1>
      <n-space v-if="activeTab === 'links'">
        <n-button :size="isMobile ? 'small' : 'medium'" @click="showMyInfoModal = true">
          <template #icon>
            <n-icon :component="PersonOutline" />
          </template>
          <span v-if="!isMobile">我的友链信息</span>
          <span v-else>我的信息</span>
        </n-button>
        <n-button type="primary" :size="isMobile ? 'small' : 'medium'" @click="handleCreate">
          <template #icon>
            <n-icon :component="AddOutline" />
          </template>
          <span v-if="!isMobile">新建友链</span>
          <span v-else>新建</span>
        </n-button>
      </n-space>
      <n-button v-else type="primary" :size="isMobile ? 'small' : 'medium'" @click="handleCreateCategory">
        <template #icon>
          <n-icon :component="AddOutline" />
        </template>
        <span v-if="!isMobile">新建分类</span>
        <span v-else>新建</span>
      </n-button>
    </div>

    <n-tabs v-model:value="activeTab" type="line" animated>
      <n-tab-pane name="links" tab="友链管理">
        <n-data-table 
          :columns="columns" 
          :data="friendLinks" 
          :loading="loading"
          :scroll-x="isMobile ? 1200 : undefined"
          :single-line="false"
          :pagination="pagination"
          @update:page="handlePageChange"
          @update:page-size="handlePageSizeChange"
        />
      </n-tab-pane>
      <n-tab-pane name="categories" tab="友链分类">
        <n-data-table 
          :columns="categoryColumns" 
          :data="categories" 
          :loading="categoryLoading"
          :scroll-x="isMobile ? 800 : undefined"
          :single-line="false"
        />
      </n-tab-pane>
    </n-tabs>

    <!-- 创建/编辑对话框 -->
    <n-modal 
      v-model:show="showModal" 
      preset="card" 
      :title="editingId ? '编辑友链' : '新建友链'" 
      :style="{ width: isMobile ? '95%' : '700px', maxWidth: isMobile ? '95vw' : '700px' }"
    >
      <n-form ref="formRef" :model="formData" :rules="rules">
        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <n-form-item label="网站名称" path="name">
              <n-input v-model:value="formData.name" placeholder="例如：無以菱" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="网址" path="url">
              <n-input v-model:value="formData.url" placeholder="https://example.com" />
            </n-form-item>
          </n-gi>
        </n-grid>

        <n-form-item label="分类" path="category_id">
          <n-select
            v-model:value="formData.category_id"
            :options="categoryOptions"
            placeholder="请选择分类"
            :loading="categories.length === 0"
          />
        </n-form-item>

        <n-form-item label="网站图标">
          <n-input v-model:value="formData.icon" placeholder="https://example.com/icon.ico" />
        </n-form-item>

        <n-form-item label="网站描述">
          <n-input
            v-model:value="formData.description"
            type="textarea"
            :rows="2"
            placeholder="网站描述"
          />
        </n-form-item>

        <n-form-item label="网站截图">
          <n-input v-model:value="formData.screenshot" placeholder="https://example.com/screenshot.jpg" />
        </n-form-item>

        <n-form-item label="Atom/RSS 地址（可选）">
          <n-input v-model:value="formData.atom_url" placeholder="https://example.com/atom.xml" />
        </n-form-item>

        <n-grid :cols="2" :x-gap="16">
          <n-gi>
            <n-form-item label="排序">
              <n-input-number v-model:value="formData.sort_order" :min="0" style="width: 100%" />
            </n-form-item>
          </n-gi>
          <n-gi>
            <n-form-item label="状态">
              <n-radio-group v-model:value="formData.status">
                <n-radio :value="1">启用</n-radio>
                <n-radio :value="0">禁用</n-radio>
              </n-radio-group>
            </n-form-item>
          </n-gi>
        </n-grid>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">
            保存
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 我的友链信息对话框 -->
    <n-modal 
      v-model:show="showMyInfoModal" 
      preset="card" 
      title="我的友链信息" 
      :style="{ width: isMobile ? '95%' : '700px', maxWidth: isMobile ? '95vw' : '700px' }"
    >
      <n-form ref="myInfoFormRef" :model="myInfoFormData" :rules="myInfoRules">
        <n-form-item label="名称" path="name">
          <n-input v-model:value="myInfoFormData.name" placeholder="例如：無以菱" />
        </n-form-item>

        <n-form-item label="描述">
          <n-input
            v-model:value="myInfoFormData.desc"
            type="textarea"
            :rows="2"
            placeholder="网站描述"
          />
        </n-form-item>

        <n-form-item label="地址" path="url">
          <n-input v-model:value="myInfoFormData.url" placeholder="https://example.com" />
        </n-form-item>

        <n-form-item label="头像">
          <n-input v-model:value="myInfoFormData.avatar" placeholder="https://example.com/avatar.png" />
        </n-form-item>

        <n-form-item label="站点图片">
          <n-input v-model:value="myInfoFormData.screenshot" placeholder="https://example.com/screenshot.jpg" />
        </n-form-item>

        <n-form-item label="RSS/Atom 订阅地址">
          <n-input v-model:value="myInfoFormData.rss" placeholder="https://example.com/rss.xml" />
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showMyInfoModal = false">取消</n-button>
          <n-button type="primary" :loading="myInfoSubmitting" @click="handleSubmitMyInfo">
            保存
          </n-button>
        </n-space>
      </template>
    </n-modal>

    <!-- 创建/编辑分类对话框 -->
    <n-modal 
      v-model:show="showCategoryModal" 
      preset="card" 
      :title="editingCategoryId ? '编辑分类' : '新建分类'" 
      :style="{ width: isMobile ? '95%' : '600px', maxWidth: isMobile ? '95vw' : '600px' }"
    >
      <n-form ref="categoryFormRef" :model="categoryFormData" :rules="categoryRules">
        <n-form-item label="分类名称" path="name">
          <n-input v-model:value="categoryFormData.name" placeholder="例如：推荐" />
        </n-form-item>

        <n-form-item label="分类描述">
          <n-input
            v-model:value="categoryFormData.description"
            type="textarea"
            :rows="2"
            placeholder="分类描述，例如：都是大佬,推荐关注"
          />
        </n-form-item>

        <n-form-item label="排序">
          <n-input-number v-model:value="categoryFormData.sort_order" :min="0" style="width: 100%" />
          <template #feedback>
            <span style="color: #909399; font-size: 12px;">数字越大越靠前</span>
          </template>
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showCategoryModal = false">取消</n-button>
          <n-button type="primary" :loading="categorySubmitting" @click="handleSubmitCategory">
            保存
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, onUnmounted, h } from 'vue'
import { useMessage, useDialog, NButton, NTag, NSpace, NImage } from 'naive-ui'
import type { DataTableColumns } from 'naive-ui'
import { AddOutline, PersonOutline } from '@vicons/ionicons5'
import { 
  getFriendLinksAdmin, 
  createFriendLink, 
  updateFriendLink, 
  deleteFriendLink, 
  getFriendLinkCategoriesAdmin,
  createFriendLinkCategory,
  updateFriendLinkCategory,
  deleteFriendLinkCategory
} from '@/api/friendlink'
import type { FriendLink, FriendLinkForm, FriendLinkCategory, FriendLinkCategoryForm } from '@/api/friendlink'
import { getFriendLinkInfo, updateFriendLinkInfo, type FriendLinkInfo } from '@/api/setting'

// 更宽松且可靠的 URL 校验：使用浏览器 URL 解析，限定 http/https
const validateUrl = (_rule: unknown, value: string) => {
  if (!value) return Promise.resolve()
  try {
    const urlObj = new URL(value)
    if (urlObj.protocol === 'http:' || urlObj.protocol === 'https:') {
      return Promise.resolve()
    }
  } catch (e) {
    // ignore
  }
  return Promise.reject('请输入有效的网址')
}

const message = useMessage()
const dialog = useDialog()

const activeTab = ref<'links' | 'categories'>('links')
const loading = ref(false)
const submitting = ref(false)
const showModal = ref(false)
const friendLinks = ref<FriendLink[]>([])
const categories = ref<FriendLinkCategory[]>([])
const editingId = ref<number | null>(null)
const isMobile = ref(false)
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

// 分类管理相关
const categoryLoading = ref(false)
const categorySubmitting = ref(false)
const showCategoryModal = ref(false)
const editingCategoryId = ref<number | null>(null)
const categoryFormData = reactive<FriendLinkCategoryForm>({
  name: '',
  description: '',
  sort_order: 0
})

// 我的友链信息相关
const showMyInfoModal = ref(false)
const myInfoSubmitting = ref(false)
const myInfoFormData = reactive<FriendLinkInfo>({
  name: '',
  desc: '',
  url: '',
  avatar: '',
  screenshot: '',
  rss: ''
})

const myInfoRules = {
  name: [{ required: true, message: '请输入名称', trigger: 'blur' }],
  url: [
    { required: true, message: '请输入地址', trigger: 'blur' },
    { validator: validateUrl, trigger: 'blur' }
  ]
}

const pagination = reactive({
  page: currentPage,
  pageSize: pageSize,
  itemCount: total,
  showSizePicker: true,
  pageSizes: [10, 20, 50],
  showQuickJumper: true
})

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

const formData = reactive<FriendLinkForm>({
  name: '',
  url: '',
  icon: '',
  description: '',
  screenshot: '',
  atom_url: '',
  category_id: 0,
  sort_order: 0,
  status: 1
})

const rules = {
  name: [{ required: true, message: '请输入网站名称', trigger: 'blur' }],
  url: [
    { required: true, message: '请输入网址', trigger: 'blur' },
    { validator: validateUrl, trigger: 'blur' }
  ],
  category_id: [{ required: true, message: '请选择分类', trigger: 'change', type: 'number' as const }]
}

const categoryRules = {
  name: [{ required: true, message: '请输入分类名称', trigger: 'blur' }]
}

const columns: DataTableColumns<FriendLink> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60
  },
  {
    title: '网站名称',
    key: 'name',
    width: 150,
    render: row => h('div', { style: 'display: flex; align-items: center; gap: 8px;' }, [
      row.icon ? h(NImage, { src: row.icon, width: 24, height: 24, objectFit: 'cover', style: 'border-radius: 4px;' }) : null,
      h('span', row.name)
    ])
  },
  {
    title: '网址',
    key: 'url',
    width: 200,
    render: row => h('a', { href: row.url, target: '_blank', rel: 'noopener noreferrer', style: 'color: #18a058; text-decoration: none;' }, row.url)
  },
  {
    title: '分类',
    key: 'category',
    width: 120,
    render: row => row.category ? row.category.name : '-'
  },
  {
    title: '描述',
    key: 'description',
    width: 200,
    ellipsis: { tooltip: true }
  },
  {
    title: '排序',
    key: 'sort_order',
    width: 80
  },
  {
    title: '状态',
    key: 'status',
    width: 80,
    render: row => h(NTag, { type: row.status === 1 ? 'success' : 'default' }, { default: () => row.status === 1 ? '启用' : '禁用' })
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
    render: row =>
      h(NSpace, null, {
        default: () => [
          h(
            NButton,
            { size: 'small', onClick: () => handleEdit(row) },
            { default: () => '编辑' }
          ),
          h(
            NButton,
            { size: 'small', type: 'error', onClick: () => handleDelete(row.id) },
            { default: () => '删除' }
          )
        ]
      })
  }
]

// 分类表格列
const categoryColumns: DataTableColumns<FriendLinkCategory> = [
  { 
    title: 'ID', 
    key: 'id', 
    width: 60
  },
  {
    title: '分类名称',
    key: 'name',
    width: 150
  },
  {
    title: '描述',
    key: 'description',
    width: 250,
    ellipsis: { tooltip: true }
  },
  {
    title: '排序',
    key: 'sort_order',
    width: 100
  },
  {
    title: '操作',
    key: 'actions',
    width: 150,
    fixed: 'right',
    render: row =>
      h(NSpace, null, {
        default: () => [
          h(
            NButton,
            { size: 'small', onClick: () => handleEditCategory(row) },
            { default: () => '编辑' }
          ),
          h(
            NButton,
            { size: 'small', type: 'error', onClick: () => handleDeleteCategory(row.id) },
            { default: () => '删除' }
          )
        ]
      })
  }
]

// 分类选项
const categoryOptions = computed(() => {
  return categories.value.map(cat => ({
    label: cat.name,
    value: cat.id
  }))
})

// 获取分类列表
async function fetchCategories() {
  try {
    const res = await getFriendLinkCategoriesAdmin()
    if (res && res.data) {
      categories.value = Array.isArray(res.data) ? res.data : []
    }
  } catch (e: any) {
    console.error('获取分类列表失败:', e)
  }
}

onMounted(() => {
  checkMobile()
  window.addEventListener('resize', checkMobile)
  fetchCategories()
  fetchFriendLinks()
  fetchMyInfo()
})

// 分类管理相关函数
function resetCategoryForm() {
  editingCategoryId.value = null
  categoryFormData.name = ''
  categoryFormData.description = ''
  categoryFormData.sort_order = 0
}

function handleCreateCategory() {
  resetCategoryForm()
  showCategoryModal.value = true
}

function handleEditCategory(category: FriendLinkCategory) {
  editingCategoryId.value = category.id
  categoryFormData.name = category.name
  categoryFormData.description = category.description || ''
  categoryFormData.sort_order = category.sort_order
  showCategoryModal.value = true
}

async function handleSubmitCategory() {
  try {
    await categoryFormRef.value?.validate()
    categorySubmitting.value = true

    if (editingCategoryId.value) {
      await updateFriendLinkCategory(editingCategoryId.value, categoryFormData)
      message.success('分类更新成功')
    } else {
      await createFriendLinkCategory(categoryFormData)
      message.success('分类创建成功')
    }

    showCategoryModal.value = false
    fetchCategories()
    // 如果友链表单正在使用分类，刷新分类选项
    if (showModal.value) {
      // 如果编辑的是当前选中的分类，保持选中状态
      if (editingCategoryId.value && formData.category_id === editingCategoryId.value) {
        // 保持选中
      }
    }
  } catch (error: any) {
    if (error.message && !error.message.includes('验证')) {
      message.error(error.message || '操作失败')
    }
  } finally {
    categorySubmitting.value = false
  }
}

function handleDeleteCategory(id: number) {
  const category = categories.value.find(c => c.id === id)
  const categoryName = category?.name || '该分类'
  
  dialog.warning({
    title: '确认删除',
    content: `确定要删除分类"${categoryName}"吗？删除后该分类下的友链需要重新分配分类。`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteFriendLinkCategory(id)
        message.success('删除成功')
        fetchCategories()
        // 如果删除的是当前选中的分类，重置友链表单的分类选择
        if (formData.category_id === id && categories.value.length > 0) {
          formData.category_id = categories.value[0].id
        }
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}

async function fetchMyInfo() {
  try {
    const res = await getFriendLinkInfo()
    if (res.data) {
      myInfoFormData.name = res.data.name || ''
      myInfoFormData.desc = res.data.desc || ''
      myInfoFormData.url = res.data.url || ''
      myInfoFormData.avatar = res.data.avatar || ''
      myInfoFormData.screenshot = res.data.screenshot || ''
      myInfoFormData.rss = res.data.rss || ''
    }
  } catch (error: any) {
    console.error('获取我的友链信息失败:', error)
  }
}

async function handleSubmitMyInfo() {
  try {
    await myInfoFormRef.value?.validate()
    myInfoSubmitting.value = true
    await updateFriendLinkInfo(myInfoFormData)
    message.success('更新成功')
    showMyInfoModal.value = false
  } catch (error: any) {
    if (error.message && !error.message.includes('验证')) {
      message.error(error.message || '操作失败')
    }
  } finally {
    myInfoSubmitting.value = false
  }
}

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})

function handlePageChange(page: number) {
  currentPage.value = page
  fetchFriendLinks()
}

function handlePageSizeChange(size: number) {
  pageSize.value = size
  currentPage.value = 1
  fetchFriendLinks()
}

async function fetchFriendLinks() {
  try {
    loading.value = true
    const res = await getFriendLinksAdmin(currentPage.value, pageSize.value)
    if (res.data) {
      friendLinks.value = res.data.list
      total.value = res.data.total
    }
  } catch (error: any) {
    message.error(error.message || '获取友链列表失败')
  } finally {
    loading.value = false
  }
}

function resetForm() {
  editingId.value = null
  formData.name = ''
  formData.url = ''
  formData.icon = ''
  formData.description = ''
  formData.screenshot = ''
  formData.atom_url = ''
  formData.category_id = categories.value.length > 0 ? categories.value[0].id : 0
  formData.sort_order = 0
  formData.status = 1
}

function handleCreate() {
  resetForm()
  showModal.value = true
}

function handleEdit(friendLink: FriendLink) {
  editingId.value = friendLink.id
  formData.name = friendLink.name
  formData.url = friendLink.url
  formData.icon = friendLink.icon || ''
  formData.description = friendLink.description || ''
  formData.screenshot = friendLink.screenshot || ''
  formData.atom_url = friendLink.atom_url || ''
  formData.category_id = friendLink.category_id
  formData.sort_order = friendLink.sort_order
  formData.status = friendLink.status
  showModal.value = true
}

async function handleSubmit() {
  try {
    await formRef.value?.validate()
    submitting.value = true
    if (editingId.value) {
      // 确保传递所有字段，包括 category_id
      const updateData: Partial<FriendLinkForm> = {
        name: formData.name,
        url: formData.url,
        icon: formData.icon || '',
        description: formData.description || '',
        screenshot: formData.screenshot || '',
        atom_url: formData.atom_url || '',
        category_id: formData.category_id,
        sort_order: formData.sort_order,
        status: formData.status
      }
      await updateFriendLink(editingId.value, updateData)
      message.success('更新成功')
    } else {
      await createFriendLink(formData)
      message.success('创建成功')
    }
    showModal.value = false
    resetForm()
    fetchFriendLinks()
  } catch (error: any) {
    if (error.message && !error.message.includes('验证')) {
      message.error(error.message || '操作失败')
    }
  } finally {
    submitting.value = false
  }
}

function handleDelete(id: number) {
  const friendLink = friendLinks.value.find(f => f.id === id)
  const friendLinkName = friendLink?.name || '该友链'
  
  dialog.warning({
    title: '确认删除',
    content: `确定要删除友链"${friendLinkName}"吗？`,
    positiveText: '确定',
    negativeText: '取消',
    onPositiveClick: async () => {
      try {
        await deleteFriendLink(id)
        message.success('删除成功')
        fetchFriendLinks()
      } catch (error: any) {
        message.error(error.message || '删除失败')
      }
    }
  })
}

const formRef = ref()
const myInfoFormRef = ref()
const categoryFormRef = ref()
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  gap: 12px;
}

.header h1 {
  margin: 0;
  font-size: 24px;
}

/* 移动端样式 */
@media (max-width: 768px) {
  .header h1 {
    font-size: 20px;
  }
  
  .friendlink-manage-page :deep(.n-data-table) {
    font-size: 13px;
  }
}
</style>

