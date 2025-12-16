<script setup lang="ts">
interface Props {
  hover?: boolean
  padding?: 'none' | 'sm' | 'md' | 'lg'
}

withDefaults(defineProps<Props>(), {
  hover: false,
  padding: 'lg'
})
</script>

<template>
  <div
    class="glass-card"
    :class="[
      hover ? 'glass-card-hover cursor-pointer' : '',
      {
        'p-0': padding === 'none',
        'p-4': padding === 'sm',
        'p-6': padding === 'md',
        'p-8': padding === 'lg'
      }
    ]"
  >
    <slot />
  </div>
</template>
