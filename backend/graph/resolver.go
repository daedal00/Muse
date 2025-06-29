package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"log"

	"github.com/daedal00/muse/backend/internal/config"
	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/repository"
	"github.com/daedal00/muse/backend/internal/repository/postgres"
	redisrepo "github.com/daedal00/muse/backend/internal/repository/redis"
	"github.com/daedal00/muse/backend/internal/spotify"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

// Resolver is the root resolver struct
type Resolver struct {
	repos                  *repository.Repositories
	spotifyServices        *spotify.Services
	SpotifyAuthService     *spotify.AuthService
	SpotifyPlaylistService *spotify.PlaylistImportService
	subscriptionMgr        *SubscriptionManager
	paginationHelper       *PaginationHelper
	config                 *config.Config
}

// NewResolver creates a new GraphQL resolver with all required dependencies
func NewResolver(cfg *config.Config) (*Resolver, error) {
	// Initialize PostgreSQL database
	postgresDB, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}
	log.Printf("✅ Connected to PostgreSQL database")

	// Initialize Redis client
	redisClient, err := database.NewRedisConnection(cfg.RedisURL)
	if err != nil {
		return nil, err
	}
	log.Printf("✅ Connected to Redis at %s", cfg.RedisURL)

	// Initialize repositories (using Redis for sessions, PostgreSQL for others)
	repos := &repository.Repositories{
		User:            postgres.NewUserRepository(postgresDB),
		Review:          postgres.NewReviewRepository(postgresDB),
		Playlist:        postgres.NewPlaylistRepository(postgresDB),
		UserPreferences: postgres.NewUserPreferencesRepository(postgresDB), // New repository for user preferences
		Session:         redisrepo.NewSessionRepository(redisClient),       // Using Redis for sessions
		SpotifyCache:    redisrepo.NewSpotifyCacheRepository(redisClient),  // New optimized Spotify cache
		MusicCache:      redisrepo.NewMusicCacheRepository(redisClient),    // Legacy music cache (backward compatibility)
	}

	// Initialize Spotify services (optional)
	var spotifyServices *spotify.Services
	var spotifyAuthService *spotify.AuthService
	var spotifyPlaylistService *spotify.PlaylistImportService

	if cfg.SpotifyClientID != "" && cfg.SpotifyClientSecret != "" {
		// Create a basic Spotify client for client credentials flow (public data)
		spotifyClient := spotify.NewClient(spotify.Config{
			ClientID:     cfg.SpotifyClientID,
			ClientSecret: cfg.SpotifyClientSecret,
			RedirectURL:  cfg.SpotifyRedirectURL, // Use config redirect URL
			Scopes: []string{
				"user-read-private",
				"user-read-email",
				"playlist-read-private",
				"playlist-read-collaborative",
				"user-library-read",
			},
		})

		// Get client credentials client for public API access
		ctx := context.Background()
		client, err := spotifyClient.GetClientCredentialsClient(ctx)
		if err != nil {
			// Log error but don't fail - Spotify is optional
			log.Printf("Warning: Failed to initialize Spotify client: %v", err)
		} else {
			spotifyServices = spotify.NewServices(client)
		}

		// Initialize auth and playlist services
		spotifyAuthService = spotify.NewAuthService(spotifyClient, postgresDB.Pool)
		spotifyPlaylistService = spotify.NewPlaylistImportService(postgresDB.Pool, spotifyAuthService)
	}

	// Initialize subscription manager
	subscriptionMgr := NewSubscriptionManager(redisClient)

	// Initialize pagination helper
	paginationHelper := NewPaginationHelper(repos)

	return &Resolver{
		repos:                  repos,
		spotifyServices:        spotifyServices,
		SpotifyAuthService:     spotifyAuthService,
		SpotifyPlaylistService: spotifyPlaylistService,
		subscriptionMgr:        subscriptionMgr,
		paginationHelper:       paginationHelper,
		config:                 cfg,
	}, nil
}

// Close properly closes all database connections
func (r *Resolver) Close() error {
	// Note: In a production system, you'd want to track both connections
	// and close them properly. For now, we'll add this placeholder.
	return nil
}
