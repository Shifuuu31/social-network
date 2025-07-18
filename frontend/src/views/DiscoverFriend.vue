<template>
  <div class="discover-friend-container">
    <h2>Discover Friends</h2>
    <div v-if="loading" class="loading-spinner">Loading users...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else>
      <div v-if="users.length" class="user-list">
        <div v-for="user in users" :key="user.id" class="user-item">
          <router-link :to="`/profile/${user.id}`" class="user-link">
            <img :src="getAvatarUrl(user.avatar_url)" class="avatar" :alt="user.nickname || user.first_name" />
            <div class="user-info">
              <h4>{{ user.nickname || user.first_name }}</h4>
              <span class="user-handle">@{{ user.nickname }}</span>
            </div>
          </router-link>
        </div>
      </div>
      <div v-else class="no-users">No users found.</div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const users = ref([])
const loading = ref(true)
const error = ref('')

function getAvatarUrl(avatarUrl) {
  if (avatarUrl && avatarUrl.startsWith('/images/')) {
    return `http://localhost:8080${avatarUrl}`
  }
  if (avatarUrl && avatarUrl.startsWith('uploads/')) {
    const filename = avatarUrl.replace('uploads/', '')
    return `http://localhost:8080/images/${filename}`
  }
  return '/images/default-avatar.png'
}

async function fetchUsers() {
  loading.value = true
  error.value = ''
  try {
    // Adjust endpoint as needed for your backend
    const response = await fetch('http://localhost:8080/users/all', {
      credentials: 'include',
      headers: { 'Accept': 'application/json' }
    })
    if (!response.ok) {
      throw new Error('Failed to fetch users')
    }
    const data = await response.json()
    users.value = Array.isArray(data) ? data : []
  } catch (err) {
    error.value = err.message || 'Failed to fetch users.'
    users.value = []
  } finally {
    loading.value = false
  }
}

onMounted(fetchUsers)
</script>

<style scoped>
.discover-friend-container {
  max-width: 600px;
  margin: 2rem auto;
  padding: 2rem;
  background: #fff;
  border-radius: 12px;
  box-shadow: 0 2px 10px rgba(0,0,0,0.07);
}
.user-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}
.user-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  border: 1px solid #dbdbdb;
  border-radius: 8px;
  background: #fafafa;
}
.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #f0f0f0;
}
.user-info h4 {
  margin: 0;
  font-size: 1.1rem;
  font-weight: 600;
}
.user-handle {
  font-size: 0.9rem;
  color: #8e8e8e;
}
.loading-spinner {
  text-align: center;
  color: #0095f6;
  font-weight: 600;
  margin: 2rem 0;
}
.error {
  color: #dc2743;
  text-align: center;
  margin: 2rem 0;
}
.no-users {
  text-align: center;
  color: #8e8e8e;
  margin: 2rem 0;
}
</style> 