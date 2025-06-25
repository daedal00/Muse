package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daedal00/muse/backend/auth"
	"github.com/daedal00/muse/backend/graph"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "8080"

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	_ = godotenv.Load(".env")
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	srv := handler.New(graph.NewExecutableSchema(
		graph.Config{Resolvers: graph.NewResolver(clientID, clientSecret)},
		),
	)

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
		baseCtx := r.Context()
		newCtx := baseCtx
		// b) extract "Bearer <jwtToken>"
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			// Parse and validate
			token, err := jwt.ParseWithClaims(tokStr, &auth.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
				// Ensure HMAC is used
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return auth.JWTSecret, nil
			})
			if err == nil && token.Valid {
				claims := token.Claims.(*auth.CustomClaims)
				// put user ID into GraphQL context
				newCtx = context.WithValue(baseCtx, graph.UserIDKey, claims.UserID)
			}
		}
		// c) call gqlgen server using r.WithContext(ctx) so resolvers can see it
		srv.ServeHTTP(w, r.WithContext(newCtx))
	}))

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
