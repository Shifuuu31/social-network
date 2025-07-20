<template>
  <div class="notifications-dropdown">
    <!-- Notification Bell Button -->
    <button 
      class="notification-bell" 
      @click="toggleDropdown"
      :class="{ 'has-unread': unreadCount > 0 }"
    >
      <span class="bell-icon">🔔</span>
      <span v-if="unreadCount > 0" class="notification-badge">
        {{ unreadCount > 99 ? '99+' : unreadCount }}
      </span>
    </button>

    <!-- Dropdown Menu -->
    <div v-if="isDropdownOpen" class="dropdown-menu" @click.stop>
      <div class="dropdown-header">
        <h3>Notifications</h3>
        <div class="header-actions">
          <button 
            v-if="unreadCount > 0" 
            @click="markAllAsRead"
            class="mark-all-btn"
            :disabled="isMarkingAllRead"
          >
            {{ isMarkingAllRead ? 'Marking...' : 'Mark all read' }}
          </button>
        </div>
      </div>

      <div class="notifications-list">
        <!-- Loading State -->
        <div v-if="notificationStore.isLoading && notifications.length === 0" class="loading-state">
          <div class="spinner"></div>
          <p>Loading notifications...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="notificationStore.error" class="error-state">
          <p>{{ notificationStore.error }}</p>
          <button @click="loadNotifications" class="retry-btn">Retry</button>
        </div>

        <!-- Empty State -->
        <div v-else-if="notifications.length === 0" class="empty-state">
          <span class="empty-icon">📭</span>
          <p>No notifications yet</p>
        </div>

        <!-- Notifications List -->
        <div v-else class="notifications-container">
          <div 
            v-for="notification in displayedNotifications" 
            :key="notification.id"
            class="notification-item"
            :class="{ 'unread': !notification.seen }"
            @click="handleNotificationClick(notification)"
          >
            <div class="notification-icon">
              <span :style="{ color: getNotificationColor(notification.type) }">
                {{ getNotificationIcon(notification.type) }}
              </span>
            </div>
            
            <div class="notification-content">
              <p class="notification-message">{{ notification.message }}</p>
              <span class="notification-time">
                {{ formatNotificationTime(notification.created_at) }}
              </span>
            </div>

            <div class="notification-actions">
              <button 
                v-if="!notification.seen"
                @click.stop="markAsRead(notification.id)"
                class="mark-read-btn"
                title="Mark as read"
              >
                ✓
              </button>
              <button 
                @click.stop="deleteNotification(notification.id)"
                class="delete-btn"
                title="Delete notification"
              >
                ✕
              </button>
            </div>
          </div>

          <!-- Load More Button -->
          <button 
            v-if="canLoadMore"
            @click="loadMoreNotifications"
            class="load-more-btn"
            :disabled="notificationStore.isLoading"
          >
            {{ notificationStore.isLoading ? 'Loading...' : 'Load more' }}
          </button>
        </div>
      </div>

      <!-- Footer -->
      <div class="dropdown-footer">
        <router-link to="/notifications" class="view-all-btn" @click="closeDropdown">
          View all notifications
        </router-link>
      </div>
    </div>

    <!-- Backdrop -->
    <div v-if="isDropdownOpen" class="dropdown-backdrop" @click="closeDropdown"></div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useNotificationStore } from '@/stores/notificationStore'

const notificationStore = useNotificationStore()

// Local state
const isDropdownOpen = ref(false)
const isMarkingAllRead = ref(false)
const displayLimit = ref(10)

// Computed properties
const notifications = computed(() => notificationStore.notifications)
const unreadCount = computed(() => notificationStore.unreadCount)

const displayedNotifications = computed(() => 
  notifications.value.slice(0, displayLimit.value)
)

const canLoadMore = computed(() => 
  notifications.value.length > displayLimit.value
)

// Methods
const toggleDropdown = () => {
  isDropdownOpen.value = !isDropdownOpen.value
  
  if (isDropdownOpen.value) {
    loadNotifications()
  }
}

