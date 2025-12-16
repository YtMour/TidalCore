<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useScrollLock } from '@vueuse/core'
import {
  Menu,
  X,
  Home,
  Timer,
  Trophy,
  User,
  LogOut,
  LogIn,
  Github
} from 'lucide-vue-next'

const userStore = useUserStore()
const route = useRoute()

const mobileMenuOpen = ref(false)
const scrollLock = useScrollLock(document.body)
const scrolled = ref(false)

// 监听滚动
onMounted(() => {
  const handleScroll = () => {
    scrolled.value = window.scrollY > 20
  }
  window.addEventListener('scroll', handleScroll)
  return () => window.removeEventListener('scroll', handleScroll)
})

watch(mobileMenuOpen, (open) => {
  scrollLock.value = open
})

watch(() => route.path, () => {
  mobileMenuOpen.value = false
})

function handleLogout() {
  userStore.logout()
  mobileMenuOpen.value = false
}

const navLinks = [
  { to: '/', label: '首页', icon: Home },
  { to: '/train', label: '训练', icon: Timer },
  { to: '/leaderboard', label: '排行榜', icon: Trophy }
]
</script>

<template>
  <div class="min-h-screen flex flex-col">
    <!-- Animated Background -->
    <div class="animated-bg"></div>
    <div class="particles-bg"></div>
    <div class="noise-overlay"></div>

    <!-- Navigation -->
    <nav
      class="sticky top-0 z-40 transition-all duration-300"
      :class="[
        scrolled
          ? 'bg-slate-900/90 backdrop-blur-xl border-b border-white/5 shadow-lg shadow-black/10'
          : 'bg-transparent border-b border-transparent'
      ]"
    >
      <div class="container-app">
        <div class="h-16 flex items-center justify-between">
          <!-- Logo -->
          <RouterLink
            to="/"
            class="flex items-center gap-2.5 text-white font-bold text-xl group"
          >
            <span class="text-2xl group-hover:scale-110 transition-transform duration-300">🌊</span>
            <span class="gradient-text">TidalCore</span>
          </RouterLink>

          <!-- Desktop Navigation -->
          <div class="hidden md:flex items-center gap-1 bg-white/[0.02] rounded-xl px-1.5 py-1.5 border border-white/[0.06]">
            <RouterLink
              v-for="link in navLinks"
              :key="link.to"
              :to="link.to"
              class="flex items-center gap-2 px-4 py-2 rounded-lg text-white/60 hover:text-white transition-all duration-200"
              :class="route.path === link.to ? 'text-white bg-white/10' : 'hover:bg-white/5'"
            >
              <component :is="link.icon" class="w-4 h-4" />
              <span>{{ link.label }}</span>
            </RouterLink>
          </div>

          <!-- Desktop User Menu -->
          <div class="hidden md:flex items-center gap-2">
            <template v-if="userStore.isLoggedIn">
              <RouterLink
                to="/dashboard"
                class="flex items-center gap-2 px-4 py-2 rounded-lg text-white/60 hover:text-white hover:bg-white/5 transition-all duration-200 border border-transparent hover:border-white/[0.08]"
                :class="{ 'text-white bg-white/10 border-white/[0.08]': route.path === '/dashboard' }"
              >
                <User class="w-4 h-4" />
                <span>我的</span>
              </RouterLink>
              <button
                @click="handleLogout"
                class="flex items-center gap-2 px-4 py-2 rounded-lg text-white/60 hover:text-white hover:bg-white/5 transition-all duration-200"
              >
                <LogOut class="w-4 h-4" />
                <span>退出</span>
              </button>
            </template>
            <template v-else>
              <RouterLink
                to="/login"
                class="group btn-gradient px-5 py-2 rounded-lg text-white text-sm font-medium flex items-center gap-2"
              >
                <LogIn class="w-4 h-4" />
                <span>登录</span>
              </RouterLink>
            </template>
          </div>

          <!-- Mobile Menu Button -->
          <button
            @click="mobileMenuOpen = !mobileMenuOpen"
            class="md:hidden w-10 h-10 rounded-lg flex items-center justify-center text-white/70 hover:text-white hover:bg-white/5 transition-all"
          >
            <Transition name="fade" mode="out-in">
              <X v-if="mobileMenuOpen" class="w-6 h-6" />
              <Menu v-else class="w-6 h-6" />
            </Transition>
          </button>
        </div>
      </div>
    </nav>

    <!-- Mobile Menu -->
    <Transition name="slide-up">
      <div
        v-if="mobileMenuOpen"
        class="fixed inset-0 z-30 md:hidden"
      >
        <div class="absolute inset-0 bg-slate-900/95 backdrop-blur-xl pt-20">
          <div class="container-app py-6 space-y-2">
            <RouterLink
              v-for="link in navLinks"
              :key="link.to"
              :to="link.to"
              class="flex items-center gap-3 px-4 py-3 rounded-lg text-lg text-white/70 hover:text-white hover:bg-white/5 transition-all"
              :class="{ 'text-white bg-white/10': route.path === link.to }"
            >
              <component :is="link.icon" class="w-5 h-5" />
              <span>{{ link.label }}</span>
            </RouterLink>

            <div class="divider my-4"></div>

            <template v-if="userStore.isLoggedIn">
              <RouterLink
                to="/dashboard"
                class="flex items-center gap-3 px-4 py-3 rounded-lg text-lg text-white/70 hover:text-white hover:bg-white/5 transition-all"
                :class="{ 'text-white bg-white/10': route.path === '/dashboard' }"
              >
                <User class="w-5 h-5" />
                <span>我的</span>
              </RouterLink>
              <button
                @click="handleLogout"
                class="w-full flex items-center gap-3 px-4 py-3 rounded-lg text-lg text-white/70 hover:text-white hover:bg-white/5 transition-all"
              >
                <LogOut class="w-5 h-5" />
                <span>退出登录</span>
              </button>
            </template>
            <template v-else>
              <RouterLink
                to="/login"
                class="flex items-center gap-3 px-4 py-3 rounded-lg text-lg text-white/70 hover:text-white hover:bg-white/5 transition-all"
              >
                <LogIn class="w-5 h-5" />
                <span>登录</span>
              </RouterLink>
              <RouterLink
                to="/register"
                class="flex items-center justify-center gap-2 mt-4 btn-gradient px-6 py-3 rounded-lg text-white font-medium"
              >
                <span>立即注册</span>
              </RouterLink>
            </template>
          </div>
        </div>
      </div>
    </Transition>

    <!-- Main Content -->
    <main class="flex-1 container-app py-8 md:py-12">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="relative border-t border-white/5 py-10 mt-auto overflow-hidden">
      <!-- Decorative gradient -->
      <div class="absolute inset-0 -z-10">
        <div class="absolute bottom-0 left-1/4 w-64 h-32 bg-gradient-to-t from-violet-500/5 to-transparent rounded-full blur-3xl"></div>
        <div class="absolute bottom-0 right-1/4 w-48 h-24 bg-gradient-to-t from-pink-500/5 to-transparent rounded-full blur-3xl"></div>
      </div>

      <div class="container-app">
        <div class="flex flex-col md:flex-row items-center justify-between gap-6">
          <!-- Logo and tagline -->
          <div class="flex flex-col items-center md:items-start gap-2">
            <div class="flex items-center gap-2 text-white/60">
              <span class="text-xl">🌊</span>
              <span class="font-semibold gradient-text">TidalCore</span>
            </div>
            <p class="text-sm text-white/40">开源盆底肌训练平台</p>
          </div>

          <!-- Links -->
          <div class="flex items-center gap-6 text-sm">
            <a
              href="https://github.com"
              target="_blank"
              rel="noopener noreferrer"
              class="flex items-center gap-1.5 text-white/40 hover:text-white/70 transition-colors"
            >
              <Github class="w-4 h-4" />
              <span>GitHub</span>
            </a>
            <span class="text-white/20">|</span>
            <span class="text-white/40">MIT License</span>
          </div>
        </div>

        <!-- Bottom bar -->
        <div class="mt-8 pt-6 border-t border-white/5 flex flex-col md:flex-row items-center justify-between gap-4 text-xs text-white/30">
          <p>Made with <span class="text-pink-400">♥</span> for health</p>
          <p>© 2024 TidalCore. All rights reserved.</p>
        </div>
      </div>
    </footer>
  </div>
</template>
