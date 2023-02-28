import { computed, reactive, ref } from 'vue'
import { defineStore } from 'pinia'
import type { Message, User } from '@/types'
import { useAuthStore } from './auth'

interface SocketMessageData {
    op: number
    data: any
}

interface SocketState {
    messages: Message[]
    users: User[]
}

export const useSocketStore = defineStore('socket', () => {
    const authStore = useAuthStore()
    console.log(authStore.loggedIn)

    const socket = ref<WebSocket | null>(null)
    const username = ref('')
    const state = reactive<SocketState>({
        messages: [],
        users: []
    })
    const pingInterval = ref<number | null>(null)

    const connect = () => {
        if (!authStore.loggedIn) return
        const ws = new WebSocket(`ws://localhost:8080/ws?username=${username.value}`)

        ws.onopen = () => {
            pingInterval.value = setInterval(() => {
                console.log('weed')
            }, 5000)
        }

        ws.onclose = () => {
            disconnect()
        }

        ws.onmessage = (data: MessageEvent) => {
            const msg: SocketMessageData = JSON.parse(data.data)
            if (!msg) return

            switch (msg.op) {
                case 0: // join
                    state.users = [...state.users, msg.data.user as User]
                    break
                case 1: // leave
                    state.users = state.users.filter((u) => u.id !== (msg.data.user as User).id)
                    break
                case 2: // message
                    state.messages = [...state.messages, msg.data.message as Message]
                    break
                case 3: // ping
                    break
                case 4: // ready
                    state.messages = msg.data.messages as Message[]
                    state.users = msg.data.users as User[]
                    break

                default:
                    break
            }
        }

        socket.value = ws
    }

    const disconnect = () => {
        socket.value = null
        pingInterval.value = null
    }

    const sendMessage = (msg: string) => {
        if (!socket.value) return
        const data = {
            text: msg
        }
        socket.value.send(JSON.stringify(data))
    }

    const setUsername = (newUsername: string) => {
        username.value = newUsername
    }

    const messages = computed(() => state.messages)
    const users = computed(() => state.users)

    return { socket, connect, disconnect, sendMessage, username, setUsername, users, messages }
})
