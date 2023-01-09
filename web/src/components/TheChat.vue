<template>
    <div class="chat">
        <div class="messages">
            <h2>Messages</h2>
            <MessageList :messages="mainStore.messages"/>
            <div class="chat-input">
                <AppInput :text="message" @input="setMessage" />
                <AppButton text="Send" @click="mainStore.sendMessage(message)"/>
            </div>
        </div>
        <div class="users">
            <h2>Users</h2>
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

.chat{
    display: flex;
    flex-direction: row;
}

.messages, .users {
    border: 1px solid rgb(184, 184, 184);
}


</style>