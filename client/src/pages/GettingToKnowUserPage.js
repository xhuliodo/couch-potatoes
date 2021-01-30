import { useEffect, useState } from "react";

import {
  Button,
  Container,
  Paper,
  Typography,
  Dialog,
  DialogTitle,
  DialogContent,
  DialogContentText,
  DialogActions,
} from "@material-ui/core";
import "../components/MovieCard.css";

import { useMovieStore } from "../context/movies";

import GenreBasedRec from "../components/GenreBasedRec";

export default function GettingToKnowUserPage(props) {
  const [skip, setSkip] = useState(0);

  const { ratedMovies, requiredMovies } = useMovieStore();

  const [open, setOpen] = useState(false);
  const handleClose = () => {
    setOpen(false);
  };
  useEffect(() => {
    // eslint-disable-next-line no-unused-expressions
    ratedMovies === requiredMovies ? setOpen(true) : null;
  }, [ratedMovies, requiredMovies]);
  const handleNext = () => {
    handleClose();
    props.history.push("/recommendations");
  };

  return (
    <Paper elevation={0} style={{ height: "90vh" }}>
      <Typography style={{ textAlign: "center" }} variant="h6">
        Please rate at least {requiredMovies} movies{" "}
        <i>(ignored movies will not count)</i>
      </Typography>
      <Typography
        style={{ textAlign: "center", fontWeight: "bold" }}
        variant="h6"
      >
        You have rated {ratedMovies} / {requiredMovies}
      </Typography>
      <Container disableGutters={true}>
        <GenreBasedRec skip={skip} setSkip={setSkip} />
        <Container maxWidth="sm">
          <Button
            style={{ marginTop: "3.4vh" }}
            color="primary"
            fullWidth
            variant="contained"
            onClick={handleNext}
            disabled={ratedMovies >= requiredMovies ? false : true}
          >
            Next
          </Button>
        </Container>
        <Dialog open={open} onClose={handleClose}>
          <DialogTitle>We've gotten to know you enough 😃</DialogTitle>
          <DialogContent>
            <DialogContentText>
              You can either click next to get recommendations based on other
              users ratings or continue to rate movies in the genre you selected
              as well (you can always click the continue button at the end to go
              to the next step)
            </DialogContentText>
            <DialogActions>
              <Button onClick={handleClose} color="secondary">
                Continue
              </Button>
              <Button onClick={handleNext} color="primary">
                Next
              </Button>
            </DialogActions>
          </DialogContent>
        </Dialog>
      </Container>
    </Paper>
  );
}
