<template>
  <div class="profile-container">
    <!-- Banner Placeholder -->
    <div class="profile-banner">
      <div class="banner-image"></div>
    </div>

    <!-- Profile Info Card -->
    <div class="profile-header">
      <div class="avatar-container">
        <img class="avatar" :src="getAvatarUrl(profileUser.avatar_url)" alt="Profile Picture" />
        <!-- Upload button for profile owner -->
        <div v-if="isOwner" class="avatar-upload">
          <input 
            type="file" 
            ref="fileInput" 
            @change="handleFileSelect" 
            accept="image/*" 
            style="display: none;"
          />
          <button @click="$refs.fileInput.click()" class="upload-btn">
            📷 Change Photo
          </button>
        </div>
        <!-- Upload progress -->
        <div v-if="isUploading" class="upload-progress">
          <div class="progress-bar">
            <div class="progress-fill" :style="{ width: uploadProgress + '%' }"></div>
          </div>
          <p>Uploading... {{ uploadProgress }}%</p>
        </div>
      </div>
      <div class="profile-info">
        <h2>{{ profileUser.nickname || profileUser.first_name }}</h2>
        <!-- <p class="role"></p> -->
        <div class="profile-buttons">
          <button v-if="isOwner" @click="toggleVisibility"> {{ profileUser.is_public ? '🔓 Public' : '🔒 Private' }} </button>
          <template v-else>
            <button v-if="followStatus === 'none'" @click="toggleFollow('follow')"> Follow </button>
            <button v-else-if="followStatus === 'pending'" disabled> Pending </button>
            <button v-else-if="followStatus === 'accepted'" @click="toggleFollow('unfollow')"> Unfollow </button>
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
          <li><strong>Full Name:</strong> {{ profileUser.first_name || 'N/A' }} {{ profileUser.last_name || '' }}</li>
          <li><strong>Nickname:</strong> {{ profileUser.nickname || 'N/A' }}</li>
          <li><strong>Gender:</strong> {{ profileUser.gender || 'N/A' }}</li>
          <li><strong>Date Of Birth:</strong> {{ formatDate(profileUser.date_of_birth) }}</li>
          <li><strong>Email:</strong> {{ profileUser.email || 'N/A' }}</li>
          <li><strong>Member Since:</strong> {{ formatDate(profileUser.created_at) }}</li>
          <li><strong>About Me:</strong> {{ profileUser.about_me || 'N/A' }}</li>
        </ul>
      </div>

      <!-- Center Column -->
      <div class="profile-center">
        <div class="tabs">
          <span :class="{ active: activeTab === 'posts' }" @click="activeTab = 'posts'">Posts</span>
          <span :class="{ active: activeTab === 'followers' }" @click="activeTab = 'followers'">Followers</span>
          <span :class="{ active: activeTab === 'following' }" @click="activeTab = 'following'">Following</span>
          <span :class="{ active: activeTab === 'closeFriends' }" @click="activeTab = 'closeFriends'">Close Friends</span>
        </div>

        <div class="tab-content">
          <div v-if="activeTab === 'posts'">
            <div v-if="loadingPosts" class="loading-posts">
              <div class="loading-spinner"></div>
              <p>Loading posts...</p>
            </div>
            <div v-else-if="userPosts.length === 0" class="no-posts">
              <div class="no-posts-icon">📝</div>
              <h3>No posts yet</h3>
              <p>Start sharing your thoughts with the world!</p>
              <button v-if="isOwner" @click="fetchUserPosts" class="refresh-btn">
                🔄 Refresh
              </button>
            </div>
            <div v-else class="posts-list">
              <div class="posts-header">
                <h3>{{ userPosts.length }} post{{ userPosts.length !== 1 ? 's' : '' }}</h3>
                <button @click="fetchUserPosts" class="refresh-btn">
                  🔄 Refresh
                </button>
              </div>
              <div v-for="post in userPosts" :key="post.id" class="post-item">
                <div class="post-header">
                  <div class="post-author">
                    <img 
                      :src="getAvatarUrl(profileUser.avatar_url)" 
                      :alt="profileUser.nickname || profileUser.first_name"
                      class="post-avatar"
                    >
                    <div class="author-info">
                      <h4 class="author-name">{{ profileUser.nickname || profileUser.first_name }}</h4>
                      <span class="post-time">{{ formatDate(post.created_at) }}</span>
                      <span class="post-privacy">{{ post.privacy }}</span>
                    </div>
                  </div>
                </div>
                <div class="post-content">
                  <p v-if="post.content">{{ post.content }}</p>
                  <div v-if="post.image_url" class="post-image-container">
                    <img 
                      :src="getImageUrl(post.image_url)" 
                      :alt="post.content || 'Post image'"
                      class="post-image"
                      @error="handleImageError"
                    >
                  </div>
                </div>
                <div class="post-stats">
                  <span class="stat">💬 {{ post.replies || 0 }} replies</span>
                  <span class="stat">👁️ {{ post.views || 0 }} views</span>
                </div>
              </div>
            </div>
          </div>
          
          <div v-if="activeTab === 'followers'">
            <div v-if="followersList.length" class="users-list">
              <div v-for="f in followersList" :key="f.id" class="user-item">
                <img 
                  :src="getAvatarUrl(f.avatar_url)" 
                  :alt="f.nickname || f.first_name"
                  class="user-avatar"
                >
                <div class="user-info">
                  <h4 class="user-name">{{ f.nickname || f.first_name }}</h4>
                  <span class="user-handle">@{{ f.nickname }}</span>
                </div>
              </div>
            </div>
            <p v-else class="no-data">No followers yet.</p>
          </div>

          <div v-if="activeTab === 'following'">
            <div v-if="followingList.length" class="users-list">
              <div v-for="f in followingList" :key="f.id" class="user-item">
                <img 
                  :src="getAvatarUrl(f.avatar_url)" 
                  :alt="f.nickname || f.first_name"
                  class="user-avatar"
                >
                <div class="user-info">
                  <h4 class="user-name">{{ f.nickname || f.first_name }}</h4>
                  <span class="user-handle">@{{ f.nickname }}</span>
                </div>
              </div>
            </div>
            <p v-else class="no-data">Not following anyone yet.</p>
          </div>

          <div v-if="activeTab === 'closeFriends'">
            <div class="close-friends-section">
              <h4>Manage Close Friends</h4>
              <div v-if="closeFriendsLoading" class="loading-posts"><div class="loading-spinner"></div><p>Loading close friends...</p></div>
              <div v-if="closeFriendsError" class="no-data">{{ closeFriendsError }}</div>
              <div v-if="Array.isArray(followersList) && followersList.length && !closeFriendsLoading" class="users-list">
                <div v-for="f in followersList" :key="f.id" class="user-item">
                  <img 
                    :src="getAvatarUrl(f.avatar_url)" 
                    :alt="f.nickname || f.first_name"
                    class="user-avatar"
                  >
                  <div class="user-info">
                    <h4 class="user-name">{{ f.nickname || f.first_name }}</h4>
                    <span class="user-handle">@{{ f.nickname }}</span>
                  </div>
                  <button
                    v-if="!isInCloseFriends(f.id)"
                    @click="addToCloseFriends(f)"
                    class="refresh-btn"
                  >Add to Close Friends</button>
                  <button
                    v-else
                    @click="removeFromCloseFriends(f)"
                    class="refresh-btn"
                  >Remove</button>
                </div>
              </div>
              <p v-else-if="!closeFriendsLoading && !closeFriendsError" class="no-data">No followers to add as close friends.</p>
              <h4 style="margin-top:2rem;">Your Close Friends</h4>
              <div v-if="closeFriendsList.length && !closeFriendsLoading" class="users-list">
                <div v-for="cf in closeFriendsList" :key="cf.id" class="user-item">
                  <img 
                    :src="getAvatarUrl(cf.avatar_url)" 
                    :alt="cf.nickname || cf.first_name"
                    class="user-avatar"
                  >
                  <div class="user-info">
                    <h4 class="user-name">{{ cf.nickname || cf.first_name }}</h4>
                    <span class="user-handle">@{{ cf.nickname }}</span>
                  </div>
                  <button @click="removeFromCloseFriends(cf)" class="refresh-btn">Remove</button>
                </div>
              </div>
              <p v-else-if="!closeFriendsLoading && !closeFriendsError" class="no-data">You have no close friends yet.</p>
            </div>
          </div>

        </div>
      </div>

      <!-- Right Column -->
      <div class="profile-right">
        <h3>Coming Soon</h3>
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
import { onMounted, ref, watch } from 'vue'
import { useProfileView } from '@/composables/useProfile'
import { addToCloseFriends as apiAddToCloseFriends, removeFromCloseFriends as apiRemoveFromCloseFriends, fetchCloseFriends as apiFetchCloseFriends } from '@/services/api.js'

