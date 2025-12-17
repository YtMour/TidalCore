<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Heatmap from '@/components/Heatmap.vue'
import { useUserStore } from '@/store/user'
import { getHeatmap, getHistory, type CheckinRecord } from '@/api/checkin'
import { Timer, Calendar, Clock, CircleCheck, ArrowRight, Pointer, Trophy, Aim, StarFilled } from '@element-plus/icons-vue'

const userStore = useUserStore()

const heatmapData = ref<Record<string, number>>({})
const recentHistory = ref<CheckinRecord[]>([])
const loading = ref(true)
const mounted = ref(false)

// 计算用户等级
const userLevel = computed(() => {
  const total = userStore.user?.total_checkin || 0
  if (total >= 365) return { name: '传奇', color: '#fbbf24', bg: 'linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(249, 115, 22, 0.1))' }
  if (total >= 180) return { name: '大师', color: '#a78bfa', bg: 'linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(168, 85, 247, 0.1))' }
  if (total >= 90) return { name: '专家', color: '#22d3ee', bg: 'linear-gradient(135deg, rgba(6, 182, 212, 0.2), rgba(34, 211, 238, 0.1))' }
  if (total >= 30) return { name: '进阶', color: '#34d399', bg: 'linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(52, 211, 153, 0.1))' }
  if (total >= 7) return { name: '新手', color: '#f472b6', bg: 'linear-gradient(135deg, rgba(236, 72, 153, 0.2), rgba(244, 114, 182, 0.1))' }
  return { name: '初学者', color: 'rgba(255, 255, 255, 0.6)', bg: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.05))' }
})

onMounted(async () => {
  setTimeout(() => {
    mounted.value = true
  }, 100)

  await userStore.fetchProfile()
  await loadData()
})

async function loadData() {
  loading.value = true
  try {
    const [checkins, history] = await Promise.all([
      getHeatmap(365),
      getHistory(10)
    ])

    const heatmap: Record<string, number> = {}
    checkins.forEach(c => {
      const datePart = c.checked_at.split('T')[0]
      if (datePart) {
        heatmap[datePart] = (heatmap[datePart] || 0) + 1
      }
    })
    heatmapData.value = heatmap
    recentHistory.value = history
  } catch {
    // 静默处理
  } finally {
    loading.value = false
  }
}

