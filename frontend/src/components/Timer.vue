<script setup lang="ts">
import { computed, onUnmounted } from 'vue'
import { useTrainingStore } from '@/store/training'
import confetti from 'canvas-confetti'

const trainingStore = useTrainingStore()

const phaseColors: Record<string, string> = {
  idle: 'from-slate-600 to-slate-700',
  contract: 'from-rose-500 to-pink-600',
  hold: 'from-amber-500 to-orange-600',
  relax: 'from-emerald-500 to-teal-600'
}

const currentColor = computed(() => phaseColors[trainingStore.phase])

const circumference = 2 * Math.PI * 120
const strokeDashoffset = computed(() => {
  if (trainingStore.phase === 'idle') return circumference
  const progress = trainingStore.phaseProgress / 100
  return circumference * (1 - progress)
})

function handleStart() {
  trainingStore.start()
}

function handleStop() {
  const cycles = trainingStore.currentCycle
  trainingStore.stop()
  if (cycles > 0) {
    celebrateCompletion()
  }
}

function celebrateCompletion() {
  confetti({
    particleCount: 100,
    spread: 70,
    origin: { y: 0.6 }
  })
}

function formatTime(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

onUnmounted(() => {
  trainingStore.clearTimers()
})
</script>

<template>
  <div class="flex flex-col items-center gap-8">
    <div class="relative w-72 h-72">
      <svg class="w-full h-full -rotate-90" viewBox="0 0 256 256">
        <circle
          cx="128"
          cy="128"
          r="120"
          fill="none"
          stroke="currentColor"
          stroke-width="8"
          class="text-white/10"
        />
        <circle
          cx="128"
          cy="128"
          r="120"
          fill="none"
          stroke="url(#gradient)"
          stroke-width="8"
          stroke-linecap="round"
          :stroke-dasharray="circumference"
          :stroke-dashoffset="strokeDashoffset"
          class="transition-all duration-200 ease-linear"
        />
        <defs>
          <linearGradient id="gradient" x1="0%" y1="0%" x2="100%" y2="100%">
            <stop offset="0%" stop-color="#a855f7" />
            <stop offset="100%" stop-color="#ec4899" />
          </linearGradient>
        </defs>
      </svg>

      <div
        class="absolute inset-0 flex flex-col items-center justify-center rounded-full bg-gradient-to-br transition-all duration-500"
        :class="currentColor"
        style="margin: 16px;"
      >
        <span class="text-6xl font-bold text-white">
          {{ trainingStore.countdown }}
        </span>
        <span class="text-xl text-white/80 mt-2">
          {{ trainingStore.phaseLabel }}
        </span>
      </div>
    </div>

    <div class="flex items-center gap-8 text-white/60">
      <div class="text-center">
        <div class="text-3xl font-bold text-white">{{ trainingStore.currentCycle }}</div>
        <div class="text-sm">/ {{ trainingStore.settings.cycles }} 循环</div>
      </div>
      <div class="w-px h-12 bg-white/20"></div>
      <div class="text-center">
        <div class="text-3xl font-bold text-white">{{ formatTime(trainingStore.totalDuration) }}</div>
        <div class="text-sm">训练时长</div>
      </div>
    </div>

    <div class="flex gap-4">
      <button
        v-if="!trainingStore.isRunning"
        @click="handleStart"
        class="px-8 py-3 rounded-full bg-gradient-to-r from-purple-600 to-pink-600 hover:from-purple-500 hover:to-pink-500 text-white font-semibold text-lg transition-all shadow-lg shadow-purple-500/30"
      >
        开始训练
      </button>
      <button
        v-else
        @click="handleStop"
        class="px-8 py-3 rounded-full bg-white/10 hover:bg-white/20 text-white font-semibold text-lg transition-all"
      >
        结束训练
      </button>
    </div>

    <div class="w-full max-w-md bg-white/5 rounded-xl p-4">
      <div class="h-2 bg-white/10 rounded-full overflow-hidden">
        <div
          class="h-full bg-gradient-to-r from-purple-500 to-pink-500 transition-all duration-300"
          :style="{ width: `${trainingStore.progress}%` }"
        ></div>
      </div>
    </div>
  </div>
</template>
