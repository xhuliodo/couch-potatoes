import { useMemo, useState } from "react";

import {
  Button,
  Container,
  Grid,
  makeStyles,
  Paper,
  Typography,
} from "@material-ui/core";
import { DoneOutline } from "@material-ui/icons";

import { withAuthenticationRequired } from "@auth0/auth0-react";

import { useMutation, useQuery } from "react-query";
import { gql } from "graphql-request";
import { useGraphqlClient } from "../utils/useGraphqlClient";

import SelectingGenre from "../components/SelectingGenre";
import AuthLoading from "../components/AuthLoading";

export const SelectingGenrePage = (props) => {
  const classes = useStyle();
  const useGenres = () => {
    return useQuery("genres", async () => {
      const data = (await graphqlClient).request(
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
  const graphqlClient = useGraphqlClient();

  const [selectedGenres, setSelectedGenres] = useState([]);

  const handleSubmit = useMutation(async ({ genres }) => {
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
  });

  const doneIcon = useMemo(
    () => () => <DoneOutline style={{ height: "0.85em" }} />,
    []
  );

  return (
    <Paper elevation={0}>
      <Typography className={classes.instruction}>
        Select at least 3 genres:
      </Typography>
      <Grid container justify="center" style={{ marginTop: "10px" }}>
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
              genres: selectedGenres,
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

const useStyle = makeStyles(() => ({
  instruction: {
    fontSize: "1.2rem",
  },
}));
