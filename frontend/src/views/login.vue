 <template>
  <div class="login-container">
    <h2>Login</h2>
    <form @submit.prevent="handleLogin">
      <label>Email
        <input type="email" v-model="form.email" required />
      </label>
      <label>Password
        <input type="password" v-model="form.password" required />
      </label>
      <button type="submit" :disabled="loading">
        {{ loading ? 'Logging in...' : 'Login' }}
      </button>
      <p class="error" v-if="error">{{ error }}</p>
    </form>
  </div>
</template>

 <script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = reactive({
  email: '',
  password: '',
})

const loading = ref(false)
const error = ref('')

async function handleLogin() {
  error.value = ''
  loading.value = true

  try {

    const res = await fetch('http://localhost:8080/auth/signin', {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      // credentials: 'include',
      body: JSON.stringify({
        email: form.email,
        password: form.password,
      }),
    })

    if (!res.ok) {
      const msg = await res.text()
      throw new Error(msg || 'Login failed.')
    }

    alert('Login successful!')
    router.push('/')
  } catch (err) {
    error.value = err.message
  } finally {
    loading.value = false
  }
}
</script>


<style scoped>
.login-container {
  max-width: 400px;
  margin: 100px auto;
  padding: 2rem;
  background: #f9f9f9;
  border-radius: 10px;
  box-shadow: 0 0 10px #ccc;
}
.input-group {
  margin-bottom: 1rem;
}
label {
  display: block;
  margin-bottom: 0.5rem;
}
input {
  width: 100%;
  padding: 0.5rem;
  border: 1px solid #aaa;
  border-radius: 5px;
}
button {
  width: 100%;
  padding: 0.75rem;
  background: #007bff;
  color: white;
  border: none;
  border-radius: 5px;
  font-weight: bold;
}
button:disabled {
  background: #999;
}
.error {
  color: red;
  margin-top: 1rem;
}
</style>
