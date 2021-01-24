import { Paper } from "@material-ui/core";
import { Card, CardWrapper } from "react-swipeable-cards";

import { useMutation } from "react-query";
import request, { gql } from "graphql-request";

export default function SecondMovieCard({ movies, nextMovie }) {
  const mutation = useMutation(async ({ movieId, userId }) => {
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
    if (likeMovie === null) {
      // TODO: you should implement some error logic, when the rating did not happen
      console.log("the rating didn't do shit");
    } else {
      console.log("the movie was rated succesfully", movieId);
      setTimeout(() => {
        nextMovie();
      }, 500);
    }
  });
  return (
    <CardWrapper style={{ paddingTop: "0px" }}>
      {movies.map((m) => (
        <Card
          key={m.movieId}
          onSwipeLeft={() => {
            console.log("you hated the movie: ", m.title);
          }}
          onSwipeRight={() => {
            const mutationData = { movieId: m.movieId, userId: 1 };

            mutation.mutate(mutationData);
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
