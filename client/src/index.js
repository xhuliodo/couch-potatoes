import ReactDOM from "react-dom";
import App from "./App";

//importing auth
import { Auth0Provider } from "@auth0/auth0-react";
import history from "./utils/history";

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
