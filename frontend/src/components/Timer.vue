<script setup lang="ts">
import { computed, onUnmounted, ref, watch } from 'vue'
import { useTrainingStore } from '@/store/training'
import { Play, Square, RotateCcw } from 'lucide-vue-next'
import confetti from 'canvas-confetti'

const trainingStore = useTrainingStore()

const isAnimating = ref(false)

// Phase configurations
const phaseConfig = {
  idle: {
    gradient: 'from-slate-700 via-slate-600 to-slate-700',
    ring: '#64748b',
    label: '准备开始',
    instruction: '点击开始按钮'
  },
  contract: {
    gradient: 'from-rose-600 via-pink-500 to-rose-600',
    ring: '#f43f5e',
    label: '收缩',
    instruction: '用力收紧'
  },
  hold: {
    gradient: 'from-amber-500 via-orange-400 to-amber-500',
    ring: '#f59e0b',
    label: '保持',
    instruction: '保持用力'
  },
  relax: {
    gradient: 'from-emerald-600 via-teal-500 to-emerald-600',
    ring: '#10b981',
    label: '放松',
    instruction: '完全放松'
  }
}

const currentPhaseConfig = computed(() => phaseConfig[trainingStore.phase])

// SVG circle calculations
const radius = 140
const circumference = 2 * Math.PI * radius

// 使用 ref 实现平滑过渡的进度
const smoothProgress = ref(0)
let animationFrame: number | null = null

// 监听进度变化，实现平滑动画
watch(() => trainingStore.phaseProgress, (newProgress) => {
  if (animationFrame) {
    cancelAnimationFrame(animationFrame)
  }

  const startProgress = smoothProgress.value
  const targetProgress = newProgress
  const startTime = performance.now()
  const duration = 100 // 100ms 平滑过渡

  function animate(currentTime: number) {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)

    // 使用 easeOutQuad 缓动函数
    const eased = 1 - (1 - progress) * (1 - progress)
    smoothProgress.value = startProgress + (targetProgress - startProgress) * eased

    if (progress < 1) {
      animationFrame = requestAnimationFrame(animate)
    }
  }

  animationFrame = requestAnimationFrame(animate)
}, { immediate: true })

const strokeDashoffset = computed(() => {
  if (trainingStore.phase === 'idle') return circumference
  const progress = smoothProgress.value / 100
  return circumference * (1 - progress)
})

// Breathing animation scale
const breatheScale = computed(() => {
  if (!trainingStore.isRunning) return 1
  if (trainingStore.phase === 'contract') return 1.05
  if (trainingStore.phase === 'relax') return 0.95
  return 1
})

function handleStart() {
  isAnimating.value = true
  setTimeout(() => {
    isAnimating.value = false
  }, 300)
  trainingStore.start()
}

function handleStop() {
  const cycles = trainingStore.currentCycle
  trainingStore.stop()
  if (cycles > 0) {
    celebrateCompletion()
  }
}

function handleReset() {
  trainingStore.reset()
}

function celebrateCompletion() {
  const colors = ['#a855f7', '#ec4899', '#f59e0b', '#10b981']
  const duration = 3000
  const end = Date.now() + duration

  function frame() {
    confetti({
      particleCount: 3,
      angle: 60,
      spread: 55,
      origin: { x: 0, y: 0.7 },
      colors
    })
    confetti({
      particleCount: 3,
      angle: 120,
      spread: 55,
      origin: { x: 1, y: 0.7 },
      colors
    })

    if (Date.now() < end) {
      requestAnimationFrame(frame)
    }
  }

  // Initial burst
  confetti({
    particleCount: 100,
    spread: 100,
    origin: { y: 0.6 },
    colors
  })

  frame()
}

function formatTime(seconds: number): string {
  const mins = Math.floor(seconds / 60)
  const secs = seconds % 60
  return `${mins}:${secs.toString().padStart(2, '0')}`
}

onUnmounted(() => {
  trainingStore.clearTimers()
  if (animationFrame) {
    cancelAnimationFrame(animationFrame)
  }
})
</script>

