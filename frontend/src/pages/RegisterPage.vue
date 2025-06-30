<!-- <template>
  <div class="register-container">
    <div class="welcome-section">
      <h1>Welcome</h1>
      <p>To the zone of happiness.</p>
    </div>

    <div class="form-section">
      <h2>Welcome.</h2>
      
      <form @submit.prevent="handleSubmit">
        <div class="form-grid">
          <input v-model="form.firstName" placeholder="Firstname" required />
          <input v-model="form.lastName" placeholder="Lastname" required />
        </div>
        <input v-model="form.email" placeholder="Email" type="email" required />
        <input v-model="form.password" placeholder="Password" type="password" required />
        <input v-model="form.dateOfBirth" type="date" placeholder="Date of Birth" required />
        <input type="file" accept="image/*" @change="handleFileUpload" />
        <input v-model="form.nickname" placeholder="Nickname (Optional)" />
        <textarea v-model="form.aboutMe" placeholder="About Me (Optional)"></textarea>

        <div class="form-actions">
          <button type="submit">Register</button>
        </div>
      </form>

    </div>

  </div>
</template>

<script setup>
import { reactive, ref } from 'vue'
import { useRouter } from 'vue-router'

const router = useRouter()

const form = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  dateOfBirth: '',
  nickname: '',
  aboutMe: ''
})

const avatarFile = ref(null)

function handleFileUpload(event) {
  avatarFile.value = event.target.files[0]
}

async function handleSubmit() {
  const formData = new FormData()
  for (const [key, value] of Object.entries(form)) {
    formData.append(key, value)
  }
  if (avatarFile.value) {
    formData.append('avatar', avatarFile.value)
  }

  try {
    const res = await fetch('/register', {
      method: 'POST',
      body: formData,
      credentials: 'include', // send/receive cookies
    })

    if (!res.ok) {
      const errorData = await res.json()
      throw new Error(errorData.message || 'Registration failed')
    }

    // Success
    router.push('/login')
  } catch (err) {
    alert(err.message)
  }
}
</script>

<style scoped>
/* Reset & base */
* {
  margin: 0; padding: 0; box-sizing: border-box;
  font-family: 'Segoe UI', Tahoma, Geneva, Verdana, sans-serif;
  color: #222;
}

body, html, #app {
  height: 100%;
  background: #fafafa;
}

/* Container */
.register-container {
  max-width: 900px;
  margin: 2rem auto;
  display: flex;
  border-radius: 12px;
  overflow: hidden;
  box-shadow: 0 10px 30px rgb(0 0 0 / 0.1);
  height: 90vh;
}

/* Left panel */
.welcome-section {
  flex: 1;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: white;
  padding: 3rem 2.5rem;
  display: flex;
  flex-direction: column;
  justify-content: center;
  align-items: center;
  text-align: center;
  user-select: none;
}

.welcome-section h1 {
  font-weight: 900;
  font-size: 3rem;
  margin-bottom: 0.5rem;
  letter-spacing: 0.12em;
  text-transform: uppercase;
  line-height: 1.1;
}

.welcome-section p {
  font-weight: 300;
  font-size: 1.3rem;
  opacity: 0.85;
  max-width: 260px;
  line-height: 1.5;
  font-style: italic;
}

/* Right panel */
.form-section {
  flex: 1;
  background: white;
  padding: 3.5rem ;
  display: flex;
  flex-direction: column;
  justify-content: center;
}

/* Heading */
.form-section h2 {
  font-weight: 700;
  font-size: 2rem;
  margin-bottom: 2rem;
  color: #5a4def;
  user-select: none;
}

/* Form */
form {
  width: 100%;
}

/* Form grid for first/last name */
.form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 1.5rem;
  margin-bottom: 1.5rem;
  border-radius: 10px;
    width: 100%;
}

/* Inputs and textarea styling */

input[type="text"],
input[type="email"],
input[type="password"],
input[type="date"],
textarea,
input[type="file"] {
  width: 100%;
  padding: 10px 18px;
  font-size: 1.1rem;
  border-radius: 10px;
  border: 2px solid #e0e0e0;
  font-weight: 400;
  transition: border-color 0.25s ease, box-shadow 0.25s ease;
  outline-offset: 2px;
  outline-color: transparent;
  font-family: inherit;
  box-shadow: inset 0 1px 3px rgb(0 0 0 / 0.06);
}

/* File input smaller padding */
input[type="file"] {
  padding: 10px 12px;
  font-size: 1rem;
}

/* Textarea height */
textarea {
  min-height: 100px;
  resize: vertical;
}

/* Placeholder */
input::placeholder,
textarea::placeholder {
  color: #999;
  font-weight: 300;
}

/* Focus style */
input:focus,
textarea:focus,
input[type="file"]:focus {
  border-color: #5a4def;
  box-shadow: 0 0 8px rgb(90 77 239 / 0.3);
  outline-color: #5a4def;
}

/* Button container */
.form-actions {
  margin-top: 2.5rem;
  text-align: right;
}

/* Submit button */
button {
  background-color: #5a4def;
  color: white;
  font-weight: 700;
  font-size: 1.15rem;
  padding: 0.9rem 3rem;
  border-radius: 50px;
  border: none;
  cursor: pointer;
  box-shadow: 0 6px 20px rgb(90 77 239 / 0.3);
  transition: background-color 0.3s ease, box-shadow 0.3s ease;
  user-select: none;
}

button:hover,
button:focus {
  background-color: #4236c9;
  box-shadow: 0 8px 28px rgb(66 54 201 / 0.6);
  outline: none;
}

/* Responsive stack for smaller screens */
@media (max-width: 700px) {
  .register-container {
    flex-direction: column;
    height: auto;
    max-width: 100vw;
    border-radius: 0;
    box-shadow: none;
  }

  .welcome-section {
    padding: 2.5rem 1.5rem;
    font-size: 1rem;
  }

  .welcome-section h1 {
    font-size: 2.4rem;
  }

  .form-section {
    padding: 2.5rem 1.5rem;
  }

  .form-grid {
    grid-template-columns: 1fr;
  }

  .form-actions {
    text-align: center;
  }

  button {
    width: 100%;
    padding: 1.2rem;
  }
}

</style> -->
