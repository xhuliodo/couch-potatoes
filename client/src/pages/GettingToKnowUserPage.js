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

import { withAuthenticationRequired } from "@auth0/auth0-react";
import AuthLoading from "../components/AuthLoading";

import { useMovieStore } from "../context/movies";

import GenreBasedRec from "../components/GenreBasedRec";
import "../components/MovieCard.scss";
import SetupSwipeHelper from "../components/SetupSwipeHelper";

export const GettingToKnowUserPage = (props) => {
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
    props.history.push("/solo");
  };

  return (
    <>
      <Paper elevation={0}>
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
        <Container style={{ marginTop: "15px" }} disableGutters={true}>
          <GenreBasedRec skip={skip} setSkip={setSkip} />
          <Container style={{ marginTop: "15px" }} maxWidth="sm">
            <Button
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
            <DialogTitle style={{ textAlign: "center" }}>
              We've gotten to know you enough 😃
            </DialogTitle>
            <DialogContent>
              <DialogContentText style={{ textAlign: "center" }}>
                You can either click next to get recommendations based on other
                users ratings or continue to rate movies in the genre you
                selected as well (you can always click the continue button at
                the end to go to the next step)
              </DialogContentText>
              <DialogActions style={{ justifyContent: "center" }}>
                <Button
                  variant="contained"
                  onClick={handleClose}
                  color="secondary"
                >
                  Continue
                </Button>
                <Button
                  variant="contained"
                  onClick={handleNext}
                  color="primary"
                >
                  Next
                </Button>
              </DialogActions>
            </DialogContent>
          </Dialog>
        </Container>
      </Paper>
      <SetupSwipeHelper />
    </>
  );
};

export default withAuthenticationRequired(GettingToKnowUserPage, {
  onRedirecting: () => <AuthLoading />,
});
