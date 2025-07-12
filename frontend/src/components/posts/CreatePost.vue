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
      
      <!-- Image Upload Section -->
      <div class="form-group">
        <label class="form-label">Image (Optional)</label>
        <div class="image-upload-container">
          <input 
            ref="fileInput"
            type="file" 
            @change="handleImageSelect"
            accept="image/jpeg,image/jpg,image/png,image/gif"
            class="file-input"
            :disabled="loading"
          />
          
          <div class="upload-area" @click="triggerFileInput">
            <div v-if="!selectedImage" class="upload-placeholder">
              <div class="upload-icon">📷</div>
              <p>Click to add an image</p>
              <span class="upload-hint">Supports JPEG, PNG, GIF (max 5MB)</span>
            </div>
            
            <div v-else class="image-preview">
              <img :src="selectedImage" alt="Preview" class="preview-image" />
              <button 
                type="button" 
                @click="removeImage" 
                class="remove-image-btn"
                :disabled="loading"
              >
                ✕
              </button>
            </div>
          </div>
        </div>
      </div>
      
      <div class="form-group">
        <label for="privacy" class="form-label">Privacy</label>
        <select 
          id="privacy" 
          v-model="privacy" 
          class="form-input"
          :disabled="loading"
        >
          <option value="public">🌍 Public - Everyone can see this post</option>
          <option value="almost_private">👥 Almost Private - Only followers can see this post</option>
          <option value="private">🔒 Private - Only your close friends can see this post</option>
        </select>
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
import { useAuth } from '@/composables/useAuth.js'

// Define emits
const emit = defineEmits(['post-created'])

// Get auth context
const { user } = useAuth()

// Form data
const title = ref('')
const content = ref('')
const privacy = ref('public')
const loading = ref(false)
const error = ref(null)
const successMessage = ref('')

// Image handling
const fileInput = ref(null)
const selectedImage = ref(null)
const selectedFile = ref(null)

// Computed properties
const canSubmit = computed(() => {
  return title.value.trim().length > 0 && 
         content.value.trim().length > 0 && 
         content.value.length <= 500
})

// Methods
function triggerFileInput() {
  if (!loading.value) {
    fileInput.value.click()
  }
}

function handleImageSelect(event) {
  const file = event.target.files[0]
  if (!file) return

  // Validate file type
  const allowedTypes = ['image/jpeg', 'image/jpg', 'image/png', 'image/gif']
  if (!allowedTypes.includes(file.type)) {
    error.value = 'Please select a valid image file (JPEG, PNG, or GIF).'
    return
  }

  // Validate file size (5MB)
  const maxSize = 5 * 1024 * 1024 // 5MB
  if (file.size > maxSize) {
    error.value = 'Image file is too large. Please select an image under 5MB.'
    return
  }

  // Create preview
  const reader = new FileReader()
  reader.onload = (e) => {
    selectedImage.value = e.target.result
    selectedFile.value = file
    error.value = null
  }
  reader.readAsDataURL(file)
}

function removeImage() {
  selectedImage.value = null
  selectedFile.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
}

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

  // Check if user is authenticated
  if (!user.value || !user.value.id) {
    error.value = 'Please log in to create a post.'
    return
  }
  
  loading.value = true
  
  try {
    // Create FormData for multipart upload
    const formData = new FormData()
    formData.append('content', content.value.trim())
    formData.append('privacy', privacy.value)
    formData.append('owner_id', user.value.id.toString()) // Use actual user ID
    formData.append('group_id', '') // Optional
    
    // Add image if selected
    if (selectedFile.value) {
      formData.append('image', selectedFile.value)
    }

    const response = await createPost(formData)
    
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
  privacy.value = 'public'
  selectedImage.value = null
  selectedFile.value = null
  if (fileInput.value) {
    fileInput.value.value = ''
  }
  error.value = null
}
</script>

<style scoped>
.create-post-container {
  background: var(--background-glass);
  border-radius: var(--border-radius);
  box-shadow: var(--shadow);
  margin-bottom: 32px;
  overflow: hidden;
  backdrop-filter: blur(8px);
  transition: box-shadow var(--transition);
}

.create-post-form {
  padding: 32px 28px 24px 28px;
  display: flex;
  flex-direction: column;
  gap: 18px;
}

.form-header {
  margin-bottom: 12px;
  padding-bottom: 10px;
  border-bottom: 1.5px solid #e5e7eb;
}

.form-header h2 {
  margin: 0;
  font-size: 1.3rem;
  font-weight: 800;
  color: var(--primary);
  letter-spacing: -0.5px;
}