const closeDropdown = () => {
  isDropdownOpen.value = false
}

const loadNotifications = async () => {
  try {
    await notificationStore.fetchNotifications({ start: 0, limit: 20 })
  } catch (error) {
    console.error('Failed to load notifications:', error)
  }
}

const loadMoreNotifications = () => {
  displayLimit.value += 10
}

const markAsRead = async (notificationId) => {
  try {
    await notificationStore.markAsSeen(notificationId)
  } catch (error) {
    console.error('Failed to mark notification as read:', error)
  }
}

const markAllAsRead = async () => {
  isMarkingAllRead.value = true
  try {
    await notificationStore.markAllAsSeen()
  } catch (error) {
    console.error('Failed to mark all notifications as read:', error)
  } finally {
    isMarkingAllRead.value = false
  }
}

const deleteNotification = async (notificationId) => {
  try {
    await notificationStore.deleteNotification(notificationId)
  } catch (error) {
    console.error('Failed to delete notification:', error)
  }
}

const handleNotificationClick = (notification) => {
  // Mark as read if unread
  if (!notification.seen) {
    markAsRead(notification.id)
  }

  // Handle notification-specific actions
  handleNotificationAction(notification)
}

const handleNotificationAction = (notification) => {
  // Close dropdown first
  closeDropdown()

  // Navigate based on notification type
  switch (notification.type) {
    case 'group_invite':
    case 'group_request':
      // Extract group ID from message or use a separate field if available
      // For now, navigate to groups page
      window.location.href = '/groups'
      break
    
    case 'follow_request':
      // Navigate to profile or followers page
      // window.location.href = '/profile'
      break
    
    case 'event_created':
      // Navigate to groups page to see events
      window.location.href = '/groups'
      break
    
    default:
      // Default action - just mark as read
      break
  }
}

// Helper methods from store
const getNotificationIcon = (type) => notificationStore.getNotificationIcon(type)
const getNotificationColor = (type) => notificationStore.getNotificationColor(type)
const formatNotificationTime = (timestamp) => notificationStore.formatNotificationTime(timestamp)

// Handle click outside to close dropdown
const handleClickOutside = (event) => {
  if (!event.target.closest('.notifications-dropdown')) {
    closeDropdown()
  }
}

// Lifecycle
onMounted(() => {
  document.addEventListener('click', handleClickOutside)
  
  // Connect to WebSocket for real-time updates
  notificationStore.connectWebSocket()
  
  // Request browser notification permission
  notificationStore.requestNotificationPermission()
  
  // Initial load of unread count
  notificationStore.fetchUnreadCount()
})

onUnmounted(() => {
  document.removeEventListener('click', handleClickOutside)
})
</script>

<style scoped>
.notifications-dropdown {
  position: relative;
  display: inline-block;
}

.notification-bell {
  position: relative;
  background: none;
  border: none;
  color: #ccc;
  cursor: pointer;
  padding: 8px;
  border-radius: 8px;
  transition: all 0.2s ease;
  font-size: 1.2rem;
}

.notification-bell:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}

.notification-bell.has-unread {
  color: #f59e0b;
}

.notification-bell.has-unread:hover {
  color: #fbbf24;
}

.bell-icon {
  display: block;
}

.notification-badge {
  position: absolute;
  top: 2px;
  right: 2px;
  background: #ef4444;
  color: #fff;
  font-size: 0.7rem;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 10px;
  line-height: 1;
  min-width: 18px;
  text-align: center;
}

.dropdown-menu {
  position: absolute;
  top: calc(100% + 10px);
  right: 0;
  width: 380px;
  max-height: 500px;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 12px;
  box-shadow: 0 10px 30px rgba(0, 0, 0, 0.5);
  z-index: 1000;
  overflow: hidden;
}

.dropdown-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 999;
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px 20px;
  border-bottom: 1px solid #333;
}

