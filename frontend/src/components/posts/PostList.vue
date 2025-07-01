<template>
  <div class="post-list">
    <div v-if="loading" class="loading">Loading posts...</div>

    <div v-else-if="posts.length === 0" class="no-posts">
      No posts found. Be the first to share something!
    </div>

    <div v-else class="posts">
      <PostItem v-for="post in posts" :key="post.id" :post="post" />
    </div>

    <div v-if="error" class="error-message">
      {{ error }}
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { getPosts } from '@/services/api'
import PostItem from './PostItem.vue'

// Props
const props = defineProps({
  userId: {
    type: Number,
    required: true
  },
  groupId: {
    type: [Number, Object],
    default: null
  },
  limit: {
    type: Number,
    default: 20
  },
  offset: {
    type: Number,
    default: 0
  },
  refreshTrigger: {
    type: Number,
    default: 0
  }
})

// State
const posts = ref([])
const loading = ref(false)
const error = ref(null)

// Fetch posts
async function fetchPosts() {
  loading.value = true
  error.value = null

  try {
    const result = await getPosts({
      user_id: props.userId,
      group_id: props.groupId,
      limit: props.limit,
      offset: props.offset
    })

    posts.value = result || []
  } catch (err) {
    console.error('Failed to load posts:', err.message)
    error.value = 'Failed to load posts. Please try again later.'
    posts.value = []
  } finally {
    loading.value = false
  }
}

// Initial load
onMounted(() => {
  fetchPosts()
})

// Reload when refreshTrigger changes
watch(
  () => props.refreshTrigger,
  () => {
    fetchPosts()
  }
)
</script>

<style scoped>
.post-list {
  margin-top: 1rem;
}

.loading,
.no-posts {
  text-align: center;
  padding: 1rem;
  color: #666;
}

.posts {
  display: flex;
  flex-direction: column;
  gap: 1rem;
}

.error-message {
  color: red;
  text-align: center;
  margin-top: 1rem;
}
</style>