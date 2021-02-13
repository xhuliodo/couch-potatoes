import { Button, Container, Grid, Paper, Tooltip } from "@material-ui/core";
import { Card, CardWrapper } from "@xhuliodo/react-swipeable-cards";

import { useMutation } from "react-query";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";
import { useEffect } from "react";
import { ThumbDown, ThumbUp, VisibilityOff } from "@material-ui/icons";
import { useMovieStore } from "../context/movies";
import { rateMovie } from "../utils/rateMovie";

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
    const mutationData = {
      movieId: movies[0].movieId,
      action,
      successFunc,
    };
    rate.mutate(mutationData);
  };

  const successFunc = () => {
    nextMovie();
    increaseRatedMovies();
  };

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
          style={{ position: "absolute", bottom: "50px" }}
          container
          justify="center"
          className="cards_container"
        >
          {startedFromTheBottomNowWeHere ? null : (
            <Tooltip placement="top" arrow title="Ignore...">
              <Button
                style={buttonStyling}
                onClick={() => {
                  nextMovie();
                  setSkip(skip + 1);
                }}
                variant="contained"
              >
                <VisibilityOff fontSize="inherit" />
              </Button>
            </Tooltip>
          )}
          <Tooltip placement="top" arrow title="Hated it!">
            <Button
              style={buttonStyling}
              onClick={() => handleRate("hate")}
              variant="contained"
              color="secondary"
            >
              <ThumbDown fontSize="inherit" />
            </Button>
          </Tooltip>
          <Tooltip placement="top" arrow title="Loved it!">
            <Button
              style={buttonStyling}
              onClick={() => handleRate("love")}
              variant="contained"
              color="primary"
            >
              <ThumbUp fontSize="inherit" />
            </Button>
          </Tooltip>
        </Grid>
        {startedFromTheBottomNowWeHere ? (
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
        ) : null}
      </>
    );
  };

  return (
    <>
      <CardWrapper
        // addEndCard={waitForMoreData.bind(this)}
        style={{ paddingTop: "0px" }}
      >
        {movies.map((m) => (
          <Card
            key={m.movieId}
            onSwipeLeft={() => handleRate("hate")}
            onSwipeRight={() => handleRate("love")}
            swipeSensitivity={60}
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
        <ActionButtons />
      </CardWrapper>
    </>
  );
}

const buttonStyling = {
  marginLeft: "1.5vw",
  marginRight: "1.5vw",
  fontSize: "35px",
  maxWidth: "100px",
  width: "30vw",
};
