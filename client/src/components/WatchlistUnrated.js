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
  const { isLoading, isError, data } = useGetWatchlistMovies({
    axiosClient,
    limit,
    skip,
  });

  return (
    <WatchlistProvider
      axiosClient={axiosClient}
      statusCode={data?.statusCode}
      isLoading={isLoading}
      isError={isError}
      data={data?.watchlist}
      skip={skip}
      setSkip={setSkip}
      limit={limit}
    />
  );
}

const useGetWatchlistMovies = ({ axiosClient, skip, limit }) => {
  return useQuery(
    ["getWatchlistMovies", skip, limit],
    async () => {
      let ulala;
      try {
        const resp = (await axiosClient).get(
          `/watchlist?limit=${limit}&skip=${skip}`
        );
        const {
          data: { data: watchlist, statusCode },
        } = await resp;
        ulala = {
          watchlist,
          statusCode,
        };
      } catch (e) {
        const {
          response: {
            data: { statusCode },
          },
        } = e;
        ulala = { watchlist: [], statusCode };
      }
      return ulala;
    },
    { cacheTime: 0 }
  );
};
