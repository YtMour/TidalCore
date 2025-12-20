<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import LeaderboardTable from '@/components/LeaderboardTable.vue'
import { getLeaderboard, type LeaderboardUser } from '@/api/checkin'
import { Trophy, Timer, TrendCharts, User, Medal, Pointer } from '@element-plus/icons-vue'
import UserAvatar from '@/components/UserAvatar.vue'

const users = ref<LeaderboardUser[]>([])
const loading = ref(true)
const mounted = ref(false)

// 称号配置 - 与 Dashboard 保持一致
const titleConfig = [
  { min: 1000, name: '海神降临', color: '#f472b6', bg: 'linear-gradient(135deg, rgba(244, 114, 182, 0.25), rgba(236, 72, 153, 0.1))', icon: '🔱', effect: 'effect-divine' },
  { min: 730, name: '深渊霸主', color: '#a78bfa', bg: 'linear-gradient(135deg, rgba(167, 139, 250, 0.25), rgba(139, 92, 246, 0.1))', icon: '🦑', effect: 'effect-abyss' },
  { min: 365, name: '深海传奇', color: '#fbbf24', bg: 'linear-gradient(135deg, rgba(251, 191, 36, 0.25), rgba(245, 158, 11, 0.1))', icon: '🌊', effect: 'effect-legend' },
  { min: 180, name: '海洋大师', color: '#38bdf8', bg: 'linear-gradient(135deg, rgba(56, 189, 248, 0.25), rgba(14, 165, 233, 0.1))', icon: '🐋', effect: 'effect-master' },
  { min: 90, name: '浪潮专家', color: '#22d3ee', bg: 'linear-gradient(135deg, rgba(34, 211, 238, 0.25), rgba(6, 182, 212, 0.1))', icon: '🐬', effect: 'effect-expert' },
  { min: 30, name: '潮汐进阶', color: '#34d399', bg: 'linear-gradient(135deg, rgba(52, 211, 153, 0.25), rgba(16, 185, 129, 0.1))', icon: '🐠', effect: 'effect-advanced' },
  { min: 7, name: '入海新手', color: '#0ea5e9', bg: 'linear-gradient(135deg, rgba(14, 165, 233, 0.25), rgba(2, 132, 199, 0.1))', icon: '🐟', effect: 'effect-beginner' },
  { min: 0, name: '初探海域', color: 'rgba(255, 255, 255, 0.6)', bg: 'linear-gradient(135deg, rgba(255, 255, 255, 0.1), rgba(255, 255, 255, 0.05))', icon: '🐚', effect: '' }
]

