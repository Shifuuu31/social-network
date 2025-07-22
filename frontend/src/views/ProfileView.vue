<template>
    <NavBarOn />
    <Profile />


</template>

<script>
import NavBarOn from "@/components/NavBarOn.vue";
import Profile from "@/components/Profile.vue"
export default {
    name: 'ProfileView',
    components: { NavBarOn, Profile, }
}
</script>

<style>
</style>