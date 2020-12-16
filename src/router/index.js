// import Vue from 'vue'
// import Router from 'vue-router'
// import Login from '../components/Login.vue'
// import Register from '../components/Register.vue'

import { createRouter, createWebHistory } from 'vue-router'
import Login from '../components/Login.vue'
import Register from '../components/Register.vue'

const routerHistory = createWebHistory()

const router = createRouter({
    history: routerHistory,
    routes: [
      {
        path: '/',
        component: Login
      },
      {
        path: '/register',
        component: Register
      }
    ]
  })

  export default router



// export default new Router({
//   routes: [
//     {
//       path: '/',
//       name: 'Login',
//       component: Login
//     },
//     {
//         path: '/',
//         name: 'Register',
//         component: Register
//       },
//   ]
// })