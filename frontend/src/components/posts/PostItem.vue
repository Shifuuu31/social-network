<template>
  <div class="post-item">
    <div class="post-header">
      <div class="post-author">
        <img
          :src="post.avatar_url ? getImageUrl(post.avatar_url) : '/default-avatar.png'"
          :alt="post.owner || 'User'"
          class="avatar"
          @error="console.log('Avatar failed to load for post:', post.id, 'avatar_url:', post.avatar_url)"
          @load="console.log('Avatar loaded successfully for post:', post.id, 'avatar_url:', post.avatar_url)"
        >
        <div class="author-info">
          <h4 class="author-name">{{ post.owner || 'Anonymous' }}</h4>
          <span class="post-time">{{ formatDate(post.created_at) }}</span>
        </div>
      </div>
    </div>

        <div class="post-content">
            <p v-if="post.content">{{ post.content }}</p>
            <div v-if="post.image_url" class="post-image-container">
                <img
                    :src="getImageUrl(post.image_url)"
                    :alt="post.content || 'Post image'"
                    class="post-image"
                    @error="handleImageError"
                    @click="togglePreview"
                >
            </div>
            
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
            :alt="comment.owner || 'User'"
            class="comment-avatar"
          >
          <div class="comment-content">
            <div class="comment-header">
              <span class="comment-author">{{ comment.owner || 'Anonymous' }}</span>
              <span class="comment-time">{{ formatDate(comment.created_at) }}</span>
            </div>

            <p class="comment-text">{{ comment.content }}</p>
            <div v-if="comment.image_url || comment.image" class="comment-image-container">
              <img 
                :src="getImageUrl(comment.image_url || comment.image)" 
                alt="Comment image"
                class="comment-image"
                @error="handleImageError"
                @load="console.log('Comment image loaded:', comment.image_url || comment.image)"
              >
            </div>
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
  <div
                class="img-preview"
                v-if="isOpenPreview"
                @click="togglePreview"
            >
                <span p="togglePreview" >
<svg xmlns="http://www.w3.org/2000/svg" height="60px" viewBox="0 -960 960 960" width="60px" fill="#e3e3e3"><path d="M440-440v240h-60v-180H200v-60h240Zm140-320v180h180v60H520v-240h60Z"/></svg>
</span>
                <img
                    :src="getImageUrl(post.image_url)"
                    :alt="post.content || 'Post image'"
                    @error="handleImageError"
                    @click.stop="togglePreview"
                >
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

