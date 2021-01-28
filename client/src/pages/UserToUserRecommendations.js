import { Button, Container, Paper, Typography } from "@material-ui/core";

import "../components/MovieCard.css";

import { useQuery } from "react-query";
import request, { gql } from "graphql-request";

import { useMovieStore } from "../context/movies";
import MovieCard from "../components/MovieCard";
import UserFeedbackMovieCard from "../components/UserFeedbackMovieCard";
import { useEffect } from "react";

const useUserBasedRec = (
  userId,
  minimumRatings,
  peopleToCompare,
  moviesToRecommend
) => {
  return useQuery(
    [
      "userBasedRec",
      userId,
      minimumRatings,
      peopleToCompare,
      moviesToRecommend,
    ],
    async () => {
      const data = await request(
        "http://localhost:4001/graphql",
        gql`
      query {
        
  recommendFromOtherUsers(
    userId: ${userId}
    minimumRatings: ${minimumRatings}
    peopleToCompare: ${peopleToCompare}
    moviesToRecommend: ${moviesToRecommend}
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
      const { recommendFromOtherUsers } = await data;
      return recommendFromOtherUsers;
    }
  );
};

export default function UserToUserRecommendations() {
  const {
    peopleToCompare,
    limit,
    movies,
    setMovies,
    nextMovie,
    requiredMovies,
  } = useMovieStore();

  const minimumRatings = requiredMovies / 1.5;

  const { isLoading, isError, data = [], error, refetch } = useUserBasedRec(
    "2",
    minimumRatings,
    peopleToCompare,
    limit
  );

  useEffect(() => {
    isLoading
      ? console.log("waiting for data")
      : isError
      ? console.log(error.message)
      : setMovies(data);
  }, [data, setMovies, isLoading, isError, error]);

  return (
    <Paper elevation={0} style={{ height: "90vh" }}>
      <Typography style={{ textAlign: "center" }} variant="h6">
        Get recommendations based on other users
      </Typography>
      <Container disableGutters={true}>
        {isLoading ? (
          <UserFeedbackMovieCard
            message={"Fetching movies..."}
            type={"loading"}
          />
        ) : isError ? (
          <UserFeedbackMovieCard message={"Something went wrong..."} />
        ) : (
          <MovieCard
            startedFromTheBottomNowWeHere={true}
            movies={movies}
            nextMovie={nextMovie}
            refetch={refetch}
          />
        )}
        <Container maxWidth="sm">
          <Button
            style={{ marginTop: "3.4vh" }}
            color="primary"
            fullWidth
            variant="contained"
            // onClick={handleNext}
          >
            Add to watchlist
          </Button>
        </Container>
      </Container>
    </Paper>
  );
}