const {
  profileUser,
  followStatus,
  isOwner,
  activeTab,
  followersList,
  followingList,
  canViewPrivateProfile,
  initProfile,
  fetchProfile,
  toggleFollow,
  toggleVisibility
} = useProfileView()

// Profile image upload state
const fileInput = ref(null)
const isUploading = ref(false)
const uploadProgress = ref(0)

// Posts state
const userPosts = ref([])
const loadingPosts = ref(false)

// Close Friends state
const closeFriendsList = ref([])
const closeFriendsLoading = ref(false)
const closeFriendsError = ref('')

async function loadCloseFriends() {
  closeFriendsLoading.value = true
  closeFriendsError.value = ''
  try {
    // Always fetch close friends for the profile user being viewed
    const response = await fetch(`http://localhost:8080/users/${profileUser.id}/close-friends`, {
      credentials: 'include',
      headers: { 'Accept': 'application/json' }
    })
    if (!response.ok) {
      throw new Error('Failed to fetch close friends')
    }
    const result = await response.json()
    console.log('DEBUG: fetchCloseFriends result:', result)
    closeFriendsList.value = Array.isArray(result) ? result : []
  } catch (err) {
    console.error('DEBUG: fetchCloseFriends error:', err)
    closeFriendsError.value = err.message || 'Failed to load close friends.'
    closeFriendsList.value = []
  } finally {
    closeFriendsLoading.value = false
  }
}

