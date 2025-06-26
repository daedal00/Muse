package graph

import (
	"context"
	"encoding/base64"
	"fmt"
	"time"

	"github.com/daedal00/muse/backend/internal/models"
	"github.com/daedal00/muse/backend/internal/repository"
)

// CursorInfo contains decoded cursor information
type CursorInfo struct {
	ID        string    `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	Position  int       `json:"position"`
}

// PaginationHelper provides improved cursor-based pagination
type PaginationHelper struct {
	repos *repository.Repositories
}

// NewPaginationHelper creates a new pagination helper
func NewPaginationHelper(repos *repository.Repositories) *PaginationHelper {
	return &PaginationHelper{repos: repos}
}

// EncodeCursor creates a cursor from an ID and created timestamp
func (p *PaginationHelper) EncodeCursor(id string, createdAt time.Time, position int) string {
	cursorData := fmt.Sprintf("%s:%d:%d", id, createdAt.Unix(), position)
	return base64.StdEncoding.EncodeToString([]byte(cursorData))
}

// DecodeCursor decodes a cursor to extract ID, timestamp, and position
func (p *PaginationHelper) DecodeCursor(cursor string) (*CursorInfo, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor format")
	}

	parts := string(data)
	// Try to parse as new format: id:timestamp:position
	var id string
	var timestamp int64
	var position int

	n, err := fmt.Sscanf(parts, "%36s:%d:%d", &id, &timestamp, &position)
	if err != nil || n != 3 {
		// Fallback to old format (just ID)
		id = parts
		timestamp = 0
		position = 0
	}

	createdAt := time.Unix(timestamp, 0)

	return &CursorInfo{
		ID:        id,
		CreatedAt: createdAt,
		Position:  position,
	}, nil
}

// GetAlbumsWithCursor fetches albums using proper cursor-based pagination
func (p *PaginationHelper) GetAlbumsWithCursor(ctx context.Context, first int, after *string) ([]*models.Album, bool, error) {
	var albums []*models.Album
	var err error

	if after != nil && *after != "" {
		cursor, err := p.DecodeCursor(*after)
		if err != nil {
			return nil, false, err
		}

		// If we have a proper cursor with timestamp, use it for efficient pagination
		if !cursor.CreatedAt.IsZero() {
			albums, err = p.getAlbumsAfterCursor(ctx, cursor, first+1)
			if err != nil {
				return nil, false, err
			}
		} else {
			// Fallback to position-based pagination
			albums, err = p.getAlbumsAfterPosition(ctx, cursor.ID, first+1)
			if err != nil {
				return nil, false, err
			}
		}
	} else {
		// No cursor, get from beginning
		albums, err = p.repos.Album.List(ctx, first+1, 0)
	}

	if err != nil {
		return nil, false, err
	}

	// Check if there's a next page
	hasNextPage := len(albums) > first
	if hasNextPage {
		albums = albums[:first]
	}

	return albums, hasNextPage, nil
}

// getAlbumsAfterCursor gets albums using timestamp-based cursor
func (p *PaginationHelper) getAlbumsAfterCursor(ctx context.Context, cursor *CursorInfo, limit int) ([]*models.Album, error) {
	// This would ideally be a custom repository method for efficient timestamp-based pagination
	// For now, we'll use a combination of approaches

	// Get all albums and filter (this is not optimal for large datasets)
	// In production, you'd want a custom SQL query with WHERE created_at < ? ORDER BY created_at DESC
	allAlbums, err := p.repos.Album.List(ctx, 1000, 0) // Get a reasonable batch
	if err != nil {
		return nil, err
	}

	var filteredAlbums []*models.Album
	found := false

	for _, album := range allAlbums {
		if album.ID.String() == cursor.ID {
			found = true
			continue
		}

		if found {
			filteredAlbums = append(filteredAlbums, album)
			if len(filteredAlbums) >= limit {
				break
			}
		}
	}

	return filteredAlbums, nil
}

// getAlbumsAfterPosition gets albums using position-based fallback
func (p *PaginationHelper) getAlbumsAfterPosition(ctx context.Context, cursorID string, limit int) ([]*models.Album, error) {
	// Find the position of the cursor ID
	albums, err := p.repos.Album.List(ctx, 1000, 0) // Get a reasonable batch
	if err != nil {
		return nil, err
	}

	position := -1
	for i, album := range albums {
		if album.ID.String() == cursorID {
			position = i
			break
		}
	}

	if position == -1 {
		// Cursor not found, return empty result
		return []*models.Album{}, nil
	}

	// Return albums after the cursor position
	startIdx := position + 1
	endIdx := startIdx + limit
	if endIdx > len(albums) {
		endIdx = len(albums)
	}

	if startIdx >= len(albums) {
		return []*models.Album{}, nil
	}

	return albums[startIdx:endIdx], nil
}

// Similar methods for other entity types...

// GetTracksWithCursor fetches tracks using proper cursor-based pagination
func (p *PaginationHelper) GetTracksWithCursor(ctx context.Context, first int, after *string) ([]*models.Track, bool, error) {
	var tracks []*models.Track
	var err error

	if after != nil && *after != "" {
		cursor, err := p.DecodeCursor(*after)
		if err != nil {
			return nil, false, err
		}

		if !cursor.CreatedAt.IsZero() {
			tracks, err = p.getTracksAfterCursor(ctx, cursor, first+1)
			if err != nil {
				return nil, false, err
			}
		} else {
			tracks, err = p.getTracksAfterPosition(ctx, cursor.ID, first+1)
			if err != nil {
				return nil, false, err
			}
		}
	} else {
		tracks, err = p.repos.Track.List(ctx, first+1, 0)
	}

	if err != nil {
		return nil, false, err
	}

	hasNextPage := len(tracks) > first
	if hasNextPage {
		tracks = tracks[:first]
	}

	return tracks, hasNextPage, nil
}

// getTracksAfterCursor gets tracks using timestamp-based cursor
func (p *PaginationHelper) getTracksAfterCursor(ctx context.Context, cursor *CursorInfo, limit int) ([]*models.Track, error) {
	allTracks, err := p.repos.Track.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}

	var filteredTracks []*models.Track
	found := false

	for _, track := range allTracks {
		if track.ID.String() == cursor.ID {
			found = true
			continue
		}

		if found {
			filteredTracks = append(filteredTracks, track)
			if len(filteredTracks) >= limit {
				break
			}
		}
	}

	return filteredTracks, nil
}

// getTracksAfterPosition gets tracks using position-based fallback
func (p *PaginationHelper) getTracksAfterPosition(ctx context.Context, cursorID string, limit int) ([]*models.Track, error) {
	tracks, err := p.repos.Track.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}

	position := -1
	for i, track := range tracks {
		if track.ID.String() == cursorID {
			position = i
			break
		}
	}

	if position == -1 {
		return []*models.Track{}, nil
	}

	startIdx := position + 1
	endIdx := startIdx + limit
	if endIdx > len(tracks) {
		endIdx = len(tracks)
	}

	if startIdx >= len(tracks) {
		return []*models.Track{}, nil
	}

	return tracks[startIdx:endIdx], nil
}

// GetReviewsWithCursor fetches reviews using proper cursor-based pagination
func (p *PaginationHelper) GetReviewsWithCursor(ctx context.Context, first int, after *string) ([]*models.Review, bool, error) {
	var reviews []*models.Review
	var err error

	if after != nil && *after != "" {
		cursor, err := p.DecodeCursor(*after)
		if err != nil {
			return nil, false, err
		}

		if !cursor.CreatedAt.IsZero() {
			reviews, err = p.getReviewsAfterCursor(ctx, cursor, first+1)
			if err != nil {
				return nil, false, err
			}
		} else {
			reviews, err = p.getReviewsAfterPosition(ctx, cursor.ID, first+1)
			if err != nil {
				return nil, false, err
			}
		}
	} else {
		reviews, err = p.repos.Review.List(ctx, first+1, 0)
	}

	if err != nil {
		return nil, false, err
	}

	hasNextPage := len(reviews) > first
	if hasNextPage {
		reviews = reviews[:first]
	}

	return reviews, hasNextPage, nil
}

// getReviewsAfterCursor gets reviews using timestamp-based cursor
func (p *PaginationHelper) getReviewsAfterCursor(ctx context.Context, cursor *CursorInfo, limit int) ([]*models.Review, error) {
	allReviews, err := p.repos.Review.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}

	var filteredReviews []*models.Review
	found := false

	for _, review := range allReviews {
		if review.ID.String() == cursor.ID {
			found = true
			continue
		}

		if found {
			filteredReviews = append(filteredReviews, review)
			if len(filteredReviews) >= limit {
				break
			}
		}
	}

	return filteredReviews, nil
}

// getReviewsAfterPosition gets reviews using position-based fallback
func (p *PaginationHelper) getReviewsAfterPosition(ctx context.Context, cursorID string, limit int) ([]*models.Review, error) {
	reviews, err := p.repos.Review.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}

	position := -1
	for i, review := range reviews {
		if review.ID.String() == cursorID {
			position = i
			break
		}
	}

	if position == -1 {
		return []*models.Review{}, nil
	}

	startIdx := position + 1
	endIdx := startIdx + limit
	if endIdx > len(reviews) {
		endIdx = len(reviews)
	}

	if startIdx >= len(reviews) {
		return []*models.Review{}, nil
	}

	return reviews[startIdx:endIdx], nil
}

// GetPlaylistsWithCursor fetches playlists using proper cursor-based pagination
func (p *PaginationHelper) GetPlaylistsWithCursor(ctx context.Context, first int, after *string) ([]*models.Playlist, bool, error) {
	var playlists []*models.Playlist
	var err error

	if after != nil && *after != "" {
		cursor, err := p.DecodeCursor(*after)
		if err != nil {
			return nil, false, err
		}

		if !cursor.CreatedAt.IsZero() {
			playlists, err = p.getPlaylistsAfterCursor(ctx, cursor, first+1)
			if err != nil {
				return nil, false, err
			}
		} else {
			playlists, err = p.getPlaylistsAfterPosition(ctx, cursor.ID, first+1)
			if err != nil {
				return nil, false, err
			}
		}
	} else {
		playlists, err = p.repos.Playlist.List(ctx, first+1, 0)
	}

	if err != nil {
		return nil, false, err
	}

	hasNextPage := len(playlists) > first
	if hasNextPage {
		playlists = playlists[:first]
	}

	return playlists, hasNextPage, nil
}

// getPlaylistsAfterCursor gets playlists using timestamp-based cursor
func (p *PaginationHelper) getPlaylistsAfterCursor(ctx context.Context, cursor *CursorInfo, limit int) ([]*models.Playlist, error) {
	allPlaylists, err := p.repos.Playlist.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}

	var filteredPlaylists []*models.Playlist
	found := false

	for _, playlist := range allPlaylists {
		if playlist.ID.String() == cursor.ID {
			found = true
			continue
		}

		if found {
			filteredPlaylists = append(filteredPlaylists, playlist)
			if len(filteredPlaylists) >= limit {
				break
			}
		}
	}

	return filteredPlaylists, nil
}

// getPlaylistsAfterPosition gets playlists using position-based fallback
func (p *PaginationHelper) getPlaylistsAfterPosition(ctx context.Context, cursorID string, limit int) ([]*models.Playlist, error) {
	playlists, err := p.repos.Playlist.List(ctx, 1000, 0)
	if err != nil {
		return nil, err
	}

	position := -1
	for i, playlist := range playlists {
		if playlist.ID.String() == cursorID {
			position = i
			break
		}
	}

	if position == -1 {
		return []*models.Playlist{}, nil
	}

	startIdx := position + 1
	endIdx := startIdx + limit
	if endIdx > len(playlists) {
		endIdx = len(playlists)
	}

	if startIdx >= len(playlists) {
		return []*models.Playlist{}, nil
	}

	return playlists[startIdx:endIdx], nil
}
