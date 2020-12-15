import { Component, OnInit } from '@angular/core';
import { Router, ActivatedRoute, ParamMap } from '@angular/router';
import { FormBuilder } from '@angular/forms';

@Component({
  selector: 'app-login',
  templateUrl: './login.component.html',
  styleUrls: ['./login.component.css']
})
export class LoginComponent implements OnInit {

  form;
  errorMessage;
  email;
  password;

  constructor(
    private route: ActivatedRoute,
    private router: Router,
    private formBuilder: FormBuilder) {
      this.email = localStorage.getItem('email');
      this.password = localStorage.getItem('password');
      this.form = this.formBuilder.group({
        email: '',
        password: ''
      })
    }

  ngOnInit(): void {
    this.errorMessage = null;

    if (this.email !== null && this.password !== null) {
      this.router.navigateByUrl('/profile');
    }

  }

  goToRegister(event: Event) {
    console.log('Click!', event)
    this.router.navigate(['register', { id: "test data from login" }]);
  }

  onSubmit(customerData) {
    this.form.reset();


    const that = this;

    fetch(
      "http://localhost:8080/login",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        body: "email=" + customerData.email.trim() + "&password=" + customerData.password.trim() + "&authType=regular"
      }).then(async function(data) {
      data.json().then(async function(data) {
        if (data.message === "Logged In") {

          localStorage.setItem('email', customerData.email.trim());
          localStorage.setItem('password', customerData.password.trim());

          that.router.navigateByUrl('/profile');
          // in-case you want to pass variables instead of saving in local storage
          // that.router.navigateByUrl('/profile', { state: { email: customerData.email.trim(), password:  customerData.password.trim()} });
        } else {
          that.errorMessage = data.message;
        }
      });
    });
  }

}
