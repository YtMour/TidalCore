<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useTrainingStore, type TrainingPhase } from '@/store/training'
import gsap from 'gsap'

const trainingStore = useTrainingStore()

// Animation state
const scale = ref(1)
const breathScale = ref(1)
const time = ref(0)

// 波浪幅度随机化状态
const waveAmplitudeMultiplier = ref(1.0)
const breathMultiplier = ref(1.0)

// 特效波纹状态 - 每个波纹有独立的时间偏移，产生不同的波动
const waveRipples = ref<{
  id: number
  offset: number
  opacity: number
  strokeWidth: number
  timeOffset: number  // 时间偏移，让每个波纹有不同的波动相位
  waveIntensity: number // 波动强度
}[]>([])
let rippleIdCounter = 0
let rippleInterval: ReturnType<typeof setInterval> | null = null

// 阶段过渡状态
const phaseTransition = ref(1) // 0-1，用于平滑过渡
let phaseTransitionTween: gsap.core.Tween | null = null

// Phase configurations
const phaseConfig = {
  idle: {
    targetScale: 1.25,
    borderColor: '#7dd3fc',
    fillColor: 'rgba(56, 189, 248, 0.55)',
    fillColorInner: 'rgba(125, 211, 252, 0.4)',
    glowColor: 'rgba(125, 211, 252, 0.7)',
    breathDuration: 2.5,
    breathIntensity: 0.06,
    waveAmplitude: 8,
    waveDepth: 3,
    rippleColor: '#7dd3fc',
    rippleType: 'none' as const,
    glowIntensity: 0.4
  },
  contract: {
    targetScale: 0.85,
    borderColor: '#fda4af',
    fillColor: 'rgba(251, 113, 133, 0.5)',
    fillColorInner: 'rgba(253, 164, 175, 0.35)',
    glowColor: 'rgba(253, 164, 175, 0.65)',
    breathDuration: 2,
    breathIntensity: 0.04,
    waveAmplitude: 5,
    waveDepth: 2,
    rippleColor: '#fda4af',
    rippleType: 'none' as const,  // 收缩阶段不显示特效
    glowIntensity: 0.6
  },
  hold: {
    targetScale: 0.85,
    borderColor: '#fcd34d',
    fillColor: 'rgba(251, 191, 36, 0.5)',
    fillColorInner: 'rgba(252, 211, 77, 0.35)',
    glowColor: 'rgba(252, 211, 77, 0.65)',
    breathDuration: 1.5,
    breathIntensity: 0.05,
    waveAmplitude: 4,
    waveDepth: 1.5,
    rippleColor: '#fcd34d',
    rippleType: 'pulse' as const,
    glowIntensity: 0.7
  },
  relax: {
    targetScale: 1.5,
    borderColor: '#6ee7b7',
    fillColor: 'rgba(52, 211, 153, 0.5)',
    fillColorInner: 'rgba(110, 231, 183, 0.35)',
    glowColor: 'rgba(110, 231, 183, 0.7)',
    breathDuration: 2,
    breathIntensity: 0.07,
    waveAmplitude: 10,
    waveDepth: 4,
    rippleColor: '#6ee7b7',
    rippleType: 'expand' as const,
    glowIntensity: 0.8
  }
}

const currentConfig = computed(() => phaseConfig[trainingStore.phase])

// Animation references
let breathingTween: gsap.core.Tween | null = null
let scaleTween: gsap.core.Tween | null = null
let animationFrame: number | null = null

const combinedScale = computed(() => scale.value * breathScale.value)

