import { Component, OnInit, Input } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder } from '@angular/forms';

import { filter } from 'rxjs/operators';

@Component({
  selector: 'app-profile',
  templateUrl: './profile.component.html',
  styleUrls: ['./profile.component.css']
})
export class ProfileComponent implements OnInit {

  email;
  password;
  errorMessage;
  usersAll;
  msgsAll;
  form;

  constructor(
    private router: Router,
    private formBuilder: FormBuilder
  ) { 

    this.form = this.formBuilder.group({
      toEmail: '',
      toMessage: ''
    })
    

    // this.email = window.history.state.email;
    // this.password = window.history.state.password;

    this.email = localStorage.getItem('email');
    this.password = localStorage.getItem('password');

  this.errorMessage = null;
  }

  ngOnInit(): void {

    // if not logged-in, go to login page
    if (this.email === null || this.password === null) {
      this.router.navigateByUrl('/');
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
  }

  signMeOut() {

    // reset email and password from local storage
    localStorage.removeItem("email");
    localStorage.removeItem("password");

    this.router.navigateByUrl('/');

  }

  changePassword() {
    this.router.navigateByUrl('/change');
  }

  onSelect(user) {

    this.form.get('toEmail').setValue(user.email);

    this.showMessages(user.email);
  }

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

  onSubmit(customerData) {
    this.form.reset();

    const that = this;

    fetch(
      "http://localhost:8080/send-message",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        body:
          "message=" + customerData.toMessage +
          "&messageFromEmail=" + this.email +
          "&messageToEmail=" + customerData.toEmail +
          "&authType=regular" +
          "&password=" + this.password
      }
    ).then(async function (data) {
      data.json().then(async function (data) {
        if (data.message === "Message sent") {
          that.errorMessage = "Message sent successfully.";
          that.showMessages(customerData.toEmail);
        } else {
          that.errorMessage = data.messge;
        }
      });
    });
  }
}
