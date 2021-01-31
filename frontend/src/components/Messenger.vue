<template>
  <div class="layout">
    <mdb-side-nav-2
        :value="true"
        :data="navigation"
        push
        slim
        :slim-collapsed="collapsed"
        @toggleSlim="collapsed = $event"
    >
      <div slot="header">
        <div
            class="d-flex align-items-center my-4"
            :class="
            collapsed ? 'justify-content-center' : 'justify-content-start'
          "
        >
          <mdb-avatar :width="40" style="flex: 0 0 auto">
            <img
                src="https://mdbootstrap.com/img/Photos/Avatars/avatar-7.jpg"
                class="img-fluid rounded-circle z-depth-1"
            />
          </mdb-avatar>
          <p
              class="m-t mb-0 ml-4 p-0"
              style="flex: 0 2 auto"
              v-show="!collapsed"
          >
            <strong
            >John Smith
              <mdb-icon color="success" icon="circle" class="ml-2" size="sm"
              /></strong>
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
        <mdb-container>
          <mdb-chat-room color="rare-wind-gradient">
            <mdb-row class="px-lg-2 px-2">
              <mdb-col md="6" xl="4" class="px-0 pt-3">
                <mdb-chat-list
                    scroll
                    :height="500"
                    scrollbarClass="scrollbar-light-blue"
                    :data="gradientChat"
                    @click="changeActiveChat"
                ></mdb-chat-list>
              </mdb-col>
              <mdb-col md="6" xl="8" class="pl-md-3 px-lg-auto px-0 pt-3">
                <mdb-chat
                    :loggedUserId="gradientChat[activeChat].loggedUserId"
                    :chat="gradientChat[activeChat].chat"
                    outline="purple"
                    :avatarWidth="50"
                    scroll
                    scrollbarClass="scrollbar-light-blue"
                    @send="sendMessage($event, gradientChat[activeChat])"
                ></mdb-chat>
              </mdb-col>
            </mdb-row>
          </mdb-chat-room>
        </mdb-container>
      </div>
    </mdb-side-nav-2>
  </div>
</template>

<script>
import {
  mdbChat,
  mdbChatList,
  mdbChatRoom,
  mdbCol,
  mdbContainer,
  mdbRow,
  mdbNavbar,
  mdbNavbarNav,
  mdbSideNav2,
  mdbAvatar,
  mdbBtn,
  mdbIcon,
  waves,
  mdbFormInline,
  mdbInput
} from "mdbvue";
import router from "@/router";

