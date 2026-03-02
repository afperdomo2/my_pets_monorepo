import { ref, computed } from 'vue'
import { defineStore } from 'pinia'
import type { User } from '@/types/user'
import { authService } from '@/services/authService'

export const useAuthStore = defineStore('auth', () => {
  const user = ref<User | null>(null)
  const loading = ref(false)
  const error = ref<string | null>(null)

  const isAuthenticated = computed(() => user.value !== null)

  async function login(email: string, password: string): Promise<void> {
    loading.value = true
    error.value = null
    try {
      const res = await authService.login({ email, password })
      user.value = res.data
    } catch (e) {
      error.value = e instanceof Error ? e.message : 'Login failed'
      throw e
    } finally {
      loading.value = false
    }
  }

  async function logout(): Promise<void> {
    try {
      await authService.logout()
    } catch {
      // Ignore network errors on logout — clear state anyway
    } finally {
      user.value = null
    }
  }

  async function fetchMe(): Promise<void> {
    const res = await authService.me()
    user.value = res.data
  }

  /**
   * Called once on app startup to restore session from the access_token cookie.
   * If the cookie is expired, attempts a silent refresh before giving up.
   */
  async function initSession(): Promise<void> {
    loading.value = true
    try {
      await fetchMe()
    } catch {
      // Access token may be expired — try refresh
      try {
        await authService.refresh()
        await fetchMe()
      } catch {
        user.value = null
      }
    } finally {
      loading.value = false
    }
  }

  return {
    user,
    loading,
    error,
    isAuthenticated,
    login,
    logout,
    fetchMe,
    initSession,
  }
})
