<template>
  <div class="notifications-page">
    <!-- Header -->
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">
          <span class="icon">🔔</span>
          Notifications
        </h1>
        <div class="header-actions">
          <button 
            v-if="unreadNotifications.length > 0"
            @click="markAllAsRead" 
            class="btn btn-outline"
          >
            Mark All Read
          </button>
          <button 
            @click="refreshNotifications" 
            class="btn btn-primary"
            :disabled="isLoading"
          >
            <span v-if="isLoading">🔄</span>
            <span v-else>↻</span>
            Refresh
          </button>
        </div>
      </div>
    </div>

    <!-- Notification Stats -->
    <div class="notification-stats">
      <div class="stat-card">
        <div class="stat-number">{{ totalNotifications }}</div>
        <div class="stat-label">Total</div>
      </div>
      <div class="stat-card highlight">
        <div class="stat-number">{{ unreadCount }}</div>
        <div class="stat-label">Unread</div>
      </div>
      <div class="stat-card">
        <div class="stat-number">{{ readCount }}</div>
        <div class="stat-label">Read</div>
      </div>
    </div>

    <!-- Filter Tabs -->
    <div class="filter-tabs">
      <button 
        :class="['tab', { active: activeFilter === 'all' }]"
        @click="setFilter('all')"
      >
        All ({{ totalNotifications }})
      </button>
      <button 
        :class="['tab', { active: activeFilter === 'unread' }]"
        @click="setFilter('unread')"
      >
        Unread ({{ unreadCount }})
      </button>
    </div>

    <!-- Loading State -->
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Loading notifications...</p>
    </div>

    <!-- Empty State -->
    <div v-else-if="filteredNotifications.length === 0" class="empty-state">
      <div class="empty-icon">📭</div>
      <h3>{{ getEmptyStateTitle() }}</h3>
      <p>{{ getEmptyStateMessage() }}</p>
    </div>

    <!-- Notifications List -->
    <div v-else class="notifications-container">
      <div 
        v-for="notification in filteredNotifications" 
        :key="notification.id"
        :class="['notification-card', { unread: !notification.seen }]"
      >
        <!-- Notification Header -->
        <div class="notification-header">
          <div class="notification-type">
            <span class="type-icon">{{ getNotificationIcon(notification.type) }}</span>
            <span class="type-label">{{ getNotificationTitle(notification.type) }}</span>
          </div>
          <div class="notification-time">
            {{ formatTimeAgo(notification.created_at) }}
          </div>
        </div>

        <!-- Notification Content -->
        <div class="notification-content">
          <p class="notification-message">{{ notification.message }}</p>
        </div>

        <!-- Notification Actions -->
        <div class="notification-actions">
          <button 
            v-if="!notification.seen"
            @click="markAsRead(notification.id)"
            class="action-btn mark-read"
          >
            <span class="icon">✓</span>
            Mark as Read
          </button>
          <button 
            @click="removeNotification(notification.id)"
            class="action-btn remove"
          >
            <span class="icon">🗑️</span>
            Remove
          </button>
        </div>

        <!-- Unread Indicator -->
        <div v-if="!notification.seen" class="unread-indicator"></div>
      </div>
    </div>

    <!-- Load More Button -->
    <div v-if="hasMoreNotifications" class="load-more-section">
      <button @click="loadMoreNotifications" class="btn btn-outline load-more-btn">
        Load More Notifications
      </button>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useNotificationStore } from '@/stores/notificationStore'
import chatService from '@/services/chatService'

const notificationStore = useNotificationStore()

// Reactive state
const isLoading = ref(false)
const activeFilter = ref('unread') // 'all', 'unread', 'read'
const hasMoreNotifications = ref(false)

// Computed properties
const totalNotifications = computed(() => notificationStore.notifications.length)
const unreadNotifications = computed(() => 
  notificationStore.notifications.filter(n => !n.seen)
)
const readNotifications = computed(() => 
  notificationStore.notifications.filter(n => n.seen)
)
const unreadCount = computed(() => unreadNotifications.value.length)
const readCount = computed(() => readNotifications.value.length)

