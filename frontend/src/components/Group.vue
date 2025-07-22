<template>
    <div class="content" v-if="groupData">

        <div class="left-section">
            <GroupMembers v-bind:isMember="isMemberOfGroup" />
            <GroupEvents v-if="this.isMemberOfGroup" />
                 
        </div>

        <div class="middle-section">
            <div class="about">
                <h2 class="about-title">{{ this.groupData.name }}</h2>
                <p class="about-text">{{ this.groupData.description }}</p>
            </div>
            <NewPost v-if="this.isMemberOfGroup" />
            <GroupPosts v-if="this.isMemberOfGroup" />
            <p class="additional-info large" v-if="!this.isMemberOfGroup">Only group members can see additional
                information.
            </p>            
        </div>
      
        <GroupJoinRequests v-bind:isAdmin="this.isAdmin" class="right-section"/>  
     
    </div>
</template>

<script>
import AllPosts from './AllPosts.vue'
import Groups from './Groups.vue';
import Notifications from './Notifications.vue';
import NewPost from './NewPost.vue';
import GroupPosts from './GroupPosts.vue';
import GroupMembers from './GroupMembers.vue';
import Modal from './Modal.vue';
import GroupEvents from './GroupEvents.vue';
import GroupJoinRequests from './GroupJoinRequests.vue';
export default {
    name: "Group",
    created() {
        this.getGroupInfo();
    },
    watch: {
        $route() {
            if (this.$route.path.includes("group")){
                this.isMemberOfGroup=false;
                this.getGroupInfo(); 
            }
            // this.getGroupInfo()
        }
    },
    data() {
        return {
            groupData: null,
            isMemberOfGroup: false,
            isAdmin:false,
        };
    },
    methods: {
        async getGroupInfo() {
            await fetch("http://localhost:8081/groupInfo?groupId=" + this.$route.params.id, {
                credentials: "include"
            })
                .then((r => r.json()))
                .then((json => {
                    // console.log("/groupInfo response", json);
                    this.groupData = json.groups[0];
                    if (json.groups[0].admin === true || json.groups[0].member === true) {
                        this.isMemberOfGroup = true
                    }
                    if(json.groups[0].admin === true){
                        this.isAdmin = true
                    }
                }));
        },
    },
    components: { AllPosts, Groups, Notifications, NewPost, GroupPosts, GroupMembers, Modal, GroupEvents, GroupJoinRequests }
}
</script>


<style scoped>
.content {
    margin-top: 50px;
    padding: 0 30px;
    display: grid;
    grid-template-columns: minmax(0, 1fr) minmax(0, 550px) minmax(0, 1fr);
    column-gap: 50px;
}



.middle-section {
    justify-self: center;
    display: flex;
    flex-direction: column;
    gap: 35px;


}


.left-section {
    justify-self: flex-end;
    display: flex;
    flex-direction: column;
    gap: 35px;
}

.right-section {
    justify-self: flex-start;
    min-width: 250px;
  
}


.box {
    height: 300px;
    width: 550px;
    border: 2px solid blue;
}

@media only screen and (max-width: 1250px) {
    .content {
        grid-template-columns: minmax(min-content, max-content) minmax(min-content, 550px);
        grid-template-rows: repeat(2, minmax(auto, max-content));
        row-gap: 35px;
        justify-content: center;
        grid-template-areas:
            "left-section middle-section"
            "right-section middle-section"
            "... middle-section";

    }


    .middle-section {
        grid-area: middle-section;
    }

    .left-section {
        grid-area: left-section;
    }

    .right-section {
        grid-area: right-section;
    }

}
</style>