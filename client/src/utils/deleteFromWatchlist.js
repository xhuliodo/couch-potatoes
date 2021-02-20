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
    console.log("you shall not remove");
  }
};