// 生成流体球形路径 - 支持时间偏移和波动强度
function generateFluidBlobPath(
  t: number,
  amplitude: number,
  depth: number,
  radiusOffset: number = 0,
  timeOffset: number = 0,
  waveIntensity: number = 1
): string {
  const numPoints = 24
  const baseRadius = 115 + radiusOffset
  const cx = 150
  const cy = 150
  const points: { x: number; y: number }[] = []

  const effectiveAmplitude = amplitude * waveAmplitudeMultiplier.value * waveIntensity
  const effectiveTime = t + timeOffset

  for (let i = 0; i < numPoints; i++) {
    const angle = (i / numPoints) * Math.PI * 2

    // 波浪叠加
    const wave1 = Math.sin(angle * 2 + effectiveTime * 0.6) * effectiveAmplitude * 0.45
    const wave2 = Math.sin(angle * 3 - effectiveTime * 0.8) * effectiveAmplitude * 0.28
    const wave3 = Math.cos(angle * 5 + effectiveTime * 1.2) * effectiveAmplitude * 0.18
    const wave4 = Math.sin(angle * 7 - effectiveTime * 1.5) * effectiveAmplitude * 0.09

    let waveSum = wave1 + wave2 + wave3 + wave4
    waveSum += effectiveAmplitude * 0.15

    if (waveSum < 0) {
      waveSum = Math.max(waveSum, -depth)
    }

    const radius = baseRadius + waveSum

    points.push({
      x: cx + Math.cos(angle) * radius,
      y: cy + Math.sin(angle) * radius
    })
  }

  const firstPoint = points[0]!
  let path = `M ${firstPoint.x} ${firstPoint.y}`

  for (let i = 0; i < numPoints; i++) {
    const p0 = points[(i - 1 + numPoints) % numPoints]!
    const p1 = points[i]!
    const p2 = points[(i + 1) % numPoints]!
    const p3 = points[(i + 2) % numPoints]!

    const cp1x = p1.x + (p2.x - p0.x) / 6
    const cp1y = p1.y + (p2.y - p0.y) / 6
    const cp2x = p2.x - (p3.x - p1.x) / 6
    const cp2y = p2.y - (p3.y - p1.y) / 6

    path += ` C ${cp1x} ${cp1y}, ${cp2x} ${cp2y}, ${p2.x} ${p2.y}`
  }

  return path
}

// 主球体路径
const blobPath = computed(() =>
  generateFluidBlobPath(time.value, currentConfig.value.waveAmplitude, currentConfig.value.waveDepth)
)

// 为每个波纹生成带波动的路径
function getRipplePath(ripple: typeof waveRipples.value[0]): string {
  return generateFluidBlobPath(
    time.value,
    currentConfig.value.waveAmplitude * 0.8,
    currentConfig.value.waveDepth * 0.6,
    ripple.offset,
    ripple.timeOffset,
    ripple.waveIntensity
  )
}

// 波浪幅度随机化
let waveAmplitudeTween: gsap.core.Tween | null = null

function startWaveAmplitudeRandomization() {
  if (waveAmplitudeTween) waveAmplitudeTween.kill()

  function animateToRandomAmplitude() {
    const targetValue = 1.0 + Math.random() * 0.25
    const duration = 2 + Math.random() * 2

    waveAmplitudeTween = gsap.to(waveAmplitudeMultiplier, {
      value: targetValue,
      duration: duration,
      ease: 'sine.inOut',
      onComplete: animateToRandomAmplitude
    })
  }

  animateToRandomAmplitude()
}

// 呼吸效果随机化
let breathRandomTween: gsap.core.Tween | null = null

function startBreathRandomization() {
  if (breathRandomTween) breathRandomTween.kill()

  function animateToRandomBreath() {
    const targetValue = 1.0 + Math.random() * 0.10
    const duration = 1.5 + Math.random() * 1.5

    breathRandomTween = gsap.to(breathMultiplier, {
      value: targetValue,
      duration: duration,
      ease: 'sine.inOut',
      onComplete: animateToRandomBreath
    })
  }

  animateToRandomBreath()
}

function startBreathing() {
  if (breathingTween) breathingTween.kill()

  const config = currentConfig.value
  const effectiveIntensity = config.breathIntensity * breathMultiplier.value

  breathingTween = gsap.to(breathScale, {
    value: 1 + effectiveIntensity,
    duration: config.breathDuration / 2,
    ease: 'sine.inOut',
    yoyo: true,
    repeat: -1
  })
}

