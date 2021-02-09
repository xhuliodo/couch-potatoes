import { Button, Container, Grid, Paper, Typography } from "@material-ui/core";

import { useMutation, useQuery } from "react-query";
import { gql } from "graphql-request";

import SelectingGenre from "../components/SelectingGenre";
import { useState } from "react";
import { useAuth0, withAuthenticationRequired } from "@auth0/auth0-react";
import AuthLoading from "../components/AuthLoading";
import { useGraphqlClient } from "../utils/useGraphqlClient";

export const SelectingGenrePage = (props) => {
  const useGenres = () => {
    return useQuery("genres", async () => {
      const data = await (await graphqlClient).request(
        gql`
          query {
            Genre {
              genreId
              name
            }
          }
        `
      );
      const { Genre } = await data;
      return Genre;
    });
  };
  const { status, data, error } = useGenres();
  const [selectedGenres, setSelectedGenres] = useState([]);

  const { user } = useAuth0();
  const graphqlClient = useGraphqlClient();

  const handleSubmit = useMutation(async ({ userId, genres, name }) => {
    const userData = await (await graphqlClient).request(
      gql`mutation {
        MergeUser(
          where: { userId: "${userId}" }
          data: { userId: "${userId}" }
        ){
          userId
        }
      }`
    );
    const { MergeUser } = userData;
    if (MergeUser?.userId) {
      const data = await (await graphqlClient).request(
        gql`
        mutation  {
          MergeUserFavoriteGenres(
            from: { userId: "${userId}" }
            to: {
              genreId_in: [${genres.map((g) => `"${g}"`)}]
            }
          ) {
            from {
              userId
            }
          }
        }
      `
      );

      const { MergeUserFavoriteGenres } = data;
      if (MergeUserFavoriteGenres.userId !== null) {
        props.history.push("/getting-to-know-2");
      }
    }
  });

  return (
    <Paper elevation={0}>
      <Typography variant="h5">Select at least 3 genres:</Typography>
      <Grid container justify="center" style={{ marginTop: "2.5vh" }}>
        {status === "loading" ? (
          <span>Fetching data</span>
        ) : status === "error" ? (
          <span>Error: {error.message}</span>
        ) : (
          data.map((g) => (
            <SelectingGenre
              key={g.genreId}
              selectedGenres={selectedGenres}
              setSelectedGenres={setSelectedGenres}
              genre={g}
            />
          ))
        )}
      </Grid>
      <Container disableGutters={true} style={{ marginTop: "2.5vh" }}>
        <Button
          disabled={selectedGenres.length < 3}
          fullWidth={true}
          size="large"
          color="primary"
          variant="contained"
          onClick={() =>
            handleSubmit.mutate({
              userId: user?.sub,
              genres: selectedGenres,
              name: user?.given_name,
            })
          }
        >
          Next
        </Button>
      </Container>
    </Paper>
  );
};

export default withAuthenticationRequired(SelectingGenrePage, {
  onRedirecting: () => <AuthLoading />,
});
