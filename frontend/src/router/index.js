import { createRouter, createWebHistory } from 'vue-router'
import Groups from '../views/Groups.vue'
import GroupView from '../views/GroupView.vue'
import CreateGroup from '../views/CreateGroup.vue'
import Notifications from '../views/Notifications.vue'

const routes = [
  {
    path: '/',
    name: 'groups',
    component: Groups
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
  },
  {
    path: '/notifications',
    name: 'Notifications',
    component: Notifications
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router