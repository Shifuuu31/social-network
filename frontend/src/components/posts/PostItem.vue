<!-- PostItem.vue -->
<template>
  
  <div class="post-item">
    <div class="post-header">
      <strong>{{ post.owner }}</strong>
      <span class="post-time">{{ formatTime(post.created_at) }}</span>
    </div>
    <div class="post-title" v-if="post.title">
      <h2>{{ post.title }}</h2>
    </div>

    <div class="post-content">
      {{ post.content }}
    </div>

    <div v-if="post.image" class="post-image">
      <img :src="post.image" alt="Post image" />
    </div>

    <div class="post-footer">
      <span>Replies: {{ post.replies }}</span>
      <span class="privacy-tag">{{ post.privacy }}</span>
    </div>
  </div>
</template>

<script setup>
import { defineProps } from 'vue'

const props = defineProps({
  post: {
    type: Object,
    required: true
  }
})
// console.log(post,"|QWWWWWWWWWWW")

function formatTime(timeStr) {
  const date = new Date(timeStr)
  return date.toLocaleString()
}
</script>

<style scoped>
.post-item {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  padding: 1rem;
  margin-bottom: 1rem;
}

.post-header {
  display: flex;
  justify-content: space-between;
  font-size: 1.1rem;
  margin-bottom: 0.5rem;
}

.post-time {
  font-size: 0.9rem;
  color: #999;
}

.post-content {
  font-size: 1rem;
  margin-bottom: 1rem;
}

.post-image img {
  max-width: 100%;
  border-radius: 8px;
}

.post-footer {
  font-size: 0.9rem;
  color: #666;
  display: flex;
  justify-content: space-between;
}

.privacy-tag {
  text-transform: capitalize;
  background: #eee;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.8rem;
}

.post-title h2 {
  font-size: 1.5rem;
  margin-top: 0;
  color: #333;
}

</style>