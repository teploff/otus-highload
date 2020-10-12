import Vue from 'vue';
import Router from 'vue-router';
import Questionnaires from '@/components/Questionnaires';
import SignUp from '@/components/SignUp';
import SignIn from '@/components/SignIn';
import SignUpSuccess from '@/components/SignUpSuccess';

Vue.use(Router);

export default new Router({
  routes: [
    {
      path: '/',
      name: 'Questionnaires',
      component: Questionnaires,
    },
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
  ],
  mode: 'history',
});
