import { useAuth0 } from "@auth0/auth0-react";
import { Button, Icon, Snackbar, Typography } from "@material-ui/core";
import { ConfirmationNumber } from "@material-ui/icons";
import { useEffect, useState } from "react";
import { useHistory } from "react-router-dom";

import "./WelcomePage.scss";

export default function WelcomePage() {
  const history = useHistory();
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
      }}
    >
      <div className="movie_screen">
        <Typography style={{ fontSize: "2rem" }} align="center">
          Get personalized 🎬 recommendations
        </Typography>
        <Typography style={{ fontSize: "1rem" }} align="center">
          insert movie/stalker/lazy pun here
        </Typography>
        <Button
          style={{ margin: "5vh auto" }}
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
