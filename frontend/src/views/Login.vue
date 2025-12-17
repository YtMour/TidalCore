<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useUserStore } from '@/store/user'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, Timer, Star } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const form = reactive({
  username: '',
  password: ''
})
const loading = ref(false)
const error = ref('')
const mounted = ref(false)

const rules = reactive<FormRules>({
  username: [{ required: true, message: '请输入用户名', trigger: 'blur' }],
  password: [{ required: true, message: '请输入密码', trigger: 'blur' }]
})

onMounted(() => {
  setTimeout(() => { mounted.value = true }, 100)
  if (route.query.expired === '1') {
    error.value = '登录已过期，请重新登录'
  }
})

function isValidRedirect(url: string): boolean {
  if (!url) return false
  return url.startsWith('/') && !url.startsWith('//')
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    loading.value = true
    error.value = ''
    try {
      await userStore.login(form)
      const redirect = route.query.redirect as string
      router.push(isValidRedirect(redirect) ? redirect : '/')
    } catch (e) {
      error.value = e instanceof Error ? e.message : '登录失败'
    } finally {
      loading.value = false
    }
  })
}
</script>

<template>
  <MainLayout>
    <div class="login-container">
      <div class="login-bg">
        <div class="login-blob login-blob-1"></div>
        <div class="login-blob login-blob-2"></div>
      </div>

      <div class="login-card-wrapper" :class="{ mounted }">
        <el-card class="login-card" shadow="always">
          <div class="login-header">
            <div class="login-icon">
              <el-icon :size="36"><User /></el-icon>
            </div>
            <h1 class="login-title">欢迎回来</h1>
            <p class="login-subtitle">登录以继续你的训练之旅</p>
          </div>

          <el-form ref="formRef" :model="form" :rules="rules" label-position="top" size="large" @submit.prevent="handleSubmit">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="form.username" placeholder="请输入用户名" :prefix-icon="User" autocomplete="username" />
            </el-form-item>

            <el-form-item label="密码" prop="password">
              <el-input v-model="form.password" type="password" placeholder="请输入密码" :prefix-icon="Lock" show-password autocomplete="current-password" />
            </el-form-item>

            <el-alert v-if="error" :title="error" type="error" show-icon :closable="false" class="login-error" />

            <el-form-item>
              <el-button type="primary" native-type="submit" :loading="loading" class="login-btn">
                {{ loading ? '登录中...' : '登录' }}
              </el-button>
            </el-form-item>
          </el-form>

          <div class="login-features">
            <el-tag effect="plain" round><el-icon><Lock /></el-icon> 安全加密</el-tag>
            <el-tag effect="plain" round><el-icon><Timer /></el-icon> 快速登录</el-tag>
            <el-tag effect="plain" round><el-icon><Star /></el-icon> 个性化记录</el-tag>
          </div>

          <el-divider />
          <div class="login-footer">
            <span>还没有账号?</span>
            <RouterLink to="/register"><el-link type="primary">立即注册</el-link></RouterLink>
          </div>
        </el-card>
      </div>
    </div>
  </MainLayout>
</template>

<style scoped>
.login-container {
  min-height: 70vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 16px;
  position: relative;
}

.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.login-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
}

.login-blob-1 {
  top: 25%;
  left: 25%;
  width: 256px;
  height: 256px;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(168, 85, 247, 0.1));
  animation: float 6s ease-in-out infinite;
}

.login-blob-2 {
  bottom: 25%;
  right: 25%;
  width: 192px;
  height: 192px;
  background: linear-gradient(135deg, rgba(236, 72, 153, 0.15), rgba(244, 63, 94, 0.1));
  animation: float 8s ease-in-out infinite;
  animation-delay: -3s;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

.login-card-wrapper {
  width: 100%;
  max-width: 420px;
  opacity: 0;
  transform: translateY(32px);
  transition: all 0.8s ease;
}

.login-card-wrapper.mounted {
  opacity: 1;
  transform: translateY(0);
}

.login-card {
  background: rgba(30, 30, 46, 0.9);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
}

.login-card :deep(.el-card__body) {
  padding: 40px;
}

.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-icon {
  width: 72px;
  height: 72px;
  margin: 0 auto 20px;
  border-radius: 16px;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 40px rgba(139, 92, 246, 0.3);
  color: white;
}

.login-title {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
}

.login-subtitle {
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
  margin: 0;
}

.login-error {
  margin-bottom: 20px;
}

.login-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  border-radius: 12px;
  background: linear-gradient(135deg, #8b5cf6, #ec4899) !important;
  border: none !important;
}

.login-btn:hover {
  opacity: 0.9;
}

.login-features {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 8px;
  margin-top: 24px;
}

.login-features .el-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.1);
}

.login-footer {
  text-align: center;
  color: rgba(255, 255, 255, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

:deep(.el-form-item__label) {
  color: rgba(255, 255, 255, 0.8);
}

:deep(.el-input__wrapper) {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 10px;
  box-shadow: none !important;
}

:deep(.el-input__wrapper:hover) {
  border-color: rgba(139, 92, 246, 0.5);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: #8b5cf6;
}

:deep(.el-input__inner) {
  color: #fff;
}

:deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.3);
}

:deep(.el-divider) {
  border-color: rgba(255, 255, 255, 0.1);
  margin: 24px 0;
}
</style>
