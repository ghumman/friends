import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder } from '@angular/forms';

@Component({
  selector: 'app-forgot',
  templateUrl: './forgot.component.html',
  styleUrls: ['./forgot.component.css']
})
export class ForgotComponent implements OnInit {

  form;
  errorMessage;
  email;
  password;

  constructor(
    private router: Router,
    private formBuilder: FormBuilder
  ) {
    this.email = localStorage.getItem('email');
    this.password = localStorage.getItem('password');

    this.form = this.formBuilder.group({
      email: ''
    })
   }

  ngOnInit(): void {
    if (this.email !== null && this.password !== null) {
      this.router.navigateByUrl('/profile');
    }
  }

  goToLogin() {
    this.router.navigate(['/']);
  }


  onSubmit(customerData) {

    this.email = customerData.email.trim();

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
