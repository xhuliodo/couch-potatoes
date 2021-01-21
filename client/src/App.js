import { useState } from "react";
import MenuBar from "./components/MenuBar";
import SelectingGenre from "./components/SelectingGenre";
import Container from "@material-ui/core/Container";

// dark theme importing
import Brightness7Icon from "@material-ui/icons/Brightness7";
import Brightness3Icon from "@material-ui/icons/Brightness3";
import { createMuiTheme, ThemeProvider } from "@material-ui/core";

// react query importing
import { QueryClient, QueryClientProvider } from "react-query";
import { ReactQueryDevtools } from "react-query/devtools";

export default function App() {
  // dark theme setup
  const [darkTheme, setDarkTheme] = useState(false);
  const icon = !darkTheme ? (
    <Brightness7Icon style={verticalAlign} />
  ) : (
    <Brightness3Icon style={verticalAlign} />
  );
  const appliedTheme = createMuiTheme(!darkTheme ? light : dark);

  // react query setup
  const queryClient = new QueryClient();

  return (
    <QueryClientProvider client={queryClient}>
      <ThemeProvider theme={appliedTheme}>
        <MenuBar
          darkThemeIcon={icon}
          darkTheme={darkTheme}
          setDarkTheme={setDarkTheme}
        />
        <Container maxWidth="md" style={{ marginTop: "5vh" }}>
          <SelectingGenre />
        </Container>
      </ThemeProvider>
      <ReactQueryDevtools initialIsOpen />
    </QueryClientProvider>
  );
}

export const light = {
  palette: {
    type: "light",
  },
};
export const dark = {
  palette: {
    type: "dark",
  },
};

export const verticalAlign = {
  verticalAlign: "middle",
};
