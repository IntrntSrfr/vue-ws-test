<template>
    <div>
        <h2>set your username</h2>
        <AppInput :text="mainStore.username" @input="mainStore.setUsername"/>
        <AppButton text="Go" @click="login"/>
    </div>
</template>

<script setup>
import AppInput from './AppInput.vue';
import AppButton from './AppButton.vue';

import {useMainStore} from '@/stores/ws'
import { useRouter } from 'vue-router';

const mainStore = useMainStore()
const router = useRouter()

const login = async () => {
    mainStore.connect()
    await router.push('/chat');
}

</script>

<style scoped>

</style>
