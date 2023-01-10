<template>
    <div class="user-list">
        <UserListItem v-for="(usr, i) in users" :key="i" :username="usr.username"/> 
    </div>
</template>

<script setup>
import UserListItem from './UserListItem.vue';

defineProps({
    users: Array
})
</script>

<style scoped>


.user-list {
    display: flex;
    flex-direction: column;

    overflow-y: auto;
}


.user + .user{
    border-top: 1px solid gray;
}
</style>