// 根据用户获取称号信息
function getUserTitle(user: LeaderboardUser | undefined) {
  if (!user) {
    return titleConfig[titleConfig.length - 1]!
  }
  // 优先使用数据库中存储的称号
  if (user.title) {
    const customTitle = titleConfig.find(t => t.name === user.title)
    if (customTitle) return customTitle
  }
  // 否则根据打卡次数自动计算
  const total = user.total_checkin || 0
  for (const title of titleConfig) {
    if (total >= title.min) {
      return title
    }
  }
  return titleConfig[titleConfig.length - 1]!
}

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
          <svg viewBox="0 0 60 60" fill="none" width="60" height="60" class="trophy-svg">
            <defs>
              <linearGradient id="trophyGold" x1="0%" y1="0%" x2="100%" y2="100%">
                <stop offset="0%" stop-color="#fcd34d" />
                <stop offset="50%" stop-color="#f59e0b" />
                <stop offset="100%" stop-color="#d97706" />
              </linearGradient>
              <linearGradient id="trophyWater" x1="0%" y1="100%" x2="0%" y2="0%">
                <stop offset="0%" stop-color="#0284c7" />
                <stop offset="100%" stop-color="#38bdf8" />
              </linearGradient>
              <clipPath id="trophyCup">
                <path d="M18 14h24v6c0 8-5 15-12 18c-7-3-12-10-12-18v-6z"/>
              </clipPath>
            </defs>
            <!-- 奖杯主体 -->
            <path d="M18 14h24v6c0 8-5 15-12 18c-7-3-12-10-12-18v-6z" fill="url(#trophyGold)" opacity="0.2"/>
            <!-- 海浪填充效果 -->
            <g clip-path="url(#trophyCup)">
              <path class="trophy-wave-1" d="M10 26 Q20 22, 30 26 T50 26 L50 50 L10 50 Z" fill="url(#trophyWater)" opacity="0.9"/>
              <path class="trophy-wave-2" d="M10 30 Q20 26, 30 30 T50 30 L50 50 L10 50 Z" fill="#22d3ee" opacity="0.6"/>
              <path class="trophy-wave-3" d="M10 34 Q20 30, 30 34 T50 34 L50 50 L10 50 Z" fill="#06b6d4" opacity="0.4"/>
            </g>
            <!-- 奖杯边框 -->
            <path d="M18 14h24v6c0 8-5 15-12 18c-7-3-12-10-12-18v-6z" stroke="url(#trophyGold)" stroke-width="2" fill="none"/>
            <!-- 把手 -->
            <path d="M16 16c-4 0-6 3-6 6s2 6 6 6" stroke="url(#trophyGold)" stroke-width="2" fill="none" stroke-linecap="round"/>
            <path d="M44 16c4 0 6 3 6 6s-2 6-6 6" stroke="url(#trophyGold)" stroke-width="2" fill="none" stroke-linecap="round"/>
            <!-- 底座 -->
            <rect x="24" y="38" width="12" height="4" rx="2" fill="url(#trophyGold)" opacity="0.8"/>
            <rect x="20" y="42" width="20" height="4" rx="2" fill="url(#trophyGold)" opacity="0.6"/>
            <!-- 闪光 -->
            <circle cx="24" cy="18" r="2" fill="white" opacity="0.5"/>
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

      <!-- Top 3 Podium (when data loaded) - 全新海洋主题设计 -->
      <section
        v-if="!loading && top3Users.length >= 3"
        class="podium-section"
        :class="{ mounted }"
      >
        <div class="podium-wrapper">
          <!-- 2nd Place - 银色浪花 -->
          <div class="podium-item second">
            <div class="podium-card">
              <div class="podium-rank-badge silver">2</div>
              <div class="podium-glow silver"></div>
              <div class="podium-avatar-wrapper">
                <UserAvatar
                  :user-id="top3Users[1]?.id"
                  :username="top3Users[1]?.username"
                  :size="64"
                  :border-radius="16"
                  variant="silver"
                  class="podium-avatar silver"
                />
                <div class="podium-medal-icon silver">
                  <el-icon :size="16"><Medal /></el-icon>
                </div>
              </div>
              <div class="podium-info">
                <div class="podium-name">{{ top3Users[1]?.display_name || top3Users[1]?.username }}</div>
                <div class="podium-title" :class="getUserTitle(top3Users[1]).effect" :style="{ color: getUserTitle(top3Users[1]).color }">
                  <span class="title-icon">{{ getUserTitle(top3Users[1]).icon }}</span>
                  <span class="title-name">{{ getUserTitle(top3Users[1]).name }}</span>
                </div>
                <div class="podium-stats">
                  <div class="stat-main">
                    <span class="stat-num">{{ top3Users[1]?.streak }}</span>
                    <span class="stat-unit">天连续</span>
                  </div>
                  <div class="stat-sub">
                    <span>最高 {{ top3Users[1]?.max_streak }}</span>
                    <span class="stat-dot">·</span>
                    <span>累计 {{ top3Users[1]?.total_checkin }}</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="podium-bar silver">
              <div class="bar-wave"></div>
              <div class="bar-rank">2</div>
            </div>
          </div>

          <!-- 1st Place - 金色深海王者 -->
          <div class="podium-item first">
            <div class="podium-card champion">
              <div class="podium-rank-badge gold">1</div>
              <div class="podium-glow gold"></div>
              <div class="champion-particles">
                <span v-for="i in 6" :key="i" class="particle"></span>
              </div>
              <div class="podium-avatar-wrapper">
                <UserAvatar
                  :user-id="top3Users[0]?.id"
                  :username="top3Users[0]?.username"
                  :size="80"
                  :border-radius="20"
                  variant="gold"
                  class="podium-avatar gold"
                />
                <div class="podium-medal-icon gold">
                  <el-icon :size="18"><Trophy /></el-icon>
                </div>
              </div>
              <div class="podium-info">
                <div class="podium-name champion-name">{{ top3Users[0]?.display_name || top3Users[0]?.username }}</div>
                <div class="podium-title champion-title" :class="getUserTitle(top3Users[0]).effect" :style="{ color: getUserTitle(top3Users[0]).color }">
                  <span class="title-icon">{{ getUserTitle(top3Users[0]).icon }}</span>
                  <span class="title-name">{{ getUserTitle(top3Users[0]).name }}</span>
                </div>
                <div class="podium-stats">
                  <div class="stat-main champion-stat">
                    <span class="stat-num">{{ top3Users[0]?.streak }}</span>
                    <span class="stat-unit">天连续</span>
                  </div>
                  <div class="stat-sub">
                    <span>最高 {{ top3Users[0]?.max_streak }}</span>
                    <span class="stat-dot">·</span>
                    <span>累计 {{ top3Users[0]?.total_checkin }}</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="podium-bar gold">
              <div class="bar-wave"></div>
              <div class="bar-shine"></div>
              <div class="bar-rank">1</div>
            </div>
          </div>

          <!-- 3rd Place - 铜色珊瑚 -->
          <div class="podium-item third">
            <div class="podium-card">
              <div class="podium-rank-badge bronze">3</div>
              <div class="podium-glow bronze"></div>
              <div class="podium-avatar-wrapper">
                <UserAvatar
                  :user-id="top3Users[2]?.id"
                  :username="top3Users[2]?.username"
                  :size="64"
                  :border-radius="16"
                  variant="bronze"
                  class="podium-avatar bronze"
                />
                <div class="podium-medal-icon bronze">
                  <el-icon :size="16"><Medal /></el-icon>
                </div>
              </div>
              <div class="podium-info">
                <div class="podium-name">{{ top3Users[2]?.display_name || top3Users[2]?.username }}</div>
                <div class="podium-title" :class="getUserTitle(top3Users[2]).effect" :style="{ color: getUserTitle(top3Users[2]).color }">
                  <span class="title-icon">{{ getUserTitle(top3Users[2]).icon }}</span>
                  <span class="title-name">{{ getUserTitle(top3Users[2]).name }}</span>
                </div>
                <div class="podium-stats">
                  <div class="stat-main">
                    <span class="stat-num">{{ top3Users[2]?.streak }}</span>
                    <span class="stat-unit">天连续</span>
                  </div>
                  <div class="stat-sub">
                    <span>最高 {{ top3Users[2]?.max_streak }}</span>
                    <span class="stat-dot">·</span>
                    <span>累计 {{ top3Users[2]?.total_checkin }}</span>
                  </div>
                </div>
              </div>
            </div>
            <div class="podium-bar bronze">
              <div class="bar-wave"></div>
              <div class="bar-rank">3</div>
            </div>
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
              <svg viewBox="0 0 60 60" fill="none" class="cta-wave-svg">
                <defs>
                  <linearGradient id="ctaLeaderGrad" x1="0%" y1="100%" x2="0%" y2="0%">
                    <stop offset="0%" stop-color="#0284c7" />
                    <stop offset="100%" stop-color="#38bdf8" />
                  </linearGradient>
                  <clipPath id="ctaLeaderClip">
                    <rect x="4" y="4" width="52" height="52" rx="12" />
                  </clipPath>
                </defs>
                <!-- 海浪涌动填充 -->
                <g clip-path="url(#ctaLeaderClip)">
                  <rect x="4" y="4" width="52" height="52" rx="12" fill="rgba(255,255,255,0.1)"/>
                  <path class="cta-wave-fill cta-wave-1" d="M0 32 Q15 26, 30 32 T60 32 L60 60 L0 60 Z" fill="url(#ctaLeaderGrad)" opacity="0.8"/>
                  <path class="cta-wave-fill cta-wave-2" d="M0 38 Q15 32, 30 38 T60 38 L60 60 L0 60 Z" fill="#22d3ee" opacity="0.5"/>
                  <path class="cta-wave-fill cta-wave-3" d="M0 44 Q15 38, 30 44 T60 44 L60 60 L0 60 Z" fill="#0ea5e9" opacity="0.3"/>
                </g>
                <!-- 波浪线条 -->
                <path class="cta-wave-line" d="M8 28 Q22 20, 36 28 T52 28" stroke="white" stroke-width="2.5" stroke-linecap="round" fill="none" opacity="0.9"/>
              </svg>
            </div>
            <h3 class="cta-title">想要潜入深海排行榜？</h3>
            <p class="cta-text">
              每天坚持训练如潮汐般规律，你也能成为深海王者！
            </p>
            <RouterLink to="/train" class="cta-btn-wrapper">
              <button class="cta-btn-ocean">
                <span class="btn-bg"></span>
                <span class="btn-wave"></span>
                <span class="btn-content">
                  <el-icon><Timer /></el-icon>
                  <span>开始训练，冲击深海</span>
                  <svg viewBox="0 0 24 24" fill="none" class="btn-arrow">
                    <path d="M5 12H19M19 12L12 5M19 12L12 19" stroke="currentColor" stroke-width="2" stroke-linecap="round" stroke-linejoin="round"/>
                  </svg>
                </span>
              </button>
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
  filter: drop-shadow(0 0 15px rgba(251, 191, 36, 0.4));
}

