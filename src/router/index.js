// import Vue from 'vue'
// import Router from 'vue-router'
// import Login from '../components/Login.vue'
// import Register from '../components/Register.vue'

import { createRouter, createWebHistory } from 'vue-router'

import Change from '../components/Change.vue'
import Forgot from '../components/Forgot.vue'
import Login from '../components/Login.vue'
import Profile from '../components/Profile.vue'
import Register from '../components/Register.vue'
import ResetPassword from '../components/ResetPassword.vue'

const routerHistory = createWebHistory()

const router = createRouter({
    history: routerHistory,
    routes: [
      {
        path: '/change',
        component: Change
      },
      {
        path: '/forgot',
        component: Forgot
      },
      {
        path: '/',
        component: Login
      },
      {
        path: '/profile',
        component: Profile
      },
      {
        path: '/register',
        component: Register
      },
      {
        path: '/reset-password',
        component: ResetPassword
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