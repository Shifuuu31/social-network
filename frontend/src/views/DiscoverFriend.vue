<template>
  <div class="discover-friend-container">
    <h2>Discover Friends</h2>
    
    <!-- Search Bar -->
    <div class="search-section">
      <input 
        v-model="searchQuery" 
        @input="handleSearch($event.target.value)"
        placeholder="Search users by name, nickname, or email..."
        class="simple-search"
      />
    </div>

    <div v-if="loading" class="loading-spinner">Loading users...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else>
      <div v-if="filteredUsers.length" class="user-list">
        <div class="results-header">
          <h3>{{ filteredUsers.length }} user{{ filteredUsers.length !== 1 ? 's' : '' }} found</h3>
          <button @click="resetFilters" class="reset-btn">Show All</button>
        </div>
        <div v-for="user in filteredUsers" :key="user.id" class="user-item">
          <router-link :to="`/profile/${user.id}`" class="user-link">
            <img :src="getAvatarUrl(user.avatar_url)" class="avatar" :alt="user.nickname || user.first_name" />
            <div class="user-info">
              <h4>{{ user.nickname || user.first_name }}</h4>
              <span class="user-handle">@{{ user.nickname }}</span>
              <span v-if="user.email" class="user-email">{{ user.email }}</span>
            </div>
          </router-link>
        </div>
      </div>
      <div v-else class="no-users">
        <div class="no-users-icon">👥</div>
        <h3>No users found</h3>
        <p v-if="searchQuery">Try adjusting your search terms</p>
        <p v-else>No users are currently available</p>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted } from 'vue'

const users = ref([])
const loading = ref(true)
const error = ref('')
const searchQuery = ref('')
const filteredUsers = ref([])

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
    filteredUsers.value = [...users.value] // Initialize filtered users with all users
  } catch (err) {
    error.value = err.message || 'Failed to fetch users.'
    users.value = []
    filteredUsers.value = []
  } finally {
    loading.value = false
  }
}

// Search functionality
function handleSearch(query) {
  searchQuery.value = query
  if (!query || query.length < 2) {
    filteredUsers.value = [...users.value]
    return
  }
  
  const searchTerm = query.toLowerCase().trim()
  filteredUsers.value = users.value.filter(user => {
    return (
      (user.nickname && user.nickname.toLowerCase().includes(searchTerm)) ||
      (user.first_name && user.first_name.toLowerCase().includes(searchTerm)) ||
      (user.last_name && user.last_name.toLowerCase().includes(searchTerm)) ||
      (user.email && user.email.toLowerCase().includes(searchTerm))
    )
  })
}



function resetFilters() {
  searchQuery.value = ''
  filteredUsers.value = [...users.value]
}

onMounted(fetchUsers)
</script>

<style scoped>
.discover-friend-container {
  max-width: 800px;
  margin: 2rem auto;
  padding: 2rem;
  background: #fff;
  border-radius: 16px;
  box-shadow: 0 4px 20px rgba(0,0,0,0.08);
}

.search-section {
  margin-bottom: 2rem;
  padding-bottom: 1.5rem;
  border-bottom: 1px solid #f0f0f0;
}

.simple-search {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #dbdbdb;
  border-radius: 12px;
  font-size: 16px;
  outline: none;
  transition: border-color 0.3s ease;
}

.simple-search:focus {
  border-color: #0095f6;
  box-shadow: 0 0 0 3px rgba(0, 149, 246, 0.1);
}

.simple-search::placeholder {
  color: #8e8e8e;
}

.results-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #f0f0f0;
}

.results-header h3 {
  margin: 0;
  color: #262626;
  font-weight: 600;
  font-size: 1.1rem;
}

.reset-btn {
  background: #0095f6;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s ease;
}

.reset-btn:hover {
  background: #0081d6;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 149, 246, 0.3);
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
  padding: 1.5rem;
  border: 1px solid #dbdbdb;
  border-radius: 12px;
  background: #fafafa;
  transition: all 0.3s ease;
  cursor: pointer;
}

