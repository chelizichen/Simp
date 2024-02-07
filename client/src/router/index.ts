import { createRouter, createWebHistory } from 'vue-router'
import ServerView from '@/views/Server.vue'
import Login from '@/views/Login.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/login'
    },
    {
      path: '/server',
      name: 'server',
      component: ServerView
    },
    {
      path: '/login',
      name: 'login',
      component: Login
    }
  ]
})

export default router
