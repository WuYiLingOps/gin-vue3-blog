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
  /* 使用主题的清新蓝绿色渐变背景 */
  background: linear-gradient(135deg, #e0f7fa 0%, #e1f5fe 25%, #f1f8f4 50%, #fef5e7 100%);
  position: relative;
}

/* 添加纹理层，与主题保持一致 */
.auth-layout::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-image: 
    radial-gradient(circle at 20% 50%, rgba(129, 212, 250, 0.08) 0%, transparent 50%),
    radial-gradient(circle at 80% 80%, rgba(102, 187, 106, 0.08) 0%, transparent 50%),
    radial-gradient(circle at 50% 10%, rgba(255, 213, 79, 0.05) 0%, transparent 50%);
  pointer-events: none;
}

.auth-container {
  width: 100%;
  max-width: 450px;
  padding: 24px;
  position: relative;
  z-index: 1;
}

.auth-content {
  /* 玻璃态效果 */
  background: rgba(255, 255, 255, 0.95);
  backdrop-filter: blur(20px) saturate(180%);
  border: 1px solid rgba(255, 255, 255, 0.5);
  border-radius: 16px;
  padding: 48px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.08);
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
  /* 使用主题的蓝绿色渐变 */
  background: linear-gradient(135deg, #0891b2 0%, #059669 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  transition: all 0.3s ease;
}

.auth-header h1:hover {
  transform: scale(1.05);
}

@media (max-width: 768px) {
  .auth-content {
    padding: 32px 24px;
  }
}

/* 深色模式支持 */
html.dark .auth-layout {
  background: linear-gradient(135deg, #0f172a 0%, #1e293b 50%, #0f172a 100%);
}

html.dark .auth-layout::before {
  background-image: 
    radial-gradient(circle at 20% 50%, rgba(59, 130, 246, 0.1) 0%, transparent 50%),
    radial-gradient(circle at 80% 80%, rgba(34, 197, 94, 0.08) 0%, transparent 50%);
}

html.dark .auth-content {
  background: rgba(30, 30, 46, 0.9);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.3);
}

html.dark .auth-header h1 {
  background: linear-gradient(135deg, #38bdf8 0%, #4ade80 100%);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}
</style>

