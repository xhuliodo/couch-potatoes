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
      <div className={classes.popup}>
        <Typography className={classes.textMargin}>
          Movie cards can be swiped, it's like tinder with no ghosting
        </Typography>
        <i className={classes.textMargin}>(click anywhere to dismiss)</i>
        <div style={{ display: "flex", justifyContent: "center" }}>
          <img className={classes.width} alt="setup helper" src={gif}></img>
        </div>
      </div>
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
    textAlign: "center",
  },
  width: {
    width: "90vw",
    maxWidth: "600px",
  },
  popup: {
    flexDirection: "column",
    display: "flex",
    justifyItems: "center",
    // backgroundColor: "rgb(0, 0, 0, 0.5)",
    borderRadius: "15px",
    width: "95vw",
    padding: "10px",
  },
}));
