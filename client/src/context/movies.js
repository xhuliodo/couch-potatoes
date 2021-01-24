import create from "zustand";

export const useMovieStore = create((set) => ({
  movies: [],
  skip: 0,
  limit: 5,
  ratedMovies: 0,
  requiredMovies: 15,
  setMovies: (ms) => set(() => ({ movies: [...ms] })),
  nextMovie: () =>
    set((state) => {
      const currentMovies = state.movies;
      currentMovies.shift();
      return { movies: currentMovies };
    }),
  increaseSkip: () => set((state) => ({ skip: state.skip + state.limit })),
  increaseRequiredMovies: () =>
    set((state) => ({ requiredMovies: state.requiredMovies++ })),
}));
