import { Paper } from "@material-ui/core";
import { Card, CardWrapper } from "react-swipeable-cards";

import { useMutation } from "react-query";
import request, { gql } from "graphql-request";
import { useEffect } from "react";

export default function SecondMovieCard({ movies, nextMovie, increaseSkip }) {
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
      setTimeout(() => {
        nextMovie();
      }, 500);
    }
  });

  useEffect(() => {
    console.log(movies.length);
    if (movies.length < 1) {
      increaseSkip();
    }
  }, [movies, rate, increaseSkip]);
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
              userId: 1,
              action: "hate",
            };
            rate.mutate(mutationData);
          }}
          onSwipeRight={() => {
            const mutationData = {
              movieId: m.movieId,
              userId: 1,
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
    </CardWrapper>
  );
}
