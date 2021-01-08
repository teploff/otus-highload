import Vue from "vue";
import Vuex from "vuex";

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
        }
    },
    mutations: {
        changeAnthroponym (state, payload) {
            state.searchPeople.anthroponym = payload
        }
    },
    actions: {}
});