const filteredNotifications = computed(() => {
  switch (activeFilter.value) {
    case 'unread':
      return unreadNotifications.value
    case 'read':
      return readNotifications.value
    default:
      return notificationStore.notifications
  }
})

// Methods
function getNotificationIcon(type) {
  const icons = {
    'group_invite': '👥',
    'group_join_request': '🙋',
    'group_event': '📅',
    'follow_request': '👤',
    'notification': '🔔',
    'group_message': '💬',
    'connection': '🔗'
  }
  return icons[type] || '🔔'
}

function getNotificationTitle(type) {
  const titles = {
    'group_invite': 'Group Invitation',
    'group_join_request': 'Join Request',
    'group_event': 'New Event',
    'follow_request': 'Follow Request',
    'notification': 'Notification',
    'group_message': 'Group Message',
    'connection': 'Connection Update'
  }
  return titles[type] || 'Notification'
}

function formatTimeAgo(timestamp) {
  if (!timestamp) return 'Just now'
  
  const date = new Date(timestamp)
  const now = new Date()
  const diffInSeconds = Math.floor((now - date) / 1000)
  
  if (diffInSeconds < 60) return 'Just now'
  if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m ago`
  if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h ago`
  if (diffInSeconds < 604800) return `${Math.floor(diffInSeconds / 86400)}d ago`
  
  return date.toLocaleDateString()
}

function getEmptyStateTitle() {
  switch (activeFilter.value) {
    case 'unread':
      return 'No unread notifications'
    case 'read':
      return 'No read notifications'
    default:
      return 'No notifications'
  }
}

function getEmptyStateMessage() {
  switch (activeFilter.value) {
    case 'unread':
      return 'All caught up! You have no unread notifications.'
    default:
      return 'You have no notifications yet. They will appear here when you receive them.'
  }
}

function setFilter(filter) {
  activeFilter.value = filter
}

async function markAsRead(notificationId) {
  try {
    await notificationStore.markAsRead(notificationId)
  } catch (error) {
    console.error('Failed to mark notification as read:', error)
  }
}

async function markAllAsRead() {
  try {
    const unreadIds = unreadNotifications.value.map(n => n.id)
    if (unreadIds.length > 0) {
      await notificationStore.markMultipleAsRead(unreadIds)
    }
  } catch (error) {
    console.error('Failed to mark all notifications as read:', error)
  }
}

async function removeNotification(notificationId) {
  try {
    await notificationStore.removeNotification(notificationId)
  } catch (error) {
    console.error('Failed to remove notification:', error)
  }
}

async function refreshNotifications() {
  isLoading.value = true
  try {
    await notificationStore.fetchNotifications()
  } catch (error) {
    console.error('Failed to refresh notifications:', error)
  } finally {
    isLoading.value = false
  }
}

async function loadMoreNotifications() {
  // Implementation for pagination if needed
  console.log('Loading more notifications...')
}

// Initialize component
onMounted(async () => {
  isLoading.value = true
  
  try {
    // Connect to WebSocket for real-time updates
    chatService.connect()
    
    // Fetch initial notifications
    await notificationStore.fetchNotifications()
    
    // Set up real-time notification listener
    chatService.onMessage((message) => {
      if (message.type === 'notification') {
        notificationStore.addNotification(message.data)
      }
    })
    
  } catch (error) {
    console.error('Failed to initialize notifications:', error)
  } finally {
    isLoading.value = false
  }
})

// Auto-switch to 'all' if unread becomes empty
watch(unreadCount, (newCount) => {
  if (newCount === 0 && activeFilter.value === 'unread') {
    activeFilter.value = 'all'
  }
})
</script>

<style scoped>
.notifications-page {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
  background: #f8f9fa;
  min-height: 100vh;
}