function animateScaleToPhase(phase: TrainingPhase) {
  if (scaleTween) scaleTween.kill()

  const config = phaseConfig[phase]
  const duration = getDurationForPhase(phase)

  // 阶段过渡动画
  if (phaseTransitionTween) phaseTransitionTween.kill()
  phaseTransition.value = 0
  phaseTransitionTween = gsap.to(phaseTransition, {
    value: 1,
    duration: 0.8,
    ease: 'power2.inOut'
  })

  if (phase === 'relax') {
    scaleTween = gsap.to(scale, {
      value: config.targetScale,
      duration: duration * 0.4,
      ease: 'power2.out',
      onComplete: () => {
        gsap.to(scale, {
          value: phaseConfig.idle.targetScale,
          duration: duration * 0.6,
          ease: 'power2.inOut'
        })
      }
    })
  } else if (phase === 'idle') {
    scaleTween = gsap.to(scale, {
      value: config.targetScale,
      duration: 0.8,
      ease: 'power2.out'
    })
  } else if (phase === 'contract') {
    scaleTween = gsap.to(scale, {
      value: config.targetScale,
      duration: duration * 0.9,
      ease: 'power1.inOut'
    })
  } else if (phase === 'hold') {
    scaleTween = gsap.to(scale, {
      value: config.targetScale,
      duration: 0.5,
      ease: 'power2.out'
    })
  }

  startBreathing()
}

function getDurationForPhase(phase: TrainingPhase): number {
  switch (phase) {
    case 'contract': return trainingStore.settings.contractTime
    case 'hold': return trainingStore.settings.holdTime
    case 'relax': return trainingStore.settings.relaxTime
    default: return 1
  }
}

// 波纹特效控制
function startRippleEffect(phase: TrainingPhase) {
  // 立即停止现有波纹
  stopRippleEffect()

  const config = phaseConfig[phase]
  if (config.rippleType === 'none') return

  // 立即创建第一个波纹，不延迟
  createRipple(config.rippleType)

  // 设置定时器创建后续波纹
  if (config.rippleType === 'pulse') {
    rippleInterval = setInterval(() => {
      if (waveRipples.value.length >= 2) return
      createRipple('pulse')
    }, 1000)
  } else if (config.rippleType === 'expand') {
    rippleInterval = setInterval(() => {
      if (waveRipples.value.length >= 2) return
      createRipple('expand')
    }, 1200)
  }
}

function createRipple(type: 'pulse' | 'expand') {
  const id = rippleIdCounter++
  const timeOffset = Math.random() * Math.PI * 2

  if (type === 'pulse') {
    // 保持：波动更明显
    waveRipples.value.push({
      id,
      offset: 5,
      opacity: 0.55,
      strokeWidth: 2.5,
      timeOffset,
      waveIntensity: 1.8  // 大幅增强波动
    })

    const ripple = waveRipples.value.find(r => r.id === id)
    if (ripple) {
      gsap.to(ripple, {
        offset: 45,
        opacity: 0,
        strokeWidth: 1,
        waveIntensity: 1.4,  // 结束时仍保持较强波动
        duration: 1.5,
        ease: 'sine.out',
        onComplete: () => {
          const idx = waveRipples.value.findIndex(r => r.id === id)
          if (idx > -1) waveRipples.value.splice(idx, 1)
        }
      })
    }
  } else if (type === 'expand') {
    // 放松：范围加大，波动加大
    waveRipples.value.push({
      id,
      offset: 3,
      opacity: 0.6,
      strokeWidth: 2.5,
      timeOffset,
      waveIntensity: 1.3
    })

    const ripple = waveRipples.value.find(r => r.id === id)
    if (ripple) {
      gsap.to(ripple, {
        offset: 90,
        opacity: 0,
        strokeWidth: 0.5,
        waveIntensity: 0.7,
        duration: 2,
        ease: 'power1.out',
        onComplete: () => {
          const idx = waveRipples.value.findIndex(r => r.id === id)
          if (idx > -1) waveRipples.value.splice(idx, 1)
        }
      })
    }
  }
}

function stopRippleEffect() {
  if (rippleInterval) {
    clearInterval(rippleInterval)
    rippleInterval = null
  }
  // 立即清除所有波纹
  waveRipples.value = []
}

// Animation loop
function animate() {
  time.value += 0.02 // 稍微降低速度
  animationFrame = requestAnimationFrame(animate)
}

// Watch phase changes
watch(() => trainingStore.phase, (newPhase) => {
  animateScaleToPhase(newPhase)
  startRippleEffect(newPhase)
}, { immediate: true })

