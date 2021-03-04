import { useCallback, useEffect, useRef, useState } from "react";
import {
  CircularProgress,
  Container,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { useMovieStore } from "../context/movies";
import { useQuery } from "react-query";
import { useMutation } from "react-query";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";
import { removeFromWatchlist } from "../utils/removeFromWatchlist";
import WatchlistCard from "./WatchlistCard";
import "./scrollbar.scss";

export default function RatedMoviesInWatchlist() {
  const classes = useStyle();

  // get data
  const graphqlClient = useGraphqlClient();
  const { limit } = useMovieStore();
  const [skip, setSkip] = useState(0);
  const { isLoading, isError, data } = useGetWatchlistHistory({
    graphqlClient,
    limit,
    skip,
  });

  // infinite scrolling
  const [movies, setMovies] = useState([]);
  useEffect(() => {
    if (!isLoading && !isError) {
      setMovies([...movies, ...data]);
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [data, isError, isLoading]);
  const increaseSkip = () => {
    setSkip(skip + limit);
  };
  const observer = useRef();
  const lastElementRef = useCallback(
    (node) => {
      if (isLoading) {
        return;
      }
      if (observer.current) observer.current.disconnect();
      observer.current = new IntersectionObserver(
        (entries) => {
          if (entries[0].isIntersecting && data.length === limit) {
            increaseSkip();
          }
        },
        { threshold: 1.0 }
      );
      if (node) observer.current.observe(node);
    },
    // eslint-disable-next-line react-hooks/exhaustive-deps
    [isLoading, data]
  );

  // actions user can take
  const remove = useMutation((mutationData) =>
    removeFromWatchlist(mutationData, graphqlClient)
  );
  const handleRemove = (movieId) => {
    const mutationData = {
      movieId,
      successFunc: () => successFunc(movieId, "deleted"),
    };
    remove.mutate(mutationData);
  };

  const successFunc = (movieId, action) => {
    switch (action) {
      case "deleted":
        setAnimation("left");
        break;

      default:
        setAnimation("up");
        break;
    }
    setDeleted([...deleted, movieId]);
  };

  const [deleted, setDeleted] = useState([]);

  const [animation, setAnimation] = useState("up");

  return (
    <Container
      maxWidth="sm"
      style={{
        height: "fit-parent",
      }}
      className="showScroll"
    >
      {movies.map((m, index) => (
        <WatchlistCard
          lastElementRef={movies.length === index + 1 ? lastElementRef : null}
          key={m.movieId}
          m={m}
          handleRemove={handleRemove}
          deleted={deleted}
          animation={animation}
        />
      ))}
      {isLoading ? (
        <div className={classes.loadingDiv}>
          <CircularProgress className={classes.loading} />
        </div>
      ) : isError ? (
        <Typography align="center">Something went wrong...</Typography>
      ) : null}
    </Container>
  );
}

const useStyle = makeStyles(() => ({
  loadingDiv: { width: "100%", display: "flex", marginBottom: "15px" },
  loading: { margin: "10px auto" },
}));

const useGetWatchlistHistory = ({ graphqlClient, skip, limit }) => {
  return useQuery(["getWatchlistHistory", skip, limit], async () => {
    const data = (await graphqlClient).request(
      gql`
          query {
            watchlistHistory(skip:${skip}, limit:${limit}) {
              movieId
              title
              releaseYear
              imdbLink
              rating
            }
          }
        `
    );
    const { watchlistHistory } = await data;
    return watchlistHistory;
  });
};
