import { useQuery } from "react-query";
import { useAxiosClient } from "../utils/useAxiosClient";
import { useMovieStore } from "../context/movies";
import MovieCard from "./MovieCard";
import "../components/MovieCard.scss";
import DataStateMovieCard from "./DataStateMovieCard";

export default function ContentBasedRec() {
  const { limit } = useMovieStore();
  const axiosClient = useAxiosClient();

  const { isLoading, isError, data, refetch } = useContentBasedRec({
    axiosClient,
    limit,
  });

  return isLoading ? (
    <DataStateMovieCard message="Fetching movies..." type="loading" />
  ) : isError ? (
    <DataStateMovieCard message="Something went wrong..." />
  ) : (
    <MovieCard
      compName="cbr"
      startedFromTheBottomNowWeHere={true}
      movies={data}
      refetch={refetch}
    />
  );
}

const useContentBasedRec = ({ axiosClient, limit: moviesToRecommend }) => {
  return useQuery(["contentBasedRec", moviesToRecommend], async () => {
    const resp = (await axiosClient).get(
      `/recommendations/content-based?limit=${moviesToRecommend}`
    );

    const {
      data: { data: contentBasedRec },
    } = await resp;
    return contentBasedRec;
    // .request(
    //   gql`
    //       query {
    //         recommendFromOtherLikedMovies(
    //           moviesToRecommend: ${moviesToRecommend}
    //         ) {
    //           movieId
    //           title
    //           releaseYear
    //           imdbLink
    //         }
    //       }
    //     `
    // );
    // const { recommendFromOtherLikedMovies } = await data;
    // return recommendFromOtherLikedMovies;
  });
};
