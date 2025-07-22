<template>
    <div class="chatbox-wrapper">
        <div class="header">
            <p>{{ this.name }}</p>
            <i class="uil uil-times close" @click.stop="$emit('closeChat', this.name)"></i>
        </div>
        <div class="content" ref="contentDiv">

            <div class="message" v-for="(message, index) in allMessages" :style="msgPosition(message)">
                <p class="message-author" v-if="displayName(message, index)">{{ message.sender.nickname }}</p>
                <p :class="getClass(message)">{{ message.content }}</p>
            </div>

        </div>

        <form @submit.prevent="sendMessage" autocomplete="off" class="send-message" @keyup.enter="sendMessage">
            <input type="text" name="sent-message" id="sent-message__input" placeholder="Send a message"
                   ref="sendMessageInput">
            <button type="submit"><i class="uil uil-message"></i></button>
            <Emojis :input="this.$refs.sendMessageInput"
                    :messagebox="this.$refs.contentDiv"></Emojis>
        </form>

    </div>

</template>

<script>
import { mapState } from 'vuex';
import Emojis from './Emojis.vue';


export default {
    props: ["name", "receiverId", "type"],
    emits: ["closeChat"],
    data() {
        return {
            previousMessages: [],
        };
    },
    created() {
        this.getPreviousMessages();
    },
    unmounted() {
        this.clearChatNewMessages();
    },
    computed: {
        allMessages() {
            return [...this.previousMessages, ...this.$store.getters.getMessages(this.receiverId, this.type)];
        },
        ...mapState({
            myID: state => state.id
        })
    },
    watch: {
        allMessages() {
            this.$nextTick(() => {
                this.$refs.contentDiv.scrollTop = this.$refs.contentDiv.scrollHeight;
            });
        }
    },
    methods: {
        async getPreviousMessages() {
            const response = await fetch("http://localhost:8081/messages", {
                credentials: "include",
                method: "POST",
                body: JSON.stringify({
                    type: this.type,
                    receiverId: this.receiverId
                })
            });
            const data = await response.json();
            this.previousMessages = data.chatMessage ? data.chatMessage : [];
        },
        async sendMessage() {
            const sendMessageInput = this.$refs.sendMessageInput;
            if (sendMessageInput.value === "") {
                return;
            }
            const msgObj = {
                receiverId: this.receiverId,
                content: sendMessageInput.value,
                type: this.type
            };
            let response = await fetch("http://localhost:8081/newMessage", {
                body: JSON.stringify(msgObj),
                method: "POST",
                credentials: "include"
            });
            const data = await response.json();
            if (data.type == "Success"){
                this.$store.dispatch("addNewChatMessage", { ...msgObj, senderId: this.myID });
            }else{
                this.$toast.open({
                message: data.message,
                type: "warning", //One of success, info, warning, error, default
              });
            }
            sendMessageInput.value = "";
        },
        clearChatNewMessages() {
            // CLEAR NEW MESSAGES
            // because chat is closed and next time we open the chat we fetch all the messages
            if (this.type === "GROUP") {
                let msgs = this.$store.state.chat.newGroupChatMessages;
                // filter new messages by removing all messages that were sent to that receiverId
                // receiverId is equal to group ID
                msgs = msgs.filter((msg) => {
                    if (msg.receiverId === this.receiverId) {
                        return false;
                    }
                });
                this.$store.commit("updateNewGroupChatMessages", msgs);
            }
            else {
                let msgs = this.$store.state.chat.newChatMessages;
                // filter new messages by removing all messages that were sent or received in that chat
                msgs = msgs.filter((msg) => {
                    if (msg.receiverId === this.receiverId || msg.senderId === this.receiverId) {
                        return false;
                    }
                });
                this.$store.commit("updateNewChatMessages", msgs);
            }
        },
        // determines whether to display sender name in chatbox
        displayName(message, index) {
            let isSentMsg = message.senderId === this.myID;
            if (isSentMsg) {
                return false;
            }
            if (index < 1) {
                return true;
            }
            let isSequentMsg = message.senderId === this.allMessages[index - 1].senderId;
            if (isSequentMsg) {
                return false;
            }
            return true;
        },
        getClass(message) {
            let isSentMsg = message.senderId === this.myID;
            return isSentMsg ? { "sent-message": true } : { "recieved-message": true };
        },
        msgPosition(message) {
            let isSentMsg = message.senderId === this.myID;
            return {
                alignSelf: isSentMsg ? "flex-end" : "flex-start"
            };
        }
    },
    components: { Emojis }
}

</script>


<style scoped>
.chatbox-wrapper {
    height: 400px;
    width: 300px;
    display: flex;
    flex-direction: column;
    box-shadow: var(--container-shadow);
    border-radius: 5px 5px 0 0;
    overflow: hidden;
    --padding: 15px;
    --msg-border-rad: 10px;
    --msg-padding: 8px;
    --msg-margin-b: 5px;
}

.header {
    display: flex;
    justify-content: space-between;
    background-color: var(--color-blue);
    color: var(--button-content);
    padding: 12px 20px;
}



.content {
    flex: 1;
    background-color: var(--color-white);
    padding: var(--padding);
    color: var(--color-lg-black);
    font-size: 14px;
    overflow-y: auto;
    overscroll-behavior: contain;
    display: flex;
    flex-direction: column;
    gap: 10px;
}


.sent-message,
.recieved-message {
    padding: var(--msg-padding);
    border-radius: var(--msg-border-rad);
    word-break: break-word;
    display: inline-block;
}

.message {
    max-width: 80%;
}

.recieved-messages {
    align-items: flex-start;
}


.recieved-message {
    background-color: rgb(212, 212, 212);
}


.message-author {
    padding-bottom: 3.5px;
}


.sent-message {
    background-color: rgb(201, 201, 201);
    /* align-self: flex-end; */

}


/* SEND MESSAGE FORM */

.send-message {

    background-color: var(--color-white);
    padding: 10px;
    display: flex;
    flex-wrap: wrap;
    gap: 5px;
    align-items: center;

}

.send-message input {
    padding: 12px 20px;
    border-radius: 30px;
    flex: 1;

}



.send-message button {
    border: none;
    background-color: inherit;
    font-size: 1.25em;
}

.send-message button:hover {
    color: var(--hover-color);
}
</style>