<template>
    <div class="btn">
        {{ text }}
    </div>
</template>

<script setup lang="ts">
defineProps({
    text: {
        type: String,
        required: true
    }
})
</script>

<style scoped>
.btn {
    position: relative;
    display: inline-block;
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
</style>
