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
	repos            *repository.Repositories
	spotifyServices  *spotify.Services
	subscriptionMgr  *SubscriptionManager
	paginationHelper *PaginationHelper
	config           *config.Config
}

// NewResolver creates a new GraphQL resolver with all required dependencies
func NewResolver(cfg *config.Config) (*Resolver, error) {
	// Initialize PostgreSQL database
	postgresDB, err := database.NewPostgresConnection(cfg.DatabaseURL)
	if err != nil {
		return nil, err
	}

	// Initialize Redis client
	redisClient, err := database.NewRedisConnection(cfg.RedisURL)
	if err != nil {
		return nil, err
	}

	// Initialize repositories (using Redis for sessions, PostgreSQL for others)
	repos := &repository.Repositories{
		User:       postgres.NewUserRepository(postgresDB),
		Artist:     postgres.NewArtistRepository(postgresDB),
		Album:      postgres.NewAlbumRepository(postgresDB),
		Track:      postgres.NewTrackRepository(postgresDB),
		Review:     postgres.NewReviewRepository(postgresDB),
		Playlist:   postgres.NewPlaylistRepository(postgresDB),
		Session:    redisrepo.NewSessionRepository(redisClient),    // Using Redis for sessions
		MusicCache: redisrepo.NewMusicCacheRepository(redisClient), // Using Redis for music caching
	}

	// Initialize Spotify services (optional)
	var spotifyServices *spotify.Services
	if cfg.SpotifyClientID != "" && cfg.SpotifyClientSecret != "" {
		// Create a basic Spotify client for client credentials flow (public data)
		spotifyClient := spotify.NewClient(spotify.Config{
			ClientID:     cfg.SpotifyClientID,
			ClientSecret: cfg.SpotifyClientSecret,
			RedirectURL:  "http://localhost:8080/callback", // Default redirect for client credentials
			Scopes:       []string{},                       // No scopes needed for client credentials flow
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
	}

	// Initialize subscription manager
	subscriptionMgr := NewSubscriptionManager(redisClient)

	// Initialize pagination helper
	paginationHelper := NewPaginationHelper(repos)

	return &Resolver{
		repos:            repos,
		spotifyServices:  spotifyServices,
		subscriptionMgr:  subscriptionMgr,
		paginationHelper: paginationHelper,
		config:           cfg,
	}, nil
}

// Close properly closes all database connections
func (r *Resolver) Close() error {
	// Note: In a production system, you'd want to track both connections
	// and close them properly. For now, we'll add this placeholder.
	return nil
}
