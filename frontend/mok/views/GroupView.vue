<template>
  <div class="group-view" v-if="currentGroup">
    <div class="group-header">
      <div class="group-header-bg">
        <img :src="currentGroup.image" :alt="currentGroup.name" />
      </div>
      <div class="group-header-content">
        <div class="container">
          <div class="group-info">
            <h1 class="group-name">{{ currentGroup.name }}</h1>
            <p class="group-description">{{ currentGroup.description }}</p>
            <div class="group-meta">
              <span class="member-count">
                <span class="icon">👥</span>
                {{ currentGroup.memberCount }} membres
              </span>
              
            </div>
          </div>
          <div class="group-actions">
            <button 
              v-if="!currentGroup.isMember"
              class="btn btn-primary"
              @click="handleJoinGroup"
              :disabled="isJoining"
            >
              <span class="icon">+</span>
              {{ isJoining ? 'Rejoindre...' : 'Rejoindre le groupe' }}
            </button>
            <button 
              v-else
              class="btn btn-secondary"
              @click="handleLeaveGroup"
              :disabled="isLeaving"
            >
              <span class="icon">✓</span>
              {{ isLeaving ? 'Quitter...' : 'Membre' }}
            </button>
            <button class="btn btn-outline" @click="toggleInviteModal">
              <span class="icon">📧</span>
              Inviter
            </button>
          </div>
        </div>
      </div>
    </div>

    <div class="group-content">
      <div class="container">
        <div class="content-tabs">
          <!-- <button 
            :class="['tab-btn', { active: activeTab === 'posts' }]"
            @click="setActiveTab('posts')"
          >
            Publications
          </button>
          <button 
            :class="['tab-btn', { active: activeTab === 'events' }]"
            @click="setActiveTab('events')"
          >
            Événements
          </button> -->
        </div>

        <div class="tab-content">
          <!-- Posts Tab -->
          <div v-if="activeTab === 'posts'" class="posts-section">
            <div class="create-post" v-if="currentGroup.isMember">
              <div class="create-post-header">
                <h3>Créer une publication</h3>
              </div>
              <form @submit.prevent="handleCreatePost" class="create-post-form">
                <input 
                  type="text" 
                  v-model="newPost.title"
                  placeholder="Titre de votre publication..."
                  class="form-input"
                  required
                />
                <textarea 
                  v-model="newPost.content"
                  placeholder="Partagez quelque chose avec le groupe..."
                  class="form-textarea"
                  rows="4"
                  required
                ></textarea>
                <div class="form-actions">
                  <button type="submit" class="btn btn-primary" :disabled="isCreatingPost">
                    {{ isCreatingPost ? 'Publication...' : 'Publier' }}
                  </button>
                </div>
              </form>
            </div>

            <div class="posts-list">
              <div v-if="isLoadingPosts" class="loading">
                <div class="spinner"></div>
                <p>Chargement des publications...</p>
              </div>
              <div v-else-if="groupPosts.length === 0" class="empty-state">
                <div class="empty-icon">📝</div>
                <h3>Aucune publication</h3>
                <p v-if="currentGroup.isMember">Soyez le premier à publier quelque chose !</p>
                <p v-else>Rejoignez le groupe pour voir et publier du contenu.</p>
              </div>
              <div v-else class="posts-grid">
                <div v-for="post in groupPosts" :key="post.id" class="post-card">
                  <div class="post-header">
                    <img :src="post.authorAvatar" :alt="post.author" class="author-avatar" />
                    <div class="post-meta">
                      <h4 class="author-name">{{ post.author }}</h4>
                      <span class="post-date">{{ formatDate(post.createdAt) }}</span>
                    </div>
                  </div>
                  <div class="post-content">
                    <h3 class="post-title">{{ post.title }}</h3>
                    <p class="post-text">{{ post.content }}</p>
                  </div>
                  <div class="post-actions">
                    
                    <button class="post-action">
                      <span class="icon">💬</span>
                      {{ post.comments }}
                    </button>
                  </div>
                </div>
              </div>
            </div>
          </div>

          <!-- Events Tab -->
          <div v-if="activeTab === 'events'" class="events-section">
            <div class="create-event" v-if="currentGroup.isMember">
              <div class="create-event-header">
                <h3>Créer un événement</h3>
              </div>
              <form @submit.prevent="handleCreateEvent" class="create-event-form">
                <input 
                  type="text" 
                  v-model="newEvent.title"
                  placeholder="Titre de l'événement..."
                  class="form-input"
                  required
                />
                <textarea 
                  v-model="newEvent.description"
                  placeholder="Description de l'événement..."
                  class="form-textarea"
                  rows="3"
                  required
                ></textarea>
                <div class="form-row">
                  <input 
                    type="datetime-local" 
                    v-model="newEvent.date"
                    class="form-input"
                    required
                  />
                  <input 
                    type="text" 
                    v-model="newEvent.location"
                    placeholder="Lieu de l'événement..."
                    class="form-input"
                    required
                  />
                </div>
                <input 
                  type="number" 
                  v-model="newEvent.maxAttendees"
                  placeholder="Nombre maximum de participants..."
                  class="form-input"
                  min="1"
                  required
                />
                <div class="form-actions">
                  <button type="submit" class="btn btn-primary" :disabled="isCreatingEvent">
                    {{ isCreatingEvent ? 'Création...' : 'Créer l\'événement' }}
                  </button>
                </div>
              </form>
            </div>

            <div class="events-list">
              <div v-if="isLoadingEvents" class="loading">
                <div class="spinner"></div>
                <p>Chargement des événements...</p>
              </div>
              <div v-else-if="groupEvents.length === 0" class="empty-state">
                <div class="empty-icon">📅</div>
                <h3>Aucun événement</h3>
                <p v-if="currentGroup.isMember">Créez le premier événement du groupe !</p>
                <p v-else>Rejoignez le groupe pour voir et créer des événements.</p>
              </div>
              <div v-else class="events-grid">
                <div v-for="event in groupEvents" :key="event.id" class="event-card">
                  <div class="event-image" v-if="event.image">
                    <img :src="event.image" :alt="event.title" />
                  </div>
                  <div class="event-content">
                    <h3 class="event-title">{{ event.title }}</h3>
                    <p class="event-description">{{ event.description }}</p>
                    <div class="event-details">
                      <div class="event-detail">
                        <span class="icon">📅</span>
                        <span>{{ formatEventDate(event.date) }}</span>
                      </div>
                      <div class="event-detail">
                        <span class="icon">📍</span>
                        <span>{{ event.location }}</span>
                      </div>
                      <div class="event-detail">
                        <span class="icon">👥</span>
                        <span>{{ event.attendees }}/{{ event.maxAttendees }} participants</span>
                      </div>
                    </div>
                    <div class="event-actions">
                      <button 
                        v-if="!event.isAttending && currentGroup.isMember"
                        class="btn btn-primary btn-sm"
                        @click="handleJoinEvent(event.id)"
                      >
                        Participer
                      </button>
                      <button 
                        v-else-if="event.isAttending"
                        class="btn btn-secondary btn-sm"
                        @click="handleLeaveEvent(event.id)"
                      >
                        Je participe
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>

    <!-- Invite Modal -->
    <div v-if="showInviteModal" class="modal-overlay" @click="toggleInviteModal">
      <div class="modal-content" @click.stop>
        <div class="modal-header">
          <h3>Inviter des personnes</h3>
          <button class="close-btn" @click="toggleInviteModal">×</button>
        </div>
        <div class="modal-body">
          <p>Partagez ce lien pour inviter des personnes à rejoindre {{ currentGroup.name }} :</p>
          <div class="invite-link">
            <input 
              type="text" 
              :value="inviteLink" 
              readonly 
              class="invite-input"
              ref="inviteLinkInput"
            />
            <!-- <button class="btn btn-primary" @click="copyInviteLink">
              Copier
            </button> -->
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, reactive } from 'vue'
import { useRoute } from 'vue-router'
import { useGroupsStore } from '@/stores/groups'

