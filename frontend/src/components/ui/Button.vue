<script setup lang="ts">
import { computed } from 'vue'
import { Loader2 } from 'lucide-vue-next'

interface Props {
  variant?: 'primary' | 'secondary' | 'ghost' | 'danger'
  size?: 'sm' | 'md' | 'lg'
  loading?: boolean
  disabled?: boolean
  fullWidth?: boolean
  icon?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  variant: 'primary',
  size: 'md',
  loading: false,
  disabled: false,
  fullWidth: false,
  icon: false
})

const classes = computed(() => {
  const base = 'inline-flex items-center justify-center font-semibold transition-all duration-300 ease-out focus:outline-none focus-visible:ring-2 focus-visible:ring-purple-500 focus-visible:ring-offset-2 focus-visible:ring-offset-slate-900 disabled:opacity-50 disabled:cursor-not-allowed'

  const variants = {
    primary: 'btn-gradient text-white rounded-full shadow-lg shadow-purple-500/25 hover:shadow-purple-500/40',
    secondary: 'bg-white/10 text-white rounded-full hover:bg-white/20 border border-white/10 hover:border-white/20',
    ghost: 'text-white/70 hover:text-white hover:bg-white/10 rounded-lg',
    danger: 'bg-red-500/20 text-red-400 rounded-full hover:bg-red-500/30 border border-red-500/30'
  }

  const sizes = {
    sm: props.icon ? 'w-8 h-8' : 'px-4 py-2 text-sm gap-1.5',
    md: props.icon ? 'w-10 h-10' : 'px-6 py-2.5 text-base gap-2',
    lg: props.icon ? 'w-12 h-12' : 'px-8 py-3 text-lg gap-2.5'
  }

  return [
    base,
    variants[props.variant],
    sizes[props.size],
    props.fullWidth ? 'w-full' : ''
  ].join(' ')
})

const isDisabled = computed(() => props.disabled || props.loading)
</script>

<template>
  <button
    :class="classes"
    :disabled="isDisabled"
  >
    <Transition name="fade" mode="out-in">
      <Loader2 v-if="loading" class="w-5 h-5 animate-spin" />
      <span v-else class="inline-flex items-center gap-2">
        <slot />
      </span>
    </Transition>
  </button>
</template>
