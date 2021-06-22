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
import { useAxiosClient } from "../utils/useAxiosClient";
import { rateMovie } from "../utils/rateMovie";
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

  const axiosClient = useAxiosClient();

  const rate = useMutation((mutationData) =>
    rateMovie(mutationData, axiosClient)
  );

  const handleRate = (action) => {
    if (movies.length === 0) {
      console.log("patience is a virtue");
    } else {
      const mutationData = {
        movieId: movies[0].movie?.movieId,
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
      const resp = (await axiosClient).post(`/watchlist/${movieId}`);
      const {
        data: { statusCode, message },
      } = await resp;
      if (!(statusCode === 201) && !message) {
        console.log("the watchlist did not get filled 😏");
      } else {
        console.log(message);
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
  const handleSkip = () => {
    if (movies.length) {
      setSkip(skip + 1);
      nextMovie();
      openRateFeedbackExported("skip");
    }
  };
  const nextMovie = () => {
    movies.shift();
  };
  const handleWatchLater = () => {
    const mutationData = {
      movieId: movies[0]?.movie?.movieId,
    };
    addToWatchlist.mutate(mutationData);
  };

  useEffect(() => {
    if (movies.length === 0) {
      refetch();
    }
  }, [movies.length, refetch]);

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
              onClick={handleWatchLater}
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
              onClick={handleSkip}
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
          key={compName + m.movie?.movieId}
          onSwipeLeft={() => handleRate("hate")}
          onSwipeRight={() => handleRate("love")}
          swipeSensitivity={50}
          style={{
            backgroundImage: `url(https://thumb.cp.dev.cloudapp.al/thumbnail_${m.movie?.movieId}.jpg)`,
            backgroundSize: "contain",
            backgroundRepeat: "no-repeat",
            backgroundPosition: "center",
          }}
        >
          <MovieTitle
            title={m.movie?.title}
            releaseYear={m.movie?.releaseYear}
          />
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
