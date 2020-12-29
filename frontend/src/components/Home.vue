<template>
  <div class="page-container">
    <md-app>
      <md-app-toolbar class="md-primary">
        <span class="md-title">My Title</span>
      </md-app-toolbar>

      <md-app-drawer md-permanent="clipped">
        <md-list>
          <md-list-item @click="followHomePage">
            <md-icon>assignment_ind</md-icon>
            <span class="md-list-item-text">Моя страница</span>
          </md-list-item>

          <md-list-item @click="followNewsPage">
            <md-icon>fiber_new</md-icon>
            <span class="md-list-item-text">Новости</span>
          </md-list-item>

          <md-list-item @click="followMessengerPage">
            <md-icon>chat</md-icon>
            <span class="md-list-item-text">Мессенджер</span>
          </md-list-item>

          <md-list-item @click="followFriendsPage">
            <md-icon>supervisor_account</md-icon>
            <span class="md-list-item-text">Друзья</span>
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
export default {
  name: 'Home',
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
  },
};
</script>

<style scoped>
.md-app {
  min-height: 350px;
  border: 1px solid rgba(#000, .12);
}

.md-drawer {
 width: 230px;
 max-width: calc(100vw - 125px);
}
</style>
