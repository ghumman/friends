import React, { useEffect, useState } from "react";
import { connect } from "react-redux";

import Button from "react-bootstrap/Button";

import Form from "react-bootstrap/Form";

import { backendAddress } from "./../config/default-variables.js"

function ResetPassword(props) {

  const [errorServerMessage, setErrorServerMessage] = useState(""); 
  const [formEmail, setFormEmail] = useState(""); 
  const [formPassword, setFormPassword] = useState(""); 
  const [formPassword2, setFormPassword2] = useState(""); 
  const [token, setToken] = useState(""); 
  const [currentBackendAddress, setCurrentBackendAddress] = useState(backendAddress)
  
  const handleChangeEmail = (event) => {
    setFormEmail(event.target.value);
  }

  const handleChangePassword = (event) => {
    setFormPassword(event.target.value);
  }

  const handleChangePassword2 = (event) => {
    setFormPassword2(event.target.value);
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    var email = formEmail;
    var password = formPassword;
    var password2 = formPassword2;

    if (password.length >= 5 && password === password2) {

      fetch(

        currentBackendAddress + '/reset-password',
        {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded"
          },
          body:
            "email=" +
            email.trim() +
            "&password=" +
            password.trim() +
            "&token=" +
            token
        }
      ).then(async function(data) {
        data.json().then(async function(data) {
          if (data.message === "Password successfully reset") {
            setErrorServerMessage("Password updated successfully.")
          } else {
            setErrorServerMessage(data.message)

          }
        });
      });
    } else {
      setErrorServerMessage("Password should be 5 characters long and match.")
    }
  }

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
    } else {
      parseUrl();
    }
  }, [])

  const parseUrl = () => {

    var url = window.location.href;
    var localToken;
    var regex = /[?&]([^=#]+)=([^&#]*)/g,
      params = {},
      match;
    while ((match = regex.exec(url))) {
      params[match[1]] = match[2];
      if (match[1] === "token") {
        localToken = match[2];
      }
    } // while ends


    if (
      (localToken !== "") &
      (localToken !== undefined) 
    ) {
      setToken(localToken);
      setErrorServerMessage("Please select new password.");
    } else {
      setErrorServerMessage("Url is tempered")
    }
  }

  const loginPressed = () => {
    props.history.push({
      pathname: "/"
    });
  }

    var errorMessage;
    if (errorServerMessage !== "") {
      errorMessage = (
        <p
          style={{
            paddingTop: 20,
            color: "red"
          }}
        >
          {errorServerMessage}
        </p>
      );
    } else {
      errorMessage = null;
    }
    return (
      <div className="App">
        <header className="App-header">
          <p>Friends</p>
          <p
            style={{
              borderStyle: "solid",
              borderColor: "white",
              borderRadius: 5,
              padding: 5
            }}
          >
            Forgot Password
          </p>
          <Form onSubmit={handleSubmit}>
          <Form.Group controlId="loginFormEmail">
              <Form.Label>Email</Form.Label>
              <Form.Control
                type="email"
                value={formEmail}
                placeholder="Email"
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
              RESET PASSWORD
            </Button>
          </Form>

          <div>
            <p
              style={{ fontSize: 14, paddingTop: 20 }}
              onClick={loginPressed}
            >
              Back to Login
            </p>
          </div>

          {errorMessage}
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
)(ResetPassword);

