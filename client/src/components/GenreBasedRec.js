import { useAuth0 } from "@auth0/auth0-react";
import request, { gql } from "graphql-request";
import { useQuery } from "react-query";

import { useMovieStore } from "../context/movies";

import MovieCard from "./MovieCard";
import UserFeedbackMovieCard from "./UserFeedbackMovieCard";

export default function GenreBasedRec({
  skip = 0,
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
              userId: "${userId}", 
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

  const { user } = useAuth0();

  const { isLoading, isError, data, refetch } = useMovies(user.sub, limit);

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
