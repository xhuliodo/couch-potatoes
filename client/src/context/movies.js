import create from "zustand";

export const useMovieStore = create((set) => ({
  movies: [],
  limit: 5,
  ratedMovies: 0,
  requiredMovies: 15,
  currentMovie: null,
  setMovies: (ms) => set(() => ({ movies: [...ms] })),
  nextMovie: () =>
    set((state) => {
      const currentMovies = state.movies;
      currentMovies.shift();
      return { movies: currentMovies };
    }),
  setCurrentMovie: () => set((state) => ({ currentMovie: state.movies[0] })),
  resetSkip: () => set(() => ({ skip: 0 })),
  increaseRatedMovies: () =>
    set((state) => ({ ratedMovies: state.ratedMovies + 1 })),
  peopleToCompare: 25,
}));
