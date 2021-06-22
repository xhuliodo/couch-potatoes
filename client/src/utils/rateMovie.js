export const rateMovie = async (
  { movieId, action, successFunc },
  axiosClient
) => {
  let data = { rateMovie: null };
  switch (action) {
    case "love":
      data = (await axiosClient).post(`/users/ratings`, {
        movieId,
        rating: 1,
      });
      // data = (await graphqlClient).request(
      //   gql`
      //       mutation {
      //         rateMovie(movieId: ${movieId}, rating:1) {
      //           movieId
      //           title
      //         }
      //       }
      //     `
      // );
      break;
    case "hate":
      data = (await axiosClient).post(`/users/ratings`, {
        movieId,
        rating: 0,
      });
      // data = (await graphqlClient).request(
      //   gql`
      //       mutation {
      //         rateMovie(movieId: ${movieId}, rating:0.1) {
      //           movieId
      //           title
      //         }
      //       }
      //     `
      // );
      break;
    default:
      console.log("are you lost? check the util before using it");
  }
  const { statusCode, message } = await data;
  if (statusCode !== 201) {
    // TODO: you should implement some feedback logic, when the rating did not happen
    console.log("the rating didn't do shit");
  } else {
    console.log(message);
    successFunc();
  }
};
