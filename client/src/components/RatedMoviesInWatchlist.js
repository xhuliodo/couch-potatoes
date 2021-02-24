import { Container } from "@material-ui/core";
import { gql } from "graphql-request";
import React, { useState } from "react";
import { useMutation, useQuery } from "react-query";
import { removeFromWatchlist } from "../utils/deleteFromWatchlist";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import WatchlistCard from "./WatchlistCard";
import "./scrollbar.scss";

export default function RatedMoviesInWatchlist() {
  const graphqlClient = useGraphqlClient();
  const useGetWatchlistHistory = () => {
    return useQuery(["getWatchlistHistory"], async () => {
      const data = (await graphqlClient).request(
        gql`
          query {
            watchlistHistory {
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
  const { isLoading, isError, data } = useGetWatchlistHistory();

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
      {isLoading ? (
        <span>fetching data</span>
      ) : isError ? (
        <span>Something went wrong...</span>
      ) : (
        data.map((m) => (
          <WatchlistCard
            key={m.movieId}
            m={m}
            //   handleRate={handleRate}
            handleRemove={handleRemove}
            deleted={deleted}
            animation={animation}
          />
        ))
      )}
    </Container>
  );
}
