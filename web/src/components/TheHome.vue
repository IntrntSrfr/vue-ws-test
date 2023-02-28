<template>
    <div>
        <h2>set your username</h2>
        <AppInput :text="mainStore.username" @input="mainStore.setUsername" />
        <AppButton text="Go" @click="login" />
    </div>
</template>

<script setup lang="ts">
import AppInput from './AppInput.vue'
import AppButton from './AppButton.vue'

import { useAuthStore } from '@/stores/auth'
import { useSocketStore } from '@/stores/ws'
import { useRouter } from 'vue-router'

const authStore = useAuthStore()
const mainStore = useSocketStore()
const router = useRouter()

const login = async () => {
    mainStore.connect()
    await router.push('/chat')
}
</script>

<style scoped></style>
