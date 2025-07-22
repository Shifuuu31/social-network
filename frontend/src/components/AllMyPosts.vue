<template>

    <!-- <button @click="getAllMyPosts">get</button>
    <button @click="showAllMyPosts">show</button> -->
    <!-- <div v-if="this.myposts !== undefined"> -->
    <!-- {{this.posts}} -->
    <Post v-if="posts" v-for="post in this.posts" :key="post.id" v-bind:postData="post" />
    <p class="additional-info large" v-else>No posts to show</p>
    <!-- </div> -->

</template>


<script>
import Post from './Post.vue'
export default {
    name: 'AllMyPosts',
    components: { Post },
    props: {
        posts: {
        type: Array,
        default: () => []
      }
    }
}
</script>


<style>
</style>