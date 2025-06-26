// src/composables/usePosts.js
import { ref } from 'vue'
import { useToast } from './useToast'

export function usePosts() {
  const { showToast } = useToast()
  const posts = ref([])
  const loading = ref(false)
  const error = ref(null)

  const fetchPosts = async (body) => {
    try {
      loading.value = true
      const response = await fetch('/api/post/feed', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(body)
      })
      if (!response.ok) throw new Error('Failed to fetch posts')
      posts.value = await response.json()
    } catch (err) {
      error.value = err.message
      showToast('Failed to load posts', 'error')
    } finally {
      loading.value = false
    }
  }

  const createPost = async (postData) => {    
    try {
      loading.value = true
      const response = await fetch('/api/post/new', {
        method: 'POST',
       body: postData
      })
      if (!response.ok) throw new Error('Failed to create post')
      showToast('Post created successfully!', 'success')
      return await response.json()
    } catch (err) {
      error.value = err.message
      showToast('Failed to create post', 'error')
      throw err
    } finally {
      loading.value = false
    }
  }

  return { posts, loading, error, fetchPosts, createPost }
}