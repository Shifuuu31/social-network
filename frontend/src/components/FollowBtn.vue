<template>
    <button :class="buttonClass" @click="handleFollow">{{ buttonText }}
        <i class="uil uil-user-plus" v-if="buttonText === 'Follow'"></i>
        <i class="uil uil-check" v-else></i>
    </button>
</template>


<script>
export default {
    /* 
    1. Check if user has sent a request to follow that user
    2. Change button text based on not following or request sent

    */


    name: 'FollowBtn',
    props: ['user'],
    emits: ["follow"],
    data() {
        return {

            // tracks button functionality
            // functionalities -> follow, send request, decline request
            buttonText: "Follow",

        }
    },

    created() {
        if (this.user.requestPending === true) {
            this.buttonText = "Request sent";
        }

        // console.log("User", this.user)
    },

    computed: {
        buttonClass() {
            return {
                btn: true,
                // active hover class
                active: this.buttonText === "Request sent"
            }
        }
    },

    methods: {
        async handleFollow() {
            if (this.buttonText === "Request sent") {
                // cancel the request
                await this.cancelFollowRequest();
                this.buttonText = "Follow"
                return
            }

            await this.followUser();

        },

        async followUser() {
            // console.log('subscribe function:')
            await fetch("http://localhost:8081/follow?userId=" + this.$route.params.id, {
                credentials: "include",
            })
                .then((r) => r.json())
                .then((json => {

                    // console.log("server response:", json)

                    if (this.user.status === "PRIVATE") {
                        this.buttonText = "Request sent"
                        this.$emit("follow", "requestSent")
                    } else {
                        this.$emit("follow", "followedUser")
                    }

                }))
        },


        async cancelFollowRequest() {
            const response = await fetch(`http://localhost:8081/cancelFollowRequest?userid=${this.$route.params.id}`, {
                credentials: 'include'
            });
            const data = await response.json();

            // console.log("/cancelFollowRequest response", data)
        }

    }
}
</script>


<style>
#followBtn {
    height: 37px;
    width: 92px;
    left: 0px;
    top: 0px;
    border-radius: 8px;
    padding: 10px, 15px, 10px, 15px;
}
</style>