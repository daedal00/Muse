package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/daedal00/muse/backend/auth"
	"github.com/daedal00/muse/backend/graph"
	"github.com/daedal00/muse/backend/internal/config"
	"github.com/daedal00/muse/backend/internal/models"
	spotifyinternal "github.com/daedal00/muse/backend/internal/spotify"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/vektah/gqlparser/v2/ast"
)

// Request logging middleware
func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Log incoming request
		log.Printf("[REQUEST] %s %s from %s - User-Agent: %s",
			r.Method, r.URL.Path, r.RemoteAddr, r.UserAgent())

		// Check for GraphQL query in body for POST requests
		if r.Method == "POST" && strings.Contains(r.Header.Get("Content-Type"), "application/json") {
			log.Printf("[GRAPHQL] Processing GraphQL request")
		}

		// Create a response writer wrapper to capture status code
		wrapped := &responseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		// Call next handler
		next.ServeHTTP(wrapped, r)

		// Log response
		duration := time.Since(start)
		log.Printf("[RESPONSE] %s %s - Status: %d - Duration: %v",
			r.Method, r.URL.Path, wrapped.statusCode, duration)
	})
}

// Response writer wrapper to capture status code
type responseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

// CORS middleware to handle cross-origin requests
func corsMiddleware(allowedOrigin string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Printf("[CORS] Request from origin: %s", origin)

		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			log.Printf("[CORS] Handling preflight OPTIONS request")
			w.WriteHeader(http.StatusOK)
			return
		}

		// Continue to next handler
		next.ServeHTTP(w, r)
	})
}

