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
          <mdb-row class="justify-content-center">
            <mdb-col sm="2">
              <mdb-flipping-card
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
                    <img
                        src="https://mdbootstrap.com/img/Photos/Avatars/img%20(20).jpg"
                        class="rounded-circle"
                    />
                  </mdb-avatar>
                  <mdb-card-body>
                    <h4 class="font-weight-bold mb-3">Marie Johnson</h4>
                    <p class="font-weight-bold blue-text">Web developer</p>
                    <a class="rotate-btn" @click="flipped=true">
                      <mdb-icon class="pr-2" icon="redo" />Click here to rotate
                    </a>
                  </mdb-card-body>
                </mdb-card>
                <mdb-card class="face back" style="height: 416px;">
                  <mdb-card-body>
                    <h4 class="font-weight-bold">About me</h4>
                    <hr />
                    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Maxime quae, dolores dicta. Blanditiis rem amet repellat, dolores nihil quae in mollitia asperiores ut rerum repellendus, voluptatum eum, officia laudantium quaerat?</p>
                    <hr />
                    <ul class="list-inline py-2">
                      <li class="list-inline-item">
                        <a class="p-2 fa-lg fb-ic">
                          <mdb-icon fab icon="facebook" />
                        </a>
                      </li>
                      <li class="list-inline-item">
                        <a class="p-2 fa-lg tw-ic">
                          <mdb-icon fab icon="twitter" />
                        </a>
                      </li>
                      <li class="list-inline-item">
                        <a class="p-2 fa-lg gplus-ic">
                          <mdb-icon fab icon="google-plus" />
                        </a>
                      </li>
                      <li class="list-inline-item">
                        <a class="p-2 fa-lg li-ic">
                          <mdb-icon fab icon="linkedin" />
                        </a>
                      </li>
                    </ul>
                    <a class="rotate-btn" @click="flipped=false">
                      <mdb-icon class="pr-2" icon="undo" />Click here to rotate back
                    </a>
                  </mdb-card-body>
                </mdb-card>
              </mdb-flipping-card>
            </mdb-col>
            <mdb-col sm="2">
              <mdb-flipping-card
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
                    <img
                        src="https://mdbootstrap.com/img/Photos/Avatars/img%20(20).jpg"
                        class="rounded-circle"
                    />
                  </mdb-avatar>
                  <mdb-card-body>
                    <h4 class="font-weight-bold mb-3">Marie Johnson</h4>
                    <p class="font-weight-bold blue-text">Web developer</p>
                    <a class="rotate-btn" @click="flipped=true">
                      <mdb-icon class="pr-2" icon="redo" />Click here to rotate
                    </a>
                  </mdb-card-body>
                </mdb-card>
                <mdb-card class="face back" style="height: 416px;">
                  <mdb-card-body>
                    <h4 class="font-weight-bold">About me</h4>
                    <hr />
                    <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Maxime quae, dolores dicta. Blanditiis rem amet repellat, dolores nihil quae in mollitia asperiores ut rerum repellendus, voluptatum eum, officia laudantium quaerat?</p>
                    <hr />
                    <ul class="list-inline py-2">
                      <li class="list-inline-item">
                        <a class="p-2 fa-lg fb-ic">
                          <mdb-icon fab icon="facebook" />
                        </a>
                      </li>
                      <li class="list-inline-item">
                        <a class="p-2 fa-lg tw-ic">
                          <mdb-icon fab icon="twitter" />
                        </a>
                      </li>
                      <li class="list-inline-item">
                        <a class="p-2 fa-lg gplus-ic">
                          <mdb-icon fab icon="google-plus" />
                        </a>
                      </li>
                      <li class="list-inline-item">
                        <a class="p-2 fa-lg li-ic">
                          <mdb-icon fab icon="linkedin" />
                        </a>
                      </li>
                    </ul>
                    <a class="rotate-btn" @click="flipped=false">
                      <mdb-icon class="pr-2" icon="undo" />Click here to rotate back
                    </a>
                  </mdb-card-body>
                </mdb-card>
              </mdb-flipping-card>
            </mdb-col>
            <mdb-col sm="2">
              <mdb-flipping-card
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
                  <img
                      src="https://mdbootstrap.com/img/Photos/Avatars/img%20(20).jpg"
                      class="rounded-circle"
                  />
                </mdb-avatar>
                <mdb-card-body>
                  <h4 class="font-weight-bold mb-3">Marie Johnson</h4>
                  <p class="font-weight-bold blue-text">Web developer</p>
                  <a class="rotate-btn" @click="flipped=true">
                    <mdb-icon class="pr-2" icon="redo" />Click here to rotate
                  </a>
                </mdb-card-body>
              </mdb-card>
              <mdb-card class="face back" style="height: 416px;">
                <mdb-card-body>
                  <h4 class="font-weight-bold">About me</h4>
                  <hr />
                  <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Maxime quae, dolores dicta. Blanditiis rem amet repellat, dolores nihil quae in mollitia asperiores ut rerum repellendus, voluptatum eum, officia laudantium quaerat?</p>
                  <hr />
                  <ul class="list-inline py-2">
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg fb-ic">
                        <mdb-icon fab icon="facebook" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg tw-ic">
                        <mdb-icon fab icon="twitter" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg gplus-ic">
                        <mdb-icon fab icon="google-plus" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg li-ic">
                        <mdb-icon fab icon="linkedin" />
                      </a>
                    </li>
                  </ul>
                  <a class="rotate-btn" @click="flipped=false">
                    <mdb-icon class="pr-2" icon="undo" />Click here to rotate back
                  </a>
                </mdb-card-body>
              </mdb-card>
            </mdb-flipping-card>
            </mdb-col>
            <mdb-col sm="2">
              <mdb-flipping-card
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
                  <img
                      src="https://mdbootstrap.com/img/Photos/Avatars/img%20(20).jpg"
                      class="rounded-circle"
                  />
                </mdb-avatar>
                <mdb-card-body>
                  <h4 class="font-weight-bold mb-3">Marie Johnson</h4>
                  <p class="font-weight-bold blue-text">Web developer</p>
                  <a class="rotate-btn" @click="flipped=true">
                    <mdb-icon class="pr-2" icon="redo" />Click here to rotate
                  </a>
                </mdb-card-body>
              </mdb-card>
              <mdb-card class="face back" style="height: 416px;">
                <mdb-card-body>
                  <h4 class="font-weight-bold">About me</h4>
                  <hr />
                  <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Maxime quae, dolores dicta. Blanditiis rem amet repellat, dolores nihil quae in mollitia asperiores ut rerum repellendus, voluptatum eum, officia laudantium quaerat?</p>
                  <hr />
                  <ul class="list-inline py-2">
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg fb-ic">
                        <mdb-icon fab icon="facebook" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg tw-ic">
                        <mdb-icon fab icon="twitter" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg gplus-ic">
                        <mdb-icon fab icon="google-plus" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg li-ic">
                        <mdb-icon fab icon="linkedin" />
                      </a>
                    </li>
                  </ul>
                  <a class="rotate-btn" @click="flipped=false">
                    <mdb-icon class="pr-2" icon="undo" />Click here to rotate back
                  </a>
                </mdb-card-body>
              </mdb-card>
            </mdb-flipping-card>
            </mdb-col>
            <mdb-col sm="2">
              <mdb-flipping-card
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
                  <img
                      src="https://mdbootstrap.com/img/Photos/Avatars/img%20(20).jpg"
                      class="rounded-circle"
                  />
                </mdb-avatar>
                <mdb-card-body>
                  <h4 class="font-weight-bold mb-3">Marie Johnson</h4>
                  <p class="font-weight-bold blue-text">Web developer</p>
                  <a class="rotate-btn" @click="flipped=true">
                    <mdb-icon class="pr-2" icon="redo" />Click here to rotate
                  </a>
                </mdb-card-body>
              </mdb-card>
              <mdb-card class="face back" style="height: 416px;">
                <mdb-card-body>
                  <h4 class="font-weight-bold">About me</h4>
                  <hr />
                  <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Maxime quae, dolores dicta. Blanditiis rem amet repellat, dolores nihil quae in mollitia asperiores ut rerum repellendus, voluptatum eum, officia laudantium quaerat?</p>
                  <hr />
                  <ul class="list-inline py-2">
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg fb-ic">
                        <mdb-icon fab icon="facebook" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg tw-ic">
                        <mdb-icon fab icon="twitter" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg gplus-ic">
                        <mdb-icon fab icon="google-plus" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg li-ic">
                        <mdb-icon fab icon="linkedin" />
                      </a>
                    </li>
                  </ul>
                  <a class="rotate-btn" @click="flipped=false">
                    <mdb-icon class="pr-2" icon="undo" />Click here to rotate back
                  </a>
                </mdb-card-body>
              </mdb-card>
            </mdb-flipping-card>
            </mdb-col>
          </mdb-row>
          <mdb-row class="justify-content-center">
            <mdb-col sm="2">
              <mdb-flipping-card
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
                  <img
                      src="https://mdbootstrap.com/img/Photos/Avatars/img%20(20).jpg"
                      class="rounded-circle"
                  />
                </mdb-avatar>
                <mdb-card-body>
                  <h4 class="font-weight-bold mb-3">Marie Johnson</h4>
                  <p class="font-weight-bold blue-text">Web developer</p>
                  <a class="rotate-btn" @click="flipped=true">
                    <mdb-icon class="pr-2" icon="redo" />Click here to rotate
                  </a>
                </mdb-card-body>
              </mdb-card>
              <mdb-card class="face back" style="height: 416px;">
                <mdb-card-body>
                  <h4 class="font-weight-bold">About me</h4>
                  <hr />
                  <p>Lorem ipsum dolor sit amet, consectetur adipisicing elit. Maxime quae, dolores dicta. Blanditiis rem amet repellat, dolores nihil quae in mollitia asperiores ut rerum repellendus, voluptatum eum, officia laudantium quaerat?</p>
                  <hr />
                  <ul class="list-inline py-2">
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg fb-ic">
                        <mdb-icon fab icon="facebook" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg tw-ic">
                        <mdb-icon fab icon="twitter" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg gplus-ic">
                        <mdb-icon fab icon="google-plus" />
                      </a>
                    </li>
                    <li class="list-inline-item">
                      <a class="p-2 fa-lg li-ic">
                        <mdb-icon fab icon="linkedin" />
                      </a>
                    </li>
                  </ul>
                  <a class="rotate-btn" @click="flipped=false">
                    <mdb-icon class="pr-2" icon="undo" />Click here to rotate back
                  </a>
                </mdb-card-body>
              </mdb-card>
            </mdb-flipping-card>
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
  mdbRow,
  mdbCol,
} from "mdbvue";

import router from "@/router";
import store from "@/store"

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
    mdbRow,
    mdbCol,
  },
name: "People",
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
      anthroponym: store.getters.searchAnthroponym,
    },
    flipped: false
  }),
  methods: {
    searchPeople() {
      this.$store.commit("changeAnthroponym", this.searchPayload.anthroponym);
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