<!-- src/components/comments/CommentForm.vue -->
<template>
  <div class="comment-form-container">
    <div class="comment-form-header">
      <div class="comment-avatar">
        <svg xmlns="http://www.w3.org/2000/svg" width="48" height="48" viewBox="0 0 24 24" fill="var(--twitter-blue)">
          <path d="M12 2C6.48 2 2 6.48 2 12s4.48 10 10 10 10-4.48 10-10S17.52 2 12 2zm0 3c1.66 0 3 1.34 3 3s-1.34 3-3 3-3-1.34-3-3 1.34-3 3-3zm0 14.2c-2.5 0-4.71-1.28-6-3.22.03-1.99 4-3.08 6-3.08 1.99 0 5.97 1.09 6 3.08-1.29 1.94-3.5 3.22-6 3.22z"/>
        </svg>
      </div>
      
      <form @submit.prevent="handleSubmit" class="comment-form">
        <textarea
          ref="textareaRef"
          v-model="commentText"
          placeholder="Post your reply"
          class="comment-input"
          rows="3"
          maxlength="280"
          @input="adjustTextareaHeight"
        ></textarea>
        
        <div class="comment-form-footer">
          <div class="character-count" :class="{ 'warning': isNearLimit, 'danger': isOverLimit }">
            {{ remainingChars }}
          </div>
          
          <button 
            type="submit" 
            class="reply-btn"
            :disabled="!canSubmit || isSubmitting"
          >
            {{ isSubmitting ? 'Replying...' : 'Reply' }}
          </button>
        </div>
        
        <div v-if="error" class="error-message">{{ error }}</div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, nextTick } from 'vue'

const props = defineProps({
  postId: {
    type: [String, Number],
    required: true
  }
})

const emit = defineEmits(['comment-added'])

const textareaRef = ref(null)
const commentText = ref('')
const isSubmitting = ref(false)
const error = ref(null)

const remainingChars = computed(() => 280 - commentText.value.length)
const isNearLimit = computed(() => remainingChars.value <= 20 && remainingChars.value > 0)
const isOverLimit = computed(() => remainingChars.value < 0)
const canSubmit = computed(() => 
  commentText.value.trim().length > 0 && 
  commentText.value.length <= 280 && 
  !isSubmitting.value
)

const adjustTextareaHeight = async () => {
  await nextTick()
  if (textareaRef.value) {
    textareaRef.value.style.height = 'auto'
    textareaRef.value.style.height = textareaRef.value.scrollHeight + 'px'
  }
}

const handleSubmit = async () => {
  if (!canSubmit.value) return
  
  try {
    isSubmitting.value = true
    error.value = null
    console.log(`Posting comment for post ID: ${props.postId} with content: "${commentText.value}"`);
    
    const response = await fetch(`/api/comments/${props.postId}/new`, {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json'
      },
      body: JSON.stringify({
        owner_id: 1, // TODO - replace with actual user ID
        post_Id: props.postId,
        content: commentText.value.trim(),
      })
    })
    
    if (!response.ok) {
      const errorData = await response.text()
      throw new Error(`Failed to post comment: ${errorData}`)
    }
    
    const responseData = await response.json()
    console.log('Comment creation response:', responseData)
    
    const newComment = {
      id: responseData.id || Date.now(),
      content: commentText.value.trim(),
      author: responseData.author || responseData.owner_name || 'You',
      createdAt: responseData.createdAt || responseData.created_at || new Date().toISOString(),
      likes: 0,
      ...responseData
    }
    
    emit('comment-added', newComment)
    
    commentText.value = ''
    adjustTextareaHeight()
    
  } catch (err) {
    console.error('Comment submission error:', err)
    error.value = err.message || 'Failed to post comment'
  } finally {
    isSubmitting.value = false
  }
}

const focus = () => {
  if (textareaRef.value) {
    textareaRef.value.focus()
  }
}

const replyTo = (username) => {
  commentText.value = `@${username.replace(/\s+/g, '')} `
  focus()
  adjustTextareaHeight()
}

defineExpose({
  focus,
  replyTo
})
</script>

<style scoped>
.comment-form-container {
  border-bottom: 1px solid var(--twitter-extra-light-gray);
  background: white;
}

.comment-form-header {
  display: flex;
  padding: 1rem;
  gap: 0.75rem;
}

.comment-avatar svg {
  width: 48px;
  height: 48px;
}

.comment-form {
  flex: 1;
  display: flex;
  flex-direction: column;
}

.comment-input {
  border: none;
  outline: none;
  font-size: 1.25rem;
  line-height: 1.5;
  resize: none;
  min-height: 50px;
  max-height: 200px;
  font-family: inherit;
  background: transparent;
  color: var(--twitter-dark);
}

.comment-input::placeholder {
  color: var(--twitter-gray);
}

.comment-form-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 0.75rem;
  padding-top: 0.75rem;
  border-top: 1px solid var(--twitter-extra-light-gray);
}

.character-count {
  font-size: 0.875rem;
  color: var(--twitter-gray);
  font-weight: 500;
}

.character-count.warning {
  color: #ff6600;
}

.character-count.danger {
  color: #ff0000;
}

.reply-btn {
  background: var(--twitter-blue);
  color: white;
  border: none;
  border-radius: 9999px;
  padding: 0.5rem 1rem;
  font-weight: 700;
  font-size: 0.9rem;
  cursor: pointer;
  transition: background-color 0.2s;
  min-width: 70px;
}

.reply-btn:hover:not(:disabled) {
  background: #1991db;
}

.reply-btn:disabled {
  background: var(--twitter-light-gray);
  cursor: not-allowed;
}

.error-message {
  color: #ff0000;
  font-size: 0.875rem;
  margin-top: 0.5rem;
  padding: 0.5rem;
  background: #fee;
  border-radius: 4px;
  border: 1px solid #fcc;
}

@media (max-width: 500px) {
  .comment-form-header {
    padding: 0.75rem;
  }
  
  .comment-avatar svg {
    width: 40px;
    height: 40px;
  }
  
  .comment-input {
    font-size: 1.1rem;
  }
}
</style>