package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daedal00/muse/backend/graph"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver()}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// Wrap query with auth-middleware
	http.Handle("/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// a) extract raw authorization header
		authHeader := r.Header.Get("Authorization")
		// b) extract "Bearer <userID>"
		// TODO: swap from userID to JWT token
		var userID string
		if strings.HasPrefix(authHeader, "Bearer ") {
			userID = strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
		}
		// c) Put userID into graph.UserIDKey context
		ctx := context.WithValue(r.Context(), graph.UserIDKey, userID)
		// d) call gqlgen server using r.WithContext(ctx) so resolvers can see it
		srv.ServeHTTP(w, r.WithContext(ctx))
	}))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
