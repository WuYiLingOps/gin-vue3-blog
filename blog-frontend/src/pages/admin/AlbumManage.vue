<!--
 * @ProjectName: go-vue3-blog
 * @FileName: AlbumManage.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 相册管理页面组件，提供相册图片的增删改查功能
 -->
<template>
  <div class="album-manage-page">
    <div class="header">
      <h1>我的相册</h1>
      <n-button type="primary" :size="isMobile ? 'small' : 'medium'" @click="handleCreate">
        <template #icon>
          <n-icon :component="AddOutline" />
        </template>
        <span v-if="!isMobile">新增照片</span>
        <span v-else>新增</span>
      </n-button>
    </div>

    <n-data-table 
      :columns="columns" 
      :data="albums" 
      :loading="loading"
      :scroll-x="isMobile ? 1200 : undefined"
      :single-line="false"
      :pagination="pagination"
      @update:page="handlePageChange"
      @update:page-size="handlePageSizeChange"
    />

    <!-- 创建/编辑对话框 -->
    <n-modal 
      v-model:show="showModal" 
      preset="card" 
      :title="editingId ? '编辑照片' : '新增照片'" 
      :style="{ width: isMobile ? '95%' : '700px', maxWidth: isMobile ? '95vw' : '700px' }"
    >
      <n-form ref="formRef" :model="formData" :rules="rules">
        <n-form-item label="照片" path="image_url">
          <ImageUpload
            v-model="formData.image_url"
            :width="400"
            :height="250"
          />
        </n-form-item>

        <n-form-item label="标题">
          <n-input v-model:value="formData.title" placeholder="照片标题（可选）" />
        </n-form-item>

        <n-form-item label="描述">
          <n-input
            v-model:value="formData.description"
            type="textarea"
            :rows="3"
            placeholder="照片描述（可选）"
          />
        </n-form-item>

        <n-form-item label="排序">
          <n-input-number v-model:value="formData.sort_order" :min="0" style="width: 100%" />
          <template #feedback>
            <n-text depth="3" style="font-size: 12px">
              数字越大越靠前，默认为 0
            </n-text>
          </template>
        </n-form-item>
      </n-form>

      <template #footer>
        <n-space justify="end">
          <n-button @click="showModal = false">取消</n-button>
          <n-button type="primary" :loading="submitting" @click="handleSubmit">
            {{ editingId ? '更新' : '创建' }}
          </n-button>
        </n-space>
      </template>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onUnmounted, h } from 'vue'
import { useMessage, NSpace, NButton } from 'naive-ui'
import type { FormInst, DataTableColumns } from 'naive-ui'
import { AddOutline } from '@vicons/ionicons5'
import { getAlbums, createAlbum, updateAlbum, deleteAlbum, type Album } from '@/api/album'
import ImageUpload from '@/components/ImageUpload.vue'

const message = useMessage()
const formRef = ref<FormInst | null>(null)
const loading = ref(false)
const submitting = ref(false)
const isMobile = ref(false)
const showModal = ref(false)
const editingId = ref<number | null>(null)

const albums = ref<Album[]>([])
const pagination = reactive({
  page: 1,
  pageSize: 10,
  pageCount: 1,
  itemCount: 0,
  showSizePicker: true,
  pageSizes: [10, 20, 50, 100],
  onChange: (page: number) => {
    pagination.page = page
    fetchAlbums()
  },
  onUpdatePageSize: (pageSize: number) => {
    pagination.pageSize = pageSize
    pagination.page = 1
    fetchAlbums()
  }
})

const formData = reactive({
  image_url: '',
  title: '',
  description: '',
  sort_order: 0
})

const rules = {
  image_url: {
    required: true,
    message: '请上传照片',
    trigger: ['blur', 'change']
  }
}

// 检测移动设备
function checkMobile() {
  isMobile.value = window.innerWidth <= 768
}

