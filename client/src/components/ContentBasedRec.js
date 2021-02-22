import "../components/MovieCard.css";

import { useQuery } from "react-query";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";

import { useMovieStore } from "../context/movies";
import MovieCard from "./MovieCard";
import DataStateMovieCard from "./DataStateMovieCard";

export default function ContentBasedRec() {
  const { recentMoviesToCompare, limit } = useMovieStore();

  const graphqlClient = useGraphqlClient();

  const useContentBasedRec = (recentMoviesToCompare, moviesToRecommend) => {
    return useQuery(
      ["contentBasedRec", recentMoviesToCompare, moviesToRecommend],
      async () => {
        const data = (await graphqlClient).request(
          gql`
            query {
              recommendFromOtherLikedMovies(
                recentMoviesToCompare: ${recentMoviesToCompare}
                moviesToRecommend: ${moviesToRecommend}
              ) {
                movieId
                title
                releaseYear
                imdbLink
              }
            }
          `
        );
        const { recommendFromOtherLikedMovies } = await data;
        return recommendFromOtherLikedMovies;
      }
    );
  };

  const { isLoading, isError, data, refetch } = useContentBasedRec(
    recentMoviesToCompare,
    limit
  );

  return isLoading ? (
    <DataStateMovieCard message="Fetching movies..." type="loading" />
  ) : isError ? (
    <DataStateMovieCard message="Something went wrong..." />
  ) : (
      <MovieCard
        startedFromTheBottomNowWeHere={true}
        movies={data}
        refetch={refetch}
      />
  );
}
