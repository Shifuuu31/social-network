import { createRouter, createWebHistory } from 'vue-router'
import SignUp from '../pages/SignUp.vue'
import SignIn from '../pages/SignIn.vue'
import Profile from '../pages/ProfileView.vue'
import { useAuth } from '@/composables/useAuth'

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

const auth = useAuth()

// Add guard
router.beforeEach(async (to, from, next) => {  
  // Run check only if route requires auth  
  if (to.meta.requiresAuth && !auth.isAuthenticated.value) {
    const success = await auth.fetchCurrentUser()

    if (!success) {
      return next('/signin')
    }
  }

  next()
})

export default router