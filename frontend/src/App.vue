<template>

    <router-view>

    </router-view>

    <router-view name="Chat"></router-view>
</template>

<script>

import Chat from './components/Chat/Chat.vue'
export default {
    name: 'App',
    components: { Chat },

    mounted() {
        window.addEventListener("load", this.createWebSocketConn)
    },

    methods: {
        createWebSocketConn() {
            if (this.$route.path === "/sign-in" || this.$route.path === "/reg" ) {
                return
            }
            this.$store.dispatch("createWebSocketConn")
        }
    }
}
</script>

<style>
</style>