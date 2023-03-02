<script setup lang="ts">
import { RouterView } from 'vue-router'
import {useSocketStore} from "@/stores/ws";
const socketStore = useSocketStore()
socketStore.connect()
</script>

<template>
    <RouterView />
</template>

<style scoped></style>
