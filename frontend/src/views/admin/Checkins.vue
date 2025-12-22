<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { getGlobalHeatmap, getVisitHeatmap, getVisitStats } from '@/api/checkin'
import { getUsers, type UsersResponse } from '@/api/admin'
import type { UserInfo } from '@/api/auth'
import { TrendCharts, Timer, Pointer, View } from '@element-plus/icons-vue'
import Heatmap from '@/components/Heatmap.vue'

const loading = ref(true)
const checkinHeatmapData = ref<Record<string, number>>({})
const visitHeatmapData = ref<Record<string, number>>({})
const heatmapType = ref<'checkin' | 'visit'>('checkin')
const users = ref<UserInfo[]>([])

// 当前显示的热力图数据
const currentHeatmapData = computed(() => {
  return heatmapType.value === 'checkin' ? checkinHeatmapData.value : visitHeatmapData.value
})

// 统计数据
const stats = ref({
  totalCheckins: 0,
  todayCheckins: 0,
  totalVisits: 0,
  todayVisits: 0
})

onMounted(async () => {
  await Promise.all([
    loadCheckinHeatmap(),
    loadVisitHeatmap(),
    loadVisitStats(),
    loadUsers()
  ])
  calculateStats()
  loading.value = false
})

async function loadCheckinHeatmap() {
  try {
    checkinHeatmapData.value = await getGlobalHeatmap(365)
  } catch (error) {
    console.error('加载打卡热力图数据失败', error)
  }
}

async function loadVisitHeatmap() {
  try {
    visitHeatmapData.value = await getVisitHeatmap(365)
  } catch (error) {
    console.error('加载访问热力图数据失败', error)
  }
}

async function loadVisitStats() {
  try {
    const visitStats = await getVisitStats()
    stats.value.totalVisits = visitStats.total_visits || 0
    stats.value.todayVisits = visitStats.today_visits || 0
  } catch (error) {
    console.error('加载访问统计失败', error)
  }
}

async function loadUsers() {
  try {
    const res: UsersResponse = await getUsers(1, 1000)
    users.value = res.users
  } catch (error) {
    console.error('加载用户数据失败', error)
  }
}

function calculateStats() {
  // 计算总打卡次数
  let total = 0
  Object.values(checkinHeatmapData.value).forEach(count => {
    total += count
  })
  stats.value.totalCheckins = total

  // 计算今日打卡
  const today = new Date()
  const todayKey = `${today.getFullYear()}-${String(today.getMonth() + 1).padStart(2, '0')}-${String(today.getDate()).padStart(2, '0')}`
  stats.value.todayCheckins = checkinHeatmapData.value[todayKey] || 0
}

// 热力图统计数据
const heatmapStats = computed(() => {
  const data = currentHeatmapData.value || {}
  const total = Object.values(data).reduce((sum, count) => sum + count, 0)
  const activeDays = Object.values(data).filter(count => count > 0).length
  return { total, activeDays }
})

// 用户打卡排行
const topUsers = computed(() => {
  return [...users.value]
    .sort((a, b) => (b.total_checkin || 0) - (a.total_checkin || 0))
    .slice(0, 10)
})

// 称号配置
const titleConfig = [
  { min: 1000, name: '海神降临', color: '#f472b6', icon: '🔱' },
  { min: 730, name: '深渊霸主', color: '#a78bfa', icon: '🦑' },
  { min: 365, name: '深海传奇', color: '#fbbf24', icon: '🌊' },
  { min: 180, name: '海洋大师', color: '#38bdf8', icon: '🐋' },
  { min: 90, name: '浪潮专家', color: '#22d3ee', icon: '🐬' },
  { min: 30, name: '潮汐进阶', color: '#34d399', icon: '🐠' },
  { min: 7, name: '入海新手', color: '#0ea5e9', icon: '🐟' },
  { min: 0, name: '初探海域', color: 'rgba(255, 255, 255, 0.6)', icon: '🐚' }
]

function getUserTitle(user: UserInfo) {
  if (user.title) {
    const found = titleConfig.find(t => t.name === user.title)
    if (found) return found
  }
  for (const title of titleConfig) {
    if ((user.total_checkin || 0) >= title.min) {
      return title
    }
  }
  return titleConfig[titleConfig.length - 1]!
}
</script>

