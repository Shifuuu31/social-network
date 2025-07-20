import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useAuthStore = defineStore('auth', () => {
  // State
  const user = ref(null)
  const isLoading = ref(false)
  const error = ref(null)

  const API_BASE = '/api'

  // Getters
  const isAuthenticated = computed(() => !!user.value)
  const currentUserId = computed(() => user.value?.id)

  // Actions
  const signUp = async (userData) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/auth/signup`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(userData)
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || 'Registration failed')
      }

      // After successful signup, user can sign in
      return true
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const signIn = async (credentials) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/auth/signin`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(credentials),
        credentials: 'include' // Include cookies for session management
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || 'Sign in failed')
      }

      // Get user information after successful sign in
      await getCurrentUser()

      return true
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const signOut = async () => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/auth/signout`, {
        method: 'DELETE',
        credentials: 'include'
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || 'Sign out failed')
      }

      // Clear user data
      user.value = null

      return true
    } catch (err) {
      error.value = err.message
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const getCurrentUser = async () => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/user/profile/me`, {
        method: 'GET',
        credentials: 'include'
      })

      if (!response.ok) {
        if (response.status === 401) {
          // User not authenticated
          user.value = null
          return null
        }
        throw new Error('Failed to get user information')
      }

      const userData = await response.json()
      user.value = userData

      return userData
    } catch (err) {
      error.value = err.message
      user.value = null
      return null
    } finally {
      isLoading.value = false
    }
  }

  const checkAuthStatus = async () => {
    // Check if user is authenticated by trying to get current user
    return await getCurrentUser()
  }

  const clearError = () => {
    error.value = null
  }

  // Initialize auth store
  const initializeAuth = async () => {
    // Check if user is already authenticated
    await checkAuthStatus()
  }

  return {
    // State
    user,
    isLoading,
    error,

    // Getters
    isAuthenticated,
    currentUserId,

    // Actions
    signUp,
    signIn,
    signOut,
    getCurrentUser,
    checkAuthStatus,
    clearError,
    initializeAuth
  }
})
