import { useState } from "react";
import {
  BottomNavigation,
  BottomNavigationAction,
  makeStyles,
  Paper,
  useTheme,
} from "@material-ui/core";
import {
  Movie,
  MovieOutlined,
  People,
  PeopleOutline,
  WatchLater,
  WatchLaterOutlined,
} from "@material-ui/icons";
import { withAuthenticationRequired } from "@auth0/auth0-react";
import AuthLoading from "../components/AuthLoading";
import { useQuery } from "react-query";
import { useAxiosClient } from "../utils/useAxiosClient";
import { useMovieStore } from "../context/movies";
import SwipeableViews from "react-swipeable-views";
import ContentBasedRec from "../components/ContentBasedRec";
import UserBasedRec from "../components/UserBasedRec";
import WatchlistContainer from "../components/WatchlistContainer";
import DataStateMovieCard from "../components/DataStateMovieCard";

export const Solo = (props) => {
  const classes = useStyles();
  const theme = useTheme();

  const { requiredMovies } = useMovieStore();

  const [nav, setNav] = useState(1);

  const handleChange = (event, newValue) => {
    setNav(newValue);
  };

  const handleChangeIndex = (index) => {
    setNav(index);
  };

  const axiosClient = useAxiosClient();

  const { isLoading, isError } = useSetupRedirect({
    axiosClient,
    requiredMovies,
    props,
  });

  return (
    <Paper elevation={0}>
      <Paper elevation={3} style={{ padding: "12px 0" }}>
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
            icon={nav === 0 ? <Movie /> : <MovieOutlined />}
            classes={{ selected: classes.selected }}
          />
          <BottomNavigationAction
            style={{ margin: "0 10px" }}
            label="Other people also liked"
            value={1}
            icon={nav === 1 ? <People /> : <PeopleOutline />}
            classes={{ selected: classes.selected }}
          />
          <BottomNavigationAction
            style={{ paddingLeft: "10px" }}
            label="Watch later"
            value={2}
            icon={nav === 2 ? <WatchLater /> : <WatchLaterOutlined />}
            classes={{ selected: classes.selected }}
          />
        </BottomNavigation>
      </Paper>

      {isLoading ? (
        <DataStateMovieCard message="Fetching movies..." type="loading" />
      ) : isError ? (
        <DataStateMovieCard message="Something went wrong..." />
      ) : (
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
            <WatchlistContainer />
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
  selected: {
    color: theme.palette.type === "dark" ? "#fff!important" : "#000!important",
  },
}));

export default withAuthenticationRequired(Solo, {
  onRedirecting: () => <AuthLoading />,
});

// redirect rule for people who have not finished the setup
const useSetupRedirect = ({ axiosClient, props }) => {
  return useQuery("setupRedirect", async () => {
    const finishedSetup = localStorage.getItem("finishedSetup");
    if (!finishedSetup) {
      const resp = await (await axiosClient).get("/users/setup");
      const {
        data: {
          data: { step, finished, message, ratingsGiven },
        },
      } = resp;

      switch (step) {
        case 1:
          props.history.push("/getting-to-know-1");
          return;
        case 2:
          props.history.push("/getting-to-know-2");
          return;
        case 3:
          localStorage.setItem("finishedSetup", finished);
          break;
        default:
          break;
      }
    }
  });
};
