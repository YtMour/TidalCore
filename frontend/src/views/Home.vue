<script setup lang="ts">
import { ref, onMounted } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Heatmap from '@/components/Heatmap.vue'
import { getGlobalHeatmap } from '@/api/checkin'

const heatmapData = ref<Record<string, number>>({})
const loading = ref(true)

onMounted(async () => {
  try {
    heatmapData.value = await getGlobalHeatmap(365)
  } catch {
    // 静默处理错误
  } finally {
    loading.value = false
  }
})
</script>

<template>
  <MainLayout>
    <div class="space-y-16">
      <section class="text-center py-16">
        <h1 class="text-5xl font-bold text-white mb-6">
          <span class="text-6xl">🌊</span>
          <br />
          TidalCore
        </h1>
        <p class="text-xl text-white/70 mb-8 max-w-2xl mx-auto">
          开源盆底肌训练平台，帮助你建立健康的训练习惯。
          <br />
          每天几分钟，坚持就是胜利。
        </p>
        <div class="flex justify-center gap-4">
          <RouterLink
            to="/train"
            class="px-8 py-3 rounded-full bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 text-white font-semibold text-lg transition-all shadow-lg shadow-purple-500/30"
          >
            开始训练
          </RouterLink>
          <RouterLink
            to="/leaderboard"
            class="px-8 py-3 rounded-full bg-white/10 hover:bg-white/20 text-white font-semibold text-lg transition-all"
          >
            查看排行
          </RouterLink>
        </div>
      </section>

      <section class="bg-white/5 rounded-2xl p-8">
        <h2 class="text-2xl font-bold text-white mb-6">全站打卡热力图</h2>
        <div v-if="loading" class="text-center text-white/50 py-8">
          加载中...
        </div>
        <Heatmap v-else :data="heatmapData" :days="365" />
      </section>

      <section class="grid md:grid-cols-3 gap-6">
        <div class="bg-white/5 rounded-xl p-6 text-center">
          <div class="text-4xl mb-4">⏱️</div>
          <h3 class="text-lg font-semibold text-white mb-2">科学计时</h3>
          <p class="text-white/60 text-sm">
            自定义收缩-保持-放松循环，专业训练节奏
          </p>
        </div>
        <div class="bg-white/5 rounded-xl p-6 text-center">
          <div class="text-4xl mb-4">🔥</div>
          <h3 class="text-lg font-semibold text-white mb-2">连续打卡</h3>
          <p class="text-white/60 text-sm">
            记录你的坚持轨迹，见证每一天的进步
          </p>
        </div>
        <div class="bg-white/5 rounded-xl p-6 text-center">
          <div class="text-4xl mb-4">🏆</div>
          <h3 class="text-lg font-semibold text-white mb-2">毅力排行</h3>
          <p class="text-white/60 text-sm">
            匿名排行榜，与全站用户一起坚持
          </p>
        </div>
      </section>
    </div>
  </MainLayout>
</template>