<template>
  <div class="flex flex-col items-center gap-8">
    <!-- Main Timer Circle -->
    <div
      class="relative w-72 h-72 sm:w-80 sm:h-80 transition-transform duration-500 ease-out"
      :style="{ transform: `scale(${breatheScale})` }"
    >
      <!-- Outer glow -->
      <div
        class="absolute inset-0 rounded-full blur-2xl opacity-30 transition-all duration-500"
        :class="trainingStore.isRunning ? 'animate-pulse-glow' : ''"
        :style="{
          background: `radial-gradient(circle, ${currentPhaseConfig.ring}40 0%, transparent 70%)`
        }"
      ></div>

      <!-- SVG Progress Ring -->
      <svg class="w-full h-full -rotate-90 relative z-10" viewBox="0 0 320 320">
        <!-- Background circle -->
        <circle
          cx="160"
          cy="160"
          :r="radius"
          fill="none"
          stroke="currentColor"
          stroke-width="6"
          class="text-white/5"
        />
        <!-- Progress circle -->
        <circle
          cx="160"
          cy="160"
          :r="radius"
          fill="none"
          :stroke="currentPhaseConfig.ring"
          stroke-width="6"
          stroke-linecap="round"
          :stroke-dasharray="circumference"
          :stroke-dashoffset="strokeDashoffset"
          class="progress-ring drop-shadow-lg"
          :style="{
            filter: `drop-shadow(0 0 10px ${currentPhaseConfig.ring}60)`
          }"
        />
      </svg>

      <!-- Inner content -->
      <div
        class="absolute inset-4 sm:inset-5 rounded-full bg-gradient-to-br flex flex-col items-center justify-center transition-all duration-500 shadow-2xl"
        :class="currentPhaseConfig.gradient"
      >
        <!-- Countdown number -->
        <div class="relative">
          <span
            class="text-7xl sm:text-8xl font-bold text-white tabular-nums tracking-tight"
            :class="{ 'animate-pulse': trainingStore.isRunning && trainingStore.countdown <= 3 }"
          >
            {{ trainingStore.countdown }}
          </span>
        </div>

        <!-- Phase label -->
        <div class="mt-2 text-center">
          <div class="text-xl sm:text-2xl font-semibold text-white/90">
            {{ currentPhaseConfig.label }}
          </div>
          <div class="text-sm text-white/60 mt-1">
            {{ currentPhaseConfig.instruction }}
          </div>
        </div>
      </div>
    </div>

    <!-- Stats Row -->
    <div class="flex items-center gap-6 sm:gap-10">
      <div class="text-center">
        <div class="text-3xl sm:text-4xl font-bold text-white tabular-nums">
          {{ trainingStore.currentCycle }}
        </div>
        <div class="text-sm text-white/50 mt-1">
          / {{ trainingStore.settings.cycles }} 循环
        </div>
      </div>

      <div class="w-px h-14 bg-gradient-to-b from-transparent via-white/20 to-transparent"></div>

      <div class="text-center">
        <div class="text-3xl sm:text-4xl font-bold text-white tabular-nums">
          {{ formatTime(trainingStore.totalDuration) }}
        </div>
        <div class="text-sm text-white/50 mt-1">训练时长</div>
      </div>
    </div>

    <!-- Control Buttons -->
    <div class="flex items-center gap-4">
      <button
        v-if="!trainingStore.isRunning"
        @click="handleStart"
        class="group relative btn-gradient px-8 py-4 rounded-xl text-white font-semibold text-lg flex items-center gap-3 overflow-hidden"
        :class="{ 'scale-95': isAnimating }"
      >
        <Play class="w-6 h-6 fill-current" />
        <span>开始训练</span>
      </button>

      <template v-else>
        <button
          @click="handleStop"
          class="px-8 py-4 rounded-xl bg-white/[0.06] hover:bg-white/[0.1] text-white font-semibold text-lg flex items-center gap-3 transition-all duration-200 border border-white/[0.08] hover:border-white/[0.12]"
        >
          <Square class="w-5 h-5 fill-current" />
          <span>结束训练</span>
        </button>
      </template>

      <button
        v-if="trainingStore.currentCycle > 0 && !trainingStore.isRunning"
        @click="handleReset"
        class="icon-btn w-12 h-12 rounded-xl bg-white/[0.04] hover:bg-white/[0.08] border border-white/[0.06]"
        title="重置"
      >
        <RotateCcw class="w-5 h-5" />
      </button>
    </div>

    <!-- Overall Progress Bar -->
    <div class="w-full max-w-md">
      <div class="flex items-center justify-between text-sm text-white/50 mb-2">
        <span>整体进度</span>
        <span class="tabular-nums">{{ Math.round(trainingStore.progress) }}%</span>
      </div>
      <div class="h-2 bg-white/[0.04] rounded-lg overflow-hidden">
        <div
          class="h-full rounded-lg progress-bar-animated"
          :style="{
            width: `${trainingStore.progress}%`,
            background: 'linear-gradient(90deg, rgb(99, 102, 241), rgb(139, 92, 246))',
            transition: 'width 0.5s cubic-bezier(0.4, 0, 0.2, 1)'
          }"
        ></div>
      </div>
    </div>
  </div>
</template>
