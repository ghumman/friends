<template>
<div>
<h2>Friends</h2>
<h3>Forgot Password</h3>

<form id="app" @submit="onSubmit">
  
  <p>
    <label for="email">Email</label>
    <input type="email" name="email" id="email" v-model="email">
  </p>

  <p>
    <input type="submit" value="FORGOT">  
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
  name: 'Forgot',
    data() {
    return {
      errorMessage: null,
      email:""
    }
  },
  mounted: function () {

    this.email = localStorage.getItem('email')


    if (this.email != null) {
        this.$router.push('/profile')
    }
},
  methods: {
    goToLogin: function() {
      this.$router.push('/')
    },
    onSubmit:function(e) {

        e.preventDefault();

        const that = this;
        fetch(
        "http://localhost:8080/forgot-password",
        {
            method: "POST",
            headers: {
            "Content-Type": "application/x-www-form-urlencoded"
            },
            body: "email=" + this.email
        }
        ).then(async function(data) {
            data.json().then(async function(data) {
                if (data.message === "Reset password is sent") {
                that.errorMessage = "Password reset link is sent to your email.";
                } else {
                that.errorMessage =  data.message;
                }
            });
        });
    }
  }
}
</script>
