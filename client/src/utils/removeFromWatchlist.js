export const removeFromWatchlist = async (
  { movieId, successFunc },
  axiosClient
) => {
  const resp = (await axiosClient).delete(`/watchlist/${movieId}`);
  const { status } = await resp;
  // console.log(headers);
  if (status === 204) {
    console.log("you removed from watchlist the movie with id", movieId);
    successFunc();
  } else {
    // TODO: you should implement some feedback logic, when the rating did not happen
    console.log(resp?.data?.message);
  }

  // .request(
  //   gql`
  //     mutation {
  //       removeFromWatchlist(movieId: ${movieId}) {
  //         movieId
  //       }
  //     }
  //   `
  // );

  // const { removeFromWatchlist } = await data;
  // if (removeFromWatchlist?.movieId) {
  //   console.log(
  //     "you removed from watchlist the movie with id",
  //     removeFromWatchlist.movieId
  //   );
  //   successFunc();
  // } else {
  //   // TODO: you should implement some feedback logic, when the rating did not happen
  //   console.log("you shall not remove");
  // }
};
