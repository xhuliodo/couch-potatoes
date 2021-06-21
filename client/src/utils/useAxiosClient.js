import { useAuth0 } from "@auth0/auth0-react";
import axios from "axios";

export const useGraphqlClient = async () => {
  const endpoint =
    process.env.REACT_APP_API_ENDPOINT || "http://localhost:4001";

  const { getIdTokenClaims } = useAuth0();
  const { __raw: token } = await getIdTokenClaims();
  const axiosClient = axios({
    baseURL: endpoint,
    headers: { Authorization: `Bearer ${token}` },
  });

  return axiosClient;
};
