<template>
    <div class="chat">
        <div class="messages">
            <MessageList :messages="[1,2,3]"/>
            <div class="chat-input">
                <AppChatInput />
                <AppButton />
            </div>
        </div>
        <div class="users">
            <UserList :users="[1,2,3]"/>
        </div>
    </div>
</template>

<script setup>
import AppButton from './AppButton.vue';
import AppChatInput from './AppChatInput.vue';
import MessageList from './MessageList.vue'
import UserList from './UserList.vue'

</script>

<style scoped>

</style>