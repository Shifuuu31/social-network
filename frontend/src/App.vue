<template>
  <div id="app">
    <Header />
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import { onMounted, onUnmounted } from 'vue'
import Header from './components/Header.vue'
import { useNotificationStore } from './stores/notificationStore'

const notificationStore = useNotificationStore()

onMounted(async () => {
  // Initialize notification system when app starts
  notificationStore.connectWebSocket()
  notificationStore.requestNotificationPermission()
  notificationStore.fetchUnreadCount()
})

onUnmounted(() => {
  // Clean up WebSocket connection when app is destroyed
  notificationStore.disconnect()
})
</script>

<style scoped>
#app {
  min-height: 100vh;
  background: #0a0a0a;
  color: #fff;
}

.main-content {
  padding-top: 80px;
}
</style>