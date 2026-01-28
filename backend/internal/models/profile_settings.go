package models

import (
	"time"

	"github.com/google/uuid"
)

// ProfileLayout represents the layout style for a user's profile
type ProfileLayout string

const (
	ProfileLayoutGrid  ProfileLayout = "GRID"
	ProfileLayoutList  ProfileLayout = "LIST"
	ProfileLayoutBento ProfileLayout = "BENTO"
)

// BackgroundStyle represents the background type for a user's profile
type BackgroundStyle string

const (
	BackgroundStyleSolid    BackgroundStyle = "SOLID"
	BackgroundStyleGradient BackgroundStyle = "GRADIENT"
	BackgroundStyleImage    BackgroundStyle = "IMAGE"
)

// BioStyle represents the display style for a user's bio
type BioStyle string

const (
	BioStyleMinimal  BioStyle = "MINIMAL"
	BioStyleDetailed BioStyle = "DETAILED"
	BioStyleQuote    BioStyle = "QUOTE"
)

// ProfileSettings represents a user's profile customization settings
type ProfileSettings struct {
	UserID uuid.UUID `json:"user_id" db:"user_id"`

	// Layout configuration
	Layout ProfileLayout `json:"layout" db:"layout"`

	// Theme colors
	PrimaryColor *string `json:"primary_color" db:"primary_color"`
	AccentColor  *string `json:"accent_color" db:"accent_color"`

	// Background customization
	BackgroundStyle *BackgroundStyle `json:"background_style" db:"background_style"`
	BackgroundValue *string          `json:"background_value" db:"background_value"`

	// Pinned content
	PinnedAlbumIDs    []uuid.UUID `json:"pinned_album_ids" db:"pinned_album_ids"`
	PinnedTrackIDs    []uuid.UUID `json:"pinned_track_ids" db:"pinned_track_ids"`
	FeaturedArtistIDs []uuid.UUID `json:"featured_artist_ids" db:"featured_artist_ids"`

	// Section ordering
	SectionsOrder []string `json:"sections_order" db:"sections_order"`

	// Visibility toggles
	ShowSpotifyStats     bool `json:"show_spotify_stats" db:"show_spotify_stats"`
	ShowListeningHistory bool `json:"show_listening_history" db:"show_listening_history"`

	// Bio style
	BioStyle *BioStyle `json:"bio_style" db:"bio_style"`

	// Custom tags (genres, vibes)
	CustomTags []string `json:"custom_tags" db:"custom_tags"`

	// Timestamps
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// DefaultProfileSettings returns a ProfileSettings with default values
func DefaultProfileSettings(userID uuid.UUID) *ProfileSettings {
	primaryColor := "#1DB954"
	accentColor := "#191414"
	bgStyle := BackgroundStyleSolid
	bgValue := "#0f172a"
	bioStyle := BioStyleMinimal

	return &ProfileSettings{
		UserID:               userID,
		Layout:               ProfileLayoutGrid,
		PrimaryColor:         &primaryColor,
		AccentColor:          &accentColor,
		BackgroundStyle:      &bgStyle,
		BackgroundValue:      &bgValue,
		PinnedAlbumIDs:       []uuid.UUID{},
		PinnedTrackIDs:       []uuid.UUID{},
		FeaturedArtistIDs:    []uuid.UUID{},
		SectionsOrder:        []string{"featuredArtists", "topTracks", "pinnedAlbums", "recentReviews", "playlists"},
		ShowSpotifyStats:     true,
		ShowListeningHistory: true,
		BioStyle:             &bioStyle,
		CustomTags:           []string{},
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
}
