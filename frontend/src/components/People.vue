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
        <mdb-flipping-card v-for="card in cards.people" v-bind:key="card.id"
            :flipped="flipped"
            innerClass="text-center h-100 w-100"
            style="max-width: 22rem; height: 416px;"
        >
          <mdb-card class="face front" style="height: 416px;">
            <mdb-card-up>
              <img
                  class="card-image-top"
                  src="https://mdbootstrap.com/img/Photos/Others/photo7.jpg"
                  alt="a photo of a house facade"
              />
            </mdb-card-up>
            <mdb-avatar class="mx-auto white" circle>
              <img v-if="card.sex === 'male'" src="../assets/boy.png" class="rounded-circle"/>
              <img v-else src="../assets/girl.png" class="rounded-circle"/>
            </mdb-avatar>
            <mdb-card-body>
              <h4 class="font-weight-bold mb-3">{{card.name}} {{card.surname}}</h4>
              <p class="font-weight-bold blue-text">{{card.email}}</p>
              <p><b>Birthday:</b> {{ moment(card.birthday).format('Do MMMM YYYY') }}</p>
              <a class="rotate-btn" @click="flipped=true">
                <mdb-icon class="pr-2" icon="redo" />Learn more
              </a>
            </mdb-card-body>
          </mdb-card>
          <mdb-card class="face back" style="height: 416px;">
            <mdb-card-body>
              <h4 class="font-weight-bold">About me</h4>
              <hr />
              <p>Hi there! I'm {{card.name}} {{card.surname}}.</p>
              <p>I'm from {{card.city}}.</p>
              <p>My interests: {{card.interests}}</p>
              <hr />
              <ul class="list-inline py-2">
                <li class="list-inline-item">
                  <mdb-tooltip material trigger="hover" :options="{placement: 'left'}">
                    <span slot="tip">Make friendship</span>
                      <mdb-btn slot="reference" v-if="card.friendshipStatus === 'noname'" tag="a" gradient="blue" floating><mdb-icon icon="plus"/></mdb-btn>
                  </mdb-tooltip>
                  <mdb-tooltip v-if="card.friendshipStatus !== 'noname'" material trigger="hover" :options="{placement: 'left'}">
                    <span slot="tip">Your fiend</span>
                    <mdb-btn slot="reference" tag="a" gradient="green" disabled floating><mdb-icon icon="check"/></mdb-btn>
                  </mdb-tooltip>
                </li>
                <li class="list-inline-item">
                  <mdb-tooltip v-if="card.friendshipStatus !== 'noname'" material trigger="hover" :options="{placement: 'right'}">
                    <span slot="tip">Start chatting</span>
                    <mdb-btn slot="reference" tag="a" gradient="peach" floating><mdb-icon icon="comment"/></mdb-btn>
                  </mdb-tooltip>
                </li>
              </ul>
              <a class="rotate-btn" @click="flipped=false">
                <mdb-icon class="pr-2" icon="undo" />Back to preview
              </a>
            </mdb-card-body>
          </mdb-card>
        </mdb-flipping-card>
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
  mdbInput,
  mdbFormInline,
  mdbCard,
  mdbCardBody,
  mdbCardUp,
  mdbFlippingCard,
  mdbPagination,
  mdbPageItem,
  mdbPageNav,
  mdbTooltip,
} from "mdbvue";

import router from "@/router";
import store from "@/store"
import {searchByAnthroponym} from "@/api/social.api";

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
    mdbCard,
    mdbCardBody,
    mdbCardUp,
    mdbFlippingCard,
    mdbPagination,
    mdbPageItem,
    mdbPageNav,
    mdbTooltip,
  },
name: "People",
  data: () => ({
    moment: require('moment'),
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
      anthroponym: store.getters.searchAnthroponym,
      limit: 10,
      offset: 0,
    },
    flipped: false,
    cards: {
      people: [],
      count: 0,
    }
  }),
  methods: {
    searchPeople() {
      this.$store.commit("changeAnthroponym", this.searchPayload.anthroponym);

      this.getPeopleByAnthroponym();
    },
    async getPeopleByAnthroponym() {
      try {
        const response = await searchByAnthroponym(this.searchPayload)

        this.cards.people = response.data.users
        this.cards.count = response.data.count
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
  created() {
    if (this.$store.getters.searchAnthroponym !== null && this.$store.getters.searchAnthroponym !== '') {
      this.getPeopleByAnthroponym();
    }
  },
  beforeDestroy() {
    this.$store.commit("changeAnthroponym", null);
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