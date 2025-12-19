<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useUserStore } from '@/store/user'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, InfoFilled } from '@element-plus/icons-vue'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const form = reactive({
  username: '',
  password: '',
  confirmPassword: ''
})
const loading = ref(false)
const error = ref('')
const mounted = ref(false)

const validateConfirmPassword = (_rule: unknown, value: string, callback: (error?: Error) => void) => {
  if (value !== form.password) {
    callback(new Error('两次密码输入不一致'))
  } else {
    callback()
  }
}

const rules = reactive<FormRules>({
  username: [
    { required: true, message: '请输入用户名', trigger: 'blur' },
    { min: 3, max: 50, message: '用户名长度为3-50个字符', trigger: 'blur' }
  ],
  password: [
    { required: true, message: '请输入密码', trigger: 'blur' },
    { min: 6, message: '密码至少6个字符', trigger: 'blur' }
  ],
  confirmPassword: [
    { required: true, message: '请确认密码', trigger: 'blur' },
    { validator: validateConfirmPassword, trigger: 'blur' }
  ]
})

onMounted(() => {
  setTimeout(() => { mounted.value = true }, 100)
})

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    loading.value = true
    error.value = ''
    try {
      await userStore.register({ username: form.username, password: form.password })
      router.push('/')
    } catch (e) {
      error.value = e instanceof Error ? e.message : '注册失败'
    } finally {
      loading.value = false
    }
  })
}
</script>

<template>
  <MainLayout>
    <div class="register-container">
      <!-- 海洋主题背景 -->
      <div class="register-bg">
        <div class="register-glow register-glow-1"></div>
        <div class="register-glow register-glow-2"></div>
        <div class="register-ripple register-ripple-1"></div>
        <div class="register-ripple register-ripple-2"></div>
      </div>

      <div class="register-card-wrapper" :class="{ mounted }">
        <el-card class="register-card" shadow="never">
          <!-- 装饰波浪 -->
          <div class="card-wave-top"></div>
          <div class="card-decoration"></div>

          <div class="register-header">
            <div class="register-icon">
              <svg viewBox="0 0 32 32" fill="none" width="32" height="32" class="wave-svg">
                <circle cx="16" cy="16" r="10" stroke="currentColor" stroke-width="2" fill="none"/>
                <path d="M8 16C8 16 11 12 14 16C17 20 20 12 23 16" stroke="currentColor" stroke-width="2" stroke-linecap="round"/>
              </svg>
            </div>
            <h1 class="register-title">加入潮汐</h1>
            <p class="register-subtitle">开启你的海洋训练之旅</p>
          </div>

          <div class="register-perks">
            <div class="perk-item">
              <span class="perk-icon">
                <svg viewBox="0 0 16 16" fill="none" width="12" height="12">
                  <path d="M2 8L6 12L14 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </span>
              <span>浪迹打卡</span>
            </div>
            <div class="perk-item">
              <span class="perk-icon">
                <svg viewBox="0 0 16 16" fill="none" width="12" height="12">
                  <path d="M2 8L6 12L14 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </span>
              <span>潮汐统计</span>
            </div>
            <div class="perk-item">
              <span class="perk-icon">
                <svg viewBox="0 0 16 16" fill="none" width="12" height="12">
                  <path d="M2 8L6 12L14 4" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                </svg>
              </span>
              <span>深海排行</span>
            </div>
          </div>

          <el-form ref="formRef" :model="form" :rules="rules" label-position="top" size="large" @submit.prevent="handleSubmit">
            <el-form-item label="用户名" prop="username">
              <el-input v-model="form.username" placeholder="3-50个字符" :prefix-icon="User" />
            </el-form-item>

            <el-form-item label="密码" prop="password">
              <el-input v-model="form.password" type="password" placeholder="至少6个字符" :prefix-icon="Lock" show-password />
            </el-form-item>

            <el-form-item label="确认密码" prop="confirmPassword">
              <el-input v-model="form.confirmPassword" type="password" placeholder="再次输入密码" :prefix-icon="Lock" show-password />
            </el-form-item>

            <el-alert v-if="error" :title="error" type="error" show-icon :closable="false" class="register-error" />

            <el-form-item>
              <el-button type="primary" native-type="submit" :loading="loading" class="register-btn">
                {{ loading ? '注册中...' : '注册账号' }}
              </el-button>
            </el-form-item>
          </el-form>

          <div class="privacy-notice">
            <el-icon><InfoFilled /></el-icon>
            <p><span class="highlight">隐私守护：</span>如深海般守护，仅需用户名和密码</p>
          </div>

          <el-divider />

          <div class="register-footer">
            <span>已有账号?</span>
            <RouterLink to="/login"><el-link type="primary">立即登录</el-link></RouterLink>
          </div>
        </el-card>
      </div>
    </div>
  </MainLayout>
