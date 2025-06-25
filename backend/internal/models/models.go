package models

import (
	"time"

	"github.com/google/uuid"
)

// User represents a user in the system
type User struct {
	ID           uuid.UUID  `json:"id" db:"id"`
	Name         string     `json:"name" db:"name"`
	Email        string     `json:"email" db:"email"`
	PasswordHash string     `json:"-" db:"password_hash"`
	Bio          *string    `json:"bio" db:"bio"`
	Avatar       *string    `json:"avatar" db:"avatar"`
	CreatedAt    time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at" db:"updated_at"`
}

// Artist represents a music artist
type Artist struct {
	ID        uuid.UUID  `json:"id" db:"id"`
	SpotifyID *string    `json:"spotify_id" db:"spotify_id"`
	Name      string     `json:"name" db:"name"`
	CreatedAt time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt time.Time  `json:"updated_at" db:"updated_at"`
}

// Album represents a music album
type Album struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	SpotifyID   *string    `json:"spotify_id" db:"spotify_id"`
	Title       string     `json:"title" db:"title"`
	ArtistID    uuid.UUID  `json:"artist_id" db:"artist_id"`
	ReleaseDate *time.Time `json:"release_date" db:"release_date"`
	CoverImage  *string    `json:"cover_image" db:"cover_image"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relations
	Artist *Artist `json:"artist,omitempty"`
}

// Track represents a music track
type Track struct {
	ID          uuid.UUID  `json:"id" db:"id"`
	SpotifyID   *string    `json:"spotify_id" db:"spotify_id"`
	Title       string     `json:"title" db:"title"`
	AlbumID     uuid.UUID  `json:"album_id" db:"album_id"`
	DurationMs  *int       `json:"duration_ms" db:"duration_ms"`
	TrackNumber *int       `json:"track_number" db:"track_number"`
	CreatedAt   time.Time  `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at" db:"updated_at"`
	
	// Relations
	Album *Album `json:"album,omitempty"`
}

// Review represents a user's review of an album
type Review struct {
	ID         uuid.UUID `json:"id" db:"id"`
	UserID     uuid.UUID `json:"user_id" db:"user_id"`
	AlbumID    uuid.UUID `json:"album_id" db:"album_id"`
	Rating     int       `json:"rating" db:"rating"`
	ReviewText *string   `json:"review_text" db:"review_text"`
	CreatedAt  time.Time `json:"created_at" db:"created_at"`
	UpdatedAt  time.Time `json:"updated_at" db:"updated_at"`
	
	// Relations
	User  *User  `json:"user,omitempty"`
	Album *Album `json:"album,omitempty"`
}

// Playlist represents a user's playlist
type Playlist struct {
	ID          uuid.UUID `json:"id" db:"id"`
	Title       string    `json:"title" db:"title"`
	Description *string   `json:"description" db:"description"`
	CoverImage  *string   `json:"cover_image" db:"cover_image"`
	CreatorID   uuid.UUID `json:"creator_id" db:"creator_id"`
	CreatedAt   time.Time `json:"created_at" db:"created_at"`
	UpdatedAt   time.Time `json:"updated_at" db:"updated_at"`
	
	// Relations
	Creator *User `json:"creator,omitempty"`
}

// PlaylistTrack represents the many-to-many relationship between playlists and tracks
type PlaylistTrack struct {
	ID         uuid.UUID `json:"id" db:"id"`
	PlaylistID uuid.UUID `json:"playlist_id" db:"playlist_id"`
	TrackID    uuid.UUID `json:"track_id" db:"track_id"`
	Position   int       `json:"position" db:"position"`
	AddedAt    time.Time `json:"added_at" db:"added_at"`
	
	// Relations
	Playlist *Playlist `json:"playlist,omitempty"`
	Track    *Track    `json:"track,omitempty"`
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

// Pagination types
type PageInfo struct {
	EndCursor   *string `json:"end_cursor"`
	HasNextPage bool    `json:"has_next_page"`
}

type Connection[T any] struct {
	TotalCount int         `json:"total_count"`
	Edges      []Edge[T]   `json:"edges"`
	PageInfo   PageInfo    `json:"page_info"`
}

type Edge[T any] struct {
	Cursor string `json:"cursor"`
	Node   T      `json:"node"`
} 