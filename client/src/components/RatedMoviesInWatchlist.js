import { CircularProgress, Container } from "@material-ui/core";
import { gql } from "graphql-request";
import React, { useCallback, useEffect, useRef, useState } from "react";
import { useMutation, useQuery } from "react-query";
import { removeFromWatchlist } from "../utils/deleteFromWatchlist";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import WatchlistCard from "./WatchlistCard";
import "./scrollbar.scss";
import { useMovieStore } from "../context/movies";

export default function RatedMoviesInWatchlist() {
  const graphqlClient = useGraphqlClient();
  const useGetWatchlistHistory = ({ skip, limit }) => {
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
  const { limit } = useMovieStore();
  const [skip, setSkip] = useState(0);
  const { isLoading, isError, data } = useGetWatchlistHistory({
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
      {movies.map((m, index) => (
        <WatchlistCard
          lastElementRef={movies.length === index + 1 ? lastElementRef : null}
          key={m.movieId}
          m={m}
          //   handleRate={handleRate}
          handleRemove={handleRemove}
          deleted={deleted}
          animation={animation}
        />
      ))}
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
