<template>
  <div class="profile-container">
    <!-- Banner Placeholder -->
    <div class="profile-banner">
      <div class="banner-image"></div>
    </div>

    <!-- Profile Info Card -->
    <div class="profile-header">
      <img class="avatar" :src="profileUser.avatar" alt="Profile Picture" />
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
import { onMounted } from 'vue'
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
  toggleFollow,
  toggleVisibility
} = useProfileView()

onMounted(
  initProfile
)
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
