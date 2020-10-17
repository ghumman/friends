import React, {useEffect} from "react";
import { connect } from "react-redux";


import Button from "react-bootstrap/Button";



const USERNAME = "username";
const PASSWORD = "password";
const LOGIN_TYPE = "login_type";
const ACCOUNT_TYPE = "account_type";
const EMAIL = "email";


function Profile(props) {

  useEffect(() => {
    if (props.authenticate.user.trim() === "") {
      props.history.push({
        pathname: "/"
      });
    }
  }, [])

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

  return (
    <div className="App">
      <header className="App-header">
        <div style={{ padding: 10 }}>
          <Button onClick={() => signMeOut()}>SIGN OUT</Button>
        </div>
        {changePassword}
        
        <p>Friends</p>
        <p>Hi {props.authenticate.user}</p>

        
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

