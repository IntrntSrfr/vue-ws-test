<template>
    <input type="text" :value="text" @input="onInput">
</template>

<script setup>

defineProps({
    text: {
        type: String,
        default: ''
    }
})

const emit = defineEmits(['input'])

const onInput = (e) => {
    emit('input', e.target.value)
}

</script>

<style scoped>

</style>
