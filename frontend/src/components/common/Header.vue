<template>
  <header class="header">
    <!-- Logo -->
    <router-link to="/" class="logo">
      Introvia
    </router-link>

    <!-- Search Bar -->
    <div class="search-container">
      <input type="text" placeholder="Search..." class="search-input" />
    </div>

    <!-- Right Icons -->
    <div class="icons">
      <!-- Notification Icon -->
      <div class="icon-container">
        <BellIcon class="icon" />
        <span class="badge red">3</span>
      </div>

      <!-- Message Icon -->
      <div class="icon-container">
        <MessageIcon class="icon" />
        <span class="badge green">2</span>
      </div>

      <!-- User Dropdown -->
      <div class="user-dropdown" @click="toggleDropdown">
        <img
          src="https://via.placeholder.com/32"
          alt="User Avatar"
          class="avatar"
        />
        <div v-if="dropdownOpen" class="dropdown-menu">
          <router-link to="/profile" class="dropdown-item">Profile</router-link>
          <button @click="logout" class="dropdown-item">Logout</button>
        </div>
      </div>
    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'
// import { BellIcon, MessageSquareIcon as MessageIcon } from 'lucide-vue-next'

const dropdownOpen = ref(false)
const router = useRouter()

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value
}

function logout() {
  localStorage.removeItem('token')
  router.push('/login')
}
</script>

<style scoped>
.header {
  background-color: white;
  padding: 12px 16px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
}

.logo {
  font-size: 20px;
  font-weight: bold;
  color: #2563eb;
  text-decoration: none;
}

.search-container {
  display: none;
  width: 33%;
}
@media (min-width: 768px) {
  .search-container {
    display: flex;
  }
}
.search-input {
  width: 100%;
  padding: 6px 10px;
  border: 1px solid #ccc;
  border-radius: 6px;
  outline: none;
}

.icons {
  display: flex;
  align-items: center;
  gap: 16px;
}

.icon-container {
  position: relative;
  cursor: pointer;
}

.icon {
  width: 24px;
  height: 24px;
  color: #444;
}

.badge {
  position: absolute;
  top: -6px;
  right: -6px;
  width: 16px;
  height: 16px;
  font-size: 10px;
  background-color: red;
  color: white;
  border-radius: 50%;
  display: flex;
  justify-content: center;
  align-items: center;
}
.badge.green {
  background-color: green;
}
.badge.red {
  background-color: red;
}

.avatar {
  width: 32px;
  height: 32px;
  border-radius: 50%;
  cursor: pointer;
  border: 1px solid #ccc;
}

.user-dropdown {
  position: relative;
}

.dropdown-menu {
  position: absolute;
  right: 0;
  top: 40px;
  width: 160px;
  background-color: white;
  border-radius: 6px;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.1);
  z-index: 10;
}

.dropdown-item {
  display: block;
  width: 100%;
  padding: 10px 16px;
  text-align: left;
  background: none;
  border: none;
  cursor: pointer;
  text-decoration: none;
  color: black;
}

.dropdown-item:hover {
  background-color: #f1f1f1;
}
</style>
