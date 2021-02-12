import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";
import { useQuery } from "react-query";

import { useMovieStore } from "../context/movies";

import MovieCard from "./MovieCard";
import UserFeedbackMovieCard from "./UserFeedbackMovieCard";

export default function GenreBasedRec({
  skip = 0,
  setSkip,
  startedFromTheBottomNowWeHere = false,
}) {
  const graphqlClient = useGraphqlClient();
  const useMovies = (limit) => {
    return useQuery(["movies", limit], async () => {
      const data = (await graphqlClient).request(
        gql`
          query {
            recommendPopularMoviesBasedOnGenre(
              limit: ${limit}, 
              skip: ${skip}) {
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

  const { isLoading, isError, data, refetch } = useMovies(limit);

  return isLoading ? (
    <UserFeedbackMovieCard message="Fetching movies..." type="loading" />
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
