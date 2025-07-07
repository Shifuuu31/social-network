<template>
  <div class="create-group-page">
    <div class="container">
      <div class="page-header">
        <h1 class="page-title">Créer un nouveau groupe</h1>
        <p class="page-subtitle">Rassemblez des personnes autour d'un intérêt commun</p>
      </div>

      <div class="create-group-form-container">
        <form @submit.prevent="handleSubmit" class="create-group-form">
          <div class="form-section">
            <h3 class="section-title">Informations générales</h3>

            <div class="form-group">
              <label for="groupName" class="form-label">Nom du groupe *</label>
              <input id="groupName" v-model="formData.name" type="text" placeholder="Entrez le nom du groupe" required
                maxlength="100" class="form-input" />
              <span class="char-count">{{ formData.name.length }}/100</span>
            </div>

            <div class="form-group">
              <label for="groupDescription" class="form-label">Description</label>
              <textarea id="groupDescription" v-model="formData.description" placeholder="Décrivez votre groupe..."
                rows="4" maxlength="500" class="form-textarea"></textarea>
              <span class="char-count">{{ formData.description.length }}/500</span>
            </div>
          </div>

          <div class="form-section">
            <h3 class="section-title">Image du groupe</h3>

            <div class="form-group">
              <div class="image-upload" @click="triggerFileInput">
                <div class="image-preview" v-if="imagePreview">
                  <img :src="imagePreview" alt="Aperçu" />
                  <button type="button" class="remove-image" @click.stop="removeImage">
                    <span class="icon">🗑️</span>
                  </button>
                </div>
                <div v-else class="upload-placeholder">
                  <span class="upload-icon">🖼️</span>
                  <p>Cliquez pour ajouter une image</p>
                  <span class="upload-hint">JPG, PNG ou GIF • Max 5MB</span>
                </div>
                <input type="file" ref="imageInput" @change="handleImageUpload" accept="image/*" class="file-input" />
              </div>
            </div>
          </div>

          <div class="form-actions">
            <router-link to="/groups" class="btn btn-secondary">
              Annuler
            </router-link>
            <button type="submit" class="btn btn-primary" :disabled="!canSubmit || isSubmitting">
              {{ isSubmitting ? 'Création en cours...' : 'Créer le groupe' }}
            </button>
          </div>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useGroupsStore } from '@/stores/groups'

const router = useRouter()
const groupsStore = useGroupsStore()

const imageInput = ref(null)
const imagePreview = ref(null)
const isSubmitting = ref(false)

const formData = reactive({
  name: '',
  description: '',
  isPublic: true,
  image: null
})

const canSubmit = computed(() => {
  return formData.name.trim().length > 5 && formData.description.trim().length > 10
})

const triggerFileInput = () => {
  imageInput.value?.click()
}

