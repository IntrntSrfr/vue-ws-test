import { createRouter, createWebHistory } from 'vue-router'
import HomeView from '../views/HomeView.vue'
import ChatView from '../views/ChatView.vue'
import AuthView from '../views/AuthView.vue'

import { useAuthStore } from '@/stores/auth'

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
        },
        {
            path: '/auth',
            name: 'auth',
            component: AuthView
        }
    ]
})

router.beforeEach((to) => {
    const authStore = useAuthStore()
    console.log(authStore.loggedIn)

    if (to.name !== 'auth' && !authStore.loggedIn) {
        return { path: '/auth' }
    }
    if (to.name !== 'chat' && authStore.loggedIn) {
        return { path: '/chat' }
    }
})

export default router
