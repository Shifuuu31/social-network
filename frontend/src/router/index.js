import { createRouter, createWebHistory } from 'vue-router'
import Home from '../views/Home.vue'
import Groups from '../views/Groups.vue'
import GroupView from '../views/GroupView.vue'
import CreateGroup from '../views/CreateGroup.vue'

const routes = [
  {
    path: '/',
    name: 'Home',
    component: Home
  },
  {
    path: '/groups',
    name: 'Groups',
    component: Groups
  },
  {
    path: '/groups/create',
    name: 'CreateGroup',
    component: CreateGroup
  },
  {
    path: '/groups/:id',
    name: 'GroupView',
    component: GroupView,
    props: true
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router