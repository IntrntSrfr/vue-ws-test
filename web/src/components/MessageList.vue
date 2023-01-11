<template>
    <div class="message-list" ref="listRef">
        <MessageListItem v-for="(msg, i) in messages" :key="i" :author="msg.username" :content="msg.text" :timestamp="msg.timestamp"/> 
    </div>
</template>

<script setup>
import { ref, watch } from 'vue';
import { useMainStore } from '../stores/ws';
import MessageListItem from './MessageListItem.vue';

const mainStore = useMainStore()

const listRef = ref(null)

defineProps({
    messages: Array
})

const scrollBottom = () => {
    setTimeout(() => {
        listRef.value.scroll({top: listRef.value.scrollHeight+100, behavior: 'smooth'})
        //listRef.value.scroll(0, listRef.value.scrollHeight+100)
    }, 10);
}

watch(() => mainStore.messages, (n, o) => {
    scrollBottom()
})

scrollBottom()

</script>

<style scoped>

.message-list {
    flex-grow: 1;
    display: flex;
    flex-direction: column;

    overflow-y: auto;
}

.message + .message{
    border-top: 1px solid gray;
}

</style>