.dropdown-header h3 {
  margin: 0;
  color: #fff;
  font-size: 1.1rem;
  font-weight: 600;
}

.mark-all-btn {
  background: none;
  border: none;
  color: #8b5cf6;
  cursor: pointer;
  font-size: 0.85rem;
  font-weight: 500;
  padding: 4px 8px;
  border-radius: 6px;
  transition: all 0.2s ease;
}

.mark-all-btn:hover:not(:disabled) {
  background: rgba(139, 92, 246, 0.1);
}

.mark-all-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.notifications-list {
  max-height: 350px;
  overflow-y: auto;
}

.loading-state,
.error-state,
.empty-state {
  padding: 40px 20px;
  text-align: center;
  color: #ccc;
}

.spinner {
  width: 24px;
  height: 24px;
  border: 2px solid rgba(255, 255, 255, 0.2);
  border-top: 2px solid #8b5cf6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 10px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.empty-icon {
  font-size: 2rem;
  display: block;
  margin-bottom: 10px;
}

.retry-btn {
  background: #8b5cf6;
  color: #fff;
  border: none;
  padding: 8px 16px;
  border-radius: 6px;
  cursor: pointer;
  font-size: 0.85rem;
  margin-top: 10px;
}

.notifications-container {
  padding: 8px 0;
}

.notification-item {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 20px;
  cursor: pointer;
  transition: background 0.2s ease;
  border-left: 3px solid transparent;
}

.notification-item:hover {
  background: rgba(255, 255, 255, 0.05);
}

.notification-item.unread {
  background: rgba(139, 92, 246, 0.05);
  border-left-color: #8b5cf6;
}

.notification-icon {
  font-size: 1.2rem;
  margin-top: 2px;
  flex-shrink: 0;
}

.notification-content {
  flex: 1;
  min-width: 0;
}

.notification-message {
  color: #fff;
  font-size: 0.9rem;
  line-height: 1.4;
  margin: 0 0 4px 0;
  word-wrap: break-word;
}

.notification-time {
  color: #666;
  font-size: 0.8rem;
}

.notification-actions {
  display: flex;
  gap: 4px;
  flex-shrink: 0;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.notification-item:hover .notification-actions {
  opacity: 1;
}

.mark-read-btn,
.delete-btn {
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  font-size: 0.8rem;
  transition: all 0.2s ease;
}

.mark-read-btn:hover {
  color: #22c55e;
  background: rgba(34, 197, 94, 0.1);
}

.delete-btn:hover {
  color: #ef4444;
  background: rgba(239, 68, 68, 0.1);
}

.load-more-btn {
  width: calc(100% - 40px);
  margin: 10px 20px;
  background: rgba(255, 255, 255, 0.1);
  color: #ccc;
  border: none;
  padding: 10px;
  border-radius: 8px;
  cursor: pointer;
  font-size: 0.85rem;
  transition: all 0.2s ease;
}

.load-more-btn:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
}

.load-more-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.dropdown-footer {
  padding: 12px 20px;
  border-top: 1px solid #333;
  text-align: center;
}

.view-all-btn {
  color: #8b5cf6;
  text-decoration: none;
  font-size: 0.9rem;
  font-weight: 500;
  padding: 8px 16px;
  border-radius: 6px;
  transition: all 0.2s ease;
  display: inline-block;
}

.view-all-btn:hover {
  background: rgba(139, 92, 246, 0.1);
}

/* Scrollbar styling */
.notifications-list::-webkit-scrollbar {
  width: 6px;
}

.notifications-list::-webkit-scrollbar-track {
  background: transparent;
}

.notifications-list::-webkit-scrollbar-thumb {
  background: #333;
  border-radius: 3px;
}

.notifications-list::-webkit-scrollbar-thumb:hover {
  background: #444;
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .dropdown-menu {
    width: 320px;
    right: -20px;
  }
  
  .notification-message {
    font-size: 0.85rem;
  }
}
</style>