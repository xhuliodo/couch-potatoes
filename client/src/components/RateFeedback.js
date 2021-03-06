import { useState } from "react";
import { makeStyles, Snackbar } from "@material-ui/core";
import {
  SentimentDissatisfied,
  SentimentVerySatisfiedRounded,
  SkipNextOutlined,
  WatchLaterOutlined,
} from "@material-ui/icons";
import "./RateFeedback.scss";

let openRateFeedbackFn;

export default function RateFeedback() {
  const classes = useStyles();

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
      <div className={classes.insideDiv}>
        {action === "love" ? (
          <SentimentVerySatisfiedRounded style={iconStyling} />
        ) : action === "hate" ? (
          <SentimentDissatisfied style={iconStyling} />
        ) : action === "watchlater" ? (
          <WatchLaterOutlined style={iconStyling} />
        ) : action === "skip" ? (
          <SkipNextOutlined style={iconStyling} />
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
    top: "auto",
    bottom: "50%!important",
    left: "50%",
  },
  insideDiv: {
    borderRadius: "15px",
    backgroundColor: theme.palette.background.paper,
    opacity: "0.85",
  },
}));

const iconStyling = {
  fontSize: "15vh",
  padding: "15px",
};
