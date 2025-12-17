<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useUserStore } from '@/store/user'
import type { FormInstance, FormRules } from 'element-plus'
import { User, Lock, UserFilled, Select, DataLine, Trophy } from '@element-plus/icons-vue'

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
      <div class="register-bg">
        <div class="register-blob register-blob-1"></div>
        <div class="register-blob register-blob-2"></div>
      </div>

      <div class="register-card-wrapper" :class="{ mounted }">
        <el-card class="register-card" shadow="always">
          <div class="register-header">
            <div class="register-icon">
              <el-icon :size="36"><UserFilled /></el-icon>
            </div>
            <h1 class="register-title">创建账号</h1>
            <p class="register-subtitle">开始你的健康训练之旅</p>
          </div>

          <div class="register-benefits">
            <el-tag type="success" effect="plain" round><el-icon><Select /></el-icon> 记录进度</el-tag>
            <el-tag type="success" effect="plain" round><el-icon><DataLine /></el-icon> 数据统计</el-tag>
            <el-tag type="success" effect="plain" round><el-icon><Trophy /></el-icon> 排行榜</el-tag>
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
              <el-button type="success" native-type="submit" :loading="loading" class="register-btn">
                {{ loading ? '注册中...' : '注册' }}
              </el-button>
            </el-form-item>
          </el-form>

          <div class="privacy-notice">
            <el-icon><InfoFilled /></el-icon>
            <p><span class="highlight">隐私保护：</span>仅需用户名和密码，不收集敏感信息</p>
          </div>

          <el-divider />
          <div class="register-footer">
            <span>已有账号?</span>
            <RouterLink to="/login"><el-link type="success">立即登录</el-link></RouterLink>
          </div>
        </el-card>
      </div>
    </div>
  </MainLayout>
</template>

<style scoped>
.register-container {
  min-height: 70vh;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 48px 16px;
  position: relative;
}

.register-bg {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.register-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
}

.register-blob-1 {
  top: 33%;
  left: 25%;
  width: 288px;
  height: 288px;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.15), rgba(20, 184, 166, 0.1));
  animation: float 6s ease-in-out infinite;
}

.register-blob-2 {
  bottom: 33%;
  right: 25%;
  width: 224px;
  height: 224px;
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.15), rgba(168, 85, 247, 0.1));
  animation: float 8s ease-in-out infinite;
  animation-delay: -2s;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-20px); }
}

.register-card-wrapper {
  width: 100%;
  max-width: 420px;
  opacity: 0;
  transform: translateY(32px);
  transition: all 0.8s ease;
}

.register-card-wrapper.mounted {
  opacity: 1;
  transform: translateY(0);
}

.register-card {
  background: rgba(30, 30, 46, 0.9);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 16px;
}

.register-card :deep(.el-card__body) {
  padding: 40px;
}

.register-header {
  text-align: center;
  margin-bottom: 24px;
}

.register-icon {
  width: 72px;
  height: 72px;
  margin: 0 auto 20px;
  border-radius: 16px;
  background: linear-gradient(135deg, #10b981, #14b8a6);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 40px rgba(16, 185, 129, 0.3);
  color: white;
}

.register-title {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 8px;
}

.register-subtitle {
  color: rgba(255, 255, 255, 0.5);
  font-size: 14px;
  margin: 0;
}

.register-benefits {
  display: flex;
  flex-wrap: wrap;
  justify-content: center;
  gap: 8px;
  margin-bottom: 24px;
}

.register-benefits .el-tag {
  display: flex;
  align-items: center;
  gap: 4px;
  background: rgba(16, 185, 129, 0.1);
  border-color: rgba(16, 185, 129, 0.2);
}

.register-error {
  margin-bottom: 20px;
}

.register-btn {
  width: 100%;
  height: 48px;
  font-size: 16px;
  border-radius: 12px;
  background: linear-gradient(135deg, #10b981, #14b8a6) !important;
  border: none !important;
}

.privacy-notice {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 16px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.1), rgba(20, 184, 166, 0.1));
  border: 1px solid rgba(16, 185, 129, 0.2);
  margin-top: 24px;
}

.privacy-notice .el-icon {
  color: #10b981;
  font-size: 20px;
  flex-shrink: 0;
  margin-top: 2px;
}

.privacy-notice p {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
  margin: 0;
}

.privacy-notice .highlight {
  color: #10b981;
  font-weight: 500;
}

.register-footer {
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
  border-color: rgba(16, 185, 129, 0.5);
}

:deep(.el-input__wrapper.is-focus) {
  border-color: #10b981;
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
