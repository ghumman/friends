<template>
<div>
<h2>Friends</h2>
<h3>Forgot Password</h3>

<form id="app" @submit="onSubmit">


  <p>
    <label for="password">Password</label>
    <input type="password" name="password" id="password" v-model="password">
  </p>

  <p>
    <label for="passwordConfirm">Password (again)</label>
    <input type="password" name="passwordConfirm" id="passwordConfirm" v-model="passwordConfirm">
  </p>

  <p>
    <input type="submit" value="RESET PASSWORD">  
  </p>

</form>
 <button v-on:click="goToLogin">Back to Login</button>

<div>
 {{errorMessage}}
</div>

</div>
</template>

<script>

export default {
  name: 'ResetPassword',
    data() {
    return {
      errorMessage: null,
      email:null,
      password:null,
      passwordConfirm: null,
      token: null
    }
  },
  mounted: function () {

    this.email = localStorage.getItem('email')
    this.password = localStorage.getItem('password')


    if (this.email != null && this.password != null) {
        this.$router.push('/profile')
    }

    // this.token = this.router.parseUrl(this.router.url).queryParams.token;
    this.token = this.$route.query.token;
    if (this.token !== null && this.token !== undefined && this.token !== "") {
      this.errorMessage = "Please select new password.";
    } else {
      this.errorMessage = "Url is tempered";
    }
},
  methods: {
    goToLogin: function() {
      this.$router.push('/')
    },
    onSubmit:function(e) {

    e.preventDefault();

    if (this.password == null || this.passwordConfirm == null) {
        this.errorMessage = "Password should be 5 characters long and match.";
        return;
    }

    this.password = this.password.trim();
    this.passwordConfirm = this.passwordConfirm.trim();

    const that = this;

    if (this.password.length >= 5 && this.password === this.passwordConfirm) {

      fetch(
        "http://localhost:8080/reset-password",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded"
          },
          body:
            "password=" +
            this.password +
            "&token=" +
            this.token
        }
      ).then(async function(data) {
        data.json().then(async function(data) {
          if (data.message === "Password successfully reset") {
            that.errorMessage = "Password updated successfully.";
          } else {
            that.errorMessage = data.message;

          }
        });
      });
    } else {
      this.errorMessage = "Password should be 5 characters long and match.";
    }

    }
  }
}
</script>
