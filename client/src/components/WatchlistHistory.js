import { useState } from "react";

import { useMovieStore } from "../context/movies";
import { useQuery } from "react-query";
import { useAxiosClient } from "../utils/useAxiosClient";
import WatchlistProvider from "./WatchlistProvider";

export default function WatchlistUnrated() {
  // const classes = useStyle();

  // get data
  const { limit } = useMovieStore();
  const axiosClient = useAxiosClient();
  const [skip, setSkip] = useState(0);
  const { isLoading, isError, data } = useGetWatchlistHistory({
    axiosClient,
    limit,
    skip,
  });

  return (
    <WatchlistProvider
      axiosClient={axiosClient}
      isLoading={isLoading}
      isError={isError}
      data={data}
      skip={skip}
      setSkip={setSkip}
      limit={limit}
    />
  );
}

const useGetWatchlistHistory = ({ axiosClient, skip, limit }) => {
  return useQuery(
    ["getWatchlistHistory", skip, limit],
    async () => {
      (await axiosClient)
        .get(`/watchlist/history?limit=${limit}&skip=${skip}`)
        .then((resp) => {
          const { data } = resp;
          return data;
        })
        .catch((err) => {
          console.log(err);
        });
      // .request(
      //   gql`
      //     query {
      //       watchlistHistory(skip:${skip}, limit:${limit}) {
      //         movieId
      //         title
      //         releaseYear
      //         imdbLink
      //         rating
      //       }
      //     }
      //   `
      // );
      // const { watchlistHistory } = await data;
      // return watchlistHistory;
    },
    { cacheTime: 0 }
  );
};
