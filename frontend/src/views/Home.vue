<template>
  <div class="layout">
    <Header />
    
    <div class="body">
      <Sidebar class="sidebar" />
      
      <main class="main">
        <div class="main-content">
          <h1 class="title">Home Feed</h1>
          
          <CreatePost @post-created="handlePostCreated" />
          
          <!-- <PostList 
            :key="postListKey" 
            :refresh-trigger="refreshTrigger"
          /> -->
          

                    <!-- <PostList 
            :key="postListKey" 
            :refresh-trigger="refreshTrigger"
          /> -->
          <PostList
          :key="postListKey"
          :refresh-trigger="refreshTrigger"
          :user-id="1"

          />


        </div>
      </main>
      
      <!-- Right sidebar for future features -->
      <!-- <aside class="right-sidebar">
        <div class="widget">
          <h3>Trending Topics</h3>
          <p>Coming soon...</p>
        </div>
      </aside> -->
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import Header from '@/components/common/Header.vue'
import Sidebar from '@/components/common/Sidebar.vue'
import PostList from '@/components/posts/PostList.vue'
import CreatePost from '@/components/posts/CreatePost.vue'

// State for refreshing post list
const postListKey = ref(0)
const refreshTrigger = ref(0)

// Handle post creation success
function handlePostCreated() {
  // Refresh the post list by changing the key
  postListKey.value += 1
  refreshTrigger.value += 1
}
</script>

<style scoped>
.layout {
  display: flex;
  flex-direction: column;
  min-height: 100vh;
  background-color: #f5f5f5;
}

.body {
  display: flex;
  flex: 1;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
}

.sidebar {
  width: 280px;
  display: none;
  background-color: white;
  border-right: 1px solid #e1e8ed;
}

.main {
  flex: 1;
  padding: 0 16px;
  min-width: 0; /* Prevents flex item from overflowing */
}

.main-content {
  max-width: 600px;
  margin: 0 auto;
}

.right-sidebar {
  width: 280px;
  padding: 16px;
  display: none;
}

.widget {
  background: white;
  border-radius: 8px;
  padding: 16px;
  box-shadow: 0 2px 4px rgba(0,0,0,0.1);
}

.widget h3 {
  margin: 0 0 12px 0;
  font-size: 18px;
  font-weight: 600;
}

.title {
  font-size: 24px;
  font-weight: bold;
  margin-bottom: 24px;
  color: #1da1f2;
  text-align: center;
}

/* Responsive Design */
@media (min-width: 768px) {
  .sidebar {
    display: block;
  }
}

@media (min-width: 1024px) {
  .right-sidebar {
    display: block;
  }
  
  .main {
    padding: 0 24px;
  }
}

@media (max-width: 768px) {
  .main {
    padding: 0 12px;
  }
  
  .title {
    font-size: 20px;
    margin-bottom: 16px;
  }
}
</style>