export default {
  components: {
    mdbRow,
    mdbCol,
    mdbChat,
    mdbChatList,
    mdbChatRoom,
    mdbContainer,
    mdbNavbar,
    mdbNavbarNav,
    mdbSideNav2,
    mdbAvatar,
    mdbBtn,
    mdbIcon,
    mdbFormInline,
    mdbInput
  },
  name: "Messenger",
  data: () => ({
    gradientChat: [
      {
        loggedUserId: 1,
        active: true,
        lastId: 4,
        chat: [
          {
            id: 0,
            name: "Brad Pitt",
            online: true,
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-6.jpg",
            messages: [
              {
                id: 0,
                date: "2019-04-21 15:00:09",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              },
              {
                id: 1,
                date: "2019-06-26 11:00",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              },
              {
                id: 4,
                date: "2019-07-01 09:20",
                content:
                    "Can you pop out and buy lemons?"
              }
            ]
          },
          {
            id: 1,
            name: "Lara Croft",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-5.jpg",
            messages: [
              {
                id: 2,
                date: "2019-06-26 11:15",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              },
              {
                id: 3,
                date: "2019-06-26 11:16",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              }
            ]
          }
        ]
      },
      {
        id: 1,
        loggedUserId: 1,
        active: false,
        lastId: 2,
        chat: [
          {
            id: 1,
            name: "Lara Croft",
            online: true,
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-5.jpg",
            messages: [
              {
                id: 0,
                date: "2019-04-21 15:00:09",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              },
              {
                id: 1,
                date: "2019-06-26 11:00",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              }
            ]
          },
          {
            id: 2,
            name: "Ashley Olsen",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-2.jpg",
            messages: [
              {
                id: 2,
                date: "2019-06-26 11:15",
                content: "Remember to bring me oranges",
                unread: true
              }
            ]
          }
        ]
      },
      {
        id: 2,
        loggedUserId: 1,
        active: false,
        lastId: 2,
        chat: [
          {
            id: 1,
            name: "Lara Croft",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-5.jpg",
            messages: [
              {
                id: 0,
                date: "2019-04-21 15:00:09",
                content:
                    "Sed ut perspiciatis unde omnis iste natus remque laudantium."
              },
              {
                id: 1,
                date: "2019-06-26 11:00",
                content: "Sed ut doloremque laudantium."
              }
            ]
          },
          {
            id: 3,
            name: "Danny Smith",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-3.jpg",
            messages: [
              {
                id: 2,
                date: "2019-06-26 11:15",
                content: "Call me later!",
                unread: false
              }
            ]
          }
        ]
      },
      {
        id: 1,
        loggedUserId: 1,
        active: false,
        lastId: 2,
        chat: [
          {
            id: 1,
            name: "Lara Croft",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-5.jpg",
            messages: [
              {
                id: 0,
                date: "2019-04-21 15:00:09",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              },
              {
                id: 1,
                date: "2019-06-26 11:00",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              }
            ]
          },
          {
            id: 6,
            name: "Alex Turner",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-1.jpg",
            messages: [
              {
                id: 2,
                date: "2019-06-26 11:15",
                content: "Hey, are you at home?",
                unread: false
              }
            ]
          }
        ]
      },
      {
        id: 2,
        loggedUserId: 1,
        active: false,
        lastId: 2,
        chat: [
          {
            id: 1,
            name: "Lara Croft",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-5.jpg",
            messages: [
              {
                id: 0,
                date: "2019-04-21 15:00:09",
                content:
                    "Sed ut perspiciatis unde omnis iste natus remque laudantium."
              },
              {
                id: 1,
                date: "2019-06-26 11:00",
                content: "Sed ut doloremque laudantium."
              }
            ]
          },
          {
            id: 7,
            name: "Kate Moss",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-4.jpg",
            messages: [
              {
                id: 2,
                date: "2019-06-26 11:15",
                content: "Wanna grab a bite later?",
                unread: true
              }
            ]
          }
        ]
      },
      {
        id: 1,
        loggedUserId: 1,
        active: false,
        lastId: 2,
        chat: [
          {
            id: 1,
            name: "Lara Croft",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-5.jpg",
            messages: [
              {
                id: 0,
                date: "2019-04-21 15:00:09",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              },
              {
                id: 1,
                date: "2019-06-26 11:00",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              }
            ]
          },
          {
            id: 10,
            name: "Meg Ryan",
            online: true,
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-12.jpg",
            messages: [
              {
                id: 2,
                date: "2019-06-26 11:15",
                content: "Sed ut perspicantium.",
                unread: false
              }
            ]
          }
        ]
      },
      {
        id: 2,
        loggedUserId: 1,
        active: false,
        lastId: 2,
        chat: [
          {
            id: 1,
            name: "Lara Croft",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-6.jpg",
            messages: [
              {
                id: 0,
                date: "2019-04-21 15:00:09",
                content:
                    "Sed ut perspiciatis unde omnis iste natus remque laudantium."
              },
              {
                id: 1,
                date: "2019-06-26 11:00",
                content: "Sed ut doloremque laudantium."
              }
            ]
          },
          {
            id: 3,
            name: "John Smith",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-13.jpg",
            messages: [
              {
                id: 2,
                date: "2019-06-26 11:15",
                content: "Sed ut!",
                unread: true
              }
            ]
          }
        ]
      },
      {
        id: 1,
        loggedUserId: 1,
        active: false,
        lastId: 2,
        chat: [
          {
            id: 1,
            name: "Lara Croft",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-5.jpg",
            messages: [
              {
                id: 0,
                date: "2019-04-21 15:00:09",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              },
              {
                id: 1,
                date: "2019-06-26 11:00",
                content:
                    "Sed ut perspiciatis unde omnis iste natus error sit voluptatem accusantium doloremque laudantium."
              }
            ]
          },
          {
            id: 13,
            name: "Olenna Gervais",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-11.jpg",
            messages: [
              {
                id: 2,
                date: "2019-06-26 11:15",
                content: "Sed ut perspicantium.",
                unread: false
              }
            ]
          }
        ]
      },
      {
        id: 2,
        loggedUserId: 1,
        active: false,
        lastId: 2,
        chat: [
          {
            id: 1,
            name: "Lara Croft",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-5.jpg",
            messages: [
              {
                id: 0,
                date: "2019-04-21 15:00:09",
                content:
                    "Sed ut perspiciatis unde omnis iste natus remque laudantium."
              },
              {
                id: 1,
                date: "2019-06-26 11:00",
                content: "Sed ut doloremque laudantium."
              }
            ]
          },
          {
            id: 11,
            name: "Max Jackson",
            img: "https://mdbootstrap.com/img/Photos/Avatars/avatar-14.jpg",
            messages: [
              {
                id: 2,
                date: "2019-06-26 11:15",
                content: "Sed laudantium!",
                unread: false
              }
            ]
          }
        ]
      }
    ],
    activeChat: 0,
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
    createMessage(e, id) {
      const { content, unread, date } = e;
      return {
        id: id + 1,
        date,
        content,
        unread
      };
    },
    sendMessage(e, chat) {
      const newMessage = this.createMessage(e, chat.lastId);
      chat.chat
          .find(el => {
            return el.id === chat.loggedUserId;
          })
          .messages.push(newMessage);

      chat.lastId += 1;
    },
    changeActiveChat(index) {
      this.activeChat = index;
      this.gradientChat.forEach((chat, i) => {
        chat.active = i === index;
      });
    },
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
.rare-wind-gradient {
  background-image: -webkit-gradient(
      linear,
      left bottom,
      left top,
      from(#a8edea),
      to(#fed6e3)
  );
  background-image: -webkit-linear-gradient(bottom, #a8edea 0%, #fed6e3 100%);
  background-image: linear-gradient(to top, #a8edea 0%, #fed6e3 100%);
}

.scrollbar-light-blue::-webkit-scrollbar-track {
  -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.1);
  box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.1);
  background-color: #f5f5f5;
  border-radius: 10px;
}

.scrollbar-light-blue::-webkit-scrollbar {
  width: 12px;
  background-color: #f5f5f5;
}

.scrollbar-light-blue::-webkit-scrollbar-thumb {
  border-radius: 10px;
  -webkit-box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.1);
  box-shadow: inset 0 0 6px rgba(0, 0, 0, 0.1);
  background-color: #82b1ff;
}

.view {
  background: center center;
  background-size: cover;
  height: 100%;
}

.navbar i {
  cursor: pointer;
  color: white;
}
</style>