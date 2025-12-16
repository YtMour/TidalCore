<script setup lang="ts">
import { useUserStore } from '@/store/user'
import { useRouter } from 'vue-router'

const userStore = useUserStore()
const router = useRouter()

function handleLogout() {
  userStore.logout()
  router.push('/')
}
</script>

<template>
  <div class="min-h-screen bg-gradient-to-br from-slate-900 via-purple-900 to-slate-900">
    <nav class="border-b border-white/10 bg-black/20 backdrop-blur-sm">
      <div class="max-w-5xl mx-auto px-4 h-16 flex items-center justify-between">
        <RouterLink to="/" class="text-xl font-bold text-white flex items-center gap-2">
          <span class="text-2xl">🌊</span>
          <span>TidalCore</span>
        </RouterLink>

        <div class="flex items-center gap-6">
          <RouterLink to="/train" class="text-white/70 hover:text-white transition">
            训练
          </RouterLink>
          <RouterLink to="/leaderboard" class="text-white/70 hover:text-white transition">
            排行榜
          </RouterLink>

          <template v-if="userStore.isLoggedIn">
            <RouterLink to="/dashboard" class="text-white/70 hover:text-white transition">
              我的
            </RouterLink>
            <button
              @click="handleLogout"
              class="text-white/70 hover:text-white transition"
            >
              退出
            </button>
          </template>
          <template v-else>
            <RouterLink
              to="/login"
              class="px-4 py-2 rounded-lg bg-purple-600 hover:bg-purple-500 text-white transition"
            >
              登录
            </RouterLink>
          </template>
        </div>
      </div>
    </nav>

    <main class="max-w-5xl mx-auto px-4 py-8">
      <slot />
    </main>

    <footer class="border-t border-white/10 py-6 text-center text-white/40 text-sm">
      <p>TidalCore - 开源盆底肌训练平台</p>
    </footer>
  </div>
</template>
