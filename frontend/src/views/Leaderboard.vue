<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import LeaderboardTable from '@/components/LeaderboardTable.vue'
import { getLeaderboard, type LeaderboardUser } from '@/api/checkin'
import { Trophy, Timer, TrendCharts, User, Medal, Pointer } from '@element-plus/icons-vue'

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
      <!-- Header - 深海排行 -->
      <section class="header-section" :class="{ mounted }">
        <!-- Decorative background - 海洋主题 -->
        <div class="header-bg">
          <div class="header-glow"></div>
          <div class="header-ripple header-ripple-1"></div>
          <div class="header-ripple header-ripple-2"></div>
        </div>

        <div class="trophy-icon">
          <svg viewBox="0 0 48 48" fill="none" width="48" height="48" class="trophy-svg">
            <circle cx="24" cy="24" r="20" fill="rgba(255,255,255,0.1)"/>
            <path d="M16 24C16 24 20 18 24 24C28 30 32 18 36 24" stroke="currentColor" stroke-width="3" stroke-linecap="round"/>
            <circle cx="24" cy="24" r="6" fill="currentColor" opacity="0.3"/>
          </svg>
        </div>

        <div class="header-eyebrow">
          <span class="eyebrow-line"></span>
          <span class="eyebrow-text">深海榜单</span>
          <span class="eyebrow-line"></span>
        </div>

        <h1 class="header-title">
          <span class="gradient-text-ocean">深海排行榜</span>
        </h1>
        <p class="header-subtitle">
          与全站用户一起潜入深海，互相激励，共同见证潮汐的力量。
        </p>
      </section>

      <!-- Top 3 Podium (when data loaded) - 海洋主题 -->
      <section
        v-if="!loading && top3Users.length >= 3"
        class="podium-section"
        :class="{ mounted }"
      >
        <div class="podium-wrapper">
          <!-- 2nd Place - 银色浪花 -->
          <div class="podium-item second">
            <div class="podium-ring"></div>
            <div class="podium-avatar silver">
              {{ top3Users[1]?.username?.[0]?.toUpperCase() || '?' }}
            </div>
            <div class="podium-medal silver">
              <el-icon :size="18"><Medal /></el-icon>
            </div>
            <div class="podium-name">{{ top3Users[1]?.username }}</div>
            <div class="podium-streak">{{ top3Users[1]?.streak }} 天</div>
            <div class="podium-bar silver"></div>
          </div>

          <!-- 1st Place - 金色深海王者 -->
          <div class="podium-item first">
            <div class="crown-wrapper">
              <div class="crown-glow"></div>
              <el-icon :size="28" color="#fbbf24"><Trophy /></el-icon>
            </div>
            <div class="podium-ring gold"></div>
            <div class="podium-avatar gold">
              {{ top3Users[0]?.username?.[0]?.toUpperCase() || '?' }}
            </div>
            <div class="podium-medal gold">
              <el-icon :size="18"><Trophy /></el-icon>
            </div>
            <div class="podium-name gold">{{ top3Users[0]?.username }}</div>
            <div class="podium-streak gold">{{ top3Users[0]?.streak }} 天</div>
            <div class="podium-bar gold"></div>
          </div>

          <!-- 3rd Place - 铜色珊瑚 -->
          <div class="podium-item third">
            <div class="podium-ring"></div>
            <div class="podium-avatar bronze">
              {{ top3Users[2]?.username?.[0]?.toUpperCase() || '?' }}
            </div>
            <div class="podium-medal bronze">
              <el-icon :size="18"><Medal /></el-icon>
            </div>
            <div class="podium-name">{{ top3Users[2]?.username }}</div>
            <div class="podium-streak">{{ top3Users[2]?.streak }} 天</div>
            <div class="podium-bar bronze"></div>
          </div>
        </div>
      </section>

      <!-- Stats Cards - 海洋主题 -->
      <section class="stats-section" :class="{ mounted }">
        <el-row :gutter="16">
          <el-col :xs="12" :md="6">
            <div class="stat-card stat-ocean">
              <div class="stat-glow"></div>
              <div class="stat-icon">
                <el-icon :size="20" color="#38bdf8"><User /></el-icon>
              </div>
              <div class="stat-value">{{ totalUsers }}</div>
              <div class="stat-label">入海用户</div>
            </div>
          </el-col>
          <el-col :xs="12" :md="6">
            <div class="stat-card stat-amber">
              <div class="stat-glow"></div>
              <div class="stat-icon">
                <el-icon :size="20" color="#fbbf24"><Pointer /></el-icon>
              </div>
              <div class="stat-value">{{ topStreak }}</div>
              <div class="stat-label">最深潜水</div>
            </div>
          </el-col>
          <el-col :xs="12" :md="6">
            <div class="stat-card stat-aqua">
              <div class="stat-glow"></div>
              <div class="stat-icon">
                <el-icon :size="20" color="#22d3ee"><TrendCharts /></el-icon>
              </div>
              <div class="stat-value">TOP 50</div>
              <div class="stat-label">深海榜单</div>
            </div>
          </el-col>
          <el-col :xs="12" :md="6">
            <div class="stat-card stat-green">
              <div class="stat-glow"></div>
              <div class="stat-icon">
                <el-icon :size="20" color="#34d399"><Timer /></el-icon>
              </div>
              <div class="stat-value">实时</div>
              <div class="stat-label">潮汐更新</div>
            </div>
          </el-col>
        </el-row>
      </section>

      <!-- Leaderboard Table -->
      <section class="table-section" :class="{ mounted }">
        <LeaderboardTable :users="users" :loading="loading" />
      </section>

      <!-- CTA - 海洋主题 -->
      <section class="cta-section" :class="{ mounted }">
        <el-card class="cta-card" shadow="never">
          <!-- Decorative elements -->
          <div class="cta-decoration cta-decoration-1"></div>
          <div class="cta-decoration cta-decoration-2"></div>
          <div class="cta-wave-bg"></div>

          <div class="cta-content">
            <div class="cta-icon">
              <svg viewBox="0 0 32 32" fill="none" width="28" height="28">
                <path d="M4 16C4 16 8 10 14 16C20 22 26 10 28 16" stroke="currentColor" stroke-width="2.5" stroke-linecap="round"/>
              </svg>
            </div>
            <h3 class="cta-title">想要潜入深海排行榜？</h3>
            <p class="cta-text">
              每天坚持训练如潮汐般规律，你也能成为深海王者！
            </p>
            <RouterLink to="/train">
              <el-button type="primary" size="large" round class="cta-btn">
                <el-icon><Timer /></el-icon>
                开始训练，冲击深海
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
  max-width: 900px;
  margin: 0 auto;
  padding: 40px 24px 80px;
  min-height: calc(100vh - 72px);
}

