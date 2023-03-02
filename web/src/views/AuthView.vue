<template>
    <RegistrationForm v-if="showRegistration" @submit="onRegister" class="slide-left" />
    <LoginForm v-else @submit="onLogin" class="slide-right" />
    <AppButton :text="buttonText" @click="() => (showRegistration = !showRegistration)" />
</template>

<script setup lang="ts">
import type { LoginFormEmit, RegisterFormEmit } from '@/types'
import { ref, computed } from 'vue'
import AppButton from '../components/AppButton.vue'
import LoginForm from '../components/LoginForm.vue'
import RegistrationForm from '../components/RegistrationForm.vue'
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
import { AxiosError } from 'axios'
import {useSocketStore} from "@/stores/ws";

const socketStore = useSocketStore()
const authStore = useAuthStore()
const router = useRouter()

const showRegistration = ref(false)
const buttonText = computed(() => {
    return showRegistration.value ? 'Already have a user?' : 'No user?'
})

const onRegister = ({ username, password, password2 }: RegisterFormEmit) => {
    if (password !== password2) return
    authStore
        .register(username, password)
        .then(() => {
            socketStore.connect();
            router.push('/chat')
        })
        .catch((err: AxiosError) => console.log(err.response?.data))
}

const onLogin = ({ username, password }: LoginFormEmit) => {
    authStore
        .login(username, password)
        .then(() => {
            socketStore.connect();
            router.push('/chat')
        })
        .catch((err: AxiosError) => console.log(err.response?.data))
}
</script>

<style scoped></style>
