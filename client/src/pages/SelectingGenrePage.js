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
import { useAxiosClient } from "../utils/useAxiosClient";

import SelectingGenre from "../components/SelectingGenre";
import AuthLoading from "../components/AuthLoading";

export const SelectingGenrePage = (props) => {
  const classes = useStyle();

  const axiosClient = useAxiosClient();
  const { status, data, error } = useGenres({ axiosClient });

  const [selectedGenres, setSelectedGenres] = useState([]);

  const handleSubmit = useMutation(async ({ genres }) => {
    const resp = (await axiosClient).post("/users/genres", {
      // genres:[`${genres.map((g) => `"${g}"`)}`],
      genres,
    });
    const {
      data: { statusCode, message },
    } = await resp;
    if (!statusCode && !message) {
      console.log("the watchlist did not get filled 😏");
    }
    if (statusCode === 201) {
      props.history.push("/getting-to-know-2");
    } else {
      console.log(message);
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

const useGenres = ({ axiosClient }) => {
  return useQuery("genres", async () => {
    const resp = (await axiosClient).get(`/genres`);
    const {
      data: { data: genres },
    } = await resp;
    return genres;
  });
};
