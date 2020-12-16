import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder } from '@angular/forms';

@Component({
  selector: 'app-reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.css']
})
export class ResetPasswordComponent implements OnInit {

  form;
  errorMessage;
  email;
  password;
  passwordConfirm;
  token;

  constructor(
    private router: Router,
    private formBuilder: FormBuilder
  ) { 
    this.email = localStorage.getItem('email');
    this.password = localStorage.getItem('password');

    this.form = this.formBuilder.group({
      email: '',
      password: '',
      passwordConfirm: ''
    })
  }

  ngOnInit(): void {
    if (this.email !== null && this.password !== null) {
      this.router.navigateByUrl('/profile');
    }

    this.token = this.router.parseUrl(this.router.url).queryParams.token;
    if (this.token !== null && this.token !== undefined && this.token !== "") {
      this.errorMessage = "Please select new password.";
    } else {
      this.errorMessage = "Url is tempered";
    }

  }


  goToLogin() {
    this.router.navigate(['/']);
  }

  onSubmit(customerData) {
    this.password = customerData.password.trim();
    this.passwordConfirm = customerData.passwordConfirm.trim();

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
