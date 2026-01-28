package graph

import (
	"context"
	"fmt"
	"time"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/daedal00/muse/backend/internal/spotify"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

func (r *Resolver) currentUserUUID(ctx context.Context) (uuid.UUID, error) {
	raw, ok := ForContext(ctx)
	if !ok {
		return uuid.UUID{}, fmt.Errorf("unauthenticated")
	}

	userID, err := uuid.Parse(raw)
	if err != nil {
		return uuid.UUID{}, fmt.Errorf("invalid user ID")
	}

	return userID, nil
}

func (r *Resolver) spotifyUserServices(ctx context.Context) (*spotify.Services, *models.SpotifyToken, error) {
	if r.config.SpotifyClientID == "" || r.config.SpotifyClientSecret == "" {
		return nil, nil, fmt.Errorf("spotify is not configured")
	}

	userID, err := r.currentUserUUID(ctx)
	if err != nil {
		return nil, nil, err
	}

	token, err := r.repos.Spotify.GetByUserID(ctx, userID)
	if err != nil {
		return nil, nil, fmt.Errorf("spotify not connected")
	}

	oauthToken := &oauth2.Token{
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		TokenType:    derefString(token.TokenType, "Bearer"),
		Expiry:       token.ExpiresAt,
	}

	spotifyClient := spotify.NewClient(spotify.Config{
		ClientID:     r.config.SpotifyClientID,
		ClientSecret: r.config.SpotifyClientSecret,
		RedirectURL:  r.config.SpotifyRedirectURL,
		Scopes:       spotifyUserScopes(),
	})

	client, refreshedToken, err := spotifyClient.GetUserClient(ctx, oauthToken)
	if err != nil {
		return nil, nil, err
	}

	if refreshedToken != nil {
		updated := token
		updated.AccessToken = refreshedToken.AccessToken
		if refreshedToken.RefreshToken != "" {
			updated.RefreshToken = refreshedToken.RefreshToken
		}
		if refreshedToken.TokenType != "" {
			updated.TokenType = &refreshedToken.TokenType
		}
		if refreshedToken.Expiry.After(time.Time{}) {
			updated.ExpiresAt = refreshedToken.Expiry
		}
		updated.UpdatedAt = time.Now()

		if err := r.repos.Spotify.Upsert(ctx, updated); err != nil {
			return nil, nil, fmt.Errorf("failed to update spotify token: %w", err)
		}
	}

	return spotify.NewServices(client), token, nil
}

func derefString(val *string, fallback string) string {
	if val == nil || *val == "" {
		return fallback
	}
	return *val
}
