import React, {useState, useEffect}  from "react";
import { connect } from "react-redux";


import Button from "react-bootstrap/Button";

import Form from "react-bootstrap/Form";

import { Link } from "react-router-dom";

const USERNAME = "username";
const PASSWORD = "password";

function Login(props) {
  const [errorServerMessage, setErrorServerMessage] = useState("");
  const [formUsername, setFormUsername] = useState("");
  const [formPassword, setFormPassword] = useState(""); 

  const handleChangeUsername = (event) => {
    // this.setState({ formUsername: event.target.value });
    setFormUsername(event.target.value); 
  }

  const handleChangePassword = (event) => {
    // this.setState({ formPassword: event.target.value });
    setFormPassword(event.target.value); 
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    console.log("Submit button is pressed.");
    var username = formUsername;
    var password = formPassword;


    fetch(
      "http://localhost:8080/login",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        body: "email=" + username.trim() + "&password=" + password.trim() + "&authType=regular"
      }).then(async function(data) {
      data.json().then(async function(data) {
        console.log("Value of data.message:");
        console.log(data.message);
        if (data.message === "Logged In") {
          console.log("Login successfull");
          props.login(formUsername, formPassword);

          /*
            save the username and password in local storage
            so that we don't have to ask for username password again from username
          */

          localStorage.setItem(USERNAME, username.trim());
          localStorage.setItem(PASSWORD, password.trim());

          props.history.push({
            pathname: "/Profile"
          });
        } else {
          console.log(data.message);
          // that.setState({ errorServerMessage: data.message });
          setErrorServerMessage(data.message);
        }
      });
    });
  }





  useEffect(()=>{
    // console.log('component mounted!')
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
  },[]) //notice the empty array here


  const needHelp = () => {
    props.history.push({
      pathname: "/Forgot"
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
          Please Login
        </p>
        <Form onSubmit={handleSubmit}>
          <Form.Group controlId="loginFormUsername">
            <Form.Label>Username</Form.Label>
            <Form.Control
              type="username"
              value={formUsername}
              placeholder="Enter Username"
              onChange={handleChangeUsername}
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

          <Button variant="primary" type="submit">
            Sign In
          </Button>
          <Link to="/register" className="btn btn-link">
            Register
          </Link>
        </Form>

        <div>
          <p style={{ fontSize: 14, padding: 20 }} onClick={needHelp}>
            Need help?
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