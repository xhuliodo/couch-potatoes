import { Backdrop, makeStyles, Typography } from "@material-ui/core";
import { useState } from "react";

import gif from "../utils/setup_helper.gif";

export default function SetupSwipeHelper() {
  const classes = useStyles();

  const [open, setOpen] = useState(true);
  const handleClose = () => {
    setOpen(!open);
  };
  return (
    <Backdrop onClick={handleClose} open={open} className={classes.backdrop}>
      <Typography className={classes.textMargin} align="center">
        Movie cards can be swiped, it's like tinder with no ghosting
        <br />
        <i>(click anywhere to dismiss)</i>
      </Typography>
      <img className={classes.width} alt="setup helper" src={gif}></img>
    </Backdrop>
  );
}

const useStyles = makeStyles((theme) => ({
  backdrop: {
    zIndex: theme.zIndex.drawer + 1,
    color: "#fff",
    flexDirection: "column",
    display: "flex",
    backgroundColor: "rgb(0, 0, 0, 0.8)",
  },
  textMargin: {
    "@media (max-width:850px)": {
      marginBottom: "-13vw",
    },
    "@media (min-width:1024px)": {
      marginBottom: "-9vh",
    },
    marginBottom: "-12vh",
  },
  width: {
    width: "90vw",
    maxWidth: "600px",
  },
}));
