<template>
    <input type="text" :value="text" @input="onInput" />
</template>

<script setup lang="ts">
defineProps({
    text: {
        type: String,
        default: ''
    }
})

const emit = defineEmits<{
    (e: 'input', text: string): void
}>()

const onInput = (e: Event) => {
    if (!e.target) return
    emit('input', (e.target as HTMLInputElement).value)
}
</script>

<style scoped></style>
