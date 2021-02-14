import { CircularProgress, makeStyles } from "@material-ui/core";

import React from "react";

export default function AuthLoading() {
  const classes = useStyles();
  return <CircularProgress className={classes.loading} />;
}

const useStyles = makeStyles((theme) => ({
  loading: {
    color: theme.palette.primary.main,
    position: "absolute",
    justifyContent: "center",
    height: "20vh!important",
    width: "20vw!important",
    top: "40vh",
    bottom: 0,
    left: "40vw",
    right: 0,
  },
}));
