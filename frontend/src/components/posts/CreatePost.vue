<template>
  <div class="create-post-container">
    <form @submit.prevent="submitPost" class="create-post-form">
      <div class="form-header">
        <h2>Create New Post</h2>
      </div>
      
      <div class="form-group">
        <label for="title" class="form-label">Title</label>
        <input 
          id="title" 
          v-model="title" 
          type="text"
          class="form-input"
          placeholder="What's your post about?"
          required 
          :disabled="loading"
        />
      </div>
      
      <div class="form-group">
        <label for="content" class="form-label">Content</label>
        <textarea 
          id="content" 
          v-model="content" 
          class="form-textarea"
          placeholder="Share your thoughts..."
          rows="4"
          required
          :disabled="loading"
        ></textarea>
        <div class="character-count">
          {{ content.length }}/500
        </div>
      </div>
      
      <div class="form-actions">
        <button 
          type="button" 
          @click="clearForm"
          class="btn btn-secondary"
          :disabled="loading"
        >
          Clear
        </button>
        
        <button 
          type="submit" 
          class="btn btn-primary"
          :disabled="loading || !canSubmit"
        >
          <span v-if="loading" class="loading-spinner"></span>
          {{ loading ? 'Posting...' : 'Create Post' }}
        </button>
      </div>
      
      <!-- Success Message -->
      <div v-if="successMessage" class="alert alert-success">
        {{ successMessage }}
      </div>
      
      <!-- Error Message -->
      <div v-if="error" class="alert alert-error">
        {{ error }}
      </div>
    </form>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { createPost } from '@/services/api.js'

// Define emits
const emit = defineEmits(['post-created'])

// Form data
const title = ref('')
const content = ref('')
const loading = ref(false)
const error = ref(null)
const successMessage = ref('')

// Computed properties
const canSubmit = computed(() => {
  return title.value.trim().length > 0 && 
         content.value.trim().length > 0 && 
         content.value.length <= 500
})

// Methods
async function submitPost() {
  error.value = null
  successMessage.value = ''





  
  if (!title.value.trim() || !content.value.trim()) {
    error.value = 'Please fill in all fields.'
    return
  }
  
  if (content.value.length > 500) {
    error.value = 'Content must be 500 characters or less.'
    return
  }
  
  loading.value = true
  
  try {




    const postData = {
      image_url :"ana/ghadi/ldar",
      ownerId: 1,
      content: content.value.trim(),
      privacy: 'public', // Or let user choose it
      groupId: null,     // Optional
      // chosenUsersIds: [], // Only used if privacy === 'private'
    }

    const response = await createPost(postData)
    
     clearForm()
    
    // Show success message
    successMessage.value = 'Post created successfully!'
    
     emit('post-created', response)
    
    // Clear success message after 3 seconds
    setTimeout(() => {
      successMessage.value = ''
    }, 3000)
    
  } catch (err) {
    console.error('Error creating post:', err)
    error.value = err.response?.data?.message || 'Failed to create post. Please try again.'
  } finally {
    loading.value = false
  }
}

function clearForm() {
  title.value = ''
  content.value = ''
  error.value = null
}
</script>

<style scoped>
.create-post-container {
  background: white;
  border-radius: 12px;
  box-shadow: 0 2px 8px rgba(0,0,0,0.1);
  margin-bottom: 24px;
  overflow: hidden;
}

.create-post-form {
  padding: 20px;
}

.form-header {
  margin-bottom: 20px;
  padding-bottom: 12px;
  border-bottom: 1px solid #e1e8ed;
}

.form-header h2 {
  margin: 0;
  font-size: 18px;
  font-weight: 600;
  color: #14171a;
}

.form-group {
  margin-bottom: 16px;
  position: relative;
}

.form-label {
  display: block;
  margin-bottom: 6px;
  font-weight: 500;
  color: #14171a;
  font-size: 14px;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 12px 16px;
  border: 2px solid #e1e8ed;
  border-radius: 8px;
  font-size: 16px;
  transition: border-color 0.2s ease;
  font-family: inherit;
  resize: none;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #1da1f2;
  box-shadow: 0 0 0 3px rgba(29, 161, 242, 0.1);
}

.form-input:disabled,
.form-textarea:disabled {
  background-color: #f7f9fa;
  cursor: not-allowed;
  opacity: 0.6;
}

.form-textarea {
  min-height: 100px;
  line-height: 1.5;
}

.character-count {
  position: absolute;
  bottom: -20px;
  right: 0;
  font-size: 12px;
  color: #657786;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 12px;
  margin-top: 24px;
}

.btn {
  padding: 10px 20px;
  border: none;
  border-radius: 6px;
  font-size: 14px;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn:disabled {
  cursor: not-allowed;
  opacity: 0.6;
}

.btn-secondary {
  background-color: #f7f9fa;
  color: #14171a;
  border: 1px solid #e1e8ed;
}

.btn-secondary:hover:not(:disabled) {
  background-color: #e1e8ed;
}

.btn-primary {
  background-color: #1da1f2;
  color: white;
}

.btn-primary:hover:not(:disabled) {
  background-color: #1991db;
}

.loading-spinner {
  width: 16px;
  height: 16px;
  border: 2px solid transparent;
  border-top: 2px solid currentColor;
  border-radius: 50%;
  animation: spin 1s linear infinite;
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.alert {
  padding: 12px 16px;
  border-radius: 6px;
  margin-top: 16px;
  font-size: 14px;
}

.alert-success {
  background-color: #d4edda;
  color: #155724;
  border: 1px solid #c3e6cb;
}

.alert-error {
  background-color: #f8d7da;
  color: #721c24;
  border: 1px solid #f5c6cb;
}

/* Mobile Responsive */
@media (max-width: 768px) {
  .create-post-form {
    padding: 16px;
  }
  
  .form-actions {
    flex-direction: column;
  }
  
  .btn {
    width: 100%;
    justify-content: center;
  }
}
</style>