<template>
  <div class="layout">
    <mdb-side-nav-2 :value="true" :data="navigation" push slim :slim-collapsed="collapsed" @toggleSlim="collapsed = $event">
      <div slot="header">
        <div
            class="d-flex align-items-center my-4"
            :class="collapsed ? 'justify-content-center' : 'justify-content-start'"
        >
          <mdb-avatar :width="40" style="flex: 0 0 auto">
            <img
                src="https://mdbootstrap.com/img/Photos/Avatars/avatar-7.jpg"
                class="img-fluid rounded-circle z-depth-1"
            />
          </mdb-avatar>
          <p class="m-t mb-0 ml-4 p-0" style="flex: 0 2 auto" v-show="!collapsed">
            <strong>John Smith<mdb-icon color="success" icon="circle" class="ml-2" size="sm"/></strong>
          </p>
        </div>
        <hr class="w-100" />
      </div>
      <div slot="content" class="mt-5 d-flex justify-content-center">
        <mdb-btn tag="a" gradient="blue" size="sm" class="mx-0" floating :icon="collapsed ? 'chevron-right' : 'chevron-left'" @click="collapsed = !collapsed"></mdb-btn>
      </div>
      <mdb-navbar
          slot="nav"
          tag="div"
          :toggler="false"
          position="top">
        <mdb-navbar-nav right>
          <mdb-form-inline class="ml-auto">
            <mdb-input v-model="searchPayload.anthroponym" class="mr-sm-1" type="text" placeholder="Search people..."/>
            <mdb-btn tag="a" size="sm" gradient="blue" floating @click="searchPeople()" v-show="searchPayload.anthroponym !== ''"><mdb-icon icon="search"/></mdb-btn>
          </mdb-form-inline>
        </mdb-navbar-nav>
        <mdb-navbar-nav class="nav-flex-icons" right>
          <mdb-btn tag="a" size="sm" gradient="blue" floating @click="logOut()"><mdb-icon icon="door-open"/></mdb-btn>
        </mdb-navbar-nav>
      </mdb-navbar>
      <div slot="main">
        <mdb-tabs
            :active="0"
            @activeTab="getActiveTabIndex"
            tabs
            card
            class="mb-5"
            justify
            color="info"
            :links="[
                { text: 'Friends', icon: 'user-friends', slot: 'friends-slot' },
                { text: 'Followers', icon: 'bullhorn', slot: 'followers-slot'}]"
        >
          <template slot="friends-slot">
            <mdb-container>
              <mdb-row v-for="(row, i) in friends" v-bind:key="i">
                <mdb-col sm="4" v-for="friend in row" v-bind:key="friend.id">
                  <mdb-card testimonial>
                    <mdb-card-up gradient="blue" class="lighten-1"></mdb-card-up>
                    <mdb-card-avatar color="white" class="mx-auto">
                      <img v-if="friend.sex==='female'" src="../assets/girl.png" class="rounded-circle">
                      <img v-else src="../assets/boy.png" class="rounded-circle">
                    </mdb-card-avatar>
                    <mdb-card-body>
                      <mdb-card-title>{{ friend.name }} {{ friend.surname }}</mdb-card-title>
                      <hr />
                      <p>
                        <mdb-icon icon="quote-left" /> Lorem ipsum dolor sit amet, consectetur adipisicing elit. Eos,
                        adipisci</p>
                      <hr />
                      <mdb-row>
                        <mdb-col>
                          <mdb-btn @click="splitUpFriendship(friend)" rounded color="red">Remove</mdb-btn>
                        </mdb-col>
                      </mdb-row>
                    </mdb-card-body>
                  </mdb-card>
                </mdb-col>
              </mdb-row>
            </mdb-container>
            <mdb-pagination class="justify-content-center" circle>
              <mdb-page-item disabled>First</mdb-page-item>
              <mdb-page-nav prev disabled></mdb-page-nav>
              <mdb-page-item active>1</mdb-page-item>
              <mdb-page-item>2</mdb-page-item>
              <mdb-page-item>3</mdb-page-item>
              <mdb-page-item>4</mdb-page-item>
              <mdb-page-item>5</mdb-page-item>
              <mdb-page-nav next></mdb-page-nav>
              <mdb-page-item>Last</mdb-page-item>
            </mdb-pagination>
          </template>
          <template slot="followers-slot">
            <mdb-container>
              <mdb-row v-for="(row, i) in followers" v-bind:key="i">
                <mdb-col sm="4" v-for="follower in row" v-bind:key="follower.id">
                  <mdb-card testimonial>
                    <mdb-card-up gradient="blue" class="lighten-1"></mdb-card-up>
                    <mdb-card-avatar color="white" class="mx-auto">
                      <img v-if="follower.sex==='female'" src="../assets/girl.png" class="rounded-circle">
                      <img v-else src="../assets/boy.png" class="rounded-circle">
                    </mdb-card-avatar>
                    <mdb-card-body>
                      <mdb-card-title>{{ follower.name }} {{ follower.surname }}</mdb-card-title>
                      <hr />
                      <p>
                        <mdb-icon icon="quote-left" /> Lorem ipsum dolor sit amet, consectetur adipisicing elit. Eos,
                        adipisci</p>
                      <hr />
                      <mdb-row>
                        <mdb-col class="col-6">
                          <mdb-btn @click="acceptFollower(follower)" rounded color="primary">Add</mdb-btn>
                        </mdb-col>
                        <mdb-col class="col-4">
                          <mdb-btn @click="refuseFollower(follower)" rounded color="red">Remove</mdb-btn>
                        </mdb-col>
                      </mdb-row>
                    </mdb-card-body>
                  </mdb-card>
                </mdb-col>
              </mdb-row>
              <mdb-pagination class="justify-content-center" circle>
                <mdb-page-item disabled>First</mdb-page-item>
                <mdb-page-nav prev disabled></mdb-page-nav>
                <mdb-page-item active>1</mdb-page-item>
                <mdb-page-item>2</mdb-page-item>
                <mdb-page-item>3</mdb-page-item>
                <mdb-page-item>4</mdb-page-item>
                <mdb-page-item>5</mdb-page-item>
                <mdb-page-nav next></mdb-page-nav>
                <mdb-page-item>Last</mdb-page-item>
              </mdb-pagination>
            </mdb-container>
          </template>
        </mdb-tabs>
      </div>
    </mdb-side-nav-2>
  </div>
