<template>
    <div class="message-list">
        <MessageListItem v-for="(msg, i) in messages" :key="i" :msg="msg"/> 
    </div>
</template>

<script setup>
import MessageListItem from './MessageListItem.vue';

defineProps({
    messages: Array
})
</script>

<style scoped>

</style>
