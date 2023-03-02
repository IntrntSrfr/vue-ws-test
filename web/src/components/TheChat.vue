<template>
    <div class="chat">
        <div class="messages">
            <h2>Messages</h2>
            <MessageList :messages="socketStore.messages" />
            <div class="chat-input">
                <div class="chat-input-inner" @click="focusInput">
                    <textarea v-model="message" ref="chatBoxRef" />
                </div>
                <div class="chat-btns">
                    <AppButton text="Send" :inactive="!validMessage" @click="sendMessage" />
                </div>
            </div>
        </div>
        <div class="users">
            <h2>Users</h2>
            <UserList :users="socketStore.users" />
        </div>
    </div>
</template>

<script setup lang="ts">
import {computed, ref} from 'vue'
import AppButton from './AppButton.vue'
import MessageList from './MessageList.vue'
import UserList from './UserList.vue'

import { useSocketStore } from '@/stores/ws'
import http from "@/http";
import type {Message} from "@/types";
import {useAuthStore} from "@/stores/auth";

const socketStore = useSocketStore()
const authStore = useAuthStore()

const chatBoxRef = ref<HTMLInputElement | null>(null)
const focusInput = () => {
    if(!chatBoxRef.value) return;
    chatBoxRef.value.focus()
}

const message = ref<string>('')

const validMessage = computed(() => {
    return !!message.value?.trim()
})

const sendMessage = () => {
    const token = authStore.token
    if(!validMessage) return;
    const tmpMsg = message.value.trim()
    message.value = ''

    http.post<Message>('/messages/', {content: tmpMsg}, {headers: {'Authorization': `Bearer ${token}`}})
        .then(resp => console.log(resp.data))
        .catch(err => console.log(err))
}

</script>

<style scoped>
.chat {
    height: 100%;
    display: flex;
    flex-direction: row;
    gap: 1em;
}

.messages,
.users {
    border: 1px solid rgb(184, 184, 184);
}

.messages {
    flex-grow: 1;
    display: flex;
    flex-direction: column;
}

.chat-input {
    display: flex;
    flex-direction: row;
    justify-content: space-between;
    background-color: var(--color-background-soft);
}

.chat-input-inner{
    flex-grow: 1;
    cursor: text;
}

.chat-input-inner textarea {
    width: 100%;
    height: 100%;

    resize: none;
    padding:  0.5em;

    background-color: transparent;
    border: none;
    outline: none;
}

.users {
    display: flex;
    flex-direction: column;
    width: 15em;
}

h2 {
    padding: 0.5em 1em;
    background-color: var(--color-background-soft);
    font-weight: bold;
}
</style>
