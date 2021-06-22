import { useAxiosClient } from "../utils/useAxiosClient";
import { useQuery } from "react-query";

import { useMovieStore } from "../context/movies";

import MovieCard from "./MovieCard";
import DataStateMovieCard from "./DataStateMovieCard";

export default function GenreBasedRec({ skip, setSkip }) {
  const { limit } = useMovieStore();

  const axiosClient = useAxiosClient();

  const { isLoading, isError, data, refetch } = useMovies({
    axiosClient,
    limit,
    skip,
  });

  return isLoading ? (
    <DataStateMovieCard message="Fetching movies..." type="loading" />
  ) : isError ? (
    <DataStateMovieCard message="Something went wrong..." />
  ) : (
    <MovieCard skip={skip} setSkip={setSkip} refetch={refetch} movies={data} />
  );
}

const useMovies = ({ axiosClient, limit, skip }) => {
  return useQuery("movies", async () => {
    const resp = (await axiosClient).get(
      `/recommendations/popular?limit=${limit}&skip=${skip}`
    );

    const {
      data: { data: popularMovie },
    } = await resp;

    return popularMovie;
  });
};
