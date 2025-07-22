<template>

    <div class="relative-wrapper" @click="toggleShowNotifications">
        <span class="link-header">Notifications
            <div class="notification-ring" v-show="hasNotifications"></div>

        </span>
        <div class="item-list__wrapper" id="notifications" v-show="showNotifications">
            <ul class="item-list">
                <li v-for="notification in allNotifications" :key="notification.id"
                    v-if="hasNotifications">

                    <div class="row1 ">
                        <img class="" src="../assets/icons/default-profile.svg">
                        <NotificationMsg :notification="notification"></NotificationMsg>
                    </div>

                    <div class="row2" v-if="notification.type === 'EVENT'">
                         <i class="uil uil-times decline" @click.stop="handleEventRequest(notification, 'NO')"></i>
                        <i class="uil uil-check accept" @click.stop="handleEventRequest(notification, 'YES')"></i>
                    </div>

                    <div class="row2" v-else>
                        <i class="uil uil-times decline" @click.stop="handleRequest(notification, 'decline')"></i>
                        <i class="uil uil-check accept" @click.stop="handleRequest(notification, 'accept')"></i>
                    </div>
                </li>

                <li v-else class="additional-info">No notifications</li>

            </ul>
        </div>
    </div>

</template>

<script>
import { mapState } from 'vuex';
import NotificationMsg from './NotificationMsg.vue';

export default {
    name: "notifications",
    data() {
        return {
            showNotifications: false,
            notificationsFromDB: {},
        };
    },
    async created() {
        await this.fetchNotifications();
    },
    computed: {
        ...mapState({
            allNotifications: state => state.notifications.allNotifications
        }),
        hasNotifications() {
            if (this.allNotifications === null) {
                return false;
            } else {
                return this.allNotifications.length > 0;
            }
        }
    },
    unmounted() {
        this.$store.commit("updateAllNotifications", []);
    },
    methods: {
        async handleEventRequest(notification, reqResponse){
            const response = await fetch(`http://localhost:8081/participate`, {
                credentials: "include",
                method: "POST",
                body: JSON.stringify({
                    requestId: notification.id,
                    eventId: notification.event.id,
                    response: reqResponse,
                })
            });
            const data = await response.json();
            if (data.type == "Success"){
                // remove the notification
                this.$store.dispatch("removeNotification", notification.id);
            }

            if (!this.hasNotifications) {
                this.toggleShowNotifications();
            }
        },
        toggleShowNotifications() {
            this.showNotifications = !this.showNotifications;
        },

        async fetchNotifications() {
            const response = await fetch("http://localhost:8081/notifications", {
                credentials: "include"
            });
            const data = await response.json();
            if (data.type === "Error") {
                this.notificationsFromDB = null;
                this.$store.commit("updateAllNotifications", null);
            } else {
                this.notificationsFromDB = data;
                this.$store.commit("updateAllNotifications", data.notifications);
            }
           
            // console.log("/notifications data", data)
        },
        async handleRequest(notification, reqResponse) {
            let endpoint;
            switch (notification.type) {
                case "FOLLOW":
                    endpoint = "responseFollowRequest";
                    break;
                case "GROUP_INVITE":
                    endpoint = "responseInviteRequest";
                    break;
                case "GROUP_REQUEST":
                    endpoint = "responseGroupRequest"
                    break;
                case "CHAT_REQUEST":
                    endpoint = "responseChatRequest"
                    break;
            }
            const response = await fetch(`http://localhost:8081/${endpoint}`, {
                credentials: "include",
                method: "POST",
                body: JSON.stringify({
                    requestId: notification.id,
                    response: reqResponse,
                    groupId: notification.group.id
                })
            });
            const data = await response.json();
            // console.log("/handleRequest data", data)
            if (notification.type === "GROUP_INVITE" && reqResponse === "accept") {
                // update user groups for live update
                this.$store.dispatch("addUserGroup", notification.group);
            }

            if (notification.type === "CHAT_REQUEST" || notification.type == "FOLLOW" && reqResponse === "accept") {
                // update user groups for live update
                // console.log("NOTIFICATION", notification)
                this.$store.dispatch("fetchChatUserList");
                this.$store.dispatch("fetchUnreadMessages");
                
            }

            // if (notification.type === "GROUP_INVITE" && reqResponse === "accept")
            // remove the notification
            this.$store.dispatch("removeNotification", notification.id);

            // this.checkNotificationsLength();

            if (!this.hasNotifications) {
                this.toggleShowNotifications();
            }
        },
     
        isDataValid(resp) {
            return resp.type === "Success" ? true : false;
        },
        additionalText(notification) {
            let a = "";
            switch (notification.type) {
                case "EVENT":
                    // console.log("Notif", notification);
                    a = `${notification.event.title}`;
                    break;
                case "GROUP_INVITE":
                    a = notification.group.name;
                    break;
            }
            // event need group name, event name
            // group invite -> who invited and to what group
            return a;
        }
    },
    components: { NotificationMsg }
}

</script>

<style>
.relative-wrapper {
    position: relative;

}


#notifications .row1 :is(img, i) {
    height: 2.25em;
    width: 2.25em;
}

#notifications {

    position: absolute;
    transform: translateX(-50%);
    left: 50%;
    font-weight: 400;
    margin-top: 10px;
    width: 400px;
    cursor: default;
}



#notifications .item-list {
    gap: 15px;
}

#notifications .item-list li {
    justify-content: space-between;
    gap: 20px;
}


.who {
    font-weight: 500;
}


.link-header::after {
    content: "";
    height: 2.5px;
    width: 0;
    display: block;

    position: absolute;

    transition: all 0.35s ease-in-out;
}

.link-header:hover::after {
    width: 100%;
    background-color: rgb(132, 148, 236);
}

.notification-ring {
    position: absolute;
    height: 10px;
    width: 10px;
    background-color: rgb(207, 59, 59);
    border-radius: 50%;
    right: -12.5px;
    top: 0;
}

.decline,
.accept {
    display: inline-block;
    transition: transform 0.25s linear;
}

.decline:hover,
.accept:hover {
    display: inline-block;
    transform: scale3d(1.125, 1.125, 1.125);
}
</style>