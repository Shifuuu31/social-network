<template>
    <div class="item-list__wrapper" id="groups">
        <h3>Members</h3>
        <ul class="item-list">
            <li v-for="member in this.groupMembers" :key="member.id">
                <img class="small" src="../assets/icons/default-profile.svg" alt="">
                <div class="item-text"><router-link :to="{ path: `/profile/${member.id}`}">{{ member.nickname }}
</router-link></div>
            </li>
        </ul>
        <button v-if="this.isMember" class="btn form-submit" @click="toggleModal();getFollowers()">Invite users<i class="uil uil-user-plus"></i></button>
        <button v-if="!this.isMember" class="btn form-submit" @click="this.joinGroup">Join group +</button>

        <Modal v-if="this.isOpen" @closeModal="toggleModal">
            <template #title>Invite users</template>
            <template #body>

                <MultiselectDropdown v-model:checkedOptions="checkedNames" placeholder="Select followers"
                    :content="listForShowing" />
                <button class="btn form-submit" @click="toggleModal() ; inviteUsersToGroup()">Invite</button>
            </template>
        </Modal>
    </div>
</template>


<script>
import Modal from './Modal.vue';
import MultiselectDropdown from './MultiselectDropdown.vue';
export default {
    name: "GroupMembers",
      props: {
        isMember: false
    },
    data() {
        return {
            groupMembers: null,
            isOpen: false,
            followers: [],
            listForShowing: [],
            allUsers: [],
            checkedNames: [],
            clearInput: false,
        };
    },
    created() {
        this.getGroupMembers();
        // this.getFollowers();
    },

    computed: {
        allFollowersNames() {
            return this.listForShowing
        }
    },

    watch: {
        $route() {
            this.getGroupMembers();
            // this.getFollowers();
        }
    },
  
    methods: {
        async getFollowers() {
            this.$store.dispatch("getMyFollowers");
            this.createFollowersListForShowing(this.$store.state.myFollowers, this.groupMembers)

        },
        async getGroupMembers() {
            await fetch("http://localhost:8081/groupMembers?groupId=" + this.$route.params.id, {
                credentials: "include"
            })
                .then((response => response.json()))
                .then((json => {
                    // console.log("GroupMembers:", json);
                    this.groupMembers = json.users;
                }));
        },
        createFollowersListForShowing(followers, members) {
            this.listForShowing = [];
            let isUserInGroup = false
            for (let i = 0; i < Object.keys(followers).length; i++) {
                for (let j = 0; j < Object.keys(members).length; j++) {
                    if (followers[i].nickname === members[j].nickname) {
                        isUserInGroup = true
                    }
                }
                if (!isUserInGroup) {
                    this.listForShowing.push(followers[i])
                }
                isUserInGroup = false
            }
        },
        toggleModal() {
            this.isOpen = !this.isOpen;
        },

        getIds() {
            let arrOfIDS = [];
            for (let name of this.checkedNames) {
                for (let obj of this.listForShowing) {
                    if (obj.nickname === name.nickname) {
                        arrOfIDS.push(obj.id)
                    }
                }
            }
            return arrOfIDS

        },
        async joinGroup(){
            await fetch("http://localhost:8081/newGroupRequest?groupId=" + this.$route.params.id, {
                credentials: 'include',
            })
            .then(response=>response.json())
            .then(json=>{
                // console.log(json);
                this.$toast.open({
                            message: json.message,
                            type: "success",
                        });
            })
        },


        async inviteUsersToGroup() {
            await fetch("http://localhost:8081/newGroupInvite", {
                method: 'POST',
                credentials: 'include',
                body: JSON.stringify({ invitations: this.getIds(), id: this.$route.params.id })
            })
                .then((response => response.json()))
                .then((json => {
                    // console.log("new group invite response:", json);
                    this.clearInput = true;
                    // this.groupMembers = json.users;
                }));
        },
    },
    components: { Modal, MultiselectDropdown }
}
</script>


<style>
</style>