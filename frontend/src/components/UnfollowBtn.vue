<template>
    <button class="btn active" @click="unfollow">Unfollow<i class="uil uil-user-plus"></i></button>
</template>

<script>
export default {
    // props: ['profileId'],
    // data() {
    //     return {
    //         userid: "",
    //         isMyProfile: false,
    //     }
    // },

    // created() {
    //     this.checkProfile();
    //     this.$store.dispatch("getMyUserID").then(() => console.log(this.$store.state.id))
    // },
    emits: ["unfollow"],

    methods: {
        async unfollow() {
            // console.log('subscribe function:')
            await fetch("http://localhost:8081/unfollow?userId=" + this.$route.params.id, {
                credentials: "include",
            })
                .then((r) => r.json())
                .then(json => {
                    this.$emit("unfollow")

                })
        },

        // async getLoggedUserId() {
        //     await fetch("http://localhost:8081/currentUser", {
        //         credentials: "include",
        //     })
        //         .then((r) => r.json())
        //         .then((json => {
        //             this.userid = json.users[0].id
        //         }))
        // },

        // async checkProfile() {
        //     await this.getLoggedUserId();
        //     this.isMyProfile = this.userid === this.profileId
        // },
    }
}
</script>