// src/router/index.js
import { createRouter, createWebHistory } from 'vue-router'
import SignUp from '../pages/SignUp.vue'
import SignIn from '../pages/SignIn.vue'
import Profile from '../pages/Profile.vue'


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
  } 
]

const router = createRouter({
  history: createWebHistory(),
  routes,
})

export default router
