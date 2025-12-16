<script setup lang="ts">
import type { LeaderboardUser } from '@/api/checkin'

defineProps<{
  users: LeaderboardUser[]
  loading?: boolean
}>()

function getMedal(index: number): string {
  const medals = ['🥇', '🥈', '🥉']
  return medals[index] || ''
}

function maskUsername(username: string): string {
  if (!username || username.length === 0) return '***'
  if (username.length === 1) return username[0] + '**'
  if (username.length === 2) return username[0] + '*' + username[1]
  return username[0] + '*'.repeat(Math.min(username.length - 2, 5)) + username[username.length - 1]
}
</script>

<template>
  <div class="bg-white/5 rounded-xl overflow-hidden">
    <div v-if="loading" class="p-8 text-center text-white/50">
      加载中...
    </div>

    <table v-else class="w-full">
      <thead>
        <tr class="border-b border-white/10 text-white/50 text-sm">
          <th class="px-4 py-3 text-left font-medium">排名</th>
          <th class="px-4 py-3 text-left font-medium">用户</th>
          <th class="px-4 py-3 text-right font-medium">连续天数</th>
          <th class="px-4 py-3 text-right font-medium">最高记录</th>
          <th class="px-4 py-3 text-right font-medium">总打卡</th>
        </tr>
      </thead>
      <tbody>
        <tr
          v-for="(user, index) in users"
          :key="user.id"
          class="border-b border-white/5 hover:bg-white/5 transition"
        >
          <td class="px-4 py-3">
            <span class="text-lg">{{ getMedal(index) }}</span>
            <span v-if="index >= 3" class="text-white/50">{{ index + 1 }}</span>
          </td>
          <td class="px-4 py-3 text-white font-medium">
            {{ maskUsername(user.username) }}
          </td>
          <td class="px-4 py-3 text-right">
            <span class="text-orange-400 font-bold">{{ user.streak }}</span>
            <span class="text-white/50 text-sm ml-1">天</span>
          </td>
          <td class="px-4 py-3 text-right text-white/70">
            {{ user.max_streak }} 天
          </td>
          <td class="px-4 py-3 text-right text-white/70">
            {{ user.total_checkin }} 次
          </td>
        </tr>

        <tr v-if="users.length === 0">
          <td colspan="5" class="px-4 py-8 text-center text-white/50">
            暂无数据
          </td>
        </tr>
      </tbody>
    </table>
  </div>
</template>
