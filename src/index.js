import React from "react";

import { Route, HashRouter as Router } from "react-router-dom";
import ReactDOM from "react-dom";
import "./index.css";
import * as serviceWorker from "./serviceWorker";

import Login from "./Account/Login";
import Forgot from "./Account/Forgot";
import ResetPassword from "./Account/ResetPassword";
import Change from "./Account/Change";
import Register from "./Account/Register";
import Profile from "./Message/Profile";

import { Provider } from "react-redux";

import { createStore } from "redux";

// eslint-disable-next-line no-unused-vars
const authentication = (state = [], action) => {
  switch (action.type) {
    case "UPDATE":
      return {
        loggedIn: true,
        user: action.username,
        userPassword: action.password
      };
    case "LOGOUT":
      return {
        loggedIn: false,
        user: "",
        userPassword: ""
      };

    default:
      return {
        loggedIn: false,
        user: "",
        userPassword: ""
      };
  }
};

const store = createStore(authentication);

const routing = (
  <Provider store={store}>
    <Router>
      <div>
        <Route exact path="/" component={Login} />
        <Route exact path="/Register" component={Register} />
        <Route path="/Forgot" component={Forgot} />
        <Route path="/reset-password" component={ResetPassword} />
        <Route path="/Change" component={Change} />
        <Route path="/Profile" component={Profile} />
      </div>
    </Router>
  </Provider>
);

const app = document.getElementById("root");
ReactDOM.render(routing, app);

// ReactDOM.render(
// <Provider store={store}>
// <App />
// </Provider>,
// document.getElementById("app")
// )

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
