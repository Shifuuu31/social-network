import { createRouter, createWebHistory } from 'vue-router'
import Groups from '../views/Groups.vue'
import GroupView from '../views/GroupView.vue'
import CreateGroup from '../views/CreateGroup.vue'
import NotificationsView from '../views/NotificationsView.vue'
import UnseenNotifications from '../components/UnseenNotifications.vue'
// import Chat from '../components/chat/ChatWindow.vue'

const routes = [
  // {
  //   path: '/chat',
  //   name: 'Chat',
  //   component: Chat
  // },
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
    component: NotificationsView
  },
  {
    path: '/notifications/unseen',
    name: 'UnseenNotifications',
    component: UnseenNotifications
  }
]

const router = createRouter({
  history: createWebHistory(),
  routes
})

export default router