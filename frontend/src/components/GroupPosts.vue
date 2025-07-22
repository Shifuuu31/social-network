<template>
    <Post v-for="post in this.groupPosts" :key="post.id" v-bind:postData="post" />
</template>


<script>
import Post from './Post.vue'
import { mapGetters } from 'vuex'
export default {
    name: 'GroupPosts',
    components: { Post },
    created() {
        this.$store.dispatch('getGroupPosts')
    },
    computed: mapGetters(['groupPosts']),
    watch: {
        $route() {
            this.$store.dispatch('getGroupPosts')
        }
    },
    methods: {}
}
</script>


<style>
</style>