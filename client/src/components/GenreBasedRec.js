import request, { gql } from "graphql-request";
import { useQuery } from "react-query";

import { useMovieStore } from "../context/movies";

import MovieCard from "./MovieCard";
import UserFeedbackMovieCard from "./UserFeedbackMovieCard";

export default function GenreBasedRec({
  skip,
  setSkip,
  startedFromTheBottomNowWeHere = false,
}) {
  const useMovies = (userId, limit) => {
    return useQuery(["movies", userId, limit], async () => {
      const data = await request(
        "http://localhost:4001/graphql",
        gql`
                query {
                    recommendPopularMoviesBasedOnGenre(
                      userId: ${userId}
                      limit: ${limit}
                      skip: ${skip}
                    ) {
                      movieId
                      posterUrl
                      title
                      releaseYear
                      imdbLink
                    }
                }
              `
      );
      const { recommendPopularMoviesBasedOnGenre } = await data;
      return recommendPopularMoviesBasedOnGenre;
    });
  };

  const { limit } = useMovieStore();

  const { isLoading, isError, data, refetch } = useMovies(
    "2",
    limit
  );

  return isLoading ? (
    <UserFeedbackMovieCard message="Fetching movies..." type={"loading"} />
  ) : isError ? (
    <UserFeedbackMovieCard message="Something went wrong..." />
  ) : (
    <MovieCard
      skip={skip}
      setSkip={setSkip}
      refetch={refetch}
      movies={data}
      startedFromTheBottomNowWeHere={startedFromTheBottomNowWeHere}
    />
  );
}
