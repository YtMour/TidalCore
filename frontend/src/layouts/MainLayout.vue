<script setup lang="ts">
import { ref, watch, onMounted, onUnmounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useScrollLock } from '@vueuse/core'
import { HomeFilled, Timer, Trophy, User, SwitchButton, Fold, Expand } from '@element-plus/icons-vue'

const userStore = useUserStore()
const route = useRoute()

const mobileMenuOpen = ref(false)
const scrollLock = useScrollLock(document.body)
const scrolled = ref(false)

// 生成气泡
const bubbles = ref<Array<{ id: number; size: number; left: number; delay: number; duration: number }>>([])
let bubbleId = 0

function createBubble() {
  const bubble = {
    id: bubbleId++,
    size: Math.random() * 20 + 8,
    left: Math.random() * 100,
    delay: Math.random() * 2,
    duration: Math.random() * 10 + 15
  }
  bubbles.value.push(bubble)

  // 移除旧气泡
  setTimeout(() => {
    bubbles.value = bubbles.value.filter(b => b.id !== bubble.id)
  }, (bubble.duration + bubble.delay) * 1000)
}

let bubbleInterval: number

onMounted(() => {
  const handleScroll = () => {
    scrolled.value = window.scrollY > 20
  }
  window.addEventListener('scroll', handleScroll)

  // 定期创建气泡
  bubbleInterval = window.setInterval(createBubble, 3000)
  // 初始气泡
  for (let i = 0; i < 5; i++) {
    setTimeout(createBubble, i * 600)
  }

  return () => window.removeEventListener('scroll', handleScroll)
})

onUnmounted(() => {
  if (bubbleInterval) clearInterval(bubbleInterval)
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
  { to: '/', label: '首页', icon: HomeFilled },
  { to: '/train', label: '训练', icon: Timer },
  { to: '/leaderboard', label: '排行榜', icon: Trophy }
]
</script>

<template>
  <div class="main-layout">
    <!-- 全屏海洋背景系统 -->
    <div class="ocean-bg"></div>
    <div class="ocean-gradient"></div>
    <div class="animated-bg"></div>
    <div class="particles-bg"></div>
    <div class="noise-overlay"></div>

    <!-- 全屏海浪动画层 -->
    <div class="wave-layer">
      <div class="wave wave-1"></div>
      <div class="wave wave-2"></div>
      <div class="wave wave-3"></div>
    </div>

    <!-- 涟漪层 -->
    <div class="ripple-layer"></div>

    <!-- 动态气泡 -->
    <div class="bubbles-container">
      <div
        v-for="bubble in bubbles"
        :key="bubble.id"
        class="bubble"
        :style="{
          width: bubble.size + 'px',
          height: bubble.size + 'px',
          left: bubble.left + '%',
          animationDelay: bubble.delay + 's',
          animationDuration: bubble.duration + 's'
        }"
      ></div>
    </div>

    <!-- Navigation - 透明悬浮 -->
    <header class="main-header" :class="{ scrolled }">
      <div class="header-content">
        <!-- Logo -->
        <RouterLink to="/" class="logo-link">
          <div class="logo-wave-icon">
            <svg viewBox="0 0 32 32" fill="none" class="wave-svg">
              <defs>
                <linearGradient id="waveGradient" x1="0%" y1="0%" x2="100%" y2="100%">
                  <stop offset="0%" stop-color="#38bdf8" />
                  <stop offset="50%" stop-color="#22d3ee" />
                  <stop offset="100%" stop-color="#0ea5e9" />
                </linearGradient>
              </defs>
              <circle cx="16" cy="16" r="14" fill="url(#waveGradient)" opacity="0.15"/>
              <path d="M6 16C6 16 9 12 13 16C17 20 21 12 25 16" stroke="url(#waveGradient)" stroke-width="2.5" stroke-linecap="round" class="wave-path"/>
            </svg>
          </div>
          <span class="logo-text">TidalCore</span>
        </RouterLink>

        <!-- Desktop Navigation -->
        <nav class="desktop-nav">
          <RouterLink
            v-for="link in navLinks"
            :key="link.to"
            :to="link.to"
            class="nav-link"
            :class="{ active: route.path === link.to }"
          >
            <el-icon :size="16"><component :is="link.icon" /></el-icon>
            <span>{{ link.label }}</span>
          </RouterLink>
        </nav>

        <!-- Desktop User Menu -->
        <div class="desktop-user">
          <template v-if="userStore.isLoggedIn">
            <RouterLink to="/dashboard" class="user-link" :class="{ active: route.path === '/dashboard' }">
              <el-icon :size="16"><User /></el-icon>
              <span>我的</span>
            </RouterLink>
            <button @click="handleLogout" class="logout-btn">
              <el-icon :size="16"><SwitchButton /></el-icon>
              <span>退出</span>
            </button>
          </template>
          <template v-else>
            <RouterLink to="/login" class="login-btn">
              <el-icon :size="16"><User /></el-icon>
              <span>登录</span>
            </RouterLink>
          </template>
        </div>

        <!-- Mobile Menu Button -->
        <button class="mobile-menu-btn" @click="mobileMenuOpen = !mobileMenuOpen">
          <el-icon :size="24">
            <Fold v-if="mobileMenuOpen" />
            <Expand v-else />
          </el-icon>
        </button>
      </div>
    </header>

    <!-- Mobile Menu Drawer -->
    <el-drawer
      v-model="mobileMenuOpen"
      direction="ttb"
      :show-close="false"
      :with-header="false"
      size="100%"
      class="mobile-drawer"
    >
      <div class="mobile-menu">
        <RouterLink
          v-for="link in navLinks"
          :key="link.to"
          :to="link.to"
          class="mobile-nav-item"
          :class="{ active: route.path === link.to }"
        >
          <el-icon :size="20"><component :is="link.icon" /></el-icon>
          <span>{{ link.label }}</span>
        </RouterLink>

        <el-divider />

        <template v-if="userStore.isLoggedIn">
          <RouterLink
            to="/dashboard"
            class="mobile-nav-item"
            :class="{ active: route.path === '/dashboard' }"
          >
            <el-icon :size="20"><User /></el-icon>
            <span>我的</span>
          </RouterLink>
          <div class="mobile-nav-item" @click="handleLogout">
            <el-icon :size="20"><SwitchButton /></el-icon>
            <span>退出登录</span>
          </div>
        </template>
        <template v-else>
          <RouterLink to="/login" class="mobile-nav-item">
            <el-icon :size="20"><User /></el-icon>
            <span>登录</span>
          </RouterLink>
          <RouterLink to="/register" class="mobile-register-btn">
            <el-button type="primary" size="large" class="w-full">立即注册</el-button>
          </RouterLink>
        </template>
      </div>
    </el-drawer>

    <!-- Main Content - 全屏无边界 -->
    <main class="main-content">
      <slot />
    </main>

    <!-- Footer - 简化透明 -->
    <footer class="main-footer">
      <div class="footer-content">
        <div class="footer-brand">
          <div class="footer-logo">
            <svg viewBox="0 0 32 32" fill="none" class="footer-wave-svg">
              <path d="M6 16C6 16 9 12 13 16C17 20 21 12 25 16" stroke="rgba(56, 189, 248, 0.5)" stroke-width="2" stroke-linecap="round"/>
            </svg>
            <span class="footer-logo-text">TidalCore</span>
          </div>
          <span class="footer-divider">·</span>
          <span class="footer-tagline">开源盆底肌训练平台</span>
        </div>

        <div class="footer-links">
          <a href="https://github.com" target="_blank" rel="noopener noreferrer" class="footer-link">GitHub</a>
          <span class="footer-dot">·</span>
          <span class="footer-license">MIT License</span>
          <span class="footer-dot">·</span>
          <span class="footer-copyright">© 2024</span>
        </div>
      </div>
    </footer>
  </div>
