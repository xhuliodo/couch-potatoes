import { makeStyles, Snackbar, useTheme } from "@material-ui/core";
import { Favorite, ThumbDown, WatchLater } from "@material-ui/icons";

import React, { useState } from "react";

import "./RateFeedback.css";

let openRateFeedbackFn;

export default function RateFeedback() {
  const classes = useStyles();
  const theme = useTheme();

  const [open, setOpen] = useState(false);
  const [action, setAction] = useState("");

  const openRateFeedback = (action) => {
    setAction(action);
    setOpen(true);
  };

  openRateFeedbackFn = openRateFeedback;

  return (
    <Snackbar
      anchorOrigin={{ horizontal: "center", vertical: "top" }}
      open={open}
      className={classes.feedback}
      autoHideDuration={400}
      onClose={() => {
        setOpen(false);
        setAction("");
      }}
      style={{ position: "absolute" }}
    >
      <div
        style={{
          borderRadius: "15px",
          backgroundColor: theme.palette.background.paper,
          opacity: "0.85",
        }}
      >
        {action === "love" ? (
          <Favorite style={iconStyling} />
        ) : action === "hate" ? (
          <ThumbDown style={iconStyling} />
        ) : action === "watchlater" ? (
          <WatchLater style={iconStyling} />
        ) : null}
      </div>
    </Snackbar>
  );
}

export function openRateFeedbackExported(action) {
  openRateFeedbackFn(action);
}

const useStyles = makeStyles((theme) => ({
  feedback: {
    // display: "flex",
    // justifyContent: "center",
    // alignSelf: "center",
    top: "auto",
    bottom: "50%!important",
    left: "50%",
    // right: "auto!important",
  },
}));

const iconStyling = {
  fontSize: "15vh",
  padding: "15px",
};
