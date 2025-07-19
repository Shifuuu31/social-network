<template>
  <div class="chat-list">
    <div class="chat-header">
      <h3>Messages</h3>
      <div class="header-controls">
        <div class="connection-status" :class="{ connected: isConnected }">
          {{ isConnected ? '●' : '○' }}
        </div>
        <button @click="toggleView" class="toggle-btn">
          {{ showNewChat ? 'Recent' : 'New Chat' }}
        </button>
      </div>
    </div>

    <!-- Recent Conversations View -->
    <div v-if="!showNewChat">
      <div v-if="loading" class="loading">
        <div class="loading-spinner"></div>
        <p>Loading conversations...</p>
      </div>

      <div v-else-if="error" class="error">
        {{ error }}
      </div>

      <div v-else-if="conversations.length === 0" class="no-conversations">
        <p>No conversations yet</p>
        <p class="hint">Start chatting with people you follow!</p>
        <button @click="showNewChat = true" class="start-chat-btn">
          Start New Chat
        </button>
      </div>

      <div v-else class="conversations">
        <div
          v-for="conversation in conversations"
          :key="conversation.id"
          class="conversation-item"
          :class="{ active: selectedUserId === getOtherUserId(conversation) }"
          @click="selectConversation(conversation)"
        >
          <div class="avatar">
            <img 
              :src="getAvatarUrl(getOtherUser(conversation).avatar_url)" 
              :alt="getOtherUser(conversation).sender_name || getOtherUser(conversation).receiver_name"
            />
          </div>
          <div class="conversation-info">
            <div class="name">
              {{ getOtherUser(conversation).sender_name || getOtherUser(conversation).receiver_name }}
            </div>
            <div class="last-message">
              {{ conversation.content }}
            </div>
            <div class="timestamp">
              {{ formatTimestamp(conversation.created_at) }}
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- New Chat View -->
    <div v-else>
      <div v-if="loadingFollowing" class="loading">
        <div class="loading-spinner"></div>
        <p>Loading people you follow...</p>
      </div>

      <div v-else-if="followingError" class="error">
        {{ followingError }}
      </div>

      <div v-else-if="followingList.length === 0" class="no-following">
        <p>You're not following anyone yet</p>
        <p class="hint">Follow people to start chatting with them!</p>
        <router-link to="/discover-friend" class="discover-btn">
          Discover People
        </router-link>
      </div>

      <div v-else class="following-list">
        <div
          v-for="user in followingList"
          :key="user.id"
          class="user-item"
          :class="{ active: selectedUserId === user.id }"
          @click="selectUser(user)"
        >
          <div class="avatar">
            <img 
              :src="getAvatarUrl(user.avatar_url)" 
              :alt="user.nickname || user.first_name"
            />
          </div>
          <div class="user-info">
            <div class="name">
              {{ user.nickname || user.first_name }}
            </div>
            <div class="status">
              {{ hasConversation(user.id) ? 'Has conversation' : 'Start new chat' }}
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script>
import { ref, onMounted, onUnmounted } from 'vue'
import { useAuth } from '../../composables/useAuth.js'
import chatService from '../../services/chatService.js'

// Debug logging utility
const DEBUG = true;
const debugLog = (method, message, data = null) => {
  if (DEBUG) {
    console.log(`🔍 [ChatList.${method}] ${message}`, data || '');
  }
};

