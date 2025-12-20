<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { LeaderboardUser } from '@/api/checkin'
import { Skeleton } from '@/components/ui'

defineProps<{
  users: LeaderboardUser[]
  loading?: boolean
}>()

const mounted = ref(false)

onMounted(() => {
  setTimeout(() => {
    mounted.value = true
  }, 100)
})

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
function getUserTitle(user: LeaderboardUser) {
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

function getMedalEmoji(index: number): string {
  const medals = ['🥇', '🥈', '🥉']
  return medals[index] || ''
}

function getMedalBg(index: number): string {
  const bgs = [
    'linear-gradient(135deg, #fbbf24, #f59e0b)',
    'linear-gradient(135deg, #94a3b8, #64748b)',
    'linear-gradient(135deg, #d97706, #b45309)'
  ]
  return bgs[index] || 'rgba(255, 255, 255, 0.08)'
}

// 获取显示名称（优先使用 display_name）
function getDisplayName(user: LeaderboardUser): string {
  return user.display_name || user.username || '未知用户'
}

function getStreakColor(streak: number): string {
  if (streak >= 30) return '#fb923c'
  if (streak >= 14) return '#fbbf24'
  if (streak >= 7) return '#facc15'
  return 'inherit'
}
</script>

<template>
  <div class="leaderboard-container">
    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div v-for="i in 5" :key="i" class="loading-row">
        <Skeleton width="2rem" height="2rem" rounded="full" />
        <Skeleton width="40%" height="1.25rem" />
        <Skeleton width="4rem" height="1.25rem" class="ml-auto" />
      </div>
    </div>

    <!-- Empty State -->
    <div v-else-if="users.length === 0" class="empty-state">
      <div class="empty-icon">🏆</div>
      <p class="empty-title">暂无排行数据</p>
      <p class="empty-desc">成为第一个打卡的人吧!</p>
    </div>

    <!-- Leaderboard List -->
    <div v-else class="leaderboard-list">
      <div
        v-for="(user, index) in users"
        :key="user.id"
        class="user-card"
        :class="{ 'top-three': index < 3, [`rank-${index + 1}`]: index < 3 }"
        :style="{ animationDelay: mounted ? `${index * 40}ms` : '0ms' }"
      >
        <!-- Rank Badge -->
        <div class="rank-badge" :class="{ 'is-medal': index < 3 }">
          <template v-if="index < 3">
            <span class="medal-emoji">{{ getMedalEmoji(index) }}</span>
          </template>
          <template v-else>
            <span class="rank-num">{{ index + 1 }}</span>
          </template>
        </div>

        <!-- User Avatar & Name -->
        <div class="user-info">
          <div
            class="user-avatar"
            :style="{ background: index < 3 ? getMedalBg(index) : undefined }"
          >
            {{ getDisplayName(user)?.[0]?.toUpperCase() || '?' }}
          </div>
          <span class="user-name">{{ getDisplayName(user) }}</span>
          <span
            class="user-title-badge"
            :class="getUserTitle(user).effect"
            :style="{ color: getUserTitle(user).color, borderColor: getUserTitle(user).color + '40' }"
          >
            <span class="title-icon">{{ getUserTitle(user).icon }}</span>
            <span class="title-text">{{ getUserTitle(user).name }}</span>
          </span>
        </div>

        <!-- Stats -->
        <div class="user-stats">
          <div class="stat-item main-stat">
            <span class="stat-value" :style="{ color: getStreakColor(user.streak) }">
              {{ user.streak }}
            </span>
            <span class="stat-label">连续</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-value">{{ user.max_streak }}</span>
            <span class="stat-label">最高</span>
          </div>
          <div class="stat-divider"></div>
          <div class="stat-item">
            <span class="stat-value">{{ user.total_checkin }}</span>
            <span class="stat-label">总计</span>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style scoped>
.leaderboard-container {
  background: var(--card-bg, rgba(30, 30, 46, 0.6));
  border: 1px solid var(--card-border, rgba(255, 255, 255, 0.06));
  border-radius: 12px;
  overflow: hidden;
}

/* Loading State */
.loading-state {
  padding: 20px;
}

.loading-row {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 12px 0;
}

.loading-row + .loading-row {
  border-top: 1px solid rgba(255, 255, 255, 0.04);
}

/* Empty State */
.empty-state {
  padding: 48px 24px;
  text-align: center;
}

.empty-icon {
  font-size: 48px;
  margin-bottom: 16px;
  opacity: 0.5;
}

.empty-title {
  font-size: 16px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
  margin: 0 0 8px;
}

.empty-desc {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.55);
  margin: 0;
}

/* Leaderboard List */
.leaderboard-list {
  padding: 8px;
}

.user-card {
  display: flex;
  align-items: center;
  gap: 12px;
  padding: 12px 16px;
  border-radius: 10px;
  background: transparent;
  transition: all 0.2s ease;
  animation: fadeIn 0.4s ease forwards;
  opacity: 0;
}

@keyframes fadeIn {
  to {
    opacity: 1;
    transform: translateY(0);
  }
  from {
    opacity: 0;
    transform: translateY(8px);
  }
}

.user-card:hover {
  background: rgba(255, 255, 255, 0.03);
}

.user-card.top-three {
  background: rgba(255, 255, 255, 0.02);
  border: 1px solid rgba(255, 255, 255, 0.04);
  margin-bottom: 8px;
}

