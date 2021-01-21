import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        socket: {
            isConnected: false,
            socketMessage: ''
        },
        searchPeople: {
            anthroponym: null,
        }
    },
    getters: {
        isWSConnected: state => {
            return state.socket.isConnected
        },
        searchAnthroponym: state => {
            return state.searchPeople.anthroponym
        },
    },
    mutations: {
        changeAnthroponym (state, anthroponym) {
            state.searchPeople.anthroponym = anthroponym
        },
        establishWSConn(state) {
            state.socket.isConnected = true
        },
        disbandWSConn(state) {
            state.socket.isConnected = false
        },
        setMsg(state, msg) {
            state.socket.socketMessage = msg
        }
    },
    actions: {},
    plugins: [createPersistedState()],
});
