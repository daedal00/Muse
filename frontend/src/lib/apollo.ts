import {
  ApolloClient,
  InMemoryCache,
  createHttpLink,
  from,
} from "@apollo/client";
import { relayStylePagination } from "@apollo/client/utilities";
import { setContext } from "@apollo/client/link/context";

const httpLink = createHttpLink({
  uri: process.env.NEXT_PUBLIC_GRAPHQL_URL || "http://localhost:8080/query",
});

const authLink = setContext((_, { headers }) => {
  // Get the authentication token from local storage if it exists
  const token =
    typeof window !== "undefined" ? localStorage.getItem("auth-token") : null;

  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : "",
    },
  };
});

export const apolloClient = new ApolloClient({
  link: from([authLink, httpLink]),
  cache: new InMemoryCache({
    typePolicies: {
      Query: {
        fields: {
          albums: relayStylePagination(),
          reviews: relayStylePagination(),
          playlists: relayStylePagination(),
        },
      },
    },
  }),
  defaultOptions: {
    watchQuery: {
      errorPolicy: "all",
    },
    query: {
      errorPolicy: "all",
    },
  },
});
