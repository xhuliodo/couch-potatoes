import {
  BottomNavigation,
  BottomNavigationAction,
  makeStyles,
  Paper,
  useTheme,
} from "@material-ui/core";
import {
  LibraryAddOutlined,
  MovieOutlined,
  PeopleOutline,
} from "@material-ui/icons";
import { useState } from "react";
import UserBasedRec from "../components/UserBasedRec";
import ContentBasedRec from "../components/ContentBasedRec";
import { withAuthenticationRequired } from "@auth0/auth0-react";
import AuthLoading from "../components/AuthLoading";
import { useQuery } from "react-query";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";
import DataStateMovieCard from "../components/DataStateMovieCard";

import SwipeableViews from "react-swipeable-views";
import GeneralWatchlist from "../components/GeneralWatchlist";

export const Solo = (props) => {
  const classes = useStyles();
  const theme = useTheme();

  const [nav, setNav] = useState(1);

  const handleChange = (event, newValue) => {
    setNav(newValue);
  };

  const handleChangeIndex = (index) => {
    setNav(index);
  };

  const graphqlClient = useGraphqlClient();

  // redirect rule for people who have not finished the setup
  const useFavoriteGenres = () => {
    return useQuery("favoriteGenres", async () => {
      const data = await (await graphqlClient).request(
        gql`
          query {
            getFavoriteGenres {
              genreId
            }
          }
        `
      );
      const { getFavoriteGenres } = data;
      if (getFavoriteGenres.length === 0) {
        props.history.push("/getting-to-know-1");
      }
    });
  };

  const { isLoading, isError } = useFavoriteGenres();

  return (
    <Paper elevation={0}>
      <Paper elevation={5} style={{ padding: "12px 0" }}>
        <BottomNavigation
          value={nav}
          onChange={handleChange}
          showLabels
          style={{ width: "fit-content", margin: "0 auto" }}
        >
          <BottomNavigationAction
            style={{ paddingRight: "10px" }}
            label="Our suggestions for you"
            value={0}
            icon={<MovieOutlined />}
            classes={{ selected: classes.selected }}
          />
          <BottomNavigationAction
            style={{ margin: "0 10px" }}
            label="Other people also liked"
            value={1}
            icon={<PeopleOutline />}
            classes={{ selected: classes.selected }}
          />
          <BottomNavigationAction
            style={{ paddingLeft: "10px" }}
            label="Watchlist"
            value={2}
            icon={<LibraryAddOutlined />}
            classes={{ selected: classes.selected }}
          />
        </BottomNavigation>
      </Paper>

      {isLoading ? (
        <DataStateMovieCard message="Fetching movies..." type="loading" />
      ) : isError ? (
        <DataStateMovieCard message="Something went wrong..." />
      ) : (
        // no animation version
        // <div style={{ marginTop: "15px" }}>
        //   {nav === "userBased" ? (
        //     <UserBasedRec />
        //   ) : nav === "genreBased" ? (
        //     <GenreBasedRec startedFromTheBottomNowWeHere={true} />
        //   ) : (
        //     <WatchlistProvider />
        //   )}
        // </div>
        <SwipeableViews
          style={{ marginTop: "15px" }}
          axis={theme.direction === "rtl" ? "x-reverse" : "x"}
          index={nav}
          disabled={true}
          onChangeIndex={handleChangeIndex}
          springConfig={{
            duration: "0.5s",
            easeFunction: "cubic-bezier(0.42, 0, 0.58, 1)",
            delay: "0s",
          }}
        >
          <Panel value={nav} index={0} dir={theme.direction}>
            <ContentBasedRec />
          </Panel>
          <Panel value={nav} index={1} dir={theme.direction}>
            <UserBasedRec />
          </Panel>
          <Panel value={nav} index={2} dir={theme.direction}>
            <GeneralWatchlist />
          </Panel>
        </SwipeableViews>
      )}
    </Paper>
  );
};

const Panel = (props) => {
  const { children, value, index } = props;
  return value === index && <div>{children}</div>;
};

const useStyles = makeStyles((theme) => ({
  selected: { color: `${theme.palette.secondary.dark}!important` },
}));

export default withAuthenticationRequired(Solo, {
  onRedirecting: () => <AuthLoading />,
});