/* 奖杯内海浪涌动动画 */
.trophy-svg .trophy-wave-1 {
  animation: trophy-surge-1 3s ease-in-out infinite;
}

.trophy-svg .trophy-wave-2 {
  animation: trophy-surge-2 3.5s ease-in-out infinite;
}

.trophy-svg .trophy-wave-3 {
  animation: trophy-surge-3 4s ease-in-out infinite;
}

@keyframes trophy-surge-1 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

@keyframes trophy-surge-2 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-3px); }
}

@keyframes trophy-surge-3 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-2px); }
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
  background: linear-gradient(
    90deg,
    #38bdf8 0%,
    #5cc8f5 10%,
    #22d3ee 20%,
    #60d0e8 30%,
    #fbbf24 40%,
    #f5c842 50%,
    #fbbf24 60%,
    #60d0e8 70%,
    #22d3ee 80%,
    #5cc8f5 90%,
    #38bdf8 100%
  );
  background-size: 300% 100%;
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  animation: ocean-shimmer 12s ease-in-out infinite;
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

/* ===== Podium Section - 全新海洋主题设计 ===== */
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
  gap: 16px;
  max-width: 800px;
  margin: 0 auto;
  padding: 0 16px;
}

@media (min-width: 768px) {
  .podium-wrapper {
    gap: 24px;
  }
}

