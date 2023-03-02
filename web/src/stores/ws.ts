import { computed, reactive, ref } from 'vue'
import { defineStore } from 'pinia'
import type { Message, User } from '@/types'
import { useAuthStore } from './auth'

enum OpCode {
    Identify = 0,
    Ping,
    PingACK,
    Action,
    Error
}

enum ActionCode {
    None = 0,
    UserReady,
    UserJoin,
    UserLeave,
    UserMessage
}

enum ErrorCode {
    Unknown = 0,
    PingTimedOut,
    AuthFailed
}

interface SocketEvent {
    op: OpCode
    data: any
    action: ActionCode
}
interface IdentifyData {
    token: string
}

interface PingData {
    sequence: number
}

interface ErrorData {
    code: ErrorCode
    message: string
}

interface ActionEvent {}
interface UserReadyData extends ActionEvent {
    messages: Message[]
    users: User[]
}
interface UserJoinData extends ActionEvent {
    user: User
}
interface UserLeaveData extends ActionEvent {
    user: User
}
interface UserMessageData extends ActionEvent {
    message: Message
}

interface SocketState {
    messages: Message[]
    users: User[]
}

const defaultState = (): SocketState =>({messages: [], users: []})

export const useSocketStore = defineStore('socket', () => {
    const authStore = useAuthStore()
    console.log(authStore.loggedIn)

    const socket = ref<WebSocket | null>(null)
    const username = ref('')
    const state = reactive<SocketState>(defaultState())
    const pingInterval = ref<number | null>(null)

    const connect = () => {
        if (!authStore.loggedIn) return
        const ws = new WebSocket(`ws://localhost:7070/ws`)

        ws.onopen = () => {
            ws.send(JSON.stringify({
                op: OpCode.Identify,
                data: {token: authStore.token},
                action: ActionCode.None
            } as SocketEvent))

            pingInterval.value = setInterval(() => {
                console.log('weed')
            }, 5000)
        }

        ws.onerror = () => {
            disconnect()
        }

        ws.onclose = () => {
            disconnect()
        }

        ws.onmessage = (data: MessageEvent) => {
            const msg: SocketEvent = JSON.parse(data.data)
            if (!msg) return
            handleEvent(msg)
        }

        socket.value = ws
    }

    const handleEvent = (evt: SocketEvent) => {
        switch (evt.op) {
            case OpCode.Action:
                handleAction(evt.action, evt.data as ActionEvent)
                break

            default:
                break
        }
    }

    const handleAction = (code: ActionCode, evt: ActionEvent) => {
        switch (code) {
            case ActionCode.UserReady:
                handleUserReady(evt as UserReadyData)
                break
            case ActionCode.UserJoin:
                handleUserJoin(evt as UserJoinData)
                break
            case ActionCode.UserLeave:
                handleUserLeave(evt as UserLeaveData)
                break
            case ActionCode.UserMessage:
                handleUserMessage(evt as UserMessageData)
                break
            default:
                break
        }
    }

    const handleUserReady = (evt: UserReadyData) => {
        state.messages = evt.messages as Message[]
        state.users = evt.users as User[]
    }

    const handleUserJoin = (evt: UserJoinData) => {
        state.users = [...state.users, evt.user as User]
    }

    const handleUserLeave = (evt: UserLeaveData) => {
        state.users = state.users.filter((u) => u.id !== (evt.user as User).id)
    }

    const handleUserMessage = (evt: UserMessageData) => {
        state.messages = [...state.messages, evt.message as Message]
    }

    const disconnect = () => {
        socket.value = null
        pingInterval.value = null
        Object.assign(state, defaultState())
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
