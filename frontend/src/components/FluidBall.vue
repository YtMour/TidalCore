<script setup lang="ts">
import { ref, computed, watch, onMounted, onUnmounted } from 'vue'
import { useTrainingStore, type TrainingPhase } from '@/store/training'
import gsap from 'gsap'

const trainingStore = useTrainingStore()

// Animation state
const scale = ref(1)
const breathScale = ref(1)
const time = ref(0)

// Phase configurations - 优化后的颜色和体积设置
// idle: 1.25, contract: 0.85, hold: 0.85, relax: 1.5
const phaseConfig = {
  idle: {
    targetScale: 1.25,
    // 清澈的海蓝色调
    borderColor: '#60a5fa',
    fillColor: 'rgba(37, 99, 235, 0.7)',
    glowColor: 'rgba(96, 165, 250, 0.6)',
    breathDuration: 2.5,
    breathIntensity: 0.06,  // 加大呼吸幅度
    waveAmplitude: 8,
    waveDepth: 4
  },
  contract: {
    targetScale: 0.85,
    // 柔和的粉红色调
    borderColor: '#fb7185',
    fillColor: 'rgba(225, 29, 72, 0.6)',
    glowColor: 'rgba(251, 113, 133, 0.5)',
    breathDuration: 2,
    breathIntensity: 0.04,
    waveAmplitude: 5,
    waveDepth: 2.5
  },
  hold: {
    targetScale: 0.85,
    // 温暖的琥珀色调
    borderColor: '#fbbf24',
    fillColor: 'rgba(217, 119, 6, 0.6)',
    glowColor: 'rgba(251, 191, 36, 0.5)',
    breathDuration: 1.5,
    breathIntensity: 0.05,
    waveAmplitude: 4,
    waveDepth: 2
  },
  relax: {
    targetScale: 1.5,
    // 清新的翠绿色调
    borderColor: '#34d399',
    fillColor: 'rgba(16, 185, 129, 0.6)',
    glowColor: 'rgba(52, 211, 153, 0.5)',
    breathDuration: 2,
    breathIntensity: 0.07,  // 放松时呼吸最大
    waveAmplitude: 10,
    waveDepth: 5
  }
}

const currentConfig = computed(() => phaseConfig[trainingStore.phase])

// Animation references
let breathingTween: gsap.core.Tween | null = null
let scaleTween: gsap.core.Tween | null = null
let animationFrame: number | null = null

// Combined scale for the ball
const combinedScale = computed(() => scale.value * breathScale.value)

// 生成海浪式流体球形路径 - 以圆为基础，凸起大于凹陷
function generateFluidBlobPath(t: number, amplitude: number, depth: number): string {
  const numPoints = 16
  const baseRadius = 115
  const cx = 150
  const cy = 150
  const points: { x: number; y: number }[] = []

  for (let i = 0; i < numPoints; i++) {
    const angle = (i / numPoints) * Math.PI * 2

    // 多层波浪叠加 - 主要向外凸起
    const wave1 = Math.sin(angle * 3 + t * 0.8) * amplitude * 0.5
    const wave2 = Math.sin(angle * 2 - t * 0.5) * amplitude * 0.3
    const wave3 = Math.cos(angle * 4 + t * 1.1) * amplitude * 0.2

    // 合并波浪，正值凸起，负值凹陷（限制凹陷深度）
    let waveSum = wave1 + wave2 + wave3

    // 限制凹陷深度，保持形体
    if (waveSum < 0) {
      waveSum = Math.max(waveSum, -depth)
    }

    const radius = baseRadius + waveSum

    points.push({
      x: cx + Math.cos(angle) * radius,
      y: cy + Math.sin(angle) * radius
    })
  }

  // 使用 Catmull-Rom 曲线生成平滑路径
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

const blobPath = computed(() =>
  generateFluidBlobPath(time.value, currentConfig.value.waveAmplitude, currentConfig.value.waveDepth)
)

// Start breathing animation
function startBreathing() {
  if (breathingTween) breathingTween.kill()

  const config = currentConfig.value
  breathingTween = gsap.to(breathScale, {
    value: 1 + config.breathIntensity,
    duration: config.breathDuration / 2,
    ease: 'sine.inOut',
    yoyo: true,
    repeat: -1
  })
}

// Animate scale change for phase transitions
function animateScaleToPhase(phase: TrainingPhase) {
  if (scaleTween) scaleTween.kill()

  const config = phaseConfig[phase]
  const duration = getDurationForPhase(phase)

  if (phase === 'relax') {
    // 放松：快速释放到最大，然后缓慢回到idle大小
    scaleTween = gsap.to(scale, {
      value: config.targetScale,
      duration: duration * 0.5,
      ease: 'power2.out',
      onComplete: () => {
        gsap.to(scale, {
          value: phaseConfig.idle.targetScale,
          duration: duration * 0.5,
          ease: 'power2.inOut'
        })
      }
    })
  } else if (phase === 'idle') {
    scaleTween = gsap.to(scale, {
      value: config.targetScale,
      duration: 0.6,
      ease: 'power2.out'
    })
  } else if (phase === 'contract') {
    // 收缩：缓慢压缩
    scaleTween = gsap.to(scale, {
      value: config.targetScale,
      duration: duration * 0.85,
      ease: 'power2.inOut'
    })
  } else if (phase === 'hold') {
    // 保持：快速调整到保持状态
    scaleTween = gsap.to(scale, {
      value: config.targetScale,
      duration: 0.3,
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

// Animation loop
function animate() {
  time.value += 0.025
  animationFrame = requestAnimationFrame(animate)
}

// Watch phase changes
watch(() => trainingStore.phase, (newPhase) => {
  animateScaleToPhase(newPhase)
}, { immediate: true })

onMounted(() => {
  startBreathing()
  animate()
})

onUnmounted(() => {
  if (breathingTween) breathingTween.kill()
  if (scaleTween) scaleTween.kill()
  if (animationFrame) cancelAnimationFrame(animationFrame)
})

// Expose for parent component
defineExpose({
  combinedScale,
  currentConfig
})
</script>

<template>
  <div class="fluid-ball-container">
    <!-- Main fluid ball SVG -->
    <svg
      class="fluid-svg"
      viewBox="0 0 300 300"
      :style="{ transform: `scale(${combinedScale})` }"
    >
      <defs>
        <!-- 内部渐变 - 更通透的效果 -->
        <radialGradient id="ballFill" cx="30%" cy="30%" r="70%">
          <stop offset="0%" stop-color="rgba(255, 255, 255, 0.2)" />
          <stop offset="50%" :stop-color="currentConfig.fillColor" />
          <stop offset="100%" :stop-color="currentConfig.fillColor" />
        </radialGradient>

        <!-- 边框发光效果 -->
        <filter id="borderGlow" x="-50%" y="-50%" width="200%" height="200%">
          <feGaussianBlur in="SourceGraphic" stdDeviation="4" result="blur" />
          <feMerge>
            <feMergeNode in="blur" />
            <feMergeNode in="blur" />
            <feMergeNode in="SourceGraphic" />
          </feMerge>
        </filter>
      </defs>

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
        stroke-width="3"
        filter="url(#borderGlow)"
      />

      <!-- 高光 - 增强立体感 -->
      <ellipse
        cx="110"
        cy="100"
        rx="40"
        ry="25"
        fill="rgba(255, 255, 255, 0.1)"
        transform="rotate(-20, 110, 100)"
      />
    </svg>

    <!-- Content overlay -->
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
}

.fluid-svg path {
  transition: fill 0.5s ease, stroke 0.5s ease;
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
