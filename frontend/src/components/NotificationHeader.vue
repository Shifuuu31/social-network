<template>
  <div class="notification-header">
    <!-- Notification Bell with Count -->
    <router-link to="/notifications" class="notification-bell" :class="{ 'has-unseen': unseenCount > 0 }">
      <span class="bell-icon">🔔</span>
      <span v-if="unseenCount > 0" class="notification-count">{{ displayCount }}</span>
    </router-link>
    
    <!-- Quick Access Dropdown (Optional) -->
    <div v-if="showQuickAccess" class="notification-dropdown" @click.stop>
      <div class="dropdown-header">
        <h4>Recent Notifications</h4>
        <router-link to="/notifications/unseen" class="view-all-link">View All Unseen</router-link>
      </div>
      <div class="dropdown-content">
        <div v-if="recentNotifications.length === 0" class="no-notifications">
          No recent notifications
        </div>
        <div v-else>
          <div 
            v-for="notification in recentNotifications.slice(0, 3)" 
            :key="notification.id"
            class="quick-notification"
          >
            <div class="notification-icon">{{ getNotificationIcon(notification.type) }}</div>
            <div class="notification-text">
              <div class="notification-title">{{ getNotificationTitle(notification.type) }}</div>
              <div class="notification-message">{{ notification.message }}</div>
            </div>
            <div v-if="!notification.seen" class="unseen-dot"></div>
          </div>
        </div>
        <div class="dropdown-footer">
          <router-link to="/notifications" class="view-all-btn">View All Notifications</router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, onUnmounted } from 'vue'
import { useNotificationStore } from '@/stores/notificationStore'
import chatService from '@/services/chatService'

const notificationStore = useNotificationStore()
const showQuickAccess = ref(false)

// Computed properties
const unseenCount = computed(() => {
  return notificationStore.notifications.filter(n => !n.seen).length
})

const displayCount = computed(() => {
  return unseenCount.value > 99 ? '99+' : unseenCount.value.toString()
})

const recentNotifications = computed(() => {
  return notificationStore.notifications
    .sort((a, b) => new Date(b.created_at) - new Date(a.created_at))
    .slice(0, 5)
})

// Methods
const toggleQuickAccess = () => {
  showQuickAccess.value = !showQuickAccess.value
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

// Lifecycle hooks
onMounted(async () => {
  // Load initial notifications
  try {
    await notificationStore.fetchNotifications()
  } catch (error) {
    console.error('Failed to load notifications in header:', error)
  }
  
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
.notification-header {
  position: relative;
  display: flex;
  align-items: center;
}

.notification-bell {
  position: relative;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 40px;
  height: 40px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.1);
  text-decoration: none;
  transition: all 0.2s ease;
  cursor: pointer;
}

.notification-bell:hover {
  background: rgba(255, 255, 255, 0.2);
  transform: scale(1.05);
}

.notification-bell.has-unseen {
  background: rgba(239, 68, 68, 0.2);
}

.notification-bell.has-unseen:hover {
  background: rgba(239, 68, 68, 0.3);
}

.bell-icon {
  font-size: 18px;
  color: #ccc;
}

.notification-bell.has-unseen .bell-icon {
  color: #ef4444;
  animation: ring 2s ease-in-out infinite;
}

.notification-count {
  position: absolute;
  top: -2px;
  right: -2px;
  background: #ef4444;
  color: white;
  font-size: 10px;
  font-weight: bold;
  padding: 2px 6px;
  border-radius: 10px;
  min-width: 16px;
  text-align: center;
  line-height: 1;
}

.notification-dropdown {
  position: absolute;
  top: 100%;
  right: 0;
  width: 320px;
  background: #1a1a1a;
  border: 1px solid #333;
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0, 0, 0, 0.4);
  z-index: 1000;
  margin-top: 8px;
}

.dropdown-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 16px;
  border-bottom: 1px solid #333;
}

.dropdown-header h4 {
  color: #fff;
  margin: 0;
  font-size: 14px;
  font-weight: 600;
}

.view-all-link {
  color: #8b5cf6;
  text-decoration: none;
  font-size: 12px;
  font-weight: 500;
}

.view-all-link:hover {
  text-decoration: underline;
}

.dropdown-content {
  max-height: 300px;
  overflow-y: auto;
}

.no-notifications {
  padding: 20px;
  text-align: center;
  color: #666;
  font-size: 14px;
}

.quick-notification {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  padding: 12px 16px;
  border-bottom: 1px solid #2a2a2a;
  transition: background-color 0.2s ease;
  position: relative;
}

.quick-notification:hover {
  background: rgba(255, 255, 255, 0.05);
}

.quick-notification:last-child {
  border-bottom: none;
}

.notification-icon {
  font-size: 16px;
  flex-shrink: 0;
  margin-top: 2px;
}

.notification-text {
  flex: 1;
  min-width: 0;
}

.notification-title {
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  margin-bottom: 2px;
}

.notification-message {
  color: #ccc;
  font-size: 11px;
  line-height: 1.3;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.unseen-dot {
  width: 6px;
  height: 6px;
  background: #ef4444;
  border-radius: 50%;
  flex-shrink: 0;
  margin-top: 4px;
}

.dropdown-footer {
  padding: 12px 16px;
  border-top: 1px solid #333;
}

.view-all-btn {
  display: block;
  width: 100%;
  padding: 8px;
  background: rgba(139, 92, 246, 0.1);
  color: #8b5cf6;
  text-decoration: none;
  text-align: center;
  border-radius: 6px;
  font-size: 12px;
  font-weight: 500;
  transition: background-color 0.2s ease;
}

.view-all-btn:hover {
  background: rgba(139, 92, 246, 0.2);
}

@keyframes ring {
  0%, 100% { transform: rotate(0deg); }
  10%, 30% { transform: rotate(-10deg); }
  20% { transform: rotate(10deg); }
}

/* Mobile responsiveness */
@media (max-width: 768px) {
  .notification-dropdown {
    width: 280px;
    right: -20px;
  }
}
</style>