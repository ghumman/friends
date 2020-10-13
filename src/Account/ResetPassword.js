import React, { useEffect, useState } from "react";
import { connect } from "react-redux";


import Button from "react-bootstrap/Button";
// import ButtonToolbar from "react-bootstrap/ButtonToolbar";


import Form from "react-bootstrap/Form";





function Login(props) {

  const [errorServerMessage, setErrorServerMessage] = useState(""); 
  const [formPassword, setFormPassword] = useState(""); 
  const [formPassword2, setFormPassword2] = useState(""); 
  const [selector, setSelector] = useState(""); 
  const [validator, setValidator] = useState(""); 
  

  const handleChangePassword = (event) => {
    setFormPassword(event.target.value);
  }

  const handleChangePassword2 = (event) => {
    setFormPassword2(event.target.value);
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    console.log("Submit button is pressed.");
    var password = formPassword;
    var password2 = formPassword2;

    if (password.length >= 5 && password === password2) {

      fetch(
        "https://localhost/resetPassword",
        {
          method: "POST",
          headers: {
            "Content-Type": "application/x-www-form-urlencoded"
          },
          body:
            "password=" +
            password.trim() +
            "&selector=" +
            selector +
            "&validator=" +
            validator
        }
      ).then(async function(data) {
        data.json().then(async function(data) {
          if (data.message === "Password updated successfully") {
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

  // in componentDidMount
  // check if user credentials are already saved,
  // in that case save it in redux and send him to profile PAGE

  useEffect(() => {

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
    var localSelector, localValidator, localReset;
    var regex = /[?&]([^=#]+)=([^&#]*)/g,
      params = {},
      match;
    while ((match = regex.exec(url))) {
      params[match[1]] = match[2];
      console.log(match[1], match[2]);
      if (match[1] === "selector") {
        console.log("match1 is selector");
        localSelector = match[2];
        console.log("localSelector");
        console.log(localSelector);
      }
      if (match[1] === "validator") {
        console.log("match1 is validator");
        localValidator = match[2];
        console.log("localValidator");
        console.log(localValidator);
      }
      if (match[1] === "reset") {
        console.log("match1 is reset");
        localReset = match[2];
        console.log("localReset");
        console.log(localReset);
      }
    } // while ends

    if (
      (localSelector !== "") &
      (localSelector !== null) &
      ((localValidator !== "") & (localValidator !== null) & (localReset === "1"))
    ) {
      setSelector(localSelector);
      setValidator(localValidator);
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
)(Login);

