import { ref, reactive } from 'vue'

export const useGroups = () => {
  const groups = ref([])
  const currentGroup = ref(null)
  const groupPosts = ref([])
  const isLoading = ref(false)
  const error = ref(null)

  const API_BASE = '/api'

  const fetchGroups = async () => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups`)
      if (!response.ok) throw new Error('Failed to fetch groups')

      const data = await response.json()
      groups.value = data
      return data
    } catch (err) {
      error.value = err.message
      console.error('Error fetching groups:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Get specific group details
  const fetchGroup = async (groupId) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/${groupId}`)
      if (!response.ok) throw new Error('Failed to fetch group')

      const data = await response.json()
      currentGroup.value = data
      return data
    } catch (err) {
      error.value = err.message
      console.error('Error fetching group:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Get posts for a specific group
  const fetchGroupPosts = async (groupId) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/${groupId}/posts`)
      if (!response.ok) throw new Error('Failed to fetch group posts')

      const data = await response.json()
      groupPosts.value = data
      return data
    } catch (err) {
      error.value = err.message
      console.error('Error fetching group posts:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Create a new post in a group
  const createGroupPost = async (groupId, postData) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/${groupId}/posts`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          // Add authorization header if needed
          // 'Authorization': `Bearer ${token}`
        },
        body: JSON.stringify(postData)
      })

      if (!response.ok) throw new Error('Failed to create post')

      const newPost = await response.json()
      groupPosts.value.unshift(newPost)
      return newPost
    } catch (err) {
      error.value = err.message
      console.error('Error creating post:', err)
    } finally {
      isLoading.value = false
    }
  }


  // Join a group
  // const joinGroup = async (groupId) => {
  //   isLoading.value = true
  //   error.value = null
  //   console.log(`Joining group with ID: ${groupId}`);

  //   try {
  //     const response = await fetch(`${API_BASE}/groups/${groupId}/group/request`, {
  //       method: 'POST',
  //       headers: {
  //         'Content-Type': 'application/json',
  //         // Add authorization header if needed
  //       },
  //       body: JSON.stringify({ 
  //         group_id: groupId, 
  //         status: 'requested'
  //       })
  //     })

  //     if (!response.ok) throw new Error('Failed to join group')

  //     const data = await response.json()
  //     // Update local group data
  //     if (currentGroup.value && currentGroup.value.id === groupId) {
  //       currentGroup.value.isMember = true
  //       currentGroup.value.memberCount += 1
  //     }
  //     return data
  //   } catch (err) {
  //     error.value = err.message
  //     console.error('Error joining group:', err)
  //   } finally {
  //     isLoading.value = false
  //   }
  // }

  // Join a group
const joinGroup = async (groupId) => {
  isLoading.value = true
  error.value = null
  console.log(`Joining group with ID: ${groupId}`);
  
  try {
    const response = await fetch(`${API_BASE}/groups/group/request`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify({
        group_id: groupId,
        status: 'requested'
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
          isMember: true,
          memberCount: group.memberCount + 1
        }
      }
      return group
    })
    
    groups.value = updatedGroups
    
    if (currentGroup.value?.id === groupId) {
      currentGroup.value = {
        ...currentGroup.value,
        isMember: true,
        memberCount: currentGroup.value.memberCount + 1
      }
    }
    
    return data
  } catch (err) {
    error.value = err.message
    console.error('Error joining group:', err)
    throw err
  } finally {
    isLoading.value = false
  }
}

  // Leave a group
  const leaveGroup = async (groupId) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/${groupId}/leave`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          // Add authorization header if needed
        }
      })

      if (!response.ok) throw new Error('Failed to leave group')

      const data = await response.json()
      // Update local group data
      if (currentGroup.value && currentGroup.value.id === groupId) {
        currentGroup.value.isMember = false
        currentGroup.value.memberCount -= 1
      }
      return data
    } catch (err) {
      error.value = err.message
      console.error('Error leaving group:', err)
    } finally {
      isLoading.value = false
    }
  }

  // Search groups
  const searchGroups = async (query) => {
    isLoading.value = true
    error.value = null

    try {
      const response = await fetch(`${API_BASE}/groups/search?q=${encodeURIComponent(query)}`)
      if (!response.ok) throw new Error('Failed to search groups')

      const data = await response.json()
      return data
    } catch (err) {
      error.value = err.message
      console.error('Error searching groups:', err)
    } finally {
      isLoading.value = false
    }
  }

  const clearError = () => {
    error.value = null
  }

  return {
    groups,
    currentGroup,
    groupPosts,
    isLoading,
    error,

    fetchGroups,
    fetchGroup,
    fetchGroupPosts,
    createGroupPost,
    joinGroup,
    leaveGroup,
    searchGroups,
    clearError
  }
}