import { Button, Container, Grid, Paper, Typography } from "@material-ui/core";

import { useQuery } from "react-query";
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
            _id
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

  const handleSubmit = () => {
    setGenres(selectedGenres);
    props.history.push("/getting-to-know-2");
  };

  return (
    <Paper elevation={0}>
      <Typography variant="h5">Select at least 3 genres:</Typography>
      {/* <Container disableGutters={true} style={{ marginTop: "2.5vh" }}> */}
      <Grid container justify="center" style={{ marginTop: "2.5vh" }}>
        {status === "loading" ? (
          <span>Fetching data</span>
        ) : status === "error" ? (
          <span>Error: {error.message}</span>
        ) : (
          data.map((g) => (
            <SelectingGenre
              key={g._id}
              selectedGenres={selectedGenres}
              setSelectedGenres={setSelectedGenres}
              genre={g}
            />
          ))
        )}
      </Grid>
      {/* </Container> */}
      <Container disableGutters={true} style={{ marginTop: "2.5vh" }}>
        <Button
          disabled={selectedGenres.length < 3}
          fullWidth={true}
          size="large"
          color="primary"
          variant="contained"
          onClick={handleSubmit}
        >
          Next
        </Button>
      </Container>
    </Paper>
  );
}
