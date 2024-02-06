import { createRouter, createWebHistory } from 'vue-router'
import ServerView from '../views/Server.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      redirect: '/server'
    },
    {
      path: '/server',
      name: 'server',
      component: ServerView
    }
  ]
})

export default router
