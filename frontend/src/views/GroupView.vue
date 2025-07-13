<template>
  <div class="group-view" v-if="groupsStore.currentGroup">
    <!-- Header -->
    <div class="group-header">
      <div class="group-header-bg">
        <img :src="groupsStore.currentGroup.image" :alt="groupsStore.currentGroup.name" />
      </div>
      <div class="group-header-content">
        <div class="container">
          <div class="group-info">
            <h1 class="group-name">{{ groupsStore.currentGroup.name }}</h1>
            <p class="group-description">{{ groupsStore.currentGroup.description }}</p>
            <div class="group-meta">
              <span class="member-count">
                <span class="icon">👥</span>
                {{ groupsStore.currentGroup.memberCount }} members
              </span>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Main Content Area -->
    <div class="group-content">
      <div class="container">
        <div class="content-layout">
          <!-- Sidebar -->
          <div class="sidebar">
            <div class="sidebar-section">
              <h3>invite </h3>
              <div class="sidebar-actions">
                <!-- Not a member -->
                <button v-if="!groupsStore.currentGroup.isMember" class="btn btn-primary sidebar-btn"
                  @click="handleJoinGroup" :disabled="isJoining">
                  <span class="icon">+</span>
                  {{ isJoining ? 'Joining...' : 'Join Group' }}
                </button>

                <!-- Requested to join -->
                <button v-else-if="groupsStore.currentGroup.isMember === 'requested'" class="btn btn-grey sidebar-btn"
                  disabled>
                  <span class="icon">⏳</span>
                  Request Sent
                </button>

                <!-- Invited to join -->
                <button v-else-if="groupsStore.currentGroup.isMember === 'invited'" class="btn btn-grey sidebar-btn"
                  @click="handleAcceptInvite" :disabled="isJoining">
                  <span class="icon">📨</span>
                  {{ isJoining ? 'Accepting...' : 'Accept Invite' }}
                </button>

                <!-- Full member -->
                <button v-else-if="groupsStore.currentGroup.isMember === 'member'"
                  class="btn btn-success sidebar-btn desactivated">
                  <span class="icon">✓</span>
                  Member
                </button>

                <button v-if="groupsStore.currentGroup.isMember === 'member'" class="btn btn-outline sidebar-btn"
                  @click="toggleInviteModal">
                  <span class="icon">📧</span>
                  Invite
                </button>
              </div>
            </div>

            <div class="sidebar-section">
              <!-- <h3>sidebar where you toggle posts and events</h3> -->
              <div class="sidebar-nav">
                <button :class="['nav-btn', { active: activeTab === 'posts' }]" @click="setActiveTab('posts')">
                  <span class="icon">📝</span>
                  Posts
                </button>
                <button :class="['nav-btn', { active: activeTab === 'events' }]" @click="setActiveTab('events')">
                  <span class="icon">📅</span>
                  Events
                </button>
              </div>
            </div>
          </div>

          <!-- Main Content -->
          <div class="main-content">
            <div class="content-header">
              <h2>{{ activeTab === 'posts' ? 'Posts' : 'Events' }}</h2>
            </div>

            <div class="tab-content">
              <!-- Posts Tab -->
              <div v-if="activeTab === 'posts'" class="posts-section">
                <div class="create-post" v-if="groupsStore.currentGroup.isMember === 'member'">
                  <div class="create-post-header">
                    <h3>create a post</h3>
                  </div>
                  <form @submit.prevent="handleCreatePost" class="create-post-form">
                    <!-- <input type="text" v-model="newPost.title" placeholder="Title of your post..." class="form-input"
                       required /> -->
                    <textarea v-model="newPost.content" placeholder="Share something with the group..." 
                     class="form-textarea" rows="4" required></textarea>
                    <div class="form-actions">
                      <button type="submit" class="btn btn-primary" :disabled="isCreatingPost">
                        {{ isCreatingPost ? 'Publishing...' : 'Publish' }}
                      </button>
                    </div>
                  </form>
                </div>

                <div class="posts-list">
                  <div v-if="isLoadingPosts" class="loading">
                    <div class="spinner"></div>
                    <p>Loading posts...</p>
                  </div>
                  <div v-else-if="groupsStore.groupPosts.length === 0" class="empty-state">
                    <div class="empty-icon">📝</div>
                    <h3>No posts yet</h3>
                    <p v-if="groupsStore.currentGroup.isMember === 'member'">Be the first to share something!</p>
                    <p v-else>Join the group to see and share content.</p>
                  </div>
                  <div v-else class="posts-grid">
                    <div v-for="post in groupsStore.groupPosts" :key="post.id" class="post-card">
                      <div class="post-header">
                        <img :src="post.authorAvatar" :alt="post.author" class="author-avatar" />
                        <div class="post-meta">
                          <h4 class="author-name">{{ post.author }}</h4>
                          <span class="post-date">{{ formatDate(post.createdAt) }}</span>
                        </div>
                      </div>
                      <div class="post-content">
                        <h3 class="post-title">{{ post.title }}</h3>
                        <p class="post-text">{{ post.content }}</p>
                      </div>
                      <div class="post-actions">
                        <button class="post-action">
                          <span class="icon">💬</span>
                          {{ post.comments }}
                        </button>
                      </div>
                    </div>
                  </div>
                </div>
              </div>

              <!-- Events Tab -->
              <div v-else-if="activeTab === 'events'" class="events-section">
                <div class="create-event" v-if="groupsStore.currentGroup.isMember === 'member'">
                  <div class="create-event-header">
                    <h3>create an event</h3>
                  </div>
                  <form @submit.prevent="handleCreateEvent" class="create-event-form">
                    <input type="text" v-model="newEvent.title" placeholder="Event title..." class="form-input"
                      required />
                    <textarea v-model="newEvent.description" placeholder="Event description..." class="form-textarea"
                      rows="3" required></textarea>
                    <div class="form-row">
                      <input type="datetime-local" v-model="newEvent.date" class="form-input" required />
                      <!-- <input type="text" v-model="newEvent.location" placeholder="Location" class="form-input" /> -->
                    </div>
                    <div class="form-actions">
                      <button type="submit" class="btn btn-primary" :disabled="isCreatingEvent">
                        {{ isCreatingEvent ? 'Creating...' : 'Create Event' }}
                      </button>
                    </div>
                  </form>
                </div>

                <div class="events-list">
                  <div v-if="isLoadingEvents" class="loading">
                    <div class="spinner"></div>
                    <p>Loading events...</p>
                  </div>
                  <div v-else-if="groupsStore.groupEvents.length === 0" class="empty-state">
                    <div class="empty-icon">📅</div>
                    <h3>No events yet</h3>
                    <p v-if="groupsStore.currentGroup.isMember === 'member'">Create the first event!</p>
                    <p v-else>Join the group to see and create events.</p>
                  </div>
                  <div v-else class="events-grid">
                    <div v-for="event in groupsStore.groupEvents" :key="event.id" class="event-card">
                      <div class="event-header">
                        <div class="event-date">
                          <span class="day">{{ formatEventDay(event.date) }}</span>
                          <span class="month">{{ formatEventMonth(event.date) }}</span>
                        </div>
                        <div class="event-meta">
                          <h4 class="event-title">{{ event.title }}</h4>
                          <p class="event-time">{{ formatEventTime(event.date) }}</p>
                          <!-- <p class="event-location" v-if="event.location">📍 {{ event.location }}</p> -->
                        </div>
                      </div>
                      <div class="event-content">
                        <p class="event-description">{{ event.description }}</p>
                      </div>
                      <div class="event-actions">
                        <button class="event-action">
                          <span class="icon">👥</span>
                          {{ event.attendees || 0 }} attending
                        </button>
                        <div class="event-buttons">
                          <button class="event-action btn-attend" @click="handleAttendEvent(event.id, 'going')"
                            :class="{ active: event.isAttending === 'going' }">
                            <span class="icon">✓</span>
                            Going
                          </button>
                          <button class="event-action btn-not-attend" @click="handleAttendEvent(event.id, 'not_going')"
                            :class="{ active: event.isAttending === 'not_going' }">
                            <span class="icon">✗</span>
                            Not Going
                          </button>
                        </div>
                      </div>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Invite Modal -->
    <div v-if="showInviteModal" class="modal-overlay" @click="toggleInviteModal">
      <div class="modal-content" @click.stop>
        <div class="modal-body">
          <input v-model="userSearch" @input="fetchUsers" placeholder="Search users by name or nickname..."
            class="form-input" />
          <div v-if="isLoadingUsers" class="loading">
            <div class="spinner"></div>Loading users...
          </div>
          <div v-else-if="usersList.length === 0" class="empty-state">No users found.</div>
          <ul v-else class="user-list">
            <li v-for="user in usersList" :key="user.id" class="user-list-item">
              <span>{{ user.nickname }} ({{ user.first_name }} {{ user.last_name }})</span>
              <button class="btn btn-primary btn-sm" @click="inviteUser(+user.id)">
                {{ invitedUserIds.includes(user.id) ? 'Invited' : 'Invite' }}
              </button>
            </li>
          </ul>
        </div>
      </div>
    </div>
  </div>

  <!-- Loading state for the entire component -->
  <div v-else-if="groupsStore.isLoading" class="loading">
    <div class="spinner"></div>
    <p>Loading group...</p>
  </div>

  <!-- Error state -->
  <div v-else-if="groupsStore.error" class="error-state">
    <h3>Error loading group</h3>
    <p>{{ groupsStore.error }}</p>
    <button class="btn btn-primary" @click="loadGroup">Retry</button>
  </div>
