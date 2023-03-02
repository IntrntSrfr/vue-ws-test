import axios from 'axios'

axios.defaults.headers.post['Content-Type'] = 'application/json'

export default axios.create({
    baseURL: 'http://localhost:7070/api'
})
