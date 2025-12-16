<script setup lang="ts">
import { useMotion } from '@vueuse/motion'
import { ref, onMounted, watch } from 'vue'

interface Props {
  value: number | string
  label: string
  icon?: string
  color?: 'purple' | 'orange' | 'emerald' | 'pink'
  animated?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  color: 'purple',
  animated: true
})

const cardRef = ref<HTMLElement>()
const displayValue = ref(0)

const colorClasses = {
  purple: 'text-purple-400',
  orange: 'text-orange-400',
  emerald: 'text-emerald-400',
  pink: 'text-pink-400'
}

const glowClasses = {
  purple: 'group-hover:shadow-purple-500/20',
  orange: 'group-hover:shadow-orange-500/20',
  emerald: 'group-hover:shadow-emerald-500/20',
  pink: 'group-hover:shadow-pink-500/20'
}

onMounted(() => {
  if (cardRef.value) {
    useMotion(cardRef.value, {
      initial: { opacity: 0, y: 20, scale: 0.95 },
      enter: {
        opacity: 1,
        y: 0,
        scale: 1,
        transition: {
          type: 'spring',
          stiffness: 250,
          damping: 25
        }
      }
    })
  }

  if (props.animated && typeof props.value === 'number') {
    animateValue(0, props.value, 1000)
  } else {
    displayValue.value = typeof props.value === 'number' ? props.value : 0
  }
})

watch(() => props.value, (newVal, oldVal) => {
  if (props.animated && typeof newVal === 'number' && typeof oldVal === 'number') {
    animateValue(oldVal, newVal, 500)
  } else {
    displayValue.value = typeof newVal === 'number' ? newVal : 0
  }
})

function animateValue(start: number, end: number, duration: number) {
  const startTime = performance.now()

  function update(currentTime: number) {
    const elapsed = currentTime - startTime
    const progress = Math.min(elapsed / duration, 1)

    // Ease out cubic
    const easeOut = 1 - Math.pow(1 - progress, 3)
    displayValue.value = Math.round(start + (end - start) * easeOut)

    if (progress < 1) {
      requestAnimationFrame(update)
    }
  }

  requestAnimationFrame(update)
}
</script>

<template>
  <div
    ref="cardRef"
    class="group stat-card glass-card p-5 transition-all duration-300 hover:scale-[1.02]"
    :class="glowClasses[color]"
  >
    <div class="flex items-start justify-between">
      <div>
        <div
          class="text-3xl font-bold tracking-tight"
          :class="colorClasses[color]"
        >
          {{ typeof value === 'number' ? displayValue : value }}
        </div>
        <div class="text-sm text-white/50 mt-1">{{ label }}</div>
      </div>
      <div
        v-if="icon"
        class="text-2xl opacity-50 group-hover:opacity-80 transition-opacity"
      >
        {{ icon }}
      </div>
    </div>
  </div>
</template>