</template>



<script setup>
import { ref, reactive, onMounted, watch } from 'vue'
import { useRoute } from 'vue-router'
import { useGroupsStore } from '@/stores/groups'
// import { debounce } from 'lodash-es'

const route = useRoute()
const groupsStore = useGroupsStore()

// State
const activeTab = ref('posts')
const isJoining = ref(false)
const isCreatingPost = ref(false)
const isCreatingEvent = ref(false)
const isLoadingPosts = ref(false)
const isLoadingEvents = ref(false)
const showInviteModal = ref(false)
const userSearch = ref('')
const usersList = ref([])
const isLoadingUsers = ref(false)
const invitedUserIds = ref([])

// Form data
const newPost = reactive({
  title: '',
  content: ''
})

const newEvent = reactive({
  title: '',
  description: '',
  date: '',
})

// Debounced user search
// const debouncedFetchUsers = debounce(fetchUsers, 300)

// Methods
const setActiveTab = (tab) => {
  activeTab.value = tab
  if (tab === 'posts' && groupsStore.groupPosts.length === 0) {
    loadPosts()
  } else if (tab === 'events' && groupsStore.groupEvents.length === 0) {
    loadEvents()
  }
}

const loadGroup = async () => {
  try {
    const groupId = parseInt(route.params.id)
    await groupsStore.fetchGroup(groupId)
  } catch (error) {
    console.error('Failed to load group:', error)
  }
}

