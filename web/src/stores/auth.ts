import { useStorage } from '@vueuse/core'
import jwtDecode from 'jwt-decode'
import { defineStore } from 'pinia'
import { computed, reactive } from 'vue'
import http from '@/http'

interface AuthState {
    id: string
    username: string
    loggedIn: boolean
}

const defaultState = (): AuthState => ({ id: '', username: '', loggedIn: false })

//localStorage.setItem('token', 'eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6IjEyMzQ1Njc4OTAiLCJ1c2VybmFtZSI6ImFsZXgiLCJleHAiOjE2Nzc2NTI3OTZ9.B6KzOK22D7dYJBOzjDTDfopGBjm-agwrQsl7Nfl7bhE')

interface AuthResponse {
    token: string
}

export const useAuthStore = () => {
    const store = defineStore('auth', () => {
        const state = reactive<AuthState>(defaultState())
        const token = useStorage<string>('token', '')

        const initialize = () => {
            if (!token.value) {
                token.value = null
                return
            }
            decodeToken()
        }

        const decodeToken = () => {
            const decoded: { exp: number; sub: string; username: string } = jwtDecode(token.value)
            if (decoded.exp < Date.now() / 1000) {
                token.value = ''
                return
            }

            state.id = decoded.sub
            state.username = decoded.username
            state.loggedIn = true
        }

        const logout = () => {
            Object.assign(state, defaultState())
            token.value = null
        }

        const login = async (username: string, password: string) => {
            const res = await http.post<AuthResponse>('/auth/login', { username, password })
            token.value = res.data.token
            decodeToken()
        }

        const register = async (username: string, password: string) => {
            const res = await http.post<AuthResponse>('/auth/register', { username, password })
            token.value = res.data.token
            decodeToken()
        }

        const loggedIn = computed(() => state.loggedIn)

        return { initialize, state, token, loggedIn, login, logout, register }
    })

    const s = store()
    s.initialize()
    return store()
}
