import { reactive, ref, computed, watch } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'

export function useProfileView() {
  const router = useRouter()
  const { user: currentUser, isAuthenticated, fetchCurrentUser } = useAuth()  

  const defaultAvatar = '/images/default-avatar.png' //fake
  const profileUser = reactive({})
  const followStatus = ref('none')
  const isRequestToMe = ref(false)
  const isOwner = ref(false)
  const activeTab = ref('posts')
  const followersList = ref([])
  const followingList = ref([])
  let targetId = null

  const canViewPrivateProfile = computed(() => {
    return isOwner.value || profileUser.is_public || followStatus.value === 'accepted'
  })

  async function initProfile() {
    if (!isAuthenticated.value) {
      await fetchCurrentUser()
    }
    // const routeID = Number(route.params.id)
    const routerID = Number(router.currentRoute.value.params.id)
    // targetId = router.options.history.state?.targetId || currentUser.value?.id // un peux compliquer
    // targetId = 1 // for testing
    targetId = routerID || currentUser.value?.id
      
    isOwner.value = currentUser.value.id === targetId

    const profileExist = await fetchProfile()
    if (!profileExist){
      targetId = currentUser.value.id
      isOwner.value = true
      await fetchProfile()
    }
  }

  async function fetchProfile() {
    try {
      const res = await fetch('http://localhost:8080/users/profile/info', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: targetId }),
      })

      if (!res.ok) {
        console.error('Error fetching profile:', res.status)
        return false
      }

      const { user: u, follow_status, is_request_to_me } = await res.json()
      Object.keys(u).forEach(key => {
        profileUser[key] = u[key]
      })
      console.log("user", u)
      console.log("profileUser", profileUser)

      followStatus.value = follow_status
      isRequestToMe.value = is_request_to_me
      return true
    }catch(err){
      console.log(err)
      return false
    }
  }

  async function toggleFollow(action) {
    try {
      const res = await fetch('http://localhost:8080/users/follow/follow-unfollow', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ target_id: targetId, action }),
      })

      if (res.ok) {
        followStatus.value = action === 'follow' ? 'pending' : 'none'
      }
    }catch(err){
      console.log(err)
    }  
  }

  async function toggleVisibility() {
    try {
      const res = await fetch('http://localhost:8080/users/profile/visibility', {
        method: 'POST',
        credentials: 'include',
      })

      if (!res.ok) {
        return alert('Failed to toggle visibility')
      }

      const updated = await res.json()
      profileUser.is_public = updated.is_public
      console.log(profileUser.is_public)
    }catch(err){
      console.log(err)
    }
  }

  async function respondToRequest(action) {
    try {
      const res = await fetch('http://localhost:8080/users/follow/accept-decline', {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          target_id: profileUser.id,
          action: action,
        }),
      })

      if (!res.ok) {
        console.error('Failed to respond to follow request:', res.status, res.text)
        return alert('Failed to respond to follow request')
      }

      isRequestToMe.value = false // hide buttons
    } catch (err) {
      console.error('Failed to respond to request:', err)
    }
  }

  async function fetchConnections(type) {
    try {
      const endpoint = type === 'followers' ? 'http://localhost:8080/users/profile/followers' : 'http://localhost:8080/users/profile/following'
        
      const res = await fetch(endpoint, {
        method: 'POST',
        credentials: 'include',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: targetId }),
      })

      if (!res.ok) return console.error(`Failed to load ${type}`)

      const data = await res.json()
      if (type === 'followers') followersList.value = data || []
      if (type === 'following') followingList.value = data || []
    }catch(err){
      console.log(err)
    }
  }

  watch(activeTab, (newTab) => {
    if (canViewPrivateProfile.value && (newTab === 'followers' || newTab === 'following')) {
      fetchConnections(newTab)
    }
  })

  return {
      profileUser,
      followStatus,
      isRequestToMe,
      isOwner,
      activeTab,
      followersList,
      followingList,
      canViewPrivateProfile,
      initProfile,
      respondToRequest,
      toggleFollow,
      toggleVisibility,
  }
}