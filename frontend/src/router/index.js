// router/index.js
import { createRouter, createWebHistory } from 'vue-router'
// import PostDetailsView from '@/views/postDetailsView.vue'
// import PostsView from '@/views/postsView.vue'
import Signin from '@/views/signin.vue'
import Signup from '@/views/signup.vue'
import Home from '@/views/Home.vue'

// Import other views as needed
// import 
const routes = [
  // {
  //   path: '/',
  //   name: 'Posts',
  //   component: PostsView
  // },
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
    

  {path:'/',name:'home',component:Home},
  { path: '/signin', name: 'Signin', component: Signin },
  { path: '/signup', name: 'Signup', component: Signup },
  

  // Add other routes as needed
]

const router = createRouter({
  history: createWebHistory(),
  routes
})



// router.beforeEach((to, from, next) => {
//   const isAuthenticated = !!localStorage.getItem('token') // adjust based on your auth method
//   if (to.meta.requiresAuth && !isAuthenticated) {
//     next('/login')
//   } else {
//     next()
//   }
// })



export default router