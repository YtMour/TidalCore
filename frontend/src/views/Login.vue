<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { Card, Input, Button } from '@/components/ui'
import { useUserStore } from '@/store/user'
import { LogIn, User, Lock, AlertCircle, Sparkles, Shield, Zap } from 'lucide-vue-next'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: ''
})
const loading = ref(false)
const error = ref('')
const mounted = ref(false)

const features = [
  { icon: Shield, text: '安全加密存储' },
  { icon: Zap, text: '快速登录体验' },
  { icon: Sparkles, text: '个性化训练记录' }
]

onMounted(() => {
  setTimeout(() => {
    mounted.value = true
  }, 100)

  if (route.query.expired === '1') {
    error.value = '登录已过期，请重新登录'
  }
})

function isValidRedirect(url: string): boolean {
  if (!url) return false
  if (url.startsWith('/') && !url.startsWith('//')) {
    return true
  }
  return false
}

async function handleSubmit() {
  if (!form.value.username || !form.value.password) {
    error.value = '请填写用户名和密码'
    return
  }

  loading.value = true
  error.value = ''

  try {
    await userStore.login(form.value)
    const redirect = route.query.redirect as string
    if (isValidRedirect(redirect)) {
      router.push(redirect)
    } else {
      router.push('/')
    }
  } catch (e) {
    error.value = e instanceof Error ? e.message : '登录失败'
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
        <div class="absolute top-1/4 left-1/4 w-64 h-64 bg-gradient-to-br from-violet-500/20 to-purple-500/10 rounded-full blur-3xl animate-float-slow"></div>
        <div class="absolute bottom-1/4 right-1/4 w-48 h-48 bg-gradient-to-br from-pink-500/15 to-rose-500/10 rounded-full blur-3xl animate-float-medium" style="animation-delay: -3s;"></div>
      </div>

      <div
        class="w-full max-w-md relative transition-all duration-1000"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'"
      >
        <Card padding="lg" class="relative overflow-hidden">
          <!-- Decorative corner gradient -->
          <div class="absolute top-0 right-0 w-40 h-40 bg-gradient-to-bl from-violet-500/10 to-transparent rounded-full blur-2xl pointer-events-none"></div>

          <div class="relative">
            <!-- Header -->
            <div class="text-center mb-8">
              <div class="w-18 h-18 mx-auto mb-5 rounded-2xl bg-gradient-to-br from-violet-500 to-pink-500 flex items-center justify-center shadow-xl shadow-violet-500/30 p-4">
                <LogIn class="w-9 h-9 text-white" />
              </div>
              <h1 class="text-2xl font-bold text-white mb-2">欢迎回来</h1>
              <p class="text-white/50">登录以继续你的训练之旅</p>
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
                  placeholder="请输入用户名"
                  autocomplete="username"
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
                  placeholder="请输入密码"
                  autocomplete="current-password"
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
                class="!rounded-xl"
              >
                <LogIn class="w-5 h-5" />
                {{ loading ? '登录中...' : '登录' }}
              </Button>
            </form>

            <!-- Features -->
            <div class="mt-6 flex items-center justify-center gap-4 flex-wrap">
              <div
                v-for="feature in features"
                :key="feature.text"
                class="flex items-center gap-1.5 text-xs text-white/40"
              >
                <component :is="feature.icon" class="w-3.5 h-3.5 text-violet-400" />
                <span>{{ feature.text }}</span>
              </div>
            </div>

            <!-- Footer -->
            <div class="mt-8 pt-6 border-t border-white/5 text-center">
              <p class="text-white/50">
                还没有账号?
                <RouterLink to="/register" class="text-violet-400 hover:text-violet-300 font-medium transition-colors animated-underline">
                  立即注册
                </RouterLink>
              </p>
            </div>
          </div>
        </Card>
      </div>
    </div>
  </MainLayout>
</template>
