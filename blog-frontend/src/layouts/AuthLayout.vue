<template>
  <div class="auth-layout">
    <div class="auth-container">
      <div class="auth-content">
        <div class="auth-header">
          <h1 @click="router.push('/')">我的博客</h1>
        </div>
        <router-view />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'

const router = useRouter()
</script>

<style scoped>
.auth-layout {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
}

.auth-container {
  width: 100%;
  max-width: 450px;
  padding: 24px;
}

.auth-content {
  background: white;
  border-radius: 12px;
  padding: 48px;
  box-shadow: 0 10px 40px rgba(0, 0, 0, 0.1);
}

.auth-header {
  text-align: center;
  margin-bottom: 32px;
}

.auth-header h1 {
  font-size: 32px;
  font-weight: bold;
  margin: 0;
  cursor: pointer;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

@media (max-width: 768px) {
  .auth-content {
    padding: 32px 24px;
  }
}
</style>