</template>

<style scoped>
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
}

/* ===== 气泡容器 ===== */
.bubbles-container {
  position: fixed;
  inset: 0;
  z-index: -2;
  pointer-events: none;
  overflow: hidden;
}

.bubbles-container .bubble {
  animation: bubble-rise linear forwards;
}

/* ===== Header - 完全透明悬浮 ===== */
.main-header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  z-index: 100;
  height: 72px;
  padding: 0 24px;
  background: transparent;
  border-bottom: 1px solid transparent;
  transition: background 0.4s var(--ease-smooth), border-color 0.4s var(--ease-smooth);
}

.main-header.scrolled {
  background: rgba(8, 15, 30, 0.6);
  backdrop-filter: blur(16px);
  -webkit-backdrop-filter: blur(16px);
  border-bottom-color: rgba(56, 189, 248, 0.08);
}

.header-content {
  max-width: 1400px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 24px;
}

/* ===== Logo ===== */
.logo-link {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  flex-shrink: 0;
}

.logo-wave-icon {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.4s var(--ease-smooth);
}

.wave-svg {
  width: 100%;
  height: 100%;
}

.wave-path {
  animation: wave-flow 2s ease-in-out infinite;
}

@keyframes wave-flow {
  0%, 100% { d: path("M6 16C6 16 9 12 13 16C17 20 21 12 25 16"); }
  50% { d: path("M6 16C6 16 9 20 13 16C17 12 21 20 25 16"); }
}

.logo-link:hover .logo-wave-icon {
  transform: scale(1.1) rotate(-5deg);
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, #38bdf8 0%, #22d3ee 50%, #0ea5e9 100%);
  background-size: 200% auto;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: gradient-flow 4s ease-in-out infinite;
}

@keyframes gradient-flow {
  0%, 100% { background-position: 0% center; }
  50% { background-position: 200% center; }
}

/* ===== Desktop Navigation ===== */
.desktop-nav {
  display: none;
  align-items: center;
  gap: 4px;
  padding: 4px;
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-lg);
}

@media (min-width: 768px) {
  .desktop-nav {
    display: flex;
  }
}

.nav-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  transition: all 0.3s var(--ease-smooth);
  position: relative;
  overflow: hidden;
}

