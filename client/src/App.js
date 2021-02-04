import { useState } from "react";

// dark theme importing
import Brightness7Icon from "@material-ui/icons/Brightness7";
import Brightness3Icon from "@material-ui/icons/Brightness3";
import {
  createMuiTheme,
  Paper,
  ThemeProvider,
  Container,
} from "@material-ui/core";

// react query importing
import { QueryClient, QueryClientProvider } from "react-query";
import { ReactQueryDevtools } from "react-query/devtools";

// app routing
import { Router, Route, Switch } from "react-router-dom";
import history from "./utils/history";

// pages and components
import MenuBar from "./components/MenuBar";
import Footer from "./components/Footer";
import WelcomePage from "./pages/WelcomePage";
import GettingToKnowUserPage from "./pages/GettingToKnowUserPage";
import SelectingGenrePage from "./pages/SelectingGenrePage";
import Solo from "./pages/Solo";

export default function App() {
  // dark theme setup
  const [darkTheme, setDarkTheme] = useState(false);
  const icon = !darkTheme ? (
    <Brightness7Icon style={verticalAlign} />
  ) : (
    <Brightness3Icon style={verticalAlign} />
  );

  const customTheme = {
    palette: {
      type: !darkTheme ? "light" : "dark",
    },
  };
  const appliedTheme = createMuiTheme(customTheme);

  // react query setup
  const queryClient = new QueryClient({
    defaultOptions: { queries: { refetchOnWindowFocus: false } },
  });

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={appliedTheme}>
        <Paper
          elevation={0}
          style={{
            minHeight: "100vh",
          }}
        >
          <MenuBar
            darkThemeIcon={icon}
            darkTheme={darkTheme}
            setDarkTheme={setDarkTheme}
          />
          <Container
            maxWidth="md"
            style={{
              marginTop: "5vh",
              height: "fit-content",
            }}
          >
            <Router history={history}>
              <Switch>
                <Route exact path="/" component={WelcomePage} />
                <Route
                  exact
                  path="/getting-to-know-1"
                  component={SelectingGenrePage}
                />
                <Route
                  exact
                  path="/getting-to-know-2"
                  component={GettingToKnowUserPage}
                />
                <Route exact path="/solo" component={Solo} />
              </Switch>
            </Router>
          </Container>
          <Footer />
        </Paper>
      </ThemeProvider>
      <ReactQueryDevtools initialIsOpen />
    </QueryClientProvider>
  );
}

export const verticalAlign = {
  verticalAlign: "middle",
};