const loadPosts = async () => {
  if (isLoadingPosts.value) return

  isLoadingPosts.value = true
  try {
    const groupId = parseInt(route.params.id)
    await groupsStore.fetchGroupPosts(groupId)
  } catch (error) {
    console.error('Failed to load posts:', error)
  } finally {
    isLoadingPosts.value = false
  }
}

const loadEvents = async () => {
  if (isLoadingEvents.value) return

  isLoadingEvents.value = true
  try {
    const groupId = parseInt(route.params.id)
    await groupsStore.fetchGroupEvents(groupId)
  } catch (error) {
    console.error('Failed to load events:', error)
  } finally {
    isLoadingEvents.value = false
  }
}

const handleAttendEvent = async (eventId, voteType) => {
  try {
    await groupsStore.attendEvent(eventId, voteType)
  } catch (error) {
    console.error('Failed to attend event:', error)
  }
}

const handleJoinGroup = async () => {
  if (isJoining.value) return

  isJoining.value = true
  try {
    const groupId = parseInt(route.params.id)
    await groupsStore.requestJoinGroup(groupId)
    // The requestJoinGroup function now updates the local state
    // so no need to manually reload the group
  } catch (error) {
    console.error('Failed to join group:', error)
  } finally {
    isJoining.value = false
  }
}

