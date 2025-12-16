<script setup lang="ts">
import { ref, watch } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Timer from '@/components/Timer.vue'
import { useTrainingStore } from '@/store/training'
import { useUserStore } from '@/store/user'
import { checkin } from '@/api/checkin'
import confetti from 'canvas-confetti'

const trainingStore = useTrainingStore()
const userStore = useUserStore()

const showSettings = ref(false)
const checkinResult = ref<{
  streak: number
  maxStreak: number
  total: number
} | null>(null)
const checkinError = ref('')

watch(() => trainingStore.isRunning, async (isRunning, wasRunning) => {
  if (wasRunning && !isRunning && trainingStore.currentCycle > 0) {
    await handleCheckin()
  }
})

async function handleCheckin() {
  if (!userStore.isLoggedIn) {
    checkinError.value = '请先登录后再打卡'
    return
  }

  try {
    const res = await checkin({
      duration: trainingStore.totalDuration,
      cycles: trainingStore.currentCycle
    })

    checkinResult.value = {
      streak: res.current_streak,
      maxStreak: res.max_streak,
      total: res.total_checkin
    }

    userStore.updateStreak(res.current_streak, res.max_streak, res.total_checkin)

    confetti({
      particleCount: 150,
      spread: 100,
      origin: { y: 0.6 }
    })
  } catch (e) {
    checkinError.value = e instanceof Error ? e.message : '打卡失败'
  }
}

function closeResult() {
  checkinResult.value = null
  checkinError.value = ''
  trainingStore.reset()
}
</script>

<template>
  <MainLayout>
    <div class="max-w-2xl mx-auto">
      <div class="text-center mb-8">
        <h1 class="text-3xl font-bold text-white mb-2">潮汐训练器</h1>
        <p class="text-white/60">跟随节奏，完成今日训练</p>
      </div>

      <Timer />

      <div class="mt-8 flex justify-center">
        <button
          @click="showSettings = !showSettings"
          class="text-white/50 hover:text-white transition flex items-center gap-2"
        >
          <span>⚙️</span>
          <span>训练设置</span>
        </button>
      </div>

      <div v-if="showSettings" class="mt-6 bg-white/5 rounded-xl p-6">
        <h3 class="text-white font-semibold mb-4">自定义训练参数</h3>
        <div class="grid grid-cols-2 gap-4">
          <div>
            <label class="block text-white/60 text-sm mb-1">收缩时间 (秒)</label>
            <input
              type="number"
              :value="trainingStore.settings.contractTime"
              @input="trainingStore.updateSettings({ contractTime: +($event.target as HTMLInputElement).value })"
              min="1"
              max="30"
              class="w-full px-3 py-2 bg-white/10 rounded-lg text-white border border-white/10 focus:border-purple-500 focus:outline-none"
            />
          </div>
          <div>
            <label class="block text-white/60 text-sm mb-1">保持时间 (秒)</label>
            <input
              type="number"
              :value="trainingStore.settings.holdTime"
              @input="trainingStore.updateSettings({ holdTime: +($event.target as HTMLInputElement).value })"
              min="1"
              max="30"
              class="w-full px-3 py-2 bg-white/10 rounded-lg text-white border border-white/10 focus:border-purple-500 focus:outline-none"
            />
          </div>
          <div>
            <label class="block text-white/60 text-sm mb-1">放松时间 (秒)</label>
            <input
              type="number"
              :value="trainingStore.settings.relaxTime"
              @input="trainingStore.updateSettings({ relaxTime: +($event.target as HTMLInputElement).value })"
              min="1"
              max="30"
              class="w-full px-3 py-2 bg-white/10 rounded-lg text-white border border-white/10 focus:border-purple-500 focus:outline-none"
            />
          </div>
          <div>
            <label class="block text-white/60 text-sm mb-1">循环次数</label>
            <input
              type="number"
              :value="trainingStore.settings.cycles"
              @input="trainingStore.updateSettings({ cycles: +($event.target as HTMLInputElement).value })"
              min="1"
              max="50"
              class="w-full px-3 py-2 bg-white/10 rounded-lg text-white border border-white/10 focus:border-purple-500 focus:outline-none"
            />
          </div>
        </div>
      </div>

      <div
        v-if="checkinResult || checkinError"
        class="fixed inset-0 bg-black/50 flex items-center justify-center z-50"
        @click.self="closeResult"
      >
        <div class="bg-slate-800 rounded-2xl p-8 max-w-sm w-full mx-4 text-center">
          <template v-if="checkinResult">
            <div class="text-6xl mb-4">🎉</div>
            <h3 class="text-2xl font-bold text-white mb-4">打卡成功!</h3>
            <div class="space-y-2 text-white/70 mb-6">
              <p>连续打卡: <span class="text-orange-400 font-bold">{{ checkinResult.streak }}</span> 天</p>
              <p>最高记录: {{ checkinResult.maxStreak }} 天</p>
              <p>累计打卡: {{ checkinResult.total }} 次</p>
            </div>
          </template>
          <template v-else>
            <div class="text-6xl mb-4">😅</div>
            <h3 class="text-xl font-bold text-white mb-4">{{ checkinError }}</h3>
          </template>
          <button
            @click="closeResult"
            class="px-6 py-2 rounded-full bg-purple-600 hover:bg-purple-500 text-white transition"
          >
            确定
          </button>
        </div>
      </div>
    </div>
  </MainLayout>
</template>
