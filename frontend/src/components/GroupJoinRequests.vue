<template>
    <div v-if="isAdmin && this.joinRequests.length > 0">
        <div class="item-list__wrapper" id="requests">
            <h3>Requests</h3>
            <ul class="item-list">
                <div v-for="notif in this.joinRequests">
                    <li>
                        <div class="row1">
                            <img class="small" :src="'http://localhost:8081/' + notif.user.avatar">
                            <div>
                                {{ notif.user.nickname }}
                            </div>
                        </div>
                        <div class="row2">
                            <i class="uil uil-times decline" @click="reactToRequest(notif.id,'decline')"></i>
                            <i class="uil uil-check accept" @click="reactToRequest(notif.id,'accept')"></i>
                        </div>
                    </li>

                </div>
            </ul>
        </div>
    </div>
</template>

<script>

export default {
    name: "GroupJoinRequests",
    props: {
        isAdmin: false
    },
    data() {
        return {
            joinRequests: []
        }
    },
    created() {
        if (this.isAdmin) {
            this.getJoinRequests()
        }
    },
    methods: {
        async getJoinRequests() {
            await fetch("http://localhost:8081/groupRequests?groupId=" + this.$route.params.id, {
                credentials: "include"
            }).then(r => r.json()).then(json => {
                // console.log(json)
                this.joinRequests = json.notifications
            })
        },
        async reactToRequest(requestId, response){
            await fetch("http://localhost:8081/responseGroupRequest", {
                method:"POST",
                credentials: "include",
                body: JSON.stringify({
                    groupId: this.$route.params.id,
                    requestId: requestId,
                    response: response
                })
            }).then(r => r.json()).then(json => {
                console.log(json)
                if (json.type === "Success"){
                    this.getJoinRequests()
                    this.$toast.open({
                            message: "Done!",
                        });
                    // delete notification
                    this.$store.dispatch("removeNotification", requestId);
                }
            })
        }
    },

}
</script>

<style scoped>

#requests .item-list {
    gap: 10px;
}

#requests .item-list li {
    justify-content: space-between;
    gap: 20px;
}
</style>