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
      isLoading={isLoading}
      isError={isError}
      data={data}
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
      (await axiosClient)
        .get(`/watchlist?limit=${limit}&skip=${skip}`)
        .then((resp) => {
          const { data } = resp;
          return data;
        })
        .catch((err) => {
          console.log(err);
        });
      // .request(
      //   gql`
      //   query {
      //     watchlist(skip:${skip}, limit:${limit}) {
      //       movieId
      //       title
      //       releaseYear
      //       imdbLink
      //     }
      //   }
      // `
      // );
      // const { watchlist } = await data;
      // return watchlist;
    },
    { cacheTime: 0 }
  );
};
