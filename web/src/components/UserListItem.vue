<template>
    <div class="user">
        {{ username }}
    </div>
</template>

<script setup lang="ts">
interface Props {
    username: string
}

defineProps<Props>()
</script>

<style scoped>
.user {
    padding: 0.5em;
    font-weight: bold;
}
</style>
