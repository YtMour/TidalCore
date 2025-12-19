<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch } from 'vue'
import { useTrainingStore } from '@/store/training'
import { Play, Square, RotateCcw } from 'lucide-vue-next'
import confetti from 'canvas-confetti'
import gsap from 'gsap'
import FluidBall from './FluidBall.vue'

const trainingStore = useTrainingStore()

// Refs for DOM elements
const countdownNumber = ref<HTMLElement | null>(null)

const isAnimating = ref(false)

// Phase configurations
const phaseConfig = {
  idle: {
    label: '准备开始',
    instruction: '点击开始按钮'
  },
  contract: {
    label: '收缩',
    instruction: '用力收紧'
  },
  hold: {
    label: '保持',
    instruction: '保持用力'
  },
  relax: {
    label: '放松',
    instruction: '完全放松'
  }
}

const currentPhaseConfig = computed(() => phaseConfig[trainingStore.phase])

// Smooth countdown display
const displayCountdown = ref(trainingStore.countdown)

// Watch countdown changes - animate number
watch(() => trainingStore.countdown, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    animateCountdown(newVal)
  }
})

// Watch isCompleted - 训练自动完成时触发庆祝效果
watch(() => trainingStore.isCompleted, (completed) => {
  if (completed) {
    celebrateCompletion()
  }
})

function animateCountdown(value: number) {
  displayCountdown.value = value

  if (countdownNumber.value) {
    // Number pop animation
    gsap.fromTo(countdownNumber.value,
      { scale: 1.2, opacity: 0.6 },
      { scale: 1, opacity: 1, duration: 0.25, ease: 'back.out(1.7)' }
    )

    // Extra emphasis on last 3 seconds
    if (value <= 3 && value > 0 && trainingStore.isRunning) {
      gsap.to(countdownNumber.value, {
        textShadow: '0 0 30px currentColor, 0 0 60px currentColor',
        duration: 0.15,
        yoyo: true,
        repeat: 1
      })
    }
  }
}

function handleStart() {
  isAnimating.value = true

  // Button press animation
  gsap.timeline()
    .to('.start-btn', { scale: 0.95, duration: 0.1 })
    .to('.start-btn', { scale: 1, duration: 0.2, ease: 'back.out(1.7)' })
    .eventCallback('onComplete', () => {
      isAnimating.value = false
    })

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
  const colors = ['#38bdf8', '#22d3ee', '#10b981', '#f59e0b']
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

onMounted(() => {
  displayCountdown.value = trainingStore.countdown
})

onUnmounted(() => {
  trainingStore.clearTimers()
  if (countdownNumber.value) {
    gsap.killTweensOf(countdownNumber.value)
  }
})
</script>

<template>
  <div class="timer-root">
    <!-- Fluid Ball Timer -->
    <div class="timer-wrapper">
      <FluidBall>
        <div class="timer-content">
          <!-- Countdown Display -->
          <div class="countdown-wrapper">
            <span ref="countdownNumber" class="countdown-number">
              {{ displayCountdown }}
            </span>
          </div>

          <!-- Phase Info -->
          <div class="phase-info">
            <div class="phase-label">{{ currentPhaseConfig.label }}</div>
            <div class="phase-instruction">{{ currentPhaseConfig.instruction }}</div>
          </div>
        </div>
      </FluidBall>
    </div>

    <!-- Stats Row -->
    <div class="stats-row">
      <div class="stat-item">
        <div class="stat-value">{{ trainingStore.currentCycle }}</div>
        <div class="stat-label">/ {{ trainingStore.settings.cycles }} 循环</div>
      </div>

      <div class="stat-divider" />

      <div class="stat-item">
        <div class="stat-value">{{ formatTime(trainingStore.totalDuration) }}</div>
        <div class="stat-label">训练时长</div>
      </div>
    </div>

    <!-- Control Buttons -->
    <div class="control-buttons">
      <button
        v-if="!trainingStore.isRunning"
        @click="handleStart"
        class="start-btn"
        :class="{ 'is-animating': isAnimating }"
      >
        <span class="btn-glow" />
        <span class="btn-content">
          <Play class="btn-icon" />
          <span class="btn-text">开始训练</span>
        </span>
      </button>

      <template v-else>
        <button @click="handleStop" class="stop-btn">
          <Square class="btn-icon filled" />
          <span class="btn-text">结束训练</span>
        </button>
      </template>

      <button
        v-if="trainingStore.currentCycle > 0 && !trainingStore.isRunning"
        @click="handleReset"
        class="reset-btn"
        title="重置"
      >
        <RotateCcw class="reset-icon" />
      </button>
    </div>
  </div>
</template>

<style scoped>
.timer-root {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 32px;
}

.timer-wrapper {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
}

/* Timer Content inside FluidBall */
.timer-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  text-align: center;
}

