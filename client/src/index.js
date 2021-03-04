import ReactDOM from "react-dom";
import App from "./App";

//importing auth
import { Auth0Provider } from "@auth0/auth0-react";
import history from "./utils/history";
import { authConfig } from "./utils/auth0";

const onRedirectCallback = (appState) => {
  history.push(
    appState && appState.returnTo ? appState.returnTo : window.location.pathname
  );
};

ReactDOM.render(
  <Auth0Provider
    domain={authConfig.domain}
    clientId={authConfig.clientId}
    redirectUri={authConfig.redirectUri}
    onRedirectCallback={onRedirectCallback}
  >
    <App />
  </Auth0Provider>,
  document.getElementById("root")
);
