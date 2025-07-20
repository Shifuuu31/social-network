<template>
  <div class="groups-page">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">Groups</h1>
        <p class="page-subtitle">Find communities that share your interests</p>
      </div>

      <div class="groups-filters">
        <div class="search-bar">
          <input type="text" placeholder="Search groups..." v-model="searchQuery" @input="handleSearchInput"
            class="search-input" />
          <button class="search-btn" @click="performSearch">
            <span class="icon">🔍</span>
          </button>
        </div>

        <div class="filter-buttons">
          <button :class="['filter-btn', { active: activeFilter === 'all' }]" @click="setFilter('all')">
            <span class="icon">🔍</span>
            Explore
          </button>
          <button :class="['filter-btn', { active: activeFilter === 'joined' }]" @click="setFilter('joined')">
            <span class="icon">👥</span>
            My Groups
          </button>
        </div>
      </div>

      <div v-if="groupsStore.isLoading" class="loading">
        <div class="spinner"></div>
        <p>Loading groups...</p>
      </div>
      <!-- TODO handle errors styling better -->
      <div v-else-if="groupsStore.error" class="error">
        <p>{{ groupsStore.error }}</p>
        <button @click="loadGroups" class="btn btn-primary">Retry</button>
      </div>

      <div v-else class="groups-grid">
        <GroupCard v-for="group in groupsStore.groups" :key="group.id" :group="group" @group-joined="handleGroupJoined"
          @group-left="handleGroupLeft" />
      </div>

      <div v-if="groupsStore.groups.length === 0 && !groupsStore.isLoading" class="empty-state">
        <div class="empty-icon">📭</div>
        <h3>No groups found</h3>
        <p v-if="activeFilter === 'joined'">
          You haven't interacted with any groups yet. Explore some groups to get started!
        </p>
        <p v-else-if="searchQuery.trim()">
          No groups found matching "{{ searchQuery.trim() }}". Try different keywords or check back later.
        </p>
        <p v-else>
          No new groups to explore. Try changing your search criteria or check back later.
        </p>
        <router-link to="/groups/create" class="btn btn-primary">
          Create a Group
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useGroupsStore } from '@/stores/groups'
import GroupCard from '@/components/GroupCard.vue'

const groupsStore = useGroupsStore()
const searchQuery = ref('')
const activeFilter = ref('all')
let searchTimeout = null

// Handle search input with debouncing
const handleSearchInput = () => {
  clearTimeout(searchTimeout)
  searchTimeout = setTimeout(() => {
    loadGroups()
  }, 200)
}

// Perform immediate search when search button is clicked
const performSearch = () => {
  clearTimeout(searchTimeout)
  loadGroups()
}

// Watch for filter changes
watch(activeFilter, () => {
  loadGroups()
})

const setFilter = (filter) => {
  activeFilter.value = filter
}

const loadGroups = async () => {
  const filterType = activeFilter.value === 'joined' ? 'user' : 'all'
  const searchTerm = searchQuery.value.trim()
  await groupsStore.fetchGroups(filterType, searchTerm)
}

const handleGroupJoined = (groupId) => {
  // Handle group joined event
  loadGroups()
}

const handleGroupLeft = (groupId) => {
  // Handle group left event
  loadGroups()
}

onMounted(() => {
  loadGroups()
})
</script>

<style scoped>
.groups-page {
  padding: 40px 20px;
}

.container {
  max-width: 1200px;
  margin: 0 auto;
}

.page-header {
  text-align: center;
  margin-bottom: 40px;
}

.page-title {
  font-size: 2.5rem;
  font-weight: 700;
  color: #fff;
  margin-bottom: 10px;
}

.page-subtitle {
  font-size: 1.1rem;
  color: #ccc;
}

.groups-filters {
  display: flex;
  gap: 20px;
  margin-bottom: 40px;
  flex-wrap: wrap;
}

.search-bar {
  flex: 1;
  min-width: 300px;
  position: relative;
}

.search-input {
  width: 100%;
  padding: 16px 50px 16px 16px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  color: #fff;
  font-size: 1rem;
  transition: border-color 0.2s ease;
}

.search-input:focus {
  outline: none;
  border-color: #8b5cf6;
}

.search-btn {
  position: absolute;
  right: 12px;
  top: 50%;
  transform: translateY(-50%);
  background: none;
  border: none;
  color: #ccc;
  cursor: pointer;
  padding: 8px;
  border-radius: 6px;
  transition: color 0.2s ease;
}

.search-btn:hover {
  color: #fff;
}

.filter-buttons {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-btn {
  padding: 12px 20px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: #ccc;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9rem;
  font-weight: 500;
}

.filter-btn:hover {
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
}

.filter-btn.active {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff;
  border-color: transparent;
}

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 24px;
  margin-bottom: 40px;
}

.loading {
  text-align: center;
  padding: 60px 20px;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 4px solid rgba(255, 255, 255, 0.1);
  border-top: 4px solid #8b5cf6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin: 0 auto 20px;
}

@keyframes spin {
  0% {
    transform: rotate(0deg);
  }

  100% {
    transform: rotate(360deg);
  }
}

.error {
  text-align: center;
  padding: 60px 20px;
  color: #ef4444;
}

.empty-state {
  text-align: center;
  padding: 60px 20px;
}

.empty-icon {
  font-size: 4rem;
  margin-bottom: 20px;
}

.empty-state h3 {
  font-size: 1.5rem;
  color: #fff;
  margin-bottom: 10px;
}

.empty-state p {
  color: #ccc;
  margin-bottom: 30px;
}

.btn {
  padding: 14px 28px;
  border-radius: 10px;
  border: none;
  font-weight: 600;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
}

.btn-primary {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff;
}

.btn-primary:hover {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
  transform: translateY(-2px);
}

@media (max-width: 768px) {
  .groups-filters {
    flex-direction: column;
  }

  .search-bar {
    min-width: auto;
  }

  .filter-buttons {
    justify-content: center;
  }

  .groups-grid {
    grid-template-columns: 1fr;
  }
}
</style>