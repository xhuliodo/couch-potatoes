import { useEffect, useState } from "react";

import {
  Button,
  Container,
  Paper,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
} from "@material-ui/core";
import "../components/MovieCard.css";

import { useQuery } from "react-query";
import request, { gql } from "graphql-request";

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
              genre_1: ${genre_1}
              genre_2: ${genre_2} 
              genre_3: ${genre_3}
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

export default function GettingToKnowUserPage(props) {
  const genres = ["1", "2", "3"];
  const {
    movies,
    setMovies,
    nextMovie,
    skip,
    limit,
    increaseSkip,
    ratedMovies,
    requiredMovies,
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

  const [open, setOpen] = useState(false);
  const handleClose = () => {
    setOpen(false);
  };
  const handleNext = () => {
    handleClose();
    props.history.push("/recommendations");
  };

  return (
    <Paper elevation={0} style={{ height: "90vh" }}>
      <Typography style={{ textAlign: "center" }} variant="h6">
        Please rate at least {requiredMovies} movies{" "}
        <i>(ignored movies will not count)</i>
      </Typography>
      <Typography
        style={{ textAlign: "center", fontWeight: "bold" }}
        variant="h6"
      >
        You have rated {ratedMovies} / {requiredMovies}
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
            setOpen={setOpen}
          />
        )}
        <Container maxWidth="sm">
          <Button
            style={{ marginTop: "3.4vh" }}
            color="primary"
            fullWidth
            variant="contained"
            onClick={handleNext}
            disabled={ratedMovies >= requiredMovies ? false : true}
          >
            Continue
          </Button>
        </Container>

        <Dialog open={open} onClose={handleClose}>
          <DialogTitle>We've gotten to know you enough 😃</DialogTitle>
          <DialogContent>
            <DialogContentText>
              You can either click next to get recommendations based on other
              users ratings or continue to rate movies in the genre you selected
              as well (you can always click the continue button at the end to go
              to the next step)
            </DialogContentText>
            <DialogActions>
              <Button onClick={handleClose} color="secondary">
                Continue
              </Button>
              <Button onClick={handleNext} color="primary">
                Next
              </Button>
            </DialogActions>
          </DialogContent>
        </Dialog>
      </Container>
    </Paper>
  );
}
