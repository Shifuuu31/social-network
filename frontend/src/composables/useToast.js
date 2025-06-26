// src/composables/useToast.js
import { ref } from 'vue'

const toasts = ref([])

export function useToast() {
  const showToast = (message, type = 'info') => {
    // For now, just console log - you can implement a proper toast UI later
    console.log(`${type.toUpperCase()}: ${message}`)
    
    // Add to toasts array for potential UI display
    const toast = {
      id: Date.now(),
      message,
      type,
      timestamp: new Date()
    }
    
    toasts.value.push(toast)
    
    // Auto remove after 5 seconds
    setTimeout(() => {
      toasts.value = toasts.value.filter(t => t.id !== toast.id)
    }, 5000)
  }

  return { 
    showToast,
    toasts: toasts.value 
  }
}