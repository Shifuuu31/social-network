import { createRouter, createWebHistory } from 'vue-router'
import SignUp from '../pages/SignUp.vue'
import SignIn from '../pages/SignIn.vue'
import Profile from '../pages/ProfileView.vue'
import { useAuth } from '@/composables/useAuth'

const routes = [
  {
    path: '/',
    name: 'Home',

  },
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

const auth = useAuth()

router.beforeEach(async (to, from, next) => {
  const isPublic = publicPaths.includes(to.path)

  // Always ensure user state is loaded once
  if (!auth.isAuthenticated.value) {
    await auth.fetchCurrentUser()
  }

  // If user is already authenticated, redirect away from public pages
  if (isPublic && auth.isAuthenticated.value) {
    return next('/')
  }

  // If user is accessing protected page, ensure auth is valid
  if (!isPublic && !auth.isAuthenticated.value) {
    await auth.logout()
    return next('/signin')
  }

  next()
})


export default router