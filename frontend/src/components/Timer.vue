<script setup lang="ts">
import { computed, onMounted, onUnmounted, ref, watch, nextTick } from 'vue'
import { useTrainingStore } from '@/store/training'
import { Play, Square, RotateCcw } from 'lucide-vue-next'
import confetti from 'canvas-confetti'
import gsap from 'gsap'

const trainingStore = useTrainingStore()

// Refs for DOM elements
const timerContainer = ref<HTMLElement | null>(null)
const innerCircle = ref<HTMLElement | null>(null)
const countdownNumber = ref<HTMLElement | null>(null)
const glowOrbs = ref<HTMLElement[]>([])
const pulseRings = ref<HTMLElement[]>([])

const isAnimating = ref(false)

// Phase configurations with enhanced colors
const phaseConfig = {
  idle: {
    primary: '#64748b',
    secondary: '#475569',
    glow: 'rgba(100, 116, 139, 0.3)',
    label: '准备开始',
    instruction: '点击开始按钮'
  },
  contract: {
    primary: '#f43f5e',
    secondary: '#ec4899',
    glow: 'rgba(244, 63, 94, 0.4)',
    label: '收缩',
    instruction: '用力收紧'
  },
  hold: {
    primary: '#f59e0b',
    secondary: '#f97316',
    glow: 'rgba(245, 158, 11, 0.4)',
    label: '保持',
    instruction: '保持用力'
  },
  relax: {
    primary: '#10b981',
    secondary: '#14b8a6',
    glow: 'rgba(16, 185, 129, 0.4)',
    label: '放松',
    instruction: '完全放松'
  }
}

const currentPhaseConfig = computed(() => phaseConfig[trainingStore.phase])

// SVG circle calculations
const radius = 140
const circumference = 2 * Math.PI * radius

// Smooth progress with GSAP
const smoothProgress = ref(0)
const displayCountdown = ref(trainingStore.countdown)

// GSAP Timeline for breathing animation
let breathingTimeline: gsap.core.Timeline | null = null
let phaseTransitionTimeline: gsap.core.Timeline | null = null

// Watch for phase changes - animate transitions
watch(() => trainingStore.phase, (newPhase, oldPhase) => {
  if (newPhase !== oldPhase) {
    animatePhaseTransition(newPhase)
  }
}, { immediate: true })

// Watch countdown changes - animate number
watch(() => trainingStore.countdown, (newVal, oldVal) => {
  if (newVal !== oldVal) {
    animateCountdown(newVal)
  }
})

// Watch progress changes - smooth animation
watch(() => trainingStore.phaseProgress, (newProgress) => {
  gsap.to(smoothProgress, {
    value: newProgress,
    duration: 0.3,
    ease: 'power2.out'
  })
})

// Watch running state
watch(() => trainingStore.isRunning, (isRunning) => {
  if (isRunning) {
    startBreathingAnimation()
  } else {
    stopBreathingAnimation()
  }
})

const strokeDashoffset = computed(() => {
  if (trainingStore.phase === 'idle') return circumference
  const progress = smoothProgress.value / 100
  return circumference * (1 - progress)
})

function animatePhaseTransition(phase: string) {
  if (phaseTransitionTimeline) {
    phaseTransitionTimeline.kill()
  }

  const config = phaseConfig[phase as keyof typeof phaseConfig]

  phaseTransitionTimeline = gsap.timeline()

  // Animate inner circle gradient
  if (innerCircle.value && phaseTransitionTimeline) {
    phaseTransitionTimeline.to(innerCircle.value, {
      '--primary-color': config.primary,
      '--secondary-color': config.secondary,
      duration: 0.5,
      ease: 'power2.inOut'
    }, 0)

    // Scale pulse on phase change
    phaseTransitionTimeline.fromTo(innerCircle.value,
      { scale: 1 },
      { scale: 1.05, duration: 0.15, ease: 'power2.out', yoyo: true, repeat: 1 },
      0
    )
  }

  // Animate glow orbs
  glowOrbs.value.forEach((orb, index) => {
    if (orb && phaseTransitionTimeline) {
      phaseTransitionTimeline.to(orb, {
        backgroundColor: config.glow,
        duration: 0.4,
        ease: 'power2.inOut'
      }, index * 0.05)
    }
  })

  // Animate pulse rings
  pulseRings.value.forEach((ring, index) => {
    if (ring && phaseTransitionTimeline) {
      phaseTransitionTimeline.fromTo(ring,
        { scale: 1, opacity: 0.6 },
        { scale: 1.5, opacity: 0, duration: 0.8, ease: 'power2.out' },
        index * 0.1
      )
    }
  })
}