.podium-item {
  flex: 1;
  max-width: 240px;
  display: flex;
  flex-direction: column;
  align-items: center;
}

.podium-item.first {
  order: 2;
}

.podium-item.second {
  order: 1;
}

.podium-item.third {
  order: 3;
}

/* Podium Card */
.podium-card {
  width: 100%;
  padding: 20px 16px;
  border-radius: 20px;
  background: var(--glass-bg);
  border: 1px solid rgba(255, 255, 255, 0.08);
  position: relative;
  overflow: hidden;
  transition: all 0.4s var(--ease-smooth);
  backdrop-filter: blur(20px);
  display: flex;
  flex-direction: column;
  align-items: center;
}

.podium-card:hover {
  transform: translateY(-4px);
  border-color: rgba(255, 255, 255, 0.15);
}

/* 第二名悬停效果 - 银色 */
.podium-item.second .podium-card:hover {
  border-color: rgba(148, 163, 184, 0.3);
  box-shadow: 0 20px 50px rgba(148, 163, 184, 0.15);
}

/* 第三名悬停效果 - 铜色 */
.podium-item.third .podium-card:hover {
  border-color: rgba(217, 119, 6, 0.3);
  box-shadow: 0 20px 50px rgba(217, 119, 6, 0.15);
}

.podium-card.champion {
  padding: 28px 20px;
  background: linear-gradient(180deg, rgba(251, 191, 36, 0.06), var(--glass-bg));
  border-color: rgba(251, 191, 36, 0.15);
}

.podium-card.champion:hover {
  border-color: rgba(251, 191, 36, 0.35);
  box-shadow: 0 20px 50px rgba(251, 191, 36, 0.2);
}

/* Rank Badge */
.podium-rank-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  width: 28px;
  height: 28px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: 700;
  z-index: 2;
}

.podium-rank-badge.gold {
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  color: #1a1a2e;
  box-shadow: 0 4px 12px rgba(251, 191, 36, 0.4);
}

