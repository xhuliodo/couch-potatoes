import { gql } from "graphql-request";

export const rateMovie = async (
  { movieId, action, successFunc },
  graphqlClient
) => {
  let data = { rateMovie: null };
  switch (action) {
    case "love":
      data = (await graphqlClient).request(
        gql`
            mutation {
              rateMovie(movieId: ${movieId}, rating:1) {
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
              rateMovie(movieId: ${movieId}, rating:0.1) {
                movieId
                title
              }
            }
          `
      );
      break;
    default:
      console.log("are you lost? check the util before using it");
  }
  const { rateMovie } = await data;
  if (rateMovie === null) {
    // TODO: you should implement some feedback logic, when the rating did not happen
    console.log("the rating didn't do shit");
  } else {
    console.log("you rated the movie with id", rateMovie.movieId);
    successFunc();
  }
};
