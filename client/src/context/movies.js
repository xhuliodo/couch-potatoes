import create from "zustand";

export const useMovieStore = create((set) => ({
  movies: [],
  setMovies: (ms) => set(() => ({ movies: [...ms] })),
  nextMovie: () =>
    set((state) => {
      const currentMovies = state.movies;
      currentMovies.shift();
      return { movies: currentMovies };
    }),
}));
