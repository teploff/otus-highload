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
              v-model="selectedEmployee"
              :md-options="people"
              md-layout="box">
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
        <div v-show="cards.count === 0">
          <md-empty-state
              v-show="cards.count === 0"
              md-rounded
              md-icon="star"
              md-label="You are the first!"
              md-description="Congratulation you are the first user in this portal!
        While there aren't other people but they will appear soon"
              style="width: 600px; height: 600px; position: center">
          </md-empty-state>
        </div>
        <md-table v-show="cards.count !== 0">
          <md-table-row>
            <div class="card-expansion">
              <md-card v-for="card in cards.questionnaires" v-bind:key="card.email">
                <md-card-media>
                  <img v-if="card.sex === 'male'" src="../assets/boy.png" alt="People">
                  <img v-else src="../assets/girl.png" alt="People">
                </md-card-media>

                <md-card-header>
                  <div class="md-title">{{ card.name }} {{ card.surname }}</div>
                  <div class="md-subhead">{{ card.email }}</div>
                </md-card-header>

                <md-card-expand>
                  <md-card-actions md-alignment="right">
                    <md-card-expand-trigger>
                      <md-button class="learn-more-button" style="color: #337ab7">Learn more</md-button>
                    </md-card-expand-trigger>
                  </md-card-actions>

                  <md-card-expand-content>
                    <md-card-content>
                      <p style="text-align: left"> <b>Sex:</b> {{ card.sex }} </p>
                      <p style="text-align: left">
                        <b>Birthday:</b> {{ new Date(card.birthday) | dateFormat('DD.MM.YYYY') }}
                      </p>
                      <p style="text-align: left"> <b>City:</b> {{ card.city }} </p>
                      <p style="text-align: left"> <b>Interests:</b> {{ card.interests }} </p>
                    </md-card-content>
                  </md-card-expand-content>
                </md-card-expand>
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
  </div>
</template>

<script>
import {apiUrl, headers} from "@/const";
import axios from "axios";

export default {
  name: 'People',
  data: () => ({
    menuVisible: false,
    countNewsNotify: 0,
    countMsgNotify: 0,
    countFriendsNotify: 0,
    selectedEmployee: null,
    people: [],
    searchPayload: {
      anthroponym: null,
      limit: 10,
      offset: 0,
    },
    cards: {
      questionnaires: null,
      count: 0,
    },
    countCardsInWindow: 10,
    page: 1,
  }),
  created: function () {
    const path = `${apiUrl}/profile/search/anthroponym`;
    console.log("hehe")

    this.searchPayload.anthroponym = this.$store.getters.searchAnthroponym
    headers.Authorization = localStorage.getItem('access_token');

    console.log("hehe")

    axios.get(path, {
      headers: headers,
      params: this.searchPayload
    })
        .then((response) => {
          this.cards = JSON.parse(JSON.stringify(response.data));
        })
        .catch((error) => {
          const err = JSON.parse(JSON.stringify(error.response));
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
  beforeDestroy() {
    this.$store.commit("changeAnthroponym", null);
  },
  methods: {
    followHomePage() {
      this.$router.push({ name: 'Home' });
    },
    followNewsPage() {
      this.$router.push({ name: 'News' });
    },
    followMessengerPage() {
      this.$router.push({ name: 'Messenger' });
    },
    followFriendsPage() {
      this.$router.push({ name: 'Friends' }).catch(() => {});
    },
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
    paginatorClick(pageNum) {
      const path = `${apiUrl}/profile/search/anthroponym`;
      headers.Authorization = localStorage.getItem('access_token');

      this.payload.offset = (pageNum - 1) * this.countCardsInWindow;

      axios.get(path, {
        headers: headers,
        params: this.searchPayload
      })
          .then((response) => {
            this.cards = JSON.parse(JSON.stringify(response.data));
          })
          .catch((error) => {
            const err = JSON.parse(JSON.stringify(error.response));
            if (err.status === 401) {
              this.refreshToken();
            }
          });
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

.md-empty-state {
  max-width: 600px;
}

.learn-more-button{
  font-weight: bolder;
}

</style>
