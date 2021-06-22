import { useQuery } from "react-query";
import { useAxiosClient } from "../utils/useAxiosClient";
import { useMovieStore } from "../context/movies";
import MovieCard from "./MovieCard";
import "../components/MovieCard.scss";
import DataStateMovieCard from "./DataStateMovieCard";

export default function UserBasedRec() {
  const {
    // peopleToCompare,
    limit,
    //  requiredMovies
  } = useMovieStore();
  // const minimumRatings = requiredMovies / 1.5;
  const axiosClient = useAxiosClient();

  const { isLoading, isError, data, refetch } = useUserBasedRec({
    axiosClient,
    limit,
  });

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

const useUserBasedRec = ({ axiosClient, limit: moviesToRecommend }) => {
  return useQuery(
    ["userBasedRec", axiosClient, moviesToRecommend],
    async () => {
      (await axiosClient)
        .get(`/recommendations/user-based?limit=${moviesToRecommend}`)
        .then((resp) => {
          const { data } = resp;
          return data;
        })
        .catch((err) => {
          console.log(err);
        });
      // const data = (await graphqlClient).request(
      //   gql`
      //     query{
      //       recommendFromOtherUsers(
      //         minimumRatings: ${minimumRatings}
      //         peopleToCompare: ${peopleToCompare}
      //         moviesToRecommend: ${moviesToRecommend}
      //         ) {
      //           movieId
      //           title
      //           releaseYear
      //           imdbLink
      //           }
      //       }
      //   `
      // );
      // const { recommendFromOtherUsers } = await data;
      // return recommendFromOtherUsers;
    }
  );
};
