import {
  Backdrop,
  CircularProgress,
  makeStyles,
  Typography,
} from "@material-ui/core";

import React from "react";

export default function AuthLoading() {
  const classes = useStyles();
  return (
    <Backdrop open className={classes.backdrop}>
      <CircularProgress className={classes.loading} />
      <br />
      <br />
      <Typography variant="h5">Logging you in, please wait</Typography>
    </Backdrop>
  );
}

const useStyles = makeStyles((theme) => ({
  loading: {
    color: theme.palette.primary.main,
    height: "15vh!important",
    width: "15vh!important",
  },
  backdrop: {
    zIndex: theme.zIndex.drawer + 1,
    color: "#fff",
    flexDirection: "column",
  },
}));
