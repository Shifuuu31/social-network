<template>
  <div class="post-item">
    <div class="post-header">
      <div class="post-author">
        <img 
          :src="post.author?.avatar || '/default-avatar.png'" 
          :alt="post.author?.name || 'User'" 
          class="avatar"
        >
        <div class="author-info">
          <h4 class="author-name">{{ post.author?.name || 'Anonymous' }}</h4>
          <span class="post-time">{{ formatDate(post.createdAt) }}</span>
        </div>
      </div>
    </div>
    
    <div class="post-content">
      <p v-if="post.content">{{ post.content }}</p>
      <img 
        v-if="post.image" 
        :src="post.image" 
        :alt="post.content || 'Post image'"
        class="post-image"
      >
    </div>
    
    <div class="post-actions">
      <button class="action-btn like-btn" :class="{ active: post.isLiked }">
        <span class="icon">❤️</span>
        <span class="count">{{ post.likesCount || 0 }}</span>
      </button>
      
      <button class="action-btn comment-btn">
        <span class="icon">💬</span>
        <span class="count">{{ post.commentsCount || 0 }}</span>
      </button>
  
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

// Props - expects a single post object
const props = defineProps({
  post: {
    type: Object,
    required: true
  }
})

// Helper function to format date
function formatDate(dateString) {
  if (!dateString) return 'Just now'
  
  const date = new Date(dateString)
  const now = new Date()
  const diffInHours = Math.floor((now - date) / (1000 * 60 * 60))
  
  if (diffInHours < 1) return 'Just now'
  if (diffInHours < 24) return `${diffInHours}h ago`
  if (diffInHours < 168) return `${Math.floor(diffInHours / 24)}d ago`
  
  return date.toLocaleDateString()
}
</script>

<style scoped>
.post-item {
  background: white;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 16px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  border: 1px solid #e1e5e9;
}

.post-header {
  margin-bottom: 12px;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 12px;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.author-info {
  flex: 1;
}

.author-name {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #1c1e21;
}

.post-time {
  font-size: 12px;
  color: #65676b;
}

.post-content {
  margin-bottom: 12px;
}

.post-content p {
  margin: 0 0 12px 0;
  color: #1c1e21;
  line-height: 1.4;
}

.post-image {
  width: 100%;
  max-height: 400px;
  object-fit: cover;
  border-radius: 8px;
}

.post-actions {
  display: flex;
  gap: 8px;
  padding-top: 8px;
  border-top: 1px solid #e4e6ea;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 4px;
  padding: 6px 12px;
  background: none;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 13px;
  color: #65676b;
  transition: background-color 0.2s;
}

.action-btn:hover {
  background-color: #f2f3f4;
}

.action-btn.active {
  color: #e41e3f;
}

.action-btn .icon {
  font-size: 16px;
}

.action-btn .count {
  font-weight: 500;
}
</style>