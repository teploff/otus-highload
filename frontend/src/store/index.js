import Vue from "vue";
import Vuex from "vuex";
import createPersistedState from "vuex-persistedstate";

Vue.use(Vuex);

export default new Vuex.Store({
    state: {
        searchPeople: {
            anthroponym: null,
        }
    },
    getters: {
        searchAnthroponym: state => {
            return state.searchPeople.anthroponym
        },
    },
    mutations: {
        changeAnthroponym (state, anthroponym) {
            state.searchPeople.anthroponym = anthroponym
        }
    },
    actions: {},
    plugins: [createPersistedState()],
});
