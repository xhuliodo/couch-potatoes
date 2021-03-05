import { useState } from "react";

import { useMovieStore } from "../context/movies";
import { useQuery } from "react-query";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";
import WatchlistProvider from "./WatchlistProvider";

export default function WatchlistUnrated() {
  // const classes = useStyle();

  // get data
  const { limit } = useMovieStore();
  const graphqlClient = useGraphqlClient();
  const [skip, setSkip] = useState(0);
  const { isLoading, isError, data } = useGetWatchlistHistory({
    graphqlClient,
    limit,
    skip,
  });

  return (
    <WatchlistProvider
      graphqlClient={graphqlClient}
      isLoading={isLoading}
      isError={isError}
      data={data}
      skip={skip}
      setSkip={setSkip}
      limit={limit}
    />
  );
}

const useGetWatchlistHistory = ({ graphqlClient, skip, limit }) => {
  return useQuery(["getWatchlistHistory", skip, limit], async () => {
    const data = (await graphqlClient).request(
      gql`
          query {
            watchlistHistory(skip:${skip}, limit:${limit}) {
              movieId
              title
              releaseYear
              imdbLink
              rating
            }
          }
        `
    );
    const { watchlistHistory } = await data;
    return watchlistHistory;
  });
};
