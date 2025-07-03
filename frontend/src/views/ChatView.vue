<template>
  <div class="chat-layout">
    <ChatSidebar />
    <div class="chat-main">
      <ChatWindow />
      <ChatInput />
    </div>
  </div>
</template>

<script setup>
import ChatSidebar from '../components/ChatSidebar.vue'
import ChatWindow from '../components/ChatWindow.vue'
import ChatInput from '../components/ChatInput.vue'
</script>
