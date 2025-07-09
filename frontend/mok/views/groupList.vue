<!-- views/GroupsList.vue -->
<template>
  <div class="groups-page">
    <div class="page-header">
      <div class="header-content">
        <h1 class="page-title">Groups</h1>
        <p class="page-subtitle">Discover and join communities</p>
      </div>
      <button class="btn btn-create" @click="showCreateGroup = true">
        <i class="icon-plus"></i>
        create a group
      </button>
    </div>

    <div class="filters-section">
      <div class="search-bar">
        <div class="search-input-wrapper">
          <i class="icon-search"></i>
          <input 
            v-model="searchQuery"
            type="text" 
            placeholder="Rechercher un groupe..."
            class="search-input"
            @input="handleSearch"
          />
        </div>
      </div>
      
      <div class="filter-tabs">
        <button 
          v-for="tab in filterTabs"
          :key="tab.key"
          :class="['filter-tab', { active: activeFilter === tab.key }]"
          @click="activeFilter = tab.key"
        >
          {{ tab.label }}
        </button>
      </div>
    </div>

    <div class="groups-container">
      <div v-if="isLoading" class="loading-state">
        <div class="spinner"></div>
        <p>Loading groups...</p>
      </div>

      <div v-else-if="error" class="error-state">
        <div class="error-icon">⚠️</div>
        <h3>Loading error</h3>
        <p>{{ error }}</p>
        <button class="btn btn-secondary" @click="loadGroups">
          Retry
        </button>
      </div>

      <div v-else-if="filteredGroups.length === 0" class="empty-state">
        <div class="empty-icon">📭</div>
        <h3>No groups found</h3>
        <p v-if="searchQuery">
          No groups match your search "{{ searchQuery }}"
        </p>
        <p v-else>
          There are no groups in this category yet.
        </p>
      </div>

      <div v-else class="groups-grid">
        <GroupCard 
          v-for="group in filteredGroups"
          :key="group.id"
          :group="group"
          @group-joined="handleGroupJoined"
          @group-left="handleGroupLeft"
        />
      </div>
    </div>

    <!-- Create Group Modal -->
    <CreateGroupModal 
      v-if="showCreateGroup"
      @close="showCreateGroup = false"
      @group-created="handleGroupCreated"
    />
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue'
import GroupCard from '@/components/GroupCard.vue'
import CreateGroupModal from '@/components/CreateGroupModal.vue'
import { useGroups } from '@/composables/useGroups'

const { groups, isLoading, error, fetchGroups, searchGroups } = useGroups()

const searchQuery = ref('')
const activeFilter = ref('all')
const showCreateGroup = ref(false)
const searchTimeout = ref(null)

const filterTabs = [
  { key: 'all', label: 'Tous les groupes' },
  { key: 'my-groups', label: 'Mes groupes' },
  { key: 'public', label: 'Public' },
  { key: 'private', label: 'Privé' }
]

const filteredGroups = computed(() => {
  let filtered = groups.value

  // Apply filter
  switch (activeFilter.value) {
    case 'my-groups':
      filtered = filtered.filter(group => group.isMember)
      break
    case 'public':
      filtered = filtered.filter(group => group.isPublic)
      break
    case 'private':
      filtered = filtered.filter(group => !group.isPublic)
      break
  }

  return filtered
})

const handleSearch = () => {
  if (searchTimeout.value) {
    clearTimeout(searchTimeout.value)
  }

  searchTimeout.value = setTimeout(async () => {
    if (searchQuery.value.trim()) {
      await searchGroups(searchQuery.value.trim())
    } else {
      await loadGroups()
    }
  }, 300)
}

const loadGroups = async () => {
  await fetchGroups()
}

const handleGroupJoined = (groupId) => {
  const group = groups.value.find(g => g.id === groupId)
  if (group) {
    group.isMember = true
    group.memberCount += 1
  }
}

const handleGroupLeft = (groupId) => {
  const group = groups.value.find(g => g.id === groupId)
  if (group) {
    group.isMember = false
    group.memberCount -= 1
  }
}

const handleGroupCreated = (newGroup) => {
  groups.value.unshift(newGroup)
  showCreateGroup.value = false
}

onMounted(() => {
  loadGroups()
})

// Watch for filter changes
watch(activeFilter, () => {
  // Could trigger additional API calls if needed
})
</script>

<style scoped>
.groups-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f0f23 0%, #1a1a2e 50%, #16213e 100%);
  color: #fff;
  padding: 24px;
}

.page-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 32px;
  max-width: 1200px;
  margin-left: auto;
  margin-right: auto;
}

.header-content h1 {
  font-size: 2.5rem;
  font-weight: 700;
  margin: 0 0 8px 0;
  background: linear-gradient(135deg, #8b5cf6, #ec4899);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
}

.page-subtitle {
  color: #999;
  font-size: 1.1rem;
  margin: 0;
}

.btn-create {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff;
  border: none;
  padding: 12px 24px;
  border-radius: 10px;
  font-weight: 600;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s ease;
}

.btn-create:hover {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
  transform: translateY(-2px);
  box-shadow: 0 8px 20px rgba(139, 92, 246, 0.3);
}

.filters-section {
  max-width: 1200px;
  margin: 0 auto 32px auto;
}

.search-bar {
  margin-bottom: 24px;
}

.search-input-wrapper {
  position: relative;
  max-width: 400px;
}

.search-input-wrapper .icon-search {
  position: absolute;
  left: 16px;
  top: 50%;
  transform: translateY(-50%);
  color: #666;
}

.search-input {
  width: 100%;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  padding: 12px 16px 12px 48px;
  color: #fff;
  font-size: 1rem;
  transition: all 0.2s ease;
}

.search-input::placeholder {
  color: #666;
}

.search-input:focus {
  outline: none;
  border-color: #8b5cf6;
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
}

.filter-tabs {
  display: flex;
  gap: 8px;
  flex-wrap: wrap;
}

.filter-tab {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #999;
  padding: 8px 16px;
  border-radius: 20px;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 0.9rem;
}

.filter-tab:hover {
  border-color: rgba(255, 255, 255, 0.4);
  color: #fff;
}

.filter-tab.active {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  border-color: transparent;
  color: #fff;
}

.groups-container {
  max-width: 1200px;
  margin: 0 auto;
}

.groups-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
  gap: 24px;
}

.loading-state,
.error-state,
.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 64px 24px;
  text-align: center;
}

.spinner {
  width: 40px;
  height: 40px;
  border: 3px solid rgba(139, 92, 246, 0.3);
  border-top: 3px solid #8b5cf6;
  border-radius: 50%;
  animation: spin 1s linear infinite;
  margin-bottom: 16px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.error-state .error-icon,
.empty-state .empty-icon {
  font-size: 3rem;
  margin-bottom: 16px;
}

.error-state h3,
.empty-state h3 {
  color: #fff;
  margin: 0 0 8px 0;
  font-size: 1.25rem;
}

.error-state p,
.empty-state p {
  color: #999;
  margin: 0 0 16px 0;
  max-width: 400px;
}

.btn {
  padding: 10px 20px;
  border-radius: 8px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn-secondary:hover {
  background: rgba(255, 255, 255, 0.15);
}

/* Icon placeholders */
.icon-plus::before { content: '+'; }
.icon-search::before { content: '🔍'; }

@media (max-width: 768px) {
  .groups-page {
    padding: 16px;
  }
  
  .page-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .groups-grid {
    grid-template-columns: 1fr;
  }
}
</style>