/* ===== Header Section - 海洋主题 ===== */
.header-section {
  text-align: center;
  position: relative;
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth);
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
  width: 400px;
  height: 200px;
  background: radial-gradient(ellipse, rgba(56, 189, 248, 0.2), rgba(251, 191, 36, 0.1), transparent 70%);
  filter: blur(60px);
}

.header-ripple {
  position: absolute;
  left: 50%;
  top: 50%;
  transform: translate(-50%, -50%);
  border-radius: 50%;
  border: 1px solid rgba(56, 189, 248, 0.2);
  animation: ripple-expand 6s ease-out infinite;
}

.header-ripple-1 {
  width: 200px;
  height: 200px;
}

.header-ripple-2 {
  width: 300px;
  height: 300px;
  animation-delay: 2s;
}

@keyframes ripple-expand {
  0% {
    transform: translate(-50%, -50%) scale(0.5);
    opacity: 0.6;
  }
  100% {
    transform: translate(-50%, -50%) scale(2);
    opacity: 0;
  }
}

.trophy-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 100px;
  height: 100px;
  border-radius: var(--radius-2xl);
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  margin-bottom: 28px;
  box-shadow:
    0 20px 40px rgba(14, 165, 233, 0.3),
    0 0 60px rgba(56, 189, 248, 0.2);
  animation: float 4s ease-in-out infinite;
  position: relative;
}

