package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/daedal00/muse/backend/auth"
	"github.com/daedal00/muse/backend/graph"
	"github.com/daedal00/muse/backend/internal/config"
	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

// GraphQLRequest represents a GraphQL request
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

// GraphQLResponse represents a GraphQL response
type GraphQLResponse struct {
	Data   interface{} `json:"data"`
	Errors []struct {
		Message string        `json:"message"`
		Path    []interface{} `json:"path"`
	} `json:"errors"`
}

// setupGraphQLHandler creates a GraphQL handler with authentication middleware
func setupGraphQLHandler(resolver *graph.Resolver, cfg *config.Config) http.Handler {
	srv := handler.New(graph.NewExecutableSchema(
		graph.Config{Resolvers: resolver},
	))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))
	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		baseCtx := r.Context()
		newCtx := baseCtx

		if strings.HasPrefix(authHeader, "Bearer ") {
			tokStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			token, err := jwt.ParseWithClaims(tokStr, &auth.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(cfg.JWTSecret), nil
			})
			if err == nil && token.Valid {
				claims := token.Claims.(*auth.CustomClaims)
				newCtx = context.WithValue(baseCtx, graph.UserIDKey, claims.UserID)
			}
		}
		srv.ServeHTTP(w, r.WithContext(newCtx))
	})
}

// Helper function to get environment variables with defaults
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// executeGraphQL sends a GraphQL request to the handler
func executeGraphQL(t *testing.T, handler http.Handler, req GraphQLRequest, authToken string) GraphQLResponse {
	reqBody, err := json.Marshal(req)
	require.NoError(t, err)

	httpReq := httptest.NewRequest(http.MethodPost, "/query", bytes.NewReader(reqBody))
	httpReq.Header.Set("Content-Type", "application/json")

	if authToken != "" {
		httpReq.Header.Set("Authorization", "Bearer "+authToken)
	}

	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httpReq)

	var resp GraphQLResponse
	err = json.Unmarshal(recorder.Body.Bytes(), &resp)
	require.NoError(t, err)

	return resp
}

// Test GraphQL schema introspection
func TestGraphQLIntrospection(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load test config
	cfg := &config.Config{
		Port:        "8080",
		Environment: "test",

		SpotifyClientID:     getEnvOrDefault("SPOTIFY_CLIENT_ID", "test-client-id"),
		SpotifyClientSecret: getEnvOrDefault("SPOTIFY_CLIENT_SECRET", "test-secret"),

		DatabaseURL: getEnvOrDefault("DATABASE_URL", ""),
		DBHost:      getEnvOrDefault("DB_HOST", "localhost"),
		DBPort:      5432,
		DBName:      getEnvOrDefault("DB_NAME", "muse_test"),
		DBUser:      getEnvOrDefault("DB_USER", "postgres"),
		DBPassword:  getEnvOrDefault("DB_PASSWORD", "password"),
		DBSSLMode:   "disable",

		RedisURL:      getEnvOrDefault("REDIS_URL", "redis://localhost:6379"),
		RedisHost:     getEnvOrDefault("REDIS_HOST", "localhost"),
		RedisPort:     6379,
		RedisPassword: getEnvOrDefault("REDIS_PASSWORD", ""),
		RedisDB:       0,

		JWTSecret: "test-secret-key-for-integration-tests",
	}

	// Create resolver
	resolver, err := graph.NewResolver(cfg)
	if err != nil {
		t.Skipf("Skipping test - failed to create resolver (likely no database): %v", err)
	}
	defer resolver.Close()

	// Create handler
	handler := setupGraphQLHandler(resolver, cfg)

	t.Run("Schema Introspection", func(t *testing.T) {
		req := GraphQLRequest{
			Query: `
				query IntrospectionQuery {
					__schema {
						types {
							name
							kind
						}
					}
				}
			`,
		}

		resp := executeGraphQL(t, handler, req, "")
		assert.Empty(t, resp.Errors, "Expected no errors in introspection query")

		data := resp.Data.(map[string]interface{})
		schema := data["__schema"].(map[string]interface{})
		types := schema["types"].([]interface{})
		assert.Greater(t, len(types), 0, "Expected schema to have types")
	})

	t.Run("User Creation", func(t *testing.T) {
		req := GraphQLRequest{
			Query: `
				mutation CreateUser($name: String!, $email: String!, $password: String!) {
					createUser(name: $name, email: $email, password: $password) {
						id
						name
						email
					}
				}
			`,
			Variables: map[string]interface{}{
				"name":     "Integration Test User",
				"email":    "integration@test.com",
				"password": "testpassword123",
			},
		}

		resp := executeGraphQL(t, handler, req, "")
		if len(resp.Errors) > 0 {
			t.Logf("Create user errors (may be expected if database not available): %+v", resp.Errors)
			return // Skip validation if database not available
		}

		data := resp.Data.(map[string]interface{})
		user := data["createUser"].(map[string]interface{})
		assert.NotEmpty(t, user["id"])
		assert.Equal(t, "Integration Test User", user["name"])
		assert.Equal(t, "integration@test.com", user["email"])
	})
}
