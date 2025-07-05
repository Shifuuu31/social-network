// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import SignUp from '../pages/SignUp.vue'
import SignIn from '../pages/SignIn.vue'
import Profile from '../pages/ProfileView.vue'
import { useAuth } from '../composables/useAuth'

const routes = [
  {
    path: '/signup',
    name: 'Signup',
    component: SignUp,
  },
  {
    path: '/signin',
    name: 'Signin',
    component: SignIn,
  },
  {
    path: '/profile',
    name: 'Profile',
    component: Profile,
    meta: { requiresAuth: true }
  } 
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

// Add guard
router.beforeEach(async (to, from, next) => {
  console.log("beffore each")
  const { isAuthenticated, fetchCurrentUser } = useAuth()

  // Run check only if route requires auth
  if (to.meta.requiresAuth && !isAuthenticated.value) {
    const success = await fetchCurrentUser()

    if (!success) {
      return next('/signin')
    }
  }

  next()
})

export default router