import { useEffect } from "react";

import { Button, Container, Grid, Paper, Typography } from "@material-ui/core";
import { VisibilityOff } from "@material-ui/icons";
import "../components/MovieCard.css";

import { useQuery } from "react-query";
import request, { gql } from "graphql-request";

import { useGenreStore } from "../context/genres";
import { useMovieStore } from "../context/movies";
import SecondMovieCard from "../components/SecondMovieCard";

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

  const { status, data = [], error } = useMovies(
    genres[0],
    genres[1],
    genres[2],
    skip,
    limit
  );

  useEffect(() => {
    status === "success" ? setMovies(data) : console.log("waiting for data");
  }, [data, setMovies, status, skip]);

  return (
    <Paper elevation={0}>
      <Typography variant="h5">
        Please rate at least ${requiredMovies} movies
        <i>(ignored movies will not count)</i>
      </Typography>
      <Container disableGutters={true}>
        {status === "loading" ? (
          <span>Fetching data</span>
        ) : status === "error" ? (
          <span>Error: {error.message}</span>
        ) : (
          <SecondMovieCard
            increaseSkip={increaseSkip}
            movies={movies}
            nextMovie={nextMovie}
          />
        )}
        <Grid container justify="center">
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
        </Grid>
      </Container>
    </Paper>
  );
}