</template>

<script>
import {
  mdbNavbar,
  mdbNavbarNav,
  mdbSideNav2,
  mdbAvatar,
  mdbBtn,
  mdbIcon,
  waves,
  mdbFormInline,
  mdbInput,
  mdbTabs,
  mdbRow,
  mdbCol,
  mdbCard,
  mdbCardBody,
  mdbCardTitle,
  mdbCardUp,
  mdbCardAvatar,
  mdbContainer,
  mdbPagination,
  mdbPageItem,
  mdbPageNav,
} from "mdbvue";

import router from "@/router";
import {confirmFriendship, getFollowers, getFriends, rejectFriendship, splitUpFriendship} from "@/api/social.api";

export default {
  components: {
    mdbNavbar,
    mdbNavbarNav,
    mdbSideNav2,
    mdbAvatar,
    mdbBtn,
    mdbIcon,
    mdbFormInline,
    mdbInput,
    mdbTabs,
    mdbRow,
    mdbCol,
    mdbCard,
    mdbCardBody,
    mdbCardTitle,
    mdbCardUp,
    mdbCardAvatar,
    mdbContainer,
    mdbPagination,
    mdbPageItem,
    mdbPageNav,
  },
  name: "Friends",
  async beforeMount() {
    try {
      const response = await getFriends()
      const friends = response.data.friends

      this.friends = [[]]

      let j = 0;
      for (let i = 0; i < friends.length; i++) {
        if ((i !== 0) && (i % 4 === 0)) {
          j++;
        }

        this.friends[j].push(friends[i])
      }
    } catch (error) {
      this.$notify.error({message: error.response.data.message, position: 'top right', timeOut: 5000});
    }
  },
  data: () => ({
    show: true,
    collapsed: false,
    navigation: [
      {
        name: "My profile",
        icon: "address-card",
        href: router.resolve({name: 'Home'}).href
      },
      {
        name: "News",
        icon: "newspaper",
        href: router.resolve({name: 'News'}).href
      },
      {
        name: "Messenger",
        icon: "comments",
        href: router.resolve({name: 'Messenger'}).href
      },
      {
        name: "Friends",
        icon: "user-friends",
        href: router.resolve({name: 'Friends'}).href
      }
    ],
    searchPayload: {
      anthroponym: '',
    },
    friends: [[]],
    followers: [[]],
  }),
  methods: {
    log(text) {
      console.log(text)
    },
    searchPeople() {
      this.$store.commit("changeAnthroponym", this.searchPayload.anthroponym);
      this.$router.push({ name: 'People' })
    },
    async getActiveTabIndex(index) {
      if (index === 0) {
        try {
          const response = await getFriends()
          const friends = response.data.friends

          this.friends = [[]]

          let j = 0;
          for (let i = 0; i < friends.length; i++) {
            if ((i !== 0) && (i%4 === 0)) {
              j++;
            }

            this.friends[j].push(friends[i])
          }
        } catch (error) {
          this.$notify.error({message: error.response.data.message, position: 'top right', timeOut: 5000});
        }
      } else {
        try {
          const response = await getFollowers()
          const followers = response.data.followers

          this.followers = [[]]

          let j = 0;
          for (let i = 0; i < followers.length; i++) {
            if ((i !== 0) && (i%4 === 0)) {
              j++;
            }

            this.followers[j].push(followers[i])
          }
        } catch (error) {
          this.$notify.error({message: error.response.data.message, position: 'top right', timeOut: 5000});
        }
      }
    },
    async acceptFollower(user) {
      try {
        await confirmFriendship([user.id])
      } catch (error) {
        this.$notify.error({message: error.response.data.message, position: 'top right', timeOut: 5000});
      }
    },
    async refuseFollower(user) {
      try {
        await rejectFriendship([user.id])
      } catch (error) {
        this.$notify.error({message: error.response.data.message, position: 'top right', timeOut: 5000});
      }
    },
    async splitUpFriendship(user) {
      try {
        await splitUpFriendship([user.id])
      } catch (error) {
        this.$notify.error({message: error.response.data.message, position: 'top right', timeOut: 5000});
      }
    },
    logOut() {
      localStorage.removeItem("accessToken");
      localStorage.removeItem("refreshToken");

      this.$router.push({name: 'SignIn'});
    },
  },
  mixins: [waves]
}
</script>

<style scoped>
.navbar i {
  cursor: pointer;
  color: white;
}
</style>
