import create from "zustand";

export const useStore = create((set) => ({
    genres: [],
    setGenres: (gs) => set(() => ({ genres: [...gs] })),
  }));