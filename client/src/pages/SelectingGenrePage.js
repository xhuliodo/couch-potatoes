import { Button, Container, Grid, Paper, Typography } from "@material-ui/core";

import { useMutation, useQuery } from "react-query";
import { gql } from "graphql-request";

import SelectingGenre from "../components/SelectingGenre";
import { useMemo, useState } from "react";
import { useAuth0, withAuthenticationRequired } from "@auth0/auth0-react";
import AuthLoading from "../components/AuthLoading";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { DoneOutline } from "@material-ui/icons";

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
    // const userData = await (await graphqlClient).request(
    //   gql`
    //     mutation {
    //       selfRegister {
    //         userId
    //       }
    //     }
    //   `
    // );
    // const { selfRegister } = userData;
    // if (selfRegister?.userId) {
    const data = await (await graphqlClient).request(
      gql`
        mutation {
          setFavoriteGenres(
            genres: [${genres.map((g) => `"${g}"`)}]
          ) {
            userId
          }
        }
      `
    );

    const { setFavoriteGenres } = data;
    if (setFavoriteGenres.userId !== null) {
      props.history.push("/getting-to-know-2");
    }
    // }
  });

  const doneIcon = useMemo(() => () => <DoneOutline />, []);

  return (
    <Paper elevation={0}>
      <Typography variant="h5">Select at least 3 genres:</Typography>
      <Grid container justify="center" style={{ marginTop: "15px" }}>
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
              doneIcon={doneIcon}
            />
          ))
        )}
      </Grid>
      <Container disableGutters={true} style={{ marginTop: "15px" }}>
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
