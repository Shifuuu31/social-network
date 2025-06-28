// router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import PostDetailsView from '@/views/postDetailsView.vue'
import PostsView from '@/views/postsView.vue'
import Signin from '@/views/signin.vue'
import Signup from '@/views/signup.vue'


// Import other views as needed
// import 
const routes = [
  {
    path: '/',
    name: 'Posts',
    component: PostsView
  },
   {
    path: '/posts',
    name: 'posts',
    component: PostsView
  },
  {
    path: '/post/:id',
    name: 'PostDetail',
    component: PostDetailsView,
    props: true
  },
    
  { path: '/signin', name: 'Signin', component: Signin },
  { path: '/signup', name: 'Signup', component: Signup },

  // Add other routes as needed
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router