</template>

<style scoped>
.register-container {
  min-height: calc(100vh - 72px);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 24px;
  position: relative;
}

/* ===== 海洋背景效果 ===== */
.register-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.register-glow {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
}

.register-glow-1 {
  top: 25%;
  left: 20%;
  width: 320px;
  height: 320px;
  background: radial-gradient(circle, rgba(34, 211, 238, 0.18), transparent 70%);
  animation: float 8s ease-in-out infinite;
}

.register-glow-2 {
  bottom: 25%;
  right: 20%;
  width: 280px;
  height: 280px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.15), transparent 70%);
  animation: float 10s ease-in-out infinite;
  animation-delay: -2s;
}

.register-ripple {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  border: 1px solid rgba(34, 211, 238, 0.15);
  animation: ripple-expand 8s ease-out infinite;
}

.register-ripple-1 {
  width: 450px;
  height: 450px;
}

.register-ripple-2 {
  width: 650px;
  height: 650px;
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

/* ===== 注册卡片 ===== */
.register-card-wrapper {
  width: 100%;
  max-width: 420px;
  opacity: 0;
  transform: translateY(32px);
  transition: all 0.8s var(--ease-smooth);
  position: relative;
  z-index: 10;
}

.register-card-wrapper.mounted {
  opacity: 1;
  transform: translateY(0);
}

.register-card {
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(34, 211, 238, 0.15) !important;
  border-radius: var(--radius-2xl) !important;
  overflow: hidden;
  position: relative;
}

.register-card :deep(.el-card__body) {
  padding: 44px;
  position: relative;
}

.card-wave-top {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 4px;
  background: linear-gradient(90deg, rgb(var(--aqua-glow)), rgb(var(--ocean-surface)), rgb(var(--seaweed-green)));
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
  background: radial-gradient(circle, rgba(34, 211, 238, 0.1), transparent 70%);
  filter: blur(40px);
  pointer-events: none;
}

/* ===== 注册头部 ===== */
.register-header {
  text-align: center;
  margin-bottom: 24px;
}

.register-icon {
  width: 76px;
  height: 76px;
  margin: 0 auto 20px;
  border-radius: var(--radius-xl);
  background: linear-gradient(135deg, rgb(var(--aqua-glow)), rgb(var(--ocean-shallow)));
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow:
    0 15px 40px rgba(34, 211, 238, 0.35),
    0 0 60px rgba(34, 211, 238, 0.2);
  color: white;
  animation: pulse-glow 3s ease-in-out infinite;
}

@keyframes pulse-glow {
  0%, 100% {
    box-shadow:
      0 15px 40px rgba(34, 211, 238, 0.35),
      0 0 60px rgba(34, 211, 238, 0.2);
  }
  50% {
    box-shadow:
      0 20px 50px rgba(34, 211, 238, 0.45),
      0 0 80px rgba(34, 211, 238, 0.3);
  }
}

.wave-svg {
  animation: wave-flow 2s ease-in-out infinite;
}

@keyframes wave-flow {
  0%, 100% { transform: translateX(-2px); }
  50% { transform: translateX(2px); }
}

.register-title {
  font-size: 26px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
}

.register-subtitle {
  color: rgba(255, 255, 255, 0.65);
  font-size: 14px;
  margin: 0;
}

/* ===== 特性标签 ===== */
.register-perks {
  display: flex;
  justify-content: center;
  gap: 20px;
  margin-bottom: 28px;
  padding: 14px 20px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgba(34, 211, 238, 0.08), rgba(56, 189, 248, 0.04));
  border: 1px solid rgba(34, 211, 238, 0.12);
}

