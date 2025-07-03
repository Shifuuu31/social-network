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
      <button 
        class="action-btn comment-btn" 
        @click="toggleComments"
        :class="{ active: showComments }"
      >
        <span class="icon">💬</span>
        <span class="count">{{ totalCommentsCount }}</span>
      </button>
    </div>

    <!-- Comments Section -->
    <div v-if="showComments" class="comments-section">
      <!-- Add Comment Form -->
      <div class="add-comment">
        <img
          :src="currentUser?.avatar || '/default-avatar.png'"
          :alt="currentUser?.name || 'You'"
          class="comment-avatar"
        >
        <div class="comment-input-container">
          <textarea
            v-model="newComment.content"
            placeholder="Write a comment..."
            class="comment-input"
            rows="2"
            @keydown.enter.prevent="handleEnterKey"
            :disabled="isSubmitting"
          ></textarea>
          
          <!-- Image Upload -->
          <div class="comment-actions">
            <input
              ref="imageInput"
              type="file"
              accept="image/*"
              @change="handleImageSelect"
              style="display: none"
            >
            <button
              type="button"
              @click="$refs.imageInput.click()"
              class="image-btn"
              :disabled="isSubmitting"
            >
              📷
            </button>
            <button
              @click="addComment"
              class="submit-btn"
              :disabled="!canSubmit || isSubmitting"
            >
              {{ isSubmitting ? 'Posting...' : 'Post' }}
            </button>
          </div>

          <!-- Image Preview -->
          <div v-if="newComment.imagePreview" class="image-preview">
            <img :src="newComment.imagePreview" alt="Preview" />
            <button @click="removeImage" class="remove-image">×</button>
          </div>

          <!-- Error Message -->
          <div v-if="commentError" class="error-message">
            {{ commentError }}
          </div>
        </div>
      </div>

      <!-- Comments List -->
      <div class="comments-list">
        <div
          v-for="comment in comments"
          :key="comment.id"
          class="comment-item"
        >
          <img
            :src="comment.author?.avatar || '/default-avatar.png'"
            :alt="comment.author?.name || 'User'"
            class="comment-avatar"
          >
          <div class="comment-content">
            <div class="comment-header">
              <span class="comment-author">{{ comment.author?.name || 'Anonymous' }}</span>
              <span class="comment-time">{{ formatDate(comment.createdAt) }}</span>
            </div>
            <p class="comment-text">{{ comment.content }}</p>
            <img 
              v-if="comment.image" 
              :src="comment.image" 
              alt="Comment image"
              class="comment-image"
            >
          </div>
        </div>

        <!-- Load More Comments -->
        <button
          v-if="hasMoreComments"
          @click="loadMoreComments"
          class="load-more-btn"
          :disabled="isLoadingComments"
        >
          {{ isLoadingComments ? 'Loading...' : 'Load more comments' }}
        </button>

        <!-- No Comments Message -->
        <div v-if="comments.length === 0 && !isLoadingComments" class="no-comments">
          No comments yet. Be the first to comment!
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'

// Props
const props = defineProps({
  post: {
    type: Object,
    required: true
  },
  currentUser: {
    type: Object,
    required: true
  }
})

// Reactive state
const showComments = ref(false)
const comments = ref([])
const isLoadingComments = ref(false)
const isSubmitting = ref(false)
const commentError = ref('')
const hasMoreComments = ref(false)
const commentsOffset = ref(0)
const commentsLimit = 10

// New comment state
const newComment = ref({
  content: '',
  image: null,
  imagePreview: null
})

// Computed properties
const totalCommentsCount = computed(() => {
  return comments.value.length + (props.post.commentsCount || 0)
})

const canSubmit = computed(() => {
  return newComment.value.content.trim().length > 0 && 
         newComment.value.content.length <= 300
})

// Methods
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

async function toggleComments() {
  showComments.value = !showComments.value
  
  if (showComments.value && comments.value.length === 0) {
    await loadComments()
  }
}

