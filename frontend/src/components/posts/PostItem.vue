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

async function toggleComments() {
  showComments.value = !showComments.value
  
  if (showComments.value && comments.value.length === 0) {
    await loadComments()
  }
}

 async function loadComments() {
  if (isLoadingComments.value) return
  
  isLoadingComments.value = true
  
  try {
    // Fixed: Match your backend route structure
    const response = await fetch(`/post/${props.post.id}/comments`)
    
    if (!response.ok) {
      throw new Error('Failed to load comments')
    }
    
    const data = await response.json()
    
    if (commentsOffset.value === 0) {
      comments.value = data || []
    } else {
      comments.value.push(...(data || []))
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
  
  isSubmitting.value = true
  commentError.value = ''
  
  try {
    const formData = new FormData()
    formData.append('content', newComment.value.content.trim())
    // Removed: Don't send post_id in FormData since it's in the URL
    // formData.append('post_id', props.post.id.toString())
    formData.append('owner_id', "1") // You should use actual user ID
    
    // Add image if selected
    if (newComment.value.image) {
      formData.append('image', newComment.value.image)
    }
    
    // Fixed: Use correct URL with post_id in path
    const response = await fetch(`/post/${props.post.id}/comments/new`, {
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
    
    // Clear file input - Fixed the this reference issue
    const imageInput = document.querySelector('input[type="file"]')
    if (imageInput) {
      imageInput.value = ''
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

<style>
/* CSS Variables for consistency */
:root {
  --primary-color: #007bff;
  --primary-hover: #0056b3;
  --success-color: #28a745;
  --danger-color: #dc3545;
  --warning-color: #ffc107;
  --text-primary: #333;
  --text-secondary: #666;
  --text-muted: #999;
  --border-color: #eee;
  --border-color-focus: #007bff;
  --background-white: #ffffff;
  --background-light: #fafafa;
  --background-hover: #f5f5f5;
  --shadow-light: 0 2px 8px rgba(0, 0, 0, 0.1);
  --shadow-focus: 0 0 0 2px rgba(0, 123, 255, 0.25);
  --radius-small: 6px;
  --radius-medium: 8px;
  --radius-large: 12px;
  --radius-circle: 50%;
  --spacing-xs: 4px;
  --spacing-sm: 8px;
  --spacing-md: 12px;
  --spacing-lg: 16px;
  --spacing-xl: 20px;
  --spacing-xxl: 24px;
  --font-size-xs: 11px;
  --font-size-sm: 12px;
  --font-size-base: 14px;
  --font-size-lg: 16px;
  --font-size-xl: 18px;
}

/* Main Post Container */
.post-item {
  background: var(--background-white);
  border-radius: var(--radius-large);
  box-shadow: var(--shadow-light);
  margin-bottom: var(--spacing-xxl);
  overflow: hidden;
  transition: box-shadow 0.3s ease;
}

.post-item:hover {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.15);
}

/* Post Header */
.post-header {
  padding: var(--spacing-lg) var(--spacing-xl);
  border-bottom: 1px solid var(--border-color);
}

.post-author {
  display: flex;
  align-items: center;
}

.avatar {
  width: 40px;
  height: 40px;
  border-radius: var(--radius-circle);
  object-fit: cover;
  margin-right: var(--spacing-md);
  border: 2px solid transparent;
  transition: border-color 0.2s ease;
}

.avatar:hover {
  border-color: var(--primary-color);
}

.author-info {
  flex: 1;
  min-width: 0; /* Prevent overflow */
}

.author-name {
  margin: 0;
  font-size: var(--font-size-base);
  font-weight: 600;
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.post-time {
  font-size: var(--font-size-sm);
  color: var(--text-secondary);
  display: block;
  margin-top: 2px;
}

/* Post Content */
.post-content {
  padding: var(--spacing-lg) var(--spacing-xl);
}

.post-content p {
  margin: 0 0 var(--spacing-md) 0;
  line-height: 1.5;
  color: var(--text-primary);
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.post-content p:last-child {
  margin-bottom: 0;
}

.post-image {
  width: 100%;
  max-width: 100%;
  height: auto;
  border-radius: var(--radius-medium);
  margin-top: var(--spacing-md);
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.2s ease;
}

.post-image:hover {
  transform: scale(1.02);
}

/* Post Actions */
.post-actions {
  display: flex;
  padding: var(--spacing-md) var(--spacing-xl);
  border-top: 1px solid var(--border-color);
  gap: var(--spacing-lg);
}

.action-btn {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  padding: var(--spacing-sm) var(--spacing-md);
  background: none;
  border: none;
  border-radius: 20px;
  cursor: pointer;
  font-size: var(--font-size-base);
  color: var(--text-secondary);
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;
}

.action-btn:hover {
  background: var(--background-hover);
  transform: translateY(-1px);
}

.action-btn:focus {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

.action-btn:active {
  transform: translateY(0);
}

.action-btn.active {
  color: var(--primary-color);
  background: #e3f2fd;
  font-weight: 500;
}

.action-btn .icon {
  font-size: var(--font-size-lg);
  transition: transform 0.2s ease;
}

.action-btn:hover .icon {
  transform: scale(1.1);
}

.action-btn .count {
  font-weight: 500;
  min-width: 20px;
  text-align: center;
}

/* Comments Section */
.comments-section {
  border-top: 1px solid var(--border-color);
  padding: var(--spacing-lg) var(--spacing-xl);
  background: var(--background-light);
  animation: slideDown 0.3s ease;
}

@keyframes slideDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Add Comment Form */
.add-comment {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-xl);
}

.comment-avatar {
  width: 32px;
  height: 32px;
  border-radius: var(--radius-circle);
  object-fit: cover;
  flex-shrink: 0;
  border: 2px solid transparent;
  transition: border-color 0.2s ease;
}

.comment-avatar:hover {
  border-color: var(--primary-color);
}

.comment-input-container {
  flex: 1;
  min-width: 0;
}

.comment-input {
  width: 100%;
  padding: var(--spacing-md);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-medium);
  resize: vertical;
  min-height: 40px;
  max-height: 120px;
  font-family: inherit;
  font-size: var(--font-size-base);
  line-height: 1.4;
  transition: all 0.2s ease;
  background: var(--background-white);
}

.comment-input:focus {
  outline: none;
  border-color: var(--border-color-focus);
  box-shadow: var(--shadow-focus);
  background: var(--background-white);
}

.comment-input::placeholder {
  color: var(--text-muted);
}

.comment-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: var(--spacing-sm);
}

.image-btn {
  background: none;
  border: none;
  font-size: var(--font-size-xl);
  cursor: pointer;
  padding: var(--spacing-xs);
  border-radius: var(--spacing-xs);
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 32px;
  height: 32px;
}

.image-btn:hover {
  background: var(--background-hover);
  transform: scale(1.1);
}

.image-btn:focus {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

.image-btn:disabled {
  opacity: 0.5;
  cursor: not-allowed;
  transform: none;
}

.submit-btn {
  padding: var(--spacing-sm) var(--spacing-lg);
  background: var(--primary-color);
  color: var(--background-white);
  border: none;
  border-radius: var(--radius-small);
  cursor: pointer;
  font-size: var(--font-size-base);
  font-weight: 500;
  transition: all 0.2s ease;
  position: relative;
  overflow: hidden;
  min-width: 80px;
}

.submit-btn:hover:not(:disabled) {
  background: var(--primary-hover);
  transform: translateY(-1px);
  box-shadow: 0 4px 8px rgba(0, 123, 255, 0.3);
}

.submit-btn:active:not(:disabled) {
  transform: translateY(0);
}

.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
  box-shadow: none;
}

.submit-btn:focus {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

/* Image Preview */
.image-preview {
  position: relative;
  margin-top: var(--spacing-sm);
  display: inline-block;
  border-radius: var(--radius-small);
  overflow: hidden;
  box-shadow: var(--shadow-light);
}

.image-preview img {
  max-width: 200px;
  max-height: 150px;
  width: auto;
  height: auto;
  border-radius: var(--radius-small);
  object-fit: cover;
}

.remove-image {
  position: absolute;
  top: -8px;
  right: -8px;
  background: var(--danger-color);
  color: var(--background-white);
  border: none;
  border-radius: var(--radius-circle);
  width: 20px;
  height: 20px;
  cursor: pointer;
  font-size: var(--font-size-sm);
  display: flex;
  align-items: center;
  justify-content: center;
  transition: all 0.2s ease;
  box-shadow: var(--shadow-light);
}

.remove-image:hover {
  background: #c82333;
  transform: scale(1.1);
}

.remove-image:focus {
  outline: 2px solid var(--danger-color);
  outline-offset: 2px;
}

/* Error Message */
.error-message {
  color: var(--danger-color);
  font-size: var(--font-size-sm);
  margin-top: var(--spacing-xs);
  padding: var(--spacing-xs) var(--spacing-sm);
  background: #f8d7da;
  border: 1px solid #f5c6cb;
  border-radius: var(--radius-small);
  animation: shake 0.5s ease-in-out;
}

@keyframes shake {
  0%, 100% { transform: translateX(0); }
  25% { transform: translateX(-5px); }
  75% { transform: translateX(5px); }
}

/* Comments List */
.comments-list {
  max-height: 400px;
  overflow-y: auto;
  scrollbar-width: thin;
  scrollbar-color: var(--border-color) transparent;
}

.comments-list::-webkit-scrollbar {
  width: 6px;
}

.comments-list::-webkit-scrollbar-track {
  background: transparent;
}

.comments-list::-webkit-scrollbar-thumb {
  background: var(--border-color);
  border-radius: 3px;
}

.comments-list::-webkit-scrollbar-thumb:hover {
  background: var(--text-secondary);
}

.comment-item {
  display: flex;
  gap: var(--spacing-md);
  margin-bottom: var(--spacing-lg);
  padding-bottom: var(--spacing-lg);
  border-bottom: 1px solid var(--border-color);
  animation: fadeIn 0.3s ease;
}

@keyframes fadeIn {
  from {
    opacity: 0;
    transform: translateY(10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.comment-item:last-child {
  border-bottom: none;
  margin-bottom: 0;
  padding-bottom: 0;
}

.comment-content {
  flex: 1;
  min-width: 0;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: var(--spacing-sm);
  margin-bottom: var(--spacing-xs);
}

.comment-author {
  font-weight: 600;
  font-size: var(--font-size-sm);
  color: var(--text-primary);
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.comment-time {
  font-size: var(--font-size-xs);
  color: var(--text-secondary);
  flex-shrink: 0;
}

.comment-text {
  margin: 0;
  font-size: var(--font-size-base);
  line-height: 1.4;
  color: var(--text-primary);
  word-wrap: break-word;
  overflow-wrap: break-word;
}

.comment-image {
  max-width: 200px;
  width: auto;
  height: auto;
  border-radius: var(--radius-small);
  margin-top: var(--spacing-sm);
  object-fit: cover;
  cursor: pointer;
  transition: transform 0.2s ease;
  box-shadow: var(--shadow-light);
}

.comment-image:hover {
  transform: scale(1.02);
}

/* Load More Button */
.load-more-btn {
  width: 100%;
  padding: var(--spacing-sm);
  background: var(--background-white);
  border: 1px solid var(--border-color);
  border-radius: var(--radius-small);
  cursor: pointer;
  color: var(--text-secondary);
  font-size: var(--font-size-base);
  margin-top: var(--spacing-md);
  transition: all 0.2s ease;
  font-weight: 500;
}

.load-more-btn:hover:not(:disabled) {
  background: var(--background-hover);
  border-color: var(--primary-color);
  color: var(--primary-color);
  transform: translateY(-1px);
}

.load-more-btn:focus {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

.load-more-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
  transform: none;
}

/* No Comments Message */
.no-comments {
  text-align: center;
  color: var(--text-secondary);
  font-style: italic;
  padding: var(--spacing-xl);
  background: var(--background-white);
  border-radius: var(--radius-medium);
  border: 1px dashed var(--border-color);
}

/* Responsive Design */
@media (max-width: 768px) {
  .post-item {
    margin-bottom: var(--spacing-lg);
    border-radius: var(--radius-medium);
  }
  
  .post-header,
  .post-content,
  .post-actions,
  .comments-section {
    padding: var(--spacing-md) var(--spacing-lg);
  }
  
  .avatar {
    width: 36px;
    height: 36px;
  }
  
  .comment-avatar {
    width: 28px;
    height: 28px;
  }
  
  .post-actions {
    gap: var(--spacing-md);
  }
  
  .action-btn {
    padding: var(--spacing-sm) var(--spacing-sm);
    gap: var(--spacing-xs);
  }
  
  .add-comment {
    gap: var(--spacing-sm);
  }
  
  .comment-input {
    padding: var(--spacing-sm);
    font-size: var(--font-size-sm);
  }
  
  .comments-list {
    max-height: 300px;
  }
  
  .image-preview img,
  .comment-image {
    max-width: 150px;
    max-height: 100px;
  }
}

@media (max-width: 480px) {
  .post-header,
  .post-content,
  .post-actions,
  .comments-section {
    padding: var(--spacing-sm) var(--spacing-md);
  }
  
  .avatar {
    width: 32px;
    height: 32px;
  }
  
  .comment-avatar {
    width: 24px;
    height: 24px;
  }
  
  .author-name {
    font-size: var(--font-size-sm);
  }
  
  .post-time {
    font-size: var(--font-size-xs);
  }
  
  .action-btn {
    padding: var(--spacing-xs) var(--spacing-sm);
    font-size: var(--font-size-sm);
  }
  
  .submit-btn {
    padding: var(--spacing-xs) var(--spacing-sm);
    font-size: var(--font-size-sm);
    min-width: 60px;
  }
  
  .image-preview img,
  .comment-image {
    max-width: 120px;
    max-height: 80px;
  }
}

/* Dark Mode Support */
@media (prefers-color-scheme: dark) {
  :root {
    --text-primary: #e1e1e1;
    --text-secondary: #b0b0b0;
    --text-muted: #888;
    --background-white: #1a1a1a;
    --background-light: #2d2d2d;
    --background-hover: #333;
    --border-color: #404040;
    --shadow-light: 0 2px 8px rgba(0, 0, 0, 0.3);
  }
  
  .comment-input {
    background: var(--background-white);
    color: var(--text-primary);
  }
  
  .comment-input::placeholder {
    color: var(--text-muted);
  }
  
  .error-message {
    background: #3d1a1a;
    border-color: #5c1e1e;
  }
}

/* Loading States */
.loading {
  opacity: 0.7;
  pointer-events: none;
}

.loading::after {
  content: '';
  position: absolute;
  top: 50%;
  left: 50%;
  width: 20px;
  height: 20px;
  margin: -10px 0 0 -10px;
  border: 2px solid transparent;
  border-top: 2px solid var(--primary-color);
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

/* Accessibility Improvements */
@media (prefers-reduced-motion: reduce) {
  *,
  *::before,
  *::after {
    animation-duration: 0.01ms !important;
    animation-iteration-count: 1 !important;
    transition-duration: 0.01ms !important;
  }
}

/* Focus visible for better accessibility */
.action-btn:focus-visible,
.submit-btn:focus-visible,
.load-more-btn:focus-visible,
.image-btn:focus-visible {
  outline: 2px solid var(--primary-color);
  outline-offset: 2px;
}

/* High contrast mode support */
@media (prefers-contrast: high) {
  .post-item {
    border: 1px solid var(--text-primary);
  }
  
  .action-btn,
  .submit-btn,
  .load-more-btn {
    border: 1px solid var(--text-primary);
  }
}
</style>

 