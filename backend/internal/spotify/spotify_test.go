package spotify

import (
	"context"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
	"golang.org/x/oauth2"
)

func TestNewClient(t *testing.T) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{spotifyauth.ScopeUserReadPrivate},
	}

	client := NewClient(config)

	assert.NotNil(t, client)
	assert.Equal(t, config.ClientID, client.clientID)
	assert.Equal(t, config.ClientSecret, client.clientSecret)
	assert.NotNil(t, client.auth)
}

func TestNewClientFromEnv(t *testing.T) {
	// Set environment variables
	os.Setenv("SPOTIFY_CLIENT_ID", "env-client-id")
	os.Setenv("SPOTIFY_CLIENT_SECRET", "env-client-secret")
	os.Setenv("SPOTIFY_REDIRECT_URL", "http://localhost:8080/auth/callback")

	defer func() {
		os.Unsetenv("SPOTIFY_CLIENT_ID")
		os.Unsetenv("SPOTIFY_CLIENT_SECRET")
		os.Unsetenv("SPOTIFY_REDIRECT_URL")
	}()

	client := NewClientFromEnv()

	assert.NotNil(t, client)
	assert.Equal(t, "env-client-id", client.clientID)
	assert.Equal(t, "env-client-secret", client.clientSecret)
}

func TestClient_GetAuthURL(t *testing.T) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{spotifyauth.ScopeUserReadPrivate},
	}

	client := NewClient(config)
	state := "test-state"

	authURL := client.GetAuthURL(state)

	assert.NotEmpty(t, authURL)
	assert.Contains(t, authURL, "accounts.spotify.com/authorize")
	assert.Contains(t, authURL, "state=test-state")
	assert.Contains(t, authURL, "response_type=code")
}

func TestClient_GetAuthURLWithPKCE(t *testing.T) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{spotifyauth.ScopeUserReadPrivate},
	}

	client := NewClient(config)
	state := "test-state"
	codeChallenge := "test-code-challenge"

	authURL := client.GetAuthURLWithPKCE(state, codeChallenge)

	assert.NotEmpty(t, authURL)
	assert.Contains(t, authURL, "accounts.spotify.com/authorize")
	assert.Contains(t, authURL, "client_id=test-client-id")
	assert.Contains(t, authURL, "state=test-state")
	assert.Contains(t, authURL, "code_challenge=test-code-challenge")
	assert.Contains(t, authURL, "code_challenge_method=S256")
}

func TestClient_GetClientCredentialsClient(t *testing.T) {
	// Skip this test if we don't have real credentials
	clientID := os.Getenv("SPOTIFY_CLIENT_ID")
	clientSecret := os.Getenv("SPOTIFY_CLIENT_SECRET")

	if clientID == "" || clientSecret == "" {
		t.Skip("Skipping integration test: SPOTIFY_CLIENT_ID and SPOTIFY_CLIENT_SECRET not set")
	}

	config := Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		RedirectURL:  "http://localhost:8080/callback",
	}

	client := NewClient(config)
	ctx := context.Background()

	spotifyClient, err := client.GetClientCredentialsClient(ctx)

	require.NoError(t, err)
	assert.NotNil(t, spotifyClient)
}

func TestClient_GetAuthorizedClient(t *testing.T) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{spotifyauth.ScopeUserReadPrivate},
	}

	client := NewClient(config)
	ctx := context.Background()

	// Mock token
	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
	}

	spotifyClient := client.GetAuthorizedClient(ctx, token)

	assert.NotNil(t, spotifyClient)
}

func TestNewServices(t *testing.T) {
	// Create a mock spotify client
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
	}

	client := NewClient(config)
	ctx := context.Background()

	// Mock token
	token := &oauth2.Token{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
	}

	spotifyClient := client.GetAuthorizedClient(ctx, token)
	services := NewServices(spotifyClient)

	assert.NotNil(t, services)
	assert.NotNil(t, services.Search)
	assert.NotNil(t, services.User)
	assert.NotNil(t, services.Playlist)
	assert.NotNil(t, services.Track)
	assert.NotNil(t, services.Album)
	assert.NotNil(t, services.Artist)
	assert.Equal(t, spotifyClient, services.GetSpotifyClient())
}

func TestNewPaginationHelper(t *testing.T) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
	}

	client := NewClient(config)
	ctx := context.Background()

	token := &oauth2.Token{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
	}

	spotifyClient := client.GetAuthorizedClient(ctx, token)
	helper := NewPaginationHelper(spotifyClient)

	assert.NotNil(t, helper)
	assert.NotNil(t, helper.client)
}

// Benchmark tests for performance monitoring
func BenchmarkClient_GetAuthURL(b *testing.B) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{spotifyauth.ScopeUserReadPrivate},
	}

	client := NewClient(config)
	state := "benchmark-state"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.GetAuthURL(state)
	}
}

func BenchmarkClient_GetAuthURLWithPKCE(b *testing.B) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{spotifyauth.ScopeUserReadPrivate},
	}

	client := NewClient(config)
	state := "benchmark-state"
	codeChallenge := "benchmark-code-challenge"

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.GetAuthURLWithPKCE(state, codeChallenge)
	}
}

func BenchmarkNewServices(b *testing.B) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
	}

	client := NewClient(config)
	ctx := context.Background()

	token := &oauth2.Token{
		AccessToken: "test-access-token",
		TokenType:   "Bearer",
	}

	spotifyClient := client.GetAuthorizedClient(ctx, token)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = NewServices(spotifyClient)
	}
}

func BenchmarkClient_GetAuthorizedClient(b *testing.B) {
	config := Config{
		ClientID:     "test-client-id",
		ClientSecret: "test-client-secret",
		RedirectURL:  "http://localhost:8080/callback",
		Scopes:       []string{spotifyauth.ScopeUserReadPrivate},
	}

	client := NewClient(config)
	ctx := context.Background()

	token := &oauth2.Token{
		AccessToken:  "test-access-token",
		TokenType:    "Bearer",
		RefreshToken: "test-refresh-token",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_ = client.GetAuthorizedClient(ctx, token)
	}
}
