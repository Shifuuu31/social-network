<template>
  <div class="notifications-page">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">Notifications</h1>
        <div class="header-actions">
          <button 
            v-if="unreadCount > 0"
            @click="markAllAsRead"
            class="btn btn-secondary"
            :disabled="isMarkingAllRead"
          >
            {{ isMarkingAllRead ? 'Marking...' : 'Mark all as read' }}
          </button>
        </div>
      </div>

      <!-- Filters -->
      <div class="notification-filters">
        <button 
          v-for="filter in filters"
          :key="filter.value"
          @click="setActiveFilter(filter.value)"
          :class="['filter-btn', { active: activeFilter === filter.value }]"
        >
          {{ filter.label }}
          <span v-if="filter.value === 'unseen' && unreadCount > 0" class="filter-count">
            {{ unreadCount }}
          </span>
        </button>
      </div>

      <!-- Notifications Content -->
      <div class="notifications-content">
        <!-- Loading State -->
        <div v-if="notificationStore.isLoading && notifications.length === 0" class="loading-state">
          <div class="spinner"></div>
          <p>Loading notifications...</p>
        </div>

        <!-- Error State -->
        <div v-else-if="notificationStore.error" class="error-state">
          <div class="error-icon">⚠️</div>
          <h3>Error loading notifications</h3>
          <p>{{ notificationStore.error }}</p>
          <button @click="loadNotifications" class="btn btn-primary">Retry</button>
        </div>

        <!-- Empty State -->
        <div v-else-if="notifications.length === 0" class="empty-state">
          <div class="empty-icon">📭</div>
          <h3>No notifications found</h3>
          <p v-if="activeFilter === 'unseen'">
            You're all caught up! No unread notifications.
          </p>
          <p v-else-if="activeFilter === 'seen'">
            No read notifications to show.
          </p>
          <p v-else>
            You haven't received any notifications yet.
          </p>
        </div>

        <!-- Notifications List -->
        <div v-else class="notifications-grid">
          <div 
            v-for="notification in notifications" 
            :key="notification.id"
            class="notification-card"
            :class="{ 'unread': !notification.seen }"
          >
            <div class="notification-header">
              <div class="notification-icon">
                <span :style="{ color: getNotificationColor(notification.type) }">
                  {{ getNotificationIcon(notification.type) }}
                </span>
              </div>
              
              <div class="notification-meta">
                <span class="notification-type">
                  {{ getNotificationTypeLabel(notification.type) }}
                </span>
                <span class="notification-time">
                  {{ formatNotificationTime(notification.created_at) }}
                </span>
              </div>

              <div class="notification-actions">
                <button 
                  v-if="!notification.seen"
                  @click="markAsRead(notification.id)"
                  class="action-btn mark-read"
                  title="Mark as read"
                >
                  ✓
                </button>
                <!-- <button 
                  @click="deleteNotification(notification.id)"
                  class="action-btn delete"
                  title="Delete notification"
                >
                  🗑️
                </button> -->
              </div>
            </div>

            <div class="notification-body">
              <p class="notification-message">{{ notification.message }}</p>
              
              <!-- Action buttons for specific notification types -->
              <div v-if="hasActionButtons(notification)" class="notification-buttons">
                <template v-if="notification.type === 'group_invite'">
                  <button @click="handleGroupInviteAction(notification, 'accept')" class="btn btn-primary btn-sm">
                    Accept Invite
                  </button>
                  <button @click="handleGroupInviteAction(notification, 'decline')" class="btn btn-secondary btn-sm">
                    Decline
                  </button>
                </template>
                
                <template v-else-if="notification.type === 'follow_request'">
                  <button @click="handleFollowRequestAction(notification, 'accept')" class="btn btn-primary btn-sm">
                    Accept
                  </button>
                  <button @click="handleFollowRequestAction(notification, 'decline')" class="btn btn-secondary btn-sm">
                    Decline
                  </button>
                </template>
                
                <template v-else-if="notification.type === 'group_request'">
                  <button @click="handleGroupRequestAction(notification, 'accept')" class="btn btn-primary btn-sm">
                    Accept Request
                  </button>
                  <button @click="handleGroupRequestAction(notification, 'decline')" class="btn btn-secondary btn-sm">
                    Decline Request
                  </button>
                </template>
                
                <template v-else-if="notification.type === 'group_join_request'">
                  <button @click="handleGroupRequestAction(notification, 'accept')" class="btn btn-primary btn-sm">
                    Accept Request
                  </button>
                  <button @click="handleGroupRequestAction(notification, 'decline')" class="btn btn-secondary btn-sm">
                    Decline Request
                  </button>
                </template>
                
                <template v-else-if="notification.type === 'event_created'">
                  <button @click="viewEvent(notification)" class="btn btn-primary btn-sm">
                    View Event
                  </button>
                </template>
              </div>
            </div>
          </div>

          <!-- Load More Button -->
          <div v-if="canLoadMore" class="load-more-container">
            <button 
              @click="loadMoreNotifications"
              class="btn btn-secondary"
              :disabled="notificationStore.isLoading"
            >
              {{ notificationStore.isLoading ? 'Loading...' : 'Load more notifications' }}
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import { useNotificationStore } from '@/stores/notificationStore'
import { useRouter } from 'vue-router'

