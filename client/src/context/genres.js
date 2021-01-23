import create from "zustand";

export const useGenreStore = create((set) => ({
    genres: [],
    setGenres: (gs) => set(() => ({ genres: [...gs] })),
  }));