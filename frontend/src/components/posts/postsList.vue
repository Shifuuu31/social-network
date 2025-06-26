<template>
  <div class="posts-container">
    <div v-if="loading" class="loading-spinner">Loading...</div>
    <div v-if="error" class="error-message">{{ error }}</div>
    
    <div v-if="!loading && !error && posts.length === 0" class="no-posts">
      No posts found. Create your first post!
    </div>
    
    <PostItem 
      v-for="post in posts" 
      :key="post.id"
      :post="post" 
    />
  </div>
</template>

<script setup>
import { onMounted } from 'vue'
import PostItem from './PostItem.vue'
import { usePosts } from '@/composables/usePosts'

const { posts, loading, error, fetchPosts } = usePosts()

// Create a wrapper function that calls fetchPosts with default params
const refreshPosts = () => {
  fetchPosts({
    type: 'privacy',
    start: 1,
    n_post: 10
  })
}

// Fetch posts when component mounts
onMounted(() => {
  refreshPosts()
})

// Expose the refresh function so parent component can call it
defineExpose({
  fetchPosts: refreshPosts,
  refreshPosts // Alternative name
})

// const handleDelete = async (postId) => { // might wanna implement this later
//   try {
//     // Call delete API
//     const response = await fetch(`http://localhost:8080/post/${postId}`, {
//       method: 'DELETE'
//     })
    
//     if (!response.ok) throw new Error('Failed to delete post')
    
//     // Remove from local state
//     posts.value = posts.value.filter(post => post.id !== postId)
//   } catch (err) {
//     console.error('Delete error:', err)
//     // You could show a toast error here
//   }
// }
</script>

<style scoped>
.posts-container {
  width: 100%;
  max-width: 100%;
  margin: 0 auto;
  padding: 1rem;
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.loading-spinner {
  text-align: center;
  padding: 2rem;
  font-size: 1.2rem;
}

.error-message {
  background: #f8d7da;
  color: #721c24;
  padding: 1rem;
  border-radius: 4px;
  margin-bottom: 1rem;
}

.no-posts {
  text-align: center;
  padding: 2rem;
  color: #6c757d;
  font-style: italic;
}
</style>