// Debug: Log post data when component mounts
onMounted(() => {
  console.log('PostItem mounted for post:', props.post.id, 'avatar_url:', props.post.avatar_url)
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
let isOpenPreview = ref(false)



const togglePreview = () => {
  console.log(isOpenPreview)
  
isOpenPreview.value =  !isOpenPreview.value

}

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
    const response = await fetch(`/post/${props.post.id}/comments`, {
      credentials: 'include'
    })
    console.log("response", response)
    
    if (!response.ok) {
      throw new Error('Failed to load comments')
    }
    
    const data = await response.json()
    console.log('Loaded comments data:', data)
    
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
    formData.append('owner_id', props.currentUser?.id?.toString() || "1") // Use actual user ID
    
    // Add image if selected
    if (newComment.value.image) {
      formData.append('image', newComment.value.image)
    }
    
    // Fixed: Use correct URL with post_id in path
    const response = await fetch(`/post/${props.post.id}/comments/new`, {
      method: 'POST',
      body: formData,
      credentials: 'include'
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

// Image handling functions
function getImageUrl(imagePath) {
  if (!imagePath) return ''
  
  console.log('getImageUrl called with:', imagePath)
  
  // If it's already a full URL, return as is
  if (imagePath.startsWith('http://') || imagePath.startsWith('https://')) {
    return imagePath
  }
  
  // Extract filename from path (e.g., "uploads/filename.jpg" -> "filename.jpg")
  const filename = imagePath.split('/').pop()
  
  // Construct the full URL to the backend server
  const fullUrl = `http://localhost:8080/images/${filename}`
  console.log('Generated image URL:', fullUrl)
  return fullUrl
}

function handleImageError(event) {
  console.error('Failed to load image:', event.target.src)
  // You can set a fallback image here
  event.target.style.display = 'none'
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
  background: var(--background-glass);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow);
  margin-bottom: 40px; /* Increased gap between posts */
  padding: 28px 24px 18px 24px;
  transition: box-shadow var(--transition), background var(--transition);
  backdrop-filter: blur(8px);
  position: relative;
}

.post-header {
  display: flex;
  align-items: center;
  margin-bottom: 12px;
}

.post-author {
  display: flex;
  align-items: center;
  gap: 14px;
}

.avatar {
  width: 48px;
  height: 48px;
  border-radius: 50%;
  object-fit: cover;
  border: 2px solid var(--primary);
  box-shadow: 0 2px 8px rgba(37,99,235,0.10);
}

.author-info {
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.author-name {
  font-size: 1.08rem;
  font-weight: 700;
  color: var(--primary-dark);
  margin: 0;
}

.post-time {
  font-size: 0.92rem;
  color: var(--text-light);
  opacity: 0.7;
}

.post-content {
            margin: 16px 0 10px 0;
            font-size: 1.08rem;
            color: var(--text-main);
            word-break: break-word;
        }

        .post-image-container {
            margin-top: 16px;
            display: flex;
            justify-content: center;
            align-items: center;
            width: 100%;
            max-height: 500px;
            overflow: hidden;
            border-radius: 16px;
            box-shadow: 0 4px 20px rgba(37,99,235,0.15);
            border: 1px solid rgba(37,99,235,0.1);
            background: #f8fafc;
        }

        .post-image {
            width: 100%;
            height: 100%;
            object-fit: cover;
            border-radius: 16px;
            transition: all 0.3s ease-in-out;
            cursor: pointer;
            opacity: 0;
            animation: fadeIn 0.5s ease-in-out forwards;
        }

        .img-preview {
            background: rgba(0, 0, 0, 0.9);
            position: fixed;
            top: 0;
            left: 0;
            width: 100vw;
            height: 100vh;
            display: flex;
            align-items: center;
            justify-content: center;
            z-index: 9999;
        }

        .img-preview span {
            position: absolute;
            top: 20px;
            right: 30px;
            font-size: 2rem;
            color: white;
            cursor: pointer;
            opacity: 50%;
            z-index: 10000;
            padding: 0.3em 0.6em;
            border-radius: 8px;
            transition: background 0.3s;
        }

        .img-preview span:hover {
            opacity: 100%
        }

        .img-preview img {
            max-width: 80%;
            max-height: 80%;
            object-fit: contain;
            border-radius: 12px;
            box-shadow: 0 0 20px rgba(0, 0, 0, 0.5);
            cursor: pointer;
        }

        @keyframes fadeIn {
            from { 
                opacity: 0; 
                transform: scale(0.95); 
            }
            to { 
                opacity: 1; 
                transform: scale(1); 
            }
        }

        .post-image:hover {
            transform: scale(1.02);
            box-shadow: 0 8px 30px rgba(37,99,235,0.25);
        }

.post-actions {
  display: flex;
  align-items: center;
  gap: 18px;
  margin-top: 8px;
  margin-bottom: 2px;
}

.action-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  background: #f4f6fb;
  color: var(--primary-dark);
  border: none;
  border-radius: 8px;
  padding: 7px 18px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition), color var(--transition), box-shadow var(--transition);
  box-shadow: 0 1px 4px rgba(37,99,235,0.06);
  outline: none;
}
.action-btn .icon {
  font-size: 1.15em;
}
.action-btn:hover, .action-btn.active {
  background: linear-gradient(90deg, var(--primary), var(--accent));
  color: #fff;
  box-shadow: 0 2px 8px rgba(37,99,235,0.13);
}

.count {
  font-size: 1em;
  font-weight: 700;
  margin-left: 2px;
}

.comments-section {
  margin-top: 18px;
  background: rgba(244,246,251,0.85);
  border-radius: 14px;
  padding: 18px 14px 10px 14px;
  box-shadow: 0 2px 8px rgba(37,99,235,0.06);
}

.add-comment {
  display: flex;
  align-items: flex-start;
  gap: 12px;
  margin-bottom: 14px;
}

.comment-avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  object-fit: cover;
  border: 1.5px solid var(--primary);
}

.comment-input-container {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.comment-input {
  width: 100%;
  border-radius: 8px;
  border: 1.5px solid #e5e7eb;
  padding: 10px 14px;
  font-size: 1rem;
  background: #ffffff;
  color: #1f2937;
  transition: border var(--transition), background var(--transition);
  resize: none;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
}

.comment-input::placeholder {
  color: #6b7280;
  opacity: 0.8;
}
.comment-input:focus {
  border: 1.5px solid var(--primary);
  background: #ffffff;
  color: #1f2937;
  outline: none;
  box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.1);
}

.comment-actions {
  display: flex;
  align-items: center;
  gap: 10px;
  margin-top: 2px;
}

.image-btn {
  background: #e5e7eb;
  color: var(--primary-dark);
  border: none;
  border-radius: 6px;
  padding: 6px 10px;
  font-size: 1.1em;
  cursor: pointer;
  transition: background var(--transition), color var(--transition);
}
.image-btn:hover {
  background: #f3f4f6;
  color: var(--primary);
}

.submit-btn {
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 7px 18px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  transition: background var(--transition), box-shadow var(--transition);
  box-shadow: 0 2px 8px rgba(37,99,235,0.08);
}
.submit-btn:hover {
  background: var(--primary-dark);
  box-shadow: 0 4px 16px rgba(37,99,235,0.12);
}
.submit-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.image-preview {
  margin-top: 8px;
  position: relative;
  width: 120px;
  height: 120px;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 3px 15px rgba(37,99,235,0.12);
  border: 1px solid rgba(37,99,235,0.08);
  background: #f8fafc;
}
.image-preview img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  border-radius: 12px;
  transition: all 0.3s ease-in-out;
}
.remove-image {
  position: absolute;
  top: 2px;
  right: 2px;
  background: #ef4444;
  color: #fff;
  border: none;
  border-radius: 50%;
  width: 22px;
  height: 22px;
  font-size: 1.1em;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 1px 4px rgba(37,99,235,0.10);
}