<template>
  <div class="checkins-page">
    <!-- 统计卡片 -->
    <div class="stats-grid">
      <div class="stat-card" v-loading="loading">
        <div class="stat-icon ocean">
          <el-icon :size="24"><Pointer /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.totalCheckins }}</span>
          <span class="stat-label">累计打卡次数</span>
        </div>
      </div>

      <div class="stat-card" v-loading="loading">
        <div class="stat-icon green">
          <el-icon :size="24"><TrendCharts /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.todayCheckins }}</span>
          <span class="stat-label">今日打卡</span>
        </div>
      </div>

      <div class="stat-card" v-loading="loading">
        <div class="stat-icon purple">
          <el-icon :size="24"><View /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.totalVisits }}</span>
          <span class="stat-label">累计访问</span>
        </div>
      </div>

      <div class="stat-card" v-loading="loading">
        <div class="stat-icon amber">
          <el-icon :size="24"><Timer /></el-icon>
        </div>
        <div class="stat-info">
          <span class="stat-value">{{ stats.todayVisits }}</span>
          <span class="stat-label">今日访问</span>
        </div>
      </div>
    </div>

    <!-- 热力图区域 -->
    <div class="section">
      <div class="heatmap-card" v-loading="loading">
        <!-- 热力图头部 -->
        <div class="heatmap-header">
          <div class="heatmap-title-area">
            <h2 class="heatmap-title">
              {{ heatmapType === 'checkin' ? '全站打卡热力图' : '站点访问热力图' }}
            </h2>
            <p class="heatmap-desc">
              {{ heatmapType === 'checkin' ? '社区成员的坚持轨迹' : '每日访客的足迹记录' }}
            </p>
          </div>

          <!-- 热力图切换按钮 -->
          <div class="heatmap-toggle">
            <button
              class="toggle-btn"
              :class="{ active: heatmapType === 'checkin' }"
              @click="heatmapType = 'checkin'"
            >
              <el-icon :size="14"><Pointer /></el-icon>
              打卡
            </button>
            <button
              class="toggle-btn"
              :class="{ active: heatmapType === 'visit' }"
              @click="heatmapType = 'visit'"
            >
              <el-icon :size="14"><View /></el-icon>
              访问
            </button>
          </div>
        </div>

        <!-- 热力图统计 -->
        <div class="heatmap-stats">
          <div class="hm-stat">
            <span class="hm-value">{{ heatmapStats.total }}</span>
            <span class="hm-label">{{ heatmapType === 'checkin' ? '次打卡' : '次访问' }}</span>
          </div>
          <div class="hm-divider"></div>
          <div class="hm-stat">
            <span class="hm-value">{{ heatmapStats.activeDays }}</span>
            <span class="hm-label">天活跃</span>
          </div>
          <div class="live-badge">
            <span class="live-dot"></span>
            实时
          </div>
        </div>

        <!-- 热力图组件 -->
        <div class="heatmap-content">
          <Heatmap :data="currentHeatmapData" :days="365" :label="heatmapType === 'checkin' ? '打卡' : '访问'" />
        </div>
      </div>
    </div>

    <!-- 打卡排行 -->
    <div class="section">
      <h2 class="section-title">打卡排行榜 TOP 10</h2>
      <div class="ranking-card" v-loading="loading">
        <div
          v-for="(user, index) in topUsers"
          :key="user.id"
          class="ranking-item"
        >
          <span class="ranking-number" :class="{ top3: index < 3 }">
            {{ index + 1 }}
          </span>
          <div class="ranking-user">
            <div class="ranking-name-row">
              <span class="ranking-name">{{ user.display_name || user.username }}</span>
              <span class="ranking-title" :style="{ color: getUserTitle(user).color }">
                {{ getUserTitle(user).icon }} {{ getUserTitle(user).name }}
              </span>
            </div>
            <span class="ranking-username">@{{ user.username }}</span>
          </div>
          <div class="ranking-stats">
            <span class="ranking-checkins">{{ user.total_checkin || 0 }} 次</span>
            <span class="ranking-streak">连续 {{ user.streak || 0 }} 天</span>
          </div>
        </div>

        <div v-if="topUsers.length === 0" class="empty-state">
          暂无打卡数据
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.checkins-page {
  max-width: 1200px;
  margin: 0 auto;
}

/* ===== 统计卡片 ===== */
.stats-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  margin-bottom: 24px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: var(--glass-bg);
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-xl);
  transition: all 0.3s var(--ease-smooth);
}

.stat-card:hover {
  transform: translateY(-2px);
  border-color: rgba(56, 189, 248, 0.2);
  box-shadow: 0 8px 32px rgba(56, 189, 248, 0.1);
}

.stat-icon {
  width: 52px;
  height: 52px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: var(--radius-lg);
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

.stat-icon.purple {
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(139, 92, 246, 0.1));
  color: rgb(var(--pearl-purple));
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
  font-size: 26px;
  font-weight: 700;
  color: #fff;
}

