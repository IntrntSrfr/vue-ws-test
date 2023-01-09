<template>
    <div class="user-list">
        <UserListItem v-for="(usr, i) in users" :key="i" :username="usr"/> 
    </div>
</template>

<script setup>
import UserListItem from './UserListItem.vue';

defineProps({
    users: Array
})
</script>

<style scoped>

</style>
