<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Heatmap from '@/components/Heatmap.vue'
import { useUserStore } from '@/store/user'
import { getHeatmap, getHistory, type CheckinRecord } from '@/api/checkin'

const userStore = useUserStore()

const heatmapData = ref<Record<string, number>>({})
const recentHistory = ref<CheckinRecord[]>([])
const loading = ref(true)

onMounted(async () => {
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
      const date = c.checked_at.split('T')[0]
      heatmap[date] = (heatmap[date] || 0) + 1
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
</script>

<template>
  <MainLayout>
    <div class="space-y-8">
      <section class="bg-white/5 rounded-2xl p-8">
        <div class="flex items-center gap-6">
          <div class="w-20 h-20 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center text-3xl">
            {{ userStore.user?.username?.[0]?.toUpperCase() || '?' }}
          </div>
          <div>
            <h1 class="text-2xl font-bold text-white">{{ userStore.user?.username }}</h1>
            <p class="text-white/50">坚持就是胜利</p>
          </div>
        </div>

        <div class="grid grid-cols-3 gap-4 mt-8">
          <div class="bg-white/5 rounded-xl p-4 text-center">
            <div class="text-3xl font-bold text-orange-400">{{ userStore.user?.streak || 0 }}</div>
            <div class="text-white/50 text-sm mt-1">连续打卡</div>
          </div>
          <div class="bg-white/5 rounded-xl p-4 text-center">
            <div class="text-3xl font-bold text-purple-400">{{ userStore.user?.max_streak || 0 }}</div>
            <div class="text-white/50 text-sm mt-1">最高记录</div>
          </div>
          <div class="bg-white/5 rounded-xl p-4 text-center">
            <div class="text-3xl font-bold text-emerald-400">{{ userStore.user?.total_checkin || 0 }}</div>
            <div class="text-white/50 text-sm mt-1">累计打卡</div>
          </div>
        </div>
      </section>

      <section class="bg-white/5 rounded-2xl p-8">
        <h2 class="text-xl font-bold text-white mb-6">我的打卡热力图</h2>
        <div v-if="loading" class="text-center text-white/50 py-8">
          加载中...
        </div>
        <Heatmap v-else :data="heatmapData" :days="365" />
      </section>

      <section class="bg-white/5 rounded-2xl p-8">
        <h2 class="text-xl font-bold text-white mb-6">最近打卡记录</h2>
        <div v-if="loading" class="text-center text-white/50 py-8">
          加载中...
        </div>
        <div v-else-if="recentHistory.length === 0" class="text-center text-white/50 py-8">
          暂无打卡记录，快去训练吧！
        </div>
        <div v-else class="space-y-3">
          <div
            v-for="record in recentHistory"
            :key="record.id"
            class="flex items-center justify-between p-4 bg-white/5 rounded-lg"
          >
            <div>
              <div class="text-white">{{ formatDate(record.checked_at) }}</div>
              <div class="text-white/50 text-sm">
                {{ record.cycles }} 个循环 · {{ formatDuration(record.duration) }}
              </div>
            </div>
            <div class="text-2xl">✅</div>
          </div>
        </div>
      </section>

      <div class="text-center">
        <RouterLink
          to="/train"
          class="inline-block px-8 py-3 rounded-full bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 text-white font-semibold transition-all shadow-lg shadow-purple-500/30"
        >
          开始今日训练
        </RouterLink>
      </div>
    </div>
  </MainLayout>
</template>
