<template>
  <form @submit.prevent="handleSubmit" class="post-form">
    <div class="form-group">
      <textarea v-model="form.content" id="content" required placeholder="What's happening?"
        class="twitter-textarea"></textarea>
    </div>

    <div class="form-group">
      <label for="image" class="image-upload-label">
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="var(--twitter-blue)">
          <path
            d="M3 5.5C3 4.119 4.119 3 5.5 3h13C19.881 3 21 4.119 21 5.5v13c0 1.381-1.119 2.5-2.5 2.5h-13C4.119 21 3 19.881 3 18.5v-13zM5.5 5c-.276 0-.5.224-.5.5v9.086l3-3 3 3 5-5 3 3V5.5c0-.276-.224-.5-.5-.5h-13zM19 15.414l-3-3-5 5-3-3-3 3V18.5c0 .276.224.5.5.5h13c.276 0 .5-.224.5-.5v-3.086zM9.75 7C8.784 7 8 7.784 8 8.75s.784 1.75 1.75 1.75 1.75-.784 1.75-1.75S10.716 7 9.75 7z" />
        </svg>
        <input type="file" id="image" accept="image/jpeg, image/jpg, image/png, image/gif" @change="handleFileChange"
          class="hidden-input" />
      </label>
      <div v-if="selectedFile" class="image-preview">
        <img :src="previewImageUrl" alt="Preview" class="preview-image" />
        <button type="button" @click="selectedFile = null" class="remove-image-btn">
          ×
        </button>
      </div>
    </div>

    <div class="privacy-section">
      <select v-model="form.privacy" id="privacy" class="twitter-select">
        <option value="public">🌍 Public</option>
        <option value="almost_private">🔒 Almost Private</option>
        <option value="private">🔐 Private</option>
      </select>

      <div v-if="form.privacy === 'private'" class="chosen-users-input">
        <select v-model="selectedFollowers" multiple class="twitter-input" @change="updateChosenUsers">
          <option v-for="follower in followers" :key="follower.id" :value="follower.id">
            <!-- {{ follower.nickname }} -->
            <span v-if="follower.nickname"> {{ follower.nickname }}</span>
          </option>
        </select>
      </div>
      <div v-if="!isLoadingFollowers && followers.length === 0 " class="no-followers">
        You don't have any followers yet.
      </div>
    </div>

    <div class="form-footer">
      <div v-if="error" class="error-message">{{ error }}</div>
      <button type="submit" :disabled="isSubmitting" class="post-button" :class="{ 'posting': isSubmitting }">
        {{ isSubmitting ? 'Posting...' : 'Post' }}
      </button>
    </div>

    <div v-if="success" class="success-message">
      <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="var(--twitter-green)">
        <path d="M9 16.17L4.83 12l-1.42 1.41L9 19 21 7l-1.41-1.41L9 16.17z" />
      </svg>
      Post created successfully!
    </div>
  </form>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'


const followers = ref([])
const selectedFollowers = ref([])
const isLoadingFollowers = ref(false)



const ownerId = 1
const groupId = 1

const emit = defineEmits(['created', 'post-created'])

const form = ref({
  content: '',
  privacy: 'public',
})

const selectedFile = ref(null)
const chosenUsers = ref('')

const isSubmitting = ref(false)
const error = ref(null)
const success = ref(false)

const previewImageUrl = computed(() => {
  return selectedFile.value ? URL.createObjectURL(selectedFile.value) : ''
})

const handleFileChange = (e) => {
  const file = e.target.files[0]
  if (file) {
    selectedFile.value = file
  }
}


onMounted(async () => {
  await fetchFollowers()
})

// const fetchFollowers = async () => {
//   try {
//     isLoadingFollowers.value = true
//     const response = await fetch(`/api/post/${ownerId}/followers`)
//     if (!response.ok) throw new Error('Failed to fetch followers')
//     followers.value = await response.json()
//   } catch (err) {
//     error.value = err.message
//     console.error('Error fetching followers:', err)
//   } finally {
//     isLoadingFollowers.value = false
//   }
// }

// const updateChosenUsers = () => {
//   chosenUsers.value = selectedFollowers.value.join(',')
// }

