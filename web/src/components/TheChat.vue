<template>
    <div class="chat">
        <div class="messages">
            <MessageList :messages="mainStore.messages"/>
            <div class="chat-input">
                <AppInput :text="message" @input="setMessage" />
                <AppButton text="Send" @click="mainStore.sendMessage(message)"/>
            </div>
        </div>
        <div class="users">
            <UserList :users="mainStore.users"/>
        </div>
    </div>
</template>

<script setup>
import { ref } from 'vue';
import AppButton from './AppButton.vue';
import AppInput from './AppInput.vue';
import MessageList from './MessageList.vue'
import UserList from './UserList.vue'

import { useMainStore } from '../stores/ws';

const mainStore = useMainStore()

const message = ref('')
const setMessage = (newMsg) => {
    message.value = newMsg
}

</script>

<style scoped>

</style>