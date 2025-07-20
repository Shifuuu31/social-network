<template>
  <div class="unseen-notifications-container">
    <!-- Header -->
    <div class="notification-header">
      <h2 class="text-2xl font-bold text-white mb-4">
        Unseen Notifications
        <span v-if="unseenCount > 0" class="unread-badge">{{ unseenCount }}</span>
      </h2>
      
      <!-- Actions -->
      <div class="header-actions mb-4">
        <button 
          @click="markAllAsRead"
          v-if="unseenNotifications.length > 0"
          class="btn-primary mr-2"
        >
          Mark All as Read
        </button>
        <button 
          @click="refreshNotifications"
          class="btn-secondary"
        >
          Refresh
        </button>
      </div>
    </div>

    <!-- Loading State -->
    <div v-if="loading" class="loading-state">
      <div class="spinner"></div>
      <p class="text-gray-400 mt-2">Loading notifications...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="unseenNotifications.length === 0" class="empty-state">
      <div class="empty-icon">🔔</div>
      <h3 class="text-xl text-gray-300 mb-2">All caught up!</h3>
      <p class="text-gray-400">You have no unseen notifications.</p>
    </div>

    <!-- Notifications List -->
    <div v-else class="notifications-list">
      <div
        v-for="notification in unseenNotifications"
        :key="notification.id"
        class="notification-card"
        :class="getNotificationTypeClass(notification.type)"
      >
        <!-- Notification Icon -->
        <div class="notification-icon">
          {{ getNotificationIcon(notification.type) }}
        </div>

        <!-- Notification Content -->
        <div class="notification-content">
          <div class="notification-title">
            {{ getNotificationTitle(notification.type) }}
            <span class="new-badge">NEW</span>
          </div>
          
          <div class="notification-message">
            {{ notification.message }}
          </div>
          
          <div class="notification-meta">
            <span class="notification-time">
              {{ formatRelativeTime(notification.created_at) }}
            </span>
            <span class="notification-type-label">
              {{ notification.type.replace('_', ' ').toUpperCase() }}
            </span>
          </div>
        </div>

        <!-- Notification Actions -->
        <div class="notification-actions">
          <button 
            @click="markAsRead(notification.id)"
            class="action-btn mark-read"
            title="Mark as read"
          >
            ✓
          </button>
          
          <button 
            @click="handleNotificationAction(notification)"
            v-if="hasActionButton(notification.type)"
            class="action-btn primary"
          >
            {{ getActionButtonText(notification.type) }}
          </button>
          
          <button 
            @click="removeNotification(notification.id)"
            class="action-btn remove"
            title="Remove notification"
          >
            ×
          </button>
        </div>
      </div>
    </div>

    <!-- Load More Button -->
    <div v-if="hasMore && unseenNotifications.length > 0" class="load-more-container">
      <button 
        @click="loadMore"
        :disabled="loadingMore"
        class="btn-secondary w-full"
      >
        {{ loadingMore ? 'Loading...' : 'Load More' }}
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useNotificationStore } from '@/stores/notificationStore'
import chatService from '@/services/chatService'

const notificationStore = useNotificationStore()

// Reactive state
const loading = ref(true)
const loadingMore = ref(false)
const hasMore = ref(true)
const currentPage = ref(1)

// Computed properties
const unseenNotifications = computed(() => {
  return notificationStore.notifications.filter(notification => !notification.seen)
})

const unseenCount = computed(() => unseenNotifications.value.length)

// Methods
const refreshNotifications = async () => {
  loading.value = true
  try {
    await notificationStore.fetchNotifications()
    currentPage.value = 1
    hasMore.value = true
  } catch (error) {
    console.error('Failed to refresh notifications:', error)
  } finally {
    loading.value = false
  }
}

const loadMore = async () => {
  if (loadingMore.value) return
  
  loadingMore.value = true
  try {
    currentPage.value++
    await notificationStore.fetchNotifications(currentPage.value)
    // Check if we have fewer notifications than expected (indicating no more pages)
    if (unseenNotifications.value.length < currentPage.value * 20) {
      hasMore.value = false
    }
  } catch (error) {
    console.error('Failed to load more notifications:', error)
    currentPage.value-- // Revert page increment on error
  } finally {
    loadingMore.value = false
  }
}

const markAsRead = async (notificationId) => {
  try {
    await notificationStore.markAsRead(notificationId)
  } catch (error) {
    console.error('Failed to mark notification as read:', error)
  }
}

const markAllAsRead = async () => {
  try {
    const unseenIds = unseenNotifications.value.map(n => n.id)
    await notificationStore.markMultipleAsRead(unseenIds)
  } catch (error) {
    console.error('Failed to mark all notifications as read:', error)
  }
}

const removeNotification = async (notificationId) => {
  try {
    await notificationStore.removeNotification(notificationId)
  } catch (error) {
    console.error('Failed to remove notification:', error)
  }
}

const handleNotificationAction = (notification) => {
  // Handle specific actions based on notification type
  switch (notification.type) {
    case 'group_invite':
      // Navigate to group invitation page or show accept/decline modal
      console.log('Handle group invite:', notification)
      break
    case 'group_join_request':
      // Navigate to group management or show approve/reject modal
      console.log('Handle join request:', notification)
      break
    case 'follow_request':
      // Navigate to follow requests page
      console.log('Handle follow request:', notification)
      break
    default:
      console.log('No specific action for notification type:', notification.type)
  }
}

