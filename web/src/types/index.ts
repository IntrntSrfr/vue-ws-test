export interface User {
    id: string
    username: string
    created: Date
}

export interface Message {
    id: string
    content: string
    author: User
    timestamp: Date
}

export interface RegisterFormEmit {
    username: string
    password: string
    password2: string
}

export interface LoginFormEmit {
    username: string
    password: string
}
