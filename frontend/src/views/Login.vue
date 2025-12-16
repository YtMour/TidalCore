<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useUserStore } from '@/store/user'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: ''
})
const loading = ref(false)
const error = ref('')

onMounted(() => {
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
    <div class="max-w-md mx-auto">
      <div class="bg-white/5 rounded-2xl p-8">
        <h1 class="text-2xl font-bold text-white text-center mb-6">登录</h1>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-white/60 text-sm mb-1">用户名</label>
            <input
              v-model="form.username"
              type="text"
              placeholder="请输入用户名"
              autocomplete="username"
              class="w-full px-4 py-3 bg-white/10 rounded-lg text-white placeholder-white/30 border border-white/10 focus:border-purple-500 focus:outline-none"
            />
          </div>

          <div>
            <label class="block text-white/60 text-sm mb-1">密码</label>
            <input
              v-model="form.password"
              type="password"
              placeholder="请输入密码"
              autocomplete="current-password"
              class="w-full px-4 py-3 bg-white/10 rounded-lg text-white placeholder-white/30 border border-white/10 focus:border-purple-500 focus:outline-none"
            />
          </div>

          <div v-if="error" class="text-red-400 text-sm text-center">
            {{ error }}
          </div>

          <button
            type="submit"
            :disabled="loading"
            class="w-full py-3 rounded-lg bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 text-white font-semibold transition-all disabled:opacity-50"
          >
            {{ loading ? '登录中...' : '登录' }}
          </button>
        </form>

        <p class="text-center text-white/50 mt-6">
          还没有账号?
          <RouterLink to="/register" class="text-purple-400 hover:text-purple-300">
            立即注册
          </RouterLink>
        </p>
      </div>
    </div>
  </MainLayout>
</template>
