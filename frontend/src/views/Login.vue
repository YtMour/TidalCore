<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useUserStore } from '@/store/user'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock } from '@element-plus/icons-vue'

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
      <!-- 海洋主题背景 -->
      <div class="login-bg">
        <div class="login-glow login-glow-1"></div>
        <div class="login-glow login-glow-2"></div>
        <div class="login-ripple login-ripple-1"></div>
        <div class="login-ripple login-ripple-2"></div>
      </div>

      <div class="login-card-wrapper" :class="{ mounted }">
        <el-card class="login-card" shadow="never">
          <!-- 装饰波浪 -->
          <div class="card-wave-top"></div>
          <div class="card-decoration"></div>

          <div class="login-header">
            <div class="login-icon">
              <svg viewBox="0 0 32 32" fill="none" width="32" height="32" class="wave-svg">
                <path d="M4 16C4 16 8 10 14 16C20 22 26 10 28 16" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
              </svg>
            </div>
            <h1 class="login-title">欢迎回来</h1>
            <p class="login-subtitle">登录以继续你的潮汐训练之旅</p>
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

          <div class="login-trust">
            <div class="trust-item">
              <span class="trust-icon">
                <svg viewBox="0 0 16 16" fill="none" width="14" height="14">
                  <path d="M8 2L10 6L14 6.5L11 9.5L12 14L8 12L4 14L5 9.5L2 6.5L6 6L8 2Z" fill="currentColor"/>
                </svg>
              </span>
              <span>安全加密</span>
            </div>
            <span class="trust-divider"></span>
            <div class="trust-item">
              <span class="trust-icon">
                <svg viewBox="0 0 16 16" fill="none" width="14" height="14">
                  <path d="M2 8C2 8 4 5 6 8C8 11 10 5 12 8C14 11 16 8 16 8" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
              </span>
              <span>潮汐同步</span>
            </div>
            <span class="trust-divider"></span>
            <div class="trust-item">
              <span class="trust-icon">
                <svg viewBox="0 0 16 16" fill="none" width="14" height="14">
                  <circle cx="8" cy="8" r="6" stroke="currentColor" stroke-width="1.5"/>
                  <path d="M8 5V8L10 10" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
                </svg>
              </span>
              <span>实时记录</span>
            </div>
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

/* ===== 海洋背景效果 ===== */
.login-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.login-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
}

.login-glow-1 {
  top: 20%;
  left: 20%;
  width: 300px;
  height: 300px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.2), transparent 70%);
  animation: float 8s ease-in-out infinite;
}

.login-glow-2 {
  bottom: 20%;
  right: 20%;
  width: 250px;
  height: 250px;
  background: radial-gradient(circle, rgba(34, 211, 238, 0.15), transparent 70%);
  animation: float 10s ease-in-out infinite;
  animation-delay: -3s;
}

.login-ripple {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  border: 1px solid rgba(56, 189, 248, 0.15);
  animation: ripple-expand 8s ease-out infinite;
}

.login-ripple-1 {
  width: 400px;
  height: 400px;
}

.login-ripple-2 {
  width: 600px;
  height: 600px;
  animation-delay: 3s;
}

@keyframes float {
  0%, 100% { transform: translateY(0) rotate(0deg); }
  50% { transform: translateY(-25px) rotate(3deg); }
}

@keyframes ripple-expand {
  0% {
    transform: translate(-50%, -50%) scale(0.5);
    opacity: 0.5;
  }
  100% {
    transform: translate(-50%, -50%) scale(1.5);
    opacity: 0;
  }
}

/* ===== 登录卡片 ===== */
.login-card-wrapper {
  width: 100%;
  max-width: 420px;
  opacity: 0;
  transform: translateY(32px);
  transition: all 0.8s var(--ease-smooth);
  position: relative;
  z-index: 10;
}

.login-card-wrapper.mounted {
  opacity: 1;
  transform: translateY(0);
}

.login-card {
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.15) !important;
  border-radius: var(--radius-2xl) !important;
  overflow: hidden;
  position: relative;
}

.login-card :deep(.el-card__body) {
  padding: 44px;
  position: relative;
}

