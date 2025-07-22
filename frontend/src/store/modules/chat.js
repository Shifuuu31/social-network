export default {

    state: () => ({
        newChatMessages: [],
        newGroupChatMessages: [],

        unreadMessages: [],
        unreadMsgsStatsFromDB: [],
        
        openChats: [],
        chatUserList: [],
    }),

    getters: {
        getMessages: ({ newChatMessages, newGroupChatMessages }, getters, { id }) => (receiverId, type) => {
            let messages = [];

            if (type === "PERSON") {
                messages = newChatMessages.filter((e) => {
                    return (e.receiverId === receiverId && e.senderId === id) || (e.receiverId === id && e.senderId === receiverId)
                })
            } else {
                messages = newGroupChatMessages.filter((msg) => {
                    return (msg.receiverId === receiverId)
                })
            }

            return messages
        },


        getUnreadMessagesCount: ({ unreadMessages }, getters, { id }) => (userId) => {
            const userUnreadMsgs = unreadMessages.filter((msg) => {
                return msg.senderId === userId && msg.receiverId === id
            })

            return userUnreadMsgs.length
        },

        getUnreadGroupMessagesCount: ({ unreadMessages }, getters) => (groupId) => {
            const userUnreadMsgs = unreadMessages.filter((msg) => {
                return msg.receiverId === groupId
            })

            return userUnreadMsgs.length
        },

        getUnreadMsgsCountFromDB: (state) => (userId) => {
            if (state.unreadMsgsStatsFromDB === null) {
                return 0
            }
            const userMsgObj = state.unreadMsgsStatsFromDB.find((msg) => msg.id === userId)
            if (userMsgObj === undefined) {
                return 0
            }

            return userMsgObj.unreadMsgCount
        },
    },

    mutations: {
        updateNewChatMessages(state, msgs) {
            state.newChatMessages = msgs
        },

        updateNewGroupChatMessages(state, msgs) {
            state.newGroupChatMessages = msgs
        },

        updateOpenChats(state, openChats) {
            state.openChats = openChats
        },

        updateUnreadMessages(state, unreadMsgs) {
            state.unreadMessages = unreadMsgs
        },

        updateUnreadMsgsFromDBCount(state, unreadMsgsFromDBStats) {
            state.unreadMsgsStatsFromDB = unreadMsgsFromDBStats
        },

        updateChatUserList(state, userList) {
            state.chatUserList = userList
        }

    },

    actions: {
        async fetchUnreadMessages({state}) {
            const response = await fetch('http://localhost:8081/unreadMessages', {
                credentials: 'include'
            });
            const data = await response.json();
            // console.log("/unReadmessages data", data)
            if (data.type === "Error") {
                state.unreadMsgsStatsFromDB = null;
            } else {
                state.unreadMsgsStatsFromDB = data.chatStats;

            }

        },

        addNewChatMessage({ commit, state }, payload) {
            let newMessages;

            if (payload["type"] === "PERSON") {
                newMessages = [...state.newChatMessages, payload]
                commit("updateNewChatMessages", newMessages)
            } else {
                newMessages = [...state.newGroupChatMessages, payload]
                commit("updateNewGroupChatMessages", newMessages)
            }
        },

        addUnreadChatMessage({ commit, state }, payload) {
            const unreadChatMsgs = state.unreadMessages
            unreadChatMsgs.push(payload)
            commit("updateUnreadMessages", unreadChatMsgs)
        },


        removeUnreadMessages({ state, commit }, payload) {
            let unreadMsgs;
            if (payload.type === "GROUP") {
                unreadMsgs = state.unreadMessages.filter((msg) => {
                    if (msg.receiverId === payload.receiverId) {
                        return false
                    } else {
                        return true
                    }

                })
            } else {
                unreadMsgs = state.unreadMessages.filter((msg) => {
                    if (msg.type === "PERSON" && msg.senderId === payload.receiverId) {
                        return false
                    } else {
                        return true
                    }
                })
            }

            commit('updateUnreadMessages', unreadMsgs);
        },

        addNewChat({commit, state}, chatBox) {
            let chats = state.openChats;
            chats.push(chatBox);
            commit("updateOpenChats", chats);
        },

        removeChat({commit, state}, name) {
            let newChats = state.openChats.filter((chat) => {
                return chat.name !== name
            });

            commit("updateOpenChats", newChats);
        },

        clearOpenChats({commit}) {
            commit("updateOpenChats", [])
        },


        async fetchChatUserList({rootState, commit, dispatch}) {
            await dispatch("getMyUserID");
           
            const response = await fetch('http://localhost:8081/chatList?userId=' + rootState.id, {
                credentials: 'include'
            });

            const data = await response.json();
            commit("updateChatUserList", data.users);
        }

    },

}


