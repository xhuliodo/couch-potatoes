export const removeFromWatchlist = async (
  { movieId, successFunc },
  axiosClient
) => {
  (await axiosClient)
    .delete(`/watchlist/${movieId}`)
    .then((resp) => {
      const {
        headers: { status },
      } = resp;
      if (status === 204) {
        console.log("you removed from watchlist the movie with id", movieId);
        successFunc();
      } else {
        // TODO: you should implement some feedback logic, when the rating did not happen
        console.log(resp?.data?.message);
      }
    })
    .catch((err) => {
      console.log(err);
    });
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
