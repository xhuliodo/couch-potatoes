import {
  BottomNavigation,
  BottomNavigationAction,
  Paper,
} from "@material-ui/core";
import { People, WatchLater } from "@material-ui/icons";
import { useState } from "react";
import GenresIcon from "../utils/icons/GenresIcon";
import UserBasedRec from "../components/UserBasedRec";
import GenreBasedRec from "../components/GenreBasedRec";
import WatchlistProvider from "../components/WatchlistProvider";
import { withAuthenticationRequired } from "@auth0/auth0-react";
import AuthLoading from "../components/AuthLoading";

export const Solo = () => {
  const [nav, setNav] = useState("watchlist");

  return (
    <Paper elevation={0} style={{ height: "fit-content" }}>
      <Paper elevation={5} style={{ padding: "12px 0" }}>
        <BottomNavigation
          value={nav}
          onChange={(event, newValue) => {
            setNav(newValue);
          }}
          showLabels
          style={{ width: "fit-content", margin: "0 auto" }}
        >
          <BottomNavigationAction
            style={{ paddingRight: "10px" }}
            label="Popular by Genre"
            value="genreBased"
            icon={<GenresIcon />}
          />
          <BottomNavigationAction
            style={{ margin: "0 10px" }}
            label="Other users also liked"
            value="userBased"
            icon={<People />}
          />
          <BottomNavigationAction
            style={{ paddingLeft: "10px" }}
            label="Watchlist"
            value="watchlist"
            icon={<WatchLater />}
          />
        </BottomNavigation>
      </Paper>
      <div style={{ marginTop: "2.5vh" }}>
        {nav === "userBased" ? (
          <UserBasedRec />
        ) : nav === "genreBased" ? (
          <GenreBasedRec startedFromTheBottomNowWeHere={true} />
        ) : (
          <WatchlistProvider />
        )}
      </div>
    </Paper>
  );
};

export default withAuthenticationRequired(Solo, {
  onRedirecting: () => <AuthLoading />,
});
