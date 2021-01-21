let wsService = {}

wsService.install = function (Vue, options) {
    let ws = null
    let reconnectInterval = options.reconnectInterval || 1000

    Vue.prototype.$wsConnect = (url) => {
        ws = new WebSocket(url)

        ws.onopen = () => {
            // Restart reconnect interval
            reconnectInterval = options.reconnectInterval || 1000
            options.store.commit('establishWSConn')
        }

        ws.onmessage = (event) => {
            // New message from the backend - use JSON.parse(event.data)
            parseMsg(event)
        }

        ws.onclose = (event) => {
            options.store.commit('disbandWSConn')

            if (event) {
                // Event.code 1000 is our normal close event
                if (event.code !== 1000) {
                    let maxReconnectInterval = options.maxReconnectInterval || 3000
                    setTimeout(() => {
                        if (reconnectInterval < maxReconnectInterval) {
                            // Reconnect interval can't be > x seconds
                            reconnectInterval += 1000
                        }
                        Vue.prototype.$wsConnect()
                    }, reconnectInterval)
                }
            }
        }

        ws.onerror = (error) => {
            options.store.commit('disbandWSConn')
            console.log(error)
            ws.close()
        }
    }

    Vue.prototype.$wsDisconnect = () => {
        // Our custom disconnect event
        options.store.commit('disbandWSConn')
        ws.close()
    }

    Vue.prototype.$wsSend = (data) => {
        // Send data to the backend - use JSON.stringify(data)
        ws.send(JSON.stringify(data))
    }

    /*
      Here we write our custom functions to not make a mess in one function
    */
    function parseMsg (params) {
        console.log(params)
        options.store.commit('', params.data)
    }
}

export default wsService