/* Header */
.page-header {
  background: white;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.header-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.page-title {
  font-size: 28px;
  font-weight: 700;
  color: #1a202c;
  margin: 0;
  display: flex;
  align-items: center;
  gap: 12px;
}

.page-title .icon {
  font-size: 32px;
}

.header-actions {
  display: flex;
  gap: 12px;
}

/* Buttons */
.btn {
  padding: 10px 16px;
  border-radius: 8px;
  font-weight: 500;
  border: none;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;
}

.btn-primary {
  background: #3b82f6;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background: #2563eb;
}

.btn-outline {
  background: transparent;
  color: #3b82f6;
  border: 2px solid #3b82f6;
}

.btn-outline:hover {
  background: #3b82f6;
  color: white;
}

.btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

/* Stats */
.notification-stats {
  display: grid;
  grid-template-columns: repeat(3, 1fr);
  gap: 16px;
  margin-bottom: 20px;
}

.stat-card {
  background: white;
  padding: 20px;
  border-radius: 12px;
  text-align: center;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.stat-card.highlight {
  background: linear-gradient(135deg, #3b82f6, #1d4ed8);
  color: white;
}

.stat-number {
  font-size: 32px;
  font-weight: 700;
  margin-bottom: 4px;
}

.stat-label {
  font-size: 14px;
  opacity: 0.8;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

/* Filter Tabs */
.filter-tabs {
  display: flex;
  background: white;
  border-radius: 12px;
  padding: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.tab {
  flex: 1;
  padding: 12px 16px;
  border: none;
  background: transparent;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-weight: 500;
  color: #6b7280;
}

.tab.active {
  background: #3b82f6;
  color: white;
}

.tab:hover:not(.active) {
  background: #f3f4f6;
}

/* Loading State */
.loading-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid #f3f4f6;
  border-top: 4px solid #3b82f6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Empty State */
.empty-state {
  text-align: center;
  padding: 60px 20px;
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
}

.empty-icon {
  font-size: 64px;
  margin-bottom: 16px;
}

.empty-state h3 {
  font-size: 20px;
  margin-bottom: 8px;
  color: #1a202c;
}

.empty-state p {
  color: #6b7280;
  margin: 0;
}

/* Notifications Container */
.notifications-container {
  space-y: 12px;
}

.notification-card {
  background: white;
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  position: relative;
  transition: all 0.2s ease;
}

.notification-card:hover {
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
}

.notification-card.unread {
  border-left: 4px solid #3b82f6;
  background: linear-gradient(90deg, #eff6ff 0%, white 20%);
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.notification-type {
  display: flex;
  align-items: center;
  gap: 8px;
}

.type-icon {
  font-size: 20px;
}

.type-label {
  font-weight: 600;
  color: #1a202c;
}

.notification-time {
  font-size: 12px;
  color: #6b7280;
}

.notification-content {
  margin-bottom: 16px;
}

.notification-message {
  color: #4b5563;
  line-height: 1.5;
  margin: 0;
}

.notification-actions {
  display: flex;
  gap: 8px;
}

.action-btn {
  padding: 6px 12px;
  border: none;
  border-radius: 6px;
  font-size: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 4px;
}

.action-btn.mark-read {
  background: #10b981;
  color: white;
}

.action-btn.mark-read:hover {
  background: #059669;
}

.action-btn.remove {
  background: #ef4444;
  color: white;
}

.action-btn.remove:hover {
  background: #dc2626;
}

.unread-indicator {
  position: absolute;
  top: 20px;
  right: 20px;
  width: 8px;
  height: 8px;
  background: #3b82f6;
  border-radius: 50%;
}

/* Load More */
.load-more-section {
  text-align: center;
  margin-top: 24px;
}

.load-more-btn {
  padding: 12px 24px;
}

/* Responsive Design */
@media (max-width: 768px) {
  .notifications-page {
    padding: 16px;
  }
  
  .header-content {
    flex-direction: column;
    gap: 16px;
    text-align: center;
  }
  
  .notification-stats {
    grid-template-columns: 1fr;
  }
  
  .filter-tabs {
    flex-direction: column;
  }
  
  .notification-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 8px;
  }
  
  .notification-actions {
    flex-wrap: wrap;
  }
}
</style>
