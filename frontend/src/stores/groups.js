import { defineStore } from 'pinia'
import { ref, computed } from 'vue'

export const useGroupsStore = defineStore('groups', () => {
  const groups = ref([])
  const currentGroup = ref(null)
  const groupPosts = ref([])
  const groupEvents = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  const API_BASE = '/api'

  // Getters
  const getGroupById = computed(() => (id) => {
    return groups.value.find(group => group.id === id)
  })

  const getCurrentGroupPosts = computed(() => {
    return groupPosts.value.filter(post => post.groupId === currentGroup.value?.id)
  })

  const getCurrentGroupEvents = computed(() => {
    return groupEvents.value.filter(event => event.groupId === currentGroup.value?.id)
  })

  // Transform functions
  const transformGroupData = (apiGroup) => {
    return {
      id: apiGroup.id,
      name: apiGroup.title,
      description: apiGroup.description,
      image: apiGroup.image_uuid ? `${API_BASE}/images/${apiGroup.image_uuid}` : '/default-group.jpg',
      isPublic: true,
      memberCount: apiGroup.member_count || 0,
      isMember: apiGroup.is_member || '',
      createdAt: apiGroup.created_at,
      creatorId: apiGroup.creator_id
    }
  }

  const transformPostData = (apiPost) => {
    return {
      id: apiPost.id,
      groupId: apiPost.group_id,
      title: apiPost.title,
      content: apiPost.content,
      author: apiPost.author_name,
      authorAvatar: apiPost.author_avatar ? `${API_BASE}/images/${apiPost.author_avatar}` : '/default-avatar.jpg',
      createdAt: apiPost.created_at,
      likes: apiPost.likes_count || 0,
      comments: apiPost.comments_count || 0
    }
  }

  const transformEventData = (apiEvent) => {
    return {
      id: apiEvent.id,
      groupId: apiEvent.group_id,
      title: apiEvent.title,
      description: apiEvent.description,
      date: apiEvent.event_time,
      location: null, // Not available in current backend model
      attendees: apiEvent.vote_count || 0,
      maxAttendees: null, // Not available in current backend model
      image: null, // Not available in current backend model
      organizer: null, // Not available in current backend model
      isAttending: apiEvent.user_vote || '', // Get from backend user_vote field
      createdAt: apiEvent.created_at
    }
  }

  // Actions
  const fetchGroups = async (filter = 'all') => {
    if (isLoading.value) return

    isLoading.value = true
    error.value = null

    const temp = filter === 'user'
      ? JSON.stringify({ user_id: '1', start: -1, n_items: 20, type: filter })
      : JSON.stringify({ start: -1, n_items: 20, type: filter })
    
    try {
      const response = await fetch(`${API_BASE}/groups/group/browse`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: temp
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      let groupsData = []

      if (Array.isArray(data)) {
        groupsData = data
      } else if (data.groups && Array.isArray(data.groups)) {
        groupsData = data.groups
      } else if (data.data && Array.isArray(data.data)) {
        groupsData = data.data
      } else {
        groupsData = []
      }

      const transformedGroups = groupsData.map(transformGroupData)
      groups.value = transformedGroups

    } catch (err) {
      error.value = err.message
      groups.value = []
    } finally {
      isLoading.value = false
    }
  }

  const fetchGroup = async (groupId) => {
    isLoading.value = true
    error.value = null

    try {
      // First check if we already have the group in our store
      const existingGroup = getGroupById.value(groupId)
      if (existingGroup && currentGroup.value?.id === groupId) {
        isLoading.value = false
        return existingGroup
      }

      // Fetch from API
      const response = await fetch(`${API_BASE}/groups/group/${groupId}`)
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      const transformedGroup = transformGroupData(data)

      // Update the groups array if this group exists in it
      const existingIndex = groups.value.findIndex(g => g.id === groupId)
      if (existingIndex !== -1) {
        groups.value[existingIndex] = transformedGroup
      } else {
        groups.value.push(transformedGroup)
      }

      // Set as current group
      currentGroup.value = transformedGroup

      return transformedGroup
    } catch (err) {
      error.value = err.message
      currentGroup.value = null
      console.error('Error fetching group:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const fetchGroupPosts = async (groupId) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/posts/feed`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ id: groupId, type: 'group', start: 1, n_post: 20 })
      })

      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      const transformedPosts = Array.isArray(data) ? data.map(transformPostData) : []

      // Clear previous posts and set new ones
      groupPosts.value = transformedPosts
      
      return transformedPosts
    } catch (err) {
      error.value = err.message
      groupPosts.value = []
      console.error('Error fetching group posts:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const fetchGroupEvents = async (groupId) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/group/events`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ group_id: groupId, start: -1, n_items: 20})
      })
      if (!response.ok) {
        throw new Error(`HTTP error! status: ${response.status}`)
      }

      const data = await response.json()
      const transformedEvents = Array.isArray(data) ? data.map(transformEventData) : []

      groupEvents.value = transformedEvents
      return transformedEvents
    } catch (err) {
      error.value = err.message
      groupEvents.value = []
      console.error('Error fetching group events:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  } 

  const createGroup = async (groupData) => {
    isLoading.value = true
    error.value = null

    try {
      const apiGroupData = {
        title: groupData.name,
        description: groupData.description,
        image_uuid: groupData.image || '',
      }

      const response = await fetch(`${API_BASE}/groups/group/new`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(apiGroupData)
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }

      const newGroupData = await response.json()
      const newGroup = transformGroupData(newGroupData)

      groups.value.unshift(newGroup)
      return newGroup
    } catch (err) {
      error.value = err.message
      console.error('Error creating group:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const createPost = async (groupId, postData) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/group/${groupId}/post`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          title: postData.title,
          content: postData.content
        })
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }

      const newPostData = await response.json()
      const newPost = transformPostData(newPostData)

      groupPosts.value.unshift(newPost)
      return newPost
    } catch (err) {
      error.value = err.message
      console.error('Error creating post:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const createEvent = async (groupId, eventData) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/group/event`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          group_id: groupId,
          title: eventData.title,
          description: eventData.description,
          event_time: new Date(eventData.date).toISOString()
        })
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }

      const newEventData = await response.json()
      const newEvent = transformEventData(newEventData)

      groupEvents.value.unshift(newEvent)
      return newEvent
    } catch (err) {
      error.value = err.message
      console.error('Error creating event:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const requestJoinGroup = async (groupId) => {
    // isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/group/request`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          user_id: 1, // TODO Replace with actual user ID from context or store
          group_id: groupId,
          status: 'requested',
          prev_status: 'requested'
        })
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || 'Failed to join group')
      }

      const data = await response.json()

      // Update local state
      const updatedGroups = groups.value.map(group => {
        if (group.id === groupId) {
          return {
            ...group,
            isMember: data.is_member || '',
            memberCount: data.is_member ? group.memberCount + 1 : group.memberCount
          }
        }
        return group
      })

      groups.value = updatedGroups

      if (currentGroup.value?.id === groupId) {
        currentGroup.value = {
          ...currentGroup.value,
          isMember: data.is_member || '',
          memberCount: data.is_member ? currentGroup.value.memberCount + 1 : currentGroup.value.memberCount
        }
      }

      // return data
    } catch (err) {
      error.value = err.message
      console.error('Error joining group:', err)
      throw err
    } finally {
      // isLoading.value = false
    }
  }

  const leaveGroup = async (groupId) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/group/${groupId}/leave`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' }
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || `HTTP error! status: ${response.status}`)
      }

      const data = await response.json()

      // Update local state
      const updatedGroups = groups.value.map(group => {
        if (group.id === groupId) {
          return {
            ...group,
            isMember: '',
            memberCount: Math.max(0, group.memberCount - 1)
          }
        }
        return group
      })

      groups.value = updatedGroups

      if (currentGroup.value?.id === groupId) {
        currentGroup.value = {
          ...currentGroup.value,
          isMember: '',
          memberCount: Math.max(0, currentGroup.value.memberCount - 1)
        }
      }

      return data
    } catch (err) {
      error.value = err.message
      console.error('Error leaving group:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const acceptGroupInvite = async (groupId) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/group/accept-decline`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          group_id: groupId,
          user_id: 1, // TODO: Replace with actual user ID
          status: 'member',
          prev_status: 'invited'
        })
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || 'Failed to accept invitation')
      }

      const data = await response.json()

      // Update local state
      const updatedGroups = groups.value.map(group => {
        if (group.id === groupId) {
          return {
            ...group,
            isMember: 'member',
            memberCount: group.memberCount + 1
          }
        }
        return group
      })

      groups.value = updatedGroups

      if (currentGroup.value?.id === groupId) {
        currentGroup.value = {
          ...currentGroup.value,
          isMember: 'member',
          memberCount: currentGroup.value.memberCount + 1
        }
      }

      return data
    } catch (err) {
      error.value = err.message
      console.error('Error accepting group invitation:', err)
      throw err
    } finally {
      isLoading.value = false
    }
  }

  const clearError = () => {
    error.value = null
  }

  const attendEvent = async (eventId, voteType = 'going') => {
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/group/event/vote`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          event_id: eventId,
          user_id: 1, // TODO: Replace with actual user ID
          vote: voteType
        })
      })

      if (!response.ok) {
        const errorData = await response.json().catch(() => ({}))
        throw new Error(errorData.message || 'Failed to vote on event')
      }

      const data = await response.json()

      // Update the event in the local state
      const eventIndex = groupEvents.value.findIndex(event => event.id === eventId)
      if (eventIndex !== -1) {
        const currentEvent = groupEvents.value[eventIndex]
        let newAttendees = currentEvent.attendees || 0
        
        // Adjust attendees count based on previous and new vote
        if (currentEvent.isAttending === 'going' && voteType !== 'going') {
          newAttendees = Math.max(0, newAttendees - 1)
        } else if (currentEvent.isAttending !== 'going' && voteType === 'going') {
          newAttendees += 1
        }
        
        groupEvents.value[eventIndex] = {
          ...currentEvent,
          attendees: newAttendees,
          isAttending: voteType
        }
      }

      return data
    } catch (err) {
      error.value = err.message
      console.error('Error voting on event:', err)
      throw err
    }
  }

  return {
    // State
    groups,
    currentGroup,
    groupPosts,
    groupEvents,
    isLoading,
    error,

    // Getters
    getGroupById,
    getCurrentGroupPosts,
    getCurrentGroupEvents,

    // Actions
    fetchGroups,
    fetchGroup,
    fetchGroupPosts,
    fetchGroupEvents,
    createGroup,
    createPost,
    createEvent,
    requestJoinGroup,
    leaveGroup,
    attendEvent,
    clearError
  }
})