import create from "zustand";

export const useMovieStore = create((set) => ({
  limit: 5,
  ratedMovies: 0,
  requiredMovies: 15,
  increaseRatedMovies: () =>
    set((state) => ({ ratedMovies: state.ratedMovies + 1 })),
  peopleToCompare: 25,
  recentMoviesToCompare: 10,
}));