.nav-link::before {
  content: '';
  position: absolute;
  bottom: 0;
  left: 50%;
  width: 0;
  height: 2px;
  background: linear-gradient(90deg, transparent, rgb(var(--ocean-surface)), transparent);
  transition: all 0.3s var(--ease-smooth);
  transform: translateX(-50%);
}

.nav-link:hover {
  color: #fff;
  background: rgba(56, 189, 248, 0.1);
}

.nav-link:hover::before {
  width: 80%;
}

.nav-link.active {
  color: #fff;
  background: rgba(56, 189, 248, 0.15);
}

.nav-link.active::before {
  width: 80%;
}

/* ===== Desktop User Menu ===== */
.desktop-user {
  display: none;
  align-items: center;
  gap: 8px;
}

@media (min-width: 768px) {
  .desktop-user {
    display: flex;
  }
}

.user-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  border: 1px solid transparent;
  transition: all 0.3s var(--ease-smooth);
}

.user-link:hover {
  color: #fff;
  background: rgba(56, 189, 248, 0.08);
  border-color: rgba(56, 189, 248, 0.15);
}

.user-link.active {
  color: rgb(var(--ocean-surface));
  background: rgba(56, 189, 248, 0.1);
  border-color: rgba(56, 189, 248, 0.2);
}

.logout-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 16px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all 0.3s var(--ease-smooth);
}

.logout-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.05);
}

.login-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 10px 20px;
  border-radius: var(--radius-md);
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  text-decoration: none;
  transition: all 0.3s var(--ease-smooth);
  box-shadow: 0 4px 15px rgba(14, 165, 233, 0.3);
}

.login-btn:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 25px rgba(14, 165, 233, 0.4);
}

/* ===== Mobile Menu Button ===== */
.mobile-menu-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 44px;
  height: 44px;
  border-radius: var(--radius-md);
  border: 1px solid rgba(56, 189, 248, 0.1);
  background: rgba(56, 189, 248, 0.05);
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.3s var(--ease-smooth);
}

.mobile-menu-btn:hover {
  color: #fff;
  background: rgba(56, 189, 248, 0.1);
  border-color: rgba(56, 189, 248, 0.2);
}

@media (min-width: 768px) {
  .mobile-menu-btn {
    display: none;
  }
}

/* ===== Mobile Drawer ===== */
.mobile-drawer :deep(.el-drawer__body) {
  background: linear-gradient(180deg, rgba(5, 8, 15, 0.98), rgba(8, 12, 22, 0.98));
  backdrop-filter: blur(20px);
  padding-top: 80px;
}

.mobile-menu {
  padding: 0 20px;
}

.mobile-nav-item {
  display: flex;
  align-items: center;
  gap: 14px;
  padding: 16px 18px;
  border-radius: var(--radius-lg);
  font-size: 16px;
  color: rgba(255, 255, 255, 0.7);
  text-decoration: none;
  cursor: pointer;
  transition: all 0.3s var(--ease-smooth);
  border: 1px solid transparent;
}

.mobile-nav-item:hover {
  color: #fff;
  background: rgba(56, 189, 248, 0.08);
  border-color: rgba(56, 189, 248, 0.1);
}

.mobile-nav-item.active {
  color: rgb(var(--ocean-surface));
  background: rgba(56, 189, 248, 0.12);
  border-color: rgba(56, 189, 248, 0.15);
}

.mobile-register-btn {
  display: block;
  margin-top: 20px;
  text-decoration: none;
}

.mobile-register-btn .el-button {
  width: 100%;
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid))) !important;
  border: none !important;
  box-shadow: 0 4px 20px rgba(14, 165, 233, 0.3);
}

/* ===== Main Content - 全屏无边界 ===== */
.main-content {
  flex: 1;
  width: 100%;
  padding-top: 72px;
  position: relative;
  z-index: 1;
}

/* ===== Footer - 简化透明 ===== */
.main-footer {
  padding: 24px;
  position: relative;
  z-index: 1;
  background: transparent;
}

.footer-content {
  max-width: 1400px;
  margin: 0 auto;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
}

@media (min-width: 768px) {
  .footer-content {
    flex-direction: row;
    justify-content: space-between;
  }
}

.footer-brand {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-wrap: wrap;
  justify-content: center;
}

.footer-logo {
  display: flex;
  align-items: center;
  gap: 8px;
}

.footer-wave-svg {
  width: 24px;
  height: 24px;
}

.footer-logo-text {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
}

.footer-divider {
  color: rgba(56, 189, 248, 0.3);
}

.footer-tagline {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.55);
}

.footer-links {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 13px;
  color: rgba(255, 255, 255, 0.55);
}

.footer-link {
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  transition: color 0.2s;
}

.footer-link:hover {
  color: rgb(var(--ocean-surface));
}

.footer-dot {
  color: rgba(56, 189, 248, 0.3);
}

.footer-license,
.footer-copyright {
  color: rgba(255, 255, 255, 0.45);
}

/* Element Plus Overrides */
:deep(.el-divider) {
  border-color: rgba(56, 189, 248, 0.1);
  margin: 16px 0;
}
</style>