const notificationStore = useNotificationStore()
const router = useRouter()

// Local state
const activeFilter = ref('all')
const isMarkingAllRead = ref(false)
const currentPage = ref(0)
const pageSize = ref(20)

// Filter options
const filters = [
  { label: 'All', value: 'all' },
  { label: 'Unread', value: 'unseen' },
  { label: 'Read', value: 'seen' }
]

// Computed properties
const notifications = computed(() => notificationStore.notifications)
const unreadCount = computed(() => notificationStore.unreadCount)

const canLoadMore = computed(() => {
  // This would be better with pagination info from the API
  return notifications.value.length >= (currentPage.value + 1) * pageSize.value
})

// Methods
const setActiveFilter = (filter) => {
  activeFilter.value = filter
  currentPage.value = 0
  loadNotifications()
}

const loadNotifications = async () => {
  try {
    await notificationStore.fetchNotifications({
      start: 0,
      limit: pageSize.value,
      type: activeFilter.value
    })
  } catch (error) {
    console.error('Failed to load notifications:', error)
  }
}

const loadMoreNotifications = async () => {
  try {
    currentPage.value++
    await notificationStore.fetchNotifications({
      start: currentPage.value * pageSize.value,
      limit: pageSize.value,
      type: activeFilter.value
    })
  } catch (error) {
    console.error('Failed to load more notifications:', error)
    currentPage.value-- // Revert on error
  }
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

// const deleteNotification = async (notificationId) => {
//   if (confirm('Are you sure you want to delete this notification?')) {
//     try {
//       await notificationStore.deleteNotification(notificationId)
//     } catch (error) {
//       console.error('Failed to delete notification:', error)
//     }
//   }
// }

// Notification type helpers
const getNotificationTypeLabel = (type) => {
  const labels = {
    'follow_request': 'Follow Request',
    'group_invite': 'Group Invitation',
    'group_request': 'Group Join Request',
    'group_join_request': 'Group Join Request',
    'event_created': 'Event Created',
    'follow_accepted': 'Follow Accepted',
    'group_accepted': 'Group Accepted'
  }
  return labels[type] || 'Notification'
}

const hasActionButtons = (notification) => {
  const actionTypes = ['group_invite', 'follow_request', 'group_request', 'group_join_request', 'event_created']
  return actionTypes.includes(notification.type) && !notification.seen
}

// Action handlers
const handleGroupInviteAction = async (notification, action) => {
  try {
    // Call the notification action API
    await notificationStore.handleNotificationAction(notification, action)
    
    // Show success message
    console.log(`Successfully ${action}ed group invitation`)
    
    // Optionally show a toast notification or alert
    alert(`Successfully ${action}ed group invitation!`)
    
  } catch (error) {
    console.error(`Failed to ${action} group invitation:`, error)
    alert(`Failed to ${action} group invitation: ${error.message}`)
  }
}

const handleFollowRequestAction = async (notification, action) => {
  try {
    // Call the notification action API
    await notificationStore.handleNotificationAction(notification, action)
    
    // Show success message
    console.log(`Successfully ${action}ed follow request`)
    
    // Optionally show a toast notification or alert
    alert(`Successfully ${action}ed follow request!`)
    
  } catch (error) {
    console.error(`Failed to ${action} follow request:`, error)
    alert(`Failed to ${action} follow request: ${error.message}`)
  }
}

const handleGroupRequestAction = async (notification, action) => {
  try {
    // Call the notification action API
    await notificationStore.handleNotificationAction(notification, action)
    
    // Show success message
    console.log(`Successfully ${action}ed group join request`)
    
    // Optionally show a toast notification or alert
    alert(`Successfully ${action}ed group join request!`)
    
  } catch (error) {
    console.error(`Failed to ${action} group join request:`, error)
    alert(`Failed to ${action} group join request: ${error.message}`)
  }
}

const viewEvent = async (notification) => {
  // Mark as read first
  await markAsRead(notification.id)
  
  // Navigate to groups page to view events
  router.push('/groups')
}

// Helper methods from store
const getNotificationIcon = (type) => notificationStore.getNotificationIcon(type)
const getNotificationColor = (type) => notificationStore.getNotificationColor(type)
const formatNotificationTime = (timestamp) => notificationStore.formatNotificationTime(timestamp)

// Watch for filter changes
watch(activeFilter, () => {
  loadNotifications()
})

// Lifecycle
onMounted(() => {
  // Connect WebSocket if not already connected
  if (!notificationStore.isConnected) {
    notificationStore.connectWebSocket()
  }
  
  // Load initial notifications
  loadNotifications()
})
</script>

<style scoped>
.notifications-page {
  padding: 40px 20px;
}

.container {
  max-width: 800px;
  margin: 0 auto;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 30px;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: #fff;
  margin: 0;
}

.header-actions {
  display: flex;
  gap: 12px;
}

.notification-filters {
  display: flex;
  gap: 8px;
  margin-bottom: 30px;
  padding: 4px;
  background: rgba(255, 255, 255, 0.05);
  border-radius: 12px;
  width: fit-content;
}

.filter-btn {
  position: relative;
  background: none;
  border: none;
  color: #ccc;
  cursor: pointer;
  padding: 10px 16px;
  border-radius: 8px;
  transition: all 0.2s ease;
  font-weight: 500;
  display: flex;
  align-items: center;
  gap: 6px;
}

.filter-btn:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}

