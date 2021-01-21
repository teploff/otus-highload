import Vue from 'vue';
import Router from 'vue-router';
import SignUp from '@/components/SignUp';
import SignIn from '@/components/SignIn';
import SignUpSuccess from '@/components/SignUpSuccess';
import Home from '@/components/Home';
import News from '@/components/News';
import Messenger from '@/components/Messenger';
import Friends from '@/components/Friends';
import People from '@/components/People';
import store from '@/store'

Vue.use(Router);

const router = new Router({
  routes: [
    {
      path: '/sign-in',
      name: 'SignIn',
      component: SignIn,
    },
    {
      path: '/sign-up',
      name: 'SignUp',
      component: SignUp,
    },
    {
      path: '/sign-up-success',
      name: 'SignUpSuccess',
      component: SignUpSuccess,
    },
    {
      path: '/',
      name: 'Home',
      component: Home,
    },
    {
      path: '/news',
      name: 'News',
      component: News,
    },
    {
      path: '/messenger',
      name: 'Messenger',
      component: Messenger,
    },
    {
      path: '/friends',
      name: 'Friends',
      component: Friends,
    },
    {
      path: '/people',
      name: 'People',
      component: People,
    },
  ],
  mode: 'history',
});

router.beforeEach((to, from, next) => {
  if (to.name !== 'SignUp' && to.name !== 'SignUpSuccess' && to.name !== 'SignIn' && !localStorage.getItem('accessToken')) {
    next({name: 'SignIn'});
  } else {
    if (to.name !== 'SignUp' && to.name !== 'SignUpSuccess' && to.name !== 'SignIn') {
      if (!store.getters.isWSConnected) {
        router.app.$wsConnect('ws://localhost:9999/ws?token=' + localStorage.getItem('accessToken'))
      }
    }

    next();
  }
});

export default router;
