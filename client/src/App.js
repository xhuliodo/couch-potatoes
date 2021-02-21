import { useEffect, useState } from "react";

// dark theme importing
import {
  createMuiTheme,
  Paper,
  ThemeProvider,
  Container,
  CssBaseline,
} from "@material-ui/core";

// react query importing
import { QueryClient, QueryClientProvider } from "react-query";
// import { ReactQueryDevtools } from "react-query/devtools";

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
  const userPref = localStorage.getItem("darkMode");
  const isDarkThemeOn =
    // eslint-disable-next-line eqeqeq
    userPref == "false" ? false : userPref == "true" ? true : false;
  const [darkTheme, setDarkTheme] = useState(isDarkThemeOn);

  const customTheme = {
    overrides: {
      MuiCssBaseline: {
        "@global": {
          "*::-webkit-scrollbar": {
            "-webkit-appearance": "none",
            width: "8px",
            height: "8px",
          },
          "*::-webkit-scrollbar-thumb": {
            backgroundColor: "grey",
          },
          "::-webkit-scrollbar-thumb": {
            // background: "#888",
            borderRadius: "15px",
            "-webkit-overflow-scrolling": "auto",
          },
          "*:hover": {
            "&::-webkit-scrollbar-thumb": {
              backgroundColor: "darkgrey",
            },
          },
          /* Equivalent alternative:
          "*:hover::-webkit-scrollbar-thumb": {
            backgroundColor: "green"
          }
           */
        },
      },
    },
    palette: {
      type: !darkTheme ? "light" : "dark",
    },
  };
  const appliedTheme = createMuiTheme(customTheme);

  // react query setup
  const queryClient = new QueryClient({
    defaultOptions: { queries: { refetchOnWindowFocus: false } },
  });

  // auto updating pwa when new update is out
  useEffect(() => {
    if (window.swUpdateReady) {
      window.swUpdateReady = false;
      window.stop();
      window.location.reload();
    }
  }, []);

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={appliedTheme}>
        <CssBaseline />
        <Paper
          elevation={0}
          style={{
            minHeight: "100vh",
            borderRadius: "0px",
          }}
        >
          <MenuBar darkTheme={darkTheme} setDarkTheme={setDarkTheme} />
          <Container
            maxWidth="md"
            style={{
              marginTop: "15px",
              height: "fit-content",
              borderRadius: "0px!important",
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
      {/* <ReactQueryDevtools initialIsOpen={false} /> */}
    </QueryClientProvider>
  );
}

export const verticalAlign = {
  verticalAlign: "middle",
};
