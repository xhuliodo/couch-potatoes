import { Container } from "@material-ui/core";
import WatchlistCard from "./WatchlistCard";

import { useQuery } from "react-query";
import { useMutation } from "react-query";
import { rateMovie } from "../utils/rateMovie";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";

export default function WatchlistProvider() {
  const graphqlClient = useGraphqlClient();
  const useGetWatchlistMovies = () => {
    return useQuery(["getWatchlistMovies"], async () => {
      const data = (await graphqlClient).request(
        gql`
          query {
            watchlist {
              movieId
              title
              posterUrl
              releaseYear
              imdbLink
            }
          }
        `
      );
      const { watchlist } = await data;
      return watchlist;
    });
  };
  const { isLoading, isError, data } = useGetWatchlistMovies();

  const rate = useMutation((mutationData) =>
    rateMovie(mutationData, graphqlClient)
  );

  const handleRate = (movieId, action) => {
    const mutationData = {
      movieId,
      action,
      successFunc: () => successFunc(movieId),
    };
    rate.mutate(mutationData);
  };

  const successFunc = (movieId) => {
    for (var i = 0; i < data.length; i++) {
      if (data[i].movieId === movieId) {
        data.splice(i, 1);
        i--;
      }
    }
  };

  return (
    <Container maxWidth="sm" style={{ marginTop: "15px" }}>
      {isLoading ? (
        <span>fetching data</span>
      ) : isError ? (
        <span>Something went wrong...</span>
      ) : (
        data.map((m) => (
          <WatchlistCard key={m.movieId} m={m} handleRate={handleRate} />
        ))
      )}
    </Container>
  );
}
