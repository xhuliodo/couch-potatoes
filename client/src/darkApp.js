import { useState } from "react";
import { ThemeProvider } from "@material-ui/core";
import { createMuiTheme } from "@material-ui/core/styles";
import Paper from "@material-ui/core/Paper";
import Typography from "@material-ui/core/Typography";
import IconButton from "@material-ui/core/IconButton";
import Brightness3Icon from "@material-ui/icons/Brightness3";
import Brightness7Icon from "@material-ui/icons/Brightness7";

function App() {
  const [darkTheme, setDarkTheme] = useState(true);
  const icon = !darkTheme ? <Brightness7Icon /> : <Brightness3Icon />;
  const appliedTheme = createMuiTheme(darkTheme ? light : dark);
  return (
    <ThemeProvider theme={appliedTheme}>
      <Paper>
        <Typography variant="h3">Hello!</Typography>
        <IconButton
          edge="end"
          color="inherit"
          aria-label="mode"
          onClick={() => setDarkTheme(!darkTheme)}
        >
          {icon}
        </IconButton>
      </Paper>
    </ThemeProvider>
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

export default App;
