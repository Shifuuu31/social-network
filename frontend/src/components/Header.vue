<template>
  <header class="header">
    <div class="header-container">
      <div class="header-left">
        <router-link to="/" class="logo">
          <span class="logo-icon">🌐</span>
          <span class="logo-text">Social NETWORK</span>
        </router-link>
      </div>
      
      <nav class="header-nav">
        <router-link to="/" class="nav-link">Accueil</router-link>
        <router-link to="/groups" class="nav-link">Groupes</router-link>
        <router-link to="/groups/create" class="nav-link nav-link-primary">
          Créer un groupe
        </router-link>
        
        <!-- Notification Bell -->
        <NotificationDropdown />
      </nav>
    </div>
  </header>
</template>

<script setup>
import NotificationDropdown from './notifications.vue'
</script>

<style scoped>
.header {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  background: rgba(10, 10, 10, 0.95);
  backdrop-filter: blur(10px);
  border-bottom: 1px solid #333;
  z-index: 100;
}

.header-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 0 20px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 70px;
}

.header-left {
  display: flex;
  align-items: center;
}

.logo {
  display: flex;
  align-items: center;
  gap: 10px;
  text-decoration: none;
  color: #fff;
  font-size: 1.5rem;
  font-weight: 700;
}

.logo-icon {
  font-size: 1.8rem;
}

.logo-text {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
}

.header-nav {
  display: flex;
  align-items: center;
  gap: 20px;
}

.nav-link {
  color: #ccc;
  text-decoration: none;
  font-weight: 500;
  padding: 8px 16px;
  border-radius: 8px;
  transition: all 0.2s ease;
}

.nav-link:hover {
  color: #fff;
  background: rgba(255, 255, 255, 0.1);
}

.nav-link.router-link-active {
  color: #8b5cf6;
}

.nav-link-primary {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff !important;
}

.nav-link-primary:hover {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
  transform: translateY(-1px);
}

@media (max-width: 768px) {
  .header-container {
    padding: 0 16px;
  }
  
  .header-nav {
    gap: 12px;
  }
  
  .nav-link {
    padding: 6px 12px;
    font-size: 0.9rem;
  }
  
  .logo-text {
    display: none;
  }
}
</style>