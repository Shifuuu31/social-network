import { ref, useSlots } from 'vue'

const BASE_URL = import.meta.env.VITE_API_URL || 'http://localhost:3000/api'

export function useApi() {
  const request = async (method, endpoint, data = null) => {
    const url = `${BASE_URL}${endpoint}`
    const options = {
      method,
      headers: {
        'Content-Type': 'application/json'
      }
    }
    
    
    if (data) {
      options.body = JSON.stringify(data)
    }
    
    const response = await fetch(url, options)
    
    if (!response.ok) {
      throw new Error(`HTTP error! status: ${response.status}`)
    }
    
    return await response.json()
  }

  return {
    get: (endpoint) => request('GET', endpoint),
    post: (endpoint, data) => request('POST', endpoint, data),
    put: (endpoint, data) => request('PUT', endpoint, data),
    del: (endpoint) => request('DELETE', endpoint)
  }
}