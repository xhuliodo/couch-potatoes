import { useAuth0 } from "@auth0/auth0-react";
import { Button, makeStyles, Typography } from "@material-ui/core";
import { Movie } from "@material-ui/icons";
import { useEffect, useState } from "react";
import { useHistory, useLocation } from "react-router-dom";

import cinema from "../utils/cinema.svg";
import "./WelcomePage.scss";

export default function WelcomePage() {
  const classes = useStyle();

  const history = useHistory();
  const location = useLocation();

  const [cinemaBackground, setCinemaBackground] = useState(`url(${cinema})`);
  useEffect(() => {
    if (location.pathname === "/") {
      setCinemaBackground(`url${cinema}`);
    } else {
      setCinemaBackground("");
    }
  }, [location]);

  const { isAuthenticated, loginWithRedirect } = useAuth0();

  const redirect = () => {
    history.push("/solo");
  };

  const handleSignUp = () => {
    if (isAuthenticated) redirect();
    else loginWithRedirect({ screen_hint: "signup" });
  };

  return (
    <div
      style={{
        backgroundImage: cinemaBackground,
      }}
      className={classes.bigDiv}
    >
      <div className="movie_screen">
        <Typography className={classes.mainText}>
          Get personalized
          <Movie className={classes.movieIcon} />
          recommendations
        </Typography>
        <Typography className={classes.pun}>
          insert movie/stalker/lazy pun here
        </Typography>
        <Button
          className={classes.actionButton}
          variant="contained"
          onClick={handleSignUp}
          size="large"
        >
          {isAuthenticated ? "Go back to app" : "Start now"}
        </Button>
      </div>
    </div>
  );
}

const useStyle = makeStyles(() => ({
  actionButton: {
    margin: "3vh auto",
    backgroundColor: "#5c4f74",
    color: "#fff",
    "&:hover": {
      backgroundColor: "#5c4f74",
      boxShadow: "none",
    },
  },
  movieIcon: {
    margin: "5px",
    marginBottom: "-13px",
    fontSize: "3rem",
  },
  mainText: {
    color: "#000",
    textAlign: "center",
    fontSize: "2rem",
  },
  pun: {
    color: "#000",
    textAlign: "center",
    fontSize: "1rem",
  },
  bigDiv: {
    display: "flex",
    flexDirection: "column",
    zIndex: "1",
    height: "100vh",
    alignItems: "center",
    justifyContent: "flex-start",
    backgroundRepeat: "no-repeat",
    backgroundSize: "cover",
    "-webkit-background-size": "cover",
    backgroundPosition: "center center",
  },
}));
