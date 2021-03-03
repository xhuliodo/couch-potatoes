import { Container } from "@material-ui/core";
import React, { useEffect, useState } from "react";
import { Route, Switch, useLocation } from "react-router-dom";

import WelcomePage from "./WelcomePage";
import GettingToKnowUserPage from "./GettingToKnowUserPage";
import SelectingGenrePage from "./SelectingGenrePage";
import Solo from "./Solo";

export default function PageRoutes({ darkTheme }) {
  // landing page conditional styling
  const location = useLocation();

  const [landingStyling, setLandingStyling] = useState({
    marginTop: "15px",
  });
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
      setLandingStyling({ marginTop: "15px" });
    }
  }, [location]);

  return (
    <Container
      maxWidth="md"
      style={{
        flexGrow: "3",
        alignContent: "flex-start",
        borderRadius: "0px!important",
        ...landingStyling,
      }}
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
