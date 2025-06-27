<!-- src/components/comments/CommentItem.vue -->
<template>
  <article class="comment">
    <div class="comment-header">
      <div class="comment-avatar">
        <svg xmlns="http://www.w3.org/2000/svg" width="40" height="40" viewBox="0 0 24 24" fill="var(--twitter-blue)">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
        </svg>
      </div>
      
      <div class="comment-content">
        <div class="comment-user">
          <span class="comment-name">{{ comment.author || 'Anonymous' }}</span>
          <span class="comment-username">@{{ comment.author ? comment.author.toLowerCase().replace(/\s+/g, '') : 'anonymous' }}</span>
          <span class="comment-dot">·</span>
          <span class="comment-time">{{ formatDate(comment.createdAt) }}</span>
        </div>
        
        <p class="comment-text">{{ comment.content }}</p>
        
        <!-- <div class="comment-actions">
          <button class="comment-action" @click="$emit('reply', comment)">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" fill="currentColor">
              <path d="M14.046 2.242l-4.148-.01h-.002c-4.374 0-7.8 3.427-7.8 7.802 0 4.098 3.186 7.206 7.465 7.37v3.828c0 .108.044.286.12.403.142.225.384.347.632.347.138 0 .277-.038.402-.118.264-.168 6.473-4.14 8.088-5.506 1.902-1.61 3.04-3.97 3.043-6.312v-.017c-.006-4.367-3.43-7.787-7.8-7.788zm3.787 12.972c-1.134.96-4.862 3.405-6.772 4.643V16.67c0-.414-.335-.75-.75-.75h-.396c-3.66 0-6.318-2.476-6.318-5.886 0-3.534 2.768-6.302 6.3-6.302l4.147.01h.002c3.532 0 6.3 2.766 6.302 6.296-.003 1.91-.942 3.844-2.514 5.176z"/>
            </svg>
          </button> -->
          
          <!-- <button class="comment-action" @click="toggleLike">
            <svg xmlns="http://www.w3.org/2000/svg" width="18" height="18" viewBox="0 0 24 24" :fill="isLiked ? '#e91e63' : 'currentColor'">
              <path d="M12 21.638h-.014C9.403 21.59 1.95 14.856 1.95 8.478c0-3.064 2.525-5.754 5.403-5.754 2.29 0 3.83 1.58 4.646 2.73.814-1.148 2.354-2.73 4.645-2.73 2.88 0 5.404 2.69 5.404 5.755 0 6.376-7.454 13.11-10.037 13.157H12z"/>
            </svg>
            <span v-if="likesCount > 0">{{ likesCount }}</span>
          </button> -->
        <!-- </div> -->
      </div>
    </div>
  </article>
</template>

<script setup>
// import { ref, computed } from 'vue'

const props = defineProps({
  comment: {
    type: Object,
    required: true
  }
})

// const emit = defineEmits(['reply', 'like'])

// const isLiked = ref(false)
// const likesCount = ref(props.comment.likes || 0)

// const toggleLike = () => {
//   isLiked.value = !isLiked.value
//   likesCount.value += isLiked.value ? 1 : -1
//   emit('like', { commentId: props.comment.id, liked: isLiked.value })
// }

const formatDate = (dateString) => {
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
.comment {
  padding: 0.75rem 1rem;
  border-bottom: 1px solid var(--twitter-extra-light-gray);
  transition: background-color 0.2s;
}

.comment:hover {
  background-color: var(--twitter-extra-extra-light-gray);
}

.comment-header {
  display: flex;
  gap: 0.75rem;
}

.comment-avatar svg {
  width: 40px;
  height: 40px;
}

.comment-content {
  flex: 1;
  min-width: 0;
}

.comment-user {
  display: flex;
  align-items: center;
  flex-wrap: wrap;
  gap: 0.25rem;
  font-size: 0.875rem;
  margin-bottom: 0.25rem;
}

.comment-name {
  font-weight: 700;
  color: var(--twitter-dark);
}

.comment-username, .comment-time {
  color: var(--twitter-gray);
}

.comment-dot {
  color: var(--twitter-gray);
  padding: 0 0.125rem;
}

.comment-text {
  margin: 0.25rem 0 0.5rem 0;
  font-size: 0.95rem;
  line-height: 1.4;
  color: var(--twitter-dark);
  word-wrap: break-word;
}

.comment-actions {
  display: flex;
  gap: 1rem;
  margin-top: 0.5rem;
}

.comment-action {
  display: flex;
  align-items: center;
  gap: 0.25rem;
  color: var(--twitter-gray);
  background: none;
  border: none;
  padding: 0.25rem;
  font-size: 0.8rem;
  cursor: pointer;
  border-radius: 9999px;
  transition: all 0.2s;
  min-width: 32px;
  height: 32px;
  justify-content: center;
}

.comment-action:hover {
  background-color: rgba(29, 161, 242, 0.1);
  color: var(--twitter-blue);
}

.comment-action:nth-child(2):hover {
  background-color: rgba(233, 30, 99, 0.1);
  color: #e91e63;
}

.comment-action svg {
  width: 16px;
  height: 16px;
}

@media (max-width: 500px) {
  .comment {
    padding: 0.5rem;
  }
  
  .comment-avatar svg {
    width: 36px;
    height: 36px;
  }
  
  .comment-text {
    font-size: 0.9rem;
  }
}
</style>