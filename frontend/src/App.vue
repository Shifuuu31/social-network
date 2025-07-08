<template>
  <div id="app">
    <Header v-if="!['/signin', '/signup'].includes(route.path)" />
    <main>
      <RouterView />
    </main>
  </div>
</template>

<script setup>
import { useRoute } from 'vue-router'
import Header from '@/components/common/Header.vue'

const route = useRoute()
</script>