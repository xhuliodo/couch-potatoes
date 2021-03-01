import { CircularProgress, Container } from "@material-ui/core";
import WatchlistCard from "./WatchlistCard";

import { useQuery } from "react-query";
import { useMutation } from "react-query";
import { rateMovie } from "../utils/rateMovie";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";
import { removeFromWatchlist } from "../utils/deleteFromWatchlist";
import { useCallback, useEffect, useRef, useState } from "react";

import "./scrollbar.scss";
import { useMovieStore } from "../context/movies";

export default function WatchlistProvider() {
  const graphqlClient = useGraphqlClient();
  const useGetWatchlistMovies = ({ skip, limit }) => {
    return useQuery(["getWatchlistMovies", skip, limit], async () => {
      const data = (await graphqlClient).request(
        gql`
          query {
            watchlist(skip:${skip}, limit:${limit}) {
              movieId
              title
              releaseYear
              imdbLink
            }
          }
        `
      );
      const { watchlist } = await data;
      return watchlist;
    });
  };
  const { limit } = useMovieStore();
  const [skip, setSkip] = useState(0);
  const { isLoading, isError, data } = useGetWatchlistMovies({
    limit,
    skip,
  });

  const [movies, setMovies] = useState([]);

  useEffect(() => {
    if (!isLoading && !isError) {
      setMovies([...movies, ...data]);
    }
  }, [data, isError, isLoading]);

  const increaseSkip = () => {
    setSkip(skip + limit);
    // refetch();
  };

  const observer = useRef();
  const lastElementRef = useCallback(
    (node) => {
      if (isLoading) {
        console.log("currently loading");
        return;
      }
      if (observer.current) observer.current.disconnect();
      observer.current = new IntersectionObserver((entries) => {
        if (entries[0].isIntersecting && data.length === limit) {
          console.log(skip);
          increaseSkip();
        }
      });
      if (node) observer.current.observe(node);
    },
    [isLoading, data]
  );

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
    // for (var i = 0; i < data.length; i++) {
    //   if (data[i].movieId === movieId) {
    //     data.splice(i, 1);
    //     i--;
    //   }
    // }
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
      {
        movies.map((m, index) => (
          <WatchlistCard
            lastElementRef={movies.length === index + 1 ? lastElementRef : null}
            key={m.movieId}
            m={m}
            handleRate={handleRate}
            handleRemove={handleRemove}
            deleted={deleted}
            animation={animation}
          />
        ))
        // )
      }
      {isLoading ? (
        <div style={{ width: "100%", display: "flex" }}>
          <CircularProgress style={{ margin: "10px auto" }} />
        </div>
      ) : isError ? (
        <span style={{ textAlign: "center" }}>Something went wrong...</span>
      ) : null}
    </Container>
  );
}
