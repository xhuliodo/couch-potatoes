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

import "./GeneralWatchlist.css";

export default function GeneralWatchlist() {
  const classes = useStyles();

  const [open, setOpen] = useState("panel1");

  const handleChange = (panel) => (event, isExpanded) => {
    setOpen(isExpanded ? panel : false);
  };

  return (
    <Container
      maxWidth="sm"
      style={{ margin: "0 auto", paddingLeft: "0", paddingRight: "0" }}
    >
      <Accordion
        style={{ boxShadow: "none" }}
        TransitionProps={{ unmountOnExit: true }}
        expanded={open === "panel1"}
        onChange={handleChange("panel1")}
      >
        <AccordionSummary
          className={classes.summary}
          expandIcon={<ExpandMore />}
          aria-controls="panel1a-content"
          style={{ margin: "0" }}
        >
          <Typography className={classes.heading}>
            Movies waiting for you
          </Typography>
        </AccordionSummary>
        <AccordionDetails className="accordionStyling">
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
        <AccordionDetails className="accordionStyling">
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
    backgroundColor: "rgba(0, 0, 0, .08)",
    borderRadius: "10px",
    // margin: "0!important",
  },
}));
