import Vue from 'vue'
import App from './App.vue'
import router from './router';
import VueMaterial from 'vue-material'
import 'vue-material/dist/vue-material.min.css'
import 'vue-material/dist/theme/default.css'
import FlashMessage from '@smartweb/vue-flash-message';

Vue.config.productionTip = false

Vue.use(VueMaterial)
Vue.use(FlashMessage);

Vue.material.locale.dateFormat = 'yyyy-MM-dd'
new Vue({
  render: h => h(App),
  router,
}).$mount('#app')
