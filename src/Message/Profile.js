import React, { useEffect, useState } from "react";
import { connect } from "react-redux";

import Button from "react-bootstrap/Button";

import Form from "react-bootstrap/Form";

const USERNAME = "username";
const PASSWORD = "password";
const LOGIN_TYPE = "login_type";
const ACCOUNT_TYPE = "account_type";
const EMAIL = "email";

function Profile(props) {

  const [toEmail, setToEmail] = useState("");
  const [toMessage, setToMessage] = useState("");
  const [errorServerMessage, setErrorServerMessage] = useState("");
  const [friends, setFriends] = useState([]);
  const [friendMessages, setFriendMessages] = useState([]);


  const handleChangeToEmail = (event) => {
    setToEmail(event.target.value);
  }

  const handleChangeToMessage = (event) => {
    setToMessage(event.target.value);
  }

  const handleSubmit = (event) => {
    event.preventDefault();

    fetch(
      "http://localhost:8080/send-message",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        body:
          "message=" + toMessage +
          "&messageFromEmail=" + localStorage.getItem(USERNAME) +
          "&messageToEmail=" + toEmail.trim() +
          "&authType=regular" +
          "&password=" + localStorage.getItem(PASSWORD)
      }
    ).then(async function (data) {
      data.json().then(async function (data) {
        if (data.message === "Message sent") {
          setErrorServerMessage("Message sent successfully.");
        } else {
          setErrorServerMessage(data.message);
        }
      });
    });
  }


  useEffect(() => {
    if (props.authenticate.user.trim() === "") {
      props.history.push({
        pathname: "/"
      });
    }
    getFriends();
  }, [])

  const getFriends = () => {
    fetch(
      "http://localhost:8080/all-friends",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        body:
          "email=" + localStorage.getItem(USERNAME) +
          "&authType=regular" +
          "&password=" + localStorage.getItem(PASSWORD)
      }
    ).then(async function (data) {
      data.json().then(async function (data) {

        if (data.message === "Friends attached") {
          setFriends(data.usersAll)
        } else {
          setErrorServerMessage(data.message);
        }
      });
    });
  }

  const signMeOut = () => {
    localStorage.setItem(USERNAME, "");
    localStorage.setItem(PASSWORD, "");
    localStorage.setItem(LOGIN_TYPE, "");
    localStorage.setItem(ACCOUNT_TYPE, "");
    localStorage.setItem(EMAIL, "");
    props.signout();
    props.history.push({
      pathname: "/"
    });
  }

  const changePasswordScreen = () => {
    props.history.push({
      pathname: "/Change"
    });
  }

  const showFriendMessages = (friendEmail) => {
    fetch(
      "http://localhost:8080/messages-user-and-friend",
      {
        method: "POST",
        headers: {
          "Content-Type": "application/x-www-form-urlencoded"
        },
        body:
          "userEmail=" + localStorage.getItem(USERNAME) +
          "&friendEmail=" + friendEmail +
          "&authType=regular" +
          "&password=" + localStorage.getItem(PASSWORD)
      }
    ).then(async function (data) {
      data.json().then(async function (data) {

        if (data.message === "Messages attached") {
          setFriendMessages(data.msgs)
        } else {
          setErrorServerMessage(data.message);
        }
      });
    });
  }

  var changePassword;
  const resultType = localStorage.getItem("login_type");

  if (resultType === "special") {
    changePassword = null;
  } else {
    changePassword = (
      <div style={{ padding: 10 }}>
        <Button onClick={() => changePasswordScreen()}>
          CHANGE PASSWORD
        </Button>
      </div>
    );
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

  const friendsData = (
    <div>
      {friends.map((usr) =>
        <div key={usr.email}
          style={{
            borderStyle: "solid",
            borderColor: "white",
            borderRadius: 5,
            padding: 5,
            margin: 5
          }}>
          <li onClick={() => showFriendMessages(usr.email)}>
            {usr.firstName}{' '}{usr.lastName}{' ['}{usr.email}{']'}
          </li>
        </div>
      )}
    </div>
  );

  const friendMessagesData = (
    <div>
      {friendMessages.map((msg) =>
        <div key={msg.sentAt}
          style={{
            borderStyle: "solid",
            borderColor: "white",
            borderRadius: 5,
            padding: 5,
            margin: 5
          }}>
          <li>
            From: {'   '}{msg.messageFromEmail}{'    '}To: {'   '}{msg.messageToEmail}
          </li>
          <div>
            {msg.message}
          </div>
          <div>
            {msg.sentAt}
          </div>
        </div>
      )}
    </div>
  );

  return (
    <div className="App">
      <header className="App-header">
        {errorMessage}
        <div style={{ padding: 10 }}>
          <Button onClick={() => signMeOut()}>SIGN OUT</Button>
        </div>
        {changePassword}

        <p>Friends</p>
        <p>Hi {localStorage.getItem(USERNAME)}</p>

        <p
          style={{
            borderStyle: "solid",
            borderColor: "white",
            borderRadius: 5,
            padding: 5
          }}
        >
          Send Message
          </p>

        <Form onSubmit={handleSubmit}>
          <Form.Group controlId="toEmail">
            <Form.Label>Email</Form.Label>
            <Form.Control
              type="email"
              value={toEmail}
              placeholder="Enter Email"
              onChange={handleChangeToEmail}
            />
          </Form.Group>

          <Form.Group controlId="toMessage">
            <Form.Label>Message</Form.Label>
            <Form.Control
              as="textarea" rows="3"
              value={toMessage}
              onChange={handleChangeToMessage}
            />
          </Form.Group>

          <Button variant="primary" type="submit">
            SEND
            </Button>
        </Form>

        {friendsData}
        {friendMessagesData}

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
      dispatch({ type: "UPDATE", username: username, password: password }),
    signout: () => dispatch({ type: "LOGOUT" })
  };
}

export default connect(
  mapStateToProps,
  mapDispatchToProps
)(Profile);