.card-wave-top {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, rgb(var(--ocean-surface)), rgb(var(--aqua-glow)), rgb(var(--ocean-shallow)));
  background-size: 200% 100%;
  animation: gradient-flow 4s linear infinite;
}

@keyframes gradient-flow {
  0% { background-position: 0% center; }
  100% { background-position: 200% center; }
}

.card-decoration {
  position: absolute;
  top: 0;
  right: 0;
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.1), transparent 70%);
  filter: blur(40px);
  pointer-events: none;
}

/* ===== 登录头部 ===== */
.login-header {
  text-align: center;
  margin-bottom: 32px;
}

.login-icon {
  width: 76px;
  height: 76px;
  margin: 0 auto 20px;
  border-radius: var(--radius-xl);
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow:
    0 15px 40px rgba(14, 165, 233, 0.35),
    0 0 60px rgba(56, 189, 248, 0.2);
  color: white;
  animation: pulse-glow 3s ease-in-out infinite;
}

@keyframes pulse-glow {
  0%, 100% {
    box-shadow:
      0 15px 40px rgba(14, 165, 233, 0.35),
      0 0 60px rgba(56, 189, 248, 0.2);
  }
  50% {
    box-shadow:
      0 20px 50px rgba(14, 165, 233, 0.45),
      0 0 80px rgba(56, 189, 248, 0.3);
  }
}

.wave-svg {
  animation: wave-flow 2s ease-in-out infinite;
}

@keyframes wave-flow {
  0%, 100% { transform: translateX(-2px); }
  50% { transform: translateX(2px); }
}

.login-title {
  font-size: 26px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
}

.login-subtitle {
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
  margin: 0;
}

/* ===== 表单错误 ===== */
.login-error {
  margin-bottom: 20px;
  background: rgba(239, 68, 68, 0.1) !important;
  border: 1px solid rgba(239, 68, 68, 0.2) !important;
}

/* ===== 登录按钮 ===== */
.login-btn {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: var(--radius-lg) !important;
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid))) !important;
  border: none !important;
  box-shadow: 0 8px 25px rgba(14, 165, 233, 0.3);
  transition: all 0.3s var(--ease-smooth) !important;
}

.login-btn:hover {
  box-shadow: 0 12px 35px rgba(14, 165, 233, 0.4);
  transform: translateY(-2px);
}

/* ===== 信任标签 ===== */
.login-trust {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 14px;
  margin-top: 28px;
  padding: 16px 20px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.08), rgba(34, 211, 238, 0.04));
  border: 1px solid rgba(56, 189, 248, 0.1);
}

.trust-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.6);
}

.trust-icon {
  color: rgb(var(--ocean-surface));
  display: flex;
  align-items: center;
}

.trust-divider {
  width: 1px;
  height: 14px;
  background: rgba(56, 189, 248, 0.2);
}

/* ===== 底部链接 ===== */
.login-footer {
  text-align: center;
  color: rgba(255, 255, 255, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.login-footer :deep(.el-link) {
  color: rgb(var(--ocean-surface)) !important;
}

.login-footer :deep(.el-link:hover) {
  color: rgb(var(--aqua-glow)) !important;
}

/* ===== Element Plus 样式覆盖 ===== */
:deep(.el-form-item__label) {
  color: rgba(255, 255, 255, 0.8);
}

:deep(.el-input__wrapper) {
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.15);
  border-radius: var(--radius-md);
  box-shadow: none !important;
  transition: all 0.3s var(--ease-smooth);
}

:deep(.el-input__wrapper:hover) {
  border-color: rgba(56, 189, 248, 0.3);
  background: rgba(56, 189, 248, 0.08);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: rgb(var(--ocean-surface));
  background: rgba(56, 189, 248, 0.1);
  box-shadow: 0 0 0 3px rgba(56, 189, 248, 0.15) !important;
}

:deep(.el-input__inner) {
  color: #fff;
}

:deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.35);
}

:deep(.el-input__prefix) {
  color: rgba(56, 189, 248, 0.6);
}

:deep(.el-divider) {
  border-color: rgba(56, 189, 248, 0.1);
  margin: 28px 0;
}
</style>
