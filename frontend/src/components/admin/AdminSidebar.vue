<script setup lang="ts">
import { useRoute } from 'vue-router'
import {
  DataAnalysis,
  User,
  Calendar,
  Setting,
  Close,
  Back,
  FolderChecked
} from '@element-plus/icons-vue'

defineProps<{
  collapsed: boolean
}>()

const emit = defineEmits<{
  close: []
}>()

const route = useRoute()

const menuItems = [
  { name: 'admin-dashboard', path: '/admin/dashboard', icon: DataAnalysis, label: '仪表盘', desc: '数据概览' },
  { name: 'admin-users', path: '/admin/users', icon: User, label: '用户管理', desc: '管理用户' },
  { name: 'admin-checkins', path: '/admin/checkins', icon: Calendar, label: '打卡数据', desc: '热力图统计' },
  { name: 'admin-backup', path: '/admin/backup', icon: FolderChecked, label: '数据备份', desc: '备份恢复' },
  { name: 'admin-settings', path: '/admin/settings', icon: Setting, label: '系统设置', desc: '配置信息' }
]

const isActive = (name: string) => route.name === name

function handleNavClick() {
  emit('close')
}
</script>

<template>
  <aside class="admin-sidebar" :class="{ collapsed }">
    <!-- Logo 区域 -->
    <div class="sidebar-header">
      <RouterLink to="/" class="logo">
        <div class="logo-icon-wrap">
          <svg class="logo-wave" viewBox="0 0 36 36" fill="none">
            <defs>
              <linearGradient id="logoGrad" x1="0%" y1="100%" x2="100%" y2="0%">
                <stop offset="0%" stop-color="#0ea5e9" />
                <stop offset="100%" stop-color="#38bdf8" />
              </linearGradient>
            </defs>
            <circle cx="18" cy="18" r="16" fill="rgba(56, 189, 248, 0.15)" />
            <path
              d="M6 20 Q12 16, 18 20 T30 20"
              stroke="url(#logoGrad)"
              stroke-width="2.5"
              stroke-linecap="round"
              fill="none"
              class="wave-path"
            />
            <path
              d="M6 24 Q12 20, 18 24 T30 24"
              stroke="url(#logoGrad)"
              stroke-width="2"
              stroke-linecap="round"
              fill="none"
              opacity="0.5"
              class="wave-path-2"
            />
          </svg>
        </div>
        <transition name="fade">
          <div v-if="!collapsed" class="logo-text-wrap">
            <span class="logo-text">TidalCore</span>
            <span class="logo-sub">管理后台</span>
          </div>
        </transition>
      </RouterLink>
      <button class="close-btn" @click="$emit('close')">
        <el-icon><Close /></el-icon>
      </button>
    </div>

    <!-- 导航菜单 -->
    <nav class="sidebar-nav">
      <div class="nav-section-label" v-if="!collapsed">导航菜单</div>

      <RouterLink
        v-for="item in menuItems"
        :key="item.name"
        :to="item.path"
        class="nav-item"
        :class="{ active: isActive(item.name) }"
        @click="handleNavClick"
      >
        <div class="nav-icon-wrap">
          <el-icon class="nav-icon" :size="20">
            <component :is="item.icon" />
          </el-icon>
        </div>
        <transition name="fade">
          <div v-if="!collapsed" class="nav-text">
            <span class="nav-label">{{ item.label }}</span>
            <span class="nav-desc">{{ item.desc }}</span>
          </div>
        </transition>
        <div v-if="isActive(item.name)" class="active-glow" />
      </RouterLink>
    </nav>

    <!-- 底部区域 -->
    <div class="sidebar-footer">
      <RouterLink to="/dashboard" class="back-btn">
        <el-icon :size="18"><Back /></el-icon>
        <transition name="fade">
          <span v-if="!collapsed">返回用户中心</span>
        </transition>
      </RouterLink>

      <transition name="fade">
        <div v-if="!collapsed" class="footer-info">
          <div class="footer-divider"></div>
          <div class="footer-brand">
            <span class="brand-icon">🌊</span>
            <span class="brand-text">TidalCore v1.0</span>
          </div>
        </div>
      </transition>
    </div>
  </aside>
</template>

<style scoped>
.admin-sidebar {
  width: 260px;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
  background: linear-gradient(
    180deg,
    rgba(8, 12, 24, 0.95) 0%,
    rgba(12, 18, 32, 0.92) 50%,
    rgba(8, 12, 24, 0.95) 100%
  );
  backdrop-filter: blur(24px);
  -webkit-backdrop-filter: blur(24px);
  border-right: 1px solid rgba(56, 189, 248, 0.08);
  display: flex;
  flex-direction: column;
  transition: width 0.3s var(--ease-smooth, cubic-bezier(0.4, 0, 0.2, 1));
  z-index: 200;
  overflow: hidden;
}

.admin-sidebar::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 200px;
  background: radial-gradient(
    ellipse 80% 50% at 50% 0%,
    rgba(56, 189, 248, 0.08) 0%,
    transparent 70%
  );
  pointer-events: none;
}

