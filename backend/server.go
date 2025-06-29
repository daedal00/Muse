package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
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
	"github.com/golang-jwt/jwt/v5"
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
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		origin := r.Header.Get("Origin")
		log.Printf("[CORS] Request from origin: %s", origin)

		// Allow both HTTP and HTTPS 127.0.0.1 during development
		allowedOrigins := []string{
			"https://127.0.0.1:3000",
			"http://127.0.0.1:3000",
		}

		var allowOrigin string = ""
		for _, allowed := range allowedOrigins {
			if origin == allowed {
				allowOrigin = origin
				break
			}
		}

		if allowOrigin == "" {
			// Default to HTTPS if no origin header (for direct API access)
			allowOrigin = "https://127.0.0.1:3000"
		}

		// Set CORS headers - allow both HTTP and HTTPS 127.0.0.1
		w.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Allow-Credentials", "true")

		// Handle preflight OPTIONS request
		if r.Method == "OPTIONS" {
			log.Printf("[CORS] Handling preflight OPTIONS request for origin: %s", allowOrigin)
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
	http.Handle("/query", corsMiddleware(loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
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

	// Add health check endpoint with CORS and logging
	http.Handle("/health", corsMiddleware(loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[HEALTH] Health check requested")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	}))))

	// Add Spotify OAuth callback endpoint
	http.Handle("/auth/spotify/callback", corsMiddleware(loggingMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("[SPOTIFY] OAuth callback requested from %s", r.RemoteAddr)

		// Get code and state from query parameters
		code := r.URL.Query().Get("code")
		state := r.URL.Query().Get("state")
		errorParam := r.URL.Query().Get("error")
		errorDescription := r.URL.Query().Get("error_description")

		// Handle OAuth error cases as per Spotify documentation
		if errorParam != "" {
			log.Printf("[SPOTIFY] ❌ OAuth error: %s", errorParam)
			if errorDescription != "" {
				log.Printf("[SPOTIFY] Error description: %s", errorDescription)
			}

			// Map Spotify errors to user-friendly messages
			var userError string
			switch errorParam {
			case "access_denied":
				userError = "access_denied"
				log.Printf("[SPOTIFY] User denied access to Spotify")
			case "invalid_client":
				userError = "invalid_client"
				log.Printf("[SPOTIFY] Invalid client configuration")
			case "invalid_request":
				userError = "invalid_request"
				log.Printf("[SPOTIFY] Invalid request parameters")
			case "unauthorized_client":
				userError = "unauthorized_client"
				log.Printf("[SPOTIFY] Client not authorized for this grant type")
			case "unsupported_response_type":
				userError = "unsupported_response_type"
				log.Printf("[SPOTIFY] Unsupported response type")
			case "invalid_scope":
				userError = "invalid_scope"
				log.Printf("[SPOTIFY] Invalid or unsupported scope")
			default:
				userError = "unknown_error"
				log.Printf("[SPOTIFY] Unknown OAuth error: %s", errorParam)
			}

			redirectURL := fmt.Sprintf("https://127.0.0.1:3000/spotify-callback?error=%s", userError)
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}

		// Validate required parameters
		if code == "" {
			log.Printf("[SPOTIFY] ❌ Missing authorization code in callback")
			redirectURL := "https://127.0.0.1:3000/spotify-callback?error=missing_code"
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}

		if state == "" {
			log.Printf("[SPOTIFY] ❌ Missing state parameter in callback")
			redirectURL := "https://127.0.0.1:3000/spotify-callback?error=missing_state"
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
			return
		}

		// Validate state parameter format and expiry before processing
		if resolver.SpotifyAuthService != nil {
			// Pre-validate state to provide better error messages
			if _, err := resolver.SpotifyAuthService.ValidateState(state); err != nil {
				log.Printf("[SPOTIFY] ❌ State validation failed: %v", err)
				var errorCode string
				if strings.Contains(err.Error(), "expired") {
					errorCode = "state_expired"
				} else if strings.Contains(err.Error(), "format") {
					errorCode = "state_invalid"
				} else {
					errorCode = "state_error"
				}
				redirectURL := fmt.Sprintf("https://127.0.0.1:3000/spotify-callback?error=%s", errorCode)
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}

			ctx := r.Context()
			log.Printf("[SPOTIFY] Processing authorization code exchange...")

			user, err := resolver.SpotifyAuthService.HandleCallback(ctx, code, state)
			if err != nil {
				log.Printf("[SPOTIFY] ❌ Callback handling failed: %v", err)

				// Provide more specific error codes based on the error type
				var errorCode string
				switch {
				case strings.Contains(err.Error(), "exchange code"):
					errorCode = "token_exchange_failed"
				case strings.Contains(err.Error(), "validation failed"):
					errorCode = "token_validation_failed"
				case strings.Contains(err.Error(), "get Spotify user"):
					errorCode = "user_profile_failed"
				case strings.Contains(err.Error(), "update user"):
					errorCode = "user_update_failed"
				default:
					errorCode = "auth_failed"
				}

				redirectURL := fmt.Sprintf("https://127.0.0.1:3000/spotify-callback?error=%s", errorCode)
				http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
				return
			}

			log.Printf("[SPOTIFY] ✅ User authenticated with Spotify successfully")
			log.Printf("[SPOTIFY] User: %s (ID: %s)", user.Email, user.ID)

			// Redirect to frontend callback page with success
			redirectURL := "https://127.0.0.1:3000/spotify-callback?success=true"
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		} else {
			log.Printf("[SPOTIFY] ❌ Spotify auth service not available")
			redirectURL := "https://127.0.0.1:3000/spotify-callback?error=service_unavailable"
			http.Redirect(w, r, redirectURL, http.StatusTemporaryRedirect)
		}
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
		log.Printf("🚀 Server ready at https://127.0.0.1:%s/", cfg.Port)
		log.Printf("🕹  GraphQL playground at https://127.0.0.1:%s/", cfg.Port)
		log.Printf("💚 Health check at https://127.0.0.1:%s/health", cfg.Port)
		log.Printf("📊 Accepting requests from https://127.0.0.1:3000 (CORS enabled)")

		if err := server.ListenAndServeTLS("cert.pem", "key.pem"); err != nil && err != http.ErrServerClosed {
			log.Fatalf("[ERROR] Failed to start HTTPS server: %v", err)
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