// async function loadComments() {
//   if (isLoadingComments.value) return
  
//   isLoadingComments.value = true
  
//   try {
//     const response = await fetch(`/api/posts/${props.post.id}/comments?offset=${commentsOffset.value}&limit=${commentsLimit}`)
    
//     if (!response.ok) {
//       throw new Error('Failed to load comments')
//     }
    
//     const data = await response.json()
    
//     if (commentsOffset.value === 0) {
//       comments.value = data.comments || []
//     } else {
//       comments.value.push(...(data.comments || []))
//     }
    
//     hasMoreComments.value = data.hasMore || false
//     commentsOffset.value += commentsLimit
    
//   } catch (error) {
//     console.error('Error loading comments:', error)
//     commentError.value = 'Failed to load comments'
//   } finally {
//     isLoadingComments.value = false
//   }
// }

async function loadComments() {
  if (isLoadingComments.value) return
  
  isLoadingComments.value = true
  
  try {
    // Updated to match your backend route
    const response = await fetch(`/post/${props.post.id}/comments`)
    
    if (!response.ok) {
      throw new Error('Failed to load comments')
    }
    
    const comments = await response.json() // Your handler returns comments directly
    
    if (commentsOffset.value === 0) {
      comments.value = comments || []
    } else {
      comments.value.push(...(comments || []))
    }
    
    // You can add pagination logic here if needed
    hasMoreComments.value = false // Set based on your pagination logic
    
  } catch (error) {
    console.error('Error loading comments:', error)
    commentError.value = 'Failed to load comments'
  } finally {
    isLoadingComments.value = false
  }
}
 async function addComment() {
  if (!canSubmit.value || isSubmitting.value) return
  
  // Validate required data
  if (!props.post || !props.post.id) {
    commentError.value = 'Post information is missing'
    return
  }
  
  // if (!props.currentUser || !props.currentUser.id) {
  //   commentError.value = 'User information is missing'
  //   return
  // }
  
  isSubmitting.value = true
  commentError.value = ''
  
  try {
    const formData = new FormData()
    formData.append('content', newComment.value.content.trim())
    formData.append('post_id', props.post.id.toString())
    // formData.append('owner_id', props.currentUser.id.toString())
    formData.append('owner_id', "1")

    
    // Add image if selected
    if (newComment.value.image) {
      formData.append('image', newComment.value.image)
    }
    
    const response = await fetch('/post/comment/new', {
      method: 'POST',
      body: formData
    })
    
    if (!response.ok) {
      const errorData = await response.json()
      throw new Error(errorData.error || 'Failed to add comment')
    }
    
    const result = await response.json()
    
    // Add the new comment to the beginning of the comments array
    if (result.comment) {
      comments.value.unshift(result.comment)
    }
    
    // Reset form
    newComment.value = {
      content: '',
      image: null,
      imagePreview: null
    }
    
    // Clear file input
    if (this.$refs.imageInput) {
      this.$refs.imageInput.value = ''
    }
    
  } catch (error) {
    console.error('Error adding comment:', error)
    commentError.value = error.message || 'Failed to add comment'
  } finally {
    isSubmitting.value = false
  }
}
 
// /dsvsdv///dsv/sd//dsv
async function loadMoreComments() {
  await loadComments()
}

function handleImageSelect(event) {
  const file = event.target.files[0]
  if (!file) return

  // Validate file size (5MB limit)
  if (file.size > 5 * 1024 * 1024) {
    commentError.value = 'Image file too large (max 5MB)'
    return
  }

  // Validate file type
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    commentError.value = 'Only JPEG, PNG and GIF images are allowed'
    return
  }

  newComment.value.image = file
  
  // Create preview
  const reader = new FileReader()
  reader.onload = (e) => {
    newComment.value.imagePreview = e.target.result
  }
  reader.readAsDataURL(file)
  
  commentError.value = ''
}

function removeImage() {
  newComment.value.image = null
  newComment.value.imagePreview = null
  if (this.$refs.imageInput) {
    this.$refs.imageInput.value = ''
  }
}

