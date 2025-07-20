<template>
  <div class="fixed top-4 right-4 w-80 space-y-2 z-50">
    <div
      v-for="(notification, index) in notificationStore.notifications"
      :key="notification.id"
      class="bg-blue-100 text-blue-900 border border-blue-300 p-3 rounded-lg shadow-md transition-opacity duration-300"
    >
      <div class="flex justify-between items-start">
        <div class="flex-1">
          <strong>{{ notification.title || getNotificationTitle(notification.type) }}</strong>
          <p class="text-sm">{{ notification.message }}</p>
          <p class="text-xs text-gray-600 mt-1">{{ formatTime(notification.created_at) }}</p>
        </div>
        <div class="flex space-x-1 ml-2">
          <button 
            v-if="!notification.seen"
            @click="markAsRead(notification.id)"
            class="text-xs bg-blue-500 text-white px-2 py-1 rounded hover:bg-blue-600"
          >
            Mark Read
          </button>
          <button 
            @click="removeNotification(notification.id)"
            class="text-xs bg-red-500 text-white px-2 py-1 rounded hover:bg-red-600"
          >
            ×
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { onMounted } from 'vue';
import { useNotificationStore } from '@/stores/notificationStore';
import chatService from '@/services/chatService';

const notificationStore = useNotificationStore()

function getNotificationTitle(type) {
  const titles = {
    'group_invite': 'Group Invitation',
    'group_join_request': 'Join Request',
    'group_event': 'New Event',
    'follow_request': 'Follow Request',
    'notification': 'Notification'
  }
  return titles[type] || 'Notification'
}

function formatTime(timestamp) {
  if (!timestamp) return ''
  // Handle both ISO strings and Unix timestamps
  const date = typeof timestamp === 'string' ? new Date(timestamp) : new Date(timestamp * 1000)
  return date.toLocaleTimeString()
}

function markAsRead(notificationId) {
  notificationStore.markAsRead(notificationId)
}

function removeNotification(notificationId) {
  notificationStore.removeNotification(notificationId)
}

onMounted(() => {
  // Connect to WebSocket for real-time notifications
  chatService.connect()
  
  // Fetch existing notifications
  notificationStore.fetchNotifications()

  // Listen for specific notification events
  chatService.onMessage((msg) => {
    if (msg.group_id && msg.type === 'group') {
      notificationStore.addNotification({
        title: `New Group Message`,
        message: `Group ${msg.group_id}: ${msg.content}`,
        type: 'group_message',
        seen: false
      });
    }
  });

  chatService.onConnectionChange((connected) => {
    notificationStore.addNotification({
      title: 'Connection Status',
      message: connected ? 'Connected to chat server' : 'Disconnected from chat server',
      type: 'connection',
      seen: false
    });
  });
});
</script>

<style scoped>
.notifications {
  position: fixed;
  top: 20px;
  right: 20px;
  width: 300px;
  max-height: 400px;
  overflow-y: auto;
}
.notification {
  background: #f0f4ff;
  color: #1e40af;
  border: 1px solid #cbd5e1;
  padding: 12px;
  border-radius: 8px;
  margin-bottom: 10px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: opacity 0.3s ease;
}
.notification:hover {
  opacity: 0.9; /* Slightly fade on hover */
}
</style>
