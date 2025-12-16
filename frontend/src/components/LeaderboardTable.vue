<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { LeaderboardUser } from '@/api/checkin'
import { Trophy, Flame, Target, Award } from 'lucide-vue-next'
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

function getMedal(index: number): { icon: string; color: string } | null {
  const medals = [
    { icon: '🥇', color: 'from-yellow-400 to-amber-500' },
    { icon: '🥈', color: 'from-gray-300 to-gray-400' },
    { icon: '🥉', color: 'from-amber-600 to-orange-700' }
  ]
  return medals[index] || null
}

function maskUsername(username: string): string {
  if (!username || username.length === 0) return '***'
  if (username.length === 1) return username[0] + '**'
  if (username.length === 2) return username[0] + '*' + username[1]
  return username[0] + '*'.repeat(Math.min(username.length - 2, 5)) + username[username.length - 1]
}

function getStreakColor(streak: number): string {
  if (streak >= 30) return 'text-orange-400'
  if (streak >= 14) return 'text-amber-400'
  if (streak >= 7) return 'text-yellow-400'
  return 'text-white/70'
}
</script>

<template>
  <div class="glass-card overflow-hidden">
    <!-- Loading State -->
    <div v-if="loading" class="p-6 space-y-4">
      <div v-for="i in 5" :key="i" class="flex items-center gap-4">
        <Skeleton width="2rem" height="2rem" rounded="full" />
        <Skeleton width="40%" height="1.25rem" />
        <div class="ml-auto flex gap-8">
          <Skeleton width="3rem" height="1.25rem" />
          <Skeleton width="3rem" height="1.25rem" />
          <Skeleton width="3rem" height="1.25rem" />
        </div>
      </div>
    </div>

    <!-- Table -->
    <div v-else class="overflow-x-auto">
      <table class="w-full">
        <thead>
          <tr class="border-b border-white/5">
            <th class="px-6 py-4 text-left text-sm font-medium text-white/50">排名</th>
            <th class="px-6 py-4 text-left text-sm font-medium text-white/50">用户</th>
            <th class="px-6 py-4 text-right text-sm font-medium text-white/50">
              <span class="inline-flex items-center gap-1.5">
                <Flame class="w-4 h-4 text-orange-400" />
                连续天数
              </span>
            </th>
            <th class="px-6 py-4 text-right text-sm font-medium text-white/50">
              <span class="inline-flex items-center gap-1.5">
                <Trophy class="w-4 h-4 text-purple-400" />
                最高记录
              </span>
            </th>
            <th class="px-6 py-4 text-right text-sm font-medium text-white/50">
              <span class="inline-flex items-center gap-1.5">
                <Target class="w-4 h-4 text-emerald-400" />
                总打卡
              </span>
            </th>
          </tr>
        </thead>
        <tbody>
          <TransitionGroup name="list">
            <tr
              v-for="(user, index) in users"
              :key="user.id"
              class="table-row border-b border-white/5 last:border-0"
              :style="{ transitionDelay: mounted ? `${index * 50}ms` : '0ms' }"
            >
              <!-- Rank -->
              <td class="px-6 py-4">
                <div class="flex items-center gap-2">
                  <template v-if="getMedal(index)">
                    <span class="text-2xl">{{ getMedal(index)!.icon }}</span>
                  </template>
                  <template v-else>
                    <span class="w-8 h-8 rounded-full bg-white/5 flex items-center justify-center text-white/50 text-sm font-medium">
                      {{ index + 1 }}
                    </span>
                  </template>
                </div>
              </td>

              <!-- Username -->
              <td class="px-6 py-4">
                <div class="flex items-center gap-3">
                  <div
                    class="w-10 h-10 rounded-full flex items-center justify-center text-white font-medium"
                    :class="index < 3 ? `bg-gradient-to-br ${getMedal(index)!.color}` : 'bg-white/10'"
                  >
                    {{ user.username?.[0]?.toUpperCase() || '?' }}
                  </div>
                  <span class="text-white font-medium">{{ maskUsername(user.username) }}</span>
                </div>
              </td>

              <!-- Streak -->
              <td class="px-6 py-4 text-right">
                <span class="font-bold text-lg" :class="getStreakColor(user.streak)">
                  {{ user.streak }}
                </span>
                <span class="text-white/40 text-sm ml-1">天</span>
              </td>

              <!-- Max Streak -->
              <td class="px-6 py-4 text-right text-white/60">
                {{ user.max_streak }} 天
              </td>

              <!-- Total -->
              <td class="px-6 py-4 text-right text-white/60">
                {{ user.total_checkin }} 次
              </td>
            </tr>
          </TransitionGroup>

          <!-- Empty State -->
          <tr v-if="users.length === 0">
            <td colspan="5" class="px-6 py-12 text-center">
              <Award class="w-12 h-12 mx-auto text-white/20 mb-3" />
              <p class="text-white/50">暂无排行数据</p>
              <p class="text-white/30 text-sm mt-1">成为第一个打卡的人吧！</p>
            </td>
          </tr>
        </tbody>
      </table>
    </div>
  </div>
</template>
