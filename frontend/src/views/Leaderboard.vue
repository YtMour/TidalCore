<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import LeaderboardTable from '@/components/LeaderboardTable.vue'
import { Card } from '@/components/ui'
import { getLeaderboard, type LeaderboardUser } from '@/api/checkin'
import { Trophy, Timer, TrendingUp, Users, Crown, Medal, Award, Flame, Sparkles } from 'lucide-vue-next'

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
    <div class="space-y-12">
      <!-- Header -->
      <section
        class="text-center relative transition-all duration-700"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <!-- Decorative background -->
        <div class="absolute inset-0 -z-10 overflow-hidden">
          <div class="absolute top-0 left-1/2 -translate-x-1/2 w-96 h-48 bg-gradient-to-b from-amber-500/20 to-transparent rounded-full blur-3xl"></div>
        </div>

        <div class="inline-flex items-center justify-center w-24 h-24 rounded-3xl bg-gradient-to-br from-amber-400 to-orange-500 mb-6 shadow-xl shadow-orange-500/30 animate-float-slow">
          <Trophy class="w-12 h-12 text-white" />
        </div>

        <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-full bg-amber-500/10 border border-amber-500/20 mb-4">
          <Crown class="w-4 h-4 text-amber-400" />
          <span class="text-sm text-amber-300">荣耀榜单</span>
        </div>

        <h1 class="text-display mb-4">
          <span class="gradient-text-gold">毅力排行榜</span>
        </h1>
        <p class="text-lg text-white/50 max-w-xl mx-auto">
          坚持打卡，见证你的毅力。与全站用户一起，互相激励，共同进步。
        </p>
      </section>

      <!-- Top 3 Podium (when data loaded) -->
      <section
        v-if="!loading && top3Users.length >= 3"
        class="transition-all duration-700 delay-100"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <div class="flex items-end justify-center gap-4 md:gap-8 max-w-2xl mx-auto">
          <!-- 2nd Place -->
          <div class="flex-1 text-center">
            <div class="w-16 h-16 md:w-20 md:h-20 mx-auto mb-3 rounded-2xl bg-gradient-to-br from-slate-400 to-slate-500 flex items-center justify-center text-2xl md:text-3xl font-bold text-white shadow-lg">
              {{ top3Users[1]?.username?.[0]?.toUpperCase() || '?' }}
            </div>
            <Medal class="w-6 h-6 mx-auto mb-1 text-slate-400" />
            <div class="text-sm font-medium text-white truncate px-2">{{ top3Users[1]?.username }}</div>
            <div class="text-xs text-white/50">{{ top3Users[1]?.streak }} 天</div>
            <div class="mt-3 h-20 md:h-24 bg-gradient-to-t from-slate-500/20 to-transparent rounded-t-xl"></div>
          </div>

          <!-- 1st Place -->
          <div class="flex-1 text-center -mt-8">
            <div class="relative">
              <Crown class="w-8 h-8 mx-auto mb-2 text-amber-400 animate-float-fast" />
              <div class="w-20 h-20 md:w-24 md:h-24 mx-auto mb-3 rounded-2xl bg-gradient-to-br from-amber-400 to-orange-500 flex items-center justify-center text-3xl md:text-4xl font-bold text-white shadow-xl shadow-amber-500/30">
                {{ top3Users[0]?.username?.[0]?.toUpperCase() || '?' }}
              </div>
            </div>
            <Trophy class="w-6 h-6 mx-auto mb-1 text-amber-400" />
            <div class="text-sm font-medium text-white truncate px-2">{{ top3Users[0]?.username }}</div>
            <div class="text-xs text-amber-400 font-semibold">{{ top3Users[0]?.streak }} 天</div>
            <div class="mt-3 h-28 md:h-32 bg-gradient-to-t from-amber-500/20 to-transparent rounded-t-xl"></div>
          </div>

          <!-- 3rd Place -->
          <div class="flex-1 text-center">
            <div class="w-16 h-16 md:w-20 md:h-20 mx-auto mb-3 rounded-2xl bg-gradient-to-br from-amber-600 to-amber-700 flex items-center justify-center text-2xl md:text-3xl font-bold text-white shadow-lg">
              {{ top3Users[2]?.username?.[0]?.toUpperCase() || '?' }}
            </div>
            <Award class="w-6 h-6 mx-auto mb-1 text-amber-600" />
            <div class="text-sm font-medium text-white truncate px-2">{{ top3Users[2]?.username }}</div>
            <div class="text-xs text-white/50">{{ top3Users[2]?.streak }} 天</div>
            <div class="mt-3 h-16 md:h-20 bg-gradient-to-t from-amber-600/20 to-transparent rounded-t-xl"></div>
          </div>
        </div>
      </section>

      <!-- Stats Cards -->
      <section
        class="grid grid-cols-2 md:grid-cols-4 gap-4 transition-all duration-700 delay-200"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <div class="glass-card glass-card-hover p-5 text-center group">
          <div class="w-10 h-10 mx-auto mb-3 rounded-xl bg-gradient-to-br from-violet-500/20 to-purple-500/10 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Users class="w-5 h-5 text-violet-400" />
          </div>
          <div class="text-2xl font-bold text-white">{{ totalUsers }}</div>
          <div class="text-xs text-white/50">参与用户</div>
        </div>
        <div class="glass-card glass-card-hover p-5 text-center group">
          <div class="w-10 h-10 mx-auto mb-3 rounded-xl bg-gradient-to-br from-amber-500/20 to-orange-500/10 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Flame class="w-5 h-5 text-amber-400" />
          </div>
          <div class="text-2xl font-bold text-white">{{ topStreak }}</div>
          <div class="text-xs text-white/50">最高连续</div>
        </div>
        <div class="glass-card glass-card-hover p-5 text-center group">
          <div class="w-10 h-10 mx-auto mb-3 rounded-xl bg-gradient-to-br from-emerald-500/20 to-teal-500/10 flex items-center justify-center group-hover:scale-110 transition-transform">
            <TrendingUp class="w-5 h-5 text-emerald-400" />
          </div>
          <div class="text-2xl font-bold text-white">TOP 50</div>
          <div class="text-xs text-white/50">排行展示</div>
        </div>
        <div class="glass-card glass-card-hover p-5 text-center group">
          <div class="w-10 h-10 mx-auto mb-3 rounded-xl bg-gradient-to-br from-pink-500/20 to-rose-500/10 flex items-center justify-center group-hover:scale-110 transition-transform">
            <Sparkles class="w-5 h-5 text-pink-400" />
          </div>
          <div class="text-2xl font-bold text-white">实时</div>
          <div class="text-xs text-white/50">数据更新</div>
        </div>
      </section>

      <!-- Leaderboard Table -->
      <section
        class="transition-all duration-700 delay-300"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <LeaderboardTable :users="users" :loading="loading" />
      </section>

      <!-- CTA -->
      <section
        class="text-center transition-all duration-700 delay-400"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <Card padding="lg" class="max-w-2xl mx-auto relative overflow-hidden">
          <!-- Decorative elements -->
          <div class="absolute top-0 left-0 w-full h-full overflow-hidden pointer-events-none">
            <div class="absolute -top-20 -left-20 w-40 h-40 bg-gradient-to-br from-amber-500/15 to-transparent rounded-full blur-3xl"></div>
            <div class="absolute -bottom-20 -right-20 w-40 h-40 bg-gradient-to-br from-orange-500/15 to-transparent rounded-full blur-3xl"></div>
          </div>

          <div class="relative">
            <div class="w-14 h-14 mx-auto mb-5 rounded-2xl bg-gradient-to-br from-amber-500 to-orange-500 flex items-center justify-center shadow-lg shadow-amber-500/30">
              <Trophy class="w-7 h-7 text-white" />
            </div>
            <h3 class="text-xl font-bold text-white mb-3">想要冲击排行榜？</h3>
            <p class="text-white/50 mb-6">
              每天坚持训练，积累连续打卡天数，你也能登上榜单！
            </p>
            <RouterLink
              to="/train"
              class="group inline-flex items-center gap-2 btn-gradient px-10 py-4 rounded-full text-white font-semibold text-lg shadow-lg shadow-violet-500/25"
            >
              <Timer class="w-5 h-5" />
              开始训练，冲击排行
              <TrendingUp class="w-4 h-4 group-hover:translate-x-1 transition-transform" />
            </RouterLink>
          </div>
        </Card>
      </section>
    </div>
  </MainLayout>
</template>