function formatDate(dateStr: string): string {
  const date = new Date(dateStr)
  return date.toLocaleDateString('zh-CN', {
    month: 'short',
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

function formatDuration(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return mins > 0 ? `${mins}分${secs}秒` : `${secs}秒`
}

function getRelativeTime(dateStr: string): string {
  const date = new Date(dateStr)
  const now = new Date()
  const diff = now.getTime() - date.getTime()
  const days = Math.floor(diff / (1000 * 60 * 60 * 24))

  if (days === 0) return '今天'
  if (days === 1) return '昨天'
  if (days < 7) return `${days}天前`
  return formatDate(dateStr)
}
</script>

<template>
  <MainLayout>
    <div class="dashboard-page">
      <!-- Profile Header -->
      <section class="profile-section" :class="{ mounted }">
        <el-card class="profile-card" shadow="never">
          <!-- Decorative background -->
          <div class="profile-decoration profile-decoration-1"></div>
          <div class="profile-decoration profile-decoration-2"></div>

          <div class="profile-content">
            <!-- Avatar with level badge -->
            <div class="avatar-wrapper">
              <div class="avatar">
                {{ userStore.user?.username?.[0]?.toUpperCase() || '?' }}
              </div>
              <!-- Level badge -->
              <div
                class="level-badge"
                :style="{ background: userLevel.bg, color: userLevel.color }"
              >
                {{ userLevel.name }}
              </div>
            </div>

            <!-- User Info -->
            <div class="user-info">
              <div class="user-name">
                <h1>{{ userStore.user?.username }}</h1>
                <el-icon :size="20" color="#fbbf24"><StarFilled /></el-icon>
              </div>
              <p class="user-motto">坚持就是胜利，每天进步一点点</p>

              <!-- Quick stats inline -->
              <div class="inline-stats">
                <span class="inline-stat orange">
                  <el-icon><Pointer /></el-icon>
                  {{ userStore.user?.streak || 0 }} 天连续
                </span>
                <span class="inline-stat purple">
                  <el-icon><Trophy /></el-icon>
                  {{ userStore.user?.max_streak || 0 }} 最高
                </span>
              </div>

              <RouterLink to="/train">
                <el-button type="primary" size="large" round class="train-btn">
                  <el-icon><Timer /></el-icon>
                  开始今日训练
                  <el-icon class="arrow-icon"><ArrowRight /></el-icon>
                </el-button>
              </RouterLink>
            </div>
          </div>
        </el-card>
      </section>

      <!-- Stats Cards -->
      <section class="stats-section" :class="{ mounted }">
        <el-row :gutter="16">
          <el-col :xs="24" :sm="8">
            <div class="stat-card stat-orange">
              <div class="stat-icon">
                <el-icon :size="28" color="#fb923c"><Pointer /></el-icon>
              </div>
              <div class="stat-value">{{ userStore.user?.streak || 0 }}</div>
              <div class="stat-label">连续打卡</div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="8">
            <div class="stat-card stat-purple">
              <div class="stat-icon">
                <el-icon :size="28" color="#a78bfa"><Trophy /></el-icon>
              </div>
              <div class="stat-value">{{ userStore.user?.max_streak || 0 }}</div>
              <div class="stat-label">最高记录</div>
            </div>
          </el-col>
          <el-col :xs="24" :sm="8">
            <div class="stat-card stat-green">
              <div class="stat-icon">
                <el-icon :size="28" color="#34d399"><Aim /></el-icon>
              </div>
              <div class="stat-value">{{ userStore.user?.total_checkin || 0 }}</div>
              <div class="stat-label">累计打卡</div>
            </div>
          </el-col>
        </el-row>
      </section>

      <!-- Heatmap Section -->
      <section class="heatmap-section" :class="{ mounted }">
        <el-card class="section-card" shadow="never">
          <div class="section-header">
            <el-icon :size="20" color="#a78bfa"><Calendar /></el-icon>
            <h2>我的打卡热力图</h2>
          </div>

          <div v-if="loading" class="loading-content">
            <el-skeleton :rows="3" animated />
          </div>
          <Heatmap v-else :data="heatmapData" :days="365" />
        </el-card>
      </section>

      <!-- Recent History -->
      <section class="history-section" :class="{ mounted }">
        <el-card class="section-card" shadow="never">
          <div class="section-header">
            <el-icon :size="20" color="#a78bfa"><Clock /></el-icon>
            <h2>最近打卡记录</h2>
          </div>

          <div v-if="loading" class="loading-content">
            <div v-for="i in 5" :key="i" class="history-skeleton">
              <el-skeleton-item variant="rect" style="width: 48px; height: 48px; border-radius: 12px" />
              <div class="skeleton-text">
                <el-skeleton-item variant="text" style="width: 60%" />
                <el-skeleton-item variant="text" style="width: 40%; height: 12px" />
              </div>
            </div>
          </div>

          <div v-else-if="recentHistory.length === 0" class="empty-state">
            <div class="empty-icon">
              <el-icon :size="32" color="rgba(255, 255, 255, 0.3)"><Timer /></el-icon>
            </div>
            <p class="empty-title">暂无打卡记录</p>
            <p class="empty-desc">完成训练后，记录将显示在这里</p>
          </div>

          <div v-else class="history-list">
            <TransitionGroup name="list">
              <div
                v-for="(record, index) in recentHistory"
                :key="record.id"
                class="history-item"
                :style="{ transitionDelay: `${index * 50}ms` }"
              >
                <div class="history-icon">
                  <el-icon :size="24" color="#34d399"><CircleCheck /></el-icon>
                </div>
                <div class="history-info">
                  <div class="history-time">{{ getRelativeTime(record.checked_at) }}</div>
                  <div class="history-detail">
                    {{ record.cycles }} 个循环 · {{ formatDuration(record.duration) }}
                  </div>
                </div>
                <div class="history-date">
                  {{ formatDate(record.checked_at) }}
                </div>
              </div>
            </TransitionGroup>
          </div>
        </el-card>
      </section>
    </div>
  </MainLayout>
</template>

<style scoped>
.dashboard-page {
  display: flex;
  flex-direction: column;
  gap: 32px;
}

/* Profile Section */
.profile-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease;
}

.profile-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.profile-card {
  background: rgba(30, 30, 46, 0.8) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 16px !important;
  overflow: hidden;
  position: relative;
}

.profile-card :deep(.el-card__body) {
  padding: 32px;
}

.profile-decoration {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  pointer-events: none;
}

.profile-decoration-1 {
  top: 0;
  right: 0;
  width: 256px;
  height: 256px;
  background: linear-gradient(225deg, rgba(139, 92, 246, 0.1), transparent);
}

.profile-decoration-2 {
  bottom: 0;
  left: 0;
  width: 192px;
  height: 192px;
  background: linear-gradient(45deg, rgba(236, 72, 153, 0.1), transparent);
}

.profile-content {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 24px;
}

@media (min-width: 640px) {
  .profile-content {
    flex-direction: row;
    align-items: flex-start;
  }
}

.avatar-wrapper {
  position: relative;
  flex-shrink: 0;
}

.avatar {
  width: 112px;
  height: 112px;
  border-radius: 24px;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 48px;
  font-weight: 700;
  color: #fff;
  box-shadow: 0 20px 40px rgba(139, 92, 246, 0.3);
}

.level-badge {
  position: absolute;
  bottom: -8px;
  right: -8px;
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 12px;
  font-weight: 600;
  border: 1px solid rgba(255, 255, 255, 0.1);
}

.user-info {
  flex: 1;
  text-align: center;
}

@media (min-width: 640px) {
  .user-info {
    text-align: left;
  }
}

.user-name {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  margin-bottom: 4px;
}

@media (min-width: 640px) {
  .user-name {
    justify-content: flex-start;
  }
}

.user-name h1 {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin: 0;
}

.user-motto {
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 8px;
}

.inline-stats {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 16px;
  font-size: 14px;
}

@media (min-width: 640px) {
  .inline-stats {
    justify-content: flex-start;
  }
}

.inline-stat {
  display: flex;
  align-items: center;
  gap: 4px;
}

.inline-stat.orange {
  color: #fb923c;
}

.inline-stat.purple {
  color: #a78bfa;
}

.train-btn {
  background: linear-gradient(135deg, #8b5cf6, #ec4899) !important;
  border: none !important;
  padding: 12px 28px !important;
  box-shadow: 0 10px 30px rgba(139, 92, 246, 0.25);
}

.train-btn .arrow-icon {
  margin-left: 8px;
  transition: transform 0.2s ease;
}

.train-btn:hover .arrow-icon {
  transform: translateX(4px);
}

/* Stats Section */
.stats-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease 0.1s;
}

.stats-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.stat-card {
  text-align: center;
  padding: 24px;
  border-radius: 16px;
  background: rgba(30, 30, 46, 0.6);
  border: 1px solid rgba(255, 255, 255, 0.06);
  margin-bottom: 16px;
  transition: all 0.3s ease;
}

.stat-card:hover {
  border-color: rgba(255, 255, 255, 0.1);
}

.stat-card:hover .stat-icon {
  transform: scale(1.1);
}

.stat-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 16px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s ease;
}

.stat-orange .stat-icon {
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.2), rgba(251, 146, 60, 0.1));
}

