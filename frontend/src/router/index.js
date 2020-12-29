import Vue from 'vue';
import Router from 'vue-router';
import SignUp from '@/components/SignUp';
import SignIn from '@/components/SignIn';
import SignUpSuccess from '@/components/SignUpSuccess';
import Home from '@/components/Home';
import News from '@/components/News';
import Messenger from '@/components/Messenger';
import Friends from '@/components/Friends';

Vue.use(Router);

export default new Router({
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
  ],
  mode: 'history',
});
