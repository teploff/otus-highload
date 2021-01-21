import Vue from 'vue'
import App from './App.vue'
import router from './router';
import VueMaterial from 'vue-material'
import 'vue-material/dist/vue-material.min.css'
import 'vue-material/dist/theme/default.css'
import FlashMessage from '@smartweb/vue-flash-message';
import Paginate from 'vuejs-paginate'
import store from './store'
import VueMoment from 'vue-moment';
import wsService from "@/service/ws";

Vue.config.productionTip = false

Vue.use(VueMaterial)
Vue.use(VueMoment)
Vue.use(FlashMessage);

Vue.use(wsService, {
  store,
})

Vue.component('paginate', Paginate);

Vue.material.locale.dateFormat = 'yyyy-MM-dd'
new Vue({
  render: h => h(App),
  router,
  store,
}).$mount('#app')
