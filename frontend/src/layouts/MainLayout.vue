<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useThemeStore } from '@/store/theme'
import { useScrollLock } from '@vueuse/core'
import { HomeFilled, Timer, Trophy, User, SwitchButton, Fold, Expand, Sunny, Moon } from '@element-plus/icons-vue'

const userStore = useUserStore()
const themeStore = useThemeStore()
const route = useRoute()

const mobileMenuOpen = ref(false)
const scrollLock = useScrollLock(document.body)
const scrolled = ref(false)

const isDark = computed(() => themeStore.mode === 'dark')

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

function toggleTheme() {
  themeStore.toggle()
}

const navLinks = [
  { to: '/', label: '首页', icon: HomeFilled },
  { to: '/train', label: '训练', icon: Timer },
  { to: '/leaderboard', label: '排行榜', icon: Trophy }
]
</script>

<template>
  <div class="main-layout">
    <!-- Animated Background -->
    <div class="animated-bg"></div>
    <div class="particles-bg"></div>
    <div class="noise-overlay"></div>

    <!-- Navigation -->
    <el-header
      class="main-header"
      :class="{ scrolled }"
    >
      <div class="header-content">
        <!-- Logo -->
        <RouterLink to="/" class="logo-link">
          <span class="logo-icon">🌊</span>
          <span class="logo-text">TidalCore</span>
        </RouterLink>

        <!-- Desktop Navigation -->
        <div class="desktop-nav">
          <el-menu
            :default-active="route.path"
            mode="horizontal"
            :ellipsis="false"
            class="nav-menu"
            router
          >
            <el-menu-item
              v-for="link in navLinks"
              :key="link.to"
              :index="link.to"
            >
              <el-icon><component :is="link.icon" /></el-icon>
              <span>{{ link.label }}</span>
            </el-menu-item>
          </el-menu>
        </div>

        <!-- Desktop User Menu -->
        <div class="desktop-user">
          <!-- Theme Toggle -->
          <el-button text @click="toggleTheme" class="theme-btn">
            <el-icon :size="18">
              <Moon v-if="isDark" />
              <Sunny v-else />
            </el-icon>
          </el-button>

          <template v-if="userStore.isLoggedIn">
            <RouterLink to="/dashboard" class="user-link" :class="{ active: route.path === '/dashboard' }">
              <el-icon><User /></el-icon>
              <span>我的</span>
            </RouterLink>
            <el-button text @click="handleLogout" class="logout-btn">
              <el-icon><SwitchButton /></el-icon>
              <span>退出</span>
            </el-button>
          </template>
          <template v-else>
            <RouterLink to="/login">
              <el-button type="primary" class="login-btn">
                <el-icon><User /></el-icon>
                <span>登录</span>
              </el-button>
            </RouterLink>
          </template>
        </div>

        <!-- Mobile Menu Button -->
        <el-button
          class="mobile-menu-btn"
          text
          @click="mobileMenuOpen = !mobileMenuOpen"
        >
          <el-icon :size="24">
            <Fold v-if="mobileMenuOpen" />
            <Expand v-else />
          </el-icon>
        </el-button>
      </div>
    </el-header>

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

        <!-- Theme Toggle in Mobile -->
        <div class="mobile-nav-item" @click="toggleTheme">
          <el-icon :size="20">
            <Moon v-if="isDark" />
            <Sunny v-else />
          </el-icon>
          <span>{{ isDark ? '浅色模式' : '深色模式' }}</span>
        </div>

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
          <div
            class="mobile-nav-item"
            @click="handleLogout"
          >
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

    <!-- Main Content -->
    <el-main class="main-content">
      <slot />
    </el-main>

    <!-- Footer -->
    <el-footer class="main-footer">
      <div class="footer-bg">
        <div class="footer-blob footer-blob-1"></div>
        <div class="footer-blob footer-blob-2"></div>
      </div>

      <div class="footer-content">
        <div class="footer-top">
          <!-- Logo and tagline -->
          <div class="footer-brand">
            <div class="footer-logo">
              <span class="logo-icon">🌊</span>
              <span class="logo-text">TidalCore</span>
            </div>
            <p class="footer-tagline">开源盆底肌训练平台</p>
          </div>

          <!-- Links -->
          <div class="footer-links">
            <a
              href="https://github.com"
              target="_blank"
              rel="noopener noreferrer"
              class="footer-link"
            >
              <el-icon><Link /></el-icon>
              <span>GitHub</span>
            </a>
            <el-divider direction="vertical" />
            <span class="footer-license">MIT License</span>
          </div>
        </div>

        <!-- Bottom bar -->
        <div class="footer-bottom">
          <p>Made with <span class="heart">♥</span> for health</p>
          <p>© 2024 TidalCore. All rights reserved.</p>
        </div>
      </div>
    </el-footer>
  </div>
</template>

<style scoped>
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  position: relative;
}

/* Header Styles */
.main-header {
  position: sticky;
  top: 0;
  z-index: 100;
  height: 64px;
  padding: 0 24px;
  background: transparent;
  border-bottom: 1px solid transparent;
  transition: all 0.3s ease;
}

