<template>
  <div id="app">
    <Header />
    <main class="main-content">
      <router-view />
    </main>
    <!-- Global notification component -->
    <Notifications />
  </div>
</template>

<script setup>
import Header from './components/Header.vue'
import Notifications from './components/notifications.vue'
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