/* Countdown */
.countdown-wrapper {
  position: relative;
}

.countdown-number {
  font-size: 72px;
  font-weight: 700;
  color: #fff;
  line-height: 1;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.4);
  font-feature-settings: 'tnum';
}

@media (min-width: 640px) {
  .countdown-number {
    font-size: 88px;
  }
}

/* Phase Info */
.phase-info {
  margin-top: 8px;
  text-align: center;
}

.phase-label {
  font-size: 20px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.95);
}

@media (min-width: 640px) {
  .phase-label {
    font-size: 24px;
  }
}

.phase-instruction {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.7);
  margin-top: 4px;
}

/* Stats Row */
.stats-row {
  display: flex;
  align-items: center;
  gap: 24px;
}

@media (min-width: 640px) {
  .stats-row {
    gap: 40px;
  }
}

.stat-item {
  text-align: center;
}

.stat-value {
  font-size: 28px;
  font-weight: 700;
  color: #fff;
  font-feature-settings: 'tnum';
}

@media (min-width: 640px) {
  .stat-value {
    font-size: 32px;
  }
}

.stat-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
  margin-top: 4px;
}

.stat-divider {
  width: 1px;
  height: 48px;
  background: linear-gradient(to bottom, transparent, rgba(255, 255, 255, 0.2), transparent);
}

/* Control Buttons */
.control-buttons {
  display: flex;
  align-items: center;
  gap: 12px;
}

.start-btn {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  padding: 16px 40px;
  border: none;
  border-radius: 14px;
  background: linear-gradient(135deg, #0ea5e9, #0284c7);
  cursor: pointer;
  overflow: hidden;
  transition: all 0.3s ease;
  box-shadow:
    0 4px 20px rgba(14, 165, 233, 0.4),
    0 0 0 1px rgba(255, 255, 255, 0.1) inset;
}

.start-btn::before {
  content: '';
  position: absolute;
  bottom: -100%;
  left: -10%;
  right: -10%;
  height: 200%;
  background: linear-gradient(180deg, transparent 40%, rgba(56, 189, 248, 0.3) 48%, rgba(34, 211, 238, 0.5) 52%, rgba(56, 189, 248, 0.3) 56%, transparent 64%);
  animation: btn-ocean-wave 3s ease-in-out infinite;
  pointer-events: none;
}

@keyframes btn-ocean-wave {
  0%, 100% { transform: translateY(30%); }
  50% { transform: translateY(20%); }
}

.start-btn:hover {
  transform: translateY(-2px);
  box-shadow:
    0 8px 30px rgba(14, 165, 233, 0.5),
    0 0 0 1px rgba(255, 255, 255, 0.15) inset;
}

.start-btn:active,
.start-btn.is-animating {
  transform: scale(0.97);
}

.btn-glow {
  position: absolute;
  inset: 0;
  background: linear-gradient(180deg, rgba(255, 255, 255, 0.15) 0%, transparent 50%);
  opacity: 0;
  transition: opacity 0.3s ease;
}

.start-btn:hover .btn-glow {
  opacity: 1;
}

.btn-content {
  position: relative;
  display: flex;
  align-items: center;
  gap: 10px;
  z-index: 1;
}

.btn-icon {
  width: 22px;
  height: 22px;
  color: #fff;
}

.start-btn .btn-icon {
  fill: currentColor;
}

.btn-text {
  font-size: 16px;
  font-weight: 600;
  color: #fff;
  letter-spacing: 0.5px;
}

.stop-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 16px 36px;
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.04);
  cursor: pointer;
  transition: all 0.2s ease;
}

.stop-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.15);
}

.stop-btn .btn-icon {
  width: 18px;
  height: 18px;
  color: rgba(255, 255, 255, 0.8);
}

.stop-btn .btn-icon.filled {
  fill: currentColor;
}

.reset-btn {
  display: flex;
  align-items: center;
  justify-content: center;
  width: 52px;
  height: 52px;
  border: 1px solid rgba(255, 255, 255, 0.08);
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.03);
  cursor: pointer;
  transition: all 0.2s ease;
}

.reset-btn:hover {
  background: rgba(255, 255, 255, 0.08);
  border-color: rgba(255, 255, 255, 0.12);
}

.reset-icon {
  width: 20px;
  height: 20px;
  color: rgba(255, 255, 255, 0.6);
  transition: color 0.2s ease;
}

.reset-btn:hover .reset-icon {
  color: #fff;
}
</style>
