<template>
  <router-view/>
</template>

<script>
import Login from './components/Login.vue'
import Register from './components/Register.vue'


export default {
  name: 'App',
  components: {
    Login
  }
}
</script>
