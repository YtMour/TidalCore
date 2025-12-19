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
}

</style>
