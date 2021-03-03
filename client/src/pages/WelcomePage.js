import { useAuth0 } from "@auth0/auth0-react";
import { Button, Snackbar, Typography } from "@material-ui/core";
import { ConfirmationNumber, Movie } from "@material-ui/icons";
import { useEffect, useState } from "react";
import { useHistory, useLocation } from "react-router-dom";

import cinema from "../utils/cinema.png";
import background from "../utils/background.svg";
import "./WelcomePage.scss";

export default function WelcomePage() {
  const history = useHistory();
  const location = useLocation();

  const [cinemaBackground, setCinemaBackground] = useState(`url(${cinema})`);
  useEffect(() => {
    if (location.pathname === "/") {
      setCinemaBackground(`url${background}`);
    } else {
      setCinemaBackground("");
    }
  }, [location]);
  const redirect = () => {
    history.push("/solo");
  };

  const action = (
    <Button variant="contained" onClick={redirect} size="small">
      GO
    </Button>
  );

  const { isAuthenticated, loginWithRedirect } = useAuth0();

  const handleSignUp = () => {
    loginWithRedirect({ screen_hint: "signup" });
  };

  const [open, setOpen] = useState(isAuthenticated);

  useEffect(() => {
    setOpen(isAuthenticated);
  }, [isAuthenticated]);

  const handleClose = (event, reason) => {
    if (reason === "clickaway") {
      return;
    }

    setOpen(false);
  };

  return (
    <div
      style={{
        display: "flex",
        flexDirection: "column",
        // flex: "2 1 auto",
        zIndex: "1",
        height: "100%",
        alignItems: "center",
        justifyContent: "flex-start",
        backgroundImage: cinemaBackground,
        
      }}
      className="responsive_background"
    >
      <div className="movie_screen">
        <Typography
          className="main_text"
          style={{ fontSize: "2rem" }}
          align="center"
        >
          Get personalized
          <Movie
            style={{
              margin: "5px",
              marginBottom: "-13px",
              fontSize: "3rem",
            }}
          />
          recommendations
        </Typography>
        <Typography
          className="main_text"
          style={{ fontSize: "1rem" }}
          align="center"
        >
          insert movie/stalker/lazy pun here
        </Typography>
        <Button
          style={{
            margin: "5vh auto",
            backgroundColor: "#8e72a7",
            color: "#fff",
          }}
          variant="contained"
          onClick={handleSignUp}
          size="large"
        >
          Start now
          {/* Start couching now */}
          {/* Get a sure{" "}
        <span style={{ fontSize: "2rem", marginLeft: "5px" }}>🎟</span> */}
          {/* Get a sure{" "}
        <ConfirmationNumber
          style={{ marginLeft: "5px", color: "#b04838" }}
          as="span"
        /> */}
        </Button>
      </div>

      <Snackbar
        open={open}
        onClose={handleClose}
        action={action}
        message="Go back to using the app?"
        anchorOrigin={{ horizontal: "center", vertical: "bottom" }}
      />
    </div>
  );
}
