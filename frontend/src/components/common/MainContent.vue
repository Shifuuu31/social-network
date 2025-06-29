<template>
  <div>
    <Header />
    <CreatePost />
    <PostList />
  </div>
</template>

<script setup>
import Header from './Header.vue'
import CreatePost from '@/components/posts/CreatePost.vue'
import PostList from '@/components/posts/PostList.vue'
</script>
