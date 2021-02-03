import { CircularProgress, makeStyles } from "@material-ui/core";

import React from "react";

export default function AuthLoading() {
  const classes = useStyles();
  return <CircularProgress className={classes.loading} />;
}

const useStyles = makeStyles((theme) => ({
  loading: {
    position: "absolute",
    display: "flex",
    justifyContent: "center",
    height: "20vh!important",
    width: "20vw!important",
    backgroundColor: "white",
    top: "40vh",
    bottom: 0,
    left: "40vw",
    right: 0,
  },
}));
