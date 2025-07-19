<template>
  <div class="chat-page">
    <div class="chat-container">
      <ChatList 
        :selectedUserId="selectedUserId"
        @select-conversation="selectConversation"
      />
      <ChatWindow 
        :selectedUserId="selectedUserId"
      />
    </div>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import ChatList from '../components/chat/ChatList.vue'
import ChatWindow from '../components/chat/ChatWindow.vue'

const selectedUserId = ref(null)
function selectConversation(userId) {
  selectedUserId.value = userId
}
</script>

<style scoped>
.chat-page {
  height: 100vh;
  background: #f8f9fa;
}

.chat-container {
  display: flex;
  height: 100%;
  max-width: 1200px;
  margin: 0 auto;
  background: white;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
}
</style> 