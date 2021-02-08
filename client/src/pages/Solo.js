import {
  BottomNavigation,
  BottomNavigationAction,
  makeStyles,
  Paper,
} from "@material-ui/core";
import { People, WatchLater } from "@material-ui/icons";
import { useState } from "react";
import GenresIcon from "../utils/icons/GenresIcon";
import UserBasedRec from "../components/UserBasedRec";
import GenreBasedRec from "../components/GenreBasedRec";
import WatchlistProvider from "../components/WatchlistProvider";
import { useAuth0, withAuthenticationRequired } from "@auth0/auth0-react";
import AuthLoading from "../components/AuthLoading";
import { useQuery } from "react-query";
import request, { gql } from "graphql-request";
import UserFeedbackMovieCard from "../components/UserFeedbackMovieCard";

export const Solo = (props) => {
  const classes = useStyles();
  const [nav, setNav] = useState("userBased");

  // redirect rule for people who have not finished the setup
  const useFavoriteGenres = (userId) => {
    return useQuery(["genres", userId], async () => {
      const data = await request(
        "http://localhost:4001/graphql",
        gql`
          query {
            User(userId: "${userId}") {
              favoriteGenres {
                genreId
              }
            }
          }
        `
      );
      const { User } = data;
      if (User.length === 0 || User[0].favoriteGenres.length < 3) {
        props.history.push("/getting-to-know-1");
      }
    });
  };

  const { user } = useAuth0();

  const { isLoading, isError } = useFavoriteGenres(user.sub);

  return isLoading ? (
    <UserFeedbackMovieCard message="Fetching movies..." type="loading" />
  ) : isError ? (
    <UserFeedbackMovieCard message="Something went wrong..." />
  ) : (
    <Paper elevation={0}>
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
            classes={{ selected: classes.selected }}
          />
          <BottomNavigationAction
            style={{ margin: "0 10px" }}
            label="Other users also liked"
            value="userBased"
            icon={<People />}
            classes={{ selected: classes.selected }}
          />
          <BottomNavigationAction
            style={{ paddingLeft: "10px" }}
            label="Watchlist"
            value="watchlist"
            icon={<WatchLater />}
            classes={{ selected: classes.selected }}
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

const useStyles = makeStyles((theme) => ({
  selected: { color: `${theme.palette.secondary.dark}!important` },
}));

export default withAuthenticationRequired(Solo, {
  onRedirecting: () => <AuthLoading />,
});