function isInCloseFriends(userId) {
  return closeFriendsList.value.some(cf => cf.id === userId)
}

async function addToCloseFriends(user) {
  if (!isInCloseFriends(user.id)) {
    try {
      await apiAddToCloseFriends(user.id)
      await loadCloseFriends()
    } catch (err) {
      alert('Failed to add close friend: ' + (err.message || 'Unknown error'))
    }
  }
}

async function removeFromCloseFriends(user) {
  try {
    await apiRemoveFromCloseFriends(user.id)
    await loadCloseFriends()
  } catch (err) {
    alert('Failed to remove close friend: ' + (err.message || 'Unknown error'))
  }
}

onMounted(async () => {
  await initProfile()
  
  // If posts tab is active and we can view the profile, load posts
  if (activeTab.value === 'posts' && canViewPrivateProfile.value) {
    fetchUserPosts()
  }
  if (activeTab.value === 'closeFriends') {
    await loadCloseFriends()
  }
})

// Watch for tab changes to load posts
watch(activeTab, async (newTab) => {
  if (newTab === 'posts' && canViewPrivateProfile.value) {
    fetchUserPosts()
  }
  if (newTab === 'closeFriends') {
    await loadCloseFriends()
  }
})

// Watch for profile user changes to load posts when profile loads
watch(() => profileUser.id, (newId) => {
  console.log('Profile user ID changed:', newId)
  if (newId && activeTab.value === 'posts' && canViewPrivateProfile.value) {
    console.log('Auto-fetching posts for new profile user')
    fetchUserPosts()
  }
})

// Watch for profile access changes to load posts when access is granted
watch(canViewPrivateProfile, (canView) => {
  console.log('Profile access changed:', canView)
  if (canView && activeTab.value === 'posts' && profileUser.id) {
    console.log('Auto-fetching posts after access granted')
    fetchUserPosts()
  }
})

