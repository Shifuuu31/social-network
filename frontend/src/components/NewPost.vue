

<template>

    <button class="start-post" @click="toggleModal">
        <span>Start post</span>
        <i class="uil uil-edit"></i>
    </button>

    <Modal v-show="isOpen" @closeModal="toggleModal">
        <template #title>Create a post</template>
        <template #body>
            <form @submit.prevent="submitNewPost" id="newpost">
                <div class="form-input" v-if="!this.isGroupPage">
                    <label for="post_privacy">Post privacy</label>
                    <div class="select-wrapper">
                        <img src="../assets/icons/angle-down.svg" class="dropdown-arrow">

                        <select id="post_privacy" v-model="newpost.privacy" required>
                            <option value="" selected hidden>Choose here</option>
                            <option value="public">Everyone</option>
                            <option value="private">Followers</option>
                            <option value="almost-private">Choosen followers</option>
                        </select>

                    </div>

                    <MultiselectDropdown v-if="newpost.privacy === 'almost-private'"
                                         v-model:checkedOptions="newpost.checkedFollowers"
                                         placeholder="Select followers"
                                         :content="getMyFollowersList" />
                </div>

                <div class="form-input">
                    <label for="description">Description</label>

                    <textarea id="description" v-model="newpost.body" rows="4" cols="50"
                              placeholder="What are you thinking?"  required></textarea>

                </div>

                <FileUpload v-model:file=newpost.image
                            @inputCleared="toggleClearInput"
                            :clearInput="clearInput"
                            labelName="Image" />

            </form>
            <button class="btn submitPost" type="submit" form="newpost">Post</button>

        </template>
    </Modal>

</template>



<script>
import Modal from './Modal.vue'
import MultiselectDropdown from './MultiselectDropdown.vue';
import FileUpload from './FileUpload.vue';
export default {
    components: {
        Modal,
        MultiselectDropdown,
        FileUpload
    },
    name: 'Newpost',
    data() {
        return {
            isOpen: false,
            isGroupPage: null,
            newpost: {
                privacy: "",
                body: "",
                checkedFollowers: null,
                image: null,
            },
            clearInput: false,
        }
    },

    created() {
        this.getMyFollowers();
        this.isGroupPageCheck()
    },

    computed: {
        getMyFollowersList() {
            return this.$store.getters.followers;
        },

    },

    methods: {
        toggleModal() {
            // if modal was open, clear the form
            if (this.isOpen) { this.clearForm(); }
            this.isOpen = !this.isOpen
        },

        toggleClearInput() {
            this.clearInput = !this.clearInput
        },



        getMyFollowers() {
            this.$store.dispatch("getMyFollowers")
        },

        clearForm() {
            this.newpost.privacy = "";
            this.newpost.body = "";
            // this.newpost.image = null;
            this.toggleClearInput();
        },

        isGroupPageCheck() {
            if (this.$route.path.includes("group")) {
                this.isGroupPage = true
            } else {
                this.isGroupPage = false
            }

        },

        async submitPost() {
            let formData = new FormData();
            formData.set("body", this.newpost.body)
            formData.set("image", this.newpost.image)
            formData.set("privacy", this.newpost.privacy)
            if (this.newpost.checkedFollowers != null){
                formData.set("checkedfollowers", this.newpost.checkedFollowers.map(x => x.id))
            }
            const response = await fetch('http://localhost:8081/newPost', {
                method: 'POST',
                credentials: 'include',
                body: formData
            })
            await this.$store.dispatch('fetchPosts')
            console.log('Post submitted', await response.json());
            this.toggleModal();
        },

        async submitGroupPost() {

            let formData = new FormData();
            formData.set('groupId', this.$route.params.id)
            formData.set('body', this.newpost.body);
            formData.set('image', this.newpost.image);

            await fetch('http://localhost:8081/newGroupPost', {
                method: 'POST',
                credentials: 'include',
                body: formData
            })
                .then((r => r.json()))
            // .then((json => console.log(json)))
            await this.$store.dispatch('getGroupPosts')
            this.toggleModal();
            console.log('Group Post Submitted');
        },

        submitNewPost() {
            if (this.isGroupPage) {
                this.submitGroupPost();
            } else {
                this.submitPost();
            }
        }
    }
}
</script>


<style scoped>
/* #newpost {
    display: flex;
    flex-direction: column;
} */


/* 
#postBtn {
    width: 600px;
    display: flex;
    justify-content: center;
    border: 1px solid #706A6A;
    border-radius: 3px;
    border: 1px solid #706A6A;
} */

.start-post {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 10px 15px;
    background-color: var(--input-bg);
    border: none;
    box-shadow: var(--container-shadow);
    /* box-shadow: var(--hover-box-shadow); */
    /* box-shadow: 0 0 2px 2px black; */
    font-family: inherit;
    font-size: 16px;
    border-radius: var(--container-border-radius);
    cursor: pointer;
    /* transition: var(--hover-box-shadow-transition); */

    transition: box-shadow 0.15s ease-in-out;

    width: 100%;
}

.start-post i {
    font-size: 1.25em;
}


.additional-info {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
}

.btns-wrapper {
    display: flex;
    flex-direction: column;
    align-items: flex-end;
    gap: 10px;
}

.submitPost {
    margin-left: auto;
}

.start-post:hover {
    /* box-shadow: none; */
    box-shadow: var(--hover-box-shadow);
}
</style>