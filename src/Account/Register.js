import React, { useState, useEffect } from "react";
import { connect } from "react-redux";


import Button from "react-bootstrap/Button";

import Form from "react-bootstrap/Form";

import { Link } from "react-router-dom";

import { backendAddress } from "./../config/default-variables.js"

const USERNAME = "username";
const PASSWORD = "password";
const BACKURL = "backurl";

function Register(props) {


  const [formFirstName, setFormFirstName] = useState(""); 
  const [formLastName, setFormLastName] = useState(""); 
  const [formPassword, setFormPassword] = useState(""); 
  const [formPassword2, setFormPassword2] = useState(""); 
  const [formEmail, setFormEmail] = useState(""); 
  const [errorMessage, setErrorMessage] = useState(""); 
  const [currentBackendAddress, setCurrentBackendAddress] = useState(backendAddress)

  const handleChangeFirstName = (event) => {
    setFormFirstName(event.target.value );
  }
  const handleChangeLastName = (event) => {
    setFormLastName(event.target.value ); 
  }

  const handleChangeEmail = (event) => {
    setFormEmail(event.target.value ); 
  }

  const handleChangePassword = (event) => {
    setFormPassword(event.target.value);
  }

  const handleChangePassword2 = (event) => {
    setFormPassword2(event.target.value ); 
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    if (
      (formFirstName.trim() === "") |
      (formLastName.trim() === "") |
      (formEmail.trim() === "") |
      (formPassword.trim() === "") |
      (formPassword2.trim() === "")
    ) {
      setErrorMessage("All fields are required. Please fill all of them.")
      return;
    }

    if (formPassword.trim() !== formPassword2.trim()) {
      setErrorMessage("Passwords don't match.");
      return;
    }

    var firstName = formFirstName;
    var lastName = formLastName;
    var email = formEmail;
    var password = formPassword;


    fetch(
      currentBackendAddress + '/add-user',
      {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        body:
          "firstName=" +
          firstName.trim() +
          "&lastName=" +
          lastName.trim() +
          "&email=" +
          email.trim() +
          "&password=" +
          password.trim() +
          "&authType=regular"

      }
    ).then(async function(data) {
      data.json().then(async function(data) {
        if (data.message === "User Created") {
          setErrorMessage(data.message );

          localStorage.setItem(USERNAME, email.trim());
          localStorage.setItem(PASSWORD, password.trim());

          props.login(formEmail, formPassword);

          props.history.push({
            pathname: "/Profile"
          });
        } else {
          setErrorMessage( data.message );
        }
      });
    });
  }

  // in componentDidMount 
  // check if user credentials are already saved,
  // in that case save it in redux and send him to profile PAGE
  useEffect(() => {
    setCurrentBackendAddress(localStorage.getItem("backurl") || backendAddress);
    const resultUsername = localStorage.getItem("username");
    const resultPassword = localStorage.getItem("password");
    if (
      resultUsername !== null &&
      resultUsername !== "" &&
      resultPassword !== null &&
      resultPassword !== ""
    ) {
      props.login(resultUsername, resultPassword);
      props.history.push({
        pathname: "/Profile"
      });
    } 
  }, [])
  // componentDidMount() {
    
  // }

    var errorMessageLocal;
    if (errorMessage.trim() === "") errorMessageLocal = null;
    else
      errorMessageLocal = (
        <p style={{ color: "red" }}>{errorMessage.trim()}</p>
      );

    return (
      <div className="App">
        <header className="App-header">
          <p>Friends</p>
          <p>We are glad you are going to become Friends User</p>
          {errorMessageLocal}

          <p
            style={{
              borderStyle: "solid",
              borderColor: "white",
              borderRadius: 5,
              padding: 5
            }}
          >
            Please Register
          </p>
          <Form onSubmit={handleSubmit}>
            <Form.Group controlId="loginFormFirstName">
              <Form.Label>First Name</Form.Label>
              <Form.Control
                type="username"
                value={formFirstName}
                placeholder="Enter First Name"
                onChange={handleChangeFirstName}
              />
            </Form.Group>

            <Form.Group controlId="loginFormLastName">
              <Form.Label>Last Name</Form.Label>
              <Form.Control
                type="username"
                value={formLastName}
                placeholder="Enter Last Name"
                onChange={handleChangeLastName}
              />
            </Form.Group>
            <Form.Group controlId="loginFormEmail">
              <Form.Label>Email</Form.Label>
              <Form.Control
                type="email"
                value={formEmail}
                placeholder="Enter Email"
                onChange={handleChangeEmail}
              />
            </Form.Group>

            <Form.Group controlId="loginFormPassword">
              <Form.Label>Password</Form.Label>
              <Form.Control
                type="password"
                value={formPassword}
                placeholder="Password"
                onChange={handleChangePassword}
              />
            </Form.Group>

            <Form.Group controlId="loginFormPassword2">
              <Form.Label>Password (again)</Form.Label>
              <Form.Control
                type="password"
                value={formPassword2}
                placeholder="Password"
                onChange={handleChangePassword2}
              />
            </Form.Group>

            <Button variant="primary" type="submit">
              Sign Up
            </Button>
            <Link to="/" className="btn btn-link">
              Cancel
            </Link>
          </Form>
          <br />
          <br />
        </header>
      </div>
    );
}


function mapStateToProps(state) {
  return {
    authenticate: state
  };
}

function mapDispatchToProps(dispatch) {
  return {
    login: (username, password) =>
      dispatch({ type: "UPDATE", username: username, password: password })
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Register);