// Fetch user posts
async function fetchUserPosts() {
  if (!profileUser.id) {
    console.log('No profile user ID available')
    return
  }
  
  console.log('Fetching posts for user:', profileUser.id)
  loadingPosts.value = true
  
  try {
    const requestBody = {
      id: profileUser.id,
      type: 'user',
      start: 0,
      n_post: 50
    }
    
    console.log('Request body:', requestBody)
    
    const response = await fetch('http://localhost:8080/post/feed', {
      method: 'POST',
      credentials: 'include',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(requestBody)
    })

    console.log('Response status:', response.status)
    
    if (!response.ok) {
      const errorText = await response.text()
      console.error('Response error:', errorText)
      throw new Error(`Failed to fetch posts: ${response.status} ${errorText}`)
    }

    const posts = await response.json()
    console.log('Fetched posts:', posts)
    userPosts.value = posts || []
  } catch (error) {
    console.error('Error fetching posts:', error)
    userPosts.value = []
  } finally {
    loadingPosts.value = false
  }
}

// Format date helper
function formatDate(dateString) {
  if (!dateString) return 'N/A'
  
  try {
    const date = new Date(dateString)
    if (isNaN(date.getTime())) return 'N/A'
    
    return date.toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    })
  } catch (error) {
    return 'N/A'
  }
}

// Get image URL helper
function getImageUrl(imageUrl) {
  if (!imageUrl) return ''
  
  if (imageUrl.startsWith('uploads/')) {
    return `http://localhost:8080/images/${imageUrl.replace('uploads/', '')}`
  }
  
  return `http://localhost:8080${imageUrl}`
}

// Handle image error
function handleImageError(event) {
  console.log('Image failed to load:', event.target.src)
  event.target.style.display = 'none'
}

// Handle file selection for profile image upload
async function handleFileSelect(event) {
  const file = event.target.files[0]
  if (!file) return

  // Validate file type
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    alert('Please select a valid image file (JPEG, PNG, or GIF)')
    return
  }

  // Validate file size (5MB limit)
  if (file.size > 5 * 1024 * 1024) {
    alert('Image file must be smaller than 5MB')
    return
  }

  await uploadProfileImage(file)
}

// Upload profile image
async function uploadProfileImage(file) {
  isUploading.value = true
  uploadProgress.value = 0

  try {
    const formData = new FormData()
    formData.append('image', file)

    const response = await fetch('http://localhost:8080/upload/profile', {
      method: 'POST',
      credentials: 'include',
      body: formData
    })

    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Upload failed')
    }

    const result = await response.json()
    
    // Simulate progress for better UX
    for (let i = 0; i <= 100; i += 10) {
      uploadProgress.value = i
      await new Promise(resolve => setTimeout(resolve, 50))
    }

    // Refresh profile data to get the updated avatar_url
    await fetchProfile()
    
    alert('Profile image updated successfully!')
  } catch (error) {
    console.error('Upload error:', error)
    alert('Failed to upload image: ' + error.message)
  } finally {
    isUploading.value = false
    uploadProgress.value = 0
    // Reset file input
    if (fileInput.value) {
      fileInput.value.value = ''
    }
  }
}

// Helper function to get the full avatar URL
function getAvatarUrl(avatarUrl) {
  if (avatarUrl && avatarUrl.startsWith('/images/')) {
    return `http://localhost:8080${avatarUrl}`
  }
  if (avatarUrl && avatarUrl.startsWith('uploads/')) {
    // Convert uploads/filename to /images/filename format
    const filename = avatarUrl.replace('uploads/', '')
    return `http://localhost:8080/images/${filename}`
  }
  return '/images/default-avatar.png'
}
</script>

<style scoped>
/* Modern Instagram-Inspired Profile Design with Enhanced Compatibility */
.profile-container {
  background: #fafafa;
  min-height: 100vh;
  font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif;
  color: #262626;
  line-height: 1.6;
  -webkit-font-smoothing: antialiased;
  -moz-osx-font-smoothing: grayscale;
}

