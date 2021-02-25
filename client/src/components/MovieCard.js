import {
  Button,
  Grid,
  LinearProgress,
  Paper,
  Tooltip,
  Typography,
} from "@material-ui/core";
import { Card, CardWrapper } from "@xhuliodo/react-swipeable-cards";

import { useMutation } from "react-query";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";
import { useEffect } from "react";
import {
  LibraryAddOutlined,
  SentimentDissatisfied,
  SentimentVerySatisfiedRounded,
  SkipNext,
} from "@material-ui/icons";
import { useMovieStore } from "../context/movies";
import { rateMovie } from "../utils/rateMovie";
import { openRateFeedbackExported } from "./RateFeedback";
import RateFeedback from "./RateFeedback";

export default function MovieCard({
  movies,
  skip,
  setSkip,
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

  const successFunc = (action) => {
    openRateFeedbackExported(action);
    nextMovie();
    increaseRatedMovies();
  };

  const waitingForMore = () => (
    <div>
      <Typography
        style={{
          color: "black",
          textAlign: "center",
          fontWeight: "bold",
          fontSize: "35px",
          marginTop: "45%",
        }}
      >
        Fetching movies...
      </Typography>
      <LinearProgress style={{ width: "85%", margin: "50px auto" }} />
    </div>
  );

  const feedback = (action) => <RateFeedback action={action} />;

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

  const nextMovie = () => {
    movies.shift();
  };

  useEffect(() => {
    if (movies.length === 0) {
      refetch();
    }
  }, [movies, rate, refetch]);

  const ActionButtons = (props) => {
    const newProps = {};

    return (
      <>
        <Grid
          {...newProps}
          style={{ position: "absolute", bottom: "5px", width: "100vw" }}
          container
          justify="center"
          className="cards_container"
        >
          <Tooltip placement="top" arrow title="Not for me">
            <Button
              className="actionButton"
              onClick={() => handleRate("hate")}
              variant="contained"
              color="secondary"
            >
              <SentimentDissatisfied fontSize="inherit" />
            </Button>
          </Tooltip>
          {startedFromTheBottomNowWeHere ? (
            <Tooltip placement="top" arrow title="Add to watchlist">
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
                <LibraryAddOutlined fontSize="inherit" />
              </Button>
            </Tooltip>
          ) : (
            <Tooltip placement="top" arrow title="Skip">
              <Button
                className="actionButton"
                onClick={() => {
                  setSkip(skip + 1);
                  nextMovie();
                  openRateFeedbackExported("skip");
                }}
                variant="contained"
              >
                <SkipNext fontSize="inherit" />
              </Button>
            </Tooltip>
          )}
          <Tooltip placement="top" arrow title="Loved it">
            <Button
              className="actionButton"
              onClick={() => handleRate("love")}
              variant="contained"
              color="primary"
            >
              <SentimentVerySatisfiedRounded fontSize="inherit" />
            </Button>
          </Tooltip>
        </Grid>
        {/* to return the watchlater strip uncomment this part
        and set null to started from the bottom before the skip button */}
        {/* {startedFromTheBottomNowWeHere ? (
          <Container
            style={{ position: "absolute", bottom: "2px" }}
            maxWidth="sm"
          >
            <Button
              style={{ marginTop: "15px" }}
              color="primary"
              fullWidth
              variant="contained"
              onClick={() => {
                const mutationData = {
                  movieId: movies[0].movieId,
                };
                addToWatchlist.mutate(mutationData);
              }}
            >
              Add to watchlist
            </Button>
          </Container>
        ) : null} */}
      </>
    );
  };

  return (
    <CardWrapper addEndCard={waitingForMore} style={{ paddingTop: "0px" }}>
      {movies.map((m) => (
        <Card
          key={m.movieId}
          onSwipeLeft={() => {
            handleRate("hate");
            feedback("hate");
          }}
          onSwipeRight={() => handleRate("love")}
          swipeSensitivity={100}
          style={{
            backgroundImage: `url(https://thumb.cp.dev.cloudapp.al/thumbnail_${m.movieId}.jpg)`,
            backgroundSize: "contain",
            backgroundRepeat: "no-repeat",
            backgroundPosition: "center",
          }}
        >
          <Paper className="secondMovie_cardText">
            <h3>
              {m.title} <i>({m.releaseYear})</i>
            </h3>
          </Paper>
        </Card>
      ))}
      <RateFeedback />
      <ActionButtons />
    </CardWrapper>
  );
}

