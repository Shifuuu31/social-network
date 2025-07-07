<template>
  <div class="modal-overlay" @click="closeModal">
    <div class="modal-content" @click.stop>
      <div class="modal-header">
        <h2>Créer un nouveau groupe</h2>
        <button class="close-btn" @click="closeModal">
          <i class="icon-x"></i>
        </button>
      </div>

      <form @submit.prevent="handleSubmit" class="group-form">
        <div class="form-group">
          <label for="groupName">Nom du groupe *</label>
          <input
            id="groupName"
            v-model="formData.name"
            type="text"
            placeholder="Entrez le nom du groupe"
            required
            maxlength="100"
            class="form-input"
          />
          <span class="char-count">{{ formData.name.length }}/100</span>
        </div>

        <div class="form-group">
          <label for="groupDescription">Description</label>
          <textarea
            id="groupDescription"
            v-model="formData.description"
            placeholder="Décrivez votre groupe..."
            rows="4"
            maxlength="500"
            class="form-textarea"
          ></textarea>
          <span class="char-count">{{ formData.description.length }}/500</span>
        </div>

        <div class="form-group">
          <label>Image du groupe</label>
          <div class="image-upload">
            <div class="image-preview" v-if="imagePreview">
              <img :src="imagePreview" alt="Aperçu" />
              <button type="button" class="remove-image" @click="removeImage">
                <i class="icon-trash"></i>
              </button>
            </div>
            <div v-else class="upload-placeholder">
              <i class="icon-image"></i>
              <p>Cliquez pour ajouter une image</p>
            </div>
            <input
              type="file"
              ref="imageInput"
              @change="handleImageUpload"
              accept="image/*"
              class="file-input"
            />
          </div>
        </div>

        <div class="form-group">
          <label>Confidentialité</label>
          <div class="privacy-options">
            <label class="privacy-option">
              <input
                type="radio"
                v-model="formData.isPublic"
                :value="true"
                name="privacy"
              />
              <div class="option-content">
                <div class="option-icon public">🌐</div>
                <div class="option-text">
                  <h4>Public</h4>
                  <p>Tout le monde peut voir le groupe et ses publications</p>
                </div>
              </div>
            </label>

            <label class="privacy-option">
              <input
                type="radio"
                v-model="formData.isPublic"
                :value="false"
                name="privacy"
              />
              <div class="option-content">
                <div class="option-icon private">🔒</div>
                <div class="option-text">
                  <h4>Privé</h4>
                  <p>Seuls les membres peuvent voir le contenu</p>
                </div>
              </div>
            </label>
          </div>
        </div>

        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="closeModal">
            Annuler
          </button>
          <button 
            type="submit" 
            class="btn btn-primary"
            :disabled="!canSubmit || isSubmitting"
          >
            {{ isSubmitting ? 'Création...' : 'Créer le groupe' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'

const emit = defineEmits(['close', 'group-created'])

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
  return formData.name.trim().length > 0
})

const handleImageUpload = (event) => {
  const file = event.target.files[0]
  if (file) {
    formData.image = file
    const reader = new FileReader()
    reader.onload = (e) => {
      imagePreview.value = e.target.result
    }
    reader.readAsDataURL(file)
  }
}

const removeImage = () => {
  formData.image = null
  imagePreview.value = null
  if (imageInput.value) {
    imageInput.value.value = ''
  }
}

const handleSubmit = async () => {
  if (!canSubmit.value) return

  isSubmitting.value = true
  
  try {
    // Create FormData for file upload
    const submitData = new FormData()
    submitData.append('name', formData.name.trim())
    submitData.append('description', formData.description.trim())
    submitData.append('isPublic', formData.isPublic)
    
    if (formData.image) {
      submitData.append('image', formData.image)
    }

    // API call to create group
    const response = await fetch('/api/groups', {
      method: 'POST',
      body: submitData,
      headers: {
        // Don't set Content-Type, let browser set it with boundary for FormData
        // 'Authorization': `Bearer ${token}` // Add if needed
      }
    })

    if (!response.ok) {
      throw new Error('Erreur lors de la création du groupe')
    }

    const newGroup = await response.json()
    emit('group-created', newGroup)
    
  } catch (error) {
    console.error('Error creating group:', error)
    // You might want to show an error message to the user
    alert('Erreur lors de la création du groupe. Veuillez réessayer.')
  } finally {
    isSubmitting.value = false
  }
}

const closeModal = () => {
  emit('close')
}
</script>

<style scoped>
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background: rgba(0, 0, 0, 0.8);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  padding: 20px;
}

.modal-content {
  background: #1a1a1a;
  border-radius: 16px;
  width: 100%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  border: 1px solid #333;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 24px;
  border-bottom: 1px solid #333;
}

.modal-header h2 {
  margin: 0;
  color: #fff;
  font-size: 1.5rem;
  font-weight: 600;
}

.close-btn {
  background: none;
  border: none;
  color: #999;
  font-size: 1.5rem;
  cursor: pointer;
  padding: 4px;
  border-radius: 4px;
  transition: color 0.2s ease;
}

.close-btn:hover {
  color: #fff;
}

.group-form {
  padding: 24px;
}

.form-group {
  margin-bottom: 24px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  color: #fff;
  font-weight: 500;
}

.form-input,
.form-textarea {
  width: 100%;
  background: rgba(255, 255, 255, 0.1);
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 12px;
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
  min-height: 100px;
}

.char-count {
  display: block;
  margin-top: 4px;
  color: #666;
  font-size: 0.85rem;
  text-align: right;
}

.image-upload {
  position: relative;
  border: 2px dashed rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  padding: 20px;
  text-align: center;
  cursor: pointer;
  transition: border-color 0.2s ease;
}

.image-upload:hover {
  border-color: rgba(255, 255, 255, 0.4);
}

.image-preview {
  position: relative;
  display: inline-block;
}

.image-preview img {
  max-width: 200px;
  max-height: 200px;
  border-radius: 8px;
  object-fit: cover;
}

.remove-image {
  position: absolute;
  top: -8px;
  right: -8px;
  background: #ef4444;
  color: white;
  border: none;
  width: 24px;
  height: 24px;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
}

.upload-placeholder {
  color: #999;
}

.upload-placeholder i {
  font-size: 2rem;
  margin-bottom: 8px;
  display: block;
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
  gap: 12px;
}

.privacy-option {
  display: flex;
  align-items: flex-start;
  padding: 16px;
  border: 1px solid rgba(255, 255, 255, 0.2);
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
}

.privacy-option:hover {
  border-color: rgba(255, 255, 255, 0.4);
}

.privacy-option input[type="radio"] {
  margin-right: 12px;
  margin-top: 4px;
}

.option-content {
  display: flex;
  align-items: flex-start;
  gap: 12px;
}

.option-icon {
  font-size: 1.5rem;
  margin-top: 2px;
}

.option-text h4 {
  margin: 0 0 4px 0;
  color: #fff;
  font-size: 1rem;
}

.option-text p {
  margin: 0;
  color: #999;
  font-size: 0.9rem;
}

.form-actions {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
  margin-top: 32px;
}

.btn {
  padding: 12px 24px;
  border-radius: 8px;
  border: none;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
  font-size: 1rem;
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

/* Icons */
.icon-x::before { content: '×'; }
.icon-image::before { content: '🖼️'; }
.icon-trash::before { content: '🗑️'; }

@media (max-width: 640px) {
  .modal-content {
    margin: 10px;
  }
  
  .modal-header,
  .group-form {
    padding: 16px;
  }
  
  .form-actions {
    flex-direction: column;
  }
}
</style>