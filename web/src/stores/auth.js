import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { authAPI } from '@/api'
import { loadUserFromStorage, saveUserToStorage } from '@/utils/storage'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(loadUserFromStorage())
  const loading = ref(false)
  const error = ref(null)
  const initialized = ref(false)

  // Getters
  const isAuthenticated = computed(() => initialized.value && user.value !== null)
  const currentUser = computed(() => user.value)
  const isInitializing = computed(() => !initialized.value)

  // Actions
  async function login(email, password, rememberMe = false) {
    loading.value = true
    error.value = null

    try {
      const response = await authAPI.login({ email, password, rememberMe })
      const userData = response.data?.user

      if (!userData || !response.data?.successMessage) {
        throw new Error('Invalid server response')
      }

      user.value = userData
      saveUserToStorage(userData)
      initialized.value = true
      return { success: true }
    } catch (err) {
      error.value = err.response?.data?.errorMessage || err.message || 'Login failed'
      return { success: false, error: error.value }
    } finally {
      loading.value = false
    }
  }

  async function logout() {
    // Clear user state immediately to prevent race conditions
    user.value = null
    saveUserToStorage(null)
    error.value = null
    loading.value = true

    try {
      await authAPI.logout()
    } catch (err) {
      console.error('Logout API error:', err)
    } finally {
      loading.value = false
    }
  }

  async function checkAuth() {
    const isFirstCheck = !initialized.value
    if (isFirstCheck) loading.value = true

    try {
      const response = await authAPI.getProfile()
      const userData = response.data?.user

      user.value = userData || null
      saveUserToStorage(userData || null)
    } catch (err) {
      if (err.response?.status === 401) {
        console.debug('Session invalid or expired')
      } else {
        console.error('Auth check failed:', err.message)
      }
      user.value = null
      saveUserToStorage(null)
    } finally {
      initialized.value = true
      if (isFirstCheck) loading.value = false
    }
  }

  async function updateProfile(updates) {
    if (!user.value) {
      return { success: false, error: 'Not authenticated' }
    }

    const previousUser = { ...user.value }
    const updatedUser = { ...user.value, ...updates }

    user.value = updatedUser
    saveUserToStorage(updatedUser)

    try {
      // TODO: Call backend API when profile update endpoint is implemented
      // await authAPI.updateProfile(updates)
      return { success: true }
    } catch (err) {
      user.value = previousUser
      saveUserToStorage(previousUser)
      return {
        success: false,
        error: err.response?.data?.errorMessage || 'Profile update failed'
      }
    }
  }

  // Initialize auth state on store creation
  checkAuth()

  return {
    // State
    user,
    loading,
    error,
    initialized,
    // Getters
    isAuthenticated,
    currentUser,
    isInitializing,
    // Actions
    login,
    logout,
    checkAuth,
    updateProfile
  }
})
