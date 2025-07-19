// router/index.js
import { createRouter, createWebHistory } from 'vue-router'
// import PostDetailsView from '@/views/postDetailsView.vue'
// import PostsView from '@/views/postsView.vue'
import Signin from '@/views/signin.vue'
import Signup from '@/views/signup.vue'
import Home from '@/views/Home.vue'
import { useAuth } from '@/composables/useAuth'
import Profile from '@/views/Profile.vue'
import Chat from '@/views/Chat.vue'


// Import other views as needed
// import 
const routes = [

  //  {
  //   path: '/posts',
  //   name: 'posts',
  //   component: PostsView
  // },
  // {
  //   path: '/post/:id',
  //   name: 'PostDetail',
  //   component: PostDetailsView,
  //   props: true
  // },
    

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

  // Add other routes as needed
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

const publicr = ['/signin', '/signup']
// Add guard
router.beforeEach(async (to, from, next) => {  
  const auth = useAuth()
  const pp = publicr.includes(to.path)
  if (pp && auth.isAuthenticated) {
    return next('/')

  }        
  // Run check only if route requires auth 
  if (!pp){                                                                                        
    const success = auth.isAuthenticated.value || await auth.fetchCurrentUser()
    if (!success) {
      await auth.logout()
      return next('/signin')
    }
  } 
  next()
})
                                                         


export default router