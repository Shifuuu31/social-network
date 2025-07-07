<!-- views/GroupDetail.vue -->
<template>
  <div class="group-detail-page">
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement du groupe...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>Erreur de chargement</h3>
      <p>{{ error }}</p>
      <button class="btn btn-secondary" @click="loadGroupData">
        Réessayer
      </button>
    </div>

    <div v-else-if="currentGroup" class="group-content">
      <!-- Group Header -->
      <div class="group-header">
        <div class="header-background">
          <img 
            :src="currentGroup.coverImage || currentGroup.image" 
            :alt="currentGroup.name"
            class="cover-image" 
          />
          <div class="header-overlay"></div>
        </div>
        
        <div class="header-content">
          <button class="back-btn" @click="goBack">
            <i class="icon-arrow-left"></i>
            Retour
          </button>
          
          <div class="group-info">
            <div class="group-avatar">
              <img :src="currentGroup.image" :alt="currentGroup.name" />
            </div>
            
            <div class="group-details">
              <h1 class="group-name">{{ currentGroup.name }}</h1>
              <p class="group-description">{{ currentGroup.description }}</p>
              
              <div class="group-meta">
                <span class="member-count">
                  <i class="icon-users"></i>
                  {{ currentGroup.memberCount }} {{ currentGroup.memberCount === 1 ? 'membre' : 'membres' }}
                </span>
                <span :class="['privacy-badge', currentGroup.isPublic ? 'public' : 'private']">
                  {{ currentGroup.isPublic ? 'Public' : 'Privé' }}
                </span>
              </div>
            </div>
            
            <div class="group-actions">
              <button 
                v-if="!currentGroup.isMember"
                class="btn btn-primary btn-join"
                @click="handleJoinGroup"
                :disabled="isJoining"
              >
                <i class="icon-plus"></i>
                {{ isJoining ? 'Rejoindre...' : 'Rejoindre le groupe' }}
              </button>
              
              <button 
                v-else
                class="btn btn-outline btn-leave"
                @click="handleLeaveGroup"
                :disabled="isLeaving"
              >
                <i class="icon-logout"></i>
                {{ isLeaving ? 'Quitter...' : 'Quitter le groupe' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Posts Section -->
      <div class="posts-section">
        <div class="section-header">
          <h2>Publications</h2>
          <button 
            v-if="currentGroup.isMember"
            class="btn btn-secondary btn-new-post"
            @click="showCreatePost = true"
          >
            <i class="icon-plus"></i>
            Nouvelle publication
          </button>
        </div>

        <!-- Create Post Form -->
        <div v-if="showCreatePost" class="create-post-form">
          <div class="post-form">
            <textarea 
              v-model="newPostContent"
              placeholder="Partagez quelque chose avec le groupe..."
              class="post-textarea"
              rows="4"
            ></textarea>
            
            <div class="post-actions">
              <div class="post-options">
                <button class="option-btn" title="Ajouter une image">
                  <i class="icon-image"></i>
                </button>
                <button class="option-btn" title="Ajouter un lien">
                  <i class="icon-link"></i>
                </button>
              </div>
              
              <div class="post-buttons">
                <button 
                  class="btn btn-secondary"
                  @click="cancelCreatePost"
                >
                  Annuler
                </button>
                <button 
                  class="btn btn-primary"
                  @click="handleCreatePost"
                  :disabled="!newPostContent.trim() || isCreatingPost"
                >
                  {{ isCreatingPost ? 'Publication...' : 'Publier' }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Posts List -->
        <div class="posts-container">
          <div v-if="isLoadingPosts" class="loading-posts">
            <div class="spinner"></div>
            <p>Chargement des publications...</p>
          </div>

          <div v-else-if="groupPosts.length === 0" class="empty-posts">
            <div class="empty-icon">📝</div>
            <h3>Aucune publication</h3>
            <p v-if="currentGroup.isMember">
              Soyez le premier à publier dans ce groupe !
            </p>
            <p v-else>
              Rejoignez le groupe pour voir les publications.
            </p>
          </div>

          <div v-else class="posts-list">
            <PostCard 
              v-for="post in groupPosts"
              :key="post.id"
              :post="post"
              @post-liked="handlePostLiked"
              @post-commented="handlePostCommented"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useGroups } from '@/composables/useGroups'
import PostCard from '@/components/PostCard.vue'

const route = useRoute()
const router = useRouter()
const { 
  currentGroup, 
  groupPosts, 
  isLoading, 
  error, 
  fetchGroup, 
  fetchGroupPosts, 
  joinGroup, 
  leaveGroup, 
  createGroupPost 
} = useGroups()

const isJoining = ref(false)
const isLeaving = ref(false)
const isLoadingPosts = ref(false)
const showCreatePost = ref(false)
const newPostContent = ref('')
const isCreatingPost = ref(false)

const groupId = computed(() => route.params.id)

const loadGroupData = async () => {
  const id = groupId.value
  if (!id) return
  
  try {
    await fetchGroup(id)
    if (currentGroup.value?.isMember) {
      isLoadingPosts.value = true
      await fetchGroupPosts(id)
      isLoadingPosts.value = false
    }
  } catch (err) {
    console.error('Failed to load group data:', err)
    isLoadingPosts.value = false
  }
}

const goBack = () => {
  router.push('/groups')
}

const handleJoinGroup = async () => {
  isJoining.value = true
  try {
    await useGroups.joinGroup(groupId.value)
    // Reload posts after joining
    await fetchGroupPosts(groupId.value)
  } catch (err) {
    console.error('Failed to join group:', err)
  } finally {
    isJoining.value = false
  }
}

const handleLeaveGroup = async () => {
  if (!confirm('Êtes-vous sûr de vouloir quitter ce groupe ?')) return
  
  isLeaving.value = true
  try {
    await leaveGroup(groupId.value)
    // Clear posts after leaving
    groupPosts.value = []
  } catch (err) {
    console.error('Failed to leave group:', err)
  } finally {
    isLeaving.value = false
  }
}

const handleCreatePost = async () => {
  if (!newPostContent.value.trim()) return
  
  isCreatingPost.value = true
  try {
    const postData = {
      content: newPostContent.value.trim(),
      type: 'text'
    }
    
    await createGroupPost(groupId.value, postData)
    cancelCreatePost()
  } catch (err) {
    console.error('Failed to create post:', err)
  } finally {
    isCreatingPost.value = false
  }
}

const cancelCreatePost = () => {
  showCreatePost.value = false
  newPostContent.value = ''
}

const handlePostLiked = (postId) => {
  const post = groupPosts.value.find(p => p.id === postId)
  if (post) {
    post.isLiked = !post.isLiked
    post.likesCount += post.isLiked ? 1 : -1
  }
}

const handlePostCommented = (postId) => {
  // Handle comment logic here
  // console.log('Comment on post:', postId)
}

onMounted(() => {
  loadGroupData()
})
</script>

<style scoped>
.group-detail-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f0f23 0%, #1a1a2e 50%, #16213e 100%);
  color: #fff;
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 24px;
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

.group-header {
  position: relative;
  height: 300px;
  overflow: hidden;
}

.header-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.header-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(180deg, rgba(0,0,0,0.3) 0%, rgba(0,0,0,0.8) 100%);
}

.header-content {
  position: relative;
  z-index: 2;
  height: 100%;
  padding: 24px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.back-btn {
  align-self: flex-start;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #fff;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: rgba(0, 0, 0, 0.7);
}

.group-info {
  display: flex;
  align-items: flex-end;
  gap: 24px;
}

.group-avatar {
  width: 80px;
  height: 80px;
  border-radius: 16px;
  overflow: hidden;
  border: 3px solid rgba(255, 255, 255, 0.2);
}

.group-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.group-details {
  flex: 1;
}

.group-name {
  font-size: 2rem;
  font-weight: 700;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
}

.group-description {
  color: #ccc;
  margin: 0 0 12px 0;
  font-size: 1.1rem;
  line-height: 1.4;
}

.group-meta {
  display: flex;
  align-items: center;
  gap: 16px;
}

.member-count {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #ccc;
}

.privacy-badge {
  padding: 4px 12px;
  border-radius: 12px;
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

.posts-section {
  max-width: 800px;
  margin: 0 auto;
  padding: 32px 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

.create-post-form {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 24px;
}

.post-textarea {
  width: 100%;
  background: transparent;
  border: none;
  color: #fff;
  font-size: 1rem;
  line-height: 1.5;
  resize: vertical;
  margin-bottom: 16px;
}

.post-textarea::placeholder {
  color: #666;
}

.post-textarea:focus {
  outline: none;
}

.post-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.post-options {
  display: flex;
  gap: 8px;
}

.option-btn {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #999;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.option-btn:hover {
  border-color: rgba(255, 255, 255, 0.4);
  color: #fff;
}

.post-buttons {
  display: flex;
  gap: 12px;
}

.btn {
  padding: 10px 20px;
  border-radius: 8px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
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
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #fff;
}

.btn-outline:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
}

.posts-container {
  margin-top: 24px;
}

.loading-posts,
.empty-posts {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48px 24px;
  text-align: center;
}

.empty-posts .empty-icon {
  font-size: 3rem;
  margin-bottom: 16px;
}

.empty-posts h3 {
  margin: 0 0 8px 0;
  font-size: 1.25rem;
}

.empty-posts p {
  color: #999;
  margin: 0;
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Icon placeholders */
.icon-arrow-left::before { content: '←'; }
.icon-users::before { content: '👥'; }
.icon-plus::before { content: '+'; }
.icon-logout::before { content: '🚪'; }
.icon-image::before { content: '🖼️'; }
.icon-link::before { content: '🔗'; }

@media (max-width: 768px) {
  .group-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .group-actions {
    width: 100%;
  }
  
  .posts-section {
    padding: 24px 16px;
  }
  
  .section-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .post-actions {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
}
</style><!-- views/GroupDetail.vue -->
<template>
  <div class="group-detail-page">
    <div v-if="isLoading" class="loading-state">
      <div class="spinner"></div>
      <p>Chargement du groupe...</p>
    </div>

    <div v-else-if="error" class="error-state">
      <div class="error-icon">⚠️</div>
      <h3>Erreur de chargement</h3>
      <p>{{ error }}</p>
      <button class="btn btn-secondary" @click="loadGroupData">
        Réessayer
      </button>
    </div>

    <div v-else-if="currentGroup" class="group-content">
      <!-- Group Header -->
      <div class="group-header">
        <div class="header-background">
          <img 
            :src="currentGroup.coverImage || currentGroup.image" 
            :alt="currentGroup.name"
            class="cover-image" 
          />
          <div class="header-overlay"></div>
        </div>
        
        <div class="header-content">
          <button class="back-btn" @click="goBack">
            <i class="icon-arrow-left"></i>
            Retour
          </button>
          
          <div class="group-info">
            <div class="group-avatar">
              <img :src="currentGroup.image" :alt="currentGroup.name" />
            </div>
            
            <div class="group-details">
              <h1 class="group-name">{{ currentGroup.name }}</h1>
              <p class="group-description">{{ currentGroup.description }}</p>
              
              <div class="group-meta">
                <span class="member-count">
                  <i class="icon-users"></i>
                  {{ currentGroup.memberCount }} {{ currentGroup.memberCount === 1 ? 'membre' : 'membres' }}
                </span>
                <span :class="['privacy-badge', currentGroup.isPublic ? 'public' : 'private']">
                  {{ currentGroup.isPublic ? 'Public' : 'Privé' }}
                </span>
              </div>
            </div>
            
            <div class="group-actions">
              <button 
                v-if="!currentGroup.isMember"
                class="btn btn-primary btn-join"
                @click="handleJoinGroup"
                :disabled="isJoining"
              >
                <i class="icon-plus"></i>
                {{ isJoining ? 'Rejoindre...' : 'Rejoindre le groupe' }}
              </button>
              
              <button 
                v-else
                class="btn btn-outline btn-leave"
                @click="handleLeaveGroup"
                :disabled="isLeaving"
              >
                <i class="icon-logout"></i>
                {{ isLeaving ? 'Quitter...' : 'Quitter le groupe' }}
              </button>
            </div>
          </div>
        </div>
      </div>

      <!-- Posts Section -->
      <div class="posts-section">
        <div class="section-header">
          <h2>Publications</h2>
          <button 
            v-if="currentGroup.isMember"
            class="btn btn-secondary btn-new-post"
            @click="showCreatePost = true"
          >
            <i class="icon-plus"></i>
            Nouvelle publication
          </button>
        </div>

        <!-- Create Post Form -->
        <div v-if="showCreatePost" class="create-post-form">
          <div class="post-form">
            <textarea 
              v-model="newPostContent"
              placeholder="Partagez quelque chose avec le groupe..."
              class="post-textarea"
              rows="4"
            ></textarea>
            
            <div class="post-actions">
              <div class="post-options">
                <button class="option-btn" title="Ajouter une image">
                  <i class="icon-image"></i>
                </button>
                <button class="option-btn" title="Ajouter un lien">
                  <i class="icon-link"></i>
                </button>
              </div>
              
              <div class="post-buttons">
                <button 
                  class="btn btn-secondary"
                  @click="cancelCreatePost"
                >
                  Annuler
                </button>
                <button 
                  class="btn btn-primary"
                  @click="handleCreatePost"
                  :disabled="!newPostContent.trim() || isCreatingPost"
                >
                  {{ isCreatingPost ? 'Publication...' : 'Publier' }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Posts List -->
        <div class="posts-container">
          <div v-if="isLoadingPosts" class="loading-posts">
            <div class="spinner"></div>
            <p>Chargement des publications...</p>
          </div>

          <div v-else-if="groupPosts.length === 0" class="empty-posts">
            <div class="empty-icon">📝</div>
            <h3>Aucune publication</h3>
            <p v-if="currentGroup.isMember">
              Soyez le premier à publier dans ce groupe !
            </p>
            <p v-else>
              Rejoignez le groupe pour voir les publications.
            </p>
          </div>

          <div v-else class="posts-list">
            <PostCard 
              v-for="post in groupPosts"
              :key="post.id"
              :post="post"
              @post-liked="handlePostLiked"
              @post-commented="handlePostCommented"
            />
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useGroups } from '@/composables/useGroups'
import PostCard from '@/components/PostCard.vue'

const route = useRoute()
const router = useRouter()
const { 
  currentGroup, 
  groupPosts, 
  isLoading, 
  error, 
  fetchGroup, 
  fetchGroupPosts, 
  joinGroup, 
  leaveGroup, 
  createGroupPost 
} = useGroups()

const isJoining = ref(false)
const isLeaving = ref(false)
const isLoadingPosts = ref(false)
const showCreatePost = ref(false)
const newPostContent = ref('')
const isCreatingPost = ref(false)

const groupId = computed(() => route.params.id)

const loadGroupData = async () => {
  const id = groupId.value
  if (!id) return
  
  try {
    await fetchGroup(id)
    if (currentGroup.value?.isMember) {
      isLoadingPosts.value = true
      await fetchGroupPosts(id)
      isLoadingPosts.value = false
    }
  } catch (err) {
    console.error('Failed to load group data:', err)
    isLoadingPosts.value = false
  }
}

const goBack = () => {
  router.push('/groups')
}

const handleJoinGroup = async () => {
  isJoining.value = true
  try {
    await joinGroup(groupId.value)
    // Reload posts after joining
    await fetchGroupPosts(groupId.value)
  } catch (err) {
    console.error('Failed to join group:', err)
  } finally {
    isJoining.value = false
  }
}

const handleLeaveGroup = async () => {
  if (!confirm('Êtes-vous sûr de vouloir quitter ce groupe ?')) return
  
  isLeaving.value = true
  try {
    await leaveGroup(groupId.value)
    // Clear posts after leaving
    groupPosts.value = []
  } catch (err) {
    console.error('Failed to leave group:', err)
  } finally {
    isLeaving.value = false
  }
}

const handleCreatePost = async () => {
  if (!newPostContent.value.trim()) return
  
  isCreatingPost.value = true
  try {
    const postData = {
      content: newPostContent.value.trim(),
      type: 'text'
    }
    
    await createGroupPost(groupId.value, postData)
    cancelCreatePost()
  } catch (err) {
    console.error('Failed to create post:', err)
  } finally {
    isCreatingPost.value = false
  }
}

const cancelCreatePost = () => {
  showCreatePost.value = false
  newPostContent.value = ''
}

const handlePostLiked = (postId) => {
  const post = groupPosts.value.find(p => p.id === postId)
  if (post) {
    post.isLiked = !post.isLiked
    post.likesCount += post.isLiked ? 1 : -1
  }
}

const handlePostCommented = (postId) => {
  // Handle comment logic here
  // console.log('Comment on post:', postId)
}

onMounted(() => {
  loadGroupData()
})
</script>

<style scoped>
.group-detail-page {
  min-height: 100vh;
  background: linear-gradient(135deg, #0f0f23 0%, #1a1a2e 50%, #16213e 100%);
  color: #fff;
}

.loading-state,
.error-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  min-height: 100vh;
  padding: 24px;
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

.group-header {
  position: relative;
  height: 300px;
  overflow: hidden;
}

.header-background {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
}

.cover-image {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.header-overlay {
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background: linear-gradient(180deg, rgba(0,0,0,0.3) 0%, rgba(0,0,0,0.8) 100%);
}

.header-content {
  position: relative;
  z-index: 2;
  height: 100%;
  padding: 24px;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
}

.back-btn {
  align-self: flex-start;
  background: rgba(0, 0, 0, 0.5);
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #fff;
  padding: 8px 16px;
  border-radius: 8px;
  cursor: pointer;
  display: flex;
  align-items: center;
  gap: 8px;
  transition: all 0.2s ease;
}

.back-btn:hover {
  background: rgba(0, 0, 0, 0.7);
}

.group-info {
  display: flex;
  align-items: flex-end;
  gap: 24px;
}

.group-avatar {
  width: 80px;
  height: 80px;
  border-radius: 16px;
  overflow: hidden;
  border: 3px solid rgba(255, 255, 255, 0.2);
}

.group-avatar img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.group-details {
  flex: 1;
}

.group-name {
  font-size: 2rem;
  font-weight: 700;
  margin: 0 0 8px 0;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.5);
}

.group-description {
  color: #ccc;
  margin: 0 0 12px 0;
  font-size: 1.1rem;
  line-height: 1.4;
}

.group-meta {
  display: flex;
  align-items: center;
  gap: 16px;
}

.member-count {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #ccc;
}

.privacy-badge {
  padding: 4px 12px;
  border-radius: 12px;
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

.posts-section {
  max-width: 800px;
  margin: 0 auto;
  padding: 32px 24px;
}

.section-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
}

.section-header h2 {
  font-size: 1.5rem;
  font-weight: 600;
  margin: 0;
}

.create-post-form {
  background: rgba(255, 255, 255, 0.05);
  border: 1px solid rgba(255, 255, 255, 0.1);
  border-radius: 12px;
  padding: 20px;
  margin-bottom: 24px;
}

.post-textarea {
  width: 100%;
  background: transparent;
  border: none;
  color: #fff;
  font-size: 1rem;
  line-height: 1.5;
  resize: vertical;
  margin-bottom: 16px;
}

.post-textarea::placeholder {
  color: #666;
}

.post-textarea:focus {
  outline: none;
}

.post-actions {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.post-options {
  display: flex;
  gap: 8px;
}

.option-btn {
  background: transparent;
  border: 1px solid rgba(255, 255, 255, 0.2);
  color: #999;
  padding: 8px;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.option-btn:hover {
  border-color: rgba(255, 255, 255, 0.4);
  color: #fff;
}

.post-buttons {
  display: flex;
  gap: 12px;
}

.btn {
  padding: 10px 20px;
  border-radius: 8px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 8px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-primary {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
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
  border: 1px solid rgba(255, 255, 255, 0.3);
  color: #fff;
}

.btn-outline:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.1);
}

.posts-container {
  margin-top: 24px;
}

.loading-posts,
.empty-posts {
  display: flex;
  flex-direction: column;
  align-items: center;
  padding: 48px 24px;
  text-align: center;
}

.empty-posts .empty-icon {
  font-size: 3rem;
  margin-bottom: 16px;
}

.empty-posts h3 {
  margin: 0 0 8px 0;
  font-size: 1.25rem;
}

.empty-posts p {
  color: #999;
  margin: 0;
}

.posts-list {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

/* Icon placeholders */
.icon-arrow-left::before { content: '←'; }
.icon-users::before { content: '👥'; }
.icon-plus::before { content: '+'; }
.icon-logout::before { content: '🚪'; }
.icon-image::before { content: '🖼️'; }
.icon-link::before { content: '🔗'; }

@media (max-width: 768px) {
  .group-info {
    flex-direction: column;
    align-items: flex-start;
    gap: 16px;
  }
  
  .group-actions {
    width: 100%;
  }
  
  .posts-section {
    padding: 24px 16px;
  }
  
  .section-header {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
  
  .post-actions {
    flex-direction: column;
    gap: 16px;
    align-items: stretch;
  }
}
</style>