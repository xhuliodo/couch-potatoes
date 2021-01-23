import { useEffect } from "react";

import { Container, Paper, Typography } from "@material-ui/core";
import "../components/MovieCard.css";

import { useQuery } from "react-query";
import request, { gql } from "graphql-request";

import { useGenreStore } from "../context/genres";
import { useMovieStore } from "../context/movies";
import SecondMovieCard from "../components/SecondMovieCard";

const useMovies = (genre_1, genre_2, genre_3) => {
  return useQuery(["movies", genre_1, genre_2, genre_3], async () => {
    const data = await request(
      "http://localhost:4001/graphql",
      gql`
        query {
            recommendPopularMoviesBasedOnGenre(
              genre_1: "${genre_1}" 
              genre_2: "${genre_2}" 
              genre_3: "${genre_3}" 
              limit: 30
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

export default function GettingToKnowUserPage() {
  const genres = useGenreStore((state) => state.genres);

  const { status, data = [], error } = useMovies(
    genres[0],
    genres[1],
    genres[2]
  );
  const { movies, setMovies, nextMovie } = useMovieStore();

  useEffect(() => {
    status === "success" ? setMovies(data) : console.log("waiting for data");
  }, [data, setMovies, status]);

  return (
    <Paper elevation={0}>
      <Typography variant="h5">
        Please rate at least 15 movies <i>(ignored movies will not count)</i>
      </Typography>
      <Container
        disableGutters={true}
        // style={{ marginTop: "2.5vh", minHeight: "80vh" }}
      >
        {status === "loading" ? (
          <span>Fetching data</span>
        ) : status === "error" ? (
          <span>Error: {error.message}</span>
        ) : (
          <SecondMovieCard movies={movies} nextMovie={nextMovie} />
        )}
      </Container>
    </Paper>
  );
}