.filter-btn.active {
  background: #8b5cf6;
  color: #fff;
}

.filter-count {
  background: rgba(255, 255, 255, 0.2);
  color: #fff;
  font-size: 0.75rem;
  font-weight: 600;
  padding: 2px 6px;
  border-radius: 8px;
  line-height: 1;
  min-width: 16px;
  text-align: center;
}

.filter-btn.active .filter-count {
  background: rgba(255, 255, 255, 0.3);
}

.loading-state,
.error-state,
.empty-state {
  text-align: center;
  padding: 80px 20px;
  color: #ccc;
}

.error-icon,
.empty-icon {
  font-size: 4rem;
  margin-bottom: 20px;
  display: block;
}

.error-state h3,
.empty-state h3 {
  color: #fff;
  font-size: 1.5rem;
  margin: 0 0 10px 0;
}

.error-state p,
.empty-state p {
  margin-bottom: 20px;
  line-height: 1.5;
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
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.notifications-grid {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.notification-card {
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 12px;
  padding: 20px;
  transition: all 0.2s ease;
  position: relative;
}

.notification-card:hover {
  transform: translateY(-2px);
  border-color: #444;
}

.notification-card.unread {
  border-left: 4px solid #8b5cf6;
  background: rgba(139, 92, 246, 0.05);
}

.notification-header {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 12px;
}

.notification-icon {
  font-size: 1.5rem;
  flex-shrink: 0;
  margin-top: 2px;
}

.notification-meta {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.notification-type {
  color: #8b5cf6;
  font-weight: 600;
  font-size: 0.9rem;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.notification-time {
  color: #666;
  font-size: 0.85rem;
}

.notification-actions {
  display: flex;
  gap: 8px;
  opacity: 0;
  transition: opacity 0.2s ease;
}

.notification-card:hover .notification-actions {
  opacity: 1;
}

.action-btn {
  background: none;
  border: none;
  cursor: pointer;
  padding: 6px;
  border-radius: 6px;
  font-size: 0.9rem;
  transition: all 0.2s ease;
}

.action-btn.mark-read {
  color: #22c55e;
}

.action-btn.mark-read:hover {
  background: rgba(34, 197, 94, 0.1);
}

.action-btn.delete {
  color: #ef4444;
}

.action-btn.delete:hover {
  background: rgba(239, 68, 68, 0.1);
}

.notification-body {
  margin-left: 36px;
}

.notification-message {
  color: #fff;
  line-height: 1.5;
  margin: 0 0 16px 0;
}

.notification-buttons {
  display: flex;
  gap: 12px;
  flex-wrap: wrap;
}

.btn {
  padding: 8px 16px;
  border-radius: 8px;
  border: none;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 0.9rem;
}

.btn:disabled {
  opacity: 0.5;
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
  color: #ccc;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
}

.btn-sm {
  padding: 6px 12px;
  font-size: 0.8rem;
}

.load-more-container {
  text-align: center;
  padding: 20px;
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .page-header {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }

  .page-title {
    font-size: 2rem;
  }

  .notification-filters {
    width: 100%;
    justify-content: center;
  }

  .notification-card {
    padding: 16px;
  }

  .notification-body {
    margin-left: 0;
    margin-top: 12px;
  }

  .notification-buttons {
    gap: 8px;
  }

  .btn-sm {
    font-size: 0.75rem;
    padding: 5px 10px;
  }
}
</style>
