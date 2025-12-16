<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()

onMounted(() => {
  if (userStore.token) {
    userStore.fetchProfile()
  }
})
</script>

<template>
  <RouterView />
</template>
