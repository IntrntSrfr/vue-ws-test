import { useStorage } from '@vueuse/core'
import jwtDecode from 'jwt-decode'
import { defineStore } from 'pinia'
import { computed, reactive } from 'vue'

interface AuthState {
    id: string
    username: string
    loggedIn: boolean
}

const defaultState = (): AuthState => ({ id: '', username: '', loggedIn: false })

//localStorage.setItem('token', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzQ1Njc4OTAiLCJ1c2VybmFtZSI6ImFsZXgiLCJleHAiOjE2Nzc1NTI3OTZ9.b2FzaP-7X8zUtrpYJUBpcS4P6zQxHFjmCUCw8_V-L6M')

export const useAuthStore = defineStore('auth', () => {
    const state = reactive<AuthState>(defaultState())
    const token = useStorage('token', '')
    if (!token.value) {
        token.value = null
        return
    }

    const decoded: { exp: number; id: string; username: string } = jwtDecode(token.value)
    if (decoded.exp < Date.now() / 1000) {
        token.value = null
        return
    }

    state.id = decoded.id
    state.username = decoded.username
    state.loggedIn = true

    const loggedIn = computed(() => state.loggedIn)

    return { state, token, loggedIn }
})
