import Vue from 'vue'
import Router from 'vue-router'
import SignUp from '@/components/SignUp';
import SignIn from '@/components/SignIn';
import Home from '@/components/Home';
import News from '@/components/News';
import Messenger from '@/components/Messenger';
import Friends from '@/components/Friends';
import People from '@/components/People';

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
  if (to.name !== 'SignUp' && to.name !== 'SignIn' && !localStorage.getItem('accessToken')) {
    next({name: 'SignIn'});
  } else {
    next();
  }
});

export default router;