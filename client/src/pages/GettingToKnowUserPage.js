import { useEffect } from "react";

import { Button, Container, Grid, Paper, Typography } from "@material-ui/core";
import { VisibilityOff } from "@material-ui/icons";
import "../components/MovieCard.css";

import { useQuery } from "react-query";
import request, { gql } from "graphql-request";

import { useGenreStore } from "../context/genres";
import { useMovieStore } from "../context/movies";
import SecondMovieCard from "../components/SecondMovieCard";
import UserFeedbackMovieCard from "../components/UserFeedbackMovieCard";

const useMovies = (genre_1, genre_2, genre_3, skip, limit) => {
  return useQuery(
    ["movies", genre_1, genre_2, genre_3, skip, limit],
    async () => {
      const data = await request(
        "http://localhost:4001/graphql",
        gql`
        query {
            recommendPopularMoviesBasedOnGenre(
              genre_1: "${genre_1}" 
              genre_2: "${genre_2}" 
              genre_3: "${genre_3}" 
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
    }
  );
};

export default function GettingToKnowUserPage() {
  const genres = useGenreStore((state) => state.genres);
  const {
    movies,
    setMovies,
    nextMovie,
    skip,
    limit,
    requiredMovies,
    increaseSkip,
  } = useMovieStore();

  const { isLoading, isError, data = [], error } = useMovies(
    genres[0],
    genres[1],
    genres[2],
    skip,
    limit
  );

  useEffect(() => {
    isLoading
      ? console.log("waiting for data")
      : isError
      ? console.log(error.message)
      : setMovies(data);
  }, [data, setMovies, isLoading, isError, error, skip]);

  return (
    <Paper elevation={0} style={{ height: "90vh" }}>
      <Typography variant="h4">
        Please rate at least {requiredMovies} movies
        <i>(ignored movies will not count)</i>
      </Typography>
      <Container disableGutters={true}>
        {isLoading ? (
          <UserFeedbackMovieCard message={"Fetching movies..."} />
        ) : isError ? (
          <UserFeedbackMovieCard message={"Something went wrong..."} />
        ) : (
          <SecondMovieCard
            increaseSkip={increaseSkip}
            movies={movies}
            nextMovie={nextMovie}
          />
        )}
        {/* <Grid container justify="center">
          <Button
            size="large"
            onClick={() => {
              nextMovie();
            }}
            startIcon={<VisibilityOff />}
            variant="contained"
          >
            Have not seen
          </Button>
        </Grid> */}
      </Container>
    </Paper>
  );
}
