<template>
  <div class="post-item">
    <!-- Post Header -->
    <div class="post-header">
      <div class="post-author">
        <img :src="post.authorAvatar" alt="avatar" class="avatar" />
        <div class="author-info">
          <h3 class="author-name">{{ post.author }}</h3>
          <p class="post-date">{{ post.date }}</p>
        </div>
      </div>
      <button class="more-btn">•••</button>
    </div>

    <!-- Post Content -->
    <p class="post-content">{{ post.content }}</p>
    <img v-if="post.image" :src="post.image" alt="Post Image" class="post-image" />

    <!-- Post Actions -->
    <div class="post-actions">
      <button class="action-btn">👍 Like</button>
      <button class="action-btn">💬 Comment</button>
      <button class="action-btn">🔁 Share</button>
    </div>
  </div>
</template>

<script setup>
defineProps({
  post: {
    type: Object,
    required: true,
    default: () => ({
      author: 'John Doe',
      authorAvatar: 'https://via.placeholder.com/40',
      date: 'June 29, 2025',
      content: 'This is a sample post. Later this will come from backend!',
      image: 'https://via.placeholder.com/600x300',
    }),
  },
})
</script>

<style scoped>
.post-item {
  background-color: white;
  border: 1px solid #ddd;
  border-radius: 8px;
  padding: 16px;
  margin-bottom: 20px;
}

.post-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
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
  margin-right: 10px;
  border: 1px solid #ccc;
}

.author-info {
  display: flex;
  flex-direction: column;
}

.author-name {
  margin: 0;
  font-size: 16px;
  font-weight: bold;
}

.post-date {
  margin: 0;
  font-size: 12px;
  color: gray;
}

.more-btn {
  background: none;
  border: none;
  font-size: 20px;
  color: gray;
  cursor: pointer;
}

.post-content {
  margin: 12px 0;
  font-size: 14px;
}

.post-image {
  width: 100%;
  max-height: 300px;
  object-fit: cover;
  border-radius: 6px;
  margin-bottom: 12px;
}

.post-actions {
  display: flex;
  justify-content: space-around;
  border-top: 1px solid #eee;
  padding-top: 10px;
}

.action-btn {
  background: none;
  border: none;
  color: #555;
  cursor: pointer;
  font-size: 14px;
}

.action-btn:hover {
  color: #000;
}
</style>