.trophy-svg {
  color: white;
}

@keyframes float {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-12px); }
}

.header-eyebrow {
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 16px;
  margin-bottom: 20px;
}

.header-eyebrow .eyebrow-line {
  width: 50px;
  height: 1px;
  background: linear-gradient(90deg, transparent, rgb(var(--ocean-surface)));
}

.header-eyebrow .eyebrow-line:last-child {
  background: linear-gradient(90deg, rgb(var(--ocean-surface)), transparent);
}

.header-eyebrow .eyebrow-text {
  font-size: 12px;
  font-weight: 600;
  text-transform: uppercase;
  letter-spacing: 3px;
  background: linear-gradient(135deg, rgb(var(--ocean-surface)), rgb(var(--aqua-glow)));
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
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

.gradient-text-ocean {
  background: linear-gradient(135deg, rgb(var(--ocean-surface)), rgb(var(--aqua-glow)), #fbbf24);
  background-size: 200% auto;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: ocean-shimmer 4s ease-in-out infinite;
}

@keyframes ocean-shimmer {
  0%, 100% { background-position: 0% center; }
  50% { background-position: 200% center; }
}

.header-subtitle {
  font-size: 16px;
  color: rgba(255, 255, 255, 0.5);
  max-width: 500px;
  margin: 0 auto;
  line-height: 1.7;
}

/* ===== Podium Section - 海洋主题 ===== */
.podium-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth) 0.1s;
}

.podium-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.podium-wrapper {
  display: flex;
  align-items: flex-end;
  justify-content: center;
  gap: 20px;
  max-width: 550px;
  margin: 0 auto;
  padding: 0 16px;
}

@media (min-width: 768px) {
  .podium-wrapper {
    gap: 36px;
  }
}

.podium-item {
  flex: 1;
  text-align: center;
  position: relative;
}

.podium-item.first {
  margin-top: -40px;
}

.crown-wrapper {
  margin-bottom: 12px;
  position: relative;
  animation: float 2.5s ease-in-out infinite;
}

.crown-glow {
  position: absolute;
  inset: -10px;
  background: radial-gradient(circle, rgba(251, 191, 36, 0.3), transparent 70%);
  filter: blur(10px);
}

.podium-ring {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  width: 80px;
  height: 80px;
  border-radius: 50%;
  border: 2px solid rgba(56, 189, 248, 0.2);
  animation: ring-pulse 3s ease-in-out infinite;
  pointer-events: none;
}

.podium-ring.gold {
  width: 100px;
  height: 100px;
  border-color: rgba(251, 191, 36, 0.3);
}

@keyframes ring-pulse {
  0%, 100% { transform: translate(-50%, -50%) scale(1); opacity: 0.5; }
  50% { transform: translate(-50%, -50%) scale(1.1); opacity: 1; }
}

.podium-avatar {
  width: 68px;
  height: 68px;
  margin: 0 auto 10px;
  border-radius: var(--radius-xl);
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 26px;
  font-weight: 700;
  color: #fff;
  position: relative;
}

@media (min-width: 768px) {
  .podium-avatar {
    width: 80px;
    height: 80px;
    font-size: 32px;
  }
}

.podium-item.first .podium-avatar {
  width: 88px;
  height: 88px;
  font-size: 36px;
}

@media (min-width: 768px) {
  .podium-item.first .podium-avatar {
    width: 100px;
    height: 100px;
    font-size: 42px;
  }
}

.podium-avatar.gold {
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  box-shadow:
    0 15px 35px rgba(251, 191, 36, 0.35),
    0 0 40px rgba(251, 191, 36, 0.2);
}

.podium-avatar.silver {
  background: linear-gradient(135deg, rgb(var(--ocean-surface)), rgb(var(--ocean-shallow)));
  box-shadow: 0 12px 28px rgba(56, 189, 248, 0.3);
}