.podium-rank-badge.silver {
  background: linear-gradient(135deg, #94a3b8, #64748b);
  color: #1a1a2e;
  box-shadow: 0 4px 12px rgba(148, 163, 184, 0.3);
}

.podium-rank-badge.bronze {
  background: linear-gradient(135deg, #d97706, #b45309);
  color: #fff;
  box-shadow: 0 4px 12px rgba(217, 119, 6, 0.3);
}

/* Podium Glow */
.podium-glow {
  position: absolute;
  top: -50%;
  left: 50%;
  transform: translateX(-50%);
  width: 150%;
  height: 100%;
  border-radius: 50%;
  filter: blur(40px);
  opacity: 0.5;
  pointer-events: none;
}

.podium-glow.gold {
  background: radial-gradient(circle, rgba(251, 191, 36, 0.3), transparent 70%);
}

.podium-glow.silver {
  background: radial-gradient(circle, rgba(148, 163, 184, 0.25), transparent 70%);
}

.podium-glow.bronze {
  background: radial-gradient(circle, rgba(217, 119, 6, 0.25), transparent 70%);
}

/* Champion Particles */
.champion-particles {
  position: absolute;
  inset: 0;
  overflow: hidden;
  pointer-events: none;
}

.champion-particles .particle {
  position: absolute;
  width: 4px;
  height: 4px;
  background: #fbbf24;
  border-radius: 50%;
  opacity: 0;
  animation: particle-rise 3s ease-out infinite;
}

.champion-particles .particle:nth-child(1) { left: 20%; animation-delay: 0s; }
.champion-particles .particle:nth-child(2) { left: 40%; animation-delay: 0.5s; }
.champion-particles .particle:nth-child(3) { left: 60%; animation-delay: 1s; }
.champion-particles .particle:nth-child(4) { left: 80%; animation-delay: 1.5s; }
.champion-particles .particle:nth-child(5) { left: 30%; animation-delay: 2s; }
.champion-particles .particle:nth-child(6) { left: 70%; animation-delay: 2.5s; }

@keyframes particle-rise {
  0% {
    bottom: 0;
    opacity: 0;
    transform: scale(0);
  }
  20% {
    opacity: 1;
    transform: scale(1);
  }
  100% {
    bottom: 100%;
    opacity: 0;
    transform: scale(0.5);
  }
}

/* Avatar Wrapper */
.podium-avatar-wrapper {
  position: relative;
  margin-bottom: 12px;
}

.podium-avatar {
  width: 64px;
  height: 64px;
  border-radius: 16px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 24px;
  font-weight: 700;
  color: #fff;
  position: relative;
  z-index: 1;
}

.podium-item.first .podium-avatar {
  width: 80px;
  height: 80px;
  font-size: 32px;
  border-radius: 20px;
}

.podium-avatar.gold {
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  box-shadow: 0 12px 30px rgba(251, 191, 36, 0.4);
}

.podium-avatar.silver {
  background: linear-gradient(135deg, #c0c7d0, #94a3b8);
  box-shadow: 0 10px 25px rgba(148, 163, 184, 0.35);
}

.podium-avatar.bronze {
  background: linear-gradient(135deg, #d97706, #b45309);
  box-shadow: 0 10px 25px rgba(217, 119, 6, 0.35);
}

/* Medal Icon */
.podium-medal-icon {
  position: absolute;
  bottom: -6px;
  right: -6px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 2;
  border: 2px solid rgba(30, 30, 46, 0.8);
}

.podium-item.first .podium-medal-icon {
  width: 32px;
  height: 32px;
  bottom: -8px;
  right: -8px;
}

.podium-medal-icon.gold {
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  color: #1a1a2e;
}

.podium-medal-icon.silver {
  background: linear-gradient(135deg, #94a3b8, #64748b);
  color: #1a1a2e;
}

.podium-medal-icon.bronze {
  background: linear-gradient(135deg, #d97706, #b45309);
  color: #fff;
}

/* Podium Info */
.podium-info {
  text-align: center;
  position: relative;
  z-index: 1;
}

.podium-name {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.9);
  margin-bottom: 6px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  max-width: 100%;
}

.podium-name.champion-name {
  font-size: 16px;
  color: #fbbf24;
}

/* Podium Title (称号) */
.podium-title {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 4px 10px;
  border-radius: 20px;
  background: rgba(255, 255, 255, 0.06);
  font-size: 12px;
  font-weight: 500;
  margin-bottom: 12px;
}

.podium-title.champion-title {
  background: rgba(251, 191, 36, 0.1);
  padding: 5px 12px;
  font-size: 13px;
}

.podium-title .title-icon {
  font-size: 14px;
}

.podium-title.champion-title .title-icon {
  font-size: 16px;
}

/* Podium Stats */
.podium-stats {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
}

.podium-stats .stat-main {
  display: flex;
  align-items: baseline;
  gap: 4px;
}

.podium-stats .stat-num {
  font-size: 28px;
  font-weight: 800;
  color: #fff;
  line-height: 1;
}

.podium-item.first .podium-stats .stat-num {
  font-size: 36px;
  background: linear-gradient(135deg, #fbbf24, #f59e0b);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.podium-stats .stat-unit {
  font-size: 12px;
  color: rgba(255, 255, 255, 0.5);
}

.podium-stats .stat-sub {
  display: flex;
  align-items: center;
  gap: 6px;
  font-size: 11px;
  color: rgba(255, 255, 255, 0.4);
}

.podium-stats .stat-dot {
  opacity: 0.5;
}

/* Podium Bar */
.podium-bar {
  width: 100%;
  margin-top: 12px;
  border-radius: 12px 12px 0 0;
  position: relative;
  overflow: hidden;
  display: flex;
  align-items: flex-end;
  justify-content: center;
}

.podium-bar .bar-wave {
  position: absolute;
  bottom: 0;
  left: 0;
  right: 0;
  height: 30%;
  background: linear-gradient(180deg, transparent, rgba(255, 255, 255, 0.05));
  animation: bar-wave-move 3s ease-in-out infinite;
}

@keyframes bar-wave-move {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-5px); }
}

.podium-bar .bar-shine {
  position: absolute;
  top: 0;
  left: -100%;
  width: 50%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.1), transparent);
  animation: bar-shine-move 4s ease-in-out infinite;
}

@keyframes bar-shine-move {
  0% { left: -100%; }
  50%, 100% { left: 200%; }
}

.podium-bar .bar-rank {
  position: relative;
  font-size: 48px;
  font-weight: 900;
  opacity: 0.5;
  line-height: 1;
  padding-bottom: 8px;
  text-shadow: 0 2px 8px rgba(0, 0, 0, 0.3);
}

.podium-bar.gold {
  height: 100px;
  background: linear-gradient(180deg, rgba(251, 191, 36, 0.66), rgba(251, 191, 36, 0.4));
  border: 1px solid rgba(251, 191, 36, 0.5);
  border-bottom: none;
}

.podium-bar.gold .bar-rank {
  color: #fbbf24;
}

.podium-bar.silver {
  height: 70px;
  background: linear-gradient(180deg, rgba(148, 163, 184, 0.66), rgba(148, 163, 184, 0.4));
  border: 1px solid rgba(148, 163, 184, 0.5);
  border-bottom: none;
}

.podium-bar.silver .bar-rank {
  color: #94a3b8;
}

.podium-bar.bronze {
  height: 50px;
  background: linear-gradient(180deg, rgba(217, 119, 6, 0.66), rgba(217, 119, 6, 0.4));
  border: 1px solid rgba(217, 119, 6, 0.5);
  border-bottom: none;
}

.podium-bar.bronze .bar-rank {
  color: #d97706;
}

@media (min-width: 768px) {
  .podium-bar.gold { height: 120px; }
  .podium-bar.silver { height: 85px; }
  .podium-bar.bronze { height: 60px; }
  .podium-bar .bar-rank { font-size: 56px; }
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
  width: 80px;
  height: 80px;
  margin: 0 auto 24px;
  border-radius: var(--radius-xl);
  background: linear-gradient(180deg, #0c4a6e, #0369a1);
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 15px 35px rgba(14, 165, 233, 0.3);
  overflow: hidden;
  position: relative;
}

.cta-wave-svg {
  width: 100%;
  height: 100%;
  position: absolute;
  inset: 0;
}

/* CTA海浪涌动动画 */
.cta-wave-svg .cta-wave-1 {
  animation: cta-surge-1 3s ease-in-out infinite;
}

.cta-wave-svg .cta-wave-2 {
  animation: cta-surge-2 3.5s ease-in-out infinite;
}

.cta-wave-svg .cta-wave-3 {
  animation: cta-surge-3 4s ease-in-out infinite;
}

.cta-wave-svg .cta-wave-line {
  animation: cta-line-surge 2.5s ease-in-out infinite;
}

@keyframes cta-surge-1 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-6px); }
}

