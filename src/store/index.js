/* eslint-disable no-unused-vars */
import { createStore } from "redux";
import counter from "../reducer/index";
import signInStatus from "../reducer/index";

const store = createStore(counter, signInStatus);