// 获取相册列表
async function fetchAlbums() {
  try {
    loading.value = true
    const res = await getAlbums(pagination.page, pagination.pageSize)
    if (res.data) {
      albums.value = res.data.list || []
      pagination.itemCount = res.data.total || 0
      pagination.pageCount = res.data.total_pages || Math.ceil((res.data.total || 0) / pagination.pageSize)
    }
  } catch (error: any) {
    message.error(error.message || '获取相册列表失败')
    console.error(error)
  } finally {
    loading.value = false
  }
}

// 表格列定义
const columns: DataTableColumns<Album> = [
  {
    title: 'ID',
    key: 'id',
    width: 80,
    fixed: 'left'
  },
  {
    title: '照片',
    key: 'image_url',
    width: 120,
    render(row) {
      return h('div', { style: { display: 'flex', alignItems: 'center', justifyContent: 'center' } }, [
        h('img', {
          src: row.image_url,
          style: {
            width: '80px',
            height: '80px',
            objectFit: 'cover',
            borderRadius: '4px',
            cursor: 'pointer'
          },
          onClick: () => {
            window.open(row.image_url, '_blank')
          }
        })
      ])
    }
  },
  {
    title: '标题',
    key: 'title',
    width: 150,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '描述',
    key: 'description',
    width: 200,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '排序',
    key: 'sort_order',
    width: 100
  },
  {
    title: '创建时间',
    key: 'created_at',
    width: 180,
    render(row) {
      return new Date(row.created_at).toLocaleString('zh-CN')
    }
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

// 处理分页变化
function handlePageChange(page: number) {
  pagination.page = page
  fetchAlbums()
}

function handlePageSizeChange(pageSize: number) {
  pagination.pageSize = pageSize
  pagination.page = 1
  fetchAlbums()
}

// 创建新照片
function handleCreate() {
  editingId.value = null
  formData.image_url = ''
  formData.title = ''
  formData.description = ''
  formData.sort_order = 0
  showModal.value = true
}

// 编辑照片
function handleEdit(album: Album) {
  editingId.value = album.id
  formData.image_url = album.image_url
  formData.title = album.title || ''
  formData.description = album.description || ''
  formData.sort_order = album.sort_order
  showModal.value = true
}

// 提交表单
async function handleSubmit() {
  try {
    // 表单验证
    await formRef.value?.validate()
    
    submitting.value = true
    
    if (editingId.value) {
      // 更新
      await updateAlbum(editingId.value, {
        image_url: formData.image_url,
        title: formData.title || undefined,
        description: formData.description || undefined,
        sort_order: formData.sort_order
      })
      message.success('更新成功')
    } else {
      // 创建
      await createAlbum({
        image_url: formData.image_url,
        title: formData.title || undefined,
        description: formData.description || undefined,
        sort_order: formData.sort_order
      })
      message.success('创建成功')
    }
    
    showModal.value = false
    fetchAlbums()
  } catch (error: any) {
    // 如果是表单验证错误，naive-ui 会自动处理
    if (error.message && !error.message.includes('required')) {
      message.error(error.message || '操作失败')
    }
    console.error('提交失败:', error)
  } finally {
    submitting.value = false
  }
}

// 删除照片
async function handleDelete(id: number) {
  try {
    await deleteAlbum(id)
    message.success('删除成功')
    fetchAlbums()
  } catch (error: any) {
    message.error(error.message || '删除失败')
    console.error(error)
  }
}

onMounted(() => {
  checkMobile()
  fetchAlbums()
  window.addEventListener('resize', checkMobile)
})

onUnmounted(() => {
  window.removeEventListener('resize', checkMobile)
})
</script>

<style scoped>
.album-manage-page {
  padding: 20px;
}

.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h1 {
  margin: 0;
  font-size: 24px;
  font-weight: 600;
}

@media (max-width: 768px) {
  .album-manage-page {
    padding: 12px;
  }

  .header {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }

  .header h1 {
    font-size: 20px;
  }
}
</style>
