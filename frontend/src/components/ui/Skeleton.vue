<script setup lang="ts">
interface Props {
  width?: string
  height?: string
  rounded?: 'sm' | 'md' | 'lg' | 'full'
}

withDefaults(defineProps<Props>(), {
  width: '100%',
  height: '1rem',
  rounded: 'md'
})
</script>

<template>
  <div
    class="skeleton"
    :class="{
      'rounded-sm': rounded === 'sm',
      'rounded-md': rounded === 'md',
      'rounded-lg': rounded === 'lg',
      'rounded-full': rounded === 'full'
    }"
    :style="{ width, height }"
  ></div>
</template>
