<template>
  <div class="post-list">
    <div v-if="loading" class="loading">Loading posts...</div>
    <div v-else-if="posts.length === 0" class="no-posts">
      No posts found. Be the first to share something!
    </div>
    <div v-else class="posts">
      <PostItem v-for="post in posts" :key="post.id" :post="post" :currentUser="currentUser" />
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
import { useAuth } from '@/composables/useAuth.js'

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
  },
  postType: {
    type: String,
    default: 'public'
  }
})

// Get current user
const { user: currentUser } = useAuth()

// State
const posts = ref([])
const loading = ref(false)
const error = ref(null)

// Helper function to determine post type
function determinePostType() {
  if (props.groupId) {
    return "group"
  }
  return props.postType
}

// Fetch posts
async function fetchPosts() {
  loading.value = true
  error.value = null
  try {
    const result = await getPosts({
      id: props.userId,
      type: determinePostType(),
      start: props.offset,
      nPost: props.limit
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

// Watch for changes in props to refetch posts
watch(
  [() => props.userId, () => props.groupId, () => props.postType, () => props.offset, () => props.limit],
  () => {
    fetchPosts()
  }
)
</script>