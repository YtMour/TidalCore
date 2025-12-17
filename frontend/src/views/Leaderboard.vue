<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import LeaderboardTable from '@/components/LeaderboardTable.vue'
import { getLeaderboard, type LeaderboardUser } from '@/api/checkin'
import { Trophy, Timer, TrendCharts, User, Medal, StarFilled, Pointer } from '@element-plus/icons-vue'

const users = ref<LeaderboardUser[]>([])
const loading = ref(true)
const mounted = ref(false)

onMounted(async () => {
  setTimeout(() => {
    mounted.value = true
  }, 100)

  try {
    users.value = await getLeaderboard(50)
  } catch {
    // 静默处理
  } finally {
    loading.value = false
  }
})

// Calculate stats
const totalUsers = computed(() => users.value.length)
const topStreak = computed(() => users.value[0]?.streak || 0)
const top3Users = computed(() => users.value.slice(0, 3))
</script>

<template>
  <MainLayout>
    <div class="leaderboard-page">
      <!-- Header -->
      <section class="header-section" :class="{ mounted }">
        <!-- Decorative background -->
        <div class="header-bg">
          <div class="header-glow"></div>
        </div>

        <div class="trophy-icon">
          <el-icon :size="48" color="#fff"><Trophy /></el-icon>
        </div>

        <el-tag type="warning" effect="plain" round class="header-badge">
          <el-icon><Medal /></el-icon>
          荣耀榜单
        </el-tag>

        <h1 class="header-title">
          <span class="gradient-text-gold">毅力排行榜</span>
        </h1>
        <p class="header-subtitle">
          坚持打卡，见证你的毅力。与全站用户一起，互相激励，共同进步。
        </p>
      </section>

      <!-- Top 3 Podium (when data loaded) -->
      <section
        v-if="!loading && top3Users.length >= 3"
        class="podium-section"
        :class="{ mounted }"
      >
        <div class="podium-wrapper">
          <!-- 2nd Place -->
          <div class="podium-item second">
            <div class="podium-avatar silver">
              {{ top3Users[1]?.username?.[0]?.toUpperCase() || '?' }}
            </div>
            <el-icon :size="24" color="#94a3b8"><Medal /></el-icon>
            <div class="podium-name">{{ top3Users[1]?.username }}</div>
            <div class="podium-streak">{{ top3Users[1]?.streak }} 天</div>
            <div class="podium-bar silver"></div>
          </div>

          <!-- 1st Place -->
          <div class="podium-item first">
            <div class="crown-icon">
              <el-icon :size="32" color="#fbbf24"><Medal /></el-icon>
            </div>
            <div class="podium-avatar gold">
              {{ top3Users[0]?.username?.[0]?.toUpperCase() || '?' }}
            </div>
            <el-icon :size="24" color="#fbbf24"><Trophy /></el-icon>
            <div class="podium-name">{{ top3Users[0]?.username }}</div>
            <div class="podium-streak gold">{{ top3Users[0]?.streak }} 天</div>
            <div class="podium-bar gold"></div>
          </div>

          <!-- 3rd Place -->
          <div class="podium-item third">
            <div class="podium-avatar bronze">
              {{ top3Users[2]?.username?.[0]?.toUpperCase() || '?' }}
            </div>
            <el-icon :size="24" color="#d97706"><StarFilled /></el-icon>
            <div class="podium-name">{{ top3Users[2]?.username }}</div>
            <div class="podium-streak">{{ top3Users[2]?.streak }} 天</div>
            <div class="podium-bar bronze"></div>
          </div>
        </div>
      </section>

      <!-- Stats Cards -->
      <section class="stats-section" :class="{ mounted }">
        <el-row :gutter="16">
          <el-col :xs="12" :md="6">
            <div class="stat-card stat-purple">
              <div class="stat-icon">
                <el-icon :size="20" color="#a78bfa"><User /></el-icon>
              </div>
              <div class="stat-value">{{ totalUsers }}</div>
              <div class="stat-label">参与用户</div>
            </div>
          </el-col>
          <el-col :xs="12" :md="6">
            <div class="stat-card stat-amber">
              <div class="stat-icon">
                <el-icon :size="20" color="#fbbf24"><Pointer /></el-icon>
              </div>
              <div class="stat-value">{{ topStreak }}</div>
              <div class="stat-label">最高连续</div>
            </div>
          </el-col>
          <el-col :xs="12" :md="6">
            <div class="stat-card stat-green">
              <div class="stat-icon">
                <el-icon :size="20" color="#34d399"><TrendCharts /></el-icon>
              </div>
              <div class="stat-value">TOP 50</div>
              <div class="stat-label">排行展示</div>
            </div>
          </el-col>
          <el-col :xs="12" :md="6">
            <div class="stat-card stat-pink">
              <div class="stat-icon">
                <el-icon :size="20" color="#f472b6"><StarFilled /></el-icon>
              </div>
              <div class="stat-value">实时</div>
              <div class="stat-label">数据更新</div>
            </div>
          </el-col>
        </el-row>
      </section>

      <!-- Leaderboard Table -->
      <section class="table-section" :class="{ mounted }">
        <LeaderboardTable :users="users" :loading="loading" />
      </section>

      <!-- CTA -->
      <section class="cta-section" :class="{ mounted }">
        <el-card class="cta-card" shadow="never">
          <!-- Decorative elements -->
          <div class="cta-decoration cta-decoration-1"></div>
          <div class="cta-decoration cta-decoration-2"></div>

          <div class="cta-content">
            <div class="cta-icon">
              <el-icon :size="28" color="#fff"><Trophy /></el-icon>
            </div>
            <h3 class="cta-title">想要冲击排行榜？</h3>
            <p class="cta-text">
              每天坚持训练，积累连续打卡天数，你也能登上榜单！
            </p>
            <RouterLink to="/train">
              <el-button type="primary" size="large" round class="cta-btn">
                <el-icon><Timer /></el-icon>
                开始训练，冲击排行
                <el-icon class="arrow-icon"><TrendCharts /></el-icon>
              </el-button>
            </RouterLink>
          </div>
        </el-card>
      </section>
    </div>
  </MainLayout>
