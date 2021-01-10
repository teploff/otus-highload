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
              v-model.trim="searchPayload.anthroponym"
              @input="searchPeople"
              :md-options="[]"
              md-layout="box">
            <label>Search people...</label>
          </md-autocomplete>

          <div class="md-toolbar-section-end">
            <md-button class="md-icon-button" @click="logOut">
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
        <md-table v-show="cards.count !== 0">
          <md-table-row>
            <div class="card-expansion">
        <md-card v-for="card in cards.news" v-bind:key="card.id">
          <md-card-header>
            <md-card-header-text>
              <div class="md-title">{{ card.owner.name }} {{ card.owner.surname }}</div>
              <div class="md-subhead">{{ card.content }}</div>
            </md-card-header-text>

            <md-card-media>
              <img v-if="card.owner.sex === 'male'" src="../assets/boy.png" alt="People">
              <img v-else src="../assets/girl.png" alt="People">
            </md-card-media>
          </md-card-header>
        </md-card>

            </div>
          </md-table-row>
          <md-table-row class="pagination-row">
            <div>
              <paginate
                  v-model="page"
                  :page-count="Math.ceil(cards.count / countCardsInWindow)"
                  :click-handler="paginatorClick"
                  :prev-text="'Prev'"
                  :next-text="'Next'"
                  :container-class="'pagination'"
                  :page-class="'page-item'"
                  :first-last-button="true"
              >
              </paginate>
            </div>
          </md-table-row>
        </md-table>
      </md-app-content>
    </md-app>
    <FlashMessage :position="'right top'"></FlashMessage>
  </div>
</template>

<script>
import {apiUrl, debounce, headers} from "@/const";
import axios from "axios";

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
    cards: {
      news: null,
      count: 0,
    },
    countCardsInWindow: 10,
    offset: 0,
    page: 1,
  }),
  beforeCreate() {
    if (this.$store.getters.accessToken === null) {
      this.$router.push({ name: 'SignIn' });
    }
  },
  created() {
    this.getNews();
  },
  methods: {
    followHomePage() {
      this.$router.push({ name: 'Home' });
    },
    followNewsPage() {
      this.$router.push({ name: 'News' }).catch(() => {});
    },
    followMessengerPage() {
      this.$router.push({ name: 'Messenger' });
    },
    followFriendsPage() {
      this.$router.push({ name: 'Friends' });
    },
    searchPeople: debounce(function (){
      this.$store.commit("changeAnthroponym", this.searchPayload.anthroponym);
      this.$router.push({ name: 'People' })
    }, 1000),
    paginatorClick(pageNum) {
      this.offset = (pageNum - 1) * this.countCardsInWindow;

      this.getNews()
    },
    getNews() {
      const path = `${apiUrl}/social/news`;
      const camelcaseKeys = require('camelcase-keys');

      headers.Authorization = this.$store.getters.accessToken

      axios.get(path, {
        headers: headers,
        params: {offset: this.offset},
        transformResponse: [(data) => {
          return camelcaseKeys(JSON.parse(data), { deep: true })}
        ]
      })
          .then((response) => {
            this.cards = response.data;
          })
          .catch((error) => {
            const err = error.response;

            if (err.status === 401) {
              this.refreshToken();
            }

            this.flashMessage.error({
              title: 'Error Message Title',
              message: err.data.message,
              position: 'center',
              icon: '../assets/error.svg',
            });
          });
    },
    refreshToken() {
      const path = `${apiUrl}/auth/token`;
      const camelcaseKeys = require('camelcase-keys');

      if (this.$store.getters.refreshToken === null) {
        this.$router.push({ name: 'SignIn' });
      }

      const payload = {
        refresh_token: this.$store.getters.refreshToken,
      };
      axios.put(path, payload, {transformResponse: [(data) => {
          return camelcaseKeys(JSON.parse(data), { deep: true })}
        ]})
          .then((response) => {
            this.tokenPair = response.data;

            this.$store.commit("changeAccessToken", this.tokenPair.accessToken);
            this.$store.commit("changeRefreshToken", this.tokenPair.refreshToken);

            this.$router.push({ name: 'People' });
          })
          .catch((error) => {
            const err = JSON.parse(JSON.stringify(error.response));
            if (err.status === 401) {
              this.$router.push({ name: 'SignIn' });
            }
            this.flashMessage.error({
              title: 'Error Message Title',
              message: err.data.message,
              position: 'center',
              icon: '../assets/error.svg',
            });
          });
    },
    logOut() {
      this.$store.commit("changeAccessToken", null);
      this.$store.commit("changeRefreshToken", null);

      this.$router.push({ name: 'SignIn' });
    },
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


.card-expansion {
  margin: 0 175px 10px 175px;
  text-align: center;
}

.md-card {
  width: 250px;
  margin: 4px;
  display: inline-block;
  vertical-align: top;
}

.pagination-row {
  text-align: center;
  margin-bottom: 75px;
}
</style>
