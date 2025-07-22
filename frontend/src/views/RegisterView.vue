<template>

    <div class="register__container">
        <NavBarOff />
        <Register />
    </div>

</template>

<script>
import Register from "@/components/Register.vue";
import NavBarOff from "@/components/NavBarOff.vue";
export default {
    name: 'RegisterView',
    components: { Register, NavBarOff, },
}
</script>


<style>
.register__container {
    display: grid;
    grid-template-rows: max-content 1fr;
    height: 100vh;
}
</style>