export default {
  name: 'ChatList',
  props: {
    selectedUserId: {
      type: Number,
      default: null
    }
  },
  emits: ['select-conversation'],
  setup(props, { emit }) {
    const { user } = useAuth()
    const conversations = ref([])
    const followingList = ref([])
    const loading = ref(false)
    const loadingFollowing = ref(false)
    const error = ref(null)
    const followingError = ref(null)
    const isConnected = ref(false)
    const showNewChat = ref(false)

    debugLog('setup', `Component initialized with user:`, user.value);

    const loadConversations = async () => {
      debugLog('loadConversations', 'Loading recent conversations');
      loading.value = true
      error.value = null
      
      try {
        const response = await chatService.getRecentConversations(20)
        debugLog('loadConversations', 'Received response:', response);
        conversations.value = response.conversations || []
        debugLog('loadConversations', `Loaded ${conversations.value.length} conversations`);
      } catch (err) {
        debugLog('loadConversations', 'Error loading conversations:', err);
        error.value = err.message
        console.error('Error loading conversations:', err)
      } finally {
        loading.value = false
      }
    }

    const loadFollowing = async () => {
      debugLog('loadFollowing', 'Loading following list');
      loadingFollowing.value = true
      followingError.value = null
      
      try {
        console.log('Loading following list...')
        debugLog('loadFollowing', 'Making request to /api/users/following');
        
        const response = await fetch('/api/users/following', {
          method: 'GET',
          credentials: 'include',
          headers: { 'Accept': 'application/json' }
        })
        
        debugLog('loadFollowing', `Response status: ${response.status} ${response.statusText}`);
        console.log('Response status:', response.status)
        console.log('Response headers:', response.headers)
        
        if (!response.ok) {
          const errorText = await response.text()
          debugLog('loadFollowing', 'Error response:', errorText);
          console.error('Error response:', errorText)
          throw new Error(`Failed to load following list: ${response.status} ${errorText}`)
        }
        
        const data = await response.json()
        debugLog('loadFollowing', 'Following data received:', data);
        console.log('Following data:', data)
        followingList.value = data.following || []
        debugLog('loadFollowing', `Following list length: ${followingList.value.length}`);
        console.log('Following list length:', followingList.value.length)
      } catch (err) {
        debugLog('loadFollowing', 'Error loading following:', err);
        console.error('Error loading following:', err)
        followingError.value = err.message
      } finally {
        loadingFollowing.value = false
      }
    }

    const connectWebSocket = async () => {
      debugLog('connectWebSocket', `Attempting to connect WebSocket for user ${user.value?.id}`);
      try {
        await chatService.connect(user.value?.id)
        isConnected.value = true
        debugLog('connectWebSocket', 'WebSocket connected successfully');
      } catch (err) {
        debugLog('connectWebSocket', 'Failed to connect WebSocket:', err);
        console.error('Failed to connect WebSocket:', err)
        isConnected.value = false
      }
    }

    const getOtherUserId = (conversation) => {
      const otherUserId = conversation.sender_id === user.value?.id 
        ? conversation.receiver_id 
        : conversation.sender_id;
      // Only log when debugging specific issues, not on every render
      // debugLog('getOtherUserId', `Getting other user ID from conversation: ${otherUserId}`);
      return otherUserId;
    }

    const getOtherUser = (conversation) => {
      const otherUserId = getOtherUserId(conversation)
      
      // Use conversation data directly to avoid recursive updates
      return {
        id: otherUserId,
        nickname: conversation.sender_id === otherUserId ? conversation.sender_name : conversation.receiver_name,
        first_name: conversation.sender_id === otherUserId ? conversation.sender_name : conversation.receiver_name,
        avatar_url: conversation.sender_id === otherUserId ? conversation.sender_avatar : null
      }
    }

    const selectConversation = (conversation) => {
      const otherUserId = getOtherUserId(conversation)
      debugLog('selectConversation', `Selecting conversation with user ${otherUserId}`);
      emit('select-conversation', otherUserId)
    }

    const selectUser = (user) => {
      debugLog('selectUser', `Selecting user ${user.id} for new conversation`);
      emit('select-conversation', user.id)
    }

    const toggleView = () => {
      debugLog('toggleView', `Toggling view from ${showNewChat.value} to ${!showNewChat.value}`);
      showNewChat.value = !showNewChat.value
    }

    const hasConversation = (userId) => {
      return conversations.value.some(conv => 
        conv.sender_id === userId || conv.receiver_id === userId
      )
    }

    const getAvatarUrl = (avatarUrl) => {
      if (!avatarUrl) return '/default-avatar.png'
      if (avatarUrl.startsWith('http')) return avatarUrl
      return `/api/images/${avatarUrl}`
    }

    const formatTimestamp = (timestamp) => {
      const date = new Date(timestamp)
      const now = new Date()
      const diffInHours = (now - date) / (1000 * 60 * 60)
      
      if (diffInHours < 24) {
        return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
      } else if (diffInHours < 48) {
        return 'Yesterday'
      } else {
        return date.toLocaleDateString()
      }
    }

    // WebSocket message handler
    const handleMessage = (message) => {
      // Refresh conversations when new message arrives
      loadConversations()
    }

    // Connection status handler
    const handleConnectionChange = (connected) => {
      isConnected.value = connected
    }

    onMounted(async () => {
      debugLog('onMounted', 'Component mounted');
      await connectWebSocket()
      await loadConversations()
      await loadFollowing()
      
      // Set up WebSocket handlers
      chatService.onMessage(handleMessage)
      chatService.onConnectionChange(handleConnectionChange)
    })

    onUnmounted(() => {
      debugLog('onUnmounted', 'Component unmounted');
      chatService.disconnect()
    })

    return {
      conversations,
      followingList,
      loading,
      loadingFollowing,
      error,
      followingError,
      isConnected,
      showNewChat,
      loadConversations,
      loadFollowing,
      selectConversation,
      selectUser,
      toggleView,
      hasConversation,
      getOtherUserId,
      getOtherUser,
      getAvatarUrl,
      formatTimestamp
    }
  }
}
</script>

