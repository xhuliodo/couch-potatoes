import { Container, makeStyles } from "@material-ui/core";
import React, { useEffect, useMemo, useState } from "react";
import { Route, Switch, useLocation } from "react-router-dom";

import WelcomePage from "./WelcomePage";
import GettingToKnowUserPage from "./GettingToKnowUserPage";
import SelectingGenrePage from "./SelectingGenrePage";
import Solo from "./Solo";

export default function PageRoutes({ darkTheme }) {
  const classes = useStyle();
  // landing page conditional styling
  const location = useLocation();

  const defaultStyle = useMemo(() => ({ marginTop: "15px" }), []);
  const [landingStyling, setLandingStyling] = useState(defaultStyle);

  useEffect(() => {
    if (location.pathname === "/") {
      setLandingStyling({
        paddingLeft: "0",
        paddingRight: "0",
        margin: "0",
        maxWidth: "100%",
        backgroundColor: `${darkTheme ? "#262626" : "#cecece"}`,
      });
    } else {
      setLandingStyling(defaultStyle);
    }
  }, [darkTheme, defaultStyle, location]);

  return (
    <Container
      maxWidth="md"
      className={classes.fullScreenSecondLevel}
      style={{ ...landingStyling }}
    >
      <Switch>
        <Route exact path="/" component={WelcomePage} />
        <Route exact path="/getting-to-know-1" component={SelectingGenrePage} />
        <Route
          exact
          path="/getting-to-know-2"
          component={GettingToKnowUserPage}
        />
        <Route exact path="/solo" component={Solo} />
      </Switch>
    </Container>
  );
}

const useStyle = makeStyles(() => ({
  fullScreenSecondLevel: {
    flexGrow: "3",
    alignContent: "flex-start",
    borderRadius: "0px!important",
    paddingBottom: "15px",
  },
}));
