import { useAuth0 } from "@auth0/auth0-react";
import { GraphQLClient } from "graphql-request";

export const useGraphqlClient = async () => {
  const endpoint =
    process.env.REACT_APP_API_ENDPOINT || "http://localhost:4001/graphql";
  const graphqlClient = new GraphQLClient(endpoint);

  const { getIdTokenClaims } = useAuth0();
  const { __raw: token } = await getIdTokenClaims();
  graphqlClient.setHeader("Authorization", `Bearer ${token}`);

  return graphqlClient;
};
