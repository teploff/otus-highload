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
        },
        news: {
            data: [],
            count: 0
        }
    },
    getters: {
        isWSConnected: state => {
            return state.socket.isConnected
        },
        searchAnthroponym: state => {
            return state.searchPeople.anthroponym
        },
        news: state => {
            return state.news
        }
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
        },
        setNews(state, payload) {
            state.news.count = payload.count
            state.news.data = payload.news
        },
        appendNews(state, news) {
            state.news.count += 1
            state.news.data.push(news)
            state.news.data = state.news.data.sort(function(a,b) {
                return new Date(b.createTime) - new Date(a.createTime);
            });
        }
    },
    actions: {},
    plugins: [createPersistedState()],
});