function animateCountdown(value: number) {
  displayCountdown.value = value

  if (countdownNumber.value) {
    // Number pop animation
    gsap.fromTo(countdownNumber.value,
      { scale: 1.3, opacity: 0.5 },
      { scale: 1, opacity: 1, duration: 0.3, ease: 'back.out(1.7)' }
    )

    // Extra emphasis on last 3 seconds
    if (value <= 3 && value > 0 && trainingStore.isRunning) {
      gsap.to(countdownNumber.value, {
        textShadow: `0 0 40px ${currentPhaseConfig.value.primary}, 0 0 80px ${currentPhaseConfig.value.primary}`,
        duration: 0.2,
        yoyo: true,
        repeat: 1
      })
    }
  }
}

function startBreathingAnimation() {
  if (breathingTimeline) {
    breathingTimeline.kill()
  }

  breathingTimeline = gsap.timeline({ repeat: -1 })

  // Breathing scale animation
  if (timerContainer.value) {
    breathingTimeline.to(timerContainer.value, {
      scale: 1.02,
      duration: 2,
      ease: 'sine.inOut',
      yoyo: true,
      repeat: -1
    })
  }

  // Glow orbs rotation and pulse
  glowOrbs.value.forEach((orb, index) => {
    if (orb) {
      gsap.to(orb, {
        rotation: 360,
        duration: 20 + index * 5,
        ease: 'none',
        repeat: -1,
        transformOrigin: '50% 50%'
      })

      gsap.to(orb, {
        scale: 1.2,
        opacity: 0.8,
        duration: 2 + index * 0.5,
        ease: 'sine.inOut',
        yoyo: true,
        repeat: -1
      })
    }
  })
}

function stopBreathingAnimation() {
  if (breathingTimeline) {
    breathingTimeline.kill()
  }

  // Reset animations
  if (timerContainer.value) {
    gsap.to(timerContainer.value, {
      scale: 1,
      duration: 0.3,
      ease: 'power2.out'
    })
  }

  glowOrbs.value.forEach((orb) => {
    if (orb) {
      gsap.killTweensOf(orb)
      gsap.to(orb, {
        scale: 1,
        opacity: 0.5,
        duration: 0.3
      })
    }
  })
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

  // Start ripple effect
  createRippleEffect()

  trainingStore.start()
}

function createRippleEffect() {
  const ripples = document.querySelectorAll('.start-ripple')
  ripples.forEach((ripple, index) => {
    gsap.fromTo(ripple,
      { scale: 0, opacity: 0.8 },
      { scale: 3, opacity: 0, duration: 0.8, delay: index * 0.1, ease: 'power2.out' }
    )
  })
}

function handleStop() {
  const cycles = trainingStore.currentCycle
  trainingStore.stop()
  if (cycles > 0) {
    celebrateCompletion()
  }
}