.podium-avatar.bronze {
  background: linear-gradient(135deg, #d97706, #b45309);
  box-shadow: 0 12px 28px rgba(217, 119, 6, 0.3);
}

.podium-medal {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  margin-bottom: 8px;
}

.podium-medal.gold {
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.2), rgba(245, 158, 11, 0.1));
  color: #fbbf24;
}

.podium-medal.silver {
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.2), rgba(14, 165, 233, 0.1));
  color: rgb(var(--ocean-surface));
}

.podium-medal.bronze {
  background: linear-gradient(135deg, rgba(217, 119, 6, 0.2), rgba(180, 83, 9, 0.1));
  color: #d97706;
}

.podium-name {
  font-size: 14px;
  font-weight: 500;
  color: rgba(255, 255, 255, 0.8);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  padding: 0 8px;
}

.podium-name.gold {
  color: #fbbf24;
  font-weight: 600;
}

.podium-streak {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 4px;
}

.podium-streak.gold {
  color: #fbbf24;
  font-weight: 600;
}

.podium-bar {
  margin-top: 16px;
  border-radius: var(--radius-lg) var(--radius-lg) 0 0;
  position: relative;
  overflow: hidden;
}

.podium-bar::before {
  content: '';
  position: absolute;
  inset: 0;
  background: linear-gradient(0deg, transparent, rgba(255, 255, 255, 0.1));
}

.podium-bar.gold {
  height: 120px;
  background: linear-gradient(to top, rgba(251, 191, 36, 0.15), rgba(251, 191, 36, 0.05));
  border: 1px solid rgba(251, 191, 36, 0.2);
  border-bottom: none;
}

@media (min-width: 768px) {
  .podium-bar.gold {
    height: 140px;
  }
}

.podium-bar.silver {
  height: 88px;
  background: linear-gradient(to top, rgba(56, 189, 248, 0.12), rgba(56, 189, 248, 0.04));
  border: 1px solid rgba(56, 189, 248, 0.15);
  border-bottom: none;
}

@media (min-width: 768px) {
  .podium-bar.silver {
    height: 100px;
  }
}

.podium-bar.bronze {
  height: 68px;
  background: linear-gradient(to top, rgba(217, 119, 6, 0.12), rgba(217, 119, 6, 0.04));
  border: 1px solid rgba(217, 119, 6, 0.15);
  border-bottom: none;
}

@media (min-width: 768px) {
  .podium-bar.bronze {
    height: 80px;
  }
}

/* ===== Stats Section - 海洋主题 ===== */
.stats-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth) 0.2s;
}

.stats-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.stat-card {
  text-align: center;
  padding: 22px 18px;
  border-radius: var(--radius-xl);
  background: var(--glass-bg);
  border: 1px solid rgba(56, 189, 248, 0.1);
  margin-bottom: 16px;
  transition: all 0.4s var(--ease-smooth);
  position: relative;
  overflow: hidden;
}

.stat-glow {
  position: absolute;
  top: -20px;
  left: 50%;
  transform: translateX(-50%);
  width: 80px;
  height: 80px;
  border-radius: 50%;
  filter: blur(30px);
  opacity: 0;
  transition: opacity 0.4s var(--ease-smooth);
  pointer-events: none;
}

.stat-ocean .stat-glow { background: rgba(56, 189, 248, 0.4); }
.stat-amber .stat-glow { background: rgba(251, 191, 36, 0.4); }
.stat-aqua .stat-glow { background: rgba(34, 211, 238, 0.4); }
.stat-green .stat-glow { background: rgba(52, 211, 153, 0.4); }

.stat-card:hover {
  border-color: rgba(56, 189, 248, 0.2);
  transform: translateY(-4px);
}

.stat-card:hover .stat-glow {
  opacity: 1;
}

.stat-card:hover .stat-icon {
  transform: scale(1.1);
}

.stat-icon {
  width: 44px;
  height: 44px;
  margin: 0 auto 12px;
  border-radius: var(--radius-lg);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: transform 0.3s var(--ease-smooth);
  position: relative;
}

