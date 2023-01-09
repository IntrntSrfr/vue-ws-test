import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ChatView from '../views/ChatView.vue'

import {useMainStore} from '../stores/ws'


const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'home',
      component: HomeView
    },
    {
        path: '/chat',
        name: 'chat',
        component: ChatView
    }
  ]
})

router.beforeEach((to, from) => {
    const mainStore = useMainStore()
    if(to.name === 'chat' && !mainStore.socket) {
        return {path: '/'}
    }
})

export default router