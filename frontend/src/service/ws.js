let wsService = {}

wsService.install = function (Vue, options) {
    Vue.prototype.$wsConnect = (url) => {
        console.log("ws connect!!")
        console.log(url)
        Vue.prototype.ws = new WebSocket(url)

        Vue.prototype.ws.onopen = () => {
            console.log('ws: onopen')
            // Restart reconnect interval
            options.store.commit('establishWSConn')
        }

        Vue.prototype.ws.onmessage = (event) => {
            console.log('ws: onmessage')
            // New message from the backend - use JSON.parse(event.data)
            parseMsg(event)
        }

        Vue.prototype.ws.onclose = (event) => {
            console.log('ws: onclose')
            options.store.commit('disbandWSConn')
            console.log(event)
            Vue.prototype.ws.close()
        }

        Vue.prototype.ws.onerror = (error) => {
            console.log('ws: onerror')
            options.store.commit('disbandWSConn')
            console.log(error)
            Vue.prototype.ws.close()
        }
    }

    Vue.prototype.$wsDisconnect = () => {
        // Our custom disconnect event
        options.store.commit('disbandWSConn')
        Vue.prototype.ws.close()
    }

    Vue.prototype.$wsSend = (data) => {
        // Send data to the backend - use JSON.stringify(data)
        Vue.prototype.ws.send(data)
    }

    /*
      Here we write our custom functions to not make a mess in one function
    */
    function parseMsg (params) {
        const camelcaseKeys = require('camelcase-keys');

        console.log(params)
        options.store.commit('appendNews', camelcaseKeys(JSON.parse(params.data), { deep: true }))
    }
}

export default wsService
