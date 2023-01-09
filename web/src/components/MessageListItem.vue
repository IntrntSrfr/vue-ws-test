<template>
    <div class="message">
        <div class="author">
            {{ author }}
        </div>
        <div class="content">
            {{ content }}
        </div>
    </div>
</template>

<script setup>

defineProps({
    author: String,
    content: String
})

</script>

<style scoped>

</style>
