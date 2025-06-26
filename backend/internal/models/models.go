package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID `json:"id" db:"id"`
	Name         string    `json:"name" db:"name"`
	Email        string    `json:"email" db:"email"`
	PasswordHash string    `json:"-" db:"password_hash"`
	Bio          *string   `json:"bio" db:"bio"`
	Avatar       *string   `json:"avatar" db:"avatar"`

	// Spotify OAuth fields
	SpotifyID           *string    `json:"spotify_id" db:"spotify_id"`
	SpotifyAccessToken  *string    `json:"-" db:"spotify_access_token"`
	SpotifyRefreshToken *string    `json:"-" db:"spotify_refresh_token"`
	SpotifyTokenExpiry  *time.Time `json:"-" db:"spotify_token_expiry"`

	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// Review represents a user's review of a Spotify item (album or track)
// We only store the review data, not the music metadata
type Review struct {
	ID          uuid.UUID `json:"id" db:"id"`
	UserID      uuid.UUID `json:"user_id" db:"user_id"`
	SpotifyID   string    `json:"spotify_id" db:"spotify_id"`     // Spotify ID of the item being reviewed
	SpotifyType string    `json:"spotify_type" db:"spotify_type"` // "album" or "track"
	Rating      int       `json:"rating" db:"rating"`
	ReviewText  *string   `json:"review_text" db:"review_text"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Relations
	User *User `json:"user,omitempty"`
}

// Playlist represents a user's playlist
// We store minimal metadata and track references by Spotify ID
type Playlist struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description" db:"description"`
	CoverImage  *string   `json:"cover_image" db:"cover_image"`
	CreatorID   uuid.UUID `json:"creator_id" db:"creator_id"`
	IsPublic    bool      `json:"is_public" db:"is_public"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`

	// Relations
	Creator *User `json:"creator,omitempty"`
}

// PlaylistTrack represents tracks in a playlist (only store Spotify IDs)
type PlaylistTrack struct {
	ID            uuid.UUID `json:"id" db:"id"`
	PlaylistID    uuid.UUID `json:"playlist_id" db:"playlist_id"`
	SpotifyID     string    `json:"spotify_id" db:"spotify_id"` // Spotify track ID
	Position      int       `json:"position" db:"position"`
	AddedAt       time.Time `json:"added_at" db:"added_at"`
	AddedByUserID uuid.UUID `json:"added_by_user_id" db:"added_by_user_id"`

	// Relations
	Playlist    *Playlist `json:"playlist,omitempty"`
	AddedByUser *User     `json:"added_by_user,omitempty"`
}

// Session represents a user session
type Session struct {
	ID        string    `json:"id" db:"id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	ExpiresAt time.Time `json:"expires_at" db:"expires_at"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`

	// Relations
	User *User `json:"user,omitempty"`
}

// UserPreferences represents user's music preferences and settings
type UserPreferences struct {
	ID                   uuid.UUID       `json:"id" db:"id"`
	UserID               uuid.UUID       `json:"user_id" db:"user_id"`
	PreferredGenres      []string        `json:"preferred_genres" db:"preferred_genres"`           // JSON array
	FavoriteArtistIDs    []string        `json:"favorite_artist_ids" db:"favorite_artist_ids"`     // Spotify IDs
	NotificationSettings map[string]bool `json:"notification_settings" db:"notification_settings"` // JSON object
	PrivacySettings      map[string]bool `json:"privacy_settings" db:"privacy_settings"`           // JSON object
	CreatedAt            time.Time       `json:"created_at" db:"created_at"`
	UpdatedAt            time.Time       `json:"updated_at" db:"updated_at"`

	// Relations
	User *User `json:"user,omitempty"`
}

// Pagination types
type PageInfo struct {
	EndCursor   *string `json:"end_cursor"`
	HasNextPage bool    `json:"has_next_page"`
}

type Connection[T any] struct {
	TotalCount int       `json:"total_count"`
	Edges      []Edge[T] `json:"edges"`
	PageInfo   PageInfo  `json:"page_info"`
}

type Edge[T any] struct {
	Cursor string `json:"cursor"`
	Node   T      `json:"node"`
}

// Spotify Data Models (for caching and API responses)
// These are not stored in the database, only used for API responses and caching

type SpotifyArtist struct {
	ID        string   `json:"id"`
	Name      string   `json:"name"`
	Images    []string `json:"images"`
	Genres    []string `json:"genres"`
	Followers int      `json:"followers"`
}

type SpotifyAlbum struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Artists     []SpotifyArtist `json:"artists"`
	Images      []string        `json:"images"`
	ReleaseDate string          `json:"release_date"`
	TotalTracks int             `json:"total_tracks"`
	TrackIds    []string        `json:"track_ids"` // Cached track IDs for performance
}

type SpotifyTrack struct {
	ID           string            `json:"id"`
	Name         string            `json:"name"`
	Artists      []SpotifyArtist   `json:"artists"`
	Album        SpotifyAlbum      `json:"album"`
	DurationMS   int               `json:"duration_ms"`
	TrackNumber  int               `json:"track_number"`
	PreviewURL   *string           `json:"preview_url"`
	ExternalURLs map[string]string `json:"external_urls"`
}

type SpotifyPlaylist struct {
	ID          string   `json:"id"`
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Images      []string `json:"images"`
	TrackCount  int      `json:"track_count"`
	IsPublic    bool     `json:"is_public"`
	OwnerID     string   `json:"owner_id"`
}

// Cache-specific data structures
type CachedUserData struct {
	UserID         uuid.UUID       `json:"user_id"`
	RecentlyPlayed []SpotifyTrack  `json:"recently_played"`
	TopTracks      []SpotifyTrack  `json:"top_tracks"`
	TopArtists     []SpotifyArtist `json:"top_artists"`
	SavedTracks    []string        `json:"saved_tracks"` // Spotify IDs
	SavedAlbums    []string        `json:"saved_albums"` // Spotify IDs
	PlaylistIDs    []string        `json:"playlist_ids"` // Spotify IDs
	LastUpdated    time.Time       `json:"last_updated"`
}

type CachedRecommendations struct {
	UserID             uuid.UUID       `json:"user_id"`
	RecommendedTracks  []SpotifyTrack  `json:"recommended_tracks"`
	RecommendedAlbums  []SpotifyAlbum  `json:"recommended_albums"`
	RecommendedArtists []SpotifyArtist `json:"recommended_artists"`
	BasedOnGenres      []string        `json:"based_on_genres"`
	BasedOnArtists     []string        `json:"based_on_artists"`
	GeneratedAt        time.Time       `json:"generated_at"`
}
