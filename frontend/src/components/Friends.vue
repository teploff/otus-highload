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
            <span class="md-list-item-text">My profile</span>
          </md-list-item>

          <md-list-item @click="followNewsPage">
            <md-icon>fiber_new</md-icon>
            <span class="md-list-item-text">News</span>
            <md-badge v-if="countNewsNotify > 0" class="md-primary" v-bind:md-content="countNewsNotify" />
          </md-list-item>

          <md-list-item @click="followMessengerPage">
            <md-icon>chat</md-icon>
            <span class="md-list-item-text">Messenger</span>
            <md-badge v-if="countMsgNotify > 0" class="md-primary" v-bind:md-content="countMsgNotify" />
          </md-list-item>

          <md-list-item @click="followFriendsPage">
            <md-icon>supervisor_account</md-icon>
            <span class="md-list-item-text">Friends</span>
            <md-badge v-if="countFriendsNotify > 0" class="md-primary" v-bind:md-content="countFriendsNotify" />
          </md-list-item>
        </md-list>
      </md-app-drawer>

      <md-app-content>
        <md-tabs md-alignment="centered">
          <md-tab id="tab-home" md-label="My Friends" md-icon="group" @click="getFriends">
            <md-empty-state
                v-show="friends.length === 0"
                md-rounded
                md-icon="access_time"
                md-label="Nothing in Snoozed"
                md-description="Anything you snooze will go here until it's time for it to return to the inbox.">
            </md-empty-state>
            <div v-show="friends.length !== 0">
              <md-table v-model="friends" md-card @md-selected="selectFriends">
                <md-table-toolbar>
                  <h1 class="md-title" id="friend-requests">Friends</h1>
                </md-table-toolbar>

                <md-table-toolbar id="table-friends-header-toolbar" slot="md-table-alternate-header" slot-scope="{}">
                  <div class="md-toolbar-section-start">{{ getAlternateLabel() }}</div>

                  <div class="md-toolbar-section-end">
                    <md-button class="md-icon-button" @click="removeFriends">
                      <md-icon>delete</md-icon>
                    </md-button>
                  </div>
                </md-table-toolbar>

                <md-table-row slot="md-table-row" slot-scope="{ item }" md-selectable="multiple" md-auto-select>
                  <md-table-cell md-label="Name" md-sort-by="name">{{ item.name }} {{ item.surname }}</md-table-cell>
                  <md-table-cell md-label="Email" md-sort-by="email">{{ item.email }}</md-table-cell>
                  <md-table-cell md-label="Sex" md-sort-by="sex">{{ item.sex }}</md-table-cell>
                  <md-table-cell md-label="Birthday" md-sort-by="birthday">{{ $moment(item.birthday).format('MMMM Do YYYY') }}</md-table-cell>
                  <md-table-cell md-label="City" md-sort-by="city">{{ item.city }}</md-table-cell>
                </md-table-row>
              </md-table>
            </div>
          </md-tab>
          <md-tab id="tab-pages" md-label="Friend requests" md-icon="group_add" @click="getFollowers">
            <md-empty-state
                v-show="followers.length === 0"
                md-rounded
                md-icon="access_time"
                md-label="Nothing in Snoozed"
                md-description="Anything you snooze will go here until it's time for it to return to the inbox.">
            </md-empty-state>

            <div v-show="followers.length !== 0">
              <md-table v-model="followers" md-card @md-selected="selectFollowers">
                <md-table-toolbar>
                  <h1 class="md-title" id="followers-requests">My followers</h1>
                </md-table-toolbar>

                <md-table-toolbar id="table-followers-header-toolbar" slot="md-table-alternate-header" slot-scope="{}">
                  <div class="md-toolbar-section-start">{{ getAlternateLabel() }}</div>

                  <div class="md-toolbar-section-end">
                    <md-button class="md-icon-button" @click="acceptFollowers">
                      <md-icon>add</md-icon>
                    </md-button>
                    <md-button class="md-icon-button" @click="removeFollowers">
                      <md-icon>delete</md-icon>
                    </md-button>
                  </div>
                </md-table-toolbar>

                <md-table-row slot="md-table-row" slot-scope="{ item }" md-selectable="multiple" md-auto-select>
                  <md-table-cell md-label="Name" md-sort-by="name">{{ item.name }} {{ item.surname }}</md-table-cell>
                  <md-table-cell md-label="Email" md-sort-by="email">{{ item.email }}</md-table-cell>
                  <md-table-cell md-label="Sex" md-sort-by="sex">{{ item.sex }}</md-table-cell>
                  <md-table-cell md-label="Birthday" md-sort-by="birthday">{{ $moment(item.birthday).format('MMMM Do YYYY') }}</md-table-cell>
                  <md-table-cell md-label="City" md-sort-by="city">{{ item.city }}</md-table-cell>
                </md-table-row>
              </md-table>
            </div>
          </md-tab>
        </md-tabs>
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
    friends: [],
    selectedFriends: [],
    followers: [],
    selectedFollowers: [],
  }),
  beforeMount() {
    this.getFriends()
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
    searchPeople: debounce(function (){
      this.$store.commit("changeAnthroponym", this.searchPayload.anthroponym);
      this.$router.push({ name: 'People' })
    }, 1000),
    refreshToken() {
      const path = `${apiUrl}/auth/token`;
      const camelcaseKeys = require('camelcase-keys');

      const payload = {
        refresh_token: localStorage.getItem('refreshToken'),
      };
      axios.put(path, payload, {transformResponse: [(data) => {
          return camelcaseKeys(JSON.parse(data), { deep: true })}
        ]})
          .then((response) => {
            this.tokenPair = response.data;

            localStorage.setItem('accessToken', this.tokenPair.accessToken);
            localStorage.setItem('refreshToken', this.tokenPair.refreshToken);

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
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");

      this.$wsDisconnect();

      this.$router.push({ name: 'SignIn' });
    },
    selectFriends(items) {
      this.selectedFriends = items
    },
    selectFollowers(items) {
      this.selectedFollowers = items
    },
    getAlternateLabel () {
    },
    getFriends() {
      const path = `${apiUrl}/social/friends`;
      const camelcaseKeys = require('camelcase-keys');

      headers.Authorization = localStorage.getItem('accessToken')

      axios.get(path, {
        headers: headers,
        transformResponse: [(data) => {
          return camelcaseKeys(JSON.parse(data), { deep: true })}
        ]
      })
          .then((response) => {
            this.friends = response.data.friends;
          })
          .catch((error) => {
            const err = error.response;

            if (err.status === 401) {
              this.refreshToken();
            }

            this.flashMessage.setStrategy('single');
            this.flashMessage.error({
              title: 'Error Message Title',
              message: err.data.message,
              position: 'center',
              icon: '../assets/error.svg',
            });
          });
    },
    getFollowers() {
      const path = `${apiUrl}/social/followers`;
      const camelcaseKeys = require('camelcase-keys');

      headers.Authorization = localStorage.getItem('accessToken')

      axios.get(path, {
        headers: headers,
        transformResponse: [(data) => {
          return camelcaseKeys(JSON.parse(data), { deep: true })}
        ]
      })
          .then((response) => {
            this.followers = response.data.followers;
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
    removeFriends() {
      const path = `${apiUrl}/social/break-friendship`;
      const camelcaseKeys = require('camelcase-keys');

      let friends_id = [];
      for (let i in this.selectedFriends) {
        friends_id.push(this.selectedFriends[i].id)
      }

      const payload = {
        friends_id: friends_id,
      };

      axios.post(path, payload, {headers: headers, transformResponse: [(data) => {
          return camelcaseKeys(JSON.parse(data), { deep: true })}
        ]})
          .then(() => {
            this.friends = this.friends.filter(function(item) {
              return !friends_id.some(function(s) { return s === item.id && s.lines === item.lines })
            });
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
    acceptFollowers() {
      const path = `${apiUrl}/social/confirm-friendship`;
      const camelcaseKeys = require('camelcase-keys');

      let followers_id = [];
      for (let i in this.selectedFollowers) {
        followers_id.push(this.selectedFollowers[i].id)
      }

      const payload = {
        friends_id: followers_id,
      };

      axios.post(path, payload, {headers: headers, transformResponse: [(data) => {
          return camelcaseKeys(JSON.parse(data), { deep: true })}
        ]})
          .then(() => {
            this.followers = this.followers.filter(function(item) {
              return !followers_id.some(function(s) { return s === item.id && s.lines === item.lines })
            });
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
    removeFollowers() {
      const path = `${apiUrl}/social/reject-friendship`;
      const camelcaseKeys = require('camelcase-keys');

      let followers_id = [];
      for (let i in this.selectedFollowers) {
        followers_id.push(this.selectedFollowers[i].id)
      }

      const payload = {
        friends_id: followers_id,
      };

      axios.post(path, payload, {headers: headers, transformResponse: [(data) => {
          return camelcaseKeys(JSON.parse(data), { deep: true })}
        ]})
          .then(() => {
            this.followers = this.followers.filter(function(item) {
              return !followers_id.some(function(s) { return s === item.id && s.lines === item.lines })
            });
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

#table-friends-header-toolbar{
  background-color: #d3e2fb;
}

#table-followers-header-toolbar {
  background-color: #d3e2fb;
}

#friend-requests {
  text-align: center;
}

#followers-requests {
  text-align: center;
}

</style>
