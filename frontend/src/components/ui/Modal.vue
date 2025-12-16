<script setup lang="ts">
import { watch } from 'vue'
import { X } from 'lucide-vue-next'
import { useScrollLock } from '@vueuse/core'

interface Props {
  modelValue: boolean
  title?: string
  size?: 'sm' | 'md' | 'lg'
  closable?: boolean
}

const props = withDefaults(defineProps<Props>(), {
  size: 'md',
  closable: true
})

const emit = defineEmits<{
  'update:modelValue': [value: boolean]
}>()

const scrollLock = useScrollLock(document.body)

watch(() => props.modelValue, (value) => {
  scrollLock.value = value
})

function close() {
  if (props.closable) {
    emit('update:modelValue', false)
  }
}

function handleBackdropClick(event: MouseEvent) {
  if (event.target === event.currentTarget) {
    close()
  }
}
</script>

<template>
  <Teleport to="body">
    <Transition name="modal">
      <div
        v-if="modelValue"
        class="fixed inset-0 z-50 flex items-center justify-center p-4"
        @click="handleBackdropClick"
      >
        <!-- Backdrop -->
        <div class="absolute inset-0 bg-black/60 backdrop-blur-sm"></div>

        <!-- Modal Content -->
        <div
          class="modal-content relative glass-card w-full transform transition-all duration-300"
          :class="{
            'max-w-sm': size === 'sm',
            'max-w-md': size === 'md',
            'max-w-lg': size === 'lg'
          }"
        >
          <!-- Header -->
          <div v-if="title || closable" class="flex items-center justify-between p-6 pb-0">
            <h3 v-if="title" class="text-xl font-semibold text-white">
              {{ title }}
            </h3>
            <button
              v-if="closable"
              @click="close"
              class="icon-btn ml-auto -mr-2 -mt-2"
            >
              <X class="w-5 h-5" />
            </button>
          </div>

          <!-- Body -->
          <div class="p-6">
            <slot />
          </div>

          <!-- Footer -->
          <div v-if="$slots.footer" class="p-6 pt-0">
            <slot name="footer" />
          </div>
        </div>
      </div>
    </Transition>
  </Teleport>
</template>
