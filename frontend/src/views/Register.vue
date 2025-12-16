<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import MainLayout from '@/layouts/MainLayout.vue'
import { useUserStore } from '@/store/user'

const router = useRouter()
const userStore = useUserStore()

const form = ref({
  username: '',
  password: '',
  confirmPassword: ''
})
const loading = ref(false)
const error = ref('')

async function handleSubmit() {
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
  error.value = ''

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
    <div class="max-w-md mx-auto">
      <div class="bg-white/5 rounded-2xl p-8">
        <h1 class="text-2xl font-bold text-white text-center mb-6">注册</h1>

        <form @submit.prevent="handleSubmit" class="space-y-4">
          <div>
            <label class="block text-white/60 text-sm mb-1">用户名</label>
            <input
              v-model="form.username"
              type="text"
              placeholder="3-50个字符"
              class="w-full px-4 py-3 bg-white/10 rounded-lg text-white placeholder-white/30 border border-white/10 focus:border-purple-500 focus:outline-none"
            />
          </div>

          <div>
            <label class="block text-white/60 text-sm mb-1">密码</label>
            <input
              v-model="form.password"
              type="password"
              placeholder="至少6个字符"
              class="w-full px-4 py-3 bg-white/10 rounded-lg text-white placeholder-white/30 border border-white/10 focus:border-purple-500 focus:outline-none"
            />
          </div>

          <div>
            <label class="block text-white/60 text-sm mb-1">确认密码</label>
            <input
              v-model="form.confirmPassword"
              type="password"
              placeholder="再次输入密码"
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
            {{ loading ? '注册中...' : '注册' }}
          </button>
        </form>

        <p class="text-center text-white/50 mt-6">
          已有账号?
          <RouterLink to="/login" class="text-purple-400 hover:text-purple-300">
            立即登录
          </RouterLink>
        </p>

        <p class="text-center text-white/30 text-xs mt-4">
          仅需用户名和密码，我们重视你的隐私
        </p>
      </div>
    </div>
  </MainLayout>
</template>
