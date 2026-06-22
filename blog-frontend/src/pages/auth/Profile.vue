<!--
 * @ProjectName: go-vue3-blog
 * @FileName: Profile.vue
 * @CreateTime: 2026-02-02 11:43:17
 * @SystemUser: Administrator
 * @Author: 無以菱
 * @Contact: huangjing510@126.com
 * @Description: 个人资料页面组件，提供用户信息查看和编辑功能
 -->
<template>
  <div class="profile-page">
    <n-grid :cols="1" :x-gap="24" :y-gap="24">
      <!-- 基本信息 -->
      <n-gi>
        <n-card title="基本信息">
          <n-form ref="profileFormRef" :model="profileForm" label-width="80px">
            <n-form-item label="头像">
              <avatar-upload
                v-model="profileForm.avatar"
                :size="100"
                :default-text="authStore.user?.nickname?.charAt(0) || 'U'"
                @success="handleAvatarSuccess"
              />
            </n-form-item>

            <n-form-item label="用户名">
              <n-input v-model:value="authStore.user!.username" disabled />
            </n-form-item>

            <n-form-item label="昵称" path="nickname">
              <n-input v-model:value="profileForm.nickname" placeholder="请输入昵称" />
            </n-form-item>

            <n-form-item label="邮箱">
              <n-input v-model:value="authStore.user!.email" disabled>
                <template #suffix>
                  <n-button text type="primary" @click="showEmailModal = true"> 修改 </n-button>
                </template>
              </n-input>
              <template #feedback>
                <span v-if="emailChangeInfo" style="font-size: 12px; color: #999">
                  今年已修改{{ emailChangeInfo.change_count }}次，还可修改{{
                    emailChangeInfo.remaining_times
                  }}次
                </span>
              </template>
            </n-form-item>

            <n-form-item label="个人简介" path="bio">
              <n-input
                v-model:value="profileForm.bio"
                type="textarea"
                :rows="4"
                placeholder="介绍一下自己吧"
              />
            </n-form-item>

            <n-form-item>
              <n-button type="primary" :loading="updating" @click="handleUpdateProfile">
                保存修改
              </n-button>
            </n-form-item>
          </n-form>
        </n-card>
      </n-gi>
    </n-grid>

    <!-- 修改邮箱弹窗 -->
    <n-modal
      v-model:show="showEmailModal"
      preset="dialog"
      title="修改邮箱"
      positive-text="确认修改"
      negative-text="取消"
      :positive-button-props="{ loading: emailUpdating }"
      @positive-click="handleUpdateEmail"
    >
      <n-alert v-if="!emailChangeInfo?.can_change" type="warning" style="margin-bottom: 16px">
        您今年的邮箱修改次数已达到上限（2次），请明年再试
      </n-alert>
      <n-form v-else>
        <n-form-item label="当前邮箱">
          <n-input :value="authStore.user?.email" disabled />
        </n-form-item>
        <n-form-item label="新邮箱" required>
          <n-input v-model:value="newEmail" placeholder="请输入新邮箱地址" />
        </n-form-item>
        <n-alert type="info" style="margin-top: 12px">
          <template #icon>
            <span>💡</span>
          </template>
          一年内只能修改两次邮箱，请谨慎操作
        </n-alert>
      </n-form>
    </n-modal>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import type { FormInst } from 'naive-ui'
import { useAuthStore } from '@/stores'
import { updateProfile, getEmailChangeInfo, updateEmail } from '@/api/auth'
import type { ProfileForm } from '@/types/auth'
import AvatarUpload from '@/components/AvatarUpload.vue'

const message = useMessage()
const authStore = useAuthStore()

const profileFormRef = ref<FormInst | null>(null)
const updating = ref(false)
const showEmailModal = ref(false)
const newEmail = ref('')
const emailUpdating = ref(false)
const emailChangeInfo = ref<{
  change_count: number
  remaining_times: number
  can_change: boolean
} | null>(null)

const profileForm = reactive<ProfileForm>({
  nickname: '',
  avatar: '',
  bio: ''
})

onMounted(async () => {
  if (authStore.user) {
    profileForm.nickname = authStore.user.nickname
    profileForm.avatar = authStore.user.avatar
    profileForm.bio = authStore.user.bio
  }

  // 获取邮箱修改信息
  await fetchEmailChangeInfo()
})

async function fetchEmailChangeInfo() {
  try {
    const res = await getEmailChangeInfo()
    if (res.data) {
      emailChangeInfo.value = res.data
    }
  } catch (error) {
    console.error('获取邮箱修改信息失败:', error)
  }
}

async function handleAvatarSuccess(url: string) {
  profileForm.avatar = url

  // 上传头像后自动保存到用户资料
  await updateProfile({ avatar: url })
  await authStore.fetchUserInfo()
  message.success('头像更新成功')
}

async function handleUpdateProfile() {
  await profileFormRef.value?.validate()
  updating.value = true

  await updateProfile(profileForm)
  await authStore.fetchUserInfo()

  // 重新加载表单数据（确保显示最新的头像）
  if (authStore.user) {
    profileForm.nickname = authStore.user.nickname
    profileForm.avatar = authStore.user.avatar
    profileForm.bio = authStore.user.bio
  }

  updating.value = false
  message.success('个人信息更新成功')
}

async function handleUpdateEmail() {
  if (!emailChangeInfo.value?.can_change) {
    message.error('今年的邮箱修改次数已达到上限')
    return false
  }

  if (!newEmail.value) {
    message.error('请输入新邮箱')
    return false
  }

  // 验证邮箱格式
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  if (!emailRegex.test(newEmail.value)) {
    message.error('邮箱格式不正确')
    return false
  }

  try {
    emailUpdating.value = true
    await updateEmail({ new_email: newEmail.value })
    await authStore.fetchUserInfo()
    await fetchEmailChangeInfo()

    message.success('邮箱修改成功')
    showEmailModal.value = false
    newEmail.value = ''
    return true
  } catch (error: any) {
    message.error(error.message || '邮箱修改失败')
    return false
  } finally {
    emailUpdating.value = false
  }
}
</script>

<style scoped>
.profile-page {
  max-width: 800px;
  margin: 0 auto;
}
</style>
