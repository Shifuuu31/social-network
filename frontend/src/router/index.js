// router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import PostDetailsView from '@/views/postDetailsView.vue'
import PostsView from '@/views/postsView.vue'
// Import other views as needed

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
  }
  // Add other routes as needed
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router