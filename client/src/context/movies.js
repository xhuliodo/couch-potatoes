import create from "zustand";

export const useMovieStore = create((set) => ({
  movies: [],
  skip: 0,
  limit: 5,
  ratedMovies: 0,
  requiredMovies: 5,
  currentMovie: null,
  setMovies: (ms) => set(() => ({ movies: [...ms] })),
  nextMovie: () =>
    set((state) => {
      const currentMovies = state.movies;
      currentMovies.shift();
      return { movies: currentMovies };
    }),
  setCurrentMovie: () => set((state) => ({ currentMovie: state.movies[0] })),
  increaseSkip: () => set((state) => ({ skip: state.skip + state.limit })),
  increaseRatedMovies: () =>
    set((state) => ({ ratedMovies: state.ratedMovies + 1 })),
}));
