import { ref } from 'vue'

const user = ref(null)
const isAuthenticated = ref(false)

async function fetchCurrentUser() {
  try {
    const res = await fetch('http://localhost:8080/users/profile/me', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' }
    })

    if (!res.ok) {
      user.value = null
      isAuthenticated.value = false
      return false
    }

    user.value = await res.json()
    console.log("user", user)
    console.log("userID", user.value.id)
    console.log("Nickname", user.value.nickname)


    isAuthenticated.value = true
    return true
  } catch (err) {
    user.value = null
    isAuthenticated.value = false
    return false
  }
}

export function useAuth() {
  return {
    user,
    isAuthenticated,
    fetchCurrentUser
  }
}
