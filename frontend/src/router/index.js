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
    path: '/profile/:id?',
    name: 'Profile',
    component: Profile,
  } 
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

const publicPaths = ['/signin', '/signup']

// Add guard
router.beforeEach(async (to, from, next) => {  
  const auth = useAuth()
  const isPublic = publicPaths.includes(to.path)

  // Run check only for unpublic paths
  if (!isPublic) {
    const success = auth.isAuthenticated.value || await auth.fetchCurrentUser()

    if (!success) {
      await auth.logout()
      return next('/signin')
    }
  }

  next()
})

export default router