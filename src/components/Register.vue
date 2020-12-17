<template>
<div>
<h2>Friends</h2>
<h3>We are glad you are going to become Friends User</h3>
<h3>Please Register</h3>

<form id="app" @submit="onSubmit">

  <p>
    <label for="firstName">First Name</label>
    <input type="text" name="firstName" id="firstName" v-model="firstName">
  </p>

  <p>
    <label for="lastName">Last Name</label>
    <input type="text" name="lastName" id="lastName" v-model="lastName">
  </p>
  
  <p>
    <label for="email">Email</label>
    <input type="email" name="email" id="email" v-model="email">
  </p>

  <p>
    <label for="password">Password</label>
    <input type="password" name="password" id="password" v-model="password">
  </p>

  <p>
    <label for="passwordConfirm">Password (again)</label>
    <input type="password" name="passwordConfirm" id="passwordConfirm" v-model="passwordConfirm">
  </p>

  <p>
    <input type="submit" value="Register">  
  </p>

</form>
 <button v-on:click="goToLogin">Cancel</button>

<div>
 {{errorMessage}}
</div>

</div>
</template>

<script>

export default {
  name: 'Register',
    data() {
    return {
      errorMessage: "",
      email:"",
      password:"",
      firstName: "",
      lastName: "",
      passwordConfirm: ""
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
    goToLogin: function() {
      this.$router.push('/')
    },
    onSubmit:function(e) {

      e.preventDefault();

      if (this.email == null || this.password == null) {
        
        this.errorMessage = "All fields are required. Please fill all of them.";
        return;

      }

      this.firstName = this.firstName.trim();
      this.lastName = this.lastName.trim();
      this.email = this.email.trim();
      this.password = this.password.trim();
      this.passwordConfirm = this.passwordConfirm.trim();

      if (
        (this.firstName === "") ||
        (this.lastName === "") ||
        (this.email === "") ||
        (this.password === "") ||
        (this.passwordConfirm === "")
      ) {
        this.errorMessage = "All fields are required. Please fill all of them.";
        return;
      }

      if (this.password !== this.passwordConfirm) {
        this.errorMessage = "Passwords don't match.";
        return;
      }

      const that = this;

      fetch(
        "http://localhost:8080/add-user",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded"
          },
          body:
            "firstName=" +
            this.firstName +
            "&lastName=" +
            this.lastName +
            "&email=" +
            this.email +
            "&password=" +
            this.password +
            "&authType=regular"

        }
      ).then(async function(data) {
        data.json().then(async function(data) {
          if (data.message === "User Created") {
            that.errorMessage = data.message;

            localStorage.setItem('email', that.email);
            localStorage.setItem('password', that.password);

            that.$router.push('/profile')

          } else {
            that.errorMessage =  data.message ;
          }
        });
      });
    }
  }
}
</script>
