import { ref, computed } from 'vue'
import { defineStore } from 'pinia'

export const useSocketStore = defineStore('socket', () => {
  const socket = ref(0)
  const doubleCount = computed(() => count.value * 2)
  function increment() {
    count.value++
  }

  return { count, doubleCount, increment }
})
