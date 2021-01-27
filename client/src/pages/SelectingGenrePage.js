import { Button, Container, Grid, Paper, Typography } from "@material-ui/core";

import { useMutation, useQuery } from "react-query";
import { request, gql } from "graphql-request";

import SelectingGenre from "../components/SelectingGenre";
import { useState } from "react";

import { useGenreStore } from "../context/genres";

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

export default function SelectingGenrePage(props) {
  const { status, data, error } = useGenres();

  const [selectedGenres, setSelectedGenres] = useState([]);

  const setGenres = useGenreStore((state) => state.setGenres);

  const handleSubmit = useMutation(async ({ userId, genres }) => {
    console.log(genres);
    const data = await request(
      "http://localhost:4001/graphql",
      gql`
        mutation  {
          MergeUserFavoriteGenres(
            from: { userId: "${userId}" }
            to: {
              genreId_in: ["${...genres}"]
            }
          ) {
            from {
              userId
            }
          }
        }
      `
    );
    console.log(data);
    setGenres(selectedGenres);
    props.history.push("/getting-to-know-2");
  });

  // setGenres(selectedGenres);

  // props.history.push("/getting-to-know-2");

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
            handleSubmit.mutate({ userId: "1", genres: selectedGenres })
          }
        >
          Next
        </Button>
      </Container>
    </Paper>
  );
}