const route = useRoute()
const groupsStore = useGroupsStore()

const activeTab = ref('posts')
const isJoining = ref(false)
const isLeaving = ref(false)
const isCreatingPost = ref(false)
const isCreatingEvent = ref(false)
const isLoadingPosts = ref(false)
const isLoadingEvents = ref(false)
const showInviteModal = ref(false)
const inviteLinkInput = ref(null)

const newPost = reactive({
  title: '',
  content: ''
})

const newEvent = reactive({
  title: '',
  description: '',
  date: '',
  location: '',
  maxAttendees: 50
})

const { currentGroup, groupPosts, groupEvents } = groupsStore

const inviteLink = computed(() => {
  return `${window.location.origin}/groups/${route.params.id}?invite=true`
})

const setActiveTab = (tab) => {
  activeTab.value = tab
  if (tab === 'posts' && groupPosts.length === 0) {
    loadPosts()
  } else if (tab === 'events' && groupEvents.length === 0) {
    loadEvents()
  }
}

const loadGroup = async () => {
  const groupId = parseInt(route.params.id)
  await groupsStore.fetchGroup(groupId)
}

const loadPosts = async () => {
  isLoadingPosts.value = true
  await groupsStore.fetchGroupPosts(route.params.id)
  isLoadingPosts.value = false
}

