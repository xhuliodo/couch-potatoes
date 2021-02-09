import "../components/MovieCard.css";

import { useQuery } from "react-query";
import { useGraphqlClient } from "../utils/useGraphqlClient";
import { gql } from "graphql-request";

import { useMovieStore } from "../context/movies";
import MovieCard from "./MovieCard";
import UserFeedbackMovieCard from "./UserFeedbackMovieCard";
import { useAuth0 } from "@auth0/auth0-react";

export default function UserBasedRec() {
  const { peopleToCompare, limit, requiredMovies } = useMovieStore();

  const minimumRatings = requiredMovies / 1.5;

  const graphqlClient = useGraphqlClient();

  const useUserBasedRec = (
    userId,
    minimumRatings,
    peopleToCompare,
    moviesToRecommend
  ) => {
    return useQuery(
      [
        "userBasedRec",
        userId,
        minimumRatings,
        peopleToCompare,
        moviesToRecommend,
      ],
      async () => {
        const data = (await graphqlClient).request(
          gql`
          query{
            recommendFromOtherUsers( 
              userId: "${userId}" 
              minimumRatings: ${minimumRatings} 
              peopleToCompare: ${peopleToCompare} 
              moviesToRecommend: ${moviesToRecommend} 
              ) { 
                movieId 
                posterUrl 
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

  const { user } = useAuth0();
  const { isLoading, isError, data, refetch } = useUserBasedRec(
    user?.sub,
    minimumRatings,
    peopleToCompare,
    limit
  );

  return isLoading ? (
    <UserFeedbackMovieCard message="Fetching movies..." type="loading" />
  ) : isError ? (
    <UserFeedbackMovieCard message="Something went wrong..." />
  ) : (
    <>
      <MovieCard
        startedFromTheBottomNowWeHere={true}
        movies={data}
        refetch={refetch}
      />
    </>
  );
}
