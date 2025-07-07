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
  background-color: #f8f9fa; /* Light soft background */
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
}

.body {
  display: flex;
  flex: 1;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  padding: 20px 0;
  gap: 24px;
}

.sidebar {
  width: 250px;
  display: none;
  background-color: white;
  border-radius: 12px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  padding: 16px;
}

.main {
  flex: 1;
  display: flex;
  justify-content: center;
}

.main-content {
  width: 100%;
  max-width: 640px;
}

.title {
  font-size: 28px;
  font-weight: 700;
  color: #1da1f2;
  text-align: center;
  margin-bottom: 24px;
  letter-spacing: -0.5px;
}

.right-sidebar {
  width: 280px;
  display: none;
  flex-shrink: 0;
}

.widget {
  background: white;
  border-radius: 12px;
  padding: 16px;
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.05);
  margin-bottom: 20px;
}

.widget h3 {
  margin: 0 0 12px 0;
  font-size: 18px;
  font-weight: 600;
  color: #333;
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
    padding: 0;
  }
}

@media (max-width: 768px) {
  .main-content {
    padding: 0 12px;
  }

  .title {
    font-size: 22px;
    margin-bottom: 16px;
  }

  .body {
    flex-direction: column;
    padding: 16px;
  }

  .sidebar,
  .right-sidebar {
    display: none;
  }
}
</style>