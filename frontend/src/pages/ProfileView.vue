<template>
  <div class="profile-container">
    <!-- Banner Placeholder -->
    <div class="profile-banner">
      <div class="banner-image"></div>
    </div>

    <!-- Profile Info Card -->
    <div class="profile-header">
      <img class="avatar" :src="profileUser.avatar || defaultAvatar" alt="Profile Picture" />
      <div class="profile-info">
        <h2>{{ profileUser.nickname || profileUser.first_name }}</h2>
        <!-- <p class="role"></p> -->
        <div class="profile-buttons">
          <button v-if="isOwner" @click="toggleVisibility">
            {{ profileUser.is_public ? '🔓 Public' : '🔒 Private' }}
          </button>
          <template v-else>
            <button
              v-if="followStatus === 'none'"
              @click="toggleFollow('follow')"
            >Follow</button>
            <button v-else-if="followStatus === 'pending'" disabled>Pending</button>
            <button
              v-else-if="followStatus === 'accepted'"
              @click="toggleFollow('unfollow')"
            >Unfollow</button>
          </template>
        </div>
      </div>
    </div>

    <!-- Main Layout -->
    <div class="profile-main" v-if="canViewPrivateProfile">
      <!-- Left Column -->
      <div class="profile-left">
        <h3>About</h3>
        <ul class="about-list">
          <li><strong>Gender:</strong> {{ profileUser.gender || 'N/A' }}</li>
          <li><strong>DOB:</strong> {{ profileUser.date_of_birth || 'N/A' }}</li>
          <li><strong>Location:</strong> {{ profileUser.location || 'N/A' }}</li>
          <li><strong>Email:</strong> {{ profileUser.email || 'N/A' }}</li>
          <li><strong>Phone:</strong> {{ profileUser.phone || 'N/A' }}</li>
        </ul>
      </div>

      <!-- Center Column -->
      <div class="profile-center">
        <div class="tabs">
          <span :class="{ active: activeTab === 'posts' }" @click="activeTab = 'posts'">Posts</span>
          <span :class="{ active: activeTab === 'followers' }" @click="activeTab = 'followers'">Followers</span>
          <span :class="{ active: activeTab === 'following' }" @click="activeTab = 'following'">Following</span>
        </div>

        <div class="tab-content">
          <div v-if="activeTab === 'posts'">
            <p>Coming soon: posts will appear here.</p>
          </div>
          <div v-if="activeTab === 'followers'">
            <p>List of followers from backend...</p>
          </div>
          <div v-if="activeTab === 'following'">
            <p>List of following from backend...</p>
          </div>
        </div>
      </div>

      <!-- Right Column -->
      <div class="profile-right">
        <h3>Chat</h3>
        <p>Coming soon…</p>
      </div>
    </div>
    <div v-else class="locked-profile">
    <div class="locked-card">
      <h3>🔒 This profile is private</h3>
      <p>You must follow this user to view their posts and profile details.</p>
    </div>
    </div>
  </div>
</template>


<script setup>
import { onMounted, reactive, ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuth } from '@/composables/useAuth'
// import Profile from './Profile.vue'

const router = useRouter()
const { user: currentUser, isAuthenticated, fetchCurrentUser } = useAuth()

const profileUser = reactive({})
const followStatus = ref('none')
const isOwner = ref(false)
let targetId = null



onMounted(async ()=> {
if (!isAuthenticated.value) {
  await fetchCurrentUser()
}
// Determine which profile we’re looking at:
targetId = router.options.history.state?.targetId || currentUser.value?.id
// targetId = 1
console.log("CurrentUserID: ", currentUser.value.id)
console.log("TargetID: ", targetId)

isOwner.value = currentUser.value.id === targetId
console.log("OWn",isOwner.value)
console.log("OWn",profileUser.is_public)
console.log("OWn",followStatus.value)

  await fetchProfile()
})

const canViewPrivateProfile = computed(() =>{
  return isOwner.value || profileUser.is_public || followStatus.value === 'accepted'
})