<style scoped>
.chat-list {
  width: 300px;
  border-right: 1px solid #e1e5e9;
  background: #f8f9fa;
  display: flex;
  flex-direction: column;
}

.chat-header {
  padding: 1rem;
  border-bottom: 1px solid #e1e5e9;
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: white;
}

.chat-header h3 {
  margin: 0;
  color: #1a1a1a;
}

.header-controls {
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.connection-status {
  font-size: 1.2rem;
  color: #dc3545;
}

.connection-status.connected {
  color: #28a745;
}

.toggle-btn {
  background: #007bff;
  color: white;
  border: none;
  border-radius: 0.5rem;
  padding: 0.25rem 0.75rem;
  font-size: 0.8rem;
  cursor: pointer;
  transition: background-color 0.2s;
}

.toggle-btn:hover {
  background: #0056b3;
}

.loading, .error, .no-conversations, .no-following {
  padding: 2rem;
  text-align: center;
  color: #6c757d;
}

.loading-spinner {
  border: 2px solid #f3f3f3;
  border-top: 2px solid #007bff;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  animation: spin 1s linear infinite;
  margin: 0 auto 1rem;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.conversations, .following-list {
  flex: 1;
  overflow-y: auto;
}

.conversation-item, .user-item {
  padding: 1rem;
  border-bottom: 1px solid #e1e5e9;
  cursor: pointer;
  transition: background-color 0.2s;
  display: flex;
  align-items: center;
  background: white;
}

.conversation-item:hover, .user-item:hover {
  background: #f8f9fa;
}

.conversation-item.active, .user-item.active {
  background: #e3f2fd;
  border-left: 3px solid #007bff;
}

.avatar {
  width: 50px;
  height: 50px;
  border-radius: 50%;
  overflow: hidden;
  margin-right: 1rem;
  flex-shrink: 0;
}

.avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.conversation-info, .user-info {
  flex: 1;
  min-width: 0;
}

.name {
  font-weight: 600;
  color: #1a1a1a;
  margin-bottom: 0.25rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.last-message, .status {
  color: #6c757d;
  font-size: 0.9rem;
  margin-bottom: 0.25rem;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.timestamp {
  color: #adb5bd;
  font-size: 0.8rem;
}

.hint {
  font-size: 0.9rem;
  color: #adb5bd;
  margin-top: 0.5rem;
}

.start-chat-btn, .discover-btn {
  background: #007bff;
  color: white;
  border: none;
  border-radius: 0.5rem;
  padding: 0.5rem 1rem;
  margin-top: 1rem;
  cursor: pointer;
  text-decoration: none;
  display: inline-block;
  transition: background-color 0.2s;
}

.start-chat-btn:hover, .discover-btn:hover {
  background: #0056b3;
}
</style> 