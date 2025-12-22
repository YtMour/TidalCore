<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUsers, type UsersResponse } from '@/api/admin'
import { User, Calendar, TrendCharts, Trophy } from '@element-plus/icons-vue'

const loading = ref(true)
const stats = ref({
  totalUsers: 0,
  totalCheckins: 0,
  activeUsers: 0,
  avgStreak: 0
})

// 快捷入口
const quickLinks = [
  { name: '用户管理', path: '/admin/users', icon: User, color: 'ocean' },
  { name: '打卡数据', path: '/admin/checkins', icon: Calendar, color: 'green' },
  { name: '系统设置', path: '/admin/settings', icon: TrendCharts, color: 'amber' }
]

onMounted(async () => {
  await loadStats()
})

async function loadStats() {
  loading.value = true
  try {
    // 获取用户数据来计算统计
    const res: UsersResponse = await getUsers(1, 1000)
    const users = res.users

    stats.value.totalUsers = res.total

    // 计算统计数据
    let totalCheckins = 0
    let totalStreak = 0
    let activeCount = 0

    users.forEach(user => {
      totalCheckins += user.total_checkin || 0
      totalStreak += user.streak || 0
      if ((user.streak || 0) > 0) {
        activeCount++
      }
    })

    stats.value.totalCheckins = totalCheckins
    stats.value.activeUsers = activeCount
    stats.value.avgStreak = users.length > 0 ? Math.round(totalStreak / users.length * 10) / 10 : 0
  } catch (error) {
    console.error('加载统计数据失败', error)
  } finally {
    loading.value = false
  }
}
</script>

<template>
  <div class="dashboard-page">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card" v-loading="loading">
        <div class="stat-icon ocean">
          <el-icon :size="24"><User /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.totalUsers }}</span>
          <span class="stat-label">总用户数</span>
        </div>
      </div>

      <div class="stat-card" v-loading="loading">
        <div class="stat-icon green">
          <el-icon :size="24"><Calendar /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.totalCheckins }}</span>
          <span class="stat-label">累计打卡次数</span>
        </div>
      </div>

      <div class="stat-card" v-loading="loading">
        <div class="stat-icon pink">
          <el-icon :size="24"><TrendCharts /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.activeUsers }}</span>
          <span class="stat-label">活跃用户</span>
        </div>
      </div>

      <div class="stat-card" v-loading="loading">
        <div class="stat-icon amber">
          <el-icon :size="24"><Trophy /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.avgStreak }}</span>
          <span class="stat-label">平均连续天数</span>
        </div>
      </div>
    </div>

    <!-- 快捷入口 -->
    <div class="section">
      <h2 class="section-title">快捷入口</h2>
      <div class="quick-links">
        <RouterLink
          v-for="link in quickLinks"
          :key="link.path"
          :to="link.path"
          class="quick-link"
          :class="link.color"
        >
          <el-icon :size="28"><component :is="link.icon" /></el-icon>
          <span>{{ link.name }}</span>
        </RouterLink>
      </div>
    </div>

    <!-- 系统信息 -->
    <div class="section">
      <h2 class="section-title">系统信息</h2>
      <div class="info-card">
        <div class="info-row">
          <span class="info-label">系统版本</span>
          <span class="info-value">TidalCore v1.0.0</span>
        </div>
        <div class="info-row">
          <span class="info-label">前端框架</span>
          <span class="info-value">Vue 3 + TypeScript</span>
        </div>
        <div class="info-row">
          <span class="info-label">后端框架</span>
          <span class="info-value">Go + Gin</span>
        </div>
        <div class="info-row">
          <span class="info-label">数据库</span>
          <span class="info-value">SQLite / PostgreSQL</span>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.dashboard-page {
  max-width: 1200px;
  margin: 0 auto;
}

/* ===== 统计卡片 ===== */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
  gap: 20px;
  margin-bottom: 32px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 24px;
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: 16px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  transform: translateY(-2px);
  border-color: rgba(56, 189, 248, 0.2);
  box-shadow: 0 8px 32px rgba(56, 189, 248, 0.1);
}

.stat-icon {
  width: 56px;
  height: 56px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 14px;
  flex-shrink: 0;
}

.stat-icon.ocean {
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.2), rgba(56, 189, 248, 0.1));
  color: rgb(var(--ocean-surface));
}

.stat-icon.green {
  background: linear-gradient(135deg, rgba(52, 211, 153, 0.2), rgba(52, 211, 153, 0.1));
  color: rgb(var(--seaweed-green));
}

.stat-icon.pink {
  background: linear-gradient(135deg, rgba(244, 114, 182, 0.2), rgba(244, 114, 182, 0.1));
  color: rgb(var(--coral-pink));
}

.stat-icon.amber {
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.2), rgba(251, 191, 36, 0.1));
  color: rgb(var(--sunset-amber));
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
}

.stat-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

/* ===== Section ===== */
.section {
  margin-bottom: 32px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 16px;
}

/* ===== 快捷入口 ===== */
.quick-links {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
  gap: 16px;
}

.quick-link {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 12px;
  padding: 24px;
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: 16px;
  text-decoration: none;
  transition: all 0.3s ease;
}

.quick-link:hover {
  transform: translateY(-2px);
  border-color: rgba(56, 189, 248, 0.3);
}

.quick-link.ocean {
  color: rgb(var(--ocean-surface));
}

.quick-link.ocean:hover {
  background: rgba(56, 189, 248, 0.1);
}

.quick-link.green {
  color: rgb(var(--seaweed-green));
}

.quick-link.green:hover {
  background: rgba(52, 211, 153, 0.1);
}

.quick-link.amber {
  color: rgb(var(--sunset-amber));
}

.quick-link.amber:hover {
  background: rgba(251, 191, 36, 0.1);
}

.quick-link span {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
}

/* ===== 信息卡片 ===== */
.info-card {
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: 16px;
  padding: 8px 0;
}

.info-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 20px;
  border-bottom: 1px solid rgba(56, 189, 248, 0.08);
}

.info-row:last-child {
  border-bottom: none;
}

.info-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
}

.info-value {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.9);
  font-weight: 500;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
    gap: 12px;
  }

  .stat-card {
    padding: 16px;
    flex-direction: column;
    text-align: center;
  }

  .stat-icon {
    width: 48px;
    height: 48px;
  }

  .stat-value {
    font-size: 24px;
  }

  .quick-links {
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
  }

  .quick-link {
    padding: 16px;
  }
}
</style>
