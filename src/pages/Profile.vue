<template>
  <div class="profile-page">
    <!-- Top Profile Banner -->
    <div class="profile-banner">
      <!-- <img :src="user.profile_img || defaultAvatar" class="banner-pic" alt="Profile Picture" /> -->
      <div class="banner-info">
        <h2>{{ user.nickname || user.first_name}}</h2>
        <p>{{ user.about_me || 'No bio yet.' }}</p>
        <div class="profile-actions">
          <!-- If it's own profile -->
          <button v-if="isOwner" @click="toggleVisibility">
            {{ user.is_public ? '🔓 Public' : '🔒 Private' }}
          </button>

          <!-- If it's someone else -->
          <div v-else>
            <button v-if="followStatus === 'none'" @click="followUser">Follow</button>
            <button v-else-if="followStatus === 'pending'" disabled>Pending</button>
            <button v-else-if="followStatus === 'accepted'" @click="unfollowUser">Unfollow</button>
            <div v-else-if="followStatus === 'incoming'">
              <button @click="respondToRequest('accepted')">Accept</button>
              <button @click="respondToRequest('declined')">Decline</button>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content Sections -->
    <div class="profile-layout">
      <div class="profile-left">
        <h3>About</h3>
        <p>{{ user.nickname}}</p>
        <p>{{ user.first_name, user.last_name }}</p>
        <p>{{ user.about_me || 'No bio provided yet.' }}</p>
        <p class="dob" v-if="user.date_of_birth">📅 {{ user.date_of_birth }}</p>
        <!-- <p>{{ user.created_at}}</p> -->
      </div>

      <div class="profile-middle">
        <h3>Posts</h3>
        <!-- Future: Add post fetching logic here -->
        <p>Coming soon...</p>
      </div>

      <div class="profile-right">
        <h3>Chat (Coming soon)</h3>
        <!-- Future: Chat or message interface -->
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted, reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const defaultAvatar = '/default-avatar.png'
const user = reactive({})
const followStatus = ref('none') // 'none', 'pending', 'accepted', 'incoming'
const isOwner = ref(false)

// Attempt to retrieve the ID from router state (hidden navigation data)
let targetId = router.options.history.state?.targetId
if (!targetId) {
  // fallback to own profile
  targetId = parseInt(localStorage.getItem('user_id'))
}

onMounted(async () => {
  const currentId = parseInt(localStorage.getItem('user_id'))
  isOwner.value = currentId === targetId

  await fetchUserInfo()
  // await fetchFollowStatus()
})

async function fetchUserInfo() {
  try {
    const res = await fetch('http://localhost:8080/users/profile/info', {
    method: 'POST',
    headers: { 'Content-Type': 'application/json' },
    credentials: "include",
    body: JSON.stringify({ id: 1 }),

      // credentials: "include", // Required if you want to send cookies
    })
    if (!res.ok) {
      console.log("status", res.status)
      alert('Failed to fetch profile info')
      return 
    } 

    const data = await res.json()
    console.log("res",res)
    console.log("data",data)
    Object.assign(user, data)
  }catch(err){
    // console.log("st", res.status)
    console.log("err",err)
    return 
  }
}

// async function fetchFollowStatus() {
//   const res = await fetch('http://localhost:8080/users/follow/status', {
//     method: 'POST',
//     headers: { 'Content-Type': 'application/json' },
//     body: JSON.stringify({ target_id: 1 })
//   })
//   if (!res.ok) return
//   const data = await res.json()
//   followStatus.value = data.status
// }

// async function followUser() {
//   const res = await fetch('http://localhost:8080/users/follow/follow-unfollow', {
//     method: 'POST',
//     headers: { 'Content-Type': 'application/json' },
//     body: JSON.stringify({ target_id: targetId, action: 'follow' })
//   })
//   if (res.ok) followStatus.value = 'pending'
// }

// async function unfollowUser() {
//   const res = await fetch('http://localhost:8080/users/follow/follow-unfollow', {
//     method: 'POST',
//     headers: { 'Content-Type': 'application/json' },
//     body: JSON.stringify({ target_id: targetId, action: 'unfollow' })
//   })
//   if (res.ok) followStatus.value = 'none'
// }

// async function respondToRequest(action) {
//   const res = await fetch('http://localhost:8080/users/follow/accept-decline', {
//     method: 'POST',
//     headers: { 'Content-Type': 'application/json' },
//     body: JSON.stringify({ target_id: targetId, action })
//   })
//   if (res.ok) followStatus.value = action === 'accepted' ? 'accepted' : 'none'
// }

async function toggleVisibility() {
  try {
    const res = await fetch('http://localhost:8080/users/profile/visibility', {
    method: 'POST',
    credentials: "include",
    headers: { 'Content-Type': 'application/json' }
    })
    console.log(2)
    if (!res.ok) return alert('Failed to toggle visibility')
    const updated = await res.json()
    console.log(updated.is_public)
    user.is_public = updated.is_public
    console.log(updated.is_public)
  }catch(err){
    console.log(err)
    return
  } 
}
</script>

<style scoped>
.profile-page {
  max-width: 1000px;
  margin: 2rem auto;
  padding: 2rem;
  background: white;
  border-radius: 8px;
  font-family: sans-serif;
}

.profile-banner {
  display: flex;
  gap: 2rem;
  align-items: center;
  padding-bottom: 2rem;
  border-bottom: 1px solid #ccc;
}

.banner-pic {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #999;
}

.banner-info h2 {
  color: black;
  font-size: 2rem;
  margin: 0 0 0.5rem;
}

.banner-info p {
  margin: 0.25rem 0;
  color: #555;
}

.profile-actions {
  margin-top: 1rem;
  display: flex;
  gap: 1rem;
}

.profile-layout {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  gap: 1.5rem;
  margin-top: 2rem;
}

.profile-left,
.profile-middle,
.profile-right {
  color: black;
  padding: 1rem;
  border: 1px solid #eee;
  border-radius: 8px;
  background: #f9f9f9;
}

button {
  padding: 0.5rem 1rem;
  background: #1e293b;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
}

button:hover {
  background: #334155;
}

button:disabled {
  background: #888;
  cursor: default;
}
</style>
