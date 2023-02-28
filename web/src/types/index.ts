export interface User {
    id: string
    username: string
}

export interface Message {
    text: string
    username: string
    timestamp: string
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
