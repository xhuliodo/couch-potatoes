import React from "react";
import { Route, Switch } from "react-router-dom";
import GettingToKnowUserPage from "../pages/GettingToKnowUserPage";
import SelectingGenrePage from "../pages/SelectingGenrePage";
import Solo from "../pages/Solo";
import WelcomePage from "../pages/WelcomePage";

export default function Routes() {
  return (
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
  );
}
