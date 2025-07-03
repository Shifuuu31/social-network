<template>
  <div class="sidebar">


    <h2 v-if="activeView === 'users'">Users</h2>
    <h2 v-if="activeView === 'groups'">Groups</h2>

    <div class="search-container">
      <input v-model="searchQuery" type="text" placeholder="Search users or groups..." class="search-input" />
      <span class="search-icon">
        <svg width="16" height="16" viewBox="0 0 24 24" fill="none" xmlns="http://www.w3.org/2000/svg">
          <path
            d="M11 19C15.4183 19 19 15.4183 19 11C19 6.58172 15.4183 3 11 3C6.58172 3 3 6.58172 3 11C3 15.4183 6.58172 19 11 19Z"
            stroke="#888" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
          <path d="M21 21L16.65 16.65" stroke="#888" stroke-width="2" stroke-linecap="round" stroke-linejoin="round" />
        </svg>
      </span>

    </div>
    <div v-if="activeView === 'users' && !filteredUsers.length" class="empty-state">
      <p v-if="searchQuery">No users match your search</p>
      <p v-else>No users available</p>
    </div>

    <div v-if="activeView === 'groups' && !filteredGroups.length" class="empty-state">
      <p v-if="searchQuery">No groups match your search</p>
      <p v-else>No groups available</p>
    </div>

    <div class="list-container">

      <div v-if="activeView === 'users'">
        <div v-for="u in filteredUsers" :key="u.id" @click="select('private', u.id)" class="chat-preview">
          <img :src="u.avatarUrl || defaultAvatar" alt="avatar" class="avatar" />
          <div class="chat-info">
            <div class="top-row">
              <strong>{{ u.name }}</strong>
              <small class="timestamp">{{ formatTimestamp(u.lastMessageTimestamp) }}</small>
            </div>
            <div class="bottom-row">
              <span class="last-message">{{ u.lastMessage || 'No messages yet' }}</span>
              <span v-if="u.unreadCount > 0" class="unread-count">{{ u.unreadCount }}</span>
            </div>
          </div>
          <div class="status-indicators">
            <span v-if="u.online" class="online-dot" title="Online"></span>
            <span v-if="u.typing" class="typing-indicator">typing...</span>
          </div>
        </div>
      </div>

      <div v-if="activeView === 'groups'">
        <div v-for="g in filteredGroups" :key="g.id" @click="select('group', g.id)" class="chat-preview">
          <img :src="g.avatarUrl || defaultGroupAvatar" alt="group avatar" class="avatar" />
          <div class="chat-info">
            <div class="top-row">
              <strong>{{ g.name }}</strong>
              <small class="timestamp">{{ formatTimestamp(g.lastMessageTimestamp) }}</small>
            </div>
            <div class="bottom-row">
              <span class="last-message">{{ g.lastMessage || 'No messages yet' }}</span>
              <span v-if="g.unreadCount > 0" class="unread-count">{{ g.unreadCount }}</span>
            </div>
          </div>
          <div class="status-indicators">
            <span v-if="g.typing" class="typing-indicator">typing...</span>
          </div>
        </div>
      </div>
    </div>

    <div class="toggle-buttons">
      <button :class="{ active: activeView === 'users' }" @click="activeView = 'users'; store.activeType= 'private' ">
        Users
      </button>
      <button :class="{ active: activeView === 'groups' }" @click="activeView = 'groups'; ; store.activeType= 'group'">
        Groups
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useChatStore } from '../stores/chatStore'


const store = useChatStore()
const activeView = ref('users')
const searchQuery = ref('')

const defaultAvatar = '/src/assets/avatar.png'
const defaultGroupAvatar = '/src/assets/group_avatar.png'

const filteredUsers = computed(() => {
  if (!searchQuery.value) return store.users
  const query = searchQuery.value.toLowerCase()
  return store.users.filter(u =>
    u.name.toLowerCase().includes(query) ||
    (u.lastMessage && u.lastMessage.toLowerCase().includes(query))
  )
})

const filteredGroups = computed(() => {
  if (!searchQuery.value) return store.groups
  const query = searchQuery.value.toLowerCase()
  return store.groups.filter(g =>
    g.name.toLowerCase().includes(query) ||
    (g.lastMessage && g.lastMessage.toLowerCase().includes(query))
  )
})

function select(type, id) {
  store.activeType = type
  store.activeTargetId = id
}

function formatTimestamp(ts) {
  if (!ts) return ''
  const date = new Date(ts)
  const now = new Date()
  const diffMs = now - date
  const diffMinutes = Math.floor(diffMs / 60000)

  if (diffMinutes < 1) return 'now'
  if (diffMinutes < 60) return `${diffMinutes}m ago`
  if (diffMinutes < 1440) return `${Math.floor(diffMinutes / 60)}h ago`
  return date.toLocaleDateString()
}
</script>

<style scoped>
.sidebar {
  display: flex;
  flex-direction: column;
  height: 100vh;
  border-right: 1px solid #ddd;
}

.search-container {
  position: relative;
  padding: 0.5em;
  border-bottom: 1px solid #eee;
}

.search-input {
  width: 100%;
  padding: 8px 30px 8px 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  outline: none;
}

.search-input:focus {
  border-color: #007bff;
}

.search-icon {
  position: absolute;
  right: 20px;
  top: 50%;
  transform: translateY(-50%);
  color: #888;
}

.list-container {
  flex: 1 1 auto;
  overflow-y: auto;
  padding: 0.5em;
}

.empty-state {
  padding: 1em;
  text-align: center;
  color: #888;
}

.toggle-buttons {
  flex-shrink: 0;
  display: flex;
  gap: 8px;
  padding: 0.5em;
  border-top: 1px solid #ccc;
  background-color: #fff;
}

.toggle-buttons button {
  padding: 0.5em 1em;
  cursor: pointer;
  background: none;
  border: 1px solid #ccc;
  border-radius: 4px;
}

.toggle-buttons button.active {
  background-color: #007bff;
  color: white;
  border-color: #007bff;
}

.chat-preview {
  display: flex;
  align-items: center;
  padding: 0.5em;
  cursor: pointer;
  border-bottom: 1px solid #eee;
}

.chat-preview:hover {
  background-color: #f0f0f0;
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 10px;
}

.chat-info {
  flex: 1;
  display: flex;
  flex-direction: column;
  overflow: hidden;
}

.top-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.timestamp {
  font-size: 0.75em;
  color: #888;
  margin-left: 10px;
  white-space: nowrap;
}

.bottom-row {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 2px;
  overflow: hidden;
}

.last-message {
  color: #555;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.unread-count {
  background-color: #007bff;
  color: white;
  font-weight: bold;
  padding: 0 6px;
  border-radius: 12px;
  font-size: 0.75em;
  margin-left: 8px;
  min-width: 20px;
  text-align: center;
}

.status-indicators {
  display: flex;
  flex-direction: column;
  align-items: flex-end;
  font-size: 0.75em;
  color: #007bff;
  min-width: 50px;
}

.online-dot {
  width: 10px;
  height: 10px;
  background-color: #28a745;
  border-radius: 50%;
  margin-bottom: 4px;
}

.typing-indicator {
  font-style: italic;
  color: #666;
}
</style>