const handleAcceptInvite = async () => {
  if (isJoining.value) return

  isJoining.value = true
  try {
    const groupId = parseInt(route.params.id)
    await groupsStore.acceptGroupInvite(groupId)
    await loadGroup()
  } catch (error) {
    console.error('Failed to accept invite:', error)
  } finally {
    isJoining.value = false
  }
}

const handleCreatePost = async () => {
  if (isCreatingPost.value) return

  isCreatingPost.value = true
  try {
    const groupId = parseInt(route.params.id)
    await groupsStore.createPost(groupId, {
      title: newPost.title,
      content: newPost.content
    })
    newPost.title = ''
    newPost.content = ''
  } catch (error) {
    console.error('Failed to create post:', error)
  } finally {
    isCreatingPost.value = false
  }
}

const handleCreateEvent = async () => {
  if (isCreatingEvent.value) return

  isCreatingEvent.value = true
  try {
    const groupId = parseInt(route.params.id)
    await groupsStore.createEvent(groupId, {
      title: newEvent.title,
      description: newEvent.description,
      date: newEvent.date,
    })
    newEvent.title = ''
    newEvent.description = ''
    newEvent.date = ''
  } catch (error) {
    console.error('Failed to create event:', error)
  } finally {
    isCreatingEvent.value = false
  }
}

const toggleInviteModal = () => {
  showInviteModal.value = !showInviteModal.value
  if (showInviteModal.value) {
    userSearch.value = ''
    usersList.value = []
    invitedUserIds.value = []
    fetchAllUsers() // Add this new function
  }
}

// Add this new function
const fetchAllUsers = async () => {
  isLoadingUsers.value = true
  try {
    const groupId = parseInt(route.params.id)
    const response = await fetch(`/api/groups/group/${groupId}/available-users`, { 
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      // credentials: 'include'
    })
    if (!response.ok) throw new Error('Failed to fetch available users')
    usersList.value = await response.json()
  } catch (e) {
    console.error('Error fetching available users:', e)
    usersList.value = []
  } finally {
    isLoadingUsers.value = false
  }
}

const fetchUsers = async () => {
  // if (!userSearch.value.trim()) {
  //   usersList.value = []
  //   return
  // }

  isLoadingUsers.value = true
  try {
    const groupId = parseInt(route.params.id)
    const response = await fetch(`/api/groups/group/${groupId}/search-users?q=${encodeURIComponent(userSearch.value)}`, {
      method: 'GET',
      headers: { 'Content-Type': 'application/json' },
      // credentials: 'include'
    })

    if (!response.ok) throw new Error('Failed to fetch users')

    const data = await response.json()
    usersList.value = data || [] // Use the filtered users from the group-specific endpoint
  } catch (e) {
    console.error('Error fetching users:', e)
    usersList.value = []
  } finally {
    isLoadingUsers.value = false
  }
}

const inviteUser = async (userId) => {
  try {
    const groupId = parseInt(route.params.id)
    await groupsStore.inviteUserToGroup(groupId, userId)
    invitedUserIds.value.push(userId)
  } catch (e) {
    console.error('Error inviting user:', e)
  }
}

// Formatting helpers
const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('default', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatEventDay = (dateString) => {
  const date = new Date(dateString)
  return date.getDate()
}

const formatEventMonth = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('fr-FR', { month: 'short' })
}

const formatEventTime = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleTimeString('fr-FR', { hour: '2-digit', minute: '2-digit' })
}

// Watchers
watch(
  () => route.params.id,
  async (newId) => {
    if (newId) {
      await loadGroup()
      await loadPosts()
      await loadEvents()
    }
  },
  { immediate: false }
)

watch(userSearch, () => {
  fetchUsers()
})

// Lifecycle hooks
onMounted(async () => {
  await loadGroup()
  await loadPosts()
  if (activeTab.value === 'events') {
    await loadEvents()
  }
})
</script>
<style scoped>
.group-view {
  min-height: 100vh;
}

.group-header {
  position: relative;
  height: 250px;
  overflow: hidden;
}

.group-header-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1;
}

.group-header-bg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  filter: brightness(0.4);
}

