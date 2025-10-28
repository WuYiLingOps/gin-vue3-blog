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
              <n-input v-model:value="authStore.user!.email" disabled />
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
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useMessage } from 'naive-ui'
import type { FormInst } from 'naive-ui'
import { useAuthStore } from '@/stores'
import { updateProfile } from '@/api/auth'
import type { ProfileForm } from '@/types/auth'
import AvatarUpload from '@/components/AvatarUpload.vue'

const message = useMessage()
const authStore = useAuthStore()

const profileFormRef = ref<FormInst | null>(null)
const updating = ref(false)

const profileForm = reactive<ProfileForm>({
  nickname: '',
  avatar: '',
  bio: ''
})

onMounted(() => {
  if (authStore.user) {
    profileForm.nickname = authStore.user.nickname
    profileForm.avatar = authStore.user.avatar
    profileForm.bio = authStore.user.bio
  }
})

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
</script>

<style scoped>
.profile-page {
  max-width: 800px;
  margin: 0 auto;
}
</style>

