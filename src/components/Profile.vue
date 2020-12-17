<template>
<div>

<div>
 {{errorMessage}}
</div>

 <button v-on:click="signMeOut">SIGN OUT</button>
  <button v-on:click="changePassword">CHANGE PASSWORD</button>

<h2>Friends</h2>
<h3>Hi {{email}}</h3>

<p>
  Send Message
  </p>

<form id="app" @submit="onSubmit">
  
  <p>
    <label for="toEmail">Email</label>
    <input type="email" name="toEmail" id="toEmail" v-model="toEmail">
  </p>

  <p>
    <label for="toMessage">Message</label>
    <input type="text" name="toMessage" id="toMessage" v-model="toMessage">
  </p>

  <p>
    <input type="submit" value="SEND">  
  </p>

</form>

    <ul>
      <li v-for="user in usersAll" @click="onSelect(user)">
          <span class="badge">{{user.email}}</span> {{user.firstName}} {{user.lastName}}
      </li>
    </ul>

    <ul>
      <li v-for="msg in msgsAll">
        <div>From: {{msg.messageFromEmail}}    To: {{msg.messageToEmail}}</div>
        <div>
        {{msg.message}}
        </div>
        <div>
        {{msg.sentAt}}
        </div>
    </li>
    </ul>



</div>
</template>

<script>

export default {
  name: 'Profile',
    data() {
        return {
            email: '',
            password: '',
            toEmail: '',
            toMessage: '',
            errorMessage: null,
            usersAll: [],
            msgsAll: []
        }
    },
    mounted: function () {

        this.email = localStorage.getItem('email')
        this.password = localStorage.getItem('password')


        if (this.email == null || this.password == null) {
            this.$router.push('/')
        }

        const that = this;
        fetch(
        "http://localhost:8080/all-friends",
        {
            method: "POST",
            headers: {
            "Content-Type": "application/x-www-form-urlencoded"
            },
            body:
            "email=" + this.email +
            "&authType=regular" +
            "&password=" + this.password
        }
        ).then(async function (data) {
        data.json().then(async function (data) {

            if (data.message === "Friends attached") {
            that.usersAll = data.usersAll;
            } else {
            // setErrorServerMessage(data.message);
            that.errorMessage = data.message;
            }
        });
        });

    },
    methods: {
        signMeOut: function() {
            localStorage.clear();
            this.$router.push('/')
        }, 
        changePassword: function() {
        this.$router.push('/change')
        },
        onSelect(user) {
            this.toEmail = user.email;
            this.showMessages(user.email)

        },
        onSubmit:function(e) {

            e.preventDefault();

            const that = this;

            fetch(
            "http://localhost:8080/send-message",
            {
                method: "POST",
                headers: {
                "Content-Type": "application/x-www-form-urlencoded"
                },
                body:
                "message=" + this.toMessage +
                "&messageFromEmail=" + this.email +
                "&messageToEmail=" + this.toEmail +
                "&authType=regular" +
                "&password=" + this.password
            }
            ).then(async function (data) {
            data.json().then(async function (data) {
                if (data.message === "Message sent") {
                that.errorMessage = "Message sent successfully.";
                that.showMessages(that.toEmail);
                } else {
                that.errorMessage = data.messge;
                }
            });
            });
        },
        showMessages (friendEmail) {
            const that = this;

            fetch(
            "http://localhost:8080/messages-user-and-friend",
            {
                method: "POST",
                headers: {
                "Content-Type": "application/x-www-form-urlencoded"
                },
                body:
                "userEmail=" + this.email +
                "&friendEmail=" + friendEmail +
                "&authType=regular" +
                "&password=" + this.password
            }
            ).then(async function (data) {
            data.json().then(async function (data) {

                if (data.message === "Messages attached") {
                that.msgsAll = data.msgs;
                } else {
                that.errorMessage = data.message;
                }
            });
            });
        }
    }
    }
</script>
