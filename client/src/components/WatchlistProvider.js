import { Container } from "@material-ui/core";
import WatchlistCard from "./WatchlistCard";

import { useQuery } from "react-query";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";
import { useAuth0 } from "@auth0/auth0-react";

export default function WatchlistProvider() {
  const graphqlClient = useGraphqlClient();
  const useGetWatchlistMovies = (userId) => {
    return useQuery(["getWatchlistMovies", userId], async () => {
      const data = (await graphqlClient).request(
        gql`
          query {
            watchlist(userId:"${userId}"){
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
  const { user } = useAuth0();
  const { isLoading, isError, data } = useGetWatchlistMovies(user?.sub);

  return (
    <Container maxWidth="sm" style={{ marginTop: "2.5vh" }}>
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
