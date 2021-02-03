import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import history from "./utils/history";

//importing auth
import { Auth0Provider } from "@auth0/auth0-react";

const onRedirectCallback = (appState) => {
  history.push(
    appState && appState.returnTo ? appState.returnTo : window.location.pathname
  );
};

ReactDOM.render(
  <Auth0Provider
    domain="dev-ps5dqqis.eu.auth0.com"
    clientId="DqbCvZtL8cn5plDla9TYlrJLhWIXpZtV"
    redirectUri={window.location.origin}
    onRedirectCallback={onRedirectCallback}
  >
    <App />
  </Auth0Provider>,
  document.getElementById("root")
);