.group-header-content {
  position: relative;
  z-index: 2;
  height: 100%;
  display: flex;
  align-items: end;
  padding: 40px 20px;
  background: linear-gradient(to top, rgba(0, 0, 0, 0.8), transparent);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.group-info {
  text-align: center;
}

.group-name {
  font-size: 2.5rem;
  font-weight: 700;
  color: #fff;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0, 0, 0, 0.5);
}

.group-description {
  font-size: 1.1rem;
  color: #ccc;
  margin-bottom: 15px;
  line-height: 1.5;
}

.group-meta {
  display: flex;
  justify-content: center;
  gap: 15px;
  align-items: center;
  flex-wrap: wrap;
}

.member-count {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #ccc;
}

.group-content {
  padding: 40px 20px;
}

.content-layout {
  display: grid;
  grid-template-columns: 250px 1fr;
  gap: 40px;
}

.sidebar {
  background: #1a1a1a;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #333;
  height: fit-content;
  position: sticky;
  top: 20px;
}

.sidebar-section {
  margin-bottom: 30px;
}

.sidebar-section:last-child {
  margin-bottom: 0;
}

.sidebar-section h3 {
  color: #fff;
  font-size: 0.9rem;
  margin-bottom: 15px;
  text-transform: lowercase;
  font-weight: 500;
}

.sidebar-actions {
  display: flex;
  flex-direction: column;
  gap: 10px;
}

.sidebar-btn {
  width: 100%;
  justify-content: flex-start;
  font-size: 0.85rem;
  padding: 10px 16px;
}

.sidebar-nav {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.nav-btn {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 16px;
  background: transparent;
  border: none;
  color: #ccc;
  cursor: pointer;
  border-radius: 8px;
  transition: all 0.2s ease;
  font-size: 0.9rem;
  text-align: left;
}

.nav-btn:hover {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.nav-btn.active {
  background: rgba(139, 92, 246, 0.2);
  color: #8b5cf6;
}

.main-content {
  flex: 1;
}

.content-header {
  margin-bottom: 30px;
}

.content-header h2 {
  color: #fff;
  font-size: 2rem;
  margin: 0;
  font-weight: 600;
}

.create-post,
.create-event {
  background: #1a1a1a;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 30px;
  border: 1px solid #333;
}

.create-post-header,
.create-event-header {
  margin-bottom: 20px;
}

.create-post-header h3,
.create-event-header h3 {
  color: #fff;
  font-size: 1.2rem;
  margin: 0;
}

.create-post-form,
.create-event-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-row {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 16px;
}

.form-input,
.form-textarea {
  padding: 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: #fff;
  font-size: 1rem;
  transition: border-color 0.2s ease;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #8b5cf6;
}

.form-textarea {
  resize: vertical;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.posts-grid,
.events-grid {
  display: grid;
  gap: 24px;
}

.post-card,
.event-card {
  background: #1a1a1a;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #333;
  transition: transform 0.2s ease;
}

.post-card:hover,
.event-card:hover {
  transform: translateY(-2px);
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 15px;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.post-meta {
  flex: 1;
}

.author-name {
  font-size: 1rem;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.post-date {
  font-size: 0.85rem;
  color: #666;
}

.post-content {
  margin-bottom: 15px;
}

.post-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 8px 0;
}

.post-text {
  color: #ccc;
  line-height: 1.5;
  margin: 0;
}

.post-actions {
  display: flex;
  gap: 16px;
}

.post-action {
  display: flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  font-size: 0.9rem;
  transition: color 0.2s ease;
}

.post-action:hover {
  color: #fff;
}

.event-header {
  display: flex;
  gap: 16px;
  margin-bottom: 15px;
}

.event-date {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: rgba(139, 92, 246, 0.2);
  border-radius: 8px;
  padding: 12px;
  min-width: 60px;
}

.event-date .day {
  font-size: 1.5rem;
  font-weight: 700;
  color: #8b5cf6;
}

.event-date .month {
  font-size: 0.8rem;
  color: #8b5cf6;
  text-transform: uppercase;
}

.event-meta {
  flex: 1;
}

.event-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 4px 0;
}

.event-time {
  font-size: 0.9rem;
  color: #ccc;
  margin: 0 0 4px 0;
}



.event-description {
  color: #ccc;
  line-height: 1.5;
  margin: 0 0 15px 0;
}

.event-actions {
  display: flex;
  gap: 16px;
  justify-content: space-between;
  align-items: center;
}

.event-action {
  display: flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  font-size: 0.9rem;
  transition: color 0.2s ease;
}

.event-action:hover {
  color: #fff;
}

.btn-attend {
  background: rgba(139, 92, 246, 0.2);
  color: #8b5cf6;
  padding: 8px 16px;
  border-radius: 6px;
}

.btn-attend:hover {
  background: rgba(139, 92, 246, 0.3);
  color: #fff;
}

.btn-attend.active {
  background: #8b5cf6;
  color: #fff;
}

.btn-not-attend {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  padding: 8px 16px;
  border-radius: 6px;
}

.btn-not-attend:hover {
  background: rgba(239, 68, 68, 0.3);
  color: #fff;
}

.btn-not-attend.active {
  background: #ef4444;
  color: #fff;
}

.event-buttons {
  display: flex;
  gap: 8px;
}

.btn-accent {
  background: linear-gradient(135deg, #f59e0b, #d97706);
  color: #fff;
}

.btn-accent:hover:not(:disabled) {
  background: linear-gradient(135deg, #d97706, #b45309);
  transform: translateY(-1px);
}

.btn-success {
  background: linear-gradient(135deg, #10b981, #059669);
  color: #fff;
}

.btn-success:hover:not(:disabled) {
  background: linear-gradient(135deg, #059669, #047857);
  transform: translateY(-1px);
}

.btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 8px;
  justify-content: center;
  font-size: 0.9rem;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
  transform: translateY(-1px);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn-secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.15);
}

.btn-grey {
  background: rgba(156, 163, 175, 0.2);
  color: #9ca3af;
  border: 1px solid #374151;
}

.btn-grey:hover:not(:disabled) {
  background: rgba(156, 163, 175, 0.3);
  color: #d1d5db;
}

.btn-outline {
  background: transparent;
  color: #8b5cf6;
  border: 1px solid #8b5cf6;
}

.btn-outline:hover {
  background: #8b5cf6;
  color: #fff;
}

.loading {
  text-align: center;
  padding: 60px 20px;
  color: #ccc;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-top: 4px solid #8b5cf6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}



@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.empty-state h3 {
  font-size: 1.5rem;
  color: #fff;
  margin-bottom: 10px;
}

.empty-state p {
  color: #ccc;
}

.error-state {
  text-align: center;
  padding: 60px 20px;
}

.error-state h3 {
  color: #ef4444;
  margin-bottom: 10px;
}

.error-state p {
  color: #ccc;
  margin-bottom: 20px;
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: #1a1a1a;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
  border: 1px solid #333;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #333;
}

.modal-header h3 {
  color: #fff;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  color: #ccc;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 4px;
  transition: color 0.2s ease;
}

.close-btn:hover {
  color: #fff;
}

.modal-body {
  padding: 20px;
}

.modal-body p {
  color: #ccc;
  margin-bottom: 20px;
}

.invite-link {
  display: flex;
  gap: 10px;
}

.invite-input {
  flex: 1;
  padding: 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: #fff;
  font-size: 0.9rem;
}

.icon {
  font-size: 1rem;
}

.user-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.user-list-item {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 8px 0;
  border-bottom: 1px solid #333;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 0.85rem;
}

@media (max-width: 768px) {
  .content-layout {
    grid-template-columns: 1fr;
    gap: 20px;
  }

  .sidebar {
    position: static;
    order: 2;
  }

  .main-content {
    order: 1;
  }

  .group-name {
    font-size: 2rem;
  }

  .sidebar-actions {
    flex-direction: row;
    gap: 8px;
  }

  .sidebar-btn {
    flex: 1;
  }

  .form-row {
    grid-template-columns: 1fr;
  }

  .invite-link {
    flex-direction: column;
  }
}




.desactivated {
  pointer-events: none;
  opacity: 0.5;
  cursor: default;
}
</style>