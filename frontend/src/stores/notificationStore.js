import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useNotificationStore = defineStore('notifications', () => {
  // State
  const notifications = ref([])
  const unreadCount = ref(0)
  const isLoading = ref(false)
  const error = ref(null)
  const socket = ref(null)
  const isConnected = ref(false)

  const API_BASE = '/api'

  // Getters
  const unreadNotifications = computed(() => 
    notifications.value.filter(n => !n.seen)
  )

  const recentNotifications = computed(() => 
    notifications.value.slice(0, 5)
  )

  // WebSocket connection
  const connectWebSocket = () => {
    try {
      const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:'
      const wsUrl = `${protocol}//${window.location.host}/connect`
      
      socket.value = new WebSocket(wsUrl)

      socket.value.onopen = () => {
        isConnected.value = true
        console.log('📡 Notification WebSocket connected')
      }

      socket.value.onmessage = (event) => {
        try {
          const data = JSON.parse(event.data)
          
          if (data.type === 'notification') {
            handleWebSocketNotification(data)
          }
        } catch (err) {
          console.error('Error parsing WebSocket message:', err)
        }
      }

      socket.value.onclose = () => {
        isConnected.value = false
        console.log('📡 Notification WebSocket disconnected')
        
        // Reconnect after 3 seconds
        setTimeout(() => {
          if (!isConnected.value) {
            connectWebSocket()
          }
        }, 3000)
      }

      socket.value.onerror = (error) => {
        console.error('WebSocket error:', error)
        isConnected.value = false
      }
    } catch (err) {
      console.error('Failed to connect to WebSocket:', err)
    }
  }

  // Handle WebSocket notification updates
  const handleWebSocketNotification = (data) => {
    const { action, notification, unread_count } = data

    // Update unread count
    unreadCount.value = unread_count || 0

    switch (action) {
      case 'new':
        if (notification) {
          // Add new notification to the beginning of the list
          notifications.value.unshift(notification)
          
          // Show browser notification if permission granted
          showBrowserNotification(notification)
        }
        break
      
      case 'update':
        // Refresh notifications or just update count
        if (notification) {
          const index = notifications.value.findIndex(n => n.id === notification.id)
          if (index !== -1) {
            notifications.value[index] = notification
          }
        }
        break
      
      case 'delete':
        if (notification) {
          notifications.value = notifications.value.filter(n => n.id !== notification.id)
        }
        break
    }
  }

  // Show browser notification
  const showBrowserNotification = (notification) => {
    if ('Notification' in window && Notification.permission === 'granted') {
      const notif = new Notification('Social Network', {
        body: notification.message,
        icon: '/favicon.ico',
        tag: `notification-${notification.id}`
      })

      // Auto close after 5 seconds
      setTimeout(() => notif.close(), 5000)
    }
  }

  // Request browser notification permission
  const requestNotificationPermission = async () => {
    if ('Notification' in window && Notification.permission === 'default') {
      const permission = await Notification.requestPermission()
      return permission === 'granted'
    }
    return Notification.permission === 'granted'
  }

  // Fetch notifications from API
  const fetchNotifications = async (params = {}) => {
    isLoading.value = true
    error.value = null

    try {
      const payload = {
        start: params.start || 0,
        n_items: params.limit || 20,
        type: params.type || 'all'
      }

      const response = await fetch(`${API_BASE}/notifications/fetch`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      
      if (params.start === 0) {
        // Replace notifications if starting from beginning
        notifications.value = data.notifications || []
      } else {
        // Append if loading more
        notifications.value.push(...(data.notifications || []))
      }
      
      unreadCount.value = data.unread_count || 0

      return data
    } catch (err) {
      error.value = err.message
      console.error('Error fetching notifications:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  // Get unread count
  const fetchUnreadCount = async () => {
    try {
      const response = await fetch(`${API_BASE}/notifications/unread-count`)
      
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      unreadCount.value = data.unread_count || 0
      
      return data.unread_count
    } catch (err) {
      console.error('Error fetching unread count:', err)
      return 0
    }
  }

  // Mark notification as seen
  const markAsSeen = async (notificationId) => {
    try {
      const response = await fetch(`${API_BASE}/notifications/mark-seen`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ notification_id: notificationId })
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      
      // Update local state
      const notification = notifications.value.find(n => n.id === notificationId)
      if (notification) {
        notification.seen = true
      }
      
      unreadCount.value = data.unread_count || 0

      return data
    } catch (err) {
      console.error('Error marking notification as seen:', err)
      throw err
    }
  }

  // Mark all notifications as seen
  const markAllAsSeen = async () => {
    try {
      const response = await fetch(`${API_BASE}/notifications/mark-all-seen`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' }
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      
      // Update local state
      notifications.value.forEach(n => n.seen = true)
      unreadCount.value = 0

      return data
    } catch (err) {
      console.error('Error marking all notifications as seen:', err)
      throw err
    }
  }

  // Handle notification actions (accept/decline)
  const handleNotificationAction = async (notification, action) => {
    try {
      // Extract necessary data from notification for different types
      const payload = {
        notification_id: notification.id,
        action: action
      }

      // Parse additional data from notification message or metadata
      if (notification.type === 'group_invite') {
        // Extract group ID from notification message
        const groupIdMatch = notification.message.match(/group.*?(\d+)/) || 
                           notification.message.match(/(?:to|for) (?:the )?(?:group )?.*?(\d+)/)
        if (groupIdMatch) {
          payload.group_id = parseInt(groupIdMatch[1])
        }
      } else if (notification.type === 'follow_request') {
        // Extract user ID from notification message
        const userIdMatch = notification.message.match(/user.*?(\d+)/) || 
                          notification.message.match(/from.*?(\d+)/)
        if (userIdMatch) {
          payload.user_id = parseInt(userIdMatch[1])
        }
      } else if (notification.type === 'group_request' || notification.type === 'group_join_request') {
        // Extract both group ID and user ID for group join requests
        const groupIdMatch = notification.message.match(/group.*?(\d+)/)
        const userIdMatch = notification.message.match(/user.*?(\d+)/) || 
                          notification.message.match(/from.*?(\d+)/)
        if (groupIdMatch) {
          payload.group_id = parseInt(groupIdMatch[1])
        }
        if (userIdMatch) {
          payload.user_id = parseInt(userIdMatch[1])
        }
      }

      const response = await fetch(`${API_BASE}/notifications/action`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }

      const data = await response.json()

      // Update notification as seen and refresh notifications list
      const notificationIndex = notifications.value.findIndex(n => n.id === notification.id)
      if (notificationIndex !== -1) {
        notifications.value[notificationIndex].seen = true
      }

      // Refresh notifications to get updated state
      await fetchNotifications({ start: 0, limit: 20 })

      return data
    } catch (err) {
      console.error('Error handling notification action:', err)
      throw err
    }
  }

  // Delete notification
  const deleteNotification = async (notificationId) => {
    try {
      const response = await fetch(`${API_BASE}/notifications/${notificationId}`, {
        method: 'DELETE'
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      
      // Update local state
      notifications.value = notifications.value.filter(n => n.id !== notificationId)
      unreadCount.value = data.unread_count || 0

      return data
    } catch (err) {
      console.error('Error deleting notification:', err)
      throw err
    }
  }

  // Get notification icon based on type
  const getNotificationIcon = (type) => {
    const icons = {
      'follow_request': '👤',
      'group_invite': '👥',
      'group_request': '📨',
      'group_join_request': '📨',
      'event_created': '📅',
      'follow_accepted': '✅',
      'group_accepted': '🎉',
      'default': '🔔'
    }
    return icons[type] || icons.default
  }

  // Get notification color based on type
  const getNotificationColor = (type) => {
    const colors = {
      'follow_request': '#3b82f6',
      'group_invite': '#8b5cf6',
      'group_request': '#f59e0b',
      'group_join_request': '#f59e0b',
      'event_created': '#10b981',
      'follow_accepted': '#22c55e',
      'group_accepted': '#8b5cf6',
      'default': '#6b7280'
    }
    return colors[type] || colors.default
  }

  // Format notification timestamp
  const formatNotificationTime = (createdAt) => {
    const date = new Date(createdAt)
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
    } else if (diffInSeconds < 604800) {
      const days = Math.floor(diffInSeconds / 86400)
      return `${days}d ago`
    } else {
      return date.toLocaleDateString()
    }
  }

  // Clear error
  const clearError = () => {
    error.value = null
  }

  // Disconnect WebSocket
  const disconnect = () => {
    if (socket.value) {
      socket.value.close()
      socket.value = null
      isConnected.value = false
    }
  }

  return {
    // State
    notifications,
    unreadCount,
    isLoading,
    error,
    isConnected,

    // Getters
    unreadNotifications,
    recentNotifications,

    // Actions
    connectWebSocket,
    disconnect,
    requestNotificationPermission,
    fetchNotifications,
    fetchUnreadCount,
    markAsSeen,
    markAllAsSeen,
    handleNotificationAction,
    deleteNotification,
    clearError,

    // Helpers
    getNotificationIcon,
    getNotificationColor,
    formatNotificationTime
  }
})