async function fetchProfile() {
  // Fetch profile info + follow status together
  const res = await fetch(
    'http://localhost:8080/users/profile/info',
    {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ id: targetId }),
    }
  )

  if (!res.ok) {
    console.error('Error fetching profile:', res.status)
    return
  }
  const { user: u, follow_status } = await res.json()
  // Object.assign(profileUser, u)
  Object.keys(u).forEach(key => {
  profileUser[key] = u[key]
  })

  console.log("PU",profileUser)
  console.log("UU",u)
  console.log("U",profileUser.date_of_birth)
  

  followStatus.value = follow_status
}

// Unified follow/unfollow:
async function toggleFollow(action) {
  const res = await fetch(
    'http://localhost:8080/users/follow/follow-unfollow',
    {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify({ target_id: targetId, action }),
    }
  )
  if (res.ok) {
    followStatus.value = action === 'follow' ? 'pending' : 'none'
  }
}

// Toggle your own visibility:
async function toggleVisibility() {
  const res = await fetch(
    'http://localhost:8080/users/profile/visibility',
    {
      method: 'POST',
      credentials: 'include',
    }
  )
  if (!res.ok) {
    return alert('Failed to toggle visibility')
  }
  const updated = await res.json()
  profileUser.is_public = updated.is_public
}
</script>
<style scoped>
.profile-container {
  background: #fff;
  color: #1e1e1e;
  max-width: 1200px;
  margin: auto;
  padding-bottom: 2rem;
  font-family: 'Segoe UI', sans-serif;
}

/* Top banner */
.profile-banner .banner-image {
  height: 200px;
  background: linear-gradient(to right, #8a2be2, #6a0dad);
  border-radius: 8px 8px 0 0;
}

/* Profile Header */
.profile-header {
  display: flex;
  align-items: center;
  padding: 1rem 2rem;
  background: #f9f9f9;
  border-radius: 0 0 8px 8px;
  margin-bottom: 2rem;
  gap: 1.5rem;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 100px;
  border: 4px solid white;
  object-fit: cover;
  background: #eee;
}

.profile-info {
  flex: 1;
}

.profile-info h2 {
  font-size: 1.8rem;
  margin-bottom: 0.2rem;
}

.role {
  color: #6a0dad;
  font-weight: bold;
  margin-bottom: 0.5rem;
}

.profile-buttons {
  display: flex;
  gap: 1rem;
  margin-top: 0.5rem;
}
.locked-profile {
  display: flex;
  justify-content: center;
  padding: 2rem;
}

.locked-card {
  background: #f5f5f5;
  border: 1px solid #ddd;
  padding: 2rem;
  border-radius: 12px;
  text-align: center;
  max-width: 500px;
  width: 100%;
  box-shadow: 0 2px 8px rgba(0,0,0,0.05);
}

.locked-card h3 {
  color: #6a0dad;
  margin-bottom: 0.5rem;
}

.locked-card p {
  color: #444;
}

button {
  background: #6a0dad;
  color: white;
  padding: 0.5rem 1.2rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
  transition: background 0.3s;
}

button:hover {
  background: #7d20c0;
}

button:disabled {
  background: #bbb;
  cursor: default;
}

/* Main Layout */
.profile-main {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  gap: 2rem;
  padding: 0 2rem;
}

/* Left Column */
.profile-left {
  background: #f5f5f5;
  border-radius: 8px;
  padding: 1rem;
}

.about-list {
  list-style: none;
  padding: 0;
}

.about-list li {
  margin-bottom: 0.5rem;
}

/* Center Column */
.profile-center {
  background: #fefefe;
  border-radius: 8px;
  padding: 1rem;
}

.tabs {
  display: flex;
  gap: 2rem;
  margin-bottom: 1rem;
  border-bottom: 2px solid #ddd;
}

.tabs span {
  padding: 0.5rem;
  cursor: pointer;
  font-weight: bold;
  color: #444;
}

.tabs span.active {
  color: #6a0dad;
  border-bottom: 3px solid #6a0dad;
}

/* Right Column */
.profile-right {
  background: #f5f5f5;
  border-radius: 8px;
  padding: 1rem;
}
</style>