@keyframes cta-surge-2 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-4px); }
}

@keyframes cta-surge-3 {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-2px); }
}

@keyframes cta-line-surge {
  0%, 100% {
    transform: translateY(0);
    d: path("M8 28 Q22 20, 36 28 T52 28");
  }
  50% {
    transform: translateY(-2px);
    d: path("M8 28 Q22 36, 36 28 T52 28");
  }
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

.cta-btn-wrapper {
  text-decoration: none;
  display: inline-block;
}

.cta-btn-ocean {
  position: relative;
  padding: 16px 36px;
  border: none;
  border-radius: var(--radius-full);
  cursor: pointer;
  overflow: hidden;
  background: transparent;
  isolation: isolate;
}

.cta-btn-ocean .btn-bg {
  position: absolute;
  inset: 0;
  background: linear-gradient(135deg, #0ea5e9, #0284c7);
  border-radius: inherit;
  z-index: -2;
  transition: transform 0.3s var(--ease-smooth);
}

.cta-btn-ocean .btn-wave {
  position: absolute;
  inset: 0;
  border-radius: inherit;
  z-index: -1;
  overflow: hidden;
}

.cta-btn-ocean .btn-wave::before {
  content: '';
  position: absolute;
  bottom: -100%;
  left: -10%;
  right: -10%;
  height: 200%;
  background: linear-gradient(180deg, transparent 45%, rgba(56, 189, 248, 0.4) 50%, rgba(34, 211, 238, 0.6) 55%, rgba(56, 189, 248, 0.3) 60%, transparent 65%);
  animation: btn-wave-rise 3s ease-in-out infinite;
}

@keyframes btn-wave-rise {
  0%, 100% { transform: translateY(25%); }
  50% { transform: translateY(15%); }
}

.cta-btn-ocean .btn-content {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  color: white;
  font-size: 16px;
  font-weight: 600;
  z-index: 1;
}

.cta-btn-ocean .btn-arrow {
  width: 18px;
  height: 18px;
  transition: transform 0.3s var(--ease-smooth);
}

.cta-btn-ocean:hover .btn-arrow {
  transform: translateX(4px);
}

.cta-btn-ocean:hover .btn-bg {
  transform: scale(1.02);
}

.cta-btn-ocean::after {
  content: '';
  position: absolute;
  inset: 0;
  border-radius: inherit;
  box-shadow: 0 8px 32px rgba(14, 165, 233, 0.4);
  transition: box-shadow 0.3s var(--ease-smooth);
}

.cta-btn-ocean:hover::after {
  box-shadow: 0 12px 40px rgba(14, 165, 233, 0.5);
}

/* ===== Podium Title Effects (称号特效) ===== */
/* 海神降临 - 神圣光芒 + 彩虹流光 */
.podium-title.effect-divine {
  position: relative;
  animation: divine-glow 2s ease-in-out infinite;
  box-shadow: 0 0 10px rgba(244, 114, 182, 0.4), 0 0 20px rgba(244, 114, 182, 0.2);
}

.podium-title.effect-divine::before {
  content: '';
  position: absolute;
  inset: -1px;
  border-radius: inherit;
  background: linear-gradient(90deg, #f472b6, #a78bfa, #38bdf8, #34d399, #fbbf24, #f472b6);
  background-size: 300% 100%;
  animation: rainbow-flow 3s linear infinite;
  z-index: -1;
  opacity: 0.7;
}

.podium-title.effect-divine .title-icon {
  animation: divine-icon 1.5s ease-in-out infinite;
}

@keyframes divine-glow {
  0%, 100% { box-shadow: 0 0 10px rgba(244, 114, 182, 0.4), 0 0 20px rgba(244, 114, 182, 0.2); }
  50% { box-shadow: 0 0 15px rgba(244, 114, 182, 0.6), 0 0 30px rgba(244, 114, 182, 0.3); }
}

@keyframes rainbow-flow {
  0% { background-position: 0% 50%; }
  100% { background-position: 300% 50%; }
}

@keyframes divine-icon {
  0%, 100% { transform: scale(1) rotate(0deg); }
  50% { transform: scale(1.15) rotate(8deg); }
}

/* 深渊霸主 - 暗紫脉动 */
.podium-title.effect-abyss {
  position: relative;
  animation: abyss-pulse 2.5s ease-in-out infinite;
  box-shadow: 0 0 8px rgba(167, 139, 250, 0.4), inset 0 0 6px rgba(139, 92, 246, 0.2);
}

.podium-title.effect-abyss .title-icon {
  animation: abyss-icon 2s ease-in-out infinite;
}

@keyframes abyss-pulse {
  0%, 100% { box-shadow: 0 0 8px rgba(167, 139, 250, 0.4), inset 0 0 6px rgba(139, 92, 246, 0.2); }
  50% { box-shadow: 0 0 15px rgba(167, 139, 250, 0.5), inset 0 0 10px rgba(139, 92, 246, 0.3); }
}

@keyframes abyss-icon {
  0%, 100% { transform: scale(1); }
  25% { transform: scale(1.1) rotate(-5deg); }
  75% { transform: scale(1.1) rotate(5deg); }
}

/* 深海传奇 - 金色闪耀 */
.podium-title.effect-legend {
  position: relative;
  overflow: hidden;
  animation: legend-shine 2s ease-in-out infinite;
  box-shadow: 0 0 8px rgba(251, 191, 36, 0.3);
}

.podium-title.effect-legend::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 50%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  animation: legend-sweep 2.5s ease-in-out infinite;
}

.podium-title.effect-legend .title-icon {
  animation: legend-wave 1.5s ease-in-out infinite;
}

@keyframes legend-shine {
  0%, 100% { box-shadow: 0 0 8px rgba(251, 191, 36, 0.3); }
  50% { box-shadow: 0 0 15px rgba(251, 191, 36, 0.5), 0 0 25px rgba(251, 191, 36, 0.2); }
}

@keyframes legend-sweep {
  0% { left: -100%; }
  50%, 100% { left: 200%; }
}

@keyframes legend-wave {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-2px); }
}

