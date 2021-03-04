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
import { Router } from "react-router-dom";
import history from "./utils/history";

// pages and components
import MenuBar from "./components/MenuBar";
// import Footer from "./components/Footer";
import PageRoutes from "./pages/PageRoutes";
import ServiceWorkerWrapper from "./components/ServiceWorkerWrapper";

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
            <PageRoutes darkTheme={darkTheme} />
          </Router>
          {/* <Footer /> */}
          <ServiceWorkerWrapper />
        </Paper>
      </ThemeProvider>
      {/* <ReactQueryDevtools initialIsOpen={false} /> */}
    </QueryClientProvider>
  );
}

export const verticalAlign = {
  verticalAlign: "middle",
};
