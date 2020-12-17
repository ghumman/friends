<template>
<div>
<h2>Friends</h2>
<h3>Please Login</h3>

<form id="app" @submit="onSubmit">
  
  <p>
    <label for="currentPassword">Current Password</label>
    <input type="password" name="currentPassword" id="currentPassword" v-model="currentPassword">
  </p>

  <p>
    <label for="newPassword">New Password</label>
    <input type="password" name="newPassword" id="newPassword" v-model="newPassword">
  </p>

  <p>
    <label for="confirmNewPassword">Confirm New Password</label>
    <input type="password" name="confirmNewPassword" id="confirmNewPassword" v-model="confirmNewPassword">
  </p>

  <p>
    <input type="submit" value="Change Password">  
  </p>

</form>

<div>
 {{errorMessage}}
</div>

</div>
</template>

<script>

export default {
  name: 'Change',
    data() {
    return {
      errorMessage: null,
      email:null,
      password:null,
      currentPassword: null,
      newPassword: null,
      confirmNewPassword: null
    }
  },
  mounted: function () {

    this.email = localStorage.getItem('email')
    this.password = localStorage.getItem('password')


    if (this.email === null || this.password === null) {
        this.$router.push('/')
    }
},
  methods: {
    onSubmit:function(e) {

    e.preventDefault();

        if (this.currentPassword === this.password) {

            if (this.newPassword !== null && this.newPassword.length >= 5 && this.newPassword === this.confirmNewPassword) {

                const that = this;
                fetch(
                "http://localhost:8080/change-password",
                {
                    method: "POST",
                    headers: {
                    "Content-Type": "application/x-www-form-urlencoded"
                    },
                    body:
                    "email=" +
                    this.email.trim() +
                    "&password=" +
                    this.currentPassword.trim() +
                    "&newPassword=" +
                    this.newPassword.trim() +
                    "&authType=regular"
                }
                ).then(async function(data) {
                data.json().then(async function(data) {
                    if (data.message === "Password changed") {
                    that.errorMessage = "Password changed successfully.";

                    localStorage.setItem('password', that.newPassword);
                    that.password = that.newPassword;
                    } else {
                    that.errorMessage = data.message;
                    }
                });
                });
            } else {
                this.errorMessage = "Password should be 5 characters long and match.";
            }
        } else {
        this.errorMessage = "Password does not match current password.";
        }
    }
  }
}
</script>