.stat-purple .stat-icon {
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(168, 85, 247, 0.1));
}

.stat-green .stat-icon {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(52, 211, 153, 0.1));
}

.stat-value {
  font-size: 32px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
}

/* Heatmap Section */
.heatmap-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.5s ease 0.2s;
}

.heatmap-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.section-card {
  background: rgba(30, 30, 46, 0.8) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 16px !important;
}

.section-card :deep(.el-card__body) {
  padding: 24px;
}

.section-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 24px;
}

.section-header h2 {
  font-size: 18px;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.loading-content {
  padding: 32px 0;
}

/* History Section */
.history-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.5s ease 0.3s;
}

.history-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.history-skeleton {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.05);
  margin-bottom: 12px;
}

.skeleton-text {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.empty-state {
  text-align: center;
  padding: 48px 0;
}

.empty-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 16px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.05);
  display: flex;
  align-items: center;
  justify-content: center;
}

.empty-title {
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 8px;
}

.empty-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.3);
  margin: 0;
}

.history-list {
  display: flex;
  flex-direction: column;
  gap: 12px;
}

.history-item {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 16px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.05);
  transition: all 0.3s ease;
}

.history-item:hover {
  background: rgba(255, 255, 255, 0.08);
}

.history-icon {
  width: 48px;
  height: 48px;
  border-radius: 12px;
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(20, 184, 166, 0.2));
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}

.history-info {
  flex: 1;
  min-width: 0;
}

.history-time {
  font-weight: 500;
  color: #fff;
}

.history-detail {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.5);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.history-date {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.3);
  flex-shrink: 0;
}

/* List transition */
.list-enter-active,
.list-leave-active {
  transition: all 0.3s ease;
}

.list-enter-from,
.list-leave-to {
  opacity: 0;
  transform: translateX(-20px);
}

</style>
