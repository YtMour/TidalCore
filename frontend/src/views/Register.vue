<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { Card, Input, Button } from '@/components/ui'
import { useUserStore } from '@/store/user'
import { UserPlus, User, Lock, Shield, AlertCircle, CheckCircle, Sparkles, Zap } from 'lucide-vue-next'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: '',
  confirmPassword: ''
})
const loading = ref(false)
const error = ref('')
const mounted = ref(false)

const benefits = [
  { icon: CheckCircle, text: '记录训练进度' },
  { icon: Sparkles, text: '个人数据统计' },
  { icon: Zap, text: '冲击排行榜' }
]

onMounted(() => {
  setTimeout(() => {
    mounted.value = true
  }, 100)
})

async function handleSubmit() {
  error.value = ''

  if (!form.value.username || !form.value.password) {
    error.value = '请填写用户名和密码'
    return
  }

  if (form.value.username.length < 3) {
    error.value = '用户名至少3个字符'
    return
  }

  if (form.value.password.length < 6) {
    error.value = '密码至少6个字符'
    return
  }

  if (form.value.password !== form.value.confirmPassword) {
    error.value = '两次密码输入不一致'
    return
  }

  loading.value = true

  try {
    await userStore.register({
      username: form.value.username,
      password: form.value.password
    })
    router.push('/')
  } catch (e) {
    error.value = e instanceof Error ? e.message : '注册失败'
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <MainLayout>
    <div class="min-h-[70vh] flex items-center justify-center py-12 relative">
      <!-- Decorative background elements -->
      <div class="absolute inset-0 overflow-hidden pointer-events-none">
        <div class="absolute top-1/3 left-1/4 w-72 h-72 bg-gradient-to-br from-emerald-500/15 to-teal-500/10 rounded-full blur-3xl animate-float-slow"></div>
        <div class="absolute bottom-1/3 right-1/4 w-56 h-56 bg-gradient-to-br from-violet-500/15 to-purple-500/10 rounded-full blur-3xl animate-float-medium" style="animation-delay: -2s;"></div>
      </div>

      <div
        class="w-full max-w-md relative transition-all duration-1000"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
      >
        <Card padding="lg" class="relative overflow-hidden">
          <!-- Decorative corner gradient -->
          <div class="absolute top-0 right-0 w-40 h-40 bg-gradient-to-bl from-emerald-500/10 to-transparent rounded-full blur-2xl pointer-events-none"></div>

          <div class="relative">
            <!-- Header -->
            <div class="text-center mb-8">
              <div class="w-18 h-18 mx-auto mb-5 rounded-2xl bg-gradient-to-br from-emerald-500 to-teal-500 flex items-center justify-center shadow-xl shadow-emerald-500/30 p-4">
                <UserPlus class="w-9 h-9 text-white" />
              </div>
              <h1 class="text-2xl font-bold text-white mb-2">创建账号</h1>
              <p class="text-white/50">开始你的健康训练之旅</p>
            </div>

            <!-- Benefits -->
            <div class="flex items-center justify-center gap-4 mb-6 flex-wrap">
              <div
                v-for="benefit in benefits"
                :key="benefit.text"
                class="flex items-center gap-1.5 text-xs text-white/50"
              >
                <component :is="benefit.icon" class="w-3.5 h-3.5 text-emerald-400" />
                <span>{{ benefit.text }}</span>
              </div>
            </div>

            <!-- Form -->
            <form @submit.prevent="handleSubmit" class="space-y-5">
              <div class="relative">
                <div class="absolute left-4 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none mt-3">
                  <User class="w-5 h-5" />
                </div>
                <Input
                  v-model="form.username"
                  type="text"
                  label="用户名"
                  placeholder="3-50个字符"
                  class="[&_input]:pl-12"
                />
              </div>

              <div class="relative">
                <div class="absolute left-4 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none mt-3">
                  <Lock class="w-5 h-5" />
                </div>
                <Input
                  v-model="form.password"
                  type="password"
                  label="密码"
                  placeholder="至少6个字符"
                  class="[&_input]:pl-12"
                />
              </div>

              <div class="relative">
                <div class="absolute left-4 top-1/2 -translate-y-1/2 text-white/30 pointer-events-none mt-3">
                  <Lock class="w-5 h-5" />
                </div>
                <Input
                  v-model="form.confirmPassword"
                  type="password"
                  label="确认密码"
                  placeholder="再次输入密码"
                  class="[&_input]:pl-12"
                />
              </div>

              <!-- Error Message -->
              <Transition name="slide-up">
                <div v-if="error" class="flex items-center gap-2 p-3 rounded-xl bg-red-500/10 border border-red-500/20">
                  <AlertCircle class="w-5 h-5 text-red-400 shrink-0" />
                  <span class="text-sm text-red-400">{{ error }}</span>
                </div>
              </Transition>

              <!-- Submit Button -->
              <Button
                type="submit"
                variant="primary"
                size="lg"
                :loading="loading"
                full-width
                class="!rounded-xl !bg-gradient-to-r !from-emerald-500 !to-teal-500"
              >
                <UserPlus class="w-5 h-5" />
                {{ loading ? '注册中...' : '注册' }}
              </Button>
            </form>

            <!-- Privacy Notice -->
            <div class="mt-6 p-4 rounded-xl bg-gradient-to-r from-emerald-500/10 to-teal-500/10 border border-emerald-500/20">
              <div class="flex items-start gap-3">
                <Shield class="w-5 h-5 text-emerald-400 shrink-0 mt-0.5" />
                <p class="text-sm text-white/60">
                  <span class="text-emerald-400 font-medium">隐私保护：</span>
                  仅需用户名和密码，我们不收集任何敏感个人信息
                </p>
              </div>
            </div>

            <!-- Footer -->
            <div class="mt-6 pt-6 border-t border-white/5 text-center">
              <p class="text-white/50">
                已有账号?
                <RouterLink to="/login" class="text-emerald-400 hover:text-emerald-300 font-medium transition-colors animated-underline">
                  立即登录
                </RouterLink>
              </p>
            </div>
          </div>
        </Card>
      </div>
    </div>
  </MainLayout>
</template>
