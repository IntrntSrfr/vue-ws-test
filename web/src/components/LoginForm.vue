<template>
    <form class="form">
        <div class="form-fields">
            <div class="form-field" v-for="(field, i) in fields" :key="i">
                <div class="field-title">{{ field.title }}</div>
                <div class="text-input">
                    <input :type="field.type" v-model="field.value" />
                </div>
            </div>
        </div>
        <div class="submit">
            <AppButton text="Login" @click="submit" />
        </div>
    </form>
</template>

<script setup lang="ts">
import type { LoginFormEmit } from '@/types'
import { reactive } from 'vue'
import AppButton from './AppButton.vue'

const emit = defineEmits<{
    (e: 'submit', data: LoginFormEmit): void
}>()

const submit = () => {
    emit('submit', {
        username: fields.username.value,
        password: fields.password.value
    })
}

const fields = reactive({
    username: { title: 'Username', type: 'text', value: '' },
    password: { title: 'Password', type: 'password', value: '' }
})
</script>

<style scoped>
.form {
    display: flex;
    flex-direction: column;
    padding: 0.5em;
    gap: 1em;
}

.form-field + .form-field {
    margin-top: 1em;
}

.field-title {
    margin-bottom: 0.25em;
    margin-left: 0.5em;
}

.text-input {
    display: flex;
    align-items: center;
    padding: 0.5em;
    border: 1px solid var(--color-border);
    text-align: center;
    border-radius: 5px;
}

.text-input svg {
    margin: 0 0.5em;
}

input {
    flex-grow: 1;
    background-color: transparent;
    border: none;
    outline: none;
    font-size: 1rem;
}

.submit {
    margin-top: 1em;
    align-self: center;
}

.submit .btn {
    font-weight: 500;
    padding: 1em 2em;
}
</style>
