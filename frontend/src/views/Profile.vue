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
          <li><strong>Gender:</strong> {{ profileUser.gender || 'N/A' }}</li>
          <li><strong>Date Of Birth:</strong> {{ profileUser.date_of_birth || 'N/A' }}</li>
          <li><strong>Email:</strong> {{ profileUser.email || 'N/A' }}</li>
          <li><strong>About Me:</strong> {{ profileUser.about_me || 'N/A' }}</li>
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
        <!-- //todo  -->
          <div v-if="activeTab === 'followers'">
            <ul v-if="followersList.length">
              <li v-for="f in followersList" :key="f.id">
                {{ f.nickname || f.first_name }} — @{{ f.nickname }}
              </li>
            </ul>
            <p v-else>No followers yet.</p>
          </div>

          <div v-if="activeTab === 'following'">
            <ul v-if="followingList.length">
              <li v-for="f in followingList" :key="f.id">
                {{ f.nickname || f.first_name }} — @{{ f.nickname }}
              </li>
            </ul>
            <p v-else>Not following anyone yet.</p>
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
import { onMounted, ref } from 'vue'
import { useProfileView } from '@/composables/useProfile'

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

onMounted(
  initProfile
)

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
    // Extract filename from /images/filename
    const filename = avatarUrl.replace('/images/', '')
    return `http://localhost:8080${avatarUrl}`
  }
  return '/images/default-avatar.png'
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

.avatar-container {
  position: relative;
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 100px;
  border: 4px solid white;
  object-fit: cover;
  background: #eee;
  transition: transform 0.2s ease;
}

.avatar:hover {
  transform: scale(1.05);
}

.avatar-upload {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 0.5rem;
}

.upload-btn {
  background: #6a0dad;
  color: white;
  padding: 0.4rem 0.8rem;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 0.8rem;
  font-weight: bold;
  transition: all 0.3s ease;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.upload-btn:hover {
  background: #7d20c0;
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0,0,0,0.15);
}

.upload-progress {
  position: absolute;
  top: 50%;
  left: 50%;
  transform: translate(-50%, -50%);
  background: rgba(0, 0, 0, 0.8);
  color: white;
  padding: 1rem;
  border-radius: 8px;
  text-align: center;
  min-width: 150px;
}

.progress-bar {
  width: 100%;
  height: 8px;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 4px;
  overflow: hidden;
  margin-bottom: 0.5rem;
}

.progress-fill {
  height: 100%;
  background: #6a0dad;
  transition: width 0.3s ease;
}

.upload-progress p {
  margin: 0;
  font-size: 0.8rem;
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
