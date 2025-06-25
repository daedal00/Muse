package postgres

import (
	"context"
	"fmt"

	"github.com/daedal00/muse/backend/internal/database"
	"github.com/daedal00/muse/backend/internal/models"
	"github.com/daedal00/muse/backend/internal/repository"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

type playlistRepository struct {
	db *database.PostgresDB
}

func NewPlaylistRepository(db *database.PostgresDB) repository.PlaylistRepository {
	return &playlistRepository{db: db}
}

func (r *playlistRepository) Create(ctx context.Context, playlist *models.Playlist) error {
	query := `
		INSERT INTO playlists (id, title, description, cover_image, creator_id, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	
	_, err := r.db.Pool.Exec(ctx, query,
		playlist.ID, playlist.Title, playlist.Description, playlist.CoverImage,
		playlist.CreatorID, playlist.CreatedAt, playlist.UpdatedAt,
	)
	
	if err != nil {
		return fmt.Errorf("failed to create playlist: %w", err)
	}
	
	return nil
}

func (r *playlistRepository) GetByID(ctx context.Context, id uuid.UUID) (*models.Playlist, error) {
	query := `
		SELECT id, title, description, cover_image, creator_id, created_at, updated_at
		FROM playlists 
		WHERE id = $1
	`
	
	playlist := &models.Playlist{}
	err := r.db.Pool.QueryRow(ctx, query, id).Scan(
		&playlist.ID, &playlist.Title, &playlist.Description, &playlist.CoverImage,
		&playlist.CreatorID, &playlist.CreatedAt, &playlist.UpdatedAt,
	)
	
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("playlist not found")
		}
		return nil, fmt.Errorf("failed to get playlist: %w", err)
	}
	
	return playlist, nil
}

func (r *playlistRepository) GetByCreatorID(ctx context.Context, creatorID uuid.UUID, limit, offset int) ([]*models.Playlist, error) {
	query := `
		SELECT id, title, description, cover_image, creator_id, created_at, updated_at
		FROM playlists 
		WHERE creator_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Pool.Query(ctx, query, creatorID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list playlists by creator: %w", err)
	}
	defer rows.Close()
	
	var playlists []*models.Playlist
	for rows.Next() {
		playlist := &models.Playlist{}
		err := rows.Scan(
			&playlist.ID, &playlist.Title, &playlist.Description, &playlist.CoverImage,
			&playlist.CreatorID, &playlist.CreatedAt, &playlist.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playlist: %w", err)
		}
		playlists = append(playlists, playlist)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating playlists: %w", err)
	}
	
	return playlists, nil
}

func (r *playlistRepository) Update(ctx context.Context, playlist *models.Playlist) error {
	query := `
		UPDATE playlists 
		SET title = $2, description = $3, cover_image = $4, updated_at = NOW()
		WHERE id = $1
	`
	
	result, err := r.db.Pool.Exec(ctx, query,
		playlist.ID, playlist.Title, playlist.Description, playlist.CoverImage,
	)
	
	if err != nil {
		return fmt.Errorf("failed to update playlist: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("playlist not found")
	}
	
	return nil
}

func (r *playlistRepository) Delete(ctx context.Context, id uuid.UUID) error {
	// Start a transaction to delete playlist tracks first, then the playlist
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// Delete all playlist tracks first
	_, err = tx.Exec(ctx, `DELETE FROM playlist_tracks WHERE playlist_id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete playlist tracks: %w", err)
	}
	
	// Delete the playlist
	result, err := tx.Exec(ctx, `DELETE FROM playlists WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete playlist: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("playlist not found")
	}
	
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

func (r *playlistRepository) List(ctx context.Context, limit, offset int) ([]*models.Playlist, error) {
	query := `
		SELECT id, title, description, cover_image, creator_id, created_at, updated_at
		FROM playlists 
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`
	
	rows, err := r.db.Pool.Query(ctx, query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list playlists: %w", err)
	}
	defer rows.Close()
	
	var playlists []*models.Playlist
	for rows.Next() {
		playlist := &models.Playlist{}
		err := rows.Scan(
			&playlist.ID, &playlist.Title, &playlist.Description, &playlist.CoverImage,
			&playlist.CreatorID, &playlist.CreatedAt, &playlist.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan playlist: %w", err)
		}
		playlists = append(playlists, playlist)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating playlists: %w", err)
	}
	
	return playlists, nil
}

// Playlist track operations

func (r *playlistRepository) AddTrack(ctx context.Context, playlistID, trackID uuid.UUID, position int) error {
	// If position is 0 or negative, append to the end
	if position <= 0 {
		var maxPosition int
		query := `SELECT COALESCE(MAX(position), 0) FROM playlist_tracks WHERE playlist_id = $1`
		err := r.db.Pool.QueryRow(ctx, query, playlistID).Scan(&maxPosition)
		if err != nil {
			return fmt.Errorf("failed to get max position: %w", err)
		}
		position = maxPosition + 1
	} else {
		// Shift existing tracks to make room
		shiftQuery := `UPDATE playlist_tracks SET position = position + 1 WHERE playlist_id = $1 AND position >= $2`
		_, err := r.db.Pool.Exec(ctx, shiftQuery, playlistID, position)
		if err != nil {
			return fmt.Errorf("failed to shift track positions: %w", err)
		}
	}
	
	insertQuery := `
		INSERT INTO playlist_tracks (id, playlist_id, track_id, position, added_at)
		VALUES ($1, $2, $3, $4, NOW())
	`
	
	_, err := r.db.Pool.Exec(ctx, insertQuery, uuid.New(), playlistID, trackID, position)
	if err != nil {
		return fmt.Errorf("failed to add track to playlist: %w", err)
	}
	
	return nil
}

func (r *playlistRepository) RemoveTrack(ctx context.Context, playlistID, trackID uuid.UUID) error {
	// Get the position of the track being removed
	var position int
	getPositionQuery := `SELECT position FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`
	err := r.db.Pool.QueryRow(ctx, getPositionQuery, playlistID, trackID).Scan(&position)
	if err != nil {
		if err == pgx.ErrNoRows {
			return fmt.Errorf("track not found in playlist")
		}
		return fmt.Errorf("failed to get track position: %w", err)
	}
	
	// Start a transaction to remove track and shift positions
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	// Remove the track
	removeQuery := `DELETE FROM playlist_tracks WHERE playlist_id = $1 AND track_id = $2`
	result, err := tx.Exec(ctx, removeQuery, playlistID, trackID)
	if err != nil {
		return fmt.Errorf("failed to remove track from playlist: %w", err)
	}
	
	if result.RowsAffected() == 0 {
		return fmt.Errorf("track not found in playlist")
	}
	
	// Shift remaining tracks down
	shiftQuery := `UPDATE playlist_tracks SET position = position - 1 WHERE playlist_id = $1 AND position > $2`
	_, err = tx.Exec(ctx, shiftQuery, playlistID, position)
	if err != nil {
		return fmt.Errorf("failed to shift track positions: %w", err)
	}
	
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
}

func (r *playlistRepository) GetTracks(ctx context.Context, playlistID uuid.UUID, limit, offset int) ([]*models.Track, error) {
	query := `
		SELECT t.id, t.spotify_id, t.title, t.album_id, t.duration_ms, t.track_number, t.created_at, t.updated_at
		FROM tracks t
		INNER JOIN playlist_tracks pt ON t.id = pt.track_id
		WHERE pt.playlist_id = $1
		ORDER BY pt.position ASC
		LIMIT $2 OFFSET $3
	`
	
	rows, err := r.db.Pool.Query(ctx, query, playlistID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get playlist tracks: %w", err)
	}
	defer rows.Close()
	
	var tracks []*models.Track
	for rows.Next() {
		track := &models.Track{}
		err := rows.Scan(
			&track.ID, &track.SpotifyID, &track.Title, &track.AlbumID,
			&track.DurationMs, &track.TrackNumber, &track.CreatedAt, &track.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan track: %w", err)
		}
		tracks = append(tracks, track)
	}
	
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error iterating tracks: %w", err)
	}
	
	return tracks, nil
}

func (r *playlistRepository) ReorderTracks(ctx context.Context, playlistID uuid.UUID, trackPositions map[uuid.UUID]int) error {
	tx, err := r.db.Pool.Begin(ctx)
	if err != nil {
		return fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	
	for trackID, position := range trackPositions {
		updateQuery := `UPDATE playlist_tracks SET position = $1 WHERE playlist_id = $2 AND track_id = $3`
		_, err := tx.Exec(ctx, updateQuery, position, playlistID, trackID)
		if err != nil {
			return fmt.Errorf("failed to update track position: %w", err)
		}
	}
	
	if err := tx.Commit(ctx); err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}
	
	return nil
} 