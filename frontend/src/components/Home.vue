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
        <md-autocomplete
            id="news"
            @input="createNews"
            v-model="news"
            :md-options="[]"
            md-layout="box"
            md-dense>
          <label>What's a new?</label>
        </md-autocomplete>
      </md-app-content>
    </md-app>
    <FlashMessage :position="'right top'"></FlashMessage>
  </div>
</template>

<script>
import {debounce} from "@/const";

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
    news: []
  }),
  methods: {
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
      this.$store.commit("changeAnthroponym", this.searchPayload.anthroponym);
      this.$router.push({ name: 'People' })
    }, 1000),
    createNews: debounce(function (){
      if (this.news.length === 0) {
        return
      }

      const payload = {
        topic: "news",
        payload: JSON.stringify({
          content: this.news,
        }),
      };

      this.$wsSend(JSON.stringify(payload))
      this.flashMessage.setStrategy('single');
      this.flashMessage.success({
        title: 'Success',
        message: 'News were published!'
      });
      this.news = []
    }, 1500),
    logOut() {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");

      this.$wsDisconnect();

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

#news {
  margin-left: auto;
  margin-right: auto;
  width: 60%;
}
</style>
