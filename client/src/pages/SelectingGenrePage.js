import { Button, Container, Grid, Paper, Typography } from "@material-ui/core";

import { useMutation, useQuery } from "react-query";
import { request, gql } from "graphql-request";

import SelectingGenre from "../components/SelectingGenre";
import { useState } from "react";
import { useAuth0, withAuthenticationRequired } from "@auth0/auth0-react";
import AuthLoading from "../components/AuthLoading";
import { useGetToken } from "../utils/getToken";

const useGenres = () => {
  return useQuery("genres", async () => {
    const data = await request(
      "http://localhost:4001/graphql",
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

export const SelectingGenrePage = (props) => {
  const { status, data, error } = useGenres();

  const [selectedGenres, setSelectedGenres] = useState([]);

  const { user } = useAuth0();
  const token = useGetToken();
  console.log(token);

  const handleSubmit = useMutation(async ({ userId, genres }) => {
    const data = await request(
      "http://localhost:4001/graphql",
      { headers: { authorization: `Bearer ${token}` } },
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
    console.log(MergeUserFavoriteGenres);
    if (MergeUserFavoriteGenres.userId !== null) {
      props.history.push("/getting-to-know-2");
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
            handleSubmit.mutate({ userId: user.sub, genres: selectedGenres })
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
