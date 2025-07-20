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
#app,
main{
  overflow: hidden;
}

.header {
display: block;
  position: none;
}


.chat-page {
  width: 100%;
  height: 90vh;

  background: #f8f9fa;
}

.chat-container {
  display: flex;
  flex-direction: row;
  width: 100%;
   height: 100%;
    overflow: auto;
  margin: 0 auto;
  background: white;
  box-shadow: 0 0 20px rgba(0, 0, 0, 0.1);
}
</style> 