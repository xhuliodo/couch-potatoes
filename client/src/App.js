import { useState } from "react";

// dark theme importing
import {
  createMuiTheme,
  Paper,
  ThemeProvider,
  CssBaseline,
  makeStyles,
  useTheme,
} from "@material-ui/core";

// react query importing
import { QueryClient, QueryClientProvider } from "react-query";

// app routing
import { Router } from "react-router-dom";
import history from "./utils/history";

// pages and components
import MenuBar from "./components/MenuBar";
import PageRoutes from "./pages/PageRoutes";
import ServiceWorkerWrapper from "./components/ServiceWorkerWrapper";

export default function App() {
  const classes = useStyle();
  const theme = useTheme();
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
            width: "8px",
            height: "8px",
          },
          "*::-webkit-scrollbar-thumb": {
            backgroundColor: "grey",
          },
          "::-webkit-scrollbar-thumb": {
            borderRadius: "15px",
          },
          "*:hover": {
            "&::-webkit-scrollbar-thumb": {
              backgroundColor: "darkgrey",
            },
          },
          body: {
            backgroundColor: !darkTheme
              ? theme.palette.background.paper
              : "#424242",
          },
        },
      },
    },
    palette: {
      type: !darkTheme ? "light" : "dark",
      primary: { main: "#8e72a7", dark: "#8e72a7" },
      secondary: { main: "#746d4f", dark: "#746d4f" },
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
        <CssBaseline />
        <Paper elevation={0} className={classes.fullScreen}>
          <Router history={history}>
            <MenuBar darkTheme={darkTheme} setDarkTheme={setDarkTheme} />
            <PageRoutes darkTheme={darkTheme} />
          </Router>
          <ServiceWorkerWrapper />
        </Paper>
      </ThemeProvider>
    </QueryClientProvider>
  );
}

const useStyle = makeStyles(() => ({
  fullScreen: {
    height: "100vh",
    display: "flex",
    flexDirection: "column",
    borderRadius: "0px",
  },
}));
