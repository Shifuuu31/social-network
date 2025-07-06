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
    const res = await fetch('http://localhost:8080/users/profile/me', {
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
  const router = useRouter()

  try {
    await fetch('http://localhost:8080/auth/signout', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' }
    })
  } catch (e) {
    console.warn("Logout failed:", e)
  }

  user.value = null
  error.value = null
  router.push('/signin')
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
