<script setup lang="ts">
import { computed, ref } from 'vue'
import { Eye, EyeOff } from 'lucide-vue-next'

interface Props {
  modelValue: string | number
  type?: 'text' | 'password' | 'email' | 'number'
  label?: string
  placeholder?: string
  error?: string
  disabled?: boolean
  autocomplete?: string
  min?: number
  max?: number
}

const props = withDefaults(defineProps<Props>(), {
  type: 'text',
  disabled: false
})

const emit = defineEmits<{
  'update:modelValue': [value: string | number]
}>()

const showPassword = ref(false)
const isFocused = ref(false)

const inputType = computed(() => {
  if (props.type === 'password') {
    return showPassword.value ? 'text' : 'password'
  }
  return props.type
})

function handleInput(event: Event) {
  const target = event.target as HTMLInputElement
  const value = props.type === 'number' ? Number(target.value) : target.value
  emit('update:modelValue', value)
}
</script>

<template>
  <div class="space-y-1.5">
    <label v-if="label" class="block text-sm font-medium text-white/70">
      {{ label }}
    </label>

    <div class="relative">
      <input
        :type="inputType"
        :value="modelValue"
        :placeholder="placeholder"
        :disabled="disabled"
        :autocomplete="autocomplete"
        :min="min"
        :max="max"
        @input="handleInput"
        @focus="isFocused = true"
        @blur="isFocused = false"
        class="input-field pr-10"
        :class="[
          error ? 'border-red-500/50 focus:border-red-500' : '',
          disabled ? 'opacity-50 cursor-not-allowed' : ''
        ]"
      />

      <button
        v-if="type === 'password'"
        type="button"
        @click="showPassword = !showPassword"
        class="absolute right-3 top-1/2 -translate-y-1/2 text-white/40 hover:text-white/70 transition-colors"
      >
        <Eye v-if="!showPassword" class="w-5 h-5" />
        <EyeOff v-else class="w-5 h-5" />
      </button>
    </div>

    <Transition name="slide-up">
      <p v-if="error" class="text-sm text-red-400 flex items-center gap-1">
        <span class="inline-block w-1 h-1 rounded-full bg-red-400"></span>
        {{ error }}
      </p>
    </Transition>
  </div>
</template>
