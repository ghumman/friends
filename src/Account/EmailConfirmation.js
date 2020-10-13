import React, {useState, useEffect} from "react";
import { connect } from "react-redux";


import Button from "react-bootstrap/Button";






function EmailConfirmation(props) {

  const [serverMessage, setServerMessage] = useState(""); 

  useEffect(() => {
    var url = window.location.href;
    var localSelectorEmail, localValidatorEmail;
    var regex = /[?&]([^=#]+)=([^&#]*)/g,
      params = {},
      match;
    while ((match = regex.exec(url))) {
      params[match[1]] = match[2];
      // console.log(match[1], match[2]);
      if (match[1] === "selectorEmail") {
        // console.log("match1 is selectorEmail");
        localSelectorEmail = match[2];
        // console.log("localSelectorEmail");
        // console.log(localSelectorEmail);
      }
      if (match[1] === "validatorEmail") {
        // console.log("match1 is validatorEmail");
        localValidatorEmail = match[2];
        // console.log("localValidatorEmail");
        // console.log(localValidatorEmail);
      }
    } // while ends

    if (
      (localSelectorEmail !== "") &
      (localSelectorEmail !== null) &
      ((localValidatorEmail !== "") & (localValidatorEmail !== null))
    ) {
      // console.log(
      //   "Inside if where localSelectorEmail and localValidatorEmail is not empty"
      // );
      try {
        fetch(
          "https://localhost/emailConfirmation",
          {
            method: "POST",
            headers: {
              "Content-Type": "application/x-www-form-urlencoded"
            },
            body:
              "selectorEmail=" +
              localSelectorEmail +
              "&validatorEmail=" +
              localValidatorEmail
          }
        ).then(async function(data) {
          // console.log("data: ");
          // console.log(data);
          data.json().then(async function(data) {
            // console.log("data.message:");
            // console.log(data.message);
            if (data.message === "Email confirmed successfully") {

              setServerMessage("Thank you for confirming your email. You are all set.");
              // Alert.alert("Email confirmed successfully.");
            } else {
              setServerMessage(data.message)
              // console.log(data.message);
              // Alert.alert("Unable to confirm email. Please try again.");
              // Alert.alert(data.message);
            }
          });
        });
      } catch (err) {
        // console.log(err);
        setServerMessage(err);
        // Alert.alert("Inside catch err");
        // Alert.alert(err);
      }
    } else {
      setServerMessage("Url is tempered" );
    }
  }, [])
    
  const mainApplication = () => {
    props.history.push({
      pathname: "/Profile"
    });
  }

  return (
    <div className="App">
      <header className="App-header">
        <p>Friends</p>
        <p>Email Confirmation</p>
        <p style={{ color: "red" }}>{serverMessage}</p>
        <Button onClick={() => mainApplication()}>
          Back to Web Application
        </Button>
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
)(EmailConfirmation);
