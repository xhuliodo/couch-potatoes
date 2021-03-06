import { useEffect } from "react";
import { Button, Grid, Paper, Typography } from "@material-ui/core";
import {
  SentimentDissatisfied,
  SentimentVerySatisfiedRounded,
  SkipNextOutlined,
  WatchLaterOutlined,
} from "@material-ui/icons";
import { Card, CardWrapper } from "@xhuliodo/react-swipeable-cards";
import { openRateFeedbackExported } from "./RateFeedback";
import RateFeedback from "./RateFeedback";
import { waitingForMore } from "../utils/waitingForMore";
import { useMutation } from "react-query";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { rateMovie } from "../utils/rateMovie";
import { gql } from "graphql-request";
import { useMovieStore } from "../context/movies";

export default function MovieCard({
  compName = "gtk",
  movies,
  skip = 0,
  setSkip = null,
  refetch,
  startedFromTheBottomNowWeHere = false,
}) {
  const { increaseRatedMovies } = useMovieStore();

  const graphqlClient = useGraphqlClient();

  const rate = useMutation((mutationData) =>
    rateMovie(mutationData, graphqlClient)
  );

  const handleRate = (action) => {
    if (movies.length === 0) {
      console.log("patience is a virtue");
    } else {
      const mutationData = {
        movieId: movies[0].movieId,
        action,
        successFunc: () => successFunc(action),
      };
      rate.mutate(mutationData);
    }
  };
  const addToWatchlist = useMutation(async ({ movieId }) => {
    if (movies.length === 0) {
      console.log("patience is a virtue");
    } else {
      const data = (await graphqlClient).request(
        gql`
        mutation {
          addToWatchlist(movieId: ${movieId}) {
            movieId
          }
        }
      `
      );

      const { addToWatchlist } = await data;
      if (addToWatchlist === null) {
        console.log("the watchlist did not get filled 😏");
      } else {
        console.log(
          "you added to playlist the movie with id",
          addToWatchlist.movieId
        );
        openRateFeedbackExported("watchlater");
        nextMovie();
      }
    }
  });
  const successFunc = (action) => {
    openRateFeedbackExported(action);
    nextMovie();
    increaseRatedMovies();
  };
  const skipMovie = () => {
    setSkip(skip + 1);
    nextMovie();
    openRateFeedbackExported("skip");
  };
  const nextMovie = () => {
    movies.shift();
  };

  useEffect(() => {
    if (movies.length === 0) {
      refetch();
    }
  }, [skip, setSkip, movies, rate, refetch]);

  const ActionButtons = () => {
    const newProps = {};

    return (
      <Grid
        {...newProps}
        style={{
          position: "absolute",
          bottom: "15px",
          width: "100vw",
          alignItems: "baseline",
        }}
        container
        justify="center"
        className="cards_container"
      >
        <div>
          <Button
            className="actionButton"
            onClick={() => handleRate("hate")}
            variant="contained"
            color="secondary"
          >
            <SentimentDissatisfied fontSize="inherit" />
            <Typography align="center" style={{ fontSize: "0.8rem" }}>
              Not for me
            </Typography>
          </Button>
          {/* <Typography align="center">Not for me</Typography> */}
        </div>

        {startedFromTheBottomNowWeHere ? (
          <div>
            <Button
              className="actionButton"
              onClick={() => {
                const mutationData = {
                  movieId: movies[0]?.movieId,
                };
                addToWatchlist.mutate(mutationData);
              }}
              variant="contained"
            >
              <WatchLaterOutlined fontSize="inherit" />
              <Typography align="center" style={{ fontSize: "0.8rem" }}>
                Watch later
              </Typography>
            </Button>
          </div>
        ) : (
          <div>
            <Button
              className="actionButton"
              onClick={skipMovie}
              variant="contained"
            >
              <SkipNextOutlined fontSize="inherit" />
              <Typography align="center" style={{ fontSize: "0.8rem" }}>
                Skip
              </Typography>
            </Button>
            {/* <Typography align="center">Skip</Typography> */}
          </div>
        )}
        <div>
          <Button
            className="actionButton"
            onClick={() => handleRate("love")}
            variant="contained"
            color="primary"
          >
            <SentimentVerySatisfiedRounded fontSize="inherit" />
            <Typography align="center" style={{ fontSize: "0.8rem" }}>
              Loved it
            </Typography>
          </Button>
          {/* <Typography align="center">Loved it</Typography> */}
        </div>
      </Grid>
    );
  };

  return (
    <CardWrapper addEndCard={waitingForMore} style={{ paddingTop: "0px" }}>
      {movies.map((m) => (
        <Card
          key={compName + m.movieId}
          onSwipeLeft={() => handleRate("hate")}
          onSwipeRight={() => handleRate("love")}
          swipeSensitivity={50}
          style={{
            backgroundImage: `url(https://thumb.cp.dev.cloudapp.al/thumbnail_${m.movieId}.jpg)`,
            backgroundSize: "contain",
            backgroundRepeat: "no-repeat",
            backgroundPosition: "center",
          }}
        >
          <MovieTitle title={m.title} releaseYear={m.releaseYear} />
        </Card>
      ))}
      <RateFeedback />
      <ActionButtons />
    </CardWrapper>
  );
}

const MovieTitle = ({ title, releaseYear }) => {
  const newProps = {};
  return (
    <Paper {...newProps} className="movie_cardText">
      <h3>
        {title} <i>({releaseYear})</i>
      </h3>
    </Paper>
  );
};
