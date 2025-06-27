<template>
  <article class="tweet">
    <div class="tweet-header">
      <div class="tweet-avatar">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="var(--twitter-blue)">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
        </svg>
      </div>
      <div class="tweet-user">
        <!-- <span class="tweet-name">{{ post.author || 'Anonymous' }}</span> -->
        <span class="tweet-username">@{{ post.owner ? post.owner.toLowerCase().replace(/\s+/g, '') : 'anonymous' }}</span>
        <span class="tweet-dot">·</span>
        <span class="tweet-time">{{ formatDate(post.created_at) }}</span>
      </div>
    </div>
    
    <p class="tweet-content">{{ post.content }}</p>
    
    <div v-if="post.image" class="tweet-image">
      <img :src="`api/${post.image}`" alt="Post image" />
    </div>
    
    <div class="tweet-actions">
      <button class="tweet-action" @click="handleReplyClick">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="var(--twitter-gray)">
          <path d="M14.046 2.242l-4.148-.01h-.002c-4.374 0-7.8 3.427-7.8 7.802 0 4.098 3.186 7.206 7.465 7.37v3.828c0 .108.044.286.12.403.142.225.384.347.632.347.138 0 .277-.038.402-.118.264-.168 6.473-4.14 8.088-5.506 1.902-1.61 3.04-3.97 3.043-6.312v-.017c-.006-4.367-3.43-7.787-7.8-7.788zm3.787 12.972c-1.134.96-4.862 3.405-6.772 4.643V16.67c0-.414-.335-.75-.75-.75h-.396c-3.66 0-6.318-2.476-6.318-5.886 0-3.534 2.768-6.302 6.3-6.302l4.147.01h.002c3.532 0 6.3 2.766 6.302 6.296-.003 1.91-.942 3.844-2.514 5.176z"/>
        </svg>
        <span>{{ post.replies || 0 }}</span>
      </button>
    </div>
  </article>
</template>

<script setup>
import { useRouter } from 'vue-router'

const props = defineProps({
  post: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['reply'])
const router = useRouter()

const handleReplyClick = () => {
  // Navigate to post detail page with comments
  router.push({
    name: 'PostDetail',
    params: { id: props.post.id }
  })
  
  // Also emit the reply event for any parent component that needs it
  emit('reply', props.post)
}

const formatDate = (dateString) => {
  // console.log(dateString);
  
  if (!dateString) return 'Unknown date'
  try {
    const date = new Date(dateString)
    const now = new Date()
    const diffInSeconds = Math.floor((now - date) / 1000)
    
    if (diffInSeconds < 60) return `${diffInSeconds}s`
    if (diffInSeconds < 3600) return `${Math.floor(diffInSeconds / 60)}m`
    if (diffInSeconds < 86400) return `${Math.floor(diffInSeconds / 3600)}h`
    if (diffInSeconds < 2592000) return `${Math.floor(diffInSeconds / 86400)}d`
    
    return date.toLocaleDateString('en-US', {
      month: 'short',
      day: 'numeric'
    })
  } catch (error) {
    return 'Invalid date'
  }
}
</script>

<style scoped>
.tweet {
  padding: 1rem;
  border-bottom: 1px solid var(--twitter-extra-light-gray);
  transition: background-color 0.2s;
}

.tweet:hover {
  background-color: var(--twitter-extra-extra-light-gray);
}

.tweet-header {
  display: flex;
  margin-bottom: 0.5rem;
}

.tweet-avatar {
  margin-right: 0.75rem;
}

.tweet-avatar svg {
  width: 48px;
  height: 48px;
}

.tweet-user {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.25rem;
  font-size: 0.9rem;
}

.tweet-name {
  font-weight: 700;
  color: var(--twitter-dark);
}

.tweet-username, .tweet-time {
  color: var(--twitter-gray);
}

.tweet-dot {
  color: var(--twitter-gray);
  padding: 0 0.25rem;
}

.tweet-content {
  margin: 0.5rem 0;
  font-size: 1.1rem;
  line-height: 1.5;
}

.tweet-image {
  margin: 0.75rem 0;
  border-radius: 1rem;
  overflow: hidden;
  border: 1px solid var(--twitter-extra-light-gray);
}

.tweet-image img {
  display: block;
  max-width: 100%;
  max-height: 500px;
  object-fit: cover;
}

.tweet-actions {
  display: flex;
  justify-content: space-between;
  margin-top: 0.75rem;
  max-width: 425px;
}

.tweet-action {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: var(--twitter-gray);
  background: none;
  border: none;
  padding: 0.5rem;
  font-size: 0.9rem;
  transition: color 0.2s;
  cursor: pointer;
}

.tweet-action:hover {
  color: var(--twitter-blue);
}

.tweet-action svg {
  width: 20px;
  height: 20px;
}

@media (max-width: 500px) {
  .tweet-avatar svg {
    width: 40px;
    height: 40px;
  }
  
  .tweet-content {
    font-size: 1rem;
  }
  
  .tweet-actions {
    max-width: 100%;
  }
}
</style>