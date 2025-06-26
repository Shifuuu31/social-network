<template>
  <div class="posts-view">
    <div class="header">
      <h1>Home</h1>
      <button 
        @click="showCreateForm = !showCreateForm" 
        class="create-button"
        :class="{ 'active': showCreateForm }"
      >
        <svg xmlns="http://www.w3.org/2000/svg" width="24" height="24" viewBox="0 0 24 24" fill="currentColor">
          <path d="M8.8 7.2H5.6V3.9c0-.4-.3-.8-.8-.8s-.7.4-.7.8v3.3H.8c-.4 0-.8.3-.8.8s.3.8.8.8h3.3v3.3c0 .4.3.8.8.8s.8-.3.8-.8V8.7H9c.4 0 .8-.3.8-.8s-.5-.7-1-.7zm15-4.9v-.1h-.1c-.1 0-9.2 1.2-14.4 11.7-3.8 7.6-3.6 9.9-3.3 9.9.3.1 3.4-6.5 6.7-9.2 5.2-1.1 6.6-3.6 6.6-3.6s-1.5.2-2.1.2c-.8 0-1.4-.2-1.7-.3 1.3-1.2 2.4-1.5 3.5-1.7.9-.2 1.8-.4 3-1.2 2.2-1.6 1.9-5.5 1.8-5.7z"/>
        </svg>
      </button>
    </div>
    
    <PostCreate 
      v-if="showCreateForm" 
      @created="handlePostCreated"
      class="create-form"
    />
    
    <PostsList ref="postsListRef" />
  </div>
</template>

<script setup>
import { ref } from 'vue'
import PostCreate from '@/components/posts/postCreate.vue'
import PostsList from '@/components/posts/postsList.vue'

const showCreateForm = ref(false)
const postsListRef = ref(null)

const handlePostCreated = (newPost) => {
  showCreateForm.value = false
  if (postsListRef.value && postsListRef.value.fetchPosts) {
    postsListRef.value.fetchPosts()
  }
}
</script>

<style scoped>
.posts-view {
  display: flex;
  flex-direction: column;
  height: 100%;
}

.header {
  position: sticky;
  top: 0;
  z-index: 10;
  background-color: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(12px);
  padding: 1rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
  border-bottom: 1px solid var(--twitter-extra-light-gray);
}

.header h1 {
  font-size: 1.5rem;
  font-weight: 800;
  margin: 0;
}

.create-button {
  width: 2.5rem;
  height: 2.5rem;
  border-radius: 50%;
  background-color: var(--twitter-blue);
  color: white;
  border: none;
  display: flex;
  align-items: center;
  justify-content: center;
  transition: background-color 0.2s;
}

.create-button:hover {
  background-color: #1a8cd8;
}

.create-button svg {
  width: 1.25rem;
  height: 1.25rem;
}

.create-button.active {
  transform: rotate(45deg);
}

.create-form {
  border-bottom: 1px solid var(--twitter-extra-light-gray);
}
</style>