const handleSubmit = async () => {
  try {
    isSubmitting.value = true
    error.value = null
    success.value = false

    const formData = new FormData()
    formData.append('owner_id', ownerId)
    formData.append('group_id', groupId)
    formData.append('content', form.value.content)
    formData.append('privacy', form.value.privacy)

    if (selectedFile.value) {
      formData.append('image', selectedFile.value)
    }

    if (form.value.privacy !== 'public') {
      const ids = chosenUsers.value
        .split(',')
        .map((id) => parseInt(id.trim()))
        .filter((id) => !isNaN(id))
      for (const id of ids) {
        formData.append('chosen_users_ids[]', id)
      }
    }

    const response = await fetch('/api/post/new', {
      method: 'POST',
      body: formData
    })

    if (!response.ok) throw new Error('Failed to create post')

    const newPost = await response.json()
    success.value = true

    // Reset
    form.value.content = ''
    form.value.privacy = 'public'
    selectedFile.value = null
    chosenUsers.value = ''

    emit('created', newPost)
    emit('post-created', newPost)

    setTimeout(() => {
      success.value = false
    }, 3000)
  } catch (err) {
    error.value = err.message
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.post-form {
  padding: 1rem;
  border-bottom: 1px solid var(--twitter-extra-light-gray);
}

.twitter-textarea {
  width: 100%;
  min-height: 100px;
  border: none;
  resize: none;
  font-size: 1.25rem;
  padding: 0.5rem;
  margin-bottom: 1rem;
  outline: none;
}

.twitter-textarea::placeholder {
  color: var(--twitter-gray);
}

.form-group {
  margin-bottom: 1rem;
}

.image-upload-label {
  display: inline-flex;
  align-items: center;
  cursor: pointer;
  color: var(--twitter-blue);
  padding: 0.5rem;
  border-radius: 50%;
  transition: background-color 0.2s;
}

.image-upload-label:hover {
  background-color: var(--twitter-light);
}

.hidden-input {
  display: none;
}

.image-preview {
  position: relative;
  margin-top: 1rem;
  max-width: 100%;
  border-radius: 1rem;
  overflow: hidden;
  border: 1px solid var(--twitter-extra-light-gray);
}

.preview-image {
  display: block;
  max-width: 100%;
  max-height: 300px;
  object-fit: contain;
}

.remove-image-btn {
  position: absolute;
  top: 0.5rem;
  right: 0.5rem;
  background: rgba(0, 0, 0, 0.6);
  color: white;
  border: none;
  width: 2rem;
  height: 2rem;
  border-radius: 50%;
  font-size: 1.2rem;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

.privacy-section {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-bottom: 1rem;
}

.twitter-select {
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  border: 1px solid var(--twitter-extra-light-gray);
  background-color: white;
  font-size: 0.9rem;
  outline: none;
}

.chosen-users-input {
  flex: 1;
}

.twitter-input[multiple] {
  height: auto;
  min-height: 38px;
  padding: 0.25rem;
}

.twitter-input option {
  padding: 0.5rem 1rem;
}

.twitter-input option:checked {
  background-color: var(--twitter-blue-light);
  color: white;
}

.loading-followers {
  color: var(--twitter-gray);
  font-size: 0.9rem;
  padding: 0.5rem;
}

.twitter-input {
  width: 100%;
  padding: 0.5rem 1rem;
  border-radius: 9999px;
  border: 1px solid var(--twitter-extra-light-gray);
  font-size: 0.9rem;
  outline: none;
}

.form-footer {
  display: flex;
  justify-content: flex-end;
  align-items: center;
  gap: 1rem;
}

.error-message {
  color: var(--twitter-red);
  font-size: 0.9rem;
  margin-right: auto;
}

.post-button {
  background-color: var(--twitter-blue);
  color: white;
  border: none;
  padding: 0.5rem 1.5rem;
  font-weight: 700;
  font-size: 1rem;
}

.post-button:hover {
  background-color: #1a8cd8;
}

.post-button:disabled {
  opacity: 0.5;
  cursor: not-allowed;
}

.success-message {
  display: flex;
  align-items: center;
  gap: 0.5rem;
  color: var(--twitter-green);
  padding: 0.5rem;
  margin-top: 1rem;
  font-weight: 500;
}

@media (max-width: 500px) {
  .privacy-section {
    flex-direction: column;
    align-items: flex-start;
  }

  .chosen-users-input {
    width: 100%;
  }
}
</style>