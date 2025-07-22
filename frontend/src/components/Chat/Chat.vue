<template>
    <div class="messaging-wrapper" ref="messagingWrapper">
        <ChatBox v-for="chat in openChats" v-bind="chat" @closeChat="removeChat" :key="chat.receiverId"></ChatBox>

        <div class=" messaging" @click="toggleShowContent">

            <div class="messaging-header">
                <p>Messaging
                <div class="notification-ring" v-show="hasUnreadMessages"></div>
                </p>
                <i class="uil uil-angle-up" :class="{ rotate: showContent }"></i>
            </div>

            <div class="messaging-content" v-show="showContent">
                <ul class="item-list">

                    <li v-for="user in chatUserList">
                        <div class="user">
                            <div class="user-picture small"></div>
                            <div class="item-text"
                                 @click.stop="openChat($event, { receiverId: user.id, type: 'PERSON' })">
                                {{ user.nickname }}</div>
                        </div>

                        <p class="unreadMessagesCount"
                           v-if="totalUnreadMessagesCount(user.id, 'PERSON') !== 0">
                            {{ totalUnreadMessagesCount(user.id, 'PERSON') }}</p>

                    </li>
                </ul>

                <ul class="item-list" v-if="userGroups !== null">
                    <li v-for="group in userGroups">

                        <div class="group">
                            <img src="../../assets/icons/users-alt.svg" alt="" class="small">

                            <div class="item-text"
                                 @click.stop="openChat($event, { receiverId: group.id, type: 'GROUP' })">
                                {{ group.name }}</div>
                        </div>

                        <p class="unreadMessagesCount"
                           v-if="totalUnreadMessagesCount(group.id, 'GROUP') !== 0">
                            {{ totalUnreadMessagesCount(group.id, 'GROUP') }}</p>

                    </li>

                </ul>

                <p class="additional-info" v-if="chatUserList === null && userGroups === null">
                    No one to message with :(</p>
            </div>

        </div>
    </div>

</template>


<script>
import { mapGetters, mapState } from 'vuex';
import ChatBox from './ChatBox.vue'
export default {
    name: '',
    components: { ChatBox },
    data() {
        return {
            showContent: false,
        }
    },

    unmounted() {
        this.$store.commit("updateUnreadMessages", [])
        this.$store.dispatch("clearOpenChats")
    },

    created() {
        this.$store.dispatch("fetchChatUserList");
        this.$store.dispatch("getUserGroups");
        this.$store.dispatch("fetchUnreadMessages");
    },

    computed: {
        ...mapState({
            userGroups: state => state.groups.userGroups,
            openChats: state => state.chat.openChats,
            chatUserList: state => state.chat.chatUserList,
            unreadMsgsStatsFromDB: state => state.chat. unreadMsgsStatsFromDB,
           
        }),

        ...mapGetters(['getUnreadMessagesCount', 'getUnreadGroupMessagesCount', 'getUnreadMsgsCountFromDB']),

        hasUnreadMessages() {
            if (this.$store.state.chat.unreadMessages.length > 0) {
                return true
            }

            if (this.unreadMsgsStatsFromDB !== null && this.unreadMsgsStatsFromDB.length > 0) {
                return true
            }

            return false
        },
    },


    methods: {
        async getUsersIFollow() {
            await this.$store.dispatch("getMyUserID");

            const response = await fetch('http://localhost:8081/chatList?userId=' + this.$store.state.id, {
                credentials: 'include'
            });

            const data = await response.json();

            this.usersIFollow = data;

        },

    
        openChat(e, obj) {
            // console.log("Trying to add a chatbox")
            // console.log(e.target.textContent)
            const found = this.openChats.some(chat => chat.name === e.target.textContent);
            if (found) {
                return
            };

            if (this.$refs.messagingWrapper.clientWidth + 300 > window.innerWidth) {
                this.openChats.shift();
            }

            this.$store.dispatch("addNewChat", {
                "name": e.target.textContent,
                ...obj
            })

            this.$store.dispatch("removeUnreadMessages", { receiverId: obj.receiverId, type: obj.type })

            if (Array.isArray(this.unreadMsgsStatsFromDB)) {
                let unreadMsgsStatsFromDB = this.unreadMsgsStatsFromDB.filter((msgStats) => msgStats.id !== obj.receiverId);
                this.$store.commit("updateUnreadMsgsFromDBCount", unreadMsgsStatsFromDB)
            }

        },


        removeChat(name) {
            this.$store.dispatch("removeChat", name)
        },


        // adds to the previous unread messages
        // unread messages from database and current session
        totalUnreadMessagesCount(receiverId, type) {
            if (type === "PERSON") {
                return this.getUnreadMessagesCount(receiverId) + this.getUnreadMsgsCountFromDB(receiverId)
            } else {
                return this.getUnreadGroupMessagesCount(receiverId) + this.getUnreadMsgsCountFromDB(receiverId)
            }
        },

        toggleShowContent() {
            // console.log("Content toggled!")
            this.showContent = !this.showContent
        },


    },


    watch: {
        $route() {
            // console.log("Route changed!")
            // keep chat list updated while user navigates around the website
            this.$store.dispatch("fetchChatUserList");
        }
    }
}
</script>


<style scoped>
.messaging-wrapper {
    position: fixed;
    bottom: 0;
    right: 0;
    display: flex;
    align-items: flex-end;
}

.messaging {

    width: 300px;
}

.messaging-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 12px 20px;
    background-color: var(--color-blue);
    border-top-left-radius: 5px;
    color: var(--color-white);
    cursor: pointer;
}

.messaging-header i {
    font-size: 20px;
    transition: transform 0.35s linear
}


.messaging-header p {
    font-weight: 300;
    letter-spacing: 0.35px;
    position: relative;
}


.notification-ring {
    position: absolute;
    height: 10px;
    width: 10px;
    background-color: rgb(207, 59, 59);
    border-radius: 50%;
    right: -15px;
    top: 0;
}

.messaging-content {
    background-color: var(--color-white);
    padding: 20px;
}

.messaging .item-list:first-of-type {
    padding-bottom: 20px;
}

.rotate {
    transform: rotate(180deg);
}

.user,
.group {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 5px;
}


.item-list li {
    justify-content: space-between;
}


.messaging-header:hover {
    background-color: var(--hover-background-color);
}


.unreadMessagesCount {
    padding: 2.5px 5px;
    color: var(--color-white);
    background-color: brown;
    font-size: 12px;
    border-radius: 5px;
}
</style>