.stat-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.6);
}

/* ===== Section ===== */
.section {
  margin-bottom: 24px;
}

.section-title {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 16px;
}

/* ===== 热力图卡片 ===== */
.heatmap-card {
  background: rgba(8, 15, 30, 0.4);
  border: 1px solid rgba(56, 189, 248, 0.08);
  border-radius: var(--radius-2xl);
  padding: 24px;
}

.heatmap-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 20px;
  flex-wrap: wrap;
  gap: 16px;
}

.heatmap-title-area {
  flex: 1;
}

.heatmap-title {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0 0 4px;
}

.heatmap-desc {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  margin: 0;
}

/* 热力图切换按钮 */
.heatmap-toggle {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 4px;
  background: rgba(56, 189, 248, 0.05);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-lg);
}

.heatmap-toggle .toggle-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 16px;
  border: none;
  border-radius: var(--radius-md);
  background: transparent;
  color: rgba(255, 255, 255, 0.6);
  font-size: 13px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.heatmap-toggle .toggle-btn:hover {
  color: rgba(255, 255, 255, 0.8);
  background: rgba(56, 189, 248, 0.08);
}

.heatmap-toggle .toggle-btn.active {
  background: rgba(56, 189, 248, 0.15);
  color: rgb(var(--ocean-surface));
}

/* 热力图统计 */
.heatmap-stats {
  display: flex;
  align-items: center;
  gap: 20px;
  margin-bottom: 20px;
  padding-bottom: 16px;
  border-bottom: 1px solid rgba(56, 189, 248, 0.08);
}

.hm-stat {
  display: flex;
  align-items: baseline;
  gap: 6px;
}

.hm-value {
  font-size: 24px;
  font-weight: 700;
  color: rgb(var(--ocean-surface));
}

.hm-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.hm-divider {
  width: 1px;
  height: 24px;
  background: rgba(56, 189, 248, 0.15);
}

.live-badge {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 4px 10px;
  background: rgba(52, 211, 153, 0.1);
  border: 1px solid rgba(52, 211, 153, 0.2);
  border-radius: var(--radius-full);
  font-size: 12px;
  color: rgb(var(--seaweed-green));
  margin-left: auto;
}

.live-dot {
  width: 6px;
  height: 6px;
  background: rgb(var(--seaweed-green));
  border-radius: 50%;
  animation: pulse-live 2s ease-in-out infinite;
}

@keyframes pulse-live {
  0%, 100% { opacity: 1; transform: scale(1); }
  50% { opacity: 0.5; transform: scale(0.8); }
}

.heatmap-content {
  /* 热力图组件容器 */
}

/* ===== 排行榜 ===== */
.ranking-card {
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  border-radius: var(--radius-xl);
  padding: 8px 0;
}

.ranking-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 14px 20px;
  border-bottom: 1px solid rgba(56, 189, 248, 0.08);
  transition: background 0.2s;
}

.ranking-item:last-child {
  border-bottom: none;
}

.ranking-item:hover {
  background: rgba(56, 189, 248, 0.05);
}

.ranking-number {
  width: 32px;
  height: 32px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
  background: rgba(56, 189, 248, 0.1);
  border-radius: 8px;
  flex-shrink: 0;
}

.ranking-number.top3 {
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.3), rgba(251, 191, 36, 0.1));
  color: rgb(var(--sunset-amber));
}

.ranking-user {
  flex: 1;
  min-width: 0;
}

.ranking-name-row {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
}

.ranking-name {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
}

.ranking-title {
  font-size: 12px;
  font-weight: 500;
}

.ranking-username {
  display: block;
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 2px;
}

.ranking-stats {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  gap: 2px;
}

.ranking-checkins {
  font-size: 14px;
  font-weight: 600;
  color: rgb(var(--ocean-surface));
}

.ranking-streak {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.empty-state {
  padding: 40px;
  text-align: center;
  color: rgba(255, 255, 255, 0.5);
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
    width: 44px;
    height: 44px;
  }

  .stat-value {
    font-size: 22px;
  }

  .heatmap-card {
    padding: 16px;
  }

  .heatmap-header {
    flex-direction: column;
    align-items: flex-start;
  }

  .heatmap-toggle {
    width: 100%;
    justify-content: center;
  }

  .heatmap-stats {
    flex-wrap: wrap;
    gap: 12px;
  }

  .live-badge {
    margin-left: 0;
  }

  .ranking-item {
    padding: 12px 16px;
  }

  .ranking-name-row {
    flex-direction: column;
    align-items: flex-start;
    gap: 4px;
  }
}
</style>
