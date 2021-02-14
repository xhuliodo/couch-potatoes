import { Container } from "@material-ui/core";
import WatchlistCard from "./WatchlistCard";

import { useQuery } from "react-query";
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

  return (
    <Container maxWidth="sm" style={{ marginTop: "15px" }}>
      {isLoading ? (
        <span>fetching data</span>
      ) : isError ? (
        <span>Something went wrong...</span>
      ) : (
        <WatchlistCard movies={data} />
      )}
    </Container>
  );
}