.user-card.rank-1 {
  background: linear-gradient(135deg, rgba(251, 191, 36, 0.08), rgba(245, 158, 11, 0.04));
  border-color: rgba(251, 191, 36, 0.15);
}

.user-card.rank-2 {
  background: linear-gradient(135deg, rgba(148, 163, 184, 0.08), rgba(100, 116, 139, 0.04));
  border-color: rgba(148, 163, 184, 0.15);
}

.user-card.rank-3 {
  background: linear-gradient(135deg, rgba(217, 119, 6, 0.08), rgba(180, 83, 9, 0.04));
  border-color: rgba(217, 119, 6, 0.15);
}

/* Rank Badge */
.rank-badge {
  flex-shrink: 0;
  width: 36px;
  height: 36px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 10px;
  background: rgba(255, 255, 255, 0.05);
}

.rank-badge.is-medal {
  background: transparent;
}

.medal-emoji {
  font-size: 24px;
}

.rank-num {
  font-size: 14px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.6);
}

/* User Info */
.user-info {
  display: flex;
  align-items: center;
  gap: 12px;
  flex: 1;
  min-width: 0;
}

.user-avatar {
  width: 40px;
  height: 40px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
  flex-shrink: 0;
}

.user-name {
  font-size: 15px;
  font-weight: 500;
  color: #fff;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
  flex-shrink: 1;
  min-width: 0;
}

.user-title-badge {
  display: inline-flex;
  align-items: center;
  gap: 4px;
  padding: 3px 8px;
  border-radius: 12px;
  font-size: 11px;
  font-weight: 500;
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid;
  flex-shrink: 0;
  white-space: nowrap;
}

.user-title-badge .title-icon {
  font-size: 12px;
}

.user-title-badge .title-text {
  line-height: 1;
}

/* ===== Title Effects (称号特效) ===== */
/* 海神降临 - 神圣光芒 + 彩虹流光 */
.user-title-badge.effect-divine {
  position: relative;
  animation: divine-glow 2s ease-in-out infinite;
  box-shadow: 0 0 10px rgba(244, 114, 182, 0.4), 0 0 20px rgba(244, 114, 182, 0.2);
}

.user-title-badge.effect-divine::before {
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

.user-title-badge.effect-divine .title-icon {
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
.user-title-badge.effect-abyss {
  position: relative;
  animation: abyss-pulse 2.5s ease-in-out infinite;
  box-shadow: 0 0 8px rgba(167, 139, 250, 0.4), inset 0 0 6px rgba(139, 92, 246, 0.2);
}

.user-title-badge.effect-abyss .title-icon {
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
.user-title-badge.effect-legend {
  position: relative;
  overflow: hidden;
  animation: legend-shine 2s ease-in-out infinite;
  box-shadow: 0 0 8px rgba(251, 191, 36, 0.3);
}

.user-title-badge.effect-legend::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 50%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  animation: legend-sweep 2.5s ease-in-out infinite;
}

.user-title-badge.effect-legend .title-icon {
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
.user-title-badge.effect-master {
  position: relative;
  animation: master-glow 2s ease-in-out infinite;
}

.user-title-badge.effect-master::after {
  content: '';
  position: absolute;
  inset: -3px;
  border-radius: inherit;
  border: 1px solid rgba(56, 189, 248, 0.3);
  animation: master-ripple 2s ease-out infinite;
  pointer-events: none;
}

.user-title-badge.effect-master .title-icon {
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
.user-title-badge.effect-expert {
  animation: expert-float 2.5s ease-in-out infinite;
}

.user-title-badge.effect-expert .title-icon {
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
.user-title-badge.effect-advanced {
  animation: advanced-breathe 3s ease-in-out infinite;
}

.user-title-badge.effect-advanced .title-icon {
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
.user-title-badge.effect-beginner {
  animation: beginner-twinkle 3s ease-in-out infinite;
}

@keyframes beginner-twinkle {
  0%, 100% { opacity: 0.85; }
  50% { opacity: 1; }
}

/* User Stats */
.user-stats {
  display: flex;
  align-items: center;
  gap: 12px;
  flex-shrink: 0;
}

.stat-item {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 2px;
  min-width: 40px;
}

.stat-item.main-stat .stat-value {
  font-size: 18px;
}

.stat-value {
  font-size: 15px;
  font-weight: 700;
  color: #fff;
}

.stat-label {
  font-size: 11px;
  color: rgba(255, 255, 255, 0.55);
}

.stat-divider {
  width: 1px;
  height: 24px;
  background: rgba(255, 255, 255, 0.08);
}

/* Responsive */
@media (max-width: 640px) {
  .user-stats {
    gap: 8px;
  }

  .stat-item {
    min-width: 32px;
  }

  .stat-item.main-stat .stat-value {
    font-size: 16px;
  }

  .stat-value {
    font-size: 13px;
  }

  .stat-label {
    font-size: 10px;
  }

  .stat-divider {
    height: 20px;
  }

  .user-avatar {
    width: 36px;
    height: 36px;
    font-size: 14px;
  }

  .user-name {
    font-size: 14px;
  }

  .user-title-badge {
    padding: 2px 6px;
    font-size: 10px;
    gap: 3px;
  }

  .user-title-badge .title-icon {
    font-size: 10px;
  }
}

</style>
