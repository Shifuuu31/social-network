<template>

    <div class="item-list__wrapper" id="followers">
        <h3>Followers</h3>
        <ul class="item-list users" v-if="this.followers">
            <li v-for="user in this.followers" :key="user.id">
                <div class="user-picture small"></div>
                <div class="item-text"><router-link :to="{ path: `/profile/${user.id}`}">{{ user.nickname }}
</router-link></div>
            </li>
        </ul>

        <p class="additional-info" v-else>No followers</p>
    </div>

</template>


<script>
export default {
    name: 'Followers',
    props: {       
        followers: {
        type: Array,
        default: () => []
      }
    }
}
</script>


<style>
</style>