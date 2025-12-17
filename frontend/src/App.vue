<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useUserStore } from '@/store/user'

const userStore = useUserStore()
const isReady = ref(false)

onMounted(async () => {
  if (userStore.token) {
    await userStore.fetchProfile()
  }
  // Small delay to ensure smooth initial render
  setTimeout(() => {
    isReady.value = true
  }, 100)
})
</script>

<template>
  <div
    class="app-wrapper dark"
    :class="{ 'app-ready': isReady }"
  >
    <RouterView v-slot="{ Component, route }">
      <Transition name="page" mode="out-in">
        <component :is="Component" :key="route.path" />
      </Transition>
    </RouterView>
  </div>
</template>

<style scoped>
.app-wrapper {
  opacity: 0;
  transition: opacity 0.3s ease;
}

.app-ready {
  opacity: 1;
}
</style>
