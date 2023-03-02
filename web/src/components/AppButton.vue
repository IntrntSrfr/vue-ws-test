<template>
    <button class="btn" :disabled="!!inactive" >
        {{ text }}
    </button>
</template>

<script setup lang="ts">
interface Props {
    text: string
    inactive?: boolean
}

defineProps<Props>()
</script>

<style scoped>
.btn {
    position: relative;
    display: inline-block;
    background-color: transparent;
    border: 1px solid rgb(30, 144, 255);
    border-radius: 5px;
    padding: 0.25em 0.5em;
    cursor: pointer;

    transition: 0.25s;
}

.btn:hover {
    border-color: rgb(71, 163, 255);
    background-color: rgb(0, 33, 66);
}

.btn:disabled{
    filter: grayscale(100%);
}
</style>
