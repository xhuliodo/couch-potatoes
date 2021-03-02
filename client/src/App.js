import { useEffect, useState } from "react";

// dark theme importing
import {
  createMuiTheme,
  Paper,
  ThemeProvider,
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
import PageRoutes from "./pages/PageRoutes";

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
            // "-webkit-appearance": "none",
            width: "8px",
            height: "8px",
          },
          "*::-webkit-scrollbar-thumb": {
            backgroundColor: "grey",
          },
          "::-webkit-scrollbar-thumb": {
            // background: "#888",
            borderRadius: "15px",
            // "-webkit-overflow-scrolling": "auto",
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
      primary: { main: "#8e72a7", dark: "#8e72a7" },
      secondary: { main: "#746d4f", dark: "#746d4f" },
      // background: { paper: "#c4c4c4" },
      // text: { primary: "#fff", secondary:"#fff" },
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
            height: "100vh",
            display: "flex",
            flexDirection: "column",
            borderRadius: "0px",
          }}
        >
          <Router history={history}>
            <MenuBar darkTheme={darkTheme} setDarkTheme={setDarkTheme} />
            <PageRoutes />
          </Router>
          {/* <Footer /> */}
        </Paper>
      </ThemeProvider>
      {/* <ReactQueryDevtools initialIsOpen={false} /> */}
    </QueryClientProvider>
  );
}

export const verticalAlign = {
  verticalAlign: "middle",
};