</template>

<style scoped>
.leaderboard-page {
  display: flex;
  flex-direction: column;
  gap: 48px;
}

/* Header Section */
.header-section {
  text-align: center;
  position: relative;
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease;
}

.header-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.header-bg {
  position: absolute;
  inset: 0;
  z-index: -1;
  overflow: hidden;
}

.header-glow {
  position: absolute;
  top: 0;
  left: 50%;
  transform: translateX(-50%);
  width: 384px;
  height: 192px;
  background: linear-gradient(to bottom, rgba(245, 158, 11, 0.2), transparent);
  border-radius: 50%;
  filter: blur(60px);
}

.trophy-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 96px;
  height: 96px;
  border-radius: 24px;
  background: linear-gradient(135deg, #fbbf24, #f97316);
  margin-bottom: 24px;
  box-shadow: 0 20px 40px rgba(249, 115, 22, 0.3);
  animation: float 3s ease-in-out infinite;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-10px); }
}

.header-badge {
  margin-bottom: 16px;
  background: rgba(245, 158, 11, 0.1) !important;
  border-color: rgba(245, 158, 11, 0.2) !important;
  color: #fcd34d !important;
}

.header-badge .el-icon {
  margin-right: 6px;
}

.header-title {
  font-size: 48px;
  font-weight: 800;
  margin: 0 0 16px;
}

@media (max-width: 768px) {
  .header-title {
    font-size: 36px;
  }
}

.gradient-text-gold {
  background: linear-gradient(135deg, #fbbf24, #f97316);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-subtitle {
  font-size: 18px;
  color: rgba(255, 255, 255, 0.5);
  max-width: 512px;
  margin: 0 auto;
}

/* Podium Section */
.podium-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease 0.1s;
}

.podium-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.podium-wrapper {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  gap: 16px;
  max-width: 512px;
  margin: 0 auto;
}

@media (min-width: 768px) {
  .podium-wrapper {
    gap: 32px;
  }
}

.podium-item {
  flex: 1;
  text-align: center;
}

.podium-item.first {
  margin-top: -32px;
}

.crown-icon {
  margin-bottom: 8px;
  animation: float 2s ease-in-out infinite;
}

.podium-avatar {
  width: 64px;
  height: 64px;
  margin: 0 auto 12px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  color: #fff;
}

@media (min-width: 768px) {
  .podium-avatar {
    width: 80px;
    height: 80px;
    font-size: 32px;
  }
}

.podium-item.first .podium-avatar {
  width: 80px;
  height: 80px;
  font-size: 32px;
}

@media (min-width: 768px) {
  .podium-item.first .podium-avatar {
    width: 96px;
    height: 96px;
    font-size: 40px;
  }
}

.podium-avatar.gold {
  background: linear-gradient(135deg, #fbbf24, #f97316);
  box-shadow: 0 10px 30px rgba(249, 115, 22, 0.3);
}

.podium-avatar.silver {
  background: linear-gradient(135deg, #94a3b8, #64748b);
  box-shadow: 0 10px 20px rgba(100, 116, 139, 0.3);
}

.podium-avatar.bronze {
  background: linear-gradient(135deg, #d97706, #b45309);
  box-shadow: 0 10px 20px rgba(180, 83, 9, 0.3);
}

.podium-name {
  font-size: 14px;
  font-weight: 500;
  color: #fff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 8px;
  margin-top: 4px;
}

.podium-streak {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.podium-streak.gold {
  color: #fbbf24;
  font-weight: 600;
}

.podium-bar {
  margin-top: 12px;
  border-radius: 12px 12px 0 0;
}

.podium-bar.gold {
  height: 112px;
  background: linear-gradient(to top, rgba(245, 158, 11, 0.2), transparent);
}

@media (min-width: 768px) {
  .podium-bar.gold {
    height: 128px;
  }
}

.podium-bar.silver {
  height: 80px;
  background: linear-gradient(to top, rgba(148, 163, 184, 0.2), transparent);
}

@media (min-width: 768px) {
  .podium-bar.silver {
    height: 96px;
  }
}

.podium-bar.bronze {
  height: 64px;
  background: linear-gradient(to top, rgba(217, 119, 6, 0.2), transparent);
}

@media (min-width: 768px) {
  .podium-bar.bronze {
    height: 80px;
  }
}

/* Stats Section */
.stats-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease 0.2s;
}

.stats-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.stat-card {
  text-align: center;
  padding: 20px;
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
  width: 40px;
  height: 40px;
  margin: 0 auto 12px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s ease;
}

.stat-purple .stat-icon {
  background: linear-gradient(135deg, rgba(139, 92, 246, 0.2), rgba(168, 85, 247, 0.1));
}

.stat-amber .stat-icon {
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.2), rgba(249, 115, 22, 0.1));
}

.stat-green .stat-icon {
  background: linear-gradient(135deg, rgba(16, 185, 129, 0.2), rgba(52, 211, 153, 0.1));
}

.stat-pink .stat-icon {
  background: linear-gradient(135deg, rgba(236, 72, 153, 0.2), rgba(244, 114, 182, 0.1));
}

.stat-value {
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

/* Table Section */
.table-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease 0.3s;
}

.table-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

/* CTA Section */
.cta-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s ease 0.4s;
}

.cta-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.cta-card {
  max-width: 640px;
  margin: 0 auto;
  background: rgba(30, 30, 46, 0.8) !important;
  border: 1px solid rgba(255, 255, 255, 0.06) !important;
  border-radius: 16px !important;
  overflow: hidden;
  position: relative;
}

.cta-card :deep(.el-card__body) {
  padding: 40px;
}

.cta-decoration {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  pointer-events: none;
}

.cta-decoration-1 {
  top: -80px;
  left: -80px;
  width: 160px;
  height: 160px;
  background: linear-gradient(135deg, rgba(245, 158, 11, 0.15), transparent);
}

.cta-decoration-2 {
  bottom: -80px;
  right: -80px;
  width: 160px;
  height: 160px;
  background: linear-gradient(135deg, rgba(249, 115, 22, 0.15), transparent);
}

.cta-content {
  position: relative;
  text-align: center;
}

.cta-icon {
  width: 56px;
  height: 56px;
  margin: 0 auto 20px;
  border-radius: 16px;
  background: linear-gradient(135deg, #fbbf24, #f97316);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 10px 30px rgba(245, 158, 11, 0.3);
}

.cta-title {
  font-size: 20px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 12px;
}

.cta-text {
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 24px;
}

.cta-btn {
  background: linear-gradient(135deg, #8b5cf6, #ec4899) !important;
  border: none !important;
  padding: 12px 40px !important;
  font-size: 16px !important;
  box-shadow: 0 10px 30px rgba(139, 92, 246, 0.25);
}

.cta-btn .arrow-icon {
  margin-left: 8px;
  transition: transform 0.2s ease;
}

.cta-btn:hover .arrow-icon {
  transform: translateX(4px);
}
</style>
