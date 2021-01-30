import { Button, Grid, Paper, Tooltip } from "@material-ui/core";
import { Card, CardWrapper } from "react-swipeable-cards";

import { useMutation } from "react-query";
import request, { gql } from "graphql-request";
import { useEffect } from "react";
import { ThumbDown, ThumbUp, VisibilityOff } from "@material-ui/icons";
import { useMovieStore } from "../context/movies";

export default function MovieCard({
  movies,
  skip,
  setSkip,
  refetch,
  startedFromTheBottomNowWeHere = false,
}) {
  const rate = useMutation(async ({ movieId, userId, action }) => {
    let data = { rateMovie: null };
    // eslint-disable-next-line default-case
    switch (action) {
      case "love":
        data = await request(
          "http://localhost:4001/graphql",
          gql`
            mutation {
              rateMovie(movieId: "${movieId}", userId: "${userId}", rating:1) {
                movieId
              }
            }
          `
        );
        break;
      case "hate":
        data = await request(
          "http://localhost:4001/graphql",
          gql`
            mutation {
              rateMovie(movieId: "${movieId}", userId: "${userId}", rating:0) {
                movieId
              }
            }
          `
        );
        break;
    }
    const { rateMovie } = await data;
    if (rateMovie === null) {
      // TODO: you should implement some error logic, when the rating did not happen
      console.log("the rating didn't do shit");
    } else {
      console.log("you rated the movie with id", rateMovie.movieId);
      nextMovie();
      increaseRatedMovies();
    }
  });

  const nextMovie = () => {
    movies.shift();
  };

  const { increaseRatedMovies } = useMovieStore();

  useEffect(() => {
    if (movies.length === 0) {
      console.log(("movies length is: ", movies.length));
      refetch();
    }
  }, [movies, rate, refetch]);

  // TODO: implement this card when the ratings are done, for the user to be forwarded to the normal page
  // const waitForMoreData = () => {
  //   let titleStyle = {
  //     textAlign: "center",
  //     fontWeight: "bold",
  //     fontSize: "40px",
  //     fontFamily: "Sans-Serif",
  //     marginTop: "50px",
  //   };
  //   return <div style={titleStyle}>fetching more movies...</div>;
  // };

  const ActionButtons = (props) => {
    const newProps = {};

    return (
      <Grid
        {...newProps}
        style={{ position: "absolute", bottom: "3px" }}
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

        <Tooltip placement="top" arrow title="Loved it!">
          <Button
            style={buttonStyling}
            onClick={() => {
              const mutationData = {
                movieId: movies[0].movieId,
                userId: 2,
                action: "love",
              };
              rate.mutate(mutationData);
            }}
            variant="contained"
            color="primary"
          >
            <ThumbUp fontSize="inherit" />
          </Button>
        </Tooltip>
        <Tooltip placement="top" arrow title="Hated it!">
          <Button
            style={buttonStyling}
            onClick={() => {
              const mutationData = {
                movieId: movies[0].movieId,
                userId: 2,
                action: "hate",
              };
              rate.mutate(mutationData);
            }}
            variant="contained"
            color="secondary"
          >
            <ThumbDown fontSize="inherit" />
          </Button>
        </Tooltip>
      </Grid>
    );
  };

  return (
    <CardWrapper
      // addEndCard={waitForMoreData.bind(this)}
      style={{ paddingTop: "0px" }}
    >
      {movies.map((m) => (
        <Card
          key={m.movieId}
          onSwipeLeft={() => {
            const mutationData = {
              movieId: m.movieId,
              userId: 2,
              action: "hate",
            };
            rate.mutate(mutationData);
          }}
          onSwipeRight={() => {
            const mutationData = {
              movieId: m.movieId,
              userId: 2,
              action: "love",
            };
            rate.mutate(mutationData);
          }}
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
  );
}

const buttonStyling = {
  marginLeft: "1.5vw",
  marginRight: "1.5vw",
  fontSize: "40px",
  maxWidth: "100px",
  width: "30vw",
};
