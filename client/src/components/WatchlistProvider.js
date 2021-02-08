import { Container } from "@material-ui/core";
import WatchlistCard from "./WatchlistCard";

import { useQuery } from "react-query";
import request, { gql } from "graphql-request";

const useGetWatchlistMovies = (userId) => {
  return useQuery(["getWatchlistMovies", userId], async () => {
    const data = request(
      "http://localhost:4001/graphql",
      gql`query {watchlist(userId:"${userId}"){movieId title posterUrl releaseYear imdbLink } } `
    );
    const { watchlist } = await data;
    return watchlist;
  });
};

export default function WatchlistProvider() {
  const { isLoading, isError, data } = useGetWatchlistMovies("2");

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
