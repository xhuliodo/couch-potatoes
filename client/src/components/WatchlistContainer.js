import { useState } from "react";
import {
  Accordion,
  AccordionDetails,
  AccordionSummary,
  Container,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { ExpandMore } from "@material-ui/icons";
import WatchlistHistory from "./WatchlistHistory";
import WatchlistUnrated from "./WatchlistUnrated";
import "./WatchlistContainer.scss";

export default function WatchlistContainer() {
  const classes = useStyles();

  const [open, setOpen] = useState("panel1");

  const handleChange = (panel) => (event, isExpanded) => {
    setOpen(isExpanded ? panel : false);
  };

  return (
    <Container maxWidth="sm">
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
          <WatchlistUnrated />
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
          <WatchlistHistory />
        </AccordionDetails>
      </Accordion>
    </Container>
  );
}

const useStyles = makeStyles((theme) => ({
  heading: {
    margin: "0 auto",
  },
  summary: {
    backgroundColor: "rgba(0, 0, 0, .08)",
    borderRadius: "10px",
  },
  
}));