// Utility functions
const getNotificationTitle = (type) => {
  const titles = {
    'group_invite': 'Group Invitation',
    'group_join_request': 'Join Request', 
    'group_event': 'New Event',
    'follow_request': 'Follow Request',
    'message': 'New Message',
    'notification': 'Notification'
  }
  return titles[type] || 'Notification'
}

const getNotificationIcon = (type) => {
  const icons = {
    'group_invite': '👥',
    'group_join_request': '🚪',
    'group_event': '📅',
    'follow_request': '👤',
    'message': '💬',
    'notification': '🔔'
  }
  return icons[type] || '🔔'
}

const getNotificationTypeClass = (type) => {
  const classes = {
    'group_invite': 'notification-group-invite',
    'group_join_request': 'notification-join-request',
    'group_event': 'notification-event',
    'follow_request': 'notification-follow',
    'message': 'notification-message'
  }
  return classes[type] || 'notification-default'
}

const hasActionButton = (type) => {
  return ['group_invite', 'group_join_request', 'follow_request'].includes(type)
}

const getActionButtonText = (type) => {
  const actions = {
    'group_invite': 'View Invite',
    'group_join_request': 'Review',
    'follow_request': 'Review'
  }
  return actions[type] || 'View'
}

const formatRelativeTime = (dateString) => {
  const date = new Date(dateString)
  const now = new Date()
  const diffInSeconds = Math.floor((now - date) / 1000)
  
  if (diffInSeconds < 60) {
    return 'Just now'
  } else if (diffInSeconds < 3600) {
    const minutes = Math.floor(diffInSeconds / 60)
    return `${minutes}m ago`
  } else if (diffInSeconds < 86400) {
    const hours = Math.floor(diffInSeconds / 3600)
    return `${hours}h ago`
  } else {
    const days = Math.floor(diffInSeconds / 86400)
    return `${days}d ago`
  }
}

// Lifecycle hooks
onMounted(async () => {
  // Initial load
  await refreshNotifications()
  
  // Set up real-time updates
  chatService.addMessageHandler((data) => {
    if (data.type === 'notification') {
      notificationStore.addNotification(data.notification)
    }
  })
  
  // Connect WebSocket if not already connected
  if (!chatService.connected) {
    chatService.connect()
  }
})

onUnmounted(() => {
  // Clean up if needed
})
</script>

<style scoped>
.unseen-notifications-container {
  @apply min-h-screen bg-gray-900 p-6;
}

.notification-header {
  @apply mb-6;
}

.unread-badge {
  @apply ml-2 bg-red-500 text-white text-sm px-2 py-1 rounded-full;
}

.header-actions {
  @apply flex flex-wrap gap-2;
}

.btn-primary {
  @apply bg-blue-600 hover:bg-blue-700 text-white px-4 py-2 rounded-lg transition-colors;
}

.btn-secondary {
  @apply bg-gray-600 hover:bg-gray-700 text-white px-4 py-2 rounded-lg transition-colors;
}

.loading-state {
  @apply flex flex-col items-center justify-center py-12;
}

.spinner {
  @apply w-8 h-8 border-4 border-blue-600 border-t-transparent rounded-full animate-spin;
}

.empty-state {
  @apply text-center py-12;
}

.empty-icon {
  @apply text-6xl mb-4;
}

.notifications-list {
  @apply space-y-4;
}

.notification-card {
  @apply bg-gray-800 border border-gray-700 rounded-lg p-4 flex items-start space-x-4 hover:bg-gray-750 transition-colors;
}

.notification-group-invite {
  @apply border-l-4 border-l-blue-500;
}

.notification-join-request {
  @apply border-l-4 border-l-green-500;
}

.notification-event {
  @apply border-l-4 border-l-purple-500;
}

.notification-follow {
  @apply border-l-4 border-l-yellow-500;
}

.notification-message {
  @apply border-l-4 border-l-pink-500;
}

.notification-default {
  @apply border-l-4 border-l-gray-500;
}

.notification-icon {
  @apply text-2xl flex-shrink-0;
}

.notification-content {
  @apply flex-1 min-w-0;
}

.notification-title {
  @apply text-white font-semibold mb-1 flex items-center;
}

.new-badge {
  @apply ml-2 bg-red-500 text-white text-xs px-2 py-1 rounded-full;
}

.notification-message {
  @apply text-gray-300 mb-2 break-words;
}

.notification-meta {
  @apply flex items-center justify-between text-xs text-gray-400;
}

.notification-type-label {
  @apply bg-gray-700 px-2 py-1 rounded;
}

.notification-actions {
  @apply flex flex-col space-y-2 flex-shrink-0;
}

.action-btn {
  @apply w-8 h-8 rounded-full flex items-center justify-center text-sm font-bold transition-colors;
}

.action-btn.mark-read {
  @apply bg-green-600 hover:bg-green-700 text-white;
}

.action-btn.primary {
  @apply bg-blue-600 hover:bg-blue-700 text-white w-auto px-3;
}

.action-btn.remove {
  @apply bg-red-600 hover:bg-red-700 text-white;
}

.load-more-container {
  @apply mt-6 text-center;
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .unseen-notifications-container {
    @apply p-4;
  }
  
  .notification-card {
    @apply flex-col space-x-0 space-y-3;
  }
  
  .notification-actions {
    @apply flex-row space-y-0 space-x-2 justify-end;
  }
}
</style>