.main-header.scrolled {
  background: rgba(15, 23, 42, 0.9);
  backdrop-filter: blur(20px);
  border-bottom-color: rgba(255, 255, 255, 0.05);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.header-content {
  max-width: 1200px;
  margin: 0 auto;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: space-between;
}

/* Logo */
.logo-link {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
}

.logo-icon {
  font-size: 28px;
  transition: transform 0.3s ease;
}

.logo-link:hover .logo-icon {
  transform: scale(1.1);
}

.logo-text {
  font-size: 20px;
  font-weight: 700;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* Desktop Navigation */
.desktop-nav {
  display: none;
}

@media (min-width: 768px) {
  .desktop-nav {
    display: block;
  }
}

.nav-menu {
  background: rgba(255, 255, 255, 0.02) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 12px;
  padding: 4px;
  border-bottom: none !important;
}

.nav-menu :deep(.el-menu-item) {
  background: transparent !important;
  border-radius: 8px;
  margin: 0 2px;
  padding: 0 16px;
  height: 36px;
  line-height: 36px;
  color: rgba(255, 255, 255, 0.6);
  border-bottom: none !important;
  transition: all 0.2s ease;
}

.nav-menu :deep(.el-menu-item:hover) {
  background: rgba(255, 255, 255, 0.05) !important;
  color: #fff;
}

.nav-menu :deep(.el-menu-item.is-active) {
  background: rgba(255, 255, 255, 0.1) !important;
  color: #fff;
}

.nav-menu :deep(.el-menu-item .el-icon) {
  margin-right: 6px;
}

/* Desktop User Menu */
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
  padding: 8px 16px;
  border-radius: 8px;
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  border: 1px solid transparent;
  transition: all 0.2s ease;
}

.user-link:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.05);
  border-color: rgba(255, 255, 255, 0.08);
}

.user-link.active {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
  border-color: rgba(255, 255, 255, 0.08);
}

.logout-btn {
  color: rgba(255, 255, 255, 0.6) !important;
}

.logout-btn:hover {
  color: #fff !important;
  background: rgba(255, 255, 255, 0.05) !important;
}

.login-btn {
  background: linear-gradient(135deg, #8b5cf6, #ec4899) !important;
  border: none !important;
  border-radius: 6px;
}

/* Theme Toggle Button */
.theme-btn {
  width: 36px;
  height: 36px;
  border-radius: 6px !important;
  color: rgba(255, 255, 255, 0.6) !important;
  transition: all 0.2s ease !important;
}

.theme-btn:hover {
  color: #fbbf24 !important;
  background: rgba(251, 191, 36, 0.1) !important;
}

/* Light mode adjustments */
html.light .theme-btn {
  color: rgba(15, 23, 42, 0.6) !important;
}

html.light .theme-btn:hover {
  color: #f59e0b !important;
  background: rgba(245, 158, 11, 0.1) !important;
}

/* Mobile Menu Button */
.mobile-menu-btn {
  display: flex;
  color: rgba(255, 255, 255, 0.7) !important;
}

.mobile-menu-btn:hover {
  color: #fff !important;
  background: rgba(255, 255, 255, 0.05) !important;
}

@media (min-width: 768px) {
  .mobile-menu-btn {
    display: none;
  }
}

/* Mobile Drawer */
.mobile-drawer :deep(.el-drawer__body) {
  background: rgba(15, 23, 42, 0.98);
  backdrop-filter: blur(20px);
  padding-top: 80px;
}

.mobile-menu {
  padding: 0 24px;
}

.mobile-nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 14px 16px;
  border-radius: 10px;
  font-size: 16px;
  color: rgba(255, 255, 255, 0.7);
  text-decoration: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.mobile-nav-item:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.05);
}

.mobile-nav-item.active {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}

.mobile-register-btn {
  display: block;
  margin-top: 16px;
  text-decoration: none;
}

.mobile-register-btn .el-button {
  width: 100%;
  background: linear-gradient(135deg, #8b5cf6, #ec4899) !important;
  border: none !important;
}

/* Main Content */
.main-content {
  flex: 1;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  padding: 32px 24px;
}

@media (min-width: 768px) {
  .main-content {
    padding: 48px 24px;
  }
}

/* Footer */
.main-footer {
  position: relative;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  padding: 40px 24px;
  height: auto;
  overflow: hidden;
}

.footer-bg {
  position: absolute;
  inset: 0;
  z-index: -1;
  pointer-events: none;
}

.footer-blob {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
}

.footer-blob-1 {
  bottom: 0;
  left: 25%;
  width: 256px;
  height: 128px;
  background: linear-gradient(to top, rgba(139, 92, 246, 0.05), transparent);
}

.footer-blob-2 {
  bottom: 0;
  right: 25%;
  width: 192px;
  height: 96px;
  background: linear-gradient(to top, rgba(236, 72, 153, 0.05), transparent);
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-top {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

@media (min-width: 768px) {
  .footer-top {
    flex-direction: row;
    justify-content: space-between;
  }
}

.footer-brand {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

@media (min-width: 768px) {
  .footer-brand {
    align-items: flex-start;
  }
}

.footer-logo {
  display: flex;
  align-items: center;
  gap: 8px;
  color: rgba(255, 255, 255, 0.6);
}

.footer-logo .logo-text {
  font-size: 16px;
  font-weight: 600;
}

.footer-tagline {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.4);
  margin: 0;
}

.footer-links {
  display: flex;
  align-items: center;
  gap: 16px;
  font-size: 14px;
}

.footer-link {
  display: flex;
  align-items: center;
  gap: 6px;
  color: rgba(255, 255, 255, 0.4);
  text-decoration: none;
  transition: color 0.2s ease;
}

.footer-link:hover {
  color: rgba(255, 255, 255, 0.7);
}

.footer-license {
  color: rgba(255, 255, 255, 0.4);
}

.footer-bottom {
  margin-top: 32px;
  padding-top: 24px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.3);
}

@media (min-width: 768px) {
  .footer-bottom {
    flex-direction: row;
    justify-content: space-between;
  }
}

.footer-bottom .heart {
  color: #ec4899;
}

/* Element Plus Overrides */
:deep(.el-divider) {
  border-color: rgba(255, 255, 255, 0.1);
  margin: 16px 0;
}

:deep(.el-divider--vertical) {
  border-color: rgba(255, 255, 255, 0.2);
}
</style>
