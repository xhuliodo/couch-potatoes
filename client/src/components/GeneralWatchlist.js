import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Container,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { ExpandMore } from "@material-ui/icons";
import React, { useState } from "react";
import RatedMoviesInWatchlist from "./RatedMoviesInWatchlist";
import WatchlistProvider from "./WatchlistProvider";

export default function GeneralWatchlist() {
  const classes = useStyles();

  const [open, setOpen] = useState("panel1");

  const handleChange = (panel) => (event, isExpanded) => {
    setOpen(isExpanded ? panel : false);
  };

  return (
    <Container
      maxWidth="lg"
      style={{ margin: "0", paddingLeft: "0", paddingRight: "0" }}
    >
      <Accordion
        TransitionProps={{ unmountOnExit: true }}
        expanded={open === "panel1"}
        onChange={handleChange("panel1")}
      >
        <AccordionSummary
          className={classes.summary}
          expandIcon={<ExpandMore />}
          aria-controls="panel1a-content"
        >
          <Typography className={classes.heading}>
            Movies waiting for you
          </Typography>
        </AccordionSummary>
        <AccordionDetails
          style={{ height: "62vh", paddingLeft: "0", paddingRight: "0" }}
        >
          <WatchlistProvider />
        </AccordionDetails>
      </Accordion>
      <Accordion
        TransitionProps={{ unmountOnExit: true }}
        expanded={open === "panel2"}
        onChange={handleChange("panel2")}
      >
        <AccordionSummary
          className={classes.summary}
          expandIcon={<ExpandMore />}
          aria-controls="panel1a-content"
        >
          <Typography className={classes.heading}>Your collection</Typography>
        </AccordionSummary>
        <AccordionDetails
          style={{ height: "62vh", paddingLeft: "0", paddingRight: "0" }}
        >
          <RatedMoviesInWatchlist />
        </AccordionDetails>
      </Accordion>
    </Container>
  );
}

const useStyles = makeStyles((theme) => ({
  heading: {
    // fontSize: theme.typography.pxToRem(15),
    margin: "0 auto",
  },
  summary: {
    backgroundColor: "rgba(0, 0, 0, .03)",
    borderRadius: "10px",
  },
}));
