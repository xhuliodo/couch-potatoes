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
import { Favorite, SkipNext, ThumbDown, WatchLater } from "@material-ui/icons";
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
          style={{ position: "absolute", bottom: "45px" }}
          container
          justify="center"
          className="cards_container"
        >
          {startedFromTheBottomNowWeHere ? (
            <Tooltip placement="top" arrow title="Add to watchlist">
              <Button
                style={buttonStyling}
                onClick={() => {
                  const mutationData = {
                    movieId: movies[0].movieId,
                  };
                  addToWatchlist.mutate(mutationData);
                }}
                variant="contained"
              >
                <WatchLater fontSize="inherit" />
              </Button>
            </Tooltip>
          ) : (
            <Tooltip placement="top" arrow title="Skip">
              <Button
                style={buttonStyling}
                onClick={() => {
                  nextMovie();
                  setSkip(skip + 1);
                }}
                variant="contained"
              >
                <SkipNext fontSize="inherit" />
              </Button>
            </Tooltip>
          )}
          <Tooltip placement="top" arrow title="Hated it">
            <Button
              style={buttonStyling}
              onClick={() => handleRate("hate")}
              variant="contained"
              color="secondary"
            >
              <ThumbDown fontSize="inherit" />
            </Button>
          </Tooltip>
          <Tooltip placement="top" arrow title="Loved it">
            <Button
              style={buttonStyling}
              onClick={() => handleRate("love")}
              variant="contained"
              color="primary"
            >
              <Favorite fontSize="inherit" />
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
            backgroundImage: `url(${m.posterUrl})`,
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

const buttonStyling = {
  marginLeft: "1.5vw",
  marginRight: "1.5vw",
  fontSize: "35px",
  maxWidth: "100px",
  width: "30vw",
};
