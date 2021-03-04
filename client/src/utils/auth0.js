const redirectUri = process.env.REACT_APP_DOMAIN || "http://localhost:3000";

export const authConfig = {
  domain: "dev-ps5dqqis.eu.auth0.com",
  clientId: "DqbCvZtL8cn5plDla9TYlrJLhWIXpZtV",
  redirectUri: `${redirectUri}/solo`,
  logoutUri: `${redirectUri}/`,
};
