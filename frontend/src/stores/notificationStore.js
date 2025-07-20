import { defineStore } from 'pinia'
import { ref } from 'vue'
import chatService from '@/services/chatService' // Adjust path as needed

export const useNotificationStore = defineStore('notification', () => {
  const notifications = ref([])

  // Add a notification locally
  const addNotification = (notification) => {
    notifications.value.push({
      id: notification.id || Date.now() + Math.random(),
      ...notification,
      timestamp: notification.timestamp ? new Date(notification.timestamp) : new Date(),
      read: notification.read ?? false
    })
  }

  // Remove a notification
  const removeNotification = async (notificationId) => {
    try {
      await fetch(`http://localhost:8080/api/notifications/${notificationId}`, {
        method: 'DELETE',
        credentials: 'include'
      })
      notifications.value = notifications.value.filter(n => n.id !== notificationId)
    } catch (err) {
      console.error('Failed to remove notification:', err)
    }
  }

  // Mark as read
  const markAsRead = async (notificationId) => {
    try {
      const res = await fetch(`http://localhost:8080/api/notifications/mark-read`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ notification_ids: [notificationId] })
      })
      if (!res.ok) throw new Error('Failed to mark as read')
      const notification = notifications.value.find(n => n.id === notificationId)
      if (notification) notification.seen = true
    } catch (error) {
      console.error('Failed to mark notification as read:', error)
    }
  }

  // Mark multiple notifications as read
  const markMultipleAsRead = async (notificationIds) => {
    try {
      const res = await fetch(`http://localhost:8080/api/notifications/mark-read`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        credentials: 'include',
        body: JSON.stringify({ notification_ids: notificationIds })
      })
      if (!res.ok) throw new Error('Failed to mark as read')
      notificationIds.forEach(id => {
        const notification = notifications.value.find(n => n.id === id)
        if (notification) notification.seen = true
      })
    } catch (error) {
      console.error('Failed to mark notifications as read:', error)
    }
  }

  // Fetch existing notifications
  const fetchNotifications = async () => {
    try {
      const res = await fetch(`http://localhost:8080/api/notifications`, {
        credentials: 'include'
      })
      if (!res.ok) throw new Error('Fetch failed')
      const data = await res.json()
      notifications.value = data.notifications || []
    } catch (error) {
      console.error('Failed to fetch notifications:', error)
    }
  }

  // Get unread count from server
  const fetchUnreadCount = async () => {
    try {
      const res = await fetch(`http://localhost:8080/api/notifications/unread-count`, {
        credentials: 'include'
      })
      if (!res.ok) throw new Error('Fetch failed')
      const data = await res.json()
      return data.unread_count || 0
    } catch (error) {
      console.error('Failed to fetch unread count:', error)
      return 0
    }
  }

  // Fetch only unseen notifications
  const fetchUnseenNotifications = async (page = 1, limit = 20) => {
    try {
      const res = await fetch(`http://localhost:8080/api/notifications?page=${page}&limit=${limit}&unseen_only=true`, {
        credentials: 'include'
      })
      if (!res.ok) throw new Error('Fetch failed')
      const data = await res.json()
      
      if (page === 1) {
        // Replace notifications for first page
        notifications.value = data.notifications || []
      } else {
        // Append for subsequent pages
        notifications.value.push(...(data.notifications || []))
      }
      
      return {
        notifications: data.notifications || [],
        total: data.total || 0,
        page: data.page || 1,
        hasMore: (data.notifications || []).length === limit
      }
    } catch (error) {
      console.error('Failed to fetch unseen notifications:', error)
      return {
        notifications: [],
        total: 0,
        page: 1,
        hasMore: false
      }
    }
  }

  // Helper: Create a follow request notification
  const createFollowRequest = (fromUser, recipientId) => {
    addNotification({
      type: 'follow_request',
      title: 'Follow Request',
      message: `${fromUser.username} has requested to follow you.`,
      userId: recipientId,
      read: false
    })
  }

  // Helper: Create a group invitation notification
  const createGroupInvitation = (group, fromUser, recipientId) => {
    addNotification({
      type: 'group_invite',
      title: 'Group Invitation',
      message: `You’ve been invited to join the group ${group.name}.`,
      groupId: group.id,
      fromUser: fromUser.username,
      userId: recipientId,
      read: false
    })
  }

  // Helper: Create a group join request notification (for group creator)
  const createGroupJoinRequest = (group, fromUser, creatorId) => {
    addNotification({
      type: 'group_join_request',
      title: 'Group Join Request',
      message: `${fromUser.username} requested to join your group ${group.name}.`,
      groupId: group.id,
      fromUser: fromUser.username,
      userId: creatorId,
      read: false
    })
  }

  // Helper: Create a group event notification
  const createGroupEvent = (group, event, recipientId) => {
    addNotification({
      type: 'group_event',
      title: 'Group Event',
      message: `A new event '${event.title}' has been created in group '${group.name}'.`,
      groupId: group.id,
      eventId: event.id,
      userId: recipientId,
      read: false
    })
  }

  // Listen for real-time notifications via chatService
  const setupNotificationListener = () => {
    chatService.onMessage((message) => {
      try {
        // Handle different message types for notifications
        if (message.type === 'notification' || message.type === 'new_notification') {
          addNotification(message.notification || message)
        }
        // Handle notification count updates
        if (message.type === 'notification_count_updated') {
          // Trigger a refresh of notifications or update count directly
          fetchNotifications()
        }
      } catch (err) {
        console.error('Error processing notification message:', err)
      }
    })
  }

  // Initialize the listener when the store is created
  setupNotificationListener()

  // Get unread count
  const unreadCount = () => {
    return notifications.value.filter(n => !n.seen).length
  }

  // Clear all
  const clearAll = async () => {
    try {
      for (const n of notifications.value) {
        await fetch(`http://localhost:8080/api/notifications/${n.id}`, { 
          method: 'DELETE',
          credentials: 'include'
        })
      }
      notifications.value = []
    } catch (err) {
      console.error('Failed to clear all notifications:', err)
    }
  }

  return {
    notifications,
    addNotification,
    removeNotification,
    markAsRead,
    markMultipleAsRead,
    fetchNotifications,
    fetchUnreadCount,
    fetchUnseenNotifications,
    createFollowRequest,
    createGroupInvitation,
    createGroupJoinRequest,
    createGroupEvent,
    setupNotificationListener,
    unreadCount,
    clearAll
  }
})