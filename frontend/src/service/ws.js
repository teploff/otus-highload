export default class WSService {
    constructor(store) {
        this.store = store
    }
    connect = (url) => {
        this.ws = new WebSocket(url + '?token=' + localStorage.getItem('accessToken'))

        this.ws.onopen = () => {
            console.log("ws connection is opened")
        }

        this.ws.onmessage = (event) => {
            console.log('ws get message')

            const camelcaseKeys = require('camelcase-keys');
            this.store.commit('appendNews', camelcaseKeys(JSON.parse(event.data), {deep: true}))
        }

        this.ws.onclose = (event) => {
            console.log('ws is closed onclose')
            console.log(event)
            this.ws.close()
        }

        this.ws.onerror = (error) => {
            console.log('ws is error')
            console.log(error)
            this.ws.close()
        }
    }

    disconnect = () => {
        this.ws.close()
    }

    send = (data) => {
        this.ws.send(data)
    }
}