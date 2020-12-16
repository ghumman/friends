import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { FormBuilder } from '@angular/forms';

@Component({
  selector: 'app-register',
  templateUrl: './register.component.html',
  styleUrls: ['./register.component.css']
})
export class RegisterComponent implements OnInit {

  form;
  errorMessage;
  email;
  password;
  firstName;
  lastName;
  passwordConfirm;

  constructor(
    private router: Router,
    private formBuilder: FormBuilder
  ) { 
    this.email = localStorage.getItem('email');
    this.password = localStorage.getItem('password');

    this.form = this.formBuilder.group({
      firstName: '',
      lastName: '',
      email: '',
      password: '',
      passwordConfirm: ''
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

    this.firstName = customerData.firstName.trim();
    this.lastName = customerData.lastName.trim();
    this.email = customerData.email.trim();
    this.password = customerData.password.trim();
    this.passwordConfirm = customerData.passwordConfirm.trim();
    
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

          that.router.navigateByUrl('/profile');

        } else {
          that.errorMessage =  data.message ;
        }
      });
    });
  }
}
