import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'

const user = ref(null)
const isAuthenticated = computed(() => !!user.value)
const isLoading = ref(false)
const error = ref(null)

async function fetchCurrentUser() {
  isLoading.value = true
  error.value = null

  try {
    const res = await fetch('/api/users/profile/me', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' }
    })

    if (!res.ok) {
      throw new Error("Unauthorized or invalid session")
    }

    user.value = await res.json()
    return true
  } catch (err) {
    user.value = null
    error.value = err.message
    return false
  } finally {
    isLoading.value = false
  }
}

async function logout() {
  try {
    await fetch('/api/auth/signout', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' }
    })
  } catch (e) {
    console.warn("Logout failed:", e)
  }

  user.value = null
  error.value = null
}

export function useAuth() {
  return {
    user,
    isAuthenticated,
    isLoading,
    error,
    fetchCurrentUser,
    logout
  }
}
