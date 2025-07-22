import router from "@/router"

export default {
    async fetchPosts() {
        await fetch("http://localhost:8081/allPosts", {
            credentials: "include",
        })
            // .then((r=>console.log(r)))
            .then((res) => res.json())
            .then((json) => {
                // console.log(json);
                const posts = json.posts;
                this.commit("updatePosts", posts);
            });
    },
    //fetch current logged in user posts.
    async fetchMyPosts() {
        let id = "";
        await fetch("http://localhost:8081/currentUser", {
            //first get my ID
            credentials: "include",
        })
            .then((r) => r.json())
            .then((json) => {
                // console.log("get id - ", json);
                id = json.users[0].id;
            });
        await fetch("http://localhost:8081/userPosts?id=" + id, {
            //then fetch all posts with this ID
            credentials: "include",
        })
            .then((r) => r.json())
            .then((r) => {
                const myposts = r.posts;
                this.commit("updateMyPosts", myposts);
                // console.log(myposts);
            });
        // console.log("here")

        // .then((json) => console.log("get posts -", json));
    },

    async getMyUserID({ commit }) {

        await fetch("http://localhost:8081/currentUser", {
            credentials: "include",
        })
            .then((r) => r.json())
            .then((json) => {
                // console.log("JSON response", json)
                commit("updateMyUserID", json.users[0].id)
            });
    },

    async getMyProfileInfo(context) {
        await context.dispatch("getMyUserID");
        const userID = context.state.id;
        await fetch("http://localhost:8081/userData?userId=" + userID, {
            credentials: "include",
        })
            .then((r) => r.json())
            .then((json) => {
                let userInfo = json.users[0];
                // console.log(userInfo);
                this.commit("updateProfileInfo", userInfo);
                // console.log("userinfo -", json);
            });
    },

    async getAllUsers() {
        await fetch("http://localhost:8081/allUsers", {
            credentials: "include",
        })
            .then((r) => r.json())
            .then((json) => {
                let users = json.users;
                this.commit("updateAllUsers", users);
                // console.log("allUsers:", json.users);
            });
    },
    async getAllGroups() {
        await fetch("http://localhost:8081/allGroups", {
            credentials: "include",
        })
            .then((r) => r.json())
            .then((json) => {
                let groups = json.groups;
                this.commit("updateAllGroups", groups);
                // console.log("Allgroups:", json.groups);
            });
    },

    async getUserGroups(context) {
        const response = await fetch(`http://localhost:8081/userGroups`, {
            credentials: 'include'
        });

        const data = await response.json();
        // console.log("/getUserGroups data", data)
        // context.state.groups.userGroups.loaded = true;

        context.commit("updateUserGroups", data.groups)
        context.commit("updateDataLoaded", "userGroups")

    },

    addUserGroup({ state, commit }, userGroup) {
        let userGroups = state.groups.userGroups;
        console.log("userGroups state", userGroups)
        if (userGroups === null) { userGroups = [] };
        userGroups.push(userGroup);

        console.log("userGroup", userGroup)
        commit("updateUserGroups", userGroups)
    },

    async getMyFollowers(context) {
        await context.dispatch("getMyProfileInfo");
        const myID = context.state.profileInfo.id;

        const response = await fetch(`http://localhost:8081/followers?userId=${myID}`, {
            credentials: 'include'
        });

        const data = await response.json();

        context.commit("updateMyFollowers", data.users)

    },


    async getGroupPosts() {
        await fetch(
            "http://localhost:8081/groupPosts?groupId=" +
            router.currentRoute.value.params.id,
            {
                credentials: "include",
            }
        )
            .then((r) => r.json())
            .then((json) => {
                // console.log(json)
                let posts = json.posts;
                this.commit("updateGroupPosts", posts);
            });
    },

    async isLoggedIn() {
        const response = await fetch('http://localhost:8081/sessionActive', {
            credentials: 'include'
        });

        const data = await response.json();

        if (data.message === "Session active") {
            // console.log("ah yes")
            return true
        } else {
            // console.log("ah no")
            return false
        }

    },

    createWebSocketConn({ commit, dispatch, state }) {
        const ws = new WebSocket("ws://localhost:8081/ws");
      
        ws.addEventListener("message", (e) => {
            const data = JSON.parse(e.data);
            if (data.action == "chat") {
                // only broadcast messages when participants(sender and reciever) chat is open

                const isParticipantsChatOpen = state.chat.openChats.some((chat) => {
                    // Chat is open with the person who sent the message
                    if (data.chatMessage.type === "PERSON" && data.chatMessage.senderId  === chat.receiverId) {
                        return true
                    }

                    if (data.chatMessage.type === "GROUP" && data.chatMessage.receiverId === chat.receiverId) {
                        return true
                    }


                })
                if (isParticipantsChatOpen) {
                    dispatch("addNewChatMessage", data.chatMessage)
                    dispatch("markMessageRead", data.chatMessage)
                } else {
                    if (data.message === "NEW") {
                        dispatch("fetchChatUserList");
                    }
    
                    dispatch("addUnreadChatMessage", data.chatMessage)
                }
            } else if (data.action == "notification") {
                dispatch("addNewNotification", data.notification);

            } else if(data.action == "groupAccept"){
                dispatch("getUserGroups");
            }

        })

        commit("updateWebSocketConn", ws)  
    }

}