<script setup lang="ts">
import { ref, watch, onMounted, computed } from 'vue'
import MainLayout from '@/layouts/MainLayout.vue'
import Timer from '@/components/Timer.vue'
import { Card, Modal, Input } from '@/components/ui'
import { useTrainingStore } from '@/store/training'
import { useUserStore } from '@/store/user'
import { checkin } from '@/api/checkin'
import { Settings, ChevronDown, ChevronUp, Sparkles, AlertCircle, Info, Zap, Target, Clock, TrendingUp } from 'lucide-vue-next'
import confetti from 'canvas-confetti'

const trainingStore = useTrainingStore()
const userStore = useUserStore()

const mounted = ref(false)
const showSettings = ref(false)
const showResultModal = ref(false)

onMounted(() => {
  setTimeout(() => {
    mounted.value = true
  }, 100)
})

// 计算预计训练时长
const estimatedDuration = computed(() => {
  const { contractTime, holdTime, relaxTime, cycles } = trainingStore.settings
  const totalSeconds = (contractTime + holdTime + relaxTime) * cycles
  const mins = Math.floor(totalSeconds / 60)
  const secs = totalSeconds % 60
  return mins > 0 ? `${mins}分${secs}秒` : `${secs}秒`
})
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
    showResultModal.value = true
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
    showResultModal.value = true

    // Enhanced confetti
    const colors = ['#a855f7', '#ec4899', '#f59e0b', '#10b981']
    confetti({
      particleCount: 150,
      spread: 100,
      origin: { y: 0.6 },
      colors
    })
  } catch (e) {
    checkinError.value = e instanceof Error ? e.message : '打卡失败'
    showResultModal.value = true
  }
}

function closeResult() {
  showResultModal.value = false
  checkinResult.value = null
  checkinError.value = ''
  trainingStore.reset()
}

function updateSetting(key: string, value: number) {
  trainingStore.updateSettings({ [key]: value })
}
</script>