.admin-sidebar.collapsed {
  width: 72px;
}

/* ===== Header ===== */
.sidebar-header {
  height: 72px;
  padding: 0 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  position: relative;
  z-index: 1;
}

.logo {
  display: flex;
  align-items: center;
  gap: 12px;
  text-decoration: none;
  overflow: hidden;
}

.logo-icon-wrap {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  position: relative;
}

.logo-wave {
  width: 36px;
  height: 36px;
}

.wave-path {
  animation: wave-float 3s ease-in-out infinite;
}

.wave-path-2 {
  animation: wave-float 3s ease-in-out infinite 0.5s;
}

@keyframes wave-float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-2px); }
}

.logo-text-wrap {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.logo-text {
  font-size: 17px;
  font-weight: 700;
  white-space: nowrap;
  background: linear-gradient(135deg, #38bdf8, #22d3ee);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.logo-sub {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
  font-weight: 500;
}

.close-btn {
  display: none;
  width: 32px;
  height: 32px;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.6);
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s;
}

.close-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

/* ===== Navigation ===== */
.sidebar-nav {
  flex: 1;
  padding: 8px 12px;
  display: flex;
  flex-direction: column;
  gap: 4px;
  overflow-y: auto;
  position: relative;
  z-index: 1;
}

.nav-section-label {
  font-size: 11px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 1.2px;
  color: rgba(255, 255, 255, 0.3);
  padding: 12px 12px 8px;
}

.nav-item {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px;
  border-radius: 12px;
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  transition: all 0.25s var(--ease-smooth, cubic-bezier(0.4, 0, 0.2, 1));
  position: relative;
  overflow: hidden;
}

.collapsed .nav-item {
  justify-content: center;
  padding: 12px;
}

.nav-icon-wrap {
  width: 40px;
  height: 40px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  flex-shrink: 0;
  transition: all 0.25s var(--ease-smooth, cubic-bezier(0.4, 0, 0.2, 1));
}

.nav-item:hover {
  color: rgba(255, 255, 255, 0.95);
  background: rgba(56, 189, 248, 0.05);
}

.nav-item:hover .nav-icon-wrap {
  background: rgba(56, 189, 248, 0.1);
  border-color: rgba(56, 189, 248, 0.2);
}

.nav-item.active {
  color: #fff;
  background: linear-gradient(
    135deg,
    rgba(56, 189, 248, 0.12) 0%,
    rgba(34, 211, 238, 0.08) 100%
  );
}

.nav-item.active .nav-icon-wrap {
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.3), rgba(34, 211, 238, 0.2));
  border-color: rgba(56, 189, 248, 0.4);
  color: #38bdf8;
  box-shadow: 0 0 20px rgba(56, 189, 248, 0.25);
}

.nav-icon {
  flex-shrink: 0;
  transition: transform 0.2s;
}

.nav-item:hover .nav-icon {
  transform: scale(1.1);
}

.nav-text {
  display: flex;
  flex-direction: column;
  gap: 2px;
  min-width: 0;
}

.nav-label {
  white-space: nowrap;
  font-size: 14px;
  font-weight: 500;
}

.nav-desc {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.35);
  white-space: nowrap;
}

.nav-item.active .nav-desc {
  color: rgba(56, 189, 248, 0.6);
}

.active-glow {
  position: absolute;
  inset: 0;
  background: radial-gradient(
    ellipse 100% 100% at 0% 50%,
    rgba(56, 189, 248, 0.15) 0%,
    transparent 70%
  );
  pointer-events: none;
}

/* ===== Footer ===== */
.sidebar-footer {
  padding: 12px;
  position: relative;
  z-index: 1;
}

.back-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  width: 100%;
  padding: 10px 12px;
  background: rgba(255, 255, 255, 0.03);
  border: 1px solid rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  color: rgba(255, 255, 255, 0.5);
  font-size: 13px;
  text-decoration: none;
  transition: all 0.2s;
}

.back-btn:hover {
  background: rgba(255, 255, 255, 0.06);
  border-color: rgba(255, 255, 255, 0.1);
  color: rgba(255, 255, 255, 0.8);
}

.collapsed .back-btn {
  padding: 10px;
}

.footer-info {
  margin-top: 12px;
}

.footer-divider {
  height: 1px;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(56, 189, 248, 0.15) 50%,
    transparent 100%
  );
  margin-bottom: 12px;
}

.footer-brand {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 6px;
}

.brand-icon {
  font-size: 14px;
}

.brand-text {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.3);
}

/* ===== Transitions ===== */
.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}

/* ===== Mobile ===== */
@media (max-width: 768px) {
  .admin-sidebar {
    transform: translateX(-100%);
    width: 280px;
  }

  .admin-sidebar.mobile-visible {
    transform: translateX(0);
  }

  .admin-sidebar.collapsed {
    width: 280px;
  }

  .close-btn {
    display: flex;
  }

  .collapsed .nav-item {
    justify-content: flex-start;
    padding: 12px;
  }

  .collapsed .back-btn {
    justify-content: flex-start;
    padding: 10px 12px;
  }
}
</style>
