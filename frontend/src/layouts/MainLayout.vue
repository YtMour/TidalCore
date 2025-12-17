<script setup lang="ts">
import { ref, watch, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useScrollLock } from '@vueuse/core'
import { HomeFilled, Timer, Trophy, User, SwitchButton, Fold, Expand } from '@element-plus/icons-vue'

const userStore = useUserStore()
const route = useRoute()

const mobileMenuOpen = ref(false)
const scrollLock = useScrollLock(document.body)
const scrolled = ref(false)

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
    <header class="main-header" :class="{ scrolled }">
      <div class="header-content">
        <!-- Logo -->
        <RouterLink to="/" class="logo-link">
          <span class="logo-icon">🌊</span>
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

    <!-- Main Content -->
    <main class="main-content">
      <slot />
    </main>

    <!-- Footer -->
    <footer class="main-footer">
      <div class="footer-content">
        <div class="footer-top">
          <div class="footer-brand">
            <div class="footer-logo">
              <span class="logo-icon">🌊</span>
              <span class="logo-text">TidalCore</span>
            </div>
            <p class="footer-tagline">开源盆底肌训练平台</p>
          </div>

          <div class="footer-links">
            <a href="https://github.com" target="_blank" rel="noopener noreferrer" class="footer-link">
              GitHub
            </a>
            <span class="footer-divider">·</span>
            <span class="footer-license">MIT License</span>
          </div>
        </div>

        <div class="footer-bottom">
          <p>Made with <span class="heart">♥</span> for health</p>
          <p>© 2024 TidalCore</p>
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

/* ===== Header ===== */
.main-header {
  position: sticky;
  top: 0;
  z-index: 100;
  height: 60px;
  padding: 0 20px;
  background: transparent;
  border-bottom: 1px solid transparent;
  transition: all 0.3s ease;
}

.main-header.scrolled {
  background: rgba(15, 23, 42, 0.95);
  backdrop-filter: blur(12px);
  border-bottom-color: rgba(255, 255, 255, 0.06);
}

.header-content {
  max-width: 1200px;
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
  gap: 8px;
  text-decoration: none;
  flex-shrink: 0;
}

.logo-icon {
  font-size: 24px;
  transition: transform 0.3s ease;
}

.logo-link:hover .logo-icon {
  transform: scale(1.1) rotate(-5deg);
}

.logo-text {
  font-size: 18px;
  font-weight: 700;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

/* ===== Desktop Navigation ===== */
.desktop-nav {
  display: none;
  align-items: center;
  gap: 4px;
  padding: 4px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 8px;
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
  padding: 8px 14px;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  transition: all 0.2s ease;
}

.nav-link:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.08);
}

.nav-link.active {
  color: #fff;
  background: rgba(139, 92, 246, 0.2);
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

/* User Link */
.user-link {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
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
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.1);
}

/* Logout Button */
.logout-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 14px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.6);
  background: transparent;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
}

.logout-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.05);
}

/* Login Button */
.login-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border-radius: 8px;
  font-size: 14px;
  font-weight: 600;
  color: #fff;
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  text-decoration: none;
  transition: all 0.2s ease;
}

.login-btn:hover {
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(139, 92, 246, 0.4);
}

/* ===== Mobile Menu Button ===== */
.mobile-menu-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 8px;
  border: none;
  background: rgba(255, 255, 255, 0.03);
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.2s ease;
}

.mobile-menu-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.08);
}

@media (min-width: 768px) {
  .mobile-menu-btn {
    display: none;
  }
}

/* ===== Mobile Drawer ===== */
.mobile-drawer :deep(.el-drawer__body) {
  background: rgba(15, 23, 42, 0.98);
  backdrop-filter: blur(20px);
  padding-top: 80px;
}

.mobile-menu {
  padding: 0 20px;
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

/* ===== Main Content ===== */
.main-content {
  flex: 1;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  padding: 24px 20px;
}

@media (min-width: 768px) {
  .main-content {
    padding: 40px 24px;
  }
}

/* ===== Footer ===== */
.main-footer {
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  padding: 32px 20px;
}

.footer-content {
  max-width: 1200px;
  margin: 0 auto;
}

.footer-top {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 16px;
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
  gap: 6px;
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
}

.footer-logo .logo-icon {
  font-size: 20px;
}

.footer-logo .logo-text {
  font-size: 16px;
}

.footer-tagline {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.4);
  margin: 0;
}

.footer-links {
  display: flex;
  align-items: center;
  gap: 12px;
  font-size: 13px;
}

.footer-link {
  color: rgba(255, 255, 255, 0.4);
  text-decoration: none;
  transition: color 0.2s ease;
}

.footer-link:hover {
  color: rgba(255, 255, 255, 0.7);
}

.footer-divider {
  color: rgba(255, 255, 255, 0.2);
}

.footer-license {
  color: rgba(255, 255, 255, 0.4);
}

.footer-bottom {
  margin-top: 24px;
  padding-top: 20px;
  border-top: 1px solid rgba(255, 255, 255, 0.05);
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
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
  border-color: rgba(255, 255, 255, 0.08);
  margin: 12px 0;
}
</style>
