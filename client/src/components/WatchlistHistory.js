import { useState } from "react";
import { useMovieStore } from "../context/movies";
import { useQuery } from "react-query";
import { useAxiosClient } from "../utils/useAxiosClient";
import WatchlistProvider from "./WatchlistProvider";

export default function WatchlistUnrated() {
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
      data={data?.watchlistHistory}
      statusCode={data?.statusCode}
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
      let ulala;
      try {
        const resp = (await axiosClient).get(
          `/watchlist/history?limit=${limit}&skip=${skip}`
        );
        const {
          data: { data: watchlistHistory, statusCode },
        } = await resp;
        ulala = {
          watchlistHistory,
          statusCode,
        };
      } catch (e) {
        const {
          response: {
            data: { statusCode },
          },
        } = e;
        ulala = { watchlistHistory: [], statusCode };
      }
      return ulala;
    },
    { cacheTime: 0 }
  );
};