const loadEvents = async () => {
  isLoadingEvents.value = true
  await groupsStore.fetchGroupEvents(route.params.id)
  isLoadingEvents.value = false
}

const handleJoinGroup = async () => {
  isJoining.value = true
  try {
    console.log('Joining group:', route.params.id);
    await groupsStore.requestJoinGroup(parseInt(route.params.id))
  } finally {
    isJoining.value = false
  }
}

const handleLeaveGroup = async () => {
  isLeaving.value = true
  try {
    await groupsStore.leaveGroup(parseInt(route.params.id))
  } finally {
    isLeaving.value = false
  }
}

const handleCreatePost = async () => {
  isCreatingPost.value = true
  try {
    await groupsStore.createPost(route.params.id, {
      title: newPost.title,
      content: newPost.content
    })
    newPost.title = ''
    newPost.content = ''
  } finally {
    isCreatingPost.value = false
  }
}

const handleCreateEvent = async () => {
  isCreatingEvent.value = true
  try {
    await groupsStore.createEvent(route.params.id, {
      title: newEvent.title,
      description: newEvent.description,
      date: newEvent.date,
      location: newEvent.location,
      maxAttendees: parseInt(newEvent.maxAttendees)
    })
    newEvent.title = ''
    newEvent.description = ''
    newEvent.date = ''
    newEvent.location = ''
    newEvent.maxAttendees = 50
  } finally {
    isCreatingEvent.value = false
  }
}

const handleJoinEvent = (eventId) => {
  // console.log('Joining event:', eventId)
  // Implement event joining logic
}

const handleLeaveEvent = (eventId) => {
  // console.log('Leaving event:', eventId)
  // Implement event leaving logic
}

const toggleInviteModal = () => {
  showInviteModal.value = !showInviteModal.value
}

// const copyInviteLink = () => {
//   if (inviteLinkInput.value) {
//     inviteLinkInput.value.select()
//     document.execCommand('copy')
//     // Show success message
//   }
// }

const formatDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('fr-FR', {
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

const formatEventDate = (dateString) => {
  const date = new Date(dateString)
  return date.toLocaleDateString('fr-FR', {
    weekday: 'long',
    day: 'numeric',
    month: 'long',
    year: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  })
}

onMounted(async () => {
  await loadGroup()
  await loadPosts()
})
</script>

<style scoped>
.group-view {
  min-height: 100vh;
}

.group-header {
  position: relative;
  height: 300px;
  overflow: hidden;
}

.group-header-bg {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  z-index: 1;
}

.group-header-bg img {
  width: 100%;
  height: 100%;
  object-fit: cover;
  filter: brightness(0.4);
}

.group-header-content {
  position: relative;
  z-index: 2;
  height: 100%;
  display: flex;
  align-items: end;
  padding: 40px 20px;
  background: linear-gradient(to top, rgba(0,0,0,0.8), transparent);
}

.container {
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  display: flex;
  justify-content: space-between;
  align-items: end;
  gap: 40px;
}

.group-info {
  flex: 1;
}

.group-name {
  font-size: 3rem;
  font-weight: 700;
  color: #fff;
  margin-bottom: 10px;
  text-shadow: 2px 2px 4px rgba(0,0,0,0.5);
}

.group-description {
  font-size: 1.1rem;
  color: #ccc;
  margin-bottom: 15px;
  line-height: 1.5;
}

.group-meta {
  display: flex;
  gap: 15px;
  align-items: center;
  flex-wrap: wrap;
}

.member-count {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #ccc;
}

.privacy-badge {
  padding: 4px 12px;
  border-radius: 20px;
  font-size: 0.85rem;
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

.group-actions {
  display: flex;
  gap: 12px;
}

.group-content {
  padding: 40px 20px;
}

.content-tabs {
  display: flex;
  gap: 8px;
  margin-bottom: 30px;
  border-bottom: 2px solid #333;
}

.tab-btn {
  padding: 12px 24px;
  background: none;
  border: none;
  color: #ccc;
  cursor: pointer;
  font-weight: 500;
  transition: all 0.2s ease;
  position: relative;
}

.tab-btn:hover {
  color: #fff;
}

.tab-btn.active {
  color: #8b5cf6;
}

.tab-btn.active::after {
  content: '';
  position: absolute;
  bottom: -2px;
  left: 0;
  right: 0;
  height: 2px;
  background: #8b5cf6;
}

.create-post,
.create-event {
  background: #1a1a1a;
  border-radius: 12px;
  padding: 24px;
  margin-bottom: 30px;
  border: 1px solid #333;
}

.create-post-header,
.create-event-header {
  margin-bottom: 20px;
}

.create-post-header h3,
.create-event-header h3 {
  color: #fff;
  font-size: 1.2rem;
  margin: 0;
}

.create-post-form,
.create-event-form {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.form-row {
  display: flex;
  gap: 16px;
}

.form-input,
.form-textarea {
  padding: 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: #fff;
  font-size: 1rem;
  transition: border-color 0.2s ease;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #8b5cf6;
}

.form-textarea {
  resize: vertical;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
}

.posts-grid,
.events-grid {
  display: grid;
  gap: 24px;
}

.post-card,
.event-card {
  background: #1a1a1a;
  border-radius: 12px;
  padding: 20px;
  border: 1px solid #333;
  transition: transform 0.2s ease;
}

.post-card:hover,
.event-card:hover {
  transform: translateY(-2px);
}

.post-header {
  display: flex;
  align-items: center;
  gap: 12px;
  margin-bottom: 15px;
}

.author-avatar {
  width: 40px;
  height: 40px;
  border-radius: 50%;
  object-fit: cover;
}

.post-meta {
  flex: 1;
}

.author-name {
  font-size: 1rem;
  font-weight: 600;
  color: #fff;
  margin: 0;
}

.post-date {
  font-size: 0.85rem;
  color: #666;
}

.post-content {
  margin-bottom: 15px;
}

.post-title {
  font-size: 1.1rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 8px 0;
}

.post-text {
  color: #ccc;
  line-height: 1.5;
  margin: 0;
}

.post-actions {
  display: flex;
  gap: 16px;
}

.post-action {
  display: flex;
  align-items: center;
  gap: 6px;
  background: none;
  border: none;
  color: #666;
  cursor: pointer;
  font-size: 0.9rem;
  transition: color 0.2s ease;
}

.post-action:hover {
  color: #fff;
}

.event-image {
  margin-bottom: 15px;
  border-radius: 8px;
  overflow: hidden;
}

.event-image img {
  width: 100%;
  height: 150px;
  object-fit: cover;
}

.event-title {
  font-size: 1.2rem;
  font-weight: 600;
  color: #fff;
  margin: 0 0 8px 0;
}

.event-description {
  color: #ccc;
  margin-bottom: 15px;
  line-height: 1.5;
}

.event-details {
  display: flex;
  flex-direction: column;
  gap: 8px;
  margin-bottom: 15px;
}

.event-detail {
  display: flex;
  align-items: center;
  gap: 8px;
  color: #ccc;
  font-size: 0.9rem;
}

.event-actions {
  display: flex;
  justify-content: flex-end;
}

.btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  text-decoration: none;
  display: inline-flex;
  align-items: center;
  gap: 6px;
  justify-content: center;
  font-size: 0.9rem;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-sm {
  padding: 8px 16px;
  font-size: 0.85rem;
}

.btn-primary {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
  transform: translateY(-1px);
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
}

.btn-secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.15);
}

.btn-outline {
  background: transparent;
  color: #8b5cf6;
  border: 1px solid #8b5cf6;
}

.btn-outline:hover {
  background: #8b5cf6;
  color: #fff;
}

.loading {
  text-align: center;
  padding: 60px 20px;
  color: #ccc;
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
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
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
}

.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: #1a1a1a;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
  border: 1px solid #333;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 20px;
  border-bottom: 1px solid #333;
}

.modal-header h3 {
  color: #fff;
  margin: 0;
}

.close-btn {
  background: none;
  border: none;
  color: #ccc;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 4px;
  transition: color 0.2s ease;
}

.close-btn:hover {
  color: #fff;
}

.modal-body {
  padding: 20px;
}

.modal-body p {
  color: #ccc;
  margin-bottom: 20px;
}

.invite-link {
  display: flex;
  gap: 10px;
}

.invite-input {
  flex: 1;
  padding: 12px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  color: #fff;
  font-size: 0.9rem;
}

.icon {
  font-size: 1rem;
}

@media (max-width: 768px) {
  .container {
    flex-direction: column;
    align-items: flex-start;
    gap: 20px;
  }
  
  .group-name {
    font-size: 2rem;
  }
  
  .group-actions {
    width: 100%;
    justify-content: center;
  }
  
  .form-row {
    flex-direction: column;
  }
  
  .invite-link {
    flex-direction: column;
  }
}
</style>