function handleReset() {
  // Reset animation
  if (innerCircle.value) {
    gsap.fromTo(innerCircle.value,
      { rotation: 0 },
      { rotation: 360, duration: 0.5, ease: 'power2.inOut' }
    )
  }
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

// Collect glow orb refs
function setGlowOrbRef(el: HTMLElement | null, index: number) {
  if (el) {
    glowOrbs.value[index] = el
  }
}

function setPulseRingRef(el: HTMLElement | null, index: number) {
  if (el) {
    pulseRings.value[index] = el
  }
}

onMounted(() => {
  // Initial animation
  nextTick(() => {
    if (timerContainer.value) {
      gsap.fromTo(timerContainer.value,
        { scale: 0.8, opacity: 0 },
        { scale: 1, opacity: 1, duration: 0.6, ease: 'back.out(1.4)' }
      )
    }
  })
})

onUnmounted(() => {
  trainingStore.clearTimers()
  if (breathingTimeline) breathingTimeline.kill()
  if (phaseTransitionTimeline) phaseTransitionTimeline.kill()
  gsap.killTweensOf('*')
})
</script>

<template>
  <div class="timer-root">
    <!-- Main Timer Container -->
    <div ref="timerContainer" class="timer-container">
      <!-- Ambient Glow Orbs -->
      <div class="glow-orbs">
        <div
          v-for="i in 3"
          :key="i"
          :ref="(el) => setGlowOrbRef(el as HTMLElement, i - 1)"
          class="glow-orb"
          :class="`orb-${i}`"
          :style="{ backgroundColor: currentPhaseConfig.glow }"
        />
      </div>

      <!-- Pulse Rings (for phase transitions) -->
      <div class="pulse-rings">
        <div
          v-for="i in 3"
          :key="i"
          :ref="(el) => setPulseRingRef(el as HTMLElement, i - 1)"
          class="pulse-ring"
          :style="{ borderColor: currentPhaseConfig.primary }"
        />
      </div>

      <!-- Progress Ring Container -->
      <div class="ring-container">
        <!-- SVG Progress Ring -->
        <svg class="progress-svg" viewBox="0 0 320 320">
          <!-- Outer glow filter -->
          <defs>
            <filter id="glow" x="-50%" y="-50%" width="200%" height="200%">
              <feGaussianBlur stdDeviation="4" result="coloredBlur"/>
              <feMerge>
                <feMergeNode in="coloredBlur"/>
                <feMergeNode in="SourceGraphic"/>
              </feMerge>
            </filter>
            <linearGradient id="progressGradient" x1="0%" y1="0%" x2="100%" y2="100%">
              <stop offset="0%" :stop-color="currentPhaseConfig.primary" />
              <stop offset="100%" :stop-color="currentPhaseConfig.secondary" />
            </linearGradient>
          </defs>

          <!-- Background track -->
          <circle
            cx="160"
            cy="160"
            :r="radius"
            fill="none"
            stroke="rgba(255, 255, 255, 0.06)"
            stroke-width="8"
          />

          <!-- Animated progress ring -->
          <circle
            ref="progressRing"
            cx="160"
            cy="160"
            :r="radius"
            fill="none"
            stroke="url(#progressGradient)"
            stroke-width="8"
            stroke-linecap="round"
            :stroke-dasharray="circumference"
            :stroke-dashoffset="strokeDashoffset"
            class="progress-ring"
            filter="url(#glow)"
          />

          <!-- Progress head glow -->
          <circle
            v-if="trainingStore.phase !== 'idle'"
            :cx="160 + radius * Math.cos((Math.PI * 2 * smoothProgress / 100) - Math.PI / 2)"
            :cy="160 + radius * Math.sin((Math.PI * 2 * smoothProgress / 100) - Math.PI / 2)"
            r="6"
            :fill="currentPhaseConfig.primary"
            filter="url(#glow)"
            class="progress-head"
          />
        </svg>

        <!-- Inner Circle -->
        <div
          ref="innerCircle"
          class="inner-circle"
          :style="{
            '--primary-color': currentPhaseConfig.primary,
            '--secondary-color': currentPhaseConfig.secondary
          }"
        >
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
      </div>
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
        <!-- Ripple effects -->
        <span class="start-ripple" />
        <span class="start-ripple" />
        <span class="start-ripple" />

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

    <!-- Overall Progress Bar -->
    <div class="progress-bar-container">
      <div class="progress-bar-header">
        <span class="progress-bar-label">整体进度</span>
        <span class="progress-bar-value">{{ Math.round(trainingStore.progress) }}%</span>
      </div>
      <div class="progress-bar-track">
        <div
          class="progress-bar-fill"
          :class="{ 'is-running': trainingStore.isRunning }"
          :style="{ width: `${trainingStore.progress}%` }"
        >
          <div class="progress-bar-shimmer" />
        </div>
        <div
          v-if="trainingStore.isRunning && trainingStore.progress > 0"
          class="progress-bar-glow"
          :style="{ left: `${trainingStore.progress}%` }"
        />
      </div>
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

/* Timer Container */
.timer-container {
  position: relative;
  width: 320px;
  height: 320px;
}

@media (min-width: 640px) {
  .timer-container {
    width: 360px;
    height: 360px;
  }
}

/* Glow Orbs */
.glow-orbs {
  position: absolute;
  inset: -40px;
  pointer-events: none;
}

.glow-orb {
  position: absolute;
  border-radius: 50%;
  filter: blur(60px);
  opacity: 0.5;
  transition: background-color 0.5s ease;
}

.orb-1 {
  top: -20%;
  left: 10%;
  width: 200px;
  height: 200px;
}

.orb-2 {
  bottom: -10%;
  right: 5%;
  width: 180px;
  height: 180px;
}

.orb-3 {
  top: 30%;
  left: -15%;
  width: 150px;
  height: 150px;
}

