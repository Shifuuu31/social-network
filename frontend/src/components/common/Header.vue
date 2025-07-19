<template>
  <header class="header">
    <!-- Home Icon -->
    <router-link to="/" class="icon-container" aria-label="Home">
      <font-awesome-icon icon="fa-solid fa-home" class="icon" />
    </router-link>

    <!-- Logo -->
    <router-link to="/" class="logo">
      Introvia
    </router-link>

    <!-- Right Icons -->
    <div class="icons">
      <!-- Groups Icon -->
      <router-link to="/groups" class="icon-container" aria-label="Groups">
        <font-awesome-icon icon="fa-solid fa-users" class="icon" />
      </router-link>



      <router-link to="/profile" class="icon-container" aria-label="Profile">
        <font-awesome-icon icon="fa-solid fa-user" class="icon" />
        <span class="badge green">1</span>
      </router-link>
  
      <!-- Friends Icon -->
      <router-link to="/discover-friend" class="icon-container" aria-label="Friends">
        <font-awesome-icon icon="fa-solid fa-user-friends" class="icon" />
      </router-link>

      <!-- Notification Icon -->
      <router-link to="/notifications" class="icon-container" aria-label="Notifications">
        <font-awesome-icon icon="fa-solid fa-bell" class="icon" />
        <span class="badge red">3</span>
      </router-link>
 
      <!-- Chat Icon -->
      <router-link to="/chat" class="icon-container" aria-label="Chat">
        <span class="chat-text">💬</span>
        <span class="badge green">2</span>
      </router-link>
      <!-- Profile Icon -->
 

    </div>
  </header>
</template>

<script setup>
import { ref } from 'vue'
import { useRouter } from 'vue-router'

const dropdownOpen = ref(false)
const router = useRouter()

function toggleDropdown() {
  dropdownOpen.value = !dropdownOpen.value
}

function logout() {
  // localStorage.removeItem('token')
  router.push('/login')
}

function goToProfile() {
  router.push('/profile')
}
</script>

<style scoped>
.header {
  background: rgba(255, 255, 255, 0.85);
  backdrop-filter: blur(8px);
  border-radius: 0;
  margin: 0;
  max-width: 100vw;
  /* width: 95%; */
  padding: 14px 48px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  box-shadow: 0 6px 24px rgba(0, 0, 0, 0.08), 0 1.5px 4px rgba(0,0,0,0.04);
  position: sticky;
  top: 0;
  z-index: 100;
  transition: box-shadow 0.2s;
}

@media (max-width: 768px) {
  .header {
    padding: 10px 8px;
  }
}

.logo {
  font-size: 1.6rem;
  font-weight: 800;
  color: #2563eb;
  letter-spacing: 1px;
  text-decoration: none;
  margin-left: 12px;
  transition: color 0.2s;
}
.logo:hover {
  color: #1d4ed8;
}

.icons {
  display: flex;
  align-items: center;
  gap: 24px;
}

.icon-container {
  position: relative;
  cursor: pointer;
  border-radius: 50%;
  padding: 8px;
  transition: background 0.18s;
  display: flex;
  align-items: center;
  justify-content: center;
}
.icon-container:hover {
  background: rgba(37, 99, 235, 0.08);
}

.icon {
  width: 26px;
  height: 26px;
  color: #2563eb;
  transition: color 0.18s;
}
.icon-container:hover .icon {
  color: #1d4ed8;
}

.chat-text {
  font-size: 20px;
  color: #2563eb;
  transition: color 0.18s;
}
.icon-container:hover .chat-text {
  color: #1d4ed8;
}

.badge {
  position: absolute;
  top: 2px;
  right: 2px;
  min-width: 14px;
  height: 16px;
  font-size: 10px;
  background: linear-gradient(90deg, #ef4444, #f59e42);
  color: white;
  border-radius: 8px;
  display: flex;
  justify-content: center;
  align-items: center;
  padding: 0 4px;
  font-weight: 700;
  box-shadow: 0 1px 4px rgba(0,0,0,0.08);
  border: 2px solid #fff;
}
.badge.green {
  background: linear-gradient(90deg, #22c55e, #3b82f6);
}
.badge.red {
  background: linear-gradient(90deg, #ef4444, #f59e42);
}

.avatar {
  width: 36px;
  height: 36px;
  border-radius: 50%;
  cursor: pointer;
  border: 2px solid #2563eb;
  box-shadow: 0 2px 8px rgba(37,99,235,0.08);
  transition: border 0.18s, box-shadow 0.18s;
}
.avatar:hover {
  border: 2px solid #1d4ed8;
  box-shadow: 0 4px 16px rgba(37,99,235,0.12);
}

.user-dropdown {
  position: relative;
  margin-left: 18px;
}

.dropdown-menu {
  position: absolute;
  right: 0;
  top: 48px;
  width: 170px;
  background: rgba(255,255,255,0.98);
  border-radius: 12px;
  box-shadow: 0 8px 32px rgba(0,0,0,0.12);
  z-index: 10;
  padding: 8px 0;
  animation: fadeIn 0.18s;
}

@keyframes fadeIn {
  from { opacity: 0; transform: translateY(-8px);}
  to { opacity: 1; transform: translateY(0);}
}

.dropdown-item {
  display: block;
  width: 100%;
  padding: 12px 20px;
  text-align: left;
  background: none;
  border: none;
  cursor: pointer;
  text-decoration: none;
  color: #222;
  font-weight: 500;
  border-radius: 8px;
  transition: background 0.16s, color 0.16s;
}

.dropdown-item:hover {
  background: #f1f5fd;
  color: #2563eb;
}
</style>
