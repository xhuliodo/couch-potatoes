import { useCallback, useEffect, useRef, useState } from "react";
import {
  CircularProgress,
  Container,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { useMutation } from "react-query";
import { removeFromWatchlist } from "../utils/removeFromWatchlist";
import WatchlistCard from "./WatchlistCard";
import "./scrollbar.scss";
import { rateMovie } from "../utils/rateMovie";

export default function WatchlistProvider({
  graphqlClient,
  isLoading,
  isError,
  data,
  skip,
  setSkip,
  limit,
}) {
  const classes = useStyle();

  // infinite scrolling
  const [movies, setMovies] = useState([]);
  useEffect(() => {
    if (!isLoading && !isError) {
      console.log("movies", { movies, skip });
      console.log("data", { data, skip });
      setMovies([...movies, ...data]);
    }
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
  const rate = useMutation((mutationData) =>
    rateMovie(mutationData, graphqlClient)
  );

  const handleRate = (movieId, action) => {
    const mutationData = {
      movieId,
      action,
      successFunc: () => successFunc(movieId),
    };
    rate.mutate(mutationData);
  };

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
          handleRate={handleRate}
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
