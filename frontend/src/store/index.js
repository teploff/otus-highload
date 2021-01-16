import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        tokenPair: {
            accessToken: null,
            refreshToken: null,
        },
        searchPeople: {
            anthroponym: null,
        }
    },
    getters: {
        accessToken: state => {
          return state.tokenPair.accessToken
        },
        refreshToken: state => {
          return state.tokenPair.refreshToken
        },
        searchAnthroponym: state => {
            return state.searchPeople.anthroponym
        },
    },
    mutations: {
        changeAccessToken (state, token) {
            state.tokenPair.accessToken = token
        },
        changeRefreshToken (state, token) {
            state.tokenPair.refreshToken = token
        },
        changeAnthroponym (state, anthroponym) {
            state.searchPeople.anthroponym = anthroponym
        }
    },
    actions: {},
    plugins: [createPersistedState()],
});