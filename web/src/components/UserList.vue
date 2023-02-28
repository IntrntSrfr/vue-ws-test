<template>
    <div class="user-list">
        <UserListItem v-for="(usr, i) in users" :key="i" :username="usr.username" />
    </div>
</template>

<script setup lang="ts">
import UserListItem from './UserListItem.vue'
import type { User } from '@/types'

interface Props {
    users: User[]
}

defineProps<Props>()
</script>

<style scoped>
.user-list {
    display: flex;
    flex-direction: column;

    overflow-y: auto;
}

.user + .user {
    border-top: 1px solid gray;
}
</style>
