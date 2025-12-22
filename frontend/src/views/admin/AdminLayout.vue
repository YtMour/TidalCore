<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import AdminSidebar from '@/components/admin/AdminSidebar.vue'
import { ElMessage } from 'element-plus'
import { Fold, Expand } from '@element-plus/icons-vue'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const sidebarCollapsed = ref(false)
const isMobile = ref(false)
const showMobileSidebar = ref(false)

// 检查屏幕尺寸
function checkScreenSize() {
  isMobile.value = window.innerWidth < 768
  if (isMobile.value) {
    sidebarCollapsed.value = true
  }
}

onMounted(async () => {
  checkScreenSize()
  window.addEventListener('resize', checkScreenSize)

  await userStore.fetchProfile()
  if (!userStore.user?.is_admin) {
    ElMessage.error('无权访问管理后台')
    router.push('/dashboard')
  }
})

// 当前页面标题
const pageTitle = computed(() => {
  const titles: Record<string, string> = {
    'admin-dashboard': '仪表盘',
    'admin-users': '用户管理',
    'admin-checkins': '打卡数据',
    'admin-settings': '系统设置'
  }
  return titles[route.name as string] || '管理后台'
})

function toggleSidebar() {
  if (isMobile.value) {
    showMobileSidebar.value = !showMobileSidebar.value
  } else {
    sidebarCollapsed.value = !sidebarCollapsed.value
  }
}

function closeMobileSidebar() {
  showMobileSidebar.value = false
}
</script>

<template>
  <div class="admin-layout" :class="{ 'sidebar-collapsed': sidebarCollapsed }">
    <!-- 海洋背景系统 -->
    <div class="ocean-bg"></div>
    <div class="ocean-gradient"></div>
    <div class="animated-bg"></div>
    <div class="particles-bg"></div>
    <div class="noise-overlay"></div>
    <div class="wave-layer">
      <div class="wave wave-1"></div>
      <div class="wave wave-2"></div>
      <div class="wave wave-3"></div>
    </div>

    <!-- 移动端遮罩 -->
    <div
      v-if="isMobile && showMobileSidebar"
      class="mobile-overlay"
      @click="closeMobileSidebar"
    />

    <!-- 侧边栏 -->
    <AdminSidebar
      :collapsed="sidebarCollapsed && !isMobile"
      :class="{ 'mobile-visible': showMobileSidebar }"
      @close="closeMobileSidebar"
    />

    <!-- 主内容区 -->
    <div class="admin-main">
      <!-- 顶部栏 -->
      <header class="admin-header">
        <div class="header-left">
          <button class="toggle-btn" @click="toggleSidebar">
            <el-icon :size="20">
              <Fold v-if="!sidebarCollapsed || showMobileSidebar" />
              <Expand v-else />
            </el-icon>
          </button>
          <h1 class="page-title">{{ pageTitle }}</h1>
        </div>
        <div class="header-right">
          <RouterLink to="/dashboard" class="back-link">
            返回用户中心
          </RouterLink>
          <div class="admin-user">
            <span class="admin-name">{{ userStore.user?.display_name || userStore.user?.username }}</span>
            <span class="admin-badge">管理员</span>
          </div>
        </div>
      </header>

      <!-- 内容区域 -->
      <main class="admin-content">
        <router-view v-slot="{ Component }">
          <transition name="fade-slide" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </main>
    </div>
  </div>
</template>

<style scoped>
.admin-layout {
  display: flex;
  min-height: 100vh;
  position: relative;
}

/* ===== 侧边栏过渡 ===== */
.admin-layout.sidebar-collapsed .admin-main {
  margin-left: 72px;
}

.admin-main {
  flex: 1;
  margin-left: 260px;
  transition: margin-left 0.3s ease;
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  position: relative;
  z-index: 1;
}

/* ===== 顶部栏 ===== */
.admin-header {
  height: 64px;
  padding: 0 24px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
  -webkit-backdrop-filter: blur(20px);
  border-bottom: 1px solid rgba(56, 189, 248, 0.1);
  position: sticky;
  top: 0;
  z-index: 100;
}

.header-left {
  display: flex;
  align-items: center;
  gap: 16px;
}

.toggle-btn {
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: rgba(56, 189, 248, 0.1);
  border: 1px solid rgba(56, 189, 248, 0.15);
  border-radius: 10px;
  color: rgba(255, 255, 255, 0.7);
  cursor: pointer;
  transition: all 0.3s var(--ease-smooth);
}

.toggle-btn:hover {
  background: rgba(56, 189, 248, 0.2);
  border-color: rgba(56, 189, 248, 0.3);
  color: rgb(var(--ocean-surface));
  transform: translateY(-1px);
}

.page-title {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.back-link {
  color: rgba(255, 255, 255, 0.6);
  text-decoration: none;
  font-size: 14px;
  transition: all 0.2s;
  padding: 6px 12px;
  border-radius: 8px;
}

.back-link:hover {
  color: rgb(var(--ocean-surface));
  background: rgba(56, 189, 248, 0.1);
}

.admin-user {
  display: flex;
  align-items: center;
  gap: 10px;
}

.admin-name {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.8);
}

.admin-badge {
  font-size: 12px;
  padding: 4px 10px;
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.2), rgba(251, 191, 36, 0.1));
  border: 1px solid rgba(251, 191, 36, 0.3);
  border-radius: 6px;
  color: rgb(var(--sunset-amber));
  font-weight: 500;
}

/* ===== 内容区域 ===== */
.admin-content {
  flex: 1;
  padding: 24px;
  overflow-y: auto;
}

/* ===== 页面切换动画 ===== */
.fade-slide-enter-active,
.fade-slide-leave-active {
  transition: all 0.3s var(--ease-smooth);
}

.fade-slide-enter-from {
  opacity: 0;
  transform: translateY(12px);
}

.fade-slide-leave-to {
  opacity: 0;
  transform: translateY(-12px);
}

/* ===== 移动端遮罩 ===== */
.mobile-overlay {
  position: fixed;
  inset: 0;
  background: rgba(5, 8, 15, 0.7);
  z-index: 199;
  backdrop-filter: blur(4px);
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .admin-main {
    margin-left: 0;
  }

  .admin-layout.sidebar-collapsed .admin-main {
    margin-left: 0;
  }

  .admin-header {
    padding: 0 16px;
  }

  .admin-content {
    padding: 16px;
  }

  .back-link {
    display: none;
  }

  .admin-name {
    display: none;
  }
}
</style>