.perk-item {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
}

.perk-icon {
  color: rgb(var(--aqua-glow));
  display: flex;
  align-items: center;
}

/* ===== 表单错误 ===== */
.register-error {
  margin-bottom: 20px;
  background: rgba(239, 68, 68, 0.1) !important;
  border: 1px solid rgba(239, 68, 68, 0.2) !important;
}

/* ===== 注册按钮 ===== */
.register-btn {
  width: 100%;
  height: 50px;
  font-size: 16px;
  font-weight: 600;
  border-radius: var(--radius-lg) !important;
  background: linear-gradient(135deg, rgb(var(--aqua-glow)), rgb(var(--ocean-shallow))) !important;
  border: none !important;
  box-shadow: 0 8px 25px rgba(34, 211, 238, 0.3);
  transition: all 0.3s var(--ease-smooth) !important;
}

.register-btn:hover {
  box-shadow: 0 12px 35px rgba(34, 211, 238, 0.4);
  transform: translateY(-2px);
}

/* ===== 隐私提示 ===== */
.privacy-notice {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  border-radius: var(--radius-lg);
  background: linear-gradient(135deg, rgba(34, 211, 238, 0.08), rgba(56, 189, 248, 0.04));
  border: 1px solid rgba(34, 211, 238, 0.15);
  margin-top: 24px;
}

.privacy-notice .el-icon {
  color: rgb(var(--aqua-glow));
  font-size: 20px;
  flex-shrink: 0;
  margin-top: 2px;
}

.privacy-notice p {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
  line-height: 1.5;
}

.privacy-notice .highlight {
  color: rgb(var(--aqua-glow));
  font-weight: 500;
}

/* ===== 底部链接 ===== */
.register-footer {
  text-align: center;
  color: rgba(255, 255, 255, 0.65);
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
}

.register-footer :deep(.el-link) {
  color: rgb(var(--aqua-glow)) !important;
}

.register-footer :deep(.el-link:hover) {
  color: rgb(var(--ocean-surface)) !important;
}

/* ===== Element Plus 样式覆盖 ===== */
:deep(.el-form-item__label) {
  color: rgba(255, 255, 255, 0.8);
}

:deep(.el-input__wrapper) {
  background: rgba(34, 211, 238, 0.05);
  border: 1px solid rgba(34, 211, 238, 0.15);
  border-radius: var(--radius-md);
  box-shadow: none !important;
  transition: all 0.3s var(--ease-smooth);
}

:deep(.el-input__wrapper:hover) {
  border-color: rgba(34, 211, 238, 0.3);
  background: rgba(34, 211, 238, 0.08);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: rgb(var(--aqua-glow));
  background: rgba(34, 211, 238, 0.1);
  box-shadow: 0 0 0 3px rgba(34, 211, 238, 0.15) !important;
}

:deep(.el-input__inner) {
  color: #fff;
}

:deep(.el-input__inner::placeholder) {
  color: rgba(255, 255, 255, 0.45);
}

:deep(.el-input__prefix) {
  color: rgba(34, 211, 238, 0.6);
}

:deep(.el-divider) {
  border-color: rgba(34, 211, 238, 0.1);
  margin: 28px 0;
}
</style>
