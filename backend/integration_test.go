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
	"time"

	"github.com/daedal00/muse/backend/auth"
	"github.com/daedal00/muse/backend/graph"
	"github.com/daedal00/muse/backend/internal/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/golang-jwt/jwt/v5"
	"github.com/vektah/gqlparser/v2/ast"
)

// GraphQL request/response types for testing
type GraphQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables,omitempty"`
}

type GraphQLResponse struct {
	Data   interface{} `json:"data,omitempty"`
	Errors []struct {
		Message string        `json:"message"`
		Path    []interface{} `json:"path,omitempty"`
	} `json:"errors,omitempty"`
}

// Helper function to get environment variable with default
func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// Helper function to detect CI environment
func isCI() bool {
	ciEnvVars := []string{"CI", "GITHUB_ACTIONS", "GITLAB_CI", "JENKINS_URL", "BUILDKITE", "CIRCLECI"}
	for _, envVar := range ciEnvVars {
		if os.Getenv(envVar) != "" {
			return true
		}
	}
	return false
}

// Setup GraphQL handler for testing (mirroring server.go)
func setupGraphQLHandler(resolver *graph.Resolver, cfg *config.Config) http.Handler {
	// Create GraphQL server (same as server.go)
	srv := handler.New(graph.NewExecutableSchema(
		graph.Config{Resolvers: resolver},
	))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})
	srv.AddTransport(transport.Websocket{
		KeepAlivePingInterval: 10 * time.Second,
	})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	// Create a simple mux for testing
	mux := http.NewServeMux()

	// Add playground handler
	mux.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// Add GraphQL handler with auth middleware (same as server.go)
	mux.Handle("/query", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract raw authorization header
		authHeader := r.Header.Get("Authorization")
		baseCtx := r.Context()
		newCtx := baseCtx

		// Extract "Bearer <jwtToken>"
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			// Parse and validate
			token, err := jwt.ParseWithClaims(tokStr, &auth.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
				// Ensure HMAC is used
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(cfg.JWTSecret), nil
			})
			if err == nil && token.Valid {
				claims := token.Claims.(*auth.CustomClaims)
				// Put user ID into GraphQL context
				newCtx = context.WithValue(baseCtx, graph.UserIDKey, claims.UserID)
			}
		}
		// Call gqlgen server using r.WithContext(ctx) so resolvers can see it
		srv.ServeHTTP(w, r.WithContext(newCtx))
	}))

	// Add health check endpoint
	mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	return mux
}

// Execute GraphQL query for testing
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

	// Load test config with better error handling
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

	// Create resolver with better error handling
	resolver, err := graph.NewResolver(cfg)
	if err != nil {
		if isCI() {
			// In CI environments, provide more detailed error information
			t.Logf("CI Environment detected. Database connection failed: %v", err)
			t.Logf("Config: Host=%s, Port=%d, User=%s, Database=%s",
				cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBName)

			// Check if it's a connection issue vs configuration issue
			if cfg.DatabaseURL != "" {
				t.Logf("Using DATABASE_URL: %s", maskSensitiveInfo(cfg.DatabaseURL))
			} else {
				dsn := cfg.GetDatabaseDSN()
				t.Logf("Generated DSN: %s", maskSensitiveInfo(dsn))
			}
		}
		t.Skipf("Skipping test - failed to create resolver (likely no database): %v", err)
	}
	defer resolver.Close()

	// Create handler
	handler := setupGraphQLHandler(resolver, cfg)

	// Test schema introspection query
	introspectionQuery := GraphQLRequest{
		Query: `
			query IntrospectionQuery {
				__schema {
					queryType { name }
					mutationType { name }
					subscriptionType { name }
					types {
						name
						kind
					}
				}
			}
		`,
	}

	resp := executeGraphQL(t, handler, introspectionQuery, "")

	// Verify response
	assert.Nil(t, resp.Errors, "Introspection query should not have errors")
	assert.NotNil(t, resp.Data, "Introspection query should return data")

	// Verify basic schema structure
	data, ok := resp.Data.(map[string]interface{})
	require.True(t, ok, "Response data should be a map")

	schema, ok := data["__schema"].(map[string]interface{})
	require.True(t, ok, "Schema should be present")

	queryType, ok := schema["queryType"].(map[string]interface{})
	require.True(t, ok, "Query type should be present")
	assert.Equal(t, "Query", queryType["name"], "Query type should be named 'Query'")

	mutationType, ok := schema["mutationType"].(map[string]interface{})
	require.True(t, ok, "Mutation type should be present")
	assert.Equal(t, "Mutation", mutationType["name"], "Mutation type should be named 'Mutation'")

	subscriptionType, ok := schema["subscriptionType"].(map[string]interface{})
	require.True(t, ok, "Subscription type should be present")
	assert.Equal(t, "Subscription", subscriptionType["name"], "Subscription type should be named 'Subscription'")

	// Verify some expected types exist
	types, ok := schema["types"].([]interface{})
	require.True(t, ok, "Types should be an array")

	expectedTypes := []string{"User", "Album", "Track", "Review", "Playlist"}
	foundTypes := make(map[string]bool)

	for _, typeInterface := range types {
		if typeMap, ok := typeInterface.(map[string]interface{}); ok {
			if typeName, ok := typeMap["name"].(string); ok {
				foundTypes[typeName] = true
			}
		}
	}

	for _, expectedType := range expectedTypes {
		assert.True(t, foundTypes[expectedType], "Type '%s' should exist in schema", expectedType)
	}

	t.Logf("✅ GraphQL schema introspection test completed successfully")
	t.Logf("Found %d types in schema", len(types))
}

// Test health check endpoint
func TestHealthCheck(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Load minimal config for health check
	cfg := &config.Config{
		Port:                "8080",
		Environment:         "test",
		SpotifyClientID:     "test-client-id",
		SpotifyClientSecret: "test-secret",
		DatabaseURL:         getEnvOrDefault("DATABASE_URL", "postgres://test:test@localhost/test"),
		JWTSecret:           "test-secret-key",
	}

	// Create resolver - if it fails, skip the test
	resolver, err := graph.NewResolver(cfg)
	if err != nil {
		t.Skipf("Skipping health check test - failed to create resolver: %v", err)
	}
	defer resolver.Close()

	// Create handler
	handler := setupGraphQLHandler(resolver, cfg)

	// Test health check
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, req)

	assert.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	assert.JSONEq(t, `{"status":"ok"}`, recorder.Body.String())
}

// maskSensitiveInfo masks passwords and sensitive information in connection strings for logging
func maskSensitiveInfo(connectionString string) string {
	// Simple masking for passwords in connection strings
	// This is a basic implementation - in production you'd want more sophisticated masking
	if len(connectionString) == 0 {
		return connectionString
	}

	// For URL format: postgres://user:password@host:port/db
	if len(connectionString) > 20 {
		return connectionString[:10] + "***MASKED***" + connectionString[len(connectionString)-10:]
	}

	// For key=value format, just show length
	return fmt.Sprintf("***MASKED_CONNECTION_STRING_LENGTH_%d***", len(connectionString))
}
