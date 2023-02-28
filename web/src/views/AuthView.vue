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

const showRegistration = ref(false)
const buttonText = computed(() => {
    return showRegistration.value ? 'Already have a user?' : 'No user?'
})

const onRegister = ({ username, password, password2 }: RegisterFormEmit) => {
    console.log(username, password, password2)
}

const onLogin = ({ username, password }: LoginFormEmit) => {
    console.log(username, password)
}
</script>

<style scoped></style>
