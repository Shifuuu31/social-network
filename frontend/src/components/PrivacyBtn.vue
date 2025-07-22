<template>
    <div class="user-profile__privacy">
        Private profile
        <label for="privateProfileInput">
            <div class="toggle-wrapper" :class="{ active: isActive }">
                <span class="toggle-button"></span>
                <input type="checkbox" id="privateProfileInput"
                       v-model="checked"
                       @change="updateProfileStatus">
            </div>
        </label>

    </div>

</template>


<script>
export default {
    props: ['status'],
    data() {
        return {
            currentUserStatus: this.status,
            checked: false
        }
    },

    computed: {
        isActive() {
            if (this.currentUserStatus === "PRIVATE") {
                this.checked = true;
                return true;
            } else {
                return false;
            }
        }
    },

    methods: {
        async updateProfileStatus() {
            this.currentUserStatus = (this.checked ? 'PRIVATE' : 'PUBLIC');
            const response = await fetch(`http://localhost:8081/changeStatus?status=${this.currentUserStatus}`, {
                credentials: "include"
            });
            // console.log("Response", await response.json())

        }
    }
}

</script>


<style>
.user-profile__privacy {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 7.5px;

    padding-top: 25px;
    border-top: 1px solid rgb(212, 212, 212);
    width: 100%;
    font-size: 14px;
}

.user-profile__privacy label {
    margin: 0;
}

.toggle-wrapper {
    width: 60px;
    height: 30px;
    /* border: 1px solid black; */
    box-shadow: 0 0 0 1px rgb(221, 221, 221);
    border-radius: 50px;
    position: relative;
    background-color: rgb(228, 228, 228);
}

.toggle-button {
    display: inline-block;
    height: 30px;
    width: 30px;
    background-color: rgb(184, 184, 184);
    border-radius: 50px;
    position: absolute;
    /* left: 0; */
    /* top: 0; */
    left: 0;
    transition: all 0.3s ease;
}

.active .toggle-button {
    transform: translateX(100%);
    background-color: rgb(0, 158, 0);
}

#privateProfileInput {
    opacity: 0;
}



/* .toggle-indicator:hover {
    transform: translateX(100%);
} */
</style>