/* 海洋大师 - 水波纹扩散 */
.podium-title.effect-master {
  position: relative;
  animation: master-glow 2s ease-in-out infinite;
}

.podium-title.effect-master::after {
  content: '';
  position: absolute;
  inset: -3px;
  border-radius: inherit;
  border: 1px solid rgba(56, 189, 248, 0.3);
  animation: master-ripple 2s ease-out infinite;
  pointer-events: none;
}

.podium-title.effect-master .title-icon {
  animation: master-swim 3s ease-in-out infinite;
}

@keyframes master-glow {
  0%, 100% { box-shadow: 0 0 6px rgba(56, 189, 248, 0.3); }
  50% { box-shadow: 0 0 12px rgba(56, 189, 248, 0.4); }
}

@keyframes master-ripple {
  0% { transform: scale(1); opacity: 0.5; }
  100% { transform: scale(1.2); opacity: 0; }
}

@keyframes master-swim {
  0%, 100% { transform: translateX(0); }
  50% { transform: translateX(2px); }
}

/* 浪潮专家 - 轻微波动 */
.podium-title.effect-expert {
  animation: expert-float 2.5s ease-in-out infinite;
}

.podium-title.effect-expert .title-icon {
  animation: expert-jump 2s ease-in-out infinite;
}

@keyframes expert-float {
  0%, 100% { transform: translateY(0); box-shadow: 0 0 5px rgba(34, 211, 238, 0.2); }
  50% { transform: translateY(-1px); box-shadow: 0 0 10px rgba(34, 211, 238, 0.3); }
}

@keyframes expert-jump {
  0%, 100% { transform: translateY(0); }
  50% { transform: translateY(-1px); }
}

/* 潮汐进阶 - 柔和呼吸 */
.podium-title.effect-advanced {
  animation: advanced-breathe 3s ease-in-out infinite;
}

.podium-title.effect-advanced .title-icon {
  animation: advanced-swim 2.5s ease-in-out infinite;
}

@keyframes advanced-breathe {
  0%, 100% { opacity: 0.9; box-shadow: 0 0 4px rgba(52, 211, 153, 0.2); }
  50% { opacity: 1; box-shadow: 0 0 8px rgba(52, 211, 153, 0.3); }
}

@keyframes advanced-swim {
  0%, 100% { transform: translateX(0); }
  50% { transform: translateX(1px); }
}

/* 入海新手 - 轻微闪烁 */
.podium-title.effect-beginner {
  animation: beginner-twinkle 3s ease-in-out infinite;
}

@keyframes beginner-twinkle {
  0%, 100% { opacity: 0.85; }
  50% { opacity: 1; }
}
</style>
