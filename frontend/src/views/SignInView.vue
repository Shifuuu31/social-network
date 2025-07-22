<template>
    <div class="sign-in__container">
        <NavBarOff />
        <SignIn />

    </div>
</template>

<script>

// @ is an alias to /src
import NavBarOff from '@/components/NavBarOff.vue';
import SignIn from '@/components/SignIn.vue';

export default {
    name: 'SignInView',
    components: {
        NavBarOff,
        SignIn,
    }
}

</script>

<style>
.sign-in__container {
    display: grid;
    grid-template-rows: max-content minmax(min-content, 1fr);
    height: 100vh;

}

.sign-in__container .sign-in__wrapper {
    place-self: center;
}
</style>
