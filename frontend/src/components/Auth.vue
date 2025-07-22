<template>

</template>


<script>

export default {
    name: 'Auth',
    created() {
        this.isLoggedIn()
    },
    methods: {
        async isLoggedIn() {
            await fetch('http://localhost:8081/sessionActive', {
                credentials: 'include'
            })
                .then((response => {
                    // console.log(response)
                    return response.json()
                }))
                // .then((json=>console.log(json)))
                .then((json => {
                    if (json.message === "Session active") {
                        // console.log('yes');
                        this.$router.push("/main")
                    } else {
                        // console.log('no');
                        this.$router.push("/sign-in")
                    }
                }))
        },

    }
}
</script>


<style>
</style>