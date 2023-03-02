<template>
    <div class="message-list" ref="listRef">
        <MessageListItem
            v-for="(msg, i) in messages"
            :key="i"
            :author="msg.author.username"
            :content="msg.content"
            :timestamp="new Date(msg.timestamp)"
        />
    </div>
</template>

<script setup lang="ts">
import type { Message } from '@/types'
import { ref, watch } from 'vue'
import { useSocketStore } from '@/stores/ws'
import MessageListItem from './MessageListItem.vue'

const mainStore = useSocketStore()
const listRef = ref<HTMLDivElement | null>(null)

interface Props {
    messages: Message[]
}

defineProps<Props>()

const scrollBottom = () => {
    if (!listRef.value) return
    listRef.value.scroll({ top: listRef.value.scrollHeight + 100, behavior: 'smooth' })
}

watch(
    () => mainStore.messages,
    (n, o) => {
        setTimeout(() => {
            scrollBottom()
        }, 10)
    }
)

scrollBottom()
</script>

<style scoped>
.message-list {
    flex-grow: 1;
    display: flex;
    flex-direction: column;

    overflow-y: auto;
}

.message + .message {
    border-top: 1px solid gray;
}
</style>