const handleImageUpload = (event) => {
  const file = event.target.files[0]
  if (file) {
    // Validate file size (5MB limit)
    if (file.size > 5 * 1024 * 1024) {
      alert('File too big, please select a file smaller than 5MB')
      resetFileInput()
      return
    }

    // Validate file type
    if (!file.type.startsWith('image/')) {
      alert('Veuillez sélectionner un fichier image valide')
      resetFileInput()
      return
    }

    formData.image = file
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

const resetFileInput = () => {
  if (imageInput.value) {
    imageInput.value.value = ''
  }
}

const removeImage = (event) => {
  event.stopPropagation()
  formData.image = null
  imagePreview.value = null
  resetFileInput()
}

const handleSubmit = async () => {
  if (!canSubmit.value) return

  isSubmitting.value = true

  try {
    const groupData = {
      name: formData.name.trim(),
      description: formData.description.trim(),
      isPublic: formData.isPublic,
      image: formData.image ? imagePreview.value : 'default-group.jpg' // Fallback image if no image is uploaded
    }

    const newGroup = await groupsStore.createGroup(groupData)

    if (newGroup) {
      router.push(`/groups/${newGroup.id}`)
    }

  } catch (error) {
    console.error('Error creating group:', error)
    alert('Erreur lors de la création du groupe. Veuillez réessayer.')
  } finally {
    isSubmitting.value = false
  }
}
</script>

<style scoped>
.create-group-page {
  padding: 40px 20px;
}

.container {
  max-width: 800px;
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

.create-group-form-container {
  background: #1a1a1a;
  border-radius: 16px;
  padding: 40px;
  border: 1px solid #333;
}

.create-group-form {
  display: flex;
  flex-direction: column;
  gap: 40px;
}

.form-section {
  display: flex;
  flex-direction: column;
  gap: 20px;
}

.section-title {
  font-size: 1.3rem;
  font-weight: 600;
  color: #fff;
  margin: 0;
  padding-bottom: 10px;
  border-bottom: 2px solid #333;
}

.form-group {
  display: flex;
  flex-direction: column;
  gap: 8px;
}

.form-label {
  font-weight: 500;
  color: #fff;
  font-size: 1rem;
}

.form-input,
.form-textarea {
  padding: 16px;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  color: #fff;
  font-size: 1rem;
  transition: border-color 0.2s ease, box-shadow 0.2s ease;
}

.form-input:focus,
.form-textarea:focus {
  outline: none;
  border-color: #8b5cf6;
  box-shadow: 0 0 0 3px rgba(139, 92, 246, 0.1);
}

.form-textarea {
  resize: vertical;
  min-height: 120px;
  font-family: inherit;
}

.char-count {
  text-align: right;
  color: #666;
  font-size: 0.85rem;
}

.image-upload {
  border: 2px dashed rgba(255, 255, 255, 0.3);
  border-radius: 12px;
  padding: 40px;
  text-align: center;
  cursor: pointer;
  transition: all 0.2s ease;
  position: relative;
}

.image-upload:hover {
  border-color: rgba(255, 255, 255, 0.5);
  background: rgba(255, 255, 255, 0.02);
}

.image-preview {
  position: relative;
  display: inline-block;
  max-width: 100%;
}

.image-preview img {
  max-width: 300px;
  max-height: 200px;
  border-radius: 8px;
  object-fit: cover;
}

.remove-image {
  position: absolute;
  top: -12px;
  right: -12px;
  background: #ef4444;
  color: white;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 1rem;
  transition: all 0.2s ease;
  z-index: 10;
}

.remove-image:hover {
  background: #dc2626;
  transform: scale(1.1);
}

.upload-placeholder {
  color: #ccc;
}

.upload-icon {
  font-size: 3rem;
  margin-bottom: 16px;
  display: block;
}

.upload-placeholder p {
  font-size: 1.1rem;
  margin: 0 0 8px 0;
  color: #fff;
}

.upload-hint {
  font-size: 0.9rem;
  color: #666;
}

.file-input {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  opacity: 0;
  cursor: pointer;
}

.privacy-options {
  display: flex;
  flex-direction: column;
  gap: 16px;
}

.privacy-option {
  display: flex;
  align-items: flex-start;
  padding: 20px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 12px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.privacy-option:hover {
  border-color: rgba(255, 255, 255, 0.4);
  background: rgba(255, 255, 255, 0.02);
}

.privacy-option:has(.privacy-radio:checked) {
  border-color: #8b5cf6;
  background: rgba(139, 92, 246, 0.1);
}

.privacy-radio {
  margin-right: 16px;
  margin-top: 4px;
  width: 18px;
  height: 18px;
  accent-color: #8b5cf6;
}

.option-content {
  display: flex;
  align-items: flex-start;
  gap: 16px;
  flex: 1;
}

.option-icon {
  font-size: 1.8rem;
  margin-top: 2px;
}

.option-text h4 {
  margin: 0 0 8px 0;
  color: #fff;
  font-size: 1.1rem;
  font-weight: 600;
}

.option-text p {
  margin: 0;
  color: #ccc;
  font-size: 0.95rem;
  line-height: 1.4;
}

.form-actions {
  display: flex;
  gap: 16px;
  justify-content: flex-end;
  padding-top: 20px;
  border-top: 1px solid #333;
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
  min-width: 140px;
}

.btn:disabled {
  opacity: 0.6;
  cursor: not-allowed;
}

.btn-secondary {
  background: rgba(255, 255, 255, 0.1);
  color: #fff;
  border: 1px solid rgba(255, 255, 255, 0.2);
}

.btn-secondary:hover:not(:disabled) {
  background: rgba(255, 255, 255, 0.15);
  transform: translateY(-1px);
}

.btn-primary {
  background: linear-gradient(135deg, #8b5cf6, #a855f7);
  color: #fff;
}

.btn-primary:hover:not(:disabled) {
  background: linear-gradient(135deg, #7c3aed, #9333ea);
  transform: translateY(-1px);
}

.icon {
  font-size: 1rem;
}

@media (max-width: 768px) {
  .create-group-form-container {
    padding: 24px;
  }

  .form-actions {
    flex-direction: column;
  }

  .btn {
    width: 100%;
  }

  .image-upload {
    padding: 24px;
  }

  .privacy-options {
    gap: 12px;
  }

  .privacy-option {
    padding: 16px;
  }
}
</style>