package spotify

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/zmb3/spotify/v2"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"

	"github.com/daedal00/muse/backend/internal/models"
)

// AuthService handles Spotify OAuth operations
type AuthService struct {
	client *Client
	db     *pgxpool.Pool
}

// NewAuthService creates a new Spotify auth service
func NewAuthService(client *Client, db *pgxpool.Pool) *AuthService {
	return &AuthService{
		client: client,
		db:     db,
	}
}

// StateInfo holds information stored in the OAuth state parameter
type StateInfo struct {
	UserID    uuid.UUID
	Timestamp time.Time
	Random    string
}

// GenerateAuthURL generates a Spotify OAuth URL with state parameter
func (s *AuthService) GenerateAuthURL(userID uuid.UUID) (string, string, error) {
	// Generate a cryptographically secure random state parameter
	stateBytes := make([]byte, 32)
	if _, err := rand.Read(stateBytes); err != nil {
		return "", "", fmt.Errorf("failed to generate state: %w", err)
	}
	randomState := base64.URLEncoding.EncodeToString(stateBytes)

	// Create state with timestamp for expiry validation (prevents replay attacks)
	timestamp := time.Now().Unix()
	stateWithUser := fmt.Sprintf("%s:%d:%s", userID.String(), timestamp, randomState)

	authURL := s.client.GetAuthURL(stateWithUser)
	return authURL, stateWithUser, nil
}

// ValidateState validates the state parameter and extracts user information
func (s *AuthService) ValidateState(state string) (*StateInfo, error) {
	// State format: "userID:timestamp:randomState"
	parts := strings.Split(state, ":")
	if len(parts) != 3 {
		return nil, fmt.Errorf("invalid state format")
	}

	// Parse user ID
	userID, err := uuid.Parse(parts[0])
	if err != nil {
		return nil, fmt.Errorf("invalid user ID in state: %w", err)
	}

	// Parse timestamp
	var timestamp int64
	if _, err := fmt.Sscanf(parts[1], "%d", &timestamp); err != nil {
		return nil, fmt.Errorf("invalid timestamp in state: %w", err)
	}

	// Check if state is expired (15 minutes max age)
	stateTime := time.Unix(timestamp, 0)
	if time.Since(stateTime) > 15*time.Minute {
		return nil, fmt.Errorf("state parameter has expired")
	}

	return &StateInfo{
		UserID:    userID,
		Timestamp: stateTime,
		Random:    parts[2],
	}, nil
}

// HandleCallback processes the OAuth callback and stores tokens
func (s *AuthService) HandleCallback(ctx context.Context, code, state string) (*models.User, error) {
	// Validate state parameter
	stateInfo, err := s.ValidateState(state)
	if err != nil {
		return nil, fmt.Errorf("invalid state parameter: %w", err)
	}

	// Exchange code for token
	token, err := s.client.ExchangeCode(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code: %w", err)
	}

	// Validate token has required scopes
	if err := s.validateTokenScopes(token); err != nil {
		return nil, fmt.Errorf("token validation failed: %w", err)
	}

	// Get user's Spotify profile
	spotifyClient := s.client.GetAuthorizedClient(ctx, token)
	spotifyUser, err := spotifyClient.CurrentUser(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get Spotify user: %w", err)
	}

	// Update user with Spotify credentials
	user, err := s.updateUserWithSpotifyCredentials(ctx, stateInfo.UserID, spotifyUser, token)
	if err != nil {
		return nil, fmt.Errorf("failed to update user: %w", err)
	}

	return user, nil
}

// validateTokenScopes validates that the received token has the required scopes
func (s *AuthService) validateTokenScopes(token *oauth2.Token) error {
	// The Go spotify library doesn't directly expose scopes in the token,
	// but we can validate by attempting to access protected resources
	// This is a basic validation - in production you might want more comprehensive scope checking

	requiredScopes := []string{
		spotifyauth.ScopeUserReadPrivate,
		spotifyauth.ScopeUserReadEmail,
		spotifyauth.ScopePlaylistReadPrivate,
	}

	// Log the scopes for debugging (optional)
	if scopeStr, ok := token.Extra("scope").(string); ok {
		scopes := strings.Fields(scopeStr)
		fmt.Printf("Received scopes: %v\n", scopes)

		// Check if all required scopes are present
		scopeMap := make(map[string]bool)
		for _, scope := range scopes {
			scopeMap[scope] = true
		}

		for _, required := range requiredScopes {
			if !scopeMap[required] {
				return fmt.Errorf("missing required scope: %s", required)
			}
		}
	}

	return nil
}

// GetUserSpotifyClient returns an authorized Spotify client for a user
func (s *AuthService) GetUserSpotifyClient(ctx context.Context, userID uuid.UUID) (*spotify.Client, error) {
	user, err := s.getUserByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("failed to get user: %w", err)
	}

	if user.SpotifyAccessToken == nil {
		return nil, fmt.Errorf("user has not connected Spotify")
	}

	// Check if token needs refresh
	token := &oauth2.Token{
		AccessToken:  *user.SpotifyAccessToken,
		RefreshToken: *user.SpotifyRefreshToken,
		Expiry:       *user.SpotifyTokenExpiry,
	}

	// If token is expired, refresh it
	if token.Expiry.Before(time.Now()) {
		token, err = s.client.auth.RefreshToken(ctx, token)
		if err != nil {
			return nil, fmt.Errorf("failed to refresh token: %w", err)
		}

		// Update user with new token
		err = s.updateUserTokens(ctx, userID, token)
		if err != nil {
			return nil, fmt.Errorf("failed to update tokens: %w", err)
		}
	}

	return s.client.GetAuthorizedClient(ctx, token), nil
}

// Helper functions

func (s *AuthService) getUserByID(ctx context.Context, userID uuid.UUID) (*models.User, error) {
	query := `
		SELECT id, name, email, password_hash, bio, avatar, 
		       spotify_id, spotify_access_token, spotify_refresh_token, spotify_token_expiry,
		       created_at, updated_at
		FROM users WHERE id = $1
	`

	var user models.User
	row := s.db.QueryRow(ctx, query, userID)

	err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Bio, &user.Avatar,
		&user.SpotifyID, &user.SpotifyAccessToken, &user.SpotifyRefreshToken, &user.SpotifyTokenExpiry,
		&user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) updateUserWithSpotifyCredentials(ctx context.Context, userID uuid.UUID, spotifyUser *spotify.PrivateUser, token *oauth2.Token) (*models.User, error) {
	query := `
		UPDATE users 
		SET spotify_id = $2, spotify_access_token = $3, spotify_refresh_token = $4, spotify_token_expiry = $5, updated_at = NOW()
		WHERE id = $1
		RETURNING id, name, email, password_hash, bio, avatar, spotify_id, created_at, updated_at
	`

	var user models.User
	row := s.db.QueryRow(ctx, query, userID, spotifyUser.ID, token.AccessToken, token.RefreshToken, token.Expiry)

	err := row.Scan(
		&user.ID, &user.Name, &user.Email, &user.PasswordHash, &user.Bio, &user.Avatar,
		&user.SpotifyID, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (s *AuthService) updateUserTokens(ctx context.Context, userID uuid.UUID, token *oauth2.Token) error {
	query := `
		UPDATE users 
		SET spotify_access_token = $2, spotify_refresh_token = $3, spotify_token_expiry = $4, updated_at = NOW()
		WHERE id = $1
	`

	_, err := s.db.Exec(ctx, query, userID, token.AccessToken, token.RefreshToken, token.Expiry)
	return err
}
