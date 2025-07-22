<template>
    <div id="navbar">
        <router-link to="/" id="nav-title">Social Network</router-link>

    </div>
</template>


<script>

export default {
    name: 'NavBarOff'
}

</script>


<style scoped>
#navbar {
    display: flex;
    padding: 10px 40px;
    background-color: var(--color-blue);
    align-items: center;
    justify-content: space-between;
    font-weight: 300;
    color: var(--color-white);
    margin-bottom: 0;

}

#navbar img {
    width: auto;
    height: 40px;
}
#nav-title {
    font-size: 24px;
    color: var(--color-white);
    text-decoration: none;
    font-weight: 400;
    position: relative;
}

#nav-title::after {
    content: "";
    height: 2.5px;
    width: 0;
    display: block;
    background-color: rgb(132, 148, 236);
    position: absolute;
    transition: width 0.35s ease-out;
}

#nav-title:hover::after {
    width: 100%;
}
</style>