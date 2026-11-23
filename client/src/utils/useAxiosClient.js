import { useAuth0 } from "@auth0/auth0-react";
import axios from "axios";
import { apiEndpoint } from "./config";

export const useAxiosClient = async () => {
  const { getIdTokenClaims } = useAuth0();
  const { __raw: token } = await getIdTokenClaims();

  const axiosClient = axios.create({
    baseURL: apiEndpoint,
    headers: {
      authorization: `Bearer ${token}`,
    },
  });

  return axiosClient;
};
