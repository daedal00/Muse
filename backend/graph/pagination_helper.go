package graph

import (
	"context"
	"encoding/base64"
	"encoding/json"
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
	cursor := CursorInfo{
		ID:        id,
		CreatedAt: createdAt,
		Position:  position,
	}

	data, _ := json.Marshal(cursor)
	return base64.StdEncoding.EncodeToString(data)
}

// DecodeCursor decodes a cursor to extract ID, timestamp, and position
func (p *PaginationHelper) DecodeCursor(cursor string) (*CursorInfo, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor: %w", err)
	}

	var cursorInfo CursorInfo
	if err := json.Unmarshal(data, &cursorInfo); err != nil {
		return nil, fmt.Errorf("invalid cursor format: %w", err)
	}

	return &cursorInfo, nil
}

// CreateCursor creates a cursor from an ID and created timestamp
func (p *PaginationHelper) CreateCursor(id string, createdAt time.Time, position int) (string, error) {
	if id == "" {
		return "", fmt.Errorf("id cannot be empty")
	}

	cursor := CursorInfo{
		ID:        id,
		CreatedAt: createdAt,
		Position:  position,
	}

	data, err := json.Marshal(cursor)
	if err != nil {
		return "", fmt.Errorf("failed to encode cursor: %w", err)
	}

	return base64.StdEncoding.EncodeToString(data), nil
}

// ParseCursor parses a cursor string and returns a CursorInfo
func (p *PaginationHelper) ParseCursor(cursor string) (*CursorInfo, error) {
	data, err := base64.StdEncoding.DecodeString(cursor)
	if err != nil {
		return nil, fmt.Errorf("invalid cursor encoding: %w", err)
	}

	var info CursorInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return nil, fmt.Errorf("invalid cursor format: %w", err)
	}

	if info.ID == "" {
		return nil, fmt.Errorf("cursor missing required ID field")
	}

	return &CursorInfo{
		ID:        info.ID,
		CreatedAt: info.CreatedAt,
		Position:  info.Position,
	}, nil
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
