import { gql } from "graphql-request";

export const removeFromWatchlist = async (
  { movieId, successFunc },
  graphqlClient
) => {
  const data = (await graphqlClient).request(
    gql`
      mutation {
        removeFromWatchlist(movieId: ${movieId}) {
          movieId
        }
      }
    `
  );

  const { removeFromWatchlist } = await data;
  if (removeFromWatchlist?.movieId) {
    console.log(
      "you removed from watchlist the movie with id",
      removeFromWatchlist.movieId
    );
    successFunc();
  } else {
    // TODO: you should implement some feedback logic, when the rating did not happen
    console.log("you shall not remove");
  }
};
