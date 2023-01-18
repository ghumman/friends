import React, { useState, useEffect } from "react";
import { connect } from "react-redux";


import Button from "react-bootstrap/Button";


import Form from "react-bootstrap/Form";

import { backendAddress } from "./../config/default-variables.js"


const PASSWORD = "password";

function Login(props) {

  const [errorCurrentPassword, setErrorCurrentPassword] = useState("");
  const [errorServerMessage, setErrorServerMessage] = useState("");
  const [globalUsername, setGlobalUsername] = useState("");
  const [globalPassword, setGlobalPassword] = useState("");
  const [formCurrentPassword, setFormCurrentPassword] = useState("");
  const [formPassword, setFormPassword] = useState("");
  const [formPassword2, setFormPassword2] = useState("");
  const [currentBackendAddress, setCurrentBackendAddress] = useState(backendAddress)



  const handleChangeCurrentPassword = (event) => {
    setFormCurrentPassword(event.target.value);
  }

  const handleChangePassword = (event) => {
    setFormPassword(event.target.value);
  }

  const handleChangePassword2 = (event) => {
    setFormPassword2(event.target.value); 
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    var password = formPassword;
    var password2 = formPassword2;
    var currentPassword = formCurrentPassword;

    if (currentPassword === globalPassword) {

      setErrorCurrentPassword("");
      if (password.length >= 5 && password === password2) {

        fetch(
          currentBackendAddress + '/change-password',
          {
            method: "POST",
            headers: {
              "Content-Type": "application/x-www-form-urlencoded"
            },
            body:
              "email=" +
              globalUsername.trim() +
              "&password=" +
              currentPassword.trim() +
              "&newPassword=" +
              password.trim() +
              "&authType=regular"
          }
        ).then(async function(data) {
          data.json().then(async function(data) {
            if (data.message === "Password changed") {
              setErrorServerMessage("Password changed successfully.");

              localStorage.setItem(PASSWORD, password);
            } else {
              setErrorServerMessage(data.message);
            }
          });
        });
      } else {
        setErrorServerMessage("Password should be 5 characters long and match.")
      }
    } else {
      setErrorCurrentPassword("Password does not match current password.");
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
      setGlobalUsername(resultUsername); 
      setGlobalPassword(resultPassword); 
    } else {
      props.login(resultUsername, resultPassword);
      props.history.push({
        pathname: "/"
      });
    }
  }, [])

  const loginPressed = () => {
    props.history.push({
      pathname: "/Profile"
    });
  }




    var errorMessageCurrent;
    if (errorCurrentPassword !== "") {
      errorMessageCurrent = (
        <p
          style={{
            paddingTop: 20,
            color: "red"
          }}
        >
          {errorCurrentPassword}
        </p>
      );
    } else {
      errorMessageCurrent = null;
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
            Change Password
          </p>
          <Form onSubmit={handleSubmit}>
            <Form.Group controlId="loginFormCurrentPassword">
              <Form.Label>Current Password</Form.Label>
              <Form.Control
                type="password"
                value={formCurrentPassword}
                placeholder="Current Password"
                onChange={handleChangeCurrentPassword}
              />
              {errorMessageCurrent}
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
              Back / Cancel
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
)(Login);
