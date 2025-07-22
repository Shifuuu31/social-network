<template>
    <div class="item-list__wrapper" id="groups" :style="groupsSectionStyle">
        <h3 v-if="isProfilePage">Groups</h3>
        <h3 v-else>My groups</h3>
        <ul class="item-list">
            <li v-for="group in groups" v-bind:key="group.id" v-if="groups !== null">
                <img class="small" src="../assets/icons/users-alt.svg" alt="">
                <div class="item-text">
                    <router-link :to="{ path: `/group/${group.id}`}">{{ group.name }}</router-link></div>
            </li>

            <p class="additional-info" v-else>{{noGroupsText}}</p>

        </ul>
        <NewGroup v-if="!isProfilePage"/>
    </div>

</template>


<script>
import NewGroup from '@/components/NewGroup.vue';

export default {
    props: ['groups'],
    name: 'Groups',
    components: { NewGroup },

    computed: {
        // what text to display if no groups in groups section
        noGroupsText() {
            if (this.isProfilePage) {
                return "User is not part of any group"
            } else {
                return "You are not part of any group"
                
            }
        },

        // groups section alignment is left on home page and middle on profile page
        groupsSectionStyle() {
            return {
                alignItems: this.isProfilePage ? 'center' : 'flex-start'
            }
        },

        isProfilePage() {
            return this.$route.name === "Profile"
        }
    }
}
</script>
