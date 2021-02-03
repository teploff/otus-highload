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
            tabs
            card
            class="mb-5"
            justify
            color="info"
            :links="[
                { text: 'Friends', icon: 'user-friends', slot: 'friends-slot' },
                { text: 'Followers', icon: 'bullhorn', slot: 'friends-slot' }]"
        >
          <template slot="friends-slot">
            <mdb-container>
              <mdb-row>
                <mdb-col sm="4">
                  <mdb-card testimonial>
                    <mdb-card-up gradient="blue" class="lighten-1"></mdb-card-up>
                    <mdb-card-avatar color="white" class="mx-auto"><img src="../assets/girl.png" class="rounded-circle"></mdb-card-avatar>
                    <mdb-card-body>
                      <mdb-card-title>Anna Doe</mdb-card-title>
                      <hr />
                      <p>
                        <mdb-icon icon="quote-left" /> Lorem ipsum dolor sit amet, consectetur adipisicing elit. Eos,
                        adipisci</p>
                      <hr />
                      <mdb-row>
                        <mdb-col class="col-6">
                          <mdb-btn rounded color="primary">Add</mdb-btn>
                        </mdb-col>
                        <mdb-col class="col-4">
                          <mdb-btn rounded color="red">Remove</mdb-btn>
                        </mdb-col>
                      </mdb-row>
                    </mdb-card-body>
                  </mdb-card>
                </mdb-col>
                <mdb-col sm="4">
                  <mdb-card testimonial>
                    <mdb-card-up gradient="blue" class="lighten-1"></mdb-card-up>
                    <mdb-card-avatar color="white" class="mx-auto"><img src="../assets/boy.png" class="rounded-circle"></mdb-card-avatar>
                    <mdb-card-body>
                      <mdb-card-title>Martha Smith</mdb-card-title>
                      <hr />
                      <p>
                        <mdb-icon icon="quote-left" /> Lorem ipsum dolor sit amet, consectetur adipisicing elit. Eos,
                        adipisci</p>
                      <hr />
                      <mdb-row>
                        <mdb-col class="col-6">
                          <mdb-btn rounded color="primary">Add</mdb-btn>
                        </mdb-col>
                        <mdb-col class="col-4">
                          <mdb-btn rounded color="red">Remove</mdb-btn>
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
              <mdb-row>
                <mdb-col sm="4">
                  <mdb-card testimonial>
                    <mdb-card-up gradient="blue" class="lighten-1"></mdb-card-up>
                    <mdb-card-avatar color="white" class="mx-auto"><img src="../assets/girl.png" class="rounded-circle"></mdb-card-avatar>
                    <mdb-card-body>
                      <mdb-card-title>Anna Doe</mdb-card-title>
                      <hr />
                      <p>
                        <mdb-icon icon="quote-left" /> Lorem ipsum dolor sit amet, consectetur adipisicing elit. Eos,
                        adipisci</p>
                      <hr />
                      <mdb-row>
                        <mdb-col class="col-6">
                          <mdb-btn rounded color="primary">Add</mdb-btn>
                        </mdb-col>
                        <mdb-col class="col-4">
                          <mdb-btn rounded color="red">Remove</mdb-btn>
                        </mdb-col>
                      </mdb-row>
                    </mdb-card-body>
                  </mdb-card>
                </mdb-col>
                <mdb-col sm="4">
                  <mdb-card testimonial>
                    <mdb-card-up gradient="blue" class="lighten-1"></mdb-card-up>
                    <mdb-card-avatar color="white" class="mx-auto"><img src="../assets/boy.png" class="rounded-circle"></mdb-card-avatar>
                    <mdb-card-body>
                      <mdb-card-title>Martha Smith</mdb-card-title>
                      <hr />
                      <p>
                        <mdb-icon icon="quote-left" /> Lorem ipsum dolor sit amet, consectetur adipisicing elit. Eos,
                        adipisci</p>
                      <hr />
                      <mdb-row>
                        <mdb-col class="col-6">
                          <mdb-btn rounded color="primary">Add</mdb-btn>
                        </mdb-col>
                        <mdb-col class="col-4">
                          <mdb-btn rounded color="red">Remove</mdb-btn>
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
    }
  }),
  methods: {
    searchPeople() {
      this.$store.commit("changeAnthroponym", this.searchPayload.anthroponym);
      this.$router.push({ name: 'People' })
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
