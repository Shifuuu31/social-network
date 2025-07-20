// router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import Signin from '@/views/signin.vue'
import Signup from '@/views/signup.vue'
import Home from '@/views/Home.vue'
import { useAuth } from '@/composables/useAuth'
import Profile from '@/views/Profile.vue'
import Chat from '@/views/Chat.vue'
 
const routes = [
  { path:'/',name:'home',component:Home},
  { path: '/signin', name: 'Signin', component: Signin },
  { path: '/signup', name: 'Signup', component: Signup },                                                                                              
  { path: '/profile/:id?', name: 'Profile', component: Profile},
  { path: '/chat', name: 'Chat', component: Chat },
  {
    path: '/discover-friend',
    name: 'DiscoverFriend',
    component: () => import('@/views/DiscoverFriend.vue')
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const publicRoutes = ['/signin', '/signup']

// Add guard with better error handling
router.beforeEach(async (to, from, next) => {  
  const hasToken = document.cookie.includes('session_token');
  const auth = useAuth()
  const isPublicRoute = publicRoutes.includes(to.path)
  
  //  if(!hasToken){
  //   auth.clearAuthState()
  //   return next('/signin')
  //  }
  // If going to public route and already authenticated, redirect to home
  if (isPublicRoute && auth.isAuthenticated.value) {
    return next('/')
  }        
  
  // If going to protected route, check authentication
  if (!isPublicRoute && hasToken) {
    try {
      // Check if already authenticated or try to fetch current user
      const isAuth = auth.isAuthenticated.value || await auth.fetchCurrentUser()
      
      if (!isAuth) {
        // Clear auth state and redirect to signin
        auth.clearAuthState()
        return next('/signin')
      }
    } catch (error) {
      console.warn('Auth check failed:', error)
      // Clear auth state and redirect to signin
      auth.clearAuthState()
      return next('/signin')
    }
} 
  
  next()
})

export default router