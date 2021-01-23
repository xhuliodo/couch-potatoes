import { Paper } from "@material-ui/core";
import { Card, CardWrapper } from "react-swipeable-cards";

import { useMutation } from "react-query";
import request, { gql } from "graphql-request";

const useLikeMovie = (movieId, userId = "1", nextMovie) => {
  return useMutation(["like-movie", movieId, userId], async () => {
    const data = await request(
      "http://localhost:4001/graphql",
      gql`
        mutation {
          likeMovie(movieId: "${movieId}", userId: "${userId}") {
            movieId
          }
        }
      `
    );
    const { likeMovie } = data;
    if (likeMovie.length === 1) {
      console.log("the movie was rated succesfully", movieId);
      setTimeout(() => {
        nextMovie();
      }, 500);
    } else {
      // TODO: you should implement some error logic, when the rating did not happen
      console.log("the rating didn't do shit");
    }
  });
};

export default function SecondMovieCard({ movies, nextMovie }) {
  const handleLikeMovie = ()=>{
    useLikeMovie(m.movieId, )

  }
  return (
    <CardWrapper style={{ paddingTop: "0px" }}>
      {movies.map((m) => (
        <Card
          key={m.movieId}
          onSwipeLeft={() => {
            console.log("you hated the movie: ", m.title);
          }}
          onSwipeRight={() => {
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
