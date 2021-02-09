import create from "zustand";

export const useUserStore = create((set) => ({
  firstName: "",
  profilePic: "",
  token: "",
  setFirstName: (fN) => set(() => ({ firstName: fN })),
  setProfilePic: (pP) => set(() => ({ profilePic: pP })),
  setToken: (t) => set(() => ({ token: t })),
}));