.error-message {
  color: #b91c1c;
  background: #fef2f2;
  border: 1.5px solid #ef4444;
  border-radius: 8px;
  padding: 8px 12px;
  margin-top: 6px;
  font-size: 0.98em;
  font-weight: 500;
}

.comments-list {
  margin-top: 10px;
  display: flex;
  flex-direction: column;
  gap: 14px;
}

.comment-item {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  background: #ffffff;
  border-radius: 10px;
  padding: 12px 16px;
  box-shadow: 0 2px 8px rgba(37,99,235,0.08);
  border: 1px solid rgba(37,99,235,0.05);
}

.comment-content {
  flex: 1;
  display: flex;
  flex-direction: column;
  gap: 2px;
}

.comment-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 0.98em;
  color: #2563eb;
  font-weight: 600;
}

.comment-author {
  color: #2563eb;
  font-weight: 700;
}

.comment-time {
  color: #6b7280;
  opacity: 0.8;
  font-size: 0.92em;
}

.comment-text {
  color: #1f2937;
  font-size: 1em;
  margin: 2px 0 0 0;
  line-height: 1.5;
  font-weight: 400;
}

.comment-image-container {
  margin-top: 8px;
  display: flex;
  justify-content: center;
  align-items: center;
  width: 100%;
  max-height: 300px;
  overflow: visible;
  border-radius: 12px;
  box-shadow: 0 3px 15px rgba(37,99,235,0.12);
  border: 1px solid rgba(37,99,235,0.08);
  background: #f8fafc;
  padding: 8px;
}

.comment-image {
  width: 100%;
  height: auto;
  max-height: 200px;
  object-fit: contain;
  border-radius: 12px;
  transition: all 0.3s ease-in-out;
  cursor: pointer;
  opacity: 0;
  animation: fadeIn 0.5s ease-in-out forwards;
}

.comment-image:hover {
  transform: scale(1.03);
  box-shadow: 0 6px 25px rgba(37,99,235,0.2);
}

.load-more-btn {
  margin: 12px auto 0 auto;
  display: block;
  background: var(--primary);
  color: #fff;
  border: none;
  border-radius: 8px;
  padding: 8px 22px;
  font-size: 1rem;
  font-weight: 600;
  cursor: pointer;
  box-shadow: 0 2px 8px rgba(37,99,235,0.08);
  transition: background var(--transition), box-shadow var(--transition);
}
.load-more-btn:hover {
  background: var(--primary-dark);
  box-shadow: 0 4px 16px rgba(37,99,235,0.12);
}
.load-more-btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.no-comments {
  text-align: center;
  color: var(--text-light);
  font-size: 1em;
  margin-top: 10px;
}

@media (max-width: 768px) {
  .post-item {
    padding: 12px 4px 8px 4px;
    border-radius: 12px;
  }
  .avatar {
    width: 38px;
    height: 38px;
  }
  .post-image-container {
    max-height: 300px;
    border-radius: 12px;
  }
  .post-image {
    border-radius: 12px;
  }
  .comment-image-container {
    max-height: 200px;
    border-radius: 10px;
    padding: 6px;
  }
  .comment-image {
    border-radius: 10px;
    max-height: 180px;
  }
  .comments-section {
    padding: 10px 2px 6px 2px;
    border-radius: 8px;
  }
  .image-preview {
    width: 100px;
    height: 100px;
  }
}
</style>

 