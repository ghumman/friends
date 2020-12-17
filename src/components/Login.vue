<template>
<div>
<h2>Friends</h2>
<h3>Please Login</h3>

<form id="app" @submit="onSubmit">
  
  <p>
    <label for="email">Email</label>
    <input type="email" name="email" id="email" v-model="email">
  </p>

  <p>
    <label for="password">Password</label>
    <input type="password" name="password" id="password" v-model="password">
  </p>

  <p>
    <input type="submit" value="Login">  
  </p>

</form>
 <button v-on:click="goToRegister">Register</button>
  <button v-on:click="goToForgot">Need help?</button>

<div>
 {{errorMessage}}
</div>

</div>
</template>

<script>

export default {
  name: 'Login',
    data() {
    return {
      errorMessage: null,
      email:null,
      password:null
    }
  },
  mounted: function () {

    this.email = localStorage.getItem('email')
    this.password = localStorage.getItem('password')


    if (this.email != null && this.password != null) {
        this.$router.push('/profile')
    }
},
  methods: {
    goToRegister: function() {
      this.$router.push('/register')
    }, 
    goToForgot: function() {
      this.$router.push('/forgot')
    },
    onSubmit:function(e) {

    e.preventDefault();
    const that = this;

    this.email = this.email.trim();
    this.password = this.password.trim();

    fetch(
      "http://localhost:8080/login",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        body: "email=" + this.email + "&password=" + this.password + "&authType=regular"
      }).then(async function(data) {
      data.json().then(async function(data) {
        if (data.message === "Logged In") {


          localStorage.setItem('email', that.email);
          localStorage.setItem('password', that.password);

          that.$router.push('/profile')
        } else {
          that.errorMessage = data.message;
        }
      });
    });

    }
  }
}
</script>
