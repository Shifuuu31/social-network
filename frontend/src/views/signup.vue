<template>
  <div class="signup-container">
    <div class="left-panel">
      <div class="left-panel-content">
        <h1>Welcome</h1>
        <p>Join our social network and stay connected.</p>
        <div class="signin-link-wrapper">
          <p>If you already have an account, <router-link to="/signin" class="signin-link">Sign In</router-link></p>
        </div>
      </div>
    </div>

    <div class="form-panel">
      <h2>Create Account</h2>
      <form @submit.prevent="onSubmit" id="signupForm" enctype="multipart/form-data">
        <div class="row">
          <div class="field">
            <label>First Name</label>
            <input name="first_name" v-model.trim="form.first_name" type="text" />
            <span class="error" data-for="first_name">{{ errors.first_name }}</span>
          </div>
          <div class="field">
            <label>Last Name</label>
            <input name="last_name" v-model.trim="form.last_name" type="text" />
            <span class="error" data-for="last_name">{{ errors.last_name }}</span>
          </div>
        </div>

        <div class="field">
          <label>Email</label>
          <input name="email" v-model.trim="form.email" type="email" />
          <span class="error" data-for="email">{{ errors.email }}</span>
        </div>

        <div class="field">
          <label>Username</label>
          <input name="username" v-model.trim="form.username" type="text" />
          <span class="error" data-for="username">{{ errors.username }}</span>
        </div>

        <div class="field">
          <label>Password</label>
          <input name="password" v-model="form.password" type="password" />
          <span class="error" data-for="password">{{ errors.password }}</span>
        </div>

        <div class="field">
          <label>Confirm Password</label>
          <input name="repeated_password" v-model="form.repeated_password" type="password" />
          <span class="error" data-for="repeated_password">{{ errors.repeated_password }}</span>
        </div>

        <div class="row">
          <div class="field">
            <label>Date of Birth</label>
            <input name="birth_date" v-model="form.birth_date" type="date" />
            <span class="error" data-for="birth_date">{{ errors.birth_date }}</span>
          </div>

          <div class="field">
            <label>Gender</label>
            <select name="gender" v-model="form.gender">
              <option disabled value="">Select Gender</option>
              <option value="male">Male</option>
              <option value="female">Female</option>
            </select>
            <span class="error" data-for="gender">{{ errors.gender }}</span>
          </div>
        </div>

        <div class="field">
          <label>Nickname (optional)</label>
          <input name="nickname" v-model.trim="form.nickname" type="text" />
        </div>

        <div class="field">
          <label>Avatar (optional)</label>
          <input name="avatar_file" type="file" @change="handleFileUpload" />
        </div>

        <div class="field">
          <label>About Me (optional)</label>
          <textarea name="about_me" v-model.trim="form.about_me"></textarea>
        </div>

        <div class="actions">
          <button type="submit" :disabled="loading">
            {{ loading ? 'Creating...' : 'Register' }}
          </button>
        </div>
        <p class="error" v-if="generalError">{{ generalError }}</p>
      </form>
    </div>
  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()
const loading = ref(false)
const generalError = ref('')
const avatarFile = ref(null)

const form = reactive({
  first_name: '',
  last_name: '',
  username: '',
  email: '',
  password: '',
  repeated_password: '',
  birth_date: '',
  gender: '',
  nickname: '',
  about_me: '',
})

const errors = reactive({})

const handleFileUpload = (e) => {
  avatarFile.value = e.target.files[0]
}

const onSubmit = async () => {
  generalError.value = ''
  Object.keys(errors).forEach(k => errors[k] = '')

  const validation = validate(form)
  Object.assign(errors, validation)
  if (Object.keys(validation).length) return

  const formData = new FormData()
  
  formData.append('user', JSON.stringify(form))
  
  if (avatarFile.value) {
    formData.append('avatar_file', avatarFile.value)
  }

  // for (const key in form) {
    // if (form[key]) formData.append(key, form[key])
  // }
  // if (avatarFile.value) {
    // formData.append('avatar_file', avatarFile.value)
  // }

  loading.value = true
  try {
    const res = await fetch('http://localhost:8080/auth/signup', {
      method: 'POST',
      body: formData
    })

    if (!res.ok) {
      const msg = await res.text()
      throw new Error(msg || 'Signup failed')
    }

    alert('Account created. Please login.')
    router.push('/signin')
  } catch (err) {
    generalError.value = err.message
  } finally {
    loading.value = false
  }
}

function validate(data) {
  const err = {}
  const nameRegex = /^[A-Za-z]{3,}$/
  const userRegex = /^[A-Za-z0-9_]{3,20}$/
  const emailRegex = /^[^\s@]+@[^\s@]+\.[^\s@]+$/
  const passRegex = /^(?=.*[A-Z])(?=.*\d)(?=.*[!@#$%^&*()_\-+=]).{8,}$/

  if (!nameRegex.test(data.first_name)) err.first_name = 'First name must be 3+ letters'
  if (!nameRegex.test(data.last_name)) err.last_name = 'Last name must be 3+ letters'
  if (!userRegex.test(data.username)) err.username = 'Username must be 3-20 alphanumeric/underscore'
  if (!emailRegex.test(data.email)) err.email = 'Invalid email'

  const birthDate = new Date(data.birth_date)
  const age = new Date().getFullYear() - birthDate.getFullYear()
  if (!data.birth_date || isNaN(birthDate)) err.birth_date = 'Invalid birth date'
  else if (age < 13 || (age === 13 && Date.now() < birthDate.setFullYear(birthDate.getFullYear() + 13)))
    err.birth_date = 'You must be at least 13'

  if (!['male', 'female'].includes(data.gender)) err.gender = 'Select male or female'
  if (!passRegex.test(data.password)) err.password = 'Password must be 8+ chars, include upper, digit & special'
  if (data.password !== data.repeated_password) err.repeated_password = 'Passwords do not match'

  return err
}
</script>

<style scoped>
.signup-container {
  display: flex;
  height: 100vh;
  font-family: sans-serif;
}
.left-panel {
  flex: 1;
  background: #4f46e5;
  color: white;
  display: flex;
  flex-direction: column;
  justify-content: space-between;
  padding: 2rem;
}
.left-panel-content {
  display: flex;
  flex-direction: column;
  align-items: center;
  flex-grow: 1;
}
.signin-link-wrapper {
  text-align: center;
  margin-top: auto;
}
.signin-link {
  color: #cbd5e1;
  text-decoration: underline;
}
.form-panel h2 {
    color: black;
}
.form-panel {
  flex: 2;
  padding: 3rem;
  background: #f3f4f6;
  overflow-y: auto;
}
form {
  display: flex;
  flex-direction: column;
  /* gap: 1rem; */
  gap: 1px;

}
.row {
  display: flex;
  /* gap: 1rem; */
  gap: 5px;

}
.field {
    color: black;
  flex: 1;
  display: flex;
  flex-direction: column;
  /* gap: 1px; */

}
input, select, textarea {
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
  background: #4f46e5;
  color: white;
  border: none;
  border-radius: 6px;
  font-weight: bold;
}
button:disabled {
  background: #9ca3af;
}
</style>
