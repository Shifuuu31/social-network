<template>
    <div class="item-list__wrapper" id="following">
        <h3>Following</h3>
        <ul class="item-list users" v-if="this.following">
            <li v-for="user in this.following" :key="user.id">
                <div class="user-picture small"></div>
                <div class="item-text"><router-link :to="{ path: `/profile/${user.id}`}">{{ user.nickname }}
</router-link></div>
            </li>
            <!--  -->
        </ul>

        <p class="additional-info" v-else>Not following anyone</p>
    </div>
    <!-- 
    <div>
        <span><b>Following</b></span>
        <div v-for="user in this.following" :key="user.id">           
                {{user.nickname}}
        </div>
    </div>


    <li>
        <div class="user-picture small"></div>
        <div class="item-text">User 2</div>
    </li>
    <li>
        <div class="user-picture small"></div>
        <div class="item-text">User 3</div>
    </li>
    <li>
        <div class="user-picture small"></div>
        <div class="item-text">User 4</div>
    </li> -->
</template>


<script>
export default {
    name: 'Following',
        props: {       
        following: {
        type: Array,
        default: () => []
      }
    }
}
</script>


<style>
</style>