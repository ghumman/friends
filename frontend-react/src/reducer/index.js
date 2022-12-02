/* eslint-disable no-unused-vars */
const counter = (state = 0, action) => {
  switch (action.type) {
    case "INCREMENT":
      return (state = state + 1);

    case "DECREMENT":
      return (state = state - 1);
    default:
      return state;
  }
};

const signInStatus = (state = false, action) => {
  switch (action.type) {
    case "SIGNIN":
      return (state = true);

    case "SIGNOUT":
      return (state = false);
    default:
      return state;
  }
};
