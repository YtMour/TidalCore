import { defineStore } from 'pinia'
import { ref, computed, watch } from 'vue'
import { login as apiLogin, register as apiRegister, getProfile } from '@/api/auth'
import type { LoginRequest, RegisterRequest, UserInfo } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<UserInfo | null>(null)
  const isLoading = ref(false)
  const lastFetchTime = ref(0)

  const isLoggedIn = computed(() => !!token.value)

  // Watch token changes and sync to localStorage
  watch(token, (newToken) => {
    if (newToken) {
      localStorage.setItem('token', newToken)
    } else {
      localStorage.removeItem('token')
    }
  })

  async function login(data: LoginRequest) {
    isLoading.value = true
    try {
      const res = await apiLogin(data)
      token.value = res.token
      user.value = res.user
      lastFetchTime.value = Date.now()
      return res
    } finally {
      isLoading.value = false
    }
  }

  async function register(data: RegisterRequest) {
    isLoading.value = true
    try {
      const res = await apiRegister(data)
      token.value = res.token
      user.value = res.user
      lastFetchTime.value = Date.now()
      return res
    } finally {
      isLoading.value = false
    }
  }

  async function fetchProfile(force = false) {
    if (!token.value) return

    // Avoid fetching too frequently (cache for 30 seconds)
    const now = Date.now()
    if (!force && lastFetchTime.value && now - lastFetchTime.value < 30000) {
      return user.value
    }

    isLoading.value = true
    try {
      const res = await getProfile()
      user.value = res
      lastFetchTime.value = now
      return res
    } catch {
      logout()
      return null
    } finally {
      isLoading.value = false
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    lastFetchTime.value = 0
  }

  function updateStreak(streak: number, maxStreak: number, totalCheckin: number) {
    if (user.value) {
      user.value = {
        ...user.value,
        streak,
        max_streak: maxStreak,
        total_checkin: totalCheckin
      }
      lastFetchTime.value = Date.now()
    }
  }

  // Refresh user data
  async function refreshUser() {
    return fetchProfile(true)
  }

  return {
    token,
    user,
    isLoggedIn,
    isLoading,
    login,
    register,
    fetchProfile,
    logout,
    updateStreak,
    refreshUser
  }
})
