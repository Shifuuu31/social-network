<template>

    <div id="navbar">

        <div id="nav-titleSearch">
            <router-link to="/" id="nav-title"><img src="../assets/logo.png"></router-link>
            <Search />
        </div>
        <ul class="nav-links">

            <li id="notifications-link">
                <Notifications />
            </li>
            <li>
                <router-link v-if="typeof user.id !== 'undefined'"
                             :to="{ name: 'Profile', params: { id: user.id } }">My profile</router-link>
            </li>
            <li @click="logout">Log out</li>
        </ul>

    </div>

</template>


<script>
import Search from './Search.vue';
import Notifications from './Notifications.vue'
export default {
    name: 'NavBarOn',
    data() {
        return {
            user: {}
        }
    },
    created() {
        this.getUserInfo()
    },
    methods: {
        async getUserInfo() {
            await fetch("http://localhost:8081/currentUser", {
                credentials: 'include',
            })
                .then((r => r.json()))
                .then((json => {
                    // console.log(json)
                    this.user = json.users[0]
                }))

        },
        async logout() {
            await fetch('http://localhost:8081/logout', {
                credentials: 'include',
                headers: {
                    'Accept': 'application/json',
                }
            })
                .then((response => response.json()))
                .then((json => { console.log(json) }))
            // console.log("logout")
            this.$store.state.wsConn.close(1000, "user logged out");
            this.$router.push("/");
        }
    },
    components: { Notifications, Search, }
}

</script>

<style scoped>
#navbar {
    position: fixed;
    top: 0;
    left: 0;
    z-index: 1;
    width: 100%;
    min-width: min-content;

    display: flex;
    align-items: center;
    justify-content: space-between;

    padding: 10px 40px;
    background-color: var(--color-blue);
    color: var(--color-white);

    position: relative;




}


#navbar a {
    color: var(--color-white);
}


#navbar img {
    width: auto;
    height: 40px;
}


#nav-title {
    font-size: 24px;
    font-weight: 400;
    position: relative;
}


.nav-links li {
    font-weight: 300;
    display: inline-block;
    margin-left: 20px;
    cursor: pointer;

    position: relative;
}

 

#nav-titleSearch {
    display: flex;
    gap: 25px;
    flex-grow: 1;
    align-items: center;


}


#navbar li:not(#notifications-link)::after,
#nav-title::after {
    content: "";
    height: 2.5px;
    width: 0;
    display: block;
    position: absolute;

    transition: all 0.35s ease-out;
}

#navbar li:not(#notifications-link):hover::after,
#nav-title:hover::after {
    width: 100%;
    background-color: rgb(132, 148, 236);
}

a:link {
    text-decoration: none;
}

a:visited {
    text-decoration: none;
}

a:hover {
    text-decoration: none;
}

a:active {
    text-decoration: none;
}
</style>