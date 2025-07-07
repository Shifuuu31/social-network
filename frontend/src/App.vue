<template>
  <div id="app">
    <Header />
    <main class="main-content">
      <router-view />
    </main>
  </div>
</template>

<script setup>
import Header from './components/Header.vue'
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