/* Pulse Rings */
.pulse-rings {
  position: absolute;
  inset: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  pointer-events: none;
}

.pulse-ring {
  position: absolute;
  width: 100%;
  height: 100%;
  border-radius: 50%;
  border: 2px solid transparent;
  opacity: 0;
}

/* Ring Container */
.ring-container {
  position: relative;
  width: 100%;
  height: 100%;
}

.progress-svg {
  width: 100%;
  height: 100%;
  transform: rotate(-90deg);
}

.progress-ring {
  transition: stroke-dashoffset 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

.progress-head {
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
}

/* Inner Circle */
.inner-circle {
  position: absolute;
  inset: 24px;
  border-radius: 50%;
  background: linear-gradient(
    135deg,
    var(--primary-color, #64748b),
    var(--secondary-color, #475569)
  );
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  box-shadow:
    0 20px 60px rgba(0, 0, 0, 0.4),
    inset 0 2px 20px rgba(255, 255, 255, 0.1);
  transition: background 0.5s ease;
}

@media (min-width: 640px) {
  .inner-circle {
    inset: 28px;
  }
}

/* Countdown */
.countdown-wrapper {
  position: relative;
}

.countdown-number {
  font-size: 80px;
  font-weight: 700;
  color: #fff;
  line-height: 1;
  text-shadow: 0 4px 20px rgba(0, 0, 0, 0.3);
  font-feature-settings: 'tnum';
}

@media (min-width: 640px) {
  .countdown-number {
    font-size: 96px;
  }
}

/* Phase Info */
.phase-info {
  margin-top: 8px;
  text-align: center;
}

.phase-label {
  font-size: 22px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.95);
}

@media (min-width: 640px) {
  .phase-label {
    font-size: 26px;
  }
}

.phase-instruction {
  font-size: 14px;
  color: rgba(255, 255, 255, 0.6);
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

.start-ripple {
  position: absolute;
  top: 50%;
  left: 50%;
  width: 100%;
  height: 100%;
  border-radius: 14px;
  background: rgba(255, 255, 255, 0.3);
  transform: translate(-50%, -50%) scale(0);
  opacity: 0;
  pointer-events: none;
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

/* Progress Bar */
.progress-bar-container {
  width: 100%;
  max-width: 400px;
}

.progress-bar-header {
  display: flex;
  justify-content: space-between;
  margin-bottom: 8px;
}

.progress-bar-label {
  font-size: 13px;
  color: rgba(255, 255, 255, 0.5);
}

.progress-bar-value {
  font-size: 13px;
  font-weight: 600;
  color: rgba(255, 255, 255, 0.7);
  font-feature-settings: 'tnum';
}

.progress-bar-track {
  position: relative;
  height: 10px;
  background: rgba(255, 255, 255, 0.06);
  border-radius: 10px;
  overflow: hidden;
  border: 1px solid rgba(255, 255, 255, 0.05);
}

.progress-bar-fill {
  height: 100%;
  border-radius: 10px;
  background: linear-gradient(90deg, #6366f1, #8b5cf6, #a855f7);
  box-shadow: 0 0 20px rgba(139, 92, 246, 0.5);
  transition: width 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  position: relative;
  overflow: hidden;
}

.progress-bar-shimmer {
  position: absolute;
  inset: 0;
  background: linear-gradient(
    90deg,
    transparent 0%,
    rgba(255, 255, 255, 0.2) 50%,
    transparent 100%
  );
  background-size: 200% 100%;
  opacity: 0;
}

.progress-bar-fill.is-running .progress-bar-shimmer {
  opacity: 1;
  animation: shimmer 2s ease-in-out infinite;
}

@keyframes shimmer {
  0% { background-position: 200% 0; }
  100% { background-position: -200% 0; }
}

.progress-bar-glow {
  position: absolute;
  top: 50%;
  transform: translate(-50%, -50%);
  width: 16px;
  height: 16px;
  border-radius: 50%;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.9) 0%, transparent 70%);
  box-shadow: 0 0 20px 6px rgba(139, 92, 246, 0.6);
  transition: left 0.5s cubic-bezier(0.4, 0, 0.2, 1);
  animation: pulse-glow 1.5s ease-in-out infinite;
}

@keyframes pulse-glow {
  0%, 100% { opacity: 0.6; transform: translate(-50%, -50%) scale(1); }
  50% { opacity: 1; transform: translate(-50%, -50%) scale(1.2); }
}
</style>
