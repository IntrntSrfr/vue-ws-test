<template>
    <div class="user">
        {{ username }}
    </div>
</template>

<script setup>

defineProps({
    username: String
})

</script>

<style scoped>

.user {
    padding: .5em;
    font-weight: bold;
}

</style>