<template>
  <MainLayout>
    <div class="max-w-3xl mx-auto">
      <!-- Header with decorative elements -->
      <div
        class="text-center mb-12 relative transition-all duration-700"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <!-- Decorative background -->
        <div class="absolute inset-0 -z-10 overflow-hidden">
          <div class="absolute top-0 left-1/2 -translate-x-1/2 w-96 h-32 bg-gradient-to-b from-violet-500/20 to-transparent rounded-full blur-3xl"></div>
        </div>

        <div class="inline-flex items-center gap-2 px-3 py-1.5 rounded-lg bg-violet-500/10 border border-violet-500/20 mb-4">
          <Zap class="w-4 h-4 text-violet-400" />
          <span class="text-sm text-violet-300">专业训练</span>
        </div>

        <h1 class="text-heading text-white mb-3">
          <span class="gradient-text text-glow-purple">潮汐训练器</span>
        </h1>
        <p class="text-white/50">跟随节奏，完成今日训练，坚持就是胜利</p>

        <!-- Quick stats -->
        <div class="flex items-center justify-center gap-6 mt-6">
          <div class="flex items-center gap-2 text-sm text-white/40">
            <Clock class="w-4 h-4" />
            <span>预计 {{ estimatedDuration }}</span>
          </div>
          <div class="flex items-center gap-2 text-sm text-white/40">
            <Target class="w-4 h-4" />
            <span>{{ trainingStore.settings.cycles }} 个循环</span>
          </div>
        </div>
      </div>

      <!-- Timer Component with enhanced container -->
      <div
        class="relative transition-all duration-700 delay-100"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <!-- Ambient glow behind timer -->
        <div
          class="absolute inset-0 -z-10 transition-opacity duration-500"
          :class="trainingStore.isRunning ? 'opacity-100' : 'opacity-0'"
        >
          <div class="absolute top-1/2 left-1/2 -translate-x-1/2 -translate-y-1/2 w-80 h-80 bg-gradient-to-br from-violet-500/20 to-pink-500/20 rounded-full blur-3xl animate-pulse-soft"></div>
        </div>

        <Timer />
      </div>

      <!-- Settings Toggle -->
      <div
        class="mt-12 flex justify-center transition-all duration-700 delay-200"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <button
          @click="showSettings = !showSettings"
          class="group flex items-center gap-2 px-5 py-2.5 rounded-lg text-white/50 hover:text-white bg-white/[0.02] hover:bg-white/[0.04] border border-white/[0.06] hover:border-white/[0.1] transition-all duration-200"
        >
          <Settings class="w-5 h-5 group-hover:rotate-45 transition-transform duration-300" />
          <span>训练设置</span>
          <ChevronUp v-if="showSettings" class="w-4 h-4 transition-transform" />
          <ChevronDown v-else class="w-4 h-4 transition-transform" />
        </button>
      </div>

      <!-- Settings Panel -->
      <Transition name="slide-up">
        <Card v-if="showSettings" class="mt-6 relative overflow-hidden">
          <!-- Decorative corner -->
          <div class="absolute top-0 right-0 w-32 h-32 bg-gradient-to-bl from-violet-500/10 to-transparent rounded-full blur-2xl pointer-events-none"></div>

          <div class="relative">
            <h3 class="text-lg font-semibold text-white mb-6 flex items-center gap-2">
              <div class="w-8 h-8 rounded-lg bg-gradient-to-br from-violet-500/20 to-purple-500/10 flex items-center justify-center">
                <Settings class="w-4 h-4 text-violet-400" />
              </div>
              自定义训练参数
            </h3>

            <div class="grid grid-cols-2 gap-6">
              <Input
                type="number"
                label="收缩时间 (秒)"
                :model-value="trainingStore.settings.contractTime"
                @update:model-value="updateSetting('contractTime', $event as number)"
                :min="1"
                :max="30"
                placeholder="1-30"
              />
              <Input
                type="number"
                label="保持时间 (秒)"
                :model-value="trainingStore.settings.holdTime"
                @update:model-value="updateSetting('holdTime', $event as number)"
                :min="1"
                :max="30"
                placeholder="1-30"
              />
              <Input
                type="number"
                label="放松时间 (秒)"
                :model-value="trainingStore.settings.relaxTime"
                @update:model-value="updateSetting('relaxTime', $event as number)"
                :min="1"
                :max="30"
                placeholder="1-30"
              />
              <Input
                type="number"
                label="循环次数"
                :model-value="trainingStore.settings.cycles"
                @update:model-value="updateSetting('cycles', $event as number)"
                :min="1"
                :max="50"
                placeholder="1-50"
              />
            </div>

            <!-- Training summary -->
            <div class="mt-6 p-4 rounded-lg bg-gradient-to-r from-violet-500/10 to-purple-500/10 border border-violet-500/20">
              <div class="flex items-start gap-3">
                <Info class="w-5 h-5 text-violet-400 shrink-0 mt-0.5" />
                <div>
                  <p class="text-sm text-white/70 mb-2">
                    <span class="text-violet-400 font-medium">训练概览：</span>
                    每个循环包含收缩、保持、放松三个阶段
                  </p>
                  <div class="flex flex-wrap gap-4 text-xs text-white/50">
                    <span class="flex items-center gap-1">
                      <Clock class="w-3.5 h-3.5" />
                      总时长: {{ estimatedDuration }}
                    </span>
                    <span class="flex items-center gap-1">
                      <Target class="w-3.5 h-3.5" />
                      循环数: {{ trainingStore.settings.cycles }}
                    </span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </Card>
      </Transition>

      <!-- Tips Section -->
      <div
        class="mt-12 transition-all duration-700 delay-300"
        :class="mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'"
      >
        <div class="grid sm:grid-cols-3 gap-4">
          <div class="p-4 rounded-xl bg-white/[0.02] border border-white/5 hover:border-white/10 transition-all">
            <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-rose-500/20 to-pink-500/10 flex items-center justify-center mb-3">
              <Zap class="w-5 h-5 text-rose-400" />
            </div>
            <h4 class="text-sm font-medium text-white mb-1">收缩阶段</h4>
            <p class="text-xs text-white/40">用力收紧盆底肌肉</p>
          </div>
          <div class="p-4 rounded-xl bg-white/[0.02] border border-white/5 hover:border-white/10 transition-all">
            <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-amber-500/20 to-orange-500/10 flex items-center justify-center mb-3">
              <Target class="w-5 h-5 text-amber-400" />
            </div>
            <h4 class="text-sm font-medium text-white mb-1">保持阶段</h4>
            <p class="text-xs text-white/40">维持收缩状态</p>
          </div>
          <div class="p-4 rounded-xl bg-white/[0.02] border border-white/5 hover:border-white/10 transition-all">
            <div class="w-10 h-10 rounded-xl bg-gradient-to-br from-emerald-500/20 to-teal-500/10 flex items-center justify-center mb-3">
              <TrendingUp class="w-5 h-5 text-emerald-400" />
            </div>
            <h4 class="text-sm font-medium text-white mb-1">放松阶段</h4>
            <p class="text-xs text-white/40">完全放松肌肉</p>
          </div>
        </div>
      </div>

      <!-- Result Modal -->
      <Modal v-model="showResultModal" size="sm">
        <div class="text-center py-4">
          <template v-if="checkinResult">
            <div class="w-20 h-20 mx-auto mb-6 rounded-full bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center">
              <Sparkles class="w-10 h-10 text-white" />
            </div>
            <h3 class="text-2xl font-bold text-white mb-2">打卡成功!</h3>
            <p class="text-white/50 mb-6">太棒了，继续保持！</p>

            <div class="grid grid-cols-3 gap-4 mb-8">
              <div class="glass-card p-4">
                <div class="text-2xl font-bold text-orange-400">{{ checkinResult.streak }}</div>
                <div class="text-xs text-white/50 mt-1">连续天数</div>
              </div>
              <div class="glass-card p-4">
                <div class="text-2xl font-bold text-purple-400">{{ checkinResult.maxStreak }}</div>
                <div class="text-xs text-white/50 mt-1">最高记录</div>
              </div>
              <div class="glass-card p-4">
                <div class="text-2xl font-bold text-emerald-400">{{ checkinResult.total }}</div>
                <div class="text-xs text-white/50 mt-1">累计打卡</div>
              </div>
            </div>
          </template>

          <template v-else>
            <div class="w-20 h-20 mx-auto mb-6 rounded-full bg-red-500/20 flex items-center justify-center">
              <AlertCircle class="w-10 h-10 text-red-400" />
            </div>
            <h3 class="text-xl font-bold text-white mb-2">{{ checkinError }}</h3>
            <p class="text-white/50 mb-6">
              <RouterLink to="/login" class="text-purple-400 hover:text-purple-300">
                点击登录
              </RouterLink>
              后即可记录打卡
            </p>
          </template>

          <button
            @click="closeResult"
            class="btn-gradient px-8 py-3 rounded-full text-white font-semibold"
          >
            确定
          </button>
        </div>
      </Modal>
    </div>
  </MainLayout>
</template>
