 <template>
    <div id="all_posts">
        <Post
              v-if="allPosts !== undefined"
              v-for="post in allPosts" :key="post.id" v-bind:postData="post" />

        <p class="additional-info large" v-else>No posts</p>
    </div>

</template>

<script>
import Post from './Post.vue'
import { mapGetters } from 'vuex'
export default {
    name: 'AllPosts',
    created() {
        this.$store.dispatch('fetchPosts')
    },


    computed: mapGetters(['allPosts']),
    components: { Post },
}
</script>

<style>
#all_posts {
    display: flex;
    flex-direction: column;
    justify-content: center;
    align-items: center;
    gap: 50px;
}
</style>