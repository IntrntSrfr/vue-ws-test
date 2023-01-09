<template>
    <div class="message-list">
        <MessageListItem v-for="(msg, i) in messages" :key="i" :author="msg.username" :content="msg.text"/> 
    </div>
</template>

<script setup>
import MessageListItem from './MessageListItem.vue';

defineProps({
    messages: Array
})
</script>

<style scoped>

.message + .message{
    border-top: 1px solid gray;
}

</style>
