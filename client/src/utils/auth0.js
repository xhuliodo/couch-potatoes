import { auth0ClientId, auth0Domain, redirectBase } from "./config";

export const authConfig = {
  domain: auth0Domain,
  clientId: auth0ClientId,
  redirectUri: `${redirectBase}/solo`,
  logoutUri: `${redirectBase}/`,
};
