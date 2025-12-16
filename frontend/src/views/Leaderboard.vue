<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import LeaderboardTable from '@/components/LeaderboardTable.vue'
import { getLeaderboard, type LeaderboardUser } from '@/api/checkin'

const users = ref<LeaderboardUser[]>([])
const loading = ref(true)

onMounted(async () => {
  try {
    users.value = await getLeaderboard(50)
  } catch {
    // 静默处理
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <MainLayout>
    <div class="space-y-8">
      <div class="text-center">
        <h1 class="text-3xl font-bold text-white mb-2">🏆 毅力排行榜</h1>
        <p class="text-white/60">坚持打卡，见证你的毅力</p>
      </div>

      <LeaderboardTable :users="users" :loading="loading" />

      <div class="text-center">
        <RouterLink
          to="/train"
          class="inline-block px-8 py-3 rounded-full bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 text-white font-semibold transition-all shadow-lg shadow-purple-500/30"
        >
          开始训练，冲击排行
        </RouterLink>
      </div>
    </div>
  </MainLayout>
</template>