/* Banner Section - Modern gradient with fallbacks */
.profile-banner {
  position: relative;
  height: 200px;
  background: #f09433; /* Fallback for older browsers */
  background: linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%);
  border-radius: 0;
  overflow: hidden;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.banner-image {
  height: 100%;
  background: #f09433; /* Fallback */
  background: linear-gradient(45deg, #f09433 0%, #e6683c 25%, #dc2743 50%, #cc2366 75%, #bc1888 100%);
  position: relative;
  z-index: 1;
  display: flex;
  align-items: center;
  justify-content: center;
}

.banner-image::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: radial-gradient(circle at center, rgba(255, 255, 255, 0.1) 0%, transparent 70%);
  animation: pulse 3s ease-in-out infinite;
}

@keyframes pulse {
  0%, 100% { opacity: 0.3; }
  50% { opacity: 0.6; }
}

/* Profile Header - Enhanced modern design */
.profile-header {
  position: relative;
  margin: 0;
  background: white;
  border-radius: 0;
  padding: 3rem 2rem;
  border-bottom: 1px solid #dbdbdb;
  display: flex;
  align-items: flex-start;
  gap: 3rem;
  z-index: 10;
  max-width: 935px;
  margin: 0 auto;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
}

.profile-header::before {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  height: 1px;
  background: linear-gradient(90deg, transparent, #dbdbdb, transparent);
}

.avatar-container {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 1.5rem;
  flex-shrink: 0;
}

.avatar {
  width: 150px;
  height: 150px;
  border-radius: 50%;
  border: 4px solid white;
  object-fit: cover;
  background: #fafafa;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.15);
  position: relative;
  cursor: pointer;
}

.avatar::before {
  content: '';
  position: absolute;
  top: -4px;
  left: -4px;
  right: -4px;
  bottom: -4px;
  border-radius: 50%;
  background: linear-gradient(45deg, #f09433, #e6683c, #dc2743, #cc2366, #bc1888);
  z-index: -1;
  opacity: 0;
  transition: opacity 0.3s ease;
}

.avatar:hover {
  transform: scale(1.05);
  box-shadow: 0 12px 35px rgba(0, 0, 0, 0.2);
}

.avatar:hover::before {
  opacity: 1;
}

.avatar-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.upload-btn {
  background: #0095f6;
  color: white;
  padding: 0.6rem 1.2rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 600;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 149, 246, 0.3);
  position: relative;
  overflow: hidden;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.upload-btn::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
}

.upload-btn:hover {
  background: #0081d6;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 149, 246, 0.4);
}

.upload-btn:hover::before {
  left: 100%;
}

.upload-btn:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(0, 149, 246, 0.3);
}

.upload-progress {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0, 0, 0, 0.95);
  color: white;
  padding: 2rem;
  border-radius: 12px;
  text-align: center;
  min-width: 220px;
  backdrop-filter: blur(15px);
  -webkit-backdrop-filter: blur(15px);
  border: 1px solid rgba(255, 255, 255, 0.1);
  box-shadow: 0 20px 40px rgba(0, 0, 0, 0.3);
  z-index: 1000;
}

.progress-bar {
  width: 100%;
  height: 6px;
  background: rgba(255, 255, 255, 0.2);
  border-radius: 3px;
  overflow: hidden;
  margin-bottom: 1rem;
}

.progress-fill {
  height: 100%;
  background: linear-gradient(90deg, #0095f6, #0081d6);
  transition: width 0.3s ease;
  border-radius: 3px;
  position: relative;
}

.progress-fill::after {
  content: '';
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.3), transparent);
  animation: shimmer 1.5s infinite;
}

@keyframes shimmer {
  0% { transform: translateX(-100%); }
  100% { transform: translateX(100%); }
}

.upload-progress p {
  margin: 0;
  font-size: 0.9rem;
  font-weight: 600;
  color: #f0f0f0;
}

.profile-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  justify-content: flex-start;
  gap: 1rem;
}

.profile-info h2 {
  font-size: 2rem;
  font-weight: 700;
  margin-bottom: 0.5rem;
  color: #262626;
  position: relative;
  line-height: 1.2;
}

