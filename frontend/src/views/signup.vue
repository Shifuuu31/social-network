 


<template>
  <div class="signup-container">
    <h2>Sign Up</h2>

    <form @submit.prevent="handleSubmit" class="form-grid">
      <label>Email
        <input type="email" v-model="form.email" required />
      </label>

      <label>Password
        <input type="password" v-model="form.password" required minlength="8" />
      </label>

      <label>First Name
        <input type="text" v-model="form.first_name" required />
      </label>

      <label>Last Name
        <input type="text" v-model="form.last_name" required />
      </label>

      <!-- <label>Date of Birth
        <input type="date" v-model="form.date_of_birth" required />
      </label> -->

      <label>Nickname (optional)
        <input type="text" v-model="form.nickname" />
      </label>

      <label>Avatar URL (optional)
        <input type="url" v-model="form.avatar_url" />
      </label>

      <label>About Me (optional)
        <textarea v-model="form.about_me" maxlength="500"></textarea>
      </label>

      <label>
        <input type="checkbox" v-model="form.is_public" />
        Make profile public
      </label>

      <button type="submit" :disabled="loading">
        {{ loading ? 'Creating...' : 'Sign Up' }}
      </button>

      <p class="error" v-if="error">{{ error }}</p>
    </form>
  </div>
</template>

<script setup>

import { useRouter } from 'vue-router'  // <-- import here
const router = useRouter()  // <-- create router instance

import { reactive, ref } from 'vue'

const form = reactive({
  email: '',
  password: '',
  first_name: '',
  last_name: '',
  date_of_birth: '2004-11-02T00:00:00Z',
  avatar_url: '',
  nickname: '',
  about_me: '',
  is_public: true,
})

const loading = ref(false)
const error = ref('')

async function handleSubmit() {
  error.value = ''
  loading.value = true

  try {
    const res = await fetch('http://localhost:8080/auth/signup', {
      method: 'POST',
      headers: {
        'Content-Type': 'application/json',
      },
      body: JSON.stringify(form),
    })

    if (!res.ok) {
      const msg = await res.text()
      throw new Error(msg || 'Failed to sign up.')
    }
    console.log(res,"reees");
     alert('Account created successfully!')

    router.push('/login') 

   } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}
</script>

<style scoped>
.signup-container {
  max-width: 600px;
  margin: auto;
  padding: 2rem;
  background: #fafafa;
  border-radius: 12px;
  box-shadow: 0 0 10px #ccc;
}
.form-grid {
  display: grid;
  gap: 1rem;
}
input, textarea {
  width: 100%;
  padding: 0.5rem;
}
button {
  padding: 0.75rem;
  background: #007bff;
  color: white;
  border: none;
  cursor: pointer;
}
button:disabled {
  background: #999;
}
.error {
  color: red;
}
</style>