func main() {
	log.Println("🚀 Starting Muse Backend Server...")

	// Load configuration
	log.Println("[CONFIG] Loading configuration...")
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("[ERROR] Failed to load configuration: %v", err)
	}
	log.Printf("[CONFIG] Server will run on port %s in %s environment", cfg.Port, cfg.Environment)

	// Initialize resolver with database and Redis connections
	log.Println("[INIT] Initializing database and Redis connections...")
	resolver, err := graph.NewResolver(cfg)
	if err != nil {
		log.Fatalf("[ERROR] Failed to initialize resolver: %v", err)
	}
	defer resolver.Close()
	log.Println("[INIT] ✅ Database and Redis connections established")

	// Create GraphQL server
	log.Println("[GRAPHQL] Setting up GraphQL server...")
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
	log.Println("[GRAPHQL] ✅ GraphQL server configured")

	// Set up routes
	log.Println("[ROUTES] Setting up HTTP routes...")
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))

	// Wrap query with CORS, logging, and auth middleware
	http.Handle("/query", corsMiddleware(cfg.FrontendURL, loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract raw authorization header
		authHeader := r.Header.Get("Authorization")
		baseCtx := r.Context()
		newCtx := baseCtx

		var userID string
		var authStatus string

		// Extract "Bearer <jwtToken>"
		if strings.HasPrefix(authHeader, "Bearer ") {
			tokStr := strings.TrimSpace(strings.TrimPrefix(authHeader, "Bearer "))
			log.Printf("[AUTH] Processing JWT token (length: %d)", len(tokStr))

			// Parse and validate
			token, err := jwt.ParseWithClaims(tokStr, &auth.CustomClaims{}, func(t *jwt.Token) (interface{}, error) {
				// Ensure HMAC is used
				if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("unexpected signing method")
				}
				return []byte(cfg.JWTSecret), nil
			})

			if err != nil {
				log.Printf("[AUTH] ❌ JWT validation failed: %v", err)
				authStatus = "invalid"
			} else if token.Valid {
				claims := token.Claims.(*auth.CustomClaims)
				userID = claims.UserID
				// Put user ID into GraphQL context
				newCtx = context.WithValue(baseCtx, graph.UserIDKey, claims.UserID)
				log.Printf("[AUTH] ✅ User authenticated: %s", userID)
				authStatus = "authenticated"
			} else {
				log.Printf("[AUTH] ❌ Invalid token")
				authStatus = "invalid"
			}
		} else if authHeader != "" {
			log.Printf("[AUTH] ❌ Invalid authorization header format")
			authStatus = "malformed"
		} else {
			log.Printf("[AUTH] No authorization header - anonymous request")
			authStatus = "anonymous"
		}

		log.Printf("[AUTH] Request status: %s, UserID: %s", authStatus, userID)

		// Call gqlgen server using r.WithContext(ctx) so resolvers can see it
		srv.ServeHTTP(w, r.WithContext(newCtx))
	}))))

	http.Handle("/spotify/callback", corsMiddleware(cfg.FrontendURL, loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		query := r.URL.Query()
		if errParam := query.Get("error"); errParam != "" {
			redirect := cfg.FrontendURL
			if redirect != "" {
				if redirectURL, err := url.Parse(redirect); err == nil {
					values := redirectURL.Query()
					values.Set("spotify", "error")
					values.Set("reason", errParam)
					redirectURL.RawQuery = values.Encode()
					http.Redirect(w, r, redirectURL.String(), http.StatusFound)
					return
				}
			}
			http.Error(w, "spotify authorization failed", http.StatusBadRequest)
			return
		}

		code := query.Get("code")
		state := query.Get("state")
		if code == "" || state == "" {
			http.Error(w, "missing code or state", http.StatusBadRequest)
			return
		}

		claims := &graph.SpotifyStateClaims{}
		parsedToken, err := jwt.ParseWithClaims(state, claims, func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(cfg.JWTSecret), nil
		})
		if err != nil || !parsedToken.Valid {
			http.Error(w, "invalid state", http.StatusBadRequest)
			return
		}

		userID, err := uuid.Parse(claims.UserID)
		if err != nil {
			http.Error(w, "invalid user id", http.StatusBadRequest)
			return
		}

		if cfg.SpotifyClientID == "" || cfg.SpotifyClientSecret == "" {
			http.Error(w, "spotify not configured", http.StatusBadRequest)
			return
		}

		spotifyClient := spotifyinternal.NewClient(spotifyinternal.Config{
			ClientID:     cfg.SpotifyClientID,
			ClientSecret: cfg.SpotifyClientSecret,
			RedirectURL:  cfg.SpotifyRedirectURL,
			Scopes:       graph.SpotifyUserScopes(),
		})

		oauthToken, err := spotifyClient.ExchangeCode(r.Context(), code)
		if err != nil {
			http.Error(w, "failed to exchange spotify code", http.StatusBadRequest)
			return
		}

		spotifyUserID := (*string)(nil)
		client := spotifyClient.GetAuthorizedClient(r.Context(), oauthToken)
		services := spotifyinternal.NewServices(client)
		profile, err := services.User.GetCurrentUser(r.Context())
		if err == nil {
			id := string(profile.ID)
			spotifyUserID = &id
		}

		var tokenType *string
		if oauthToken.TokenType != "" {
			val := oauthToken.TokenType
			tokenType = &val
		}

		var scope *string
		if rawScope := oauthToken.Extra("scope"); rawScope != nil {
			if scopeStr, ok := rawScope.(string); ok {
				scope = &scopeStr
			}
		}

		now := time.Now()
		dbToken := &models.SpotifyToken{
			UserID:        userID,
			SpotifyUserID: spotifyUserID,
			AccessToken:   oauthToken.AccessToken,
			RefreshToken:  oauthToken.RefreshToken,
			TokenType:     tokenType,
			Scope:         scope,
			ExpiresAt:     oauthToken.Expiry,
			CreatedAt:     now,
			UpdatedAt:     now,
		}

		if err := resolver.Repos().Spotify.Upsert(r.Context(), dbToken); err != nil {
			http.Error(w, "failed to save spotify token", http.StatusInternalServerError)
			return
		}

		redirect := resolver.SanitizeRedirectURI(claims.RedirectURI)
		if redirect == "" {
			redirect = cfg.FrontendURL
		}

		redirectURL, err := url.Parse(redirect)
		if err != nil {
			http.Redirect(w, r, cfg.FrontendURL, http.StatusFound)
			return
		}
		values := redirectURL.Query()
		values.Set("spotify", "connected")
		redirectURL.RawQuery = values.Encode()
		http.Redirect(w, r, redirectURL.String(), http.StatusFound)
	}))))

	// Add health check endpoint with CORS and logging
	http.Handle("/health", corsMiddleware(cfg.FrontendURL, loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HEALTH] Health check requested")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))))

	log.Println("[ROUTES] ✅ HTTP routes configured")

	// Set up graceful shutdown
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      nil,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Start server in a goroutine
	go func() {
		log.Printf("🚀 Server ready at http://localhost:%s/", cfg.Port)
		log.Printf("🕹  GraphQL playground at http://localhost:%s/", cfg.Port)
		log.Printf("💚 Health check at http://localhost:%s/health", cfg.Port)
		log.Printf("📊 Accepting requests from http://localhost:3000 (CORS enabled)")

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Failed to start server: %v", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	log.Println("🛑 Shutting down server...")

	// Give the server 30 seconds to shutdown gracefully
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Printf("[ERROR] Server forced to shutdown: %v", err)
	}

	log.Println("✅ Server exited")
}
