import { useCallback, useEffect, useRef, useState } from "react";
import {
  CircularProgress,
  Container,
  makeStyles,
  Typography,
} from "@material-ui/core";
import { useMutation } from "react-query";
import { removeFromWatchlist } from "../utils/removeFromWatchlist";
import { rateMovie } from "../utils/rateMovie";
import WatchlistCard from "./WatchlistCard";

export default function WatchlistProvider({
  axiosClient,
  statusCode,
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
    if (!isLoading && !isError && statusCode === 200) {
      // console.log("movies", { movies, skip });
      // console.log("data", { data, skip });
      setMovies([...movies, ...data]);
    }
  }, [data, isError, isLoading, skip, statusCode]);
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
          if (
            (entries[0].isIntersecting && data.length === limit) ||
            statusCode === 200
          ) {
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
    rateMovie(mutationData, axiosClient)
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
    removeFromWatchlist(mutationData, axiosClient)
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
    <Container maxWidth="sm" className={classes.mainContainer}>
      {movies.map((m, index) => (
        <WatchlistCard
          lastElementRef={movies.length === index + 1 ? lastElementRef : null}
          key={m.movie.movieId}
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
  mainContainer: {
    height: "fit-parent",
    overflowX: "hidden",
    overflowY: "scroll",
  },
}));