function handleEnterKey(event) {
  if (!event.shiftKey) {
    addComment()
  }
}

 
</script>

<style scoped>
.post-item {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  margin-bottom: 24px;
  overflow: hidden;
}

.post-header {
  padding: 16px 20px;
  border-bottom: 1px solid #eee;
}

.post-author {
  display: flex;
  align-items: center;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
  margin-right: 12px;
}

.author-info {
  flex: 1;
}

.author-name {
  margin: 0;
  font-size: 14px;
  font-weight: 600;
  color: #333;
}

.post-time {
  font-size: 12px;
  color: #666;
}

.post-content {
  padding: 16px 20px;
}

.post-content p {
  margin: 0 0 12px 0;
  line-height: 1.5;
  color: #333;
}

.post-image {
  width: 100%;
  border-radius: 8px;
  margin-top: 12px;
}

.post-actions {
  display: flex;
  padding: 12px 20px;
  border-top: 1px solid #eee;
  gap: 16px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 8px 12px;
  background: none;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: 14px;
  color: #666;
  transition: all 0.2s;
}

.action-btn:hover {
  background: #f5f5f5;
}

.action-btn.active {
  color: #007bff;
  background: #e3f2fd;
}

.comments-section {
  border-top: 1px solid #eee;
  padding: 16px 20px;
  background: #fafafa;
}

.add-comment {
  display: flex;
  gap: 12px;
  margin-bottom: 20px;
}

.comment-avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  object-fit: cover;
  flex-shrink: 0;
}

.comment-input-container {
  flex: 1;
}

.comment-input {
  width: 100%;
  padding: 12px;
  border: 1px solid #ddd;
  border-radius: 8px;
  resize: vertical;
  min-height: 40px;
  font-family: inherit;
  font-size: 14px;
}

.comment-input:focus {
  outline: none;
  border-color: #007bff;
  box-shadow: 0 0 0 2px rgba(0, 123, 255, 0.25);
}

.comment-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 8px;
}

.image-btn {
  background: none;
  border: none;
  font-size: 18px;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
}

.image-btn:hover {
  background: #eee;
}

.submit-btn {
  padding: 8px 16px;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 6px;
  cursor: pointer;
  font-size: 14px;
  font-weight: 500;
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.submit-btn:hover:not(:disabled) {
  background: #0056b3;
}

.image-preview {
  position: relative;
  margin-top: 8px;
  display: inline-block;
}

.image-preview img {
  max-width: 200px;
  max-height: 150px;
  border-radius: 6px;
}

.remove-image {
  position: absolute;
  top: -8px;
  right: -8px;
  background: #ff4757;
  color: white;
  border: none;
  border-radius: 50%;
  width: 20px;
  height: 20px;
  cursor: pointer;
  font-size: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.error-message {
  color: #dc3545;
  font-size: 12px;
  margin-top: 4px;
}

.comments-list {
  max-height: 400px;
  overflow-y: auto;
}

.comment-item {
  display: flex;
  gap: 12px;
  margin-bottom: 16px;
  padding-bottom: 16px;
  border-bottom: 1px solid #eee;
}

.comment-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.comment-content {
  flex: 1;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 4px;
}

.comment-author {
  font-weight: 600;
  font-size: 13px;
  color: #333;
}

.comment-time {
  font-size: 11px;
  color: #666;
}

.comment-text {
  margin: 0;
  font-size: 14px;
  line-height: 1.4;
  color: #333;
}

.comment-image {
  max-width: 200px;
  border-radius: 6px;
  margin-top: 8px;
}

.load-more-btn {
  width: 100%;
  padding: 8px;
  background: none;
  border: 1px solid #ddd;
  border-radius: 6px;
  cursor: pointer;
  color: #666;
  font-size: 14px;
  margin-top: 12px;
}

.load-more-btn:hover {
  background: #f5f5f5;
}

.no-comments {
  text-align: center;
  color: #666;
  font-style: italic;
  padding: 20px;
}
</style>