.stat-ocean .stat-icon {
  background: linear-gradient(135deg, rgba(56, 189, 248, 0.2), rgba(14, 165, 233, 0.1));
  border: 1px solid rgba(56, 189, 248, 0.2);
}

.stat-amber .stat-icon {
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.2), rgba(245, 158, 11, 0.1));
  border: 1px solid rgba(251, 191, 36, 0.2);
}

.stat-aqua .stat-icon {
  background: linear-gradient(135deg, rgba(34, 211, 238, 0.2), rgba(6, 182, 212, 0.1));
  border: 1px solid rgba(34, 211, 238, 0.2);
}

.stat-green .stat-icon {
  background: linear-gradient(135deg, rgba(52, 211, 153, 0.2), rgba(16, 185, 129, 0.1));
  border: 1px solid rgba(52, 211, 153, 0.2);
}

.stat-value {
  font-size: 26px;
  font-weight: 700;
  color: #fff;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

/* ===== Table Section ===== */
.table-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth) 0.3s;
}

.table-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

/* ===== CTA Section - 海洋主题 ===== */
.cta-section {
  opacity: 0;
  transform: translateY(16px);
  transition: all 0.7s var(--ease-smooth) 0.4s;
}

.cta-section.mounted {
  opacity: 1;
  transform: translateY(0);
}

.cta-card {
  max-width: 680px;
  margin: 0 auto;
  background: var(--glass-bg) !important;
  backdrop-filter: blur(20px);
  border: 1px solid rgba(56, 189, 248, 0.15) !important;
  border-radius: var(--radius-2xl) !important;
  overflow: hidden;
  position: relative;
}

.cta-card :deep(.el-card__body) {
  padding: 48px;
}

.cta-decoration {
  position: absolute;
  border-radius: 50%;
  filter: blur(80px);
  pointer-events: none;
}

.cta-decoration-1 {
  top: -100px;
  left: -100px;
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(56, 189, 248, 0.2), transparent 70%);
}

.cta-decoration-2 {
  bottom: -100px;
  right: -100px;
  width: 200px;
  height: 200px;
  background: radial-gradient(circle, rgba(251, 191, 36, 0.15), transparent 70%);
}

.cta-wave-bg {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 80px;
  background: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1440 100'%3E%3Cpath fill='rgba(56, 189, 248, 0.05)' d='M0,50 C360,100 720,0 1080,50 C1260,75 1380,25 1440,50 L1440,100 L0,100 Z'/%3E%3C/svg%3E");
  background-size: 1440px 100px;
  background-repeat: repeat-x;
  animation: wave-move 12s linear infinite;
  pointer-events: none;
}

@keyframes wave-move {
  0% { background-position-x: 0; }
  100% { background-position-x: 1440px; }
}

.cta-content {
  position: relative;
  text-align: center;
}

.cta-icon {
  width: 64px;
  height: 64px;
  margin: 0 auto 24px;
  border-radius: var(--radius-xl);
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid)));
  display: flex;
  align-items: center;
  justify-content: center;
  color: white;
  box-shadow: 0 15px 35px rgba(14, 165, 233, 0.3);
}

.cta-title {
  font-size: 22px;
  font-weight: 700;
  color: #fff;
  margin: 0 0 12px;
}

.cta-text {
  color: rgba(255, 255, 255, 0.5);
  margin: 0 0 28px;
  font-size: 15px;
}

.cta-btn {
  background: linear-gradient(135deg, rgb(var(--ocean-shallow)), rgb(var(--ocean-mid))) !important;
  border: none !important;
  padding: 14px 40px !important;
  font-size: 16px !important;
  box-shadow: 0 10px 30px rgba(14, 165, 233, 0.3);
  transition: all 0.3s var(--ease-smooth) !important;
}

.cta-btn:hover {
  box-shadow: 0 15px 40px rgba(14, 165, 233, 0.4);
  transform: translateY(-2px);
}

.cta-btn .arrow-icon {
  margin-left: 8px;
  transition: transform 0.2s var(--ease-smooth);
}

.cta-btn:hover .arrow-icon {
  transform: translateX(4px);
}
</style>
