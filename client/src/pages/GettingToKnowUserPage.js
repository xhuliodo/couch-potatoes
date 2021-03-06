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
  makeStyles,
} from "@material-ui/core";
// auth redirect import
import { withAuthenticationRequired } from "@auth0/auth0-react";
import AuthLoading from "../components/AuthLoading";
// global state import
import { useMovieStore } from "../context/movies";

import GenreBasedRec from "../components/GenreBasedRec";
import "../components/MovieCard.scss";
import SetupSwipeHelper from "../components/SetupSwipeHelper";

export const GettingToKnowUserPage = (props) => {
  const classes = useStyle();

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
        <Typography align="center">
          To give you tailored recommendations, we need to know your taste.
        </Typography>
        <Typography className={classes.boldInfoText} align="center">
          {requiredMovies - ratedMovies > 0
            ? `${requiredMovies - ratedMovies} ratings left`
            : "no more rating needed"}
        </Typography>
        <Container style={{ marginTop: "10px" }} disableGutters={true}>
          <GenreBasedRec skip={skip} setSkip={setSkip} />
          <Container maxWidth="xs">
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
              mission accomplished ;)
            </DialogTitle>
            <DialogContent>
              <DialogContentText style={{ textAlign: "center" }}>
                You can keep rating popular movies or get some tailored
                recommendations
              </DialogContentText>
              <DialogActions style={{ justifyContent: "center" }}>
                <Button
                  variant="contained"
                  onClick={handleClose}
                  color="secondary"
                >
                  Keep going
                </Button>
                <Button
                  variant="contained"
                  onClick={handleNext}
                  color="primary"
                >
                  Give me
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

const useStyle = makeStyles(() => ({
  boldInfoText: {
    fontWeight: "bolder",
    fontSize: "1.2rem",
  },
}));

export default withAuthenticationRequired(GettingToKnowUserPage, {
  onRedirecting: () => <AuthLoading />,
});
