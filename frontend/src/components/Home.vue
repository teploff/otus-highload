<template>
  <div class="page-container">
    <md-app md-waterfall md-mode="flexible">
      <md-app-toolbar class="md-large md-primary">
        <div class="md-toolbar-row">
          <div class="md-toolbar-section-start">
            <md-button class="md-icon-button" @click="menuVisible = !menuVisible">
              <md-icon>menu</md-icon>
            </md-button>
          </div>

          <md-autocomplete
              class="search"
              v-model.trim="selectedHuman"
              @input="searchPeople"
              md-layout="box"
              :md-options="stub">
            <label>Search people...</label>
          </md-autocomplete>

          <div class="md-toolbar-section-end">
            <md-button class="md-icon-button">
              <md-icon>login</md-icon>
            </md-button>
          </div>
        </div>
      </md-app-toolbar>

      <md-app-drawer :md-active.sync="menuVisible">
        <md-list>
          <md-list-item @click="followHomePage">
            <md-icon>assignment_ind</md-icon>
            <span class="md-list-item-text">Моя страница</span>
          </md-list-item>

          <md-list-item @click="followNewsPage">
            <md-icon>fiber_new</md-icon>
            <span class="md-list-item-text">Новости</span>
            <md-badge v-if="countNewsNotify > 0" class="md-primary" v-bind:md-content="countNewsNotify" />
          </md-list-item>

          <md-list-item @click="followMessengerPage">
            <md-icon>chat</md-icon>
            <span class="md-list-item-text">Мессенджер</span>
            <md-badge v-if="countMsgNotify > 0" class="md-primary" v-bind:md-content="countMsgNotify" />
          </md-list-item>

          <md-list-item @click="followFriendsPage">
            <md-icon>supervisor_account</md-icon>
            <span class="md-list-item-text">Друзья</span>
            <md-badge v-if="countFriendsNotify > 0" class="md-primary" v-bind:md-content="countFriendsNotify" />
          </md-list-item>
        </md-list>
      </md-app-drawer>

      <md-app-content>
        Home payload!
      </md-app-content>
    </md-app>
  </div>
</template>

<script>
import axios from 'axios';
import { apiUrl, debounce } from '@/const'

export default {
  name: 'Home',
  data: () => ({
    menuVisible: false,
    countNewsNotify: 0,
    countMsgNotify: 0,
    countFriendsNotify: 0,
    searchPayload: {
      anthroponym: null,
    },
    selectedHuman: null,
    people: {
      questionnaires: [],
      count: 0
    },
    stub: []
  }),
  methods: {
    refreshToken() {
      const path = `${apiUrl}/auth/token`;
      const refreshToken = localStorage.getItem('refresh_token');

      if (refreshToken === null) {
        this.$router.push({ name: 'SignIn' });
      }

      const payload = {
        refresh_token: refreshToken,
      };
      axios.put(path, payload)
          .then((response) => {
            const tokenPair = JSON.parse(JSON.stringify(response.data));

            localStorage.setItem('access_token', tokenPair.access_token);
            localStorage.setItem('refresh_token', tokenPair.refresh_token);

            this.$router.push({ name: 'Home' });
          })
          .catch((error) => {
            const err = JSON.parse(JSON.stringify(error.response));
            if (err.status === 401) {
              this.$router.push({ name: 'SignIn' });
            }
            console.log(err);
          });
    },
    followHomePage() {
      this.$router.push({ name: 'Home' }).catch(() => {});
    },
    followNewsPage() {
      this.$router.push({ name: 'News' });
    },
    followMessengerPage() {
      this.$router.push({ name: 'Messenger' });
    },
    followFriendsPage() {
      this.$router.push({ name: 'Friends' });
    },
    searchPeople: debounce(function (){
      this.$store.commit("changeAnthroponym", this.selectedHuman);
      this.$router.push({ name: 'People' })
    }, 1000),
  },
};
</script>

<style scoped>
.md-app {
  max-height: 100vh;
  min-height: 100vh;
  border: 1px solid rgba(#000, .12);
}

.md-drawer {
 width: 230px;
 max-width: calc(100vw - 125px);
}

.search {
  max-width: 500px;
}

.md-toolbar {
  height: 50px;
  padding: inherit;
}

</style>
