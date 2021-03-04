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
        <Typography className={classes.infoText} align="center">
          Rate at least {requiredMovies} movies
          <br />
        </Typography>
        <Typography align="center">
          <i>(ignored movies will not count)</i>
        </Typography>

        <Typography className={classes.boldInfoText} align="center">
          {requiredMovies - ratedMovies > 0
            ? `${requiredMovies - ratedMovies} ratings left`
            : "mission accomplished ;)"}
        </Typography>

        <Container style={{ marginTop: "10px" }} disableGutters={true}>
          <GenreBasedRec skip={skip} setSkip={setSkip} />
          <Container maxWidth="xs">
            <Button
              style={{ marginTop: "15px" }}
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
                You can continue rating popular movies or get some tailored
                recommendations
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

const useStyle = makeStyles(() => ({
  infoText: {
    fontSize: "1.2rem",
  },
  boldInfoText: {
    fontWeight: "bold",
    fontSize: "1.2rem",
  },
}));

export default withAuthenticationRequired(GettingToKnowUserPage, {
  onRedirecting: () => <AuthLoading />,
});
