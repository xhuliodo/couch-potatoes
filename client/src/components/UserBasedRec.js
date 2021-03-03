import "../components/MovieCard.scss";

import { useQuery } from "react-query";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";

import { useMovieStore } from "../context/movies";
import MovieCard from "./MovieCard";
import DataStateMovieCard from "./DataStateMovieCard";

export default function UserBasedRec() {
  const { peopleToCompare, limit, requiredMovies } = useMovieStore();

  const minimumRatings = requiredMovies / 1.5;

  const graphqlClient = useGraphqlClient();

  const useUserBasedRec = (
    minimumRatings,
    peopleToCompare,
    moviesToRecommend
  ) => {
    return useQuery(
      ["userBasedRec", minimumRatings, peopleToCompare, moviesToRecommend],
      async () => {
        const data = (await graphqlClient).request(
          gql`
          query{
            recommendFromOtherUsers( 
              minimumRatings: ${minimumRatings} 
              peopleToCompare: ${peopleToCompare} 
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
        const { recommendFromOtherUsers } = await data;
        return recommendFromOtherUsers;
      }
    );
  };

  const { isLoading, isError, data, refetch } = useUserBasedRec(
    minimumRatings,
    peopleToCompare,
    limit
  );

  return isLoading ? (
    <DataStateMovieCard message="Fetching movies..." type="loading" />
  ) : isError ? (
    <DataStateMovieCard message="Something went wrong..." />
  ) : (
    <MovieCard
      compName="ubr"
      startedFromTheBottomNowWeHere={true}
      movies={data}
      refetch={refetch}
    />
  );
}