.profile-info h2::after {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 30px;
  height: 2px;
  background: linear-gradient(90deg, #0095f6, #0081d6);
  border-radius: 1px;
}

.profile-buttons {
  display: flex;
  gap: 0.8rem;
  margin-top: 0.5rem;
  flex-wrap: wrap;
  align-items: center;
}

.profile-buttons button {
  background: #0095f6;
  color: white;
  padding: 0.7rem 1.2rem;
  border: none;
  border-radius: 8px;
  cursor: pointer;
  font-weight: 600;
  font-size: 0.9rem;
  transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  box-shadow: 0 2px 8px rgba(0, 149, 246, 0.2);
  position: relative;
  overflow: hidden;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.profile-buttons button::before {
  content: '';
  position: absolute;
  top: 0;
  left: -100%;
  width: 100%;
  height: 100%;
  background: linear-gradient(90deg, transparent, rgba(255, 255, 255, 0.2), transparent);
  transition: left 0.5s ease;
}

.profile-buttons button:hover {
  background: #0081d6;
  transform: translateY(-1px);
  box-shadow: 0 4px 12px rgba(0, 149, 246, 0.3);
}

.profile-buttons button:hover::before {
  left: 100%;
}

.profile-buttons button:active {
  transform: translateY(0);
  box-shadow: 0 2px 8px rgba(0, 149, 246, 0.2);
}

.profile-buttons button:disabled {
  background: #dbdbdb;
  color: #8e8e8e;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
  opacity: 0.6;
}

.profile-buttons button:disabled::before {
  display: none;
}

/* Locked Profile - Enhanced design */
.locked-profile {
  display: flex;
  justify-content: center;
  padding: 4rem 2rem;
  min-height: 60vh;
  align-items: center;
}

.locked-card {
  background: white;
  border: 1px solid #dbdbdb;
  padding: 3rem;
  border-radius: 12px;
  text-align: center;
  max-width: 500px;
  width: 100%;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  transition: all 0.3s ease;
}

.locked-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
}

.locked-card h3 {
  color: #262626;
  margin-bottom: 1rem;
  font-size: 1.5rem;
  font-weight: 600;
}

.locked-card p {
  color: #8e8e8e;
  font-size: 1.1rem;
  line-height: 1.6;
}

/* Main Layout - Enhanced grid system */
.profile-main {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  gap: 2rem;
  padding: 2rem;
  max-width: 935px;
  margin: 0 auto;
}

/* Left Column - Enhanced styling */
.profile-left {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid #dbdbdb;
  height: fit-content;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.profile-left:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.profile-left h3 {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #262626;
  position: relative;
}

.profile-left h3::after {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 20px;
  height: 2px;
  background: linear-gradient(90deg, #0095f6, #0081d6);
  border-radius: 1px;
}

.about-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.about-list li {
  margin-bottom: 0.8rem;
  padding: 0.8rem;
  background: #fafafa;
  border-radius: 8px;
  border-left: 3px solid #0095f6;
  transition: all 0.2s ease;
}

.about-list li:hover {
  background: #f0f0f0;
  transform: translateX(2px);
}

.about-list li strong {
  color: #262626;
  font-weight: 600;
  display: block;
  margin-bottom: 0.2rem;
}

/* Center Column - Enhanced tabs */
.profile-center {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid #dbdbdb;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.profile-center:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.tabs {
  display: flex;
  gap: 0;
  margin-bottom: 1.5rem;
  border-bottom: 1px solid #dbdbdb;
  background: white;
  padding: 0;
  border-radius: 8px 8px 0 0;
  overflow: hidden;
}

.tabs span {
  flex: 1;
  padding: 1rem;
  cursor: pointer;
  font-weight: 600;
  color: #8e8e8e;
  text-align: center;
  border-radius: 0;
  transition: all 0.2s ease;
  position: relative;
  border-bottom: 2px solid transparent;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.tabs span:hover {
  color: #262626;
  background: #fafafa;
}

.tabs span.active {
  background: white;
  color: #262626;
  border-bottom-color: #262626;
}

.tab-content {
  min-height: 300px;
}

.tab-content div {
  animation: fadeIn 0.3s ease-in-out;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Posts List - Enhanced design */
.loading-posts, .no-posts, .no-data {
  text-align: center;
  padding: 3rem 2rem;
  color: #8e8e8e;
  font-style: italic;
}

.loading-spinner {
  width: 40px;
  height: 40px;
  border: 3px solid #f3f3f3;
  border-top: 3px solid #0095f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.no-posts-icon {
  font-size: 3rem;
  margin-bottom: 1rem;
  opacity: 0.6;
}

.no-posts h3 {
  margin: 0 0 0.5rem 0;
  color: #262626;
  font-weight: 600;
}

.posts-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 1.5rem;
  padding-bottom: 1rem;
  border-bottom: 1px solid #dbdbdb;
}

.posts-header h3 {
  margin: 0;
  color: #262626;
  font-weight: 600;
}

.refresh-btn {
  background: #0095f6;
  color: white;
  border: none;
  padding: 0.5rem 1rem;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.9rem;
  font-weight: 500;
  transition: all 0.2s ease;
  user-select: none;
  -webkit-user-select: none;
  -moz-user-select: none;
  -ms-user-select: none;
}

.refresh-btn:hover {
  background: #0081d6;
  transform: translateY(-1px);
  box-shadow: 0 2px 8px rgba(0, 149, 246, 0.3);
}

.refresh-btn:active {
  transform: translateY(0);
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.post-item {
  background: white;
  border: 1px solid #dbdbdb;
  border-radius: 12px;
  padding: 1.5rem;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.post-item:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
  transform: translateY(-2px);
}

.post-header {
  margin-bottom: 1rem;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 0.8rem;
}

.post-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #f0f0f0;
}

.author-info {
  flex: 1;
}

.author-name {
  font-weight: 600;
  margin: 0 0 0.2rem 0;
  color: #262626;
}

.post-time {
  font-size: 0.8rem;
  color: #8e8e8e;
  margin-right: 0.5rem;
}

.post-privacy {
  font-size: 0.8rem;
  color: #0095f6;
  background: rgba(0, 149, 246, 0.1);
  padding: 0.2rem 0.5rem;
  border-radius: 4px;
  font-weight: 500;
}

.post-content {
  margin-bottom: 1rem;
}

.post-content p {
  margin: 0 0 1rem 0;
  line-height: 1.6;
  color: #262626;
  font-size: 0.95rem;
}

.post-image-container {
  margin-top: 1rem;
  border-radius: 8px;
  overflow: hidden;
}

.post-image {
  width: 100%;
  max-height: 400px;
  object-fit: cover;
  border-radius: 8px;
  transition: transform 0.3s ease;
}

.post-image:hover {
  transform: scale(1.02);
}

.post-stats {
  display: flex;
  gap: 1rem;
  padding-top: 1rem;
  border-top: 1px solid #f0f0f0;
}

.stat {
  font-size: 0.9rem;
  color: #8e8e8e;
  display: flex;
  align-items: center;
  gap: 0.3rem;
}

/* Users List - Enhanced design */
.users-list {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.user-item {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem;
  background: white;
  border: 1px solid #dbdbdb;
  border-radius: 12px;
  transition: all 0.3s ease;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
}

.user-item:hover {
  background: #fafafa;
  transform: translateX(4px);
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.user-avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid #f0f0f0;
}

.user-info {
  flex: 1;
}

.user-name {
  font-weight: 600;
  margin: 0 0 0.2rem 0;
  color: #262626;
}

.user-handle {
  font-size: 0.9rem;
  color: #8e8e8e;
}

/* Right Column - Enhanced styling */
.profile-right {
  background: white;
  border-radius: 12px;
  padding: 1.5rem;
  border: 1px solid #dbdbdb;
  height: fit-content;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.05);
  transition: all 0.3s ease;
}

.profile-right:hover {
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.1);
}

.profile-right h3 {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #262626;
  position: relative;
}

.profile-right h3::after {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 20px;
  height: 2px;
  background: linear-gradient(90deg, #0095f6, #0081d6);
  border-radius: 1px;
}

/* Close Friends Section - Enhanced styling */
.close-friends-section {
  margin-top: 2rem;
  padding-top: 1.5rem;
  border-top: 1px solid #dbdbdb;
}

.close-friends-section h4 {
  font-size: 1.1rem;
  font-weight: 600;
  margin-bottom: 1rem;
  color: #262626;
  position: relative;
}

.close-friends-section h4::after {
  content: '';
  position: absolute;
  bottom: -5px;
  left: 0;
  width: 20px;
  height: 2px;
  background: linear-gradient(90deg, #0095f6, #0081d6);
  border-radius: 1px;
}

/* Enhanced Responsive Design */
@media (max-width: 1024px) {
  .profile-main {
    grid-template-columns: 1fr;
    gap: 1.5rem;
    padding: 1.5rem;
  }
  
  .profile-header {
    flex-direction: column;
    align-items: center;
    text-align: center;
    padding: 2rem 1.5rem;
  }
  
  .profile-banner {
    height: 150px;
  }
  
  .avatar {
    width: 120px;
    height: 120px;
  }
  
  .profile-info h2 {
    font-size: 1.6rem;
  }
  
  .profile-buttons {
    justify-content: center;
  }
}

@media (max-width: 768px) {
  .profile-container {
    padding: 0;
  }
  
  .profile-main {
    padding: 1rem;
    gap: 1rem;
  }
  
  .profile-banner {
    height: 120px;
  }
  
  .profile-header {
    padding: 1.5rem 1rem;
    gap: 2rem;
  }
  
  .avatar {
    width: 100px;
    height: 100px;
  }
  
  .profile-info h2 {
    font-size: 1.4rem;
  }
  
  .profile-buttons {
    justify-content: center;
    gap: 0.5rem;
  }
  
  .profile-buttons button {
    padding: 0.6rem 1rem;
    font-size: 0.85rem;
  }
  
  .tabs {
    flex-direction: column;
    gap: 0;
  }
  
  .tabs span {
    flex: none;
    border-bottom: 1px solid #dbdbdb;
  }
  
  .tabs span.active {
    border-bottom-color: #262626;
  }
  
  .post-item {
    padding: 1rem;
  }
  
  .user-item {
    padding: 0.8rem;
  }
  
  .user-avatar {
    width: 40px;
    height: 40px;
  }
}

@media (max-width: 480px) {
  .profile-header {
    padding: 1rem;
    gap: 1.5rem;
  }
  
  .avatar {
    width: 80px;
    height: 80px;
  }
  
  .profile-info h2 {
    font-size: 1.2rem;
  }
  
  .profile-buttons {
    flex-direction: column;
    width: 100%;
  }
  
  .profile-buttons button {
    width: 100%;
  }
  
  .posts-header {
    flex-direction: column;
    gap: 1rem;
    align-items: flex-start;
  }
  
  .post-stats {
    flex-direction: column;
    gap: 0.5rem;
  }
}

/* Enhanced scrollbar styling */
::-webkit-scrollbar {
  width: 8px;
}

::-webkit-scrollbar-track {
  background: #fafafa;
  border-radius: 4px;
}

::-webkit-scrollbar-thumb {
  background: #dbdbdb;
  border-radius: 4px;
  transition: background 0.2s ease;
}

::-webkit-scrollbar-thumb:hover {
  background: #c7c7c7;
}

/* Firefox scrollbar */
* {
  scrollbar-width: thin;
  scrollbar-color: #dbdbdb #fafafa;
}

/* Smooth scrolling with fallback */
html {
  scroll-behavior: smooth;
}

/* Focus styles for accessibility */
button:focus,
input:focus {
  outline: 2px solid #0095f6;
  outline-offset: 2px;
}

/* Reduced motion for users who prefer it */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
  
  .avatar:hover,
  .post-item:hover,
  .user-item:hover {
    transform: none;
  }
}

/* High contrast mode support */
@media (prefers-contrast: high) {
  .profile-container {
    background: white;
  }
  
  .profile-header,
  .profile-left,
  .profile-center,
  .profile-right {
    border: 2px solid #000;
  }
  
  .tabs span.active {
    border-bottom-color: #000;
  }
}

/* Print styles */
@media print {
  .profile-container {
    background: white;
  }
  
  .profile-banner,
  .upload-btn,
  .profile-buttons,
  .refresh-btn {
    display: none;
  }
  
  .profile-header {
    border: 1px solid #000;
  }
}
</style>
