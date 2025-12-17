import { ref, watch } from 'vue'
import { defineStore } from 'pinia'

export type ThemeMode = 'light' | 'dark'

export const useThemeStore = defineStore('theme', () => {
  const mode = ref<ThemeMode>(
    (localStorage.getItem('theme') as ThemeMode) || 'dark'
  )

  // 初始化时应用主题
  function initTheme() {
    applyTheme(mode.value)
  }

  function applyTheme(theme: ThemeMode) {
    const root = document.documentElement
    if (theme === 'dark') {
      root.classList.add('dark')
      root.classList.remove('light')
    } else {
      root.classList.add('light')
      root.classList.remove('dark')
    }
  }

  function toggle() {
    mode.value = mode.value === 'dark' ? 'light' : 'dark'
  }

  function setMode(newMode: ThemeMode) {
    mode.value = newMode
  }

  // 监听变化并持久化
  watch(mode, (newMode) => {
    localStorage.setItem('theme', newMode)
    applyTheme(newMode)
  })

  return {
    mode,
    toggle,
    setMode,
    initTheme
  }
})
