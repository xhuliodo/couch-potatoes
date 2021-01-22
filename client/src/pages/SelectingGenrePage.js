import { Button, Container, Paper, Typography } from "@material-ui/core";

import { useQuery } from "react-query";
import { request, gql } from "graphql-request";
import SelectingGenre from "../components/SelectingGenre";
import { useState } from "react";

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

export default function SelectingGenrePage() {
  const { status, data, error } = useGenres();

  const [selectedGenres, setSelectedGenres] = useState([]);

  return (
    <Paper elevation={0}>
      <Typography variant="h5">Select at least 3 genres:</Typography>
      <Container disableGutters={true} style={{ marginTop: "2.5vh" }}>
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
      </Container>
      <Container disableGutters={true} style={{ marginTop: "2.5vh" }}>
        <Button
          disabled={selectedGenres.length < 3}
          fullWidth={true}
          size="large"
          color="primary"
          variant="contained"
        >
          Next
        </Button>
      </Container>
    </Paper>
  );
}
