<template>
  <div id="app">
    <!-- Uncomment these when ready -->
    <!-- <AppHeader /> -->
    <main>
      <RouterView />
    </main>
    <!-- <AppFooter /> -->
    <!-- <ToastContainer /> -->
  </div>
</template>

<script setup>
import { RouterView } from 'vue-router'
// import PostTest from '@/components/posts/PostTest.vue'
// import AppHeader from '@/components/AppHeader.vue'
// import AppFooter from '@/components/AppFooter.vue'
// import ToastContainer from '@/components/ui/ToastContainer.vue'
</script>







<!-- //////////////////////////////////////////////////////// -->


<!-- <template>
  <div>
    <Header v-if="isLoggedIn" />
    <router-view />
  </div>
</template>

<script setup>
import Header from '@/components/common/Header.vue'

const isLoggedIn = !!localStorage.getItem('token')
</script> -->