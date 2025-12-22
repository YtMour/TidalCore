<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getUsers, type UsersResponse } from '@/api/admin'
import { getVisitStats, type VisitStats } from '@/api/checkin'
import { InfoFilled } from '@element-plus/icons-vue'

const loading = ref(true)

// 系统统计
const systemStats = ref({
  totalUsers: 0,
  totalAdmins: 0,
  totalCheckins: 0
})

// 访问统计
const visitStats = ref<VisitStats>({
  today_visits: 0,
  total_visits: 0
})

// 称号配置
const titleConfig = [
  { min: 1000, name: '海神降临', color: '#f472b6', icon: '🔱', description: '累计打卡1000次以上' },
  { min: 730, name: '深渊霸主', color: '#a78bfa', icon: '🦑', description: '累计打卡730次以上' },
  { min: 365, name: '深海传奇', color: '#fbbf24', icon: '🌊', description: '累计打卡365次以上' },
  { min: 180, name: '海洋大师', color: '#38bdf8', icon: '🐋', description: '累计打卡180次以上' },
  { min: 90, name: '浪潮专家', color: '#22d3ee', icon: '🐬', description: '累计打卡90次以上' },
  { min: 30, name: '潮汐进阶', color: '#34d399', icon: '🐠', description: '累计打卡30次以上' },
  { min: 7, name: '入海新手', color: '#0ea5e9', icon: '🐟', description: '累计打卡7次以上' },
  { min: 0, name: '初探海域', color: 'rgba(255, 255, 255, 0.6)', icon: '🐚', description: '初始称号' }
]

onMounted(async () => {
  await Promise.all([
    loadSystemStats(),
    loadVisitStats()
  ])
  loading.value = false
})

async function loadSystemStats() {
  try {
    const res: UsersResponse = await getUsers(1, 1000)
    systemStats.value.totalUsers = res.total

    let adminCount = 0
    let totalCheckins = 0
    res.users.forEach(user => {
      if (user.is_admin) adminCount++
      totalCheckins += user.total_checkin || 0
    })

    systemStats.value.totalAdmins = adminCount
    systemStats.value.totalCheckins = totalCheckins
  } catch (error) {
    console.error('加载系统统计失败', error)
  }
}

async function loadVisitStats() {
  try {
    visitStats.value = await getVisitStats()
  } catch (error) {
    console.error('加载访问统计失败', error)
  }
}
</script>

<template>
  <div class="settings-page">
    <!-- 系统信息 -->
    <div class="section">
      <h2 class="section-title">系统信息</h2>
      <div class="info-card" v-loading="loading">
        <div class="info-row">
          <span class="info-label">系统名称</span>
          <span class="info-value">TidalCore</span>
        </div>
        <div class="info-row">
          <span class="info-label">系统版本</span>
          <span class="info-value">v1.0.0</span>
        </div>
        <div class="info-row">
          <span class="info-label">前端技术栈</span>
          <span class="info-value">Vue 3 + TypeScript + Element Plus</span>
        </div>
        <div class="info-row">
          <span class="info-label">后端技术栈</span>
          <span class="info-value">Go + Gin + GORM</span>
        </div>
        <div class="info-row">
          <span class="info-label">数据库</span>
          <span class="info-value">SQLite / PostgreSQL</span>
        </div>
      </div>
    </div>

    <!-- 数据统计 -->
    <div class="section">
      <h2 class="section-title">数据统计</h2>
      <div class="stats-grid" v-loading="loading">
        <div class="stat-item">
          <span class="stat-value">{{ systemStats.totalUsers }}</span>
          <span class="stat-label">注册用户</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ systemStats.totalAdmins }}</span>
          <span class="stat-label">管理员</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ systemStats.totalCheckins }}</span>
          <span class="stat-label">累计打卡</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ visitStats.total_visits }}</span>
          <span class="stat-label">总访问量</span>
        </div>
        <div class="stat-item">
          <span class="stat-value">{{ visitStats.today_visits }}</span>
          <span class="stat-label">今日访问</span>
        </div>
      </div>
    </div>

    <!-- 称号配置 -->
    <div class="section">
      <h2 class="section-title">
        称号配置
        <el-tooltip content="称号根据用户累计打卡次数自动计算，管理员可在用户管理中手动设置" placement="top">
          <el-icon class="info-icon"><InfoFilled /></el-icon>
        </el-tooltip>
      </h2>
      <div class="titles-card">
        <div
          v-for="title in titleConfig"
          :key="title.name"
          class="title-row"
        >
          <div class="title-info">
            <span class="title-badge" :style="{ color: title.color }">
              {{ title.icon }} {{ title.name }}
            </span>
            <span class="title-desc">{{ title.description }}</span>
          </div>
          <div class="title-requirement">
            <span class="requirement-value">{{ title.min }}+</span>
            <span class="requirement-label">次打卡</span>
          </div>
        </div>
      </div>
    </div>

    <!-- 关于 -->
    <div class="section">
      <h2 class="section-title">关于</h2>
      <div class="about-card">
        <p class="about-text">
          TidalCore 是一个开源的盆底肌训练平台，帮助用户进行科学的盆底肌锻炼。
        </p>
        <div class="about-links">
          <a href="https://github.com" target="_blank" class="about-link">
            GitHub 仓库
          </a>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.settings-page {
  max-width: 800px;
  margin: 0 auto;
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
  display: flex;
  align-items: center;
  gap: 8px;
}

.info-icon {
  color: rgba(255, 255, 255, 0.4);
  cursor: help;
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

/* ===== 统计网格 ===== */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 16px;
}

.stat-item {
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: 12px;
  padding: 20px;
  text-align: center;
}

.stat-item .stat-value {
  display: block;
  font-size: 28px;
  font-weight: 700;
  color: rgb(var(--ocean-surface));
  margin-bottom: 4px;
}

.stat-item .stat-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

/* ===== 称号卡片 ===== */
.titles-card {
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: 16px;
  padding: 8px 0;
}

.title-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 14px 20px;
  border-bottom: 1px solid rgba(56, 189, 248, 0.08);
  transition: background 0.2s;
}

.title-row:last-child {
  border-bottom: none;
}

.title-row:hover {
  background: rgba(56, 189, 248, 0.05);
}

.title-info {
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.title-badge {
  font-size: 15px;
  font-weight: 600;
}

.title-desc {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.title-requirement {
  text-align: right;
}

.requirement-value {
  display: block;
  font-size: 16px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
}

.requirement-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

/* ===== 关于卡片 ===== */
.about-card {
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: 16px;
  padding: 20px;
}

.about-text {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.6;
  margin: 0 0 16px;
}

.about-links {
  display: flex;
  gap: 16px;
}

.about-link {
  font-size: 14px;
  color: rgb(var(--ocean-surface));
  text-decoration: none;
  transition: opacity 0.2s;
}

.about-link:hover {
  opacity: 0.8;
}

/* ===== 响应式 ===== */
@media (max-width: 768px) {
  .stats-grid {
    grid-template-columns: repeat(2, 1fr);
  }

  .stat-item {
    padding: 16px;
  }

  .stat-item .stat-value {
    font-size: 24px;
  }

  .title-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }

  .title-requirement {
    text-align: left;
  }
}
</style>
