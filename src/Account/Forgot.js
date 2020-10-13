import React, { useEffect, useState } from "react";
import { connect } from "react-redux";



import Button from "react-bootstrap/Button";
// import ButtonToolbar from "react-bootstrap/ButtonToolbar";


import Form from "react-bootstrap/Form";





// function App() {
// class Login extends React.Component {
function Login(props) {

  const [errorServerMessage, setErrorServerMessage] = useState("");
  const [formEmail, setFormEmail] = useState("");

  const handleChangeEmail = (event) => {
    setFormEmail(event.target.value);
  }

  const handleSubmit = (event) => {
    event.preventDefault();
    console.log("Submit button is pressed.");
    var email = formEmail;


    fetch(
      "https://localhost/forgot",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        body: "email=" + email.trim()
      }
    ).then(async function(data) {
      data.json().then(async function(data) {
        console.log("Value of data.message:");
        console.log(data.message);
        if (data.message === "Password is sent to your email") {
          setErrorServerMessage("Password reset link is sent to your email.");
        } else {
          setErrorServerMessage(data.message);
        }
      });
    });
  }


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
    }
  }, [])

  const cancelPressed = () => {
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
                placeholder="Enter Email"
                onChange={handleChangeEmail}
              />
            </Form.Group>

            <Button variant="primary" type="submit">
              FORGOT
            </Button>
          </Form>

          <div>
            <p
              style={{ fontSize: 14, paddingTop: 20 }}
              onClick={cancelPressed}
            >
              Cancel
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
