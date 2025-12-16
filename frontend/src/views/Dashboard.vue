<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Heatmap from '@/components/Heatmap.vue'
import { Card, Skeleton } from '@/components/ui'
import { useUserStore } from '@/store/user'
import { getHeatmap, getHistory, type CheckinRecord } from '@/api/checkin'
import { Timer, Calendar, Clock, CheckCircle, ArrowRight, Flame, Trophy, Target, Sparkles } from 'lucide-vue-next'

const userStore = useUserStore()

const heatmapData = ref<Record<string, number>>({})
const recentHistory = ref<CheckinRecord[]>([])
const loading = ref(true)
const mounted = ref(false)

// 计算用户等级
const userLevel = computed(() => {
  const total = userStore.user?.total_checkin || 0
  if (total >= 365) return { name: '传奇', color: 'text-amber-400', bg: 'from-amber-500/20 to-orange-500/10' }
  if (total >= 180) return { name: '大师', color: 'text-violet-400', bg: 'from-violet-500/20 to-purple-500/10' }
  if (total >= 90) return { name: '专家', color: 'text-cyan-400', bg: 'from-cyan-500/20 to-blue-500/10' }
  if (total >= 30) return { name: '进阶', color: 'text-emerald-400', bg: 'from-emerald-500/20 to-teal-500/10' }
  if (total >= 7) return { name: '新手', color: 'text-pink-400', bg: 'from-pink-500/20 to-rose-500/10' }
  return { name: '初学者', color: 'text-white/60', bg: 'from-white/10 to-white/5' }
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
    <div class="space-y-8">
      <!-- Profile Header -->
      <section
        class="transition-all duration-700"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <Card class="relative overflow-hidden">
          <!-- Decorative background -->
          <div class="absolute top-0 right-0 w-64 h-64 bg-gradient-to-bl from-violet-500/10 to-transparent rounded-full blur-3xl pointer-events-none"></div>
          <div class="absolute bottom-0 left-0 w-48 h-48 bg-gradient-to-tr from-pink-500/10 to-transparent rounded-full blur-3xl pointer-events-none"></div>

          <div class="relative flex flex-col sm:flex-row items-center sm:items-start gap-6">
            <!-- Avatar with level badge -->
            <div class="relative">
              <div class="w-28 h-28 rounded-3xl bg-gradient-to-br from-violet-500 to-pink-500 flex items-center justify-center text-5xl font-bold text-white shadow-xl shadow-violet-500/30">
                {{ userStore.user?.username?.[0]?.toUpperCase() || '?' }}
              </div>
              <!-- Level badge -->
              <div
                class="absolute -bottom-2 -right-2 px-3 py-1 rounded-full text-xs font-semibold border border-white/10 bg-gradient-to-br"
                :class="[userLevel.color, userLevel.bg]"
              >
                {{ userLevel.name }}
              </div>
            </div>

            <!-- User Info -->
            <div class="text-center sm:text-left flex-1">
              <div class="flex items-center justify-center sm:justify-start gap-2 mb-1">
                <h1 class="text-2xl font-bold text-white">
                  {{ userStore.user?.username }}
                </h1>
                <Sparkles class="w-5 h-5 text-amber-400" />
              </div>
              <p class="text-white/50 mb-2">坚持就是胜利，每天进步一点点</p>

              <!-- Quick stats inline -->
              <div class="flex items-center justify-center sm:justify-start gap-4 mb-4 text-sm">
                <span class="flex items-center gap-1 text-orange-400">
                  <Flame class="w-4 h-4" />
                  {{ userStore.user?.streak || 0 }} 天连续
                </span>
                <span class="flex items-center gap-1 text-violet-400">
                  <Trophy class="w-4 h-4" />
                  {{ userStore.user?.max_streak || 0 }} 最高
                </span>
              </div>

              <RouterLink
                to="/train"
                class="group inline-flex items-center gap-2 btn-gradient btn-gradient-glow px-7 py-3 rounded-full text-white font-medium shadow-lg shadow-violet-500/25"
              >
                <Timer class="w-5 h-5" />
                开始今日训练
                <ArrowRight class="w-4 h-4 group-hover:translate-x-1 transition-transform" />
              </RouterLink>
            </div>
          </div>
        </Card>
      </section>

      <!-- Stats Cards -->
      <section
        class="grid grid-cols-1 sm:grid-cols-3 gap-4 transition-all duration-700 delay-100"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <div class="glass-card glass-card-hover p-6 text-center group">
          <div class="w-14 h-14 mx-auto mb-4 rounded-2xl bg-gradient-to-br from-orange-500/20 to-amber-500/10 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Flame class="w-7 h-7 text-orange-400" />
          </div>
          <div class="text-3xl font-bold text-white mb-1">{{ userStore.user?.streak || 0 }}</div>
          <div class="text-sm text-white/50">连续打卡</div>
        </div>

        <div class="glass-card glass-card-hover p-6 text-center group">
          <div class="w-14 h-14 mx-auto mb-4 rounded-2xl bg-gradient-to-br from-violet-500/20 to-purple-500/10 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Trophy class="w-7 h-7 text-violet-400" />
          </div>
          <div class="text-3xl font-bold text-white mb-1">{{ userStore.user?.max_streak || 0 }}</div>
          <div class="text-sm text-white/50">最高记录</div>
        </div>

        <div class="glass-card glass-card-hover p-6 text-center group">
          <div class="w-14 h-14 mx-auto mb-4 rounded-2xl bg-gradient-to-br from-emerald-500/20 to-teal-500/10 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Target class="w-7 h-7 text-emerald-400" />
          </div>
          <div class="text-3xl font-bold text-white mb-1">{{ userStore.user?.total_checkin || 0 }}</div>
          <div class="text-sm text-white/50">累计打卡</div>
        </div>
      </section>

      <!-- Heatmap Section -->
      <section
        class="transition-all duration-500 delay-200"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <Card>
          <div class="flex items-center gap-3 mb-6">
            <Calendar class="w-5 h-5 text-purple-400" />
            <h2 class="text-lg font-semibold text-white">我的打卡热力图</h2>
          </div>

          <div v-if="loading" class="py-8 space-y-4">
            <Skeleton height="1rem" width="40%" />
            <div class="flex gap-1">
              <Skeleton v-for="i in 15" :key="i" width="3rem" height="5rem" />
            </div>
          </div>
          <Heatmap v-else :data="heatmapData" :days="365" />
        </Card>
      </section>

      <!-- Recent History -->
      <section
        class="transition-all duration-500 delay-300"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <Card>
          <div class="flex items-center gap-3 mb-6">
            <Clock class="w-5 h-5 text-purple-400" />
            <h2 class="text-lg font-semibold text-white">最近打卡记录</h2>
          </div>

          <div v-if="loading" class="space-y-3">
            <div v-for="i in 5" :key="i" class="flex items-center gap-4 p-4 rounded-lg bg-white/5">
              <Skeleton width="3rem" height="3rem" rounded="lg" />
              <div class="flex-1 space-y-2">
                <Skeleton height="1rem" width="60%" />
                <Skeleton height="0.75rem" width="40%" />
              </div>
            </div>
          </div>

          <div v-else-if="recentHistory.length === 0" class="text-center py-12">
            <div class="w-16 h-16 mx-auto mb-4 rounded-full bg-white/5 flex items-center justify-center">
              <Timer class="w-8 h-8 text-white/30" />
            </div>
            <p class="text-white/50 mb-2">暂无打卡记录</p>
            <p class="text-white/30 text-sm">完成训练后，记录将显示在这里</p>
          </div>

          <div v-else class="space-y-3">
            <TransitionGroup name="list">
              <div
                v-for="(record, index) in recentHistory"
                :key="record.id"
                class="flex items-center gap-4 p-4 rounded-xl bg-white/5 hover:bg-white/8 transition-all duration-300"
                :style="{ transitionDelay: `${index * 50}ms` }"
              >
                <div class="w-12 h-12 rounded-xl bg-gradient-to-br from-emerald-500/20 to-teal-500/20 flex items-center justify-center shrink-0">
                  <CheckCircle class="w-6 h-6 text-emerald-400" />
                </div>
                <div class="flex-1 min-w-0">
                  <div class="text-white font-medium">{{ getRelativeTime(record.checked_at) }}</div>
                  <div class="text-white/50 text-sm truncate">
                    {{ record.cycles }} 个循环 · {{ formatDuration(record.duration) }}
                  </div>
                </div>
                <div class="text-xs text-white/30 shrink-0">
                  {{ formatDate(record.checked_at) }}
                </div>
              </div>
            </TransitionGroup>
          </div>
        </Card>
      </section>
    </div>
  </MainLayout>
</template>
