package graph

//go:generate go run github.com/99designs/gqlgen generate

import (
	"context"
	"log"

	"github.com/daedal00/muse/backend/graph/model"
	"github.com/daedal00/muse/backend/internal/spotify"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{
	users []*model.User
	passwordHashes map[string]string

	albums []*model.Album
	reviews []*model.Review
	playlists []*model.Playlist
	
	// Spotify client and services for external API calls
	spotifyClient *spotify.Client
	spotifyServices *spotify.Services
}

func NewResolver(clientID, clientSecret string) *Resolver {
	// Create Spotify client with the new implementation
	spClient := spotify.NewClient(spotify.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/callback",
		Scopes: []string{
			// Add scopes as needed for your application
		},
	})

	// Get client credentials client for public data access
	ctx := context.Background()
	publicClient, err := spClient.GetClientCredentialsClient(ctx)
	if err != nil {
		log.Printf("Warning: Failed to initialize Spotify client: %v", err)
		// Continue without Spotify for now
	}

	var services *spotify.Services
	if publicClient != nil {
		services = spotify.NewServices(publicClient)
	}

	return &Resolver{
		users: make([]*model.User, 0),
		passwordHashes: make(map[string]string),
		albums: make([]*model.Album, 0),
		reviews: make([]*model.Review, 0),
		playlists: make([]*model.Playlist, 0),
		spotifyClient: spClient,
		spotifyServices: services,
	}
}