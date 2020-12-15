import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder } from '@angular/forms';

@Component({
  selector: 'app-change',
  templateUrl: './change.component.html',
  styleUrls: ['./change.component.css']
})
export class ChangeComponent implements OnInit {

  form;
  email;
  password;
  currentPassword;
  newPassword;
  confirmNewPassword;
  errorMessage; 

  constructor(
    private router: Router,
    private formBuilder: FormBuilder
  ) {
    this.email = localStorage.getItem('email');
    this.password = localStorage.getItem('password');
    this.form = this.formBuilder.group({
      currentPassword: '',
      newPassword: '',
      confirmNewPassword: ''
    })
   }

  ngOnInit(): void {
    // if not logged-in, go to login page
    if (this.email === null) {
      this.router.navigateByUrl('/');
    }

    this.errorMessage = null;
  }

  onSubmit(customerData) {
    this.form.reset();
    this.newPassword = customerData.newPassword;
    this.confirmNewPassword = customerData.confirmNewPassword;
    this.currentPassword = customerData.currentPassword;

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
