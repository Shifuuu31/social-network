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
  width: 100%;
  max-width: var(--max-width);
  margin: 0 auto;
  min-height: 100vh;
  padding: 2rem;
  background: var(--bg-dark);
  color: var(--text-light);
}

/* Top banner */
.profile-banner .banner-image {
  height: 200px;
  background: linear-gradient(to right, var(--accent-purple), var(--bg-purple));
  border-radius: 8px 8px 0 0;
}

/* Profile Header */
.profile-header {
  display: flex;
  align-items: center;
  padding: 1rem 2rem;
  background: #1a1a1a;
  border-radius: 0 0 8px 8px;
  margin-bottom: 2rem;
  gap: 1.5rem;
}

.avatar {
  width: 100px;
  height: 100px;
  border-radius: 100px;
  border: 4px solid var(--accent-purple);
  object-fit: cover;
  background: #2e2e2e;
}

.profile-info {
  flex: 1;
}

.profile-info h2 {
  font-size: 1.8rem;
  margin-bottom: 0.2rem;
}

.profile-buttons {
  display: flex;
  gap: 1rem;
  margin-top: 0.5rem;
}

button {
  background: var(--accent-purple);
  color: white;
  padding: 0.5rem 1.2rem;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-weight: bold;
  transition: background 0.3s;
}

button:hover {
  background: var(--accent-hover);
}

button:disabled {
  background: #555;
  cursor: not-allowed;
}

/* Main Layout */
.profile-main {
  display: grid;
  grid-template-columns: 1fr 2fr 1fr;
  gap: 2rem;
}

/* Left Column */
.profile-left,
.profile-right {
  background: #1a1a1a;
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
  background: #212121;
  border-radius: 8px;
  padding: 1rem;
}

.tabs {
  display: flex;
  gap: 2rem;
  margin-bottom: 1rem;
  border-bottom: 2px solid #333;
}

.tabs span {
  padding: 0.5rem;
  cursor: pointer;
  font-weight: bold;
  color: var(--text-muted);
}

.tabs span.active {
  color: var(--accent-purple);
  border-bottom: 3px solid var(--accent-purple);
}

/* Locked Card */
.locked-profile {
  display: flex;
  justify-content: center;
  padding: 2rem;
}

.locked-card {
  background: #1a1a1a;
  border: 1px solid #333;
  padding: 2rem;
  border-radius: 12px;
  text-align: center;
  max-width: 500px;
  width: 100%;
  box-shadow: 0 2px 8px rgba(0,0,0,0.2);
}

.locked-card h3 {
  color: var(--accent-purple);
  margin-bottom: 0.5rem;
}
</style>
