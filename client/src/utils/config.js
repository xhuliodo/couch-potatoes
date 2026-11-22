export const auth0Domain = process.env.REACT_APP_AUTH0_DOMAIN;
export const auth0ClientId = process.env.REACT_APP_AUTH0_CLIENT_ID;

export const redirectBase =
  process.env.REACT_APP_REDIRECT_BASE || "http://localhost:3000";

export const apiEndpoint =
  process.env.REACT_APP_API_ENDPOINT || "http://localhost:4001";

export const thumbnailsEndpoint =
  process.env.REACT_APP_THUMBNAILS_ENDPOINT || "http://localhost:4003";

export const thumbnailUrl = (movieId) =>
  `${thumbnailsEndpoint}/thumbnail_${movieId}.jpg`;
