import {
  ApolloClient,
  InMemoryCache,
  createHttpLink,
  from,
} from "@apollo/client";
import { setContext } from "@apollo/client/link/context";
import { onError } from "@apollo/client/link/error";

const httpLink = createHttpLink({
  uri: process.env.NEXT_PUBLIC_GRAPHQL_URL || "https://127.0.0.1:8080/query",
  // Add fetch options for development with self-signed certificates
  fetchOptions: {
    // This helps with CORS in development
    mode: "cors",
    credentials: "include",
  },
});

const authLink = setContext((_, { headers }) => {
  // Get the authentication token from local storage if it exists
  const token =
    typeof window !== "undefined" ? localStorage.getItem("auth-token") : null;

  return {
    headers: {
      ...headers,
      authorization: token ? `Bearer ${token}` : "",
      // Add additional headers for better CORS handling
      "Content-Type": "application/json",
    },
  };
});

// Add error handling link
const errorLink = onError(
  ({ graphQLErrors, networkError, operation, forward }) => {
    if (graphQLErrors) {
      graphQLErrors.forEach(({ message, locations, path }) =>
        console.error(
          `[GraphQL error]: Message: ${message}, Location: ${locations}, Path: ${path}`
        )
      );
    }

    if (networkError) {
      console.error(`[Network error]: ${networkError}`);

      // Provide helpful error messages for common HTTPS/CORS issues
      if (
        networkError.message.includes("CORS") ||
        networkError.message.includes("fetch")
      ) {
        console.error(
          "CORS/Network Error - Make sure:\n" +
            "1. Backend server is running on https://127.0.0.1:8080\n" +
            "2. You have accepted the self-signed certificate by visiting https://127.0.0.1:8080 in your browser\n" +
            "3. Backend CORS is configured for https://127.0.0.1:3000"
        );
      }
    }
  }
);

export const apolloClient = new ApolloClient({
  link: from([errorLink, authLink, httpLink]),
  cache: new InMemoryCache(),
  defaultOptions: {
    watchQuery: {
      errorPolicy: "all",
    },
    query: {
      errorPolicy: "all",
    },
  },
});
