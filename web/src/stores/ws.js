import { computed, reactive, ref } from 'vue'
import { defineStore } from 'pinia'

export const useMainStore = defineStore('main', () => {
  const socket = ref(null)
  const username = ref('')
  const state = reactive({
    messages: [],
    users: []
  })

  const connect = () => {
    const ws = new WebSocket(`ws://localhost:8080/ws?username=${username.value}`)

    ws.onopen = (e) => {
        console.log("open", e);
    }
    
    ws.onclose = (e) => {
        console.log("close", e);
        socket.value = null
    }
    
    ws.onmessage = (data) => {
        console.log(data.data);
        const msg = JSON.parse(data.data)
        if(!msg?.op) return

        switch (msg.op) {
            case 0: // join
                break;
            case 1: // leave
                break;
            case 2: // message
                break;
            case 3: // ping 
                break;
            case 4: // ready
                state.messages = msg.data.messages
                state.users = msg.data.users

                break;
        
            default:
                break;
        }
        console.log("message", msg);
    }

    socket.value = ws
  }

  const sendMessage = (msg) => {
    let data = {
        text: msg
    }
    socket.value.send(JSON.stringify(data))
  }

  const setUsername = (newUsername) => {
    username.value = newUsername
  }

  const messages = computed(() => state.messages)
  const users = computed(() => state.users)
  
  return { socket, connect, sendMessage, username, setUsername, users, messages }
})
