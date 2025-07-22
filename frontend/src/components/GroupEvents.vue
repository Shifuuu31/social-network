<template>
    <div class="item-list__wrapper" id="groups">
        <h3>Events</h3>
        <ul class="item-list">
            <li v-for="event in this.groupEvents">
                <img class="small" src="../assets/icons/event.svg" alt="">
                <div class="item-text" style="cursor:pointer" @click="this.showEvent(event)">{{ event.title }}</div>

                <Modal v-if="this.eventIsOpen" @closeModal="toggleEventModal">
                    <template #title>
                        {{ this.eventData.title }}
                    </template>
                    <template #body>
                        <div>{{ this.eventData.date }}</div>
                        <div>{{ this.eventData.content }}</div>
                        <form @submit.prevent="participateRespond(); toggleEventModal();" id="event-desc">
                            <div class="form-input">
                            <label for="goingEvent">Going</label>

                            <div class="select-wrapper">
                                <img src="../assets/icons/angle-down.svg" class="dropdown-arrow">
                                <select v-model="this.formData.going" name="goingEvent" id="goingEvent" required>
                                    <option value="" selected hidden>{{this.eventData.going}}</option>
                                    <option value="YES">Yes</option>
                                    <option value="NO">No</option>
                                </select>

                            </div>

                        </div>
                        </form>
                        <button class="btn form-submit" form="event-desc">Confirm</button>
                    </template>
                </Modal>
            </li>
        </ul>
        <button class="btn form-submit" @click="toggleModal">New event<i class="uil uil-plus"></i></button>

        <Modal v-if="this.isOpen" @closeModal="toggleModal">
            <template #title>Create new event</template>
            <template #body>
                <form @submit.prevent="createNewEvent(); toggleModal();" id="new-event" ref="openForm">
                    <div class="form-input">
                        <label for="title">Title</label>
                        <input type="text" v-model="this.formData.title" id="title" required>
                    </div>
                    <div class="form-input">
                        <label for="description">Description</label>
                        <textarea v-model="this.formData.content" cols="30" rows="3"
                                  placeholder="What is this about?" id="description" required></textarea>
                    </div>
                    <div class="form-input">
                        <label for="date">Date</label>
                        <input v-model="this.formData.date" type="date" id="date" required>
                    </div>
                    <div class="form-input">
                        <label for="going">Going</label>

                        <div class="select-wrapper">
                            <img src="../assets/icons/angle-down.svg" class="dropdown-arrow">
                            <select v-model="this.formData.going" name="going" id="going" required>
                                <option value="" selected hidden>Choose here</option>
                                <option value="YES">Yes</option>
                                <option value="NO">No</option>
                            </select>

                        </div>

                    </div>
                </form>

                <button class="btn form-submit" form="new-event">Create</button>
            </template>

        </Modal>
    </div>
</template>


<script>
import Modal from './Modal.vue';
export default {
    name: "GroupEvents",
    data() {
        return {
            groupEvents: [],
            isOpen: false,
            eventIsOpen: false,
            eventData: {
                id:"",
                title: "",
                date: "",
                content: "",
                going: "",
            },
            formData: {
                id:"",
                title: "",
                content: "",
                date: null,
                going: ""
            },
        };
    },
    created() {
        this.getGroupEvents();
    },
    watch: {
        $route() {
            this.getGroupEvents()
        }
    },
    methods: {
        async getGroupEvents() {
            await fetch("http://localhost:8081/groupEvents?groupId=" + this.$route.params.id, {
                credentials: "include"
            })
                .then((response => response.json()))
                .then((json => {
                    this.groupEvents = json.events
                }));
        },
        async participateRespond(){
            const response = await fetch(`http://localhost:8081/participate`, {
                credentials: "include",
                method: "POST",
                body: JSON.stringify({
                    eventId: this.eventData.id,
                    response: this.formData.going,
                })
            });
            const data = await response.json();
            if (data.type == "Success"){
                // reload events
                this.getGroupEvents();
                // reload notifications
                this.fetchNotifications()
            }
            
        },
        async fetchNotifications() {
            const response = await fetch("http://localhost:8081/notifications", {
                credentials: "include"
            });
            const data = await response.json();
            this.notificationsFromDB = data;
            this.$store.commit("updateAllNotifications", data.notifications);
            // console.log("/notifications data", data)
        },
        async createNewEvent() {

            await fetch('http://localhost:8081/newEvent', {
                method: 'POST',
                credentials: 'include',
                body: JSON.stringify({
                    groupID: this.$route.params.id,
                    title: this.formData.title,
                    content: this.formData.content,
                    date: this.formData.date,
                    going: this.formData.going
                })
            })
                .then((r => r.json()))
                // .then((json => {
                //     // console.log("newEvent", json);
                // }))
            this.getGroupEvents()
        },
        toggleModal() {
            if (this.isOpen) {
                // clear form data
                this.formData = {
                    id:"",
                    title: "",
                    content: "",
                    date: null,
                    going: ""
                }
            }
            this.isOpen = !this.isOpen;
        },
        toggleEventModal() {
            console.log("toggle", this.eventIsOpen)
            if (this.eventIsOpen) {
                // clear form data
                this.formData = {
                    id:"",
                    title: "",
                    content: "",
                    date: null,
                    going: ""
                }
            }
            this.eventIsOpen = !this.eventIsOpen;
        },

        showEvent(event) {
            this.eventData.id = event.id
            this.eventData.title = event.title
            this.eventData.date = event.date
            this.eventData.content = event.content
            this.eventData.going = event.going

            this.eventIsOpen = true  
        },
        closeEvent() {
            this.eventIsOpen = false
        },
    },
    components: { Modal }
}
</script>


<style>
</style>