<template>
  <div class="login-container">
    <div class="left-panel">
      <h1>Welcome Back</h1>
      <p>We're happy to see you again!</p>
      <p>Don't have an account? <router-link to="/signup" class="signup-link">Sign Up</router-link></p>
    </div>

    <div class="form-wrapper">
      <div class="form-panel">
        <h2>Sign In</h2>
        <form @submit.prevent="onSubmit" id="signinForm">
          <div class="field">
            <label>Nickname or Email</label>
            <input name="identifier" v-model.trim="form.identifier" type="text" />
            <span class="error">{{ errors.identifier }}</span>
          </div>

          <div class="field">
            <label>Password</label>
            <input name="password" v-model="form.password" type="password" />
            <span class="error">{{ errors.password }}</span>
          </div>

          <div class="actions">
            <button type="submit" :disabled="loading">
              {{ loading ? 'Signing in...' : 'Sign In' }}
            </button>
          </div>

          <p class="error" v-if="generalError">{{ generalError }}</p>
        </form>
      </div>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const generalError = ref('')

const form = reactive({
  identifier: '',
  password: ''
})

const errors = reactive({
  identifier: '',
  password: ''
})

const onSubmit = async () => {
  generalError.value = ''
  errors.identifier = ''
  errors.password = ''

  const validation = validate(form)
  if (Object.keys(validation).length) {
    Object.assign(errors, validation)
    return
  }

  const isEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(form.identifier)
  const payload = {
    nickname: isEmail ? '' : form.identifier,
    email: isEmail ? form.identifier : '',
    password: form.password
  }

  loading.value = true
  try {
    const res = await fetch('http://localhost:8080/auth/signin', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      credentials:'include',

      body: JSON.stringify(payload)
    })

    if (!res.ok) {
      const msg = await res.text()
      throw new Error(msg || 'Signin failed')
    }

    const user = await res.json()
    // saveUser(user)
    alert('Login successful!')
    router.push('/')
  } catch (err) {
    generalError.value = err.message
  } finally {
    loading.value = false
  }
}

function saveUser(user) {
  localStorage.setItem('nickname', user.nickname)
  localStorage.setItem('email', user.email)
  localStorage.setItem('profile_img', user.profile_img || '')
  localStorage.setItem('token', user.token)
}

function validate({ identifier, password }) {
  const errors = {}

  if (!identifier) {
    errors.identifier = 'Email or nickname is required'
  } else {
    const isEmail = /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(identifier)
    const isnickname = /^[A-Za-z0-9_]{3,20}$/.test(identifier)
    if (!isEmail && !isnickname) {
      errors.identifier = 'Enter a valid email or nickname'
    }
  }

  const passRegex = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_\-+=]).{8,}$/
  if (!passRegex.test(password)) {
    errors.password = 'Password must be 8+ chars, include upper, digit & special'
  }

  return errors
}
</script>

<style scoped>
.login-container {
  display: flex;
  height: 100vh;
  font-family: sans-serif;
}
.left-panel {
  flex: 1;
  background: #1e293b;
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  padding: 2rem;
}
.signup-link {
  margin-top: 1rem;
  color: #cbd5e1;
  text-decoration: underline;
}
.form-wrapper {
  flex: 2;
  display: flex;
  justify-content: center;
  align-items: center;
  background: white;
  padding: 3rem;
}
.form-panel {
  width: 100%;
  max-width: 400px;
  background: white;
  padding: 2rem;
  /* border-radius: 8px;  */
  /* box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1); */
}
.form-panel h2 {
  color: black;
  margin-bottom: 1rem;
}
form {
  display: flex;
  flex-direction: column;
  gap: 5px;
}
.field {
  color: black;
  display: flex;
  flex-direction: column;
}
input {
  padding: 0.5rem;
  border-radius: 4px;
  border: 1px solid #ccc;
}
.error {
  color: red;
  font-size: 0.8rem;
  height: 1rem;
}
.actions {
  margin-top: 1rem;
}
button {
  padding: 0.75rem;
  background: #1e293b;
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: bold;
}
button:disabled {
  background: #9ca3af;
}
</style>
