import { gql } from "graphql-request";

export const rateMovie = async (
  { movieId, action, successFunc },
  graphqlClient
) => {
  let data = { rateMovie: null };
  // eslint-disable-next-line default-case
  switch (action) {
    case "love":
      data = (await graphqlClient).request(
        gql`
            mutation {
              rateMovie(movieId: "${movieId}", rating:1) {
                movieId
                title
              }
            }
          `
      );
      break;
    case "hate":
      data = (await graphqlClient).request(
        gql`
            mutation {
              rateMovie(movieId: "${movieId}", rating:0) {
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