.form-group {
  margin-bottom: 0;
  position: relative;
  display: flex;
  flex-direction: column;
  gap: 6px;
}

.form-label {
  font-weight: 600;
  color: var(--primary-dark);
  font-size: 1rem;
  margin-bottom: 2px;
}

.form-input,
.form-textarea {
  width: 100%;
  padding: 12px 16px;
  border: 1.5px solid #e5e7eb;
  border-radius: 10px;
  background: #f4f6fb;
  color: var(--text-main);
  font-size: 1rem;
  font-family: inherit;
  transition: border var(--transition), background var(--transition);
  resize: none;
}
.form-input:focus,
.form-textarea:focus {
  border: 1.5px solid var(--primary);
  background: #fff;
  outline: none;
}

.character-count {
  position: absolute;
  bottom: 8px;
  right: 16px;
  font-size: 0.85em;
  color: var(--accent);
  opacity: 0.7;
}

/* Image Upload Styles */
.image-upload-container {
  position: relative;
}

.file-input {
  position: absolute;
  opacity: 0;
  width: 0;
  height: 0;
  pointer-events: none;
}

.upload-area {
  border: 2px dashed #d1d5db;
  border-radius: 12px;
  padding: 24px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s ease;
  background: #f9fafb;
  min-height: 120px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.upload-area:hover {
  border-color: var(--primary);
  background: #f0f4ff;
}

.upload-placeholder {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 8px;
}

.upload-icon {
  font-size: 2rem;
  margin-bottom: 8px;
}

.upload-placeholder p {
  margin: 0;
  font-weight: 600;
  color: var(--primary-dark);
  font-size: 1rem;
}

.upload-hint {
  font-size: 0.85rem;
  color: #6b7280;
  opacity: 0.8;
}

.image-preview {
  position: relative;
  width: 100%;
  max-width: 300px;
  margin: 0 auto;
}

.preview-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.remove-image-btn {
  position: absolute;
  top: -8px;
  right: -8px;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background: #ef4444;
  color: white;
  border: none;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 14px;
  font-weight: bold;
  box-shadow: 0 2px 8px rgba(239, 68, 68, 0.3);
  transition: all 0.2s ease;
}

.remove-image-btn:hover {
  background: #dc2626;
  transform: scale(1.1);
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 16px;
  margin-top: 12px;
}

.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  padding: 0.75em 1.8em;
  font-size: 1.1rem;
  font-weight: 600;
  border-radius: 12px;
  border: none;
  cursor: pointer;
  transition: background 0.3s ease, transform 0.15s ease, box-shadow 0.3s ease;
  font-family: inherit;
  box-shadow: 0 3px 12px rgba(0, 0, 0, 0.07);
}

.btn-primary {
  background: linear-gradient(135deg, #6366f1, #8b5cf6);
  color: #ffffff;
  box-shadow: 0 5px 18px rgba(124, 58, 237, 0.3);
}

.btn-primary:hover,
.btn-primary:focus {
  background: linear-gradient(135deg, #4f46e5, #7c3aed);
  transform: translateY(-2px) scale(1.03);
  box-shadow: 0 8px 28px rgba(124, 58, 237, 0.4);
}

.btn-secondary {
  background: #e5e7eb;
  color: var(--primary-dark);
  box-shadow: none;
  border: 1.5px solid #e5e7eb;
}
.btn-secondary:hover, .btn-secondary:focus {
  background: #f3f4f6;
  color: var(--primary);
  border: 1.5px solid var(--primary);
}

.alert {
  margin-top: 18px;
  padding: 12px 18px;
  border-radius: 8px;
  font-size: 1rem;
  font-weight: 500;
  box-shadow: 0 2px 8px rgba(37,99,235,0.08);
}
.alert-success {
  background: #e0f7e9;
  color: #15803d;
  border: 1.5px solid #22c55e;
}
.alert-error {
  background: #fef2f2;
  color: #b91c1c;
  border: 1.5px solid #ef4444;
}

.loading-spinner {
  display: inline-block;
  width: 1em;
  height: 1em;
  border: 2px solid #fff;
  border-top: 2px solid var(--primary-dark);
  border-radius: 50%;
  animation: spin 0.7s linear infinite;
  margin-right: 8px;
  vertical-align: middle;
}
@keyframes spin {
  to { transform: rotate(360deg); }
}

@media (max-width: 768px) {
  .create-post-form {
    padding: 18px 8px 12px 8px;
  }
  
  .upload-area {
    padding: 16px;
    min-height: 100px;
  }
  
  .preview-image {
    height: 150px;
  }
}
</style>