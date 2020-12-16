import { createApp } from 'vue'
// import Vue, { createApp } from 'vue'
import App from './App.vue'
import router from './router'
import './index.css'
// import VueResource from 'vue-resource'

// Vue.use(Router)


// createApp(App).mount('#app')

const app = createApp(App)
app.use(router)
app.mount('#app')

// new Vue({
//     router, 
//     component: {App},
//   }).$mount('#app')

/* eslint-disable no-new */
// new Vue({
//     el: '#app',
//     router,
//     components: { App },
//     template: '<App/>'
//   })