.user-item:hover {
  background: #f0f0f0;
  transform: translateY(-2px);
  box-shadow: 0 4px 16px rgba(0,0,0,0.1);
  border-color: #0095f6;
}

.user-link {
  display: flex;
  align-items: center;
  gap: 1rem;
  text-decoration: none;
  color: inherit;
  width: 100%;
}

.avatar {
  width: 60px;
  height: 60px;
  border-radius: 50%;
  object-fit: cover;
  border: 3px solid #f0f0f0;
  transition: border-color 0.2s ease;
}

.user-item:hover .avatar {
  border-color: #0095f6;
}

.user-info {
  flex: 1;
  min-width: 0;
}

.user-info h4 {
  margin: 0 0 0.3rem 0;
  font-size: 1.2rem;
  font-weight: 600;
  color: #262626;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.user-handle {
  display: block;
  font-size: 0.9rem;
  color: #0095f6;
  font-weight: 500;
  margin-bottom: 0.2rem;
}

.user-email {
  display: block;
  font-size: 0.8rem;
  color: #8e8e8e;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.loading-spinner {
  text-align: center;
  color: #0095f6;
  font-weight: 600;
  margin: 3rem 0;
  font-size: 1.1rem;
}

.error {
  color: #dc2743;
  text-align: center;
  margin: 3rem 0;
  padding: 2rem;
  background: #fff5f5;
  border-radius: 12px;
  border: 1px solid #fed7d7;
}

.no-users {
  text-align: center;
  color: #8e8e8e;
  margin: 3rem 0;
  padding: 3rem 2rem;
}

.no-users-icon {
  font-size: 4rem;
  margin-bottom: 1rem;
  opacity: 0.6;
}

.no-users h3 {
  margin: 0 0 0.5rem 0;
  color: #262626;
  font-weight: 600;
}

.no-users p {
  margin: 0.5rem 0;
  font-size: 1rem;
}

/* Responsive Design */
@media (max-width: 768px) {
  .discover-friend-container {
    margin: 1rem;
    padding: 1.5rem;
    border-radius: 12px;
  }
  
  .results-header {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
  
  .user-item {
    padding: 1rem;
  }
  
  .avatar {
    width: 50px;
    height: 50px;
  }
  
  .user-info h4 {
    font-size: 1.1rem;
  }
  
  .no-users {
    padding: 2rem 1rem;
  }
  
  .no-users-icon {
    font-size: 3rem;
  }
}

/* Dark Mode Support */
@media (prefers-color-scheme: dark) {
  .discover-friend-container {
    background: #1a1a1a;
    box-shadow: 0 4px 20px rgba(0,0,0,0.3);
  }
  
  .search-section {
    border-bottom-color: #333333;
  }
  
  .results-header {
    border-bottom-color: #333333;
  }
  
  .results-header h3 {
    color: #ffffff;
  }
  
  .user-item {
    background: #2a2a2a;
    border-color: #333333;
  }
  
  .user-item:hover {
    background: #333333;
    border-color: #0095f6;
  }
  
  .user-info h4 {
    color: #ffffff;
  }
  
  .user-handle {
    color: #0095f6;
  }
  
  .user-email {
    color: #a0a0a0;
  }
  
  .error {
    background: #2a1a1a;
    border-color: #4a2a2a;
    color: #ff6b6b;
  }
  
  .no-users h3 {
    color: #ffffff;
  }
  
  .no-users p {
    color: #a0a0a0;
  }
}

/* Focus styles for accessibility */
.user-item:focus {
  outline: 2px solid #0095f6;
  outline-offset: 2px;
}

.reset-btn:focus {
  outline: 2px solid #0095f6;
  outline-offset: 2px;
}

/* Reduced motion for users who prefer it */
@media (prefers-reduced-motion: reduce) {
  .user-item,
  .avatar,
  .reset-btn {
    transition: none;
  }
  
  .user-item:hover {
    transform: none;
  }
}
</style> 