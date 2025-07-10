<template>
  <div class="group-card">
    <div class="group-image">
      <img :src="groupImage" :alt="group.name" @error="handleImageError" />
      
    </div>
    
    <div class="group-content">
      <h3 class="group-name">{{ group.name }}</h3>
      <p class="group-description">{{ group.description }}</p>
      
      <div class="group-stats">
        <span class="member-count">
          <span class="icon">👥</span>
          {{ group.memberCount }} {{ group.memberCount === 1 ? 'membre' : 'membres' }}
        </span>
        <span class="created-date">
          <span class="icon">📅</span>
          {{ formatDate(group.createdAt) }}
        </span>
      </div>
      
      <div class="group-actions">
        <button 
          class="btn btn-secondary btn-view"
          @click="viewGroup"
        >
          <span class="icon">👁</span>
          Voir
        </button>
        
        <!-- Not a member -->
        <button 
          v-if="!group.isMember"
          class="btn btn-primary btn-join"
          @click="handleJoin"
          :disabled="isJoining"
        >
          <span class="icon">+</span>
          {{ isJoining ? 'Rejoindre...' : 'Rejoindre' }}
        </button>
        
        <!-- Requested to join -->
        <button 
          v-else-if="group.isMember === 'requested'"
          class="btn btn-grey btn-requested"
          disabled
        >
          <span class="icon">⏳</span>
          Demande envoyée
        </button>
        
        <!-- Invited to join -->
        <button 
          v-else-if="group.isMember === 'invited'"
          class="btn btn-grey btn-invited"
          @click="handleAcceptInvite"
          :disabled="isJoining"
        >
          <span class="icon">📨</span>
          {{ isJoining ? 'Accepter...' : 'Accepter l\'invitation' }}
        </button>
        
        <!-- Full member -->
        <button 
          v-else-if="group.isMember === 'member'"
          class="btn btn-success btn-joined desactivated"
        >
          <span class="icon">✓</span>
          Membre
        </button>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGroupsStore } from '@/stores/groups'

const props = defineProps({
  group: {
    type: Object,
    required: true
  }
})

const emit = defineEmits(['group-joined', 'group-left'])

const router = useRouter()
const groupsStore = useGroupsStore()

const isJoining = ref(false)
const isLeaving = ref(false)
const imageError = ref(false)

// Computed property for group image with fallback
const groupImage = computed(() => {
  if (imageError.value) {
    return '/default-group.jpg'
  }
  
  // If image is already a full URL, use it as is
  if (props.group.image && (props.group.image.startsWith('http') || props.group.image.startsWith('/'))) {
    return props.group.image
  }
  
  // If no image or just a placeholder, use default
  if (!props.group.image || props.group.image === 'image' || props.group.image === 'default') {
    return '/default-group.jpg'
  }
  
  // Otherwise, construct API URL
  return `/api/images/${props.group.image}`
})

const handleImageError = () => {
  imageError.value = true
}

const formatDate = (dateString) => {
  if (!dateString) return ''
  
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString('fr-FR', {
      year: 'numeric',
      month: 'short',
      day: 'numeric'
    })
  } catch (error) {
    return ''
  }
}

const viewGroup = () => {
  router.push(`/groups/${props.group.id}`)
}

const handleJoin = async () => {
  isJoining.value = true
  try {
    await groupsStore.requestJoinGroup(props.group.id)
    emit('group-joined', props.group.id)
    // Show success message
    alert('Join request sent successfully!')
  } catch (error) {
    console.error('Failed to join group:', error)
    alert(`Failed to join group: ${error.message}`)
  } finally {
    isJoining.value = false
  }
}

const handleAcceptInvite = async () => {
  isJoining.value = true
  try {
    await groupsStore.acceptGroupInvite(props.group.id)
    emit('group-joined', props.group.id)
    alert('Invitation accepted successfully!')
  } catch (error) {
    console.error('Failed to accept invitation:', error)
    alert(`Failed to accept invitation: ${error.message}`)
  } finally {
    isJoining.value = false
  }
}

const handleLeave = async () => {
  isLeaving.value = true
  try {
    await groupsStore.leaveGroup(props.group.id)
    emit('group-left', props.group.id)
  } catch (error) {
    console.error('Failed to leave group:', error)
  } finally {
    isLeaving.value = false
  }
}
</script>

<style scoped>
.group-card {
  background: #1a1a1a;
  border-radius: 12px;
  overflow: hidden;
  transition: transform 0.2s ease, box-shadow 0.2s ease;
  border: 1px solid #333;
}

.group-card:hover {
  transform: translateY(-4px);
  box-shadow: 0 8px 25px rgba(0, 0, 0, 0.3);
}

.group-image {
  position: relative;
  height: 160px;
  overflow: hidden;
}

.group-image img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.group-privacy {
  position: absolute;
  top: 12px;
  right: 12px;
}

.privacy-badge {
  padding: 4px 8px;
  border-radius: 12px;
  font-size: 0.75rem;
  font-weight: 500;
  text-transform: uppercase;
  letter-spacing: 0.5px;
}

.privacy-badge.public {
  background: rgba(34, 197, 94, 0.2);
  color: #22c55e;
  border: 1px solid rgba(34, 197, 94, 0.3);
}

.privacy-badge.private {
  background: rgba(239, 68, 68, 0.2);
  color: #ef4444;
  border: 1px solid rgba(239, 68, 68, 0.3);
}

.group-content {
  padding: 16px;
}

.group-name {
  font-size: 1.1rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 8px 0;
  line-height: 1.3;
}

.group-description {
  color: #999;
  font-size: 0.9rem;
  line-height: 1.4;
  margin: 0 0 12px 0;
  display: -webkit-box;
  line-clamp: 2;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
  overflow: hidden;
}

.group-stats {
  margin-bottom: 16px;
  display: flex;
  flex-direction: column;
  gap: 4px;
}

.member-count,
.created-date {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #666;
  font-size: 0.85rem;
}

.group-actions {
  display: flex;
  gap: 8px;
}

.btn {
  padding: 8px 16px;
  border-radius: 8px;
  border: none;
  font-size: 0.85rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 6px;
  flex: 1;
  justify-content: center;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn-secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.15);
}

.btn-primary {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
  transform: translateY(-1px);
}

.btn-grey {
  background: rgba(107, 114, 128, 0.3);
  color: #9ca3af;
  border: 1px solid rgba(107, 114, 128, 0.4);
}

.btn-grey:hover:not(:disabled) {
  background: rgba(107, 114, 128, 0.4);
  color: #d1d5db;
}

.btn-success {
  background: linear-gradient(135deg, #22c55e, #16a34a);
  color: #fff;
}

.btn-success:hover:not(:disabled) {
  background: linear-gradient(135deg, #16a34a, #15803d);
}

.btn-joined {
  background: linear-gradient(135deg, #22c55e, #16a34a);
}

.btn-joined:hover:not(:disabled) {
  background: linear-gradient(135deg, #16a34a, #15803d);
}

.icon {
  font-size: 0.9rem;
}


.desactivated {
  pointer-events: none;
  opacity: 0.5;
  cursor: default;
}
</style>