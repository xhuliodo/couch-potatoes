import request, { gql } from "graphql-request";

export const rateMovie = async ({ movieId, userId, action, successFunc }) => {
  let data = { rateMovie: null };
  // eslint-disable-next-line default-case
  switch (action) {
    case "love":
      data = await request(
        "http://localhost:4001/graphql",
        gql`
            mutation {
              rateMovie(movieId: "${movieId}", userId: "${userId}", rating:1) {
                movieId
                title
              }
            }
          `
      );
      break;
    case "hate":
      data = await request(
        "http://localhost:4001/graphql",
        gql`
            mutation {
              rateMovie(movieId: "${movieId}", userId: "${userId}", rating:0) {
                movieId
                title
              }
            }
          `
      );
      break;
  }
  const { rateMovie } = await data;
  if (rateMovie === null) {
    // TODO: you should implement some error logic, when the rating did not happen
    console.log("the rating didn't do shit");
  } else {
    console.log("you rated the movie with id", rateMovie.movieId);
    successFunc();
  }
};
