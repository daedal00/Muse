package graph

import (
	"net/url"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	spotifyauth "github.com/zmb3/spotify/v2/auth"
)

// SpotifyStateClaims stores auth flow state for callback validation.
type SpotifyStateClaims struct {
	UserID      string `json:"user_id"`
	RedirectURI string `json:"redirect_uri,omitempty"`
	jwt.RegisteredClaims
}

func spotifyUserScopes() []string {
	return []string{
		spotifyauth.ScopeUserTopRead,
		spotifyauth.ScopeUserLibraryRead,
		spotifyauth.ScopePlaylistReadPrivate,
		spotifyauth.ScopePlaylistReadCollaborative,
		spotifyauth.ScopeUserReadRecentlyPlayed,
		spotifyauth.ScopeUserReadPrivate,
	}
}

// SpotifyUserScopes exposes the configured Spotify OAuth scopes.
func SpotifyUserScopes() []string {
	return spotifyUserScopes()
}

func (r *Resolver) sanitizeRedirectURI(raw string) string {
	if raw == "" {
		return ""
	}

	parsed, err := url.Parse(raw)
	if err != nil {
		return ""
	}

	base, err := url.Parse(r.config.FrontendURL)
	if err != nil || base.Host == "" {
		return ""
	}

	if !strings.EqualFold(parsed.Host, base.Host) {
		return ""
	}

	if parsed.Scheme != "http" && parsed.Scheme != "https" {
		return ""
	}

	return parsed.String()
}

func spotifyStateTTL() time.Duration {
	return 10 * time.Minute
}
