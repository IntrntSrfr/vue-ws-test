<template>
    <div class="message">
        <div class="header">
            <div class="author">
                {{ author }}
            </div>
            <div class="timestamp">
                {{ moment(timestamp).calendar() }}
            </div>
        </div>
        <div class="content">
            {{ content }}
        </div>
    </div>
</template>

<script setup lang="ts">
import moment from 'moment' // consider moving to list to reduce imports?

defineProps<{
    author: string
    timestamp: Date
    content: string
}>()
</script>

<style scoped>
.message {
    display: flex;
    flex-direction: column;
    padding: 0.5em;
}

.header {
    display: flex;
    flex-direction: row;
    align-items: center;
}

.author {
    font-weight: bold;
}

.timestamp {
    margin-left: 0.25em;
    color: rgba(235, 235, 235, 0.64);
    font-size: 0.75rem;
    font-weight: lighter;
}
</style>
