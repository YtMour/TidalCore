import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { login as apiLogin, register as apiRegister, getProfile } from '@/api/auth'
import type { LoginRequest, RegisterRequest } from '@/api/auth'

export const useUserStore = defineStore('user', () => {
  const token = ref(localStorage.getItem('token') || '')
  const user = ref<{
    id: number
    username: string
    streak: number
    max_streak: number
    total_checkin: number
  } | null>(null)

  const isLoggedIn = computed(() => !!token.value)

  async function login(data: LoginRequest) {
    const res = await apiLogin(data)
    token.value = res.token
    user.value = res.user
    localStorage.setItem('token', res.token)
    return res
  }

  async function register(data: RegisterRequest) {
    const res = await apiRegister(data)
    token.value = res.token
    user.value = res.user
    localStorage.setItem('token', res.token)
    return res
  }

  async function fetchProfile() {
    if (!token.value) return
    try {
      const res = await getProfile()
      user.value = res as typeof user.value
    } catch {
      logout()
    }
  }

  function logout() {
    token.value = ''
    user.value = null
    localStorage.removeItem('token')
  }

  function updateStreak(streak: number, maxStreak: number, totalCheckin: number) {
    if (user.value) {
      user.value.streak = streak
      user.value.max_streak = maxStreak
      user.value.total_checkin = totalCheckin
    }
  }

  return {
    token,
    user,
    isLoggedIn,
    login,
    register,
    fetchProfile,
    logout,
    updateStreak
  }
})
