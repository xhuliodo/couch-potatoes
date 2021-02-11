import React from "react";
import ReactDOM from "react-dom";
import "./index.css";
import App from "./App";
import history from "./utils/history";

// pwa service worker
import serviceWorker from "./serviceWorker";

//importing auth
import { Auth0Provider } from "@auth0/auth0-react";

const onRedirectCallback = (appState) => {
  history.push(
    appState && appState.returnTo ? appState.returnTo : window.location.pathname
  );
};

const redirectUri = process.env.REACT_APP_DOMAIN || "http://localhost:3000";

ReactDOM.render(
  <Auth0Provider
    domain="dev-ps5dqqis.eu.auth0.com"
    clientId="DqbCvZtL8cn5plDla9TYlrJLhWIXpZtV"
    redirectUri={`${redirectUri}/solo`}
    onRedirectCallback={onRedirectCallback}
  >
    <App />
  </Auth0Provider>,
  document.getElementById("root")
);

serviceWorker.register();