onMounted(() => {
  startBreathing()
  startWaveAmplitudeRandomization()
  startBreathRandomization()
  animate()
})

onUnmounted(() => {
  if (breathingTween) breathingTween.kill()
  if (scaleTween) scaleTween.kill()
  if (waveAmplitudeTween) waveAmplitudeTween.kill()
  if (breathRandomTween) breathRandomTween.kill()
  if (phaseTransitionTween) phaseTransitionTween.kill()
  if (animationFrame) cancelAnimationFrame(animationFrame)
  stopRippleEffect()
})

defineExpose({
  combinedScale,
  currentConfig
})
</script>

<template>
  <div class="fluid-ball-container">
    <svg
      class="fluid-svg"
      viewBox="0 0 300 300"
      :style="{ transform: `scale(${combinedScale})` }"
    >
      <defs>
        <radialGradient id="ballFill" cx="35%" cy="35%" r="65%">
          <stop offset="0%" stop-color="rgba(255, 255, 255, 0.35)" />
          <stop offset="40%" :stop-color="currentConfig.fillColorInner" />
          <stop offset="100%" :stop-color="currentConfig.fillColor" />
        </radialGradient>

        <filter id="outerGlow" x="-100%" y="-100%" width="300%" height="300%">
          <feGaussianBlur in="SourceGraphic" stdDeviation="8" result="blur" />
          <feColorMatrix in="blur" type="matrix" values="1 0 0 0 0  0 1 0 0 0  0 0 1 0 0  0 0 0 0.6 0" result="glow" />
          <feMerge>
            <feMergeNode in="glow" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>

        <filter id="borderGlow" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur in="SourceGraphic" stdDeviation="3" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>

        <filter id="rippleGlow" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur in="SourceGraphic" stdDeviation="3" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>

      <!-- 波纹层 - 带流体波动的波纹 -->
      <g class="ripple-layer">
        <path
          v-for="ripple in waveRipples"
          :key="ripple.id"
          :d="getRipplePath(ripple)"
          fill="none"
          :stroke="currentConfig.rippleColor"
          :stroke-width="ripple.strokeWidth"
          :opacity="ripple.opacity"
          filter="url(#rippleGlow)"
          stroke-linecap="round"
        />
      </g>

      <!-- 外发光层 -->
      <path
        :d="blobPath"
        :fill="currentConfig.glowColor"
        filter="url(#outerGlow)"
        :opacity="currentConfig.glowIntensity"
      />

      <!-- 背景填充 -->
      <path
        :d="blobPath"
        fill="url(#ballFill)"
      />

      <!-- 发光边框 -->
      <path
        :d="blobPath"
        fill="none"
        :stroke="currentConfig.borderColor"
        stroke-width="2.5"
        filter="url(#borderGlow)"
      />

      <!-- 主高光 -->
      <ellipse
        cx="115"
        cy="105"
        rx="35"
        ry="20"
        fill="rgba(255, 255, 255, 0.18)"
        transform="rotate(-25, 115, 105)"
      />

      <!-- 次高光 -->
      <ellipse
        cx="130"
        cy="115"
        rx="15"
        ry="8"
        fill="rgba(255, 255, 255, 0.25)"
        transform="rotate(-25, 130, 115)"
      />
    </svg>

    <div class="content-overlay">
      <slot />
    </div>
  </div>
</template>

<style scoped>
.fluid-ball-container {
  position: relative;
  width: 280px;
  height: 280px;
  display: flex;
  align-items: center;
  justify-content: center;
}

@media (min-width: 640px) {
  .fluid-ball-container {
    width: 320px;
    height: 320px;
  }
}

.fluid-svg {
  width: 100%;
  height: 100%;
  transition: transform 0.05s linear;
  overflow: visible;
}

.fluid-svg path {
  transition: fill 0.8s ease, stroke 0.8s ease;
}

.ripple-layer path {
  transition: stroke 0.8s ease;
}

.content-overlay {
  position: absolute;
  inset: 0;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  z-index: 10;
  pointer-events: none;
}

.content-overlay > :deep(*) {
  